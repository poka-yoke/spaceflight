package capcom

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Permission represents a Permission for a Security Group
type Permission struct {
	sgid, cidr, protocol string
	port int64
}

// NewPermission returns a pointer to a new Permission object
func NewPermission(origin, protocol string, port int64) (*Permission, error) {
	perm := Permission{protocol: protocol, port:port}
	switch {
	case strings.HasPrefix(origin, "sg-"):
		// It's a security group
		perm.sgid = origin

	case isCIDR(origin):
		// It's a valid CIDR
		perm.cidr = origin
	default:
		err := fmt.Errorf(
			"%s is neither sgid nor IP range in CIDR notation",
			origin,
		)
		return nil, err
	}
	return &perm, nil
}

// buildIPPermission provides an IpPermission object fully populated
func (p *Permission) buildIPPermission() (
	perm *ec2.IpPermission,
) {
	perm = &ec2.IpPermission{
		FromPort:   &p.port,
		ToPort:     &p.port,
		IpProtocol: &p.protocol,
	}
	if p.sgid != "" {
		perm.UserIdGroupPairs = []*ec2.UserIdGroupPair{{GroupId: &p.sgid}}
	}
	if p.cidr != "" {
		perm.IpRanges = []*ec2.IpRange{{CidrIp: &p.cidr}}
	}
	return perm
}

// AddToSG adds the permission to the specified Security Group ID
// using the service
func (p *Permission) AddToSG(svc ec2iface.EC2API, sgid string) bool {
	return authorizeAccessToSecurityGroup(svc, p, sgid)
}

// RemoveToSG removes the permission to the specified Security Group
// ID using the service
func (p *Permission) RemoveToSG(svc ec2iface.EC2API, sgid string) bool {
	return revokeAccessToSecurityGroup(svc, p, sgid)
}

// CreateSG creates a new security group. If a vpcid is specified the security
// group will be in that VPC
func CreateSG(
	name string,
	description string,
	vpcid string,
	svc ec2iface.EC2API,
) string {
	if description == "" {
		log.Fatal("Not a valid description")
	}
	params := &ec2.CreateSecurityGroupInput{
		Description: aws.String(description),
		GroupName:   aws.String(name),
	}
	if vpcid != "" {
		params.VpcId = aws.String(vpcid)
	}
	if err := params.Validate(); err != nil {
		log.Panic(err.Error())
	}
	res, err := svc.CreateSecurityGroup(params)
	if err != nil {
		log.Panic(err.Error())
	}
	return *res.GroupId
}

// AuthorizeAccessToSecurityGroup adds the specified permissions to the Ingress
// list of the destination security group on protocol and port
func authorizeAccessToSecurityGroup(
	svc ec2iface.EC2API,
	perm *Permission,
	destination string,
) bool {
	out, error := svc.AuthorizeSecurityGroupIngress(
		&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:       &destination,
			IpPermissions: []*ec2.IpPermission{perm.buildIPPermission()},
		})
	if error != nil {
		log.Panic(error)
	}
	if !(out != nil) {
		return false
	}
	return true
}

// RevokeAccessToSecurityGroup adds the specified permissions to the Ingress
// list of the destination security group on protocol and port
func revokeAccessToSecurityGroup(
	svc ec2iface.EC2API,
	perm *Permission,
	destination string,
) bool {
	out, error := svc.RevokeSecurityGroupIngress(
		&ec2.RevokeSecurityGroupIngressInput{
			GroupId:       &destination,
			IpPermissions: []*ec2.IpPermission{perm.buildIPPermission()},
		})
	if error != nil {
		log.Panic(error)
	}
	if !(out != nil) {
		return false
	}
	return true
}

// FindSecurityGroupsWithRange returns a list of SGIDs where the CIDR
// passed in matches any of the rules
func FindSecurityGroupsWithRange(
	svc ec2iface.EC2API,
	cidr string,
) (
	out []SearchResult,
	err error,
) {
	// IP we are searching for in the Security Groups
	searchIP, _, err := net.ParseCIDR(cidr)
	if err != nil {
		err = fmt.Errorf("%s is not a valid CIDR\n", cidr)
		return
	}
	// Obtain and traverse AWS's Security Group structure
	for _, sg := range getSecurityGroups(svc).SecurityGroups {
		for _, perm := range sg.IpPermissions {
			for _, ipRange := range perm.IpRanges {
				cont, err := NetworkContainsIPCheck(
					*ipRange.CidrIp,
					searchIP,
				)
				if err != nil {
					log.Printf(
						"Invalid CIDR %s in SG %s (%s)\n",
						*ipRange.CidrIp,
						*sg.GroupName,
						*sg.GroupId,
					)
				}
				if cont {
					out = append(out, SearchResult{
						GroupID:  *sg.GroupId,
						Protocol: *perm.IpProtocol,
						Port:     *perm.ToPort,
						Source:   *ipRange.CidrIp,
					})
				}
			}
		}
	}
	return
}

// getSecurityGroups retrieves the list of all Security Groups in the account
func getSecurityGroups(svc ec2iface.EC2API) *ec2.DescribeSecurityGroupsOutput {
	res, err := svc.DescribeSecurityGroups(nil)
	if err != nil {
		log.Panic(err)
	}
	return res
}

// ListSecurityGroups prints all available Security groups accessible
// by the account on svc
func ListSecurityGroups(svc ec2iface.EC2API) (out []string) {
	for _, sg := range getSecurityGroups(svc).SecurityGroups {
		out = append(out, fmt.Sprintf("* %10s %20s %s\n",
			*sg.GroupId,
			*sg.GroupName,
			*sg.Description),
		)
	}
	return
}

// NetworkContainsIPCheck returns true if the subnet expresed in the
// CIDR in contains the IP object
func NetworkContainsIPCheck(cidr string, searchIP net.IP) (out bool, err error) {
	ip, sub, err := net.ParseCIDR(cidr)
	if err != nil {
		err = fmt.Errorf(
			"Failed parsing CIDR %s\n",
			cidr,
		)
		return
	}
	out = ip.Equal(searchIP) || sub.Contains(searchIP)
	return
}

func isCIDR(origin string) bool {
	_, _, err := net.ParseCIDR(origin)
	if err != nil {
		return false
	}
	return true
}

// FindSGByName gets an array of sgids for a name search
func FindSGByName(name string, vpc string, svc ec2iface.EC2API) (ret []string) {
	res, err := svc.DescribeSecurityGroups(
		&ec2.DescribeSecurityGroupsInput{
			Filters: []*ec2.Filter{
				{
					Name:   aws.String("group-name"),
					Values: []*string{&name},
				},
			},
		})
	if err != nil {
		log.Panic(err.Error())
	}
	for _, sg := range res.SecurityGroups {
		ret = append(ret, *sg.GroupId)
	}
	return ret
}
