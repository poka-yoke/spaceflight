package capcom

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

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
func AuthorizeAccessToSecurityGroup(
	svc ec2iface.EC2API,
	perm *ec2.IpPermission,
	destination string,
) bool {
	out, error := svc.AuthorizeSecurityGroupIngress(
		&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId:       &destination,
			IpPermissions: []*ec2.IpPermission{perm},
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
func RevokeAccessToSecurityGroup(
	svc ec2iface.EC2API,
	perm *ec2.IpPermission,
	destination string,
) bool {
	out, error := svc.RevokeSecurityGroupIngress(
		&ec2.RevokeSecurityGroupIngressInput{
			GroupId:       &destination,
			IpPermissions: []*ec2.IpPermission{perm},
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

// BuildIPPermission provides an IpPermission object fully populated
func BuildIPPermission(
	origin string,
	proto string,
	port int64,
) (
	perm *ec2.IpPermission,
	err error,
) {
	perm = &ec2.IpPermission{
		FromPort:   &port,
		ToPort:     &port,
		IpProtocol: &proto,
	}
	switch {
	case strings.HasPrefix(origin, "sg-"):
		// It's a security group
		perm.UserIdGroupPairs = []*ec2.UserIdGroupPair{{GroupId: &origin}}

	case isCIDR(origin):
		// It's a valid CIDR
		perm.IpRanges = []*ec2.IpRange{{CidrIp: &origin}}
	default:
		err = fmt.Errorf(
			"%s is neither sgid nor IP range in CIDR notation",
			origin,
		)
	}
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
				&ec2.Filter{
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

// Init initializes connection to AWS API
func Init() ec2iface.EC2API {
	region := "us-east-1"
	return ec2.New(
		session.New(
			&aws.Config{
				Region: aws.String(region),
			},
		),
	)
}
