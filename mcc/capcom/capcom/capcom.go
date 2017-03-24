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

// FindSecurityGroupsWithRange returns a list of SGIDs where the CIDR
// passed in matches any of the rules
func FindSecurityGroupsWithRange(
	svc ec2iface.EC2API,
	cidr string,
) (
	out []string,
	err error,
) {
	// IP we are searching for in the Security Groups
	searchIP, _, err := net.ParseCIDR(cidr)
	if err != nil {
		err = fmt.Errorf("%s is not a valid CIDR\n", cidr)
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
					out = append(out, *sg.GroupId)
				}
			}
		}
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
	perm = &ec2.IpPermission{}
	perm.FromPort = &port
	perm.ToPort = &port
	perm.IpProtocol = &proto
	if strings.HasPrefix(origin, "sg-") {
		// It's a security group
		perm.UserIdGroupPairs = []*ec2.UserIdGroupPair{{GroupId: &origin}}
	} else {
		_, _, err = net.ParseCIDR(origin)
		if err != nil {
			// It's not a valid CIDR, and it wasn't an SGID before
			err = fmt.Errorf(
				"%s is neither sgid nor IP range in CIDR notation",
				origin,
			)
			return
		}
		// It's a valid CIDR
		perm.IpRanges = []*ec2.IpRange{{CidrIp: &origin}}
	}
	return
}

// AuthorizeAccessToSecurityGroup adds the specified permissions to the Ingress
// list of the destination security group on protocol and port
func AuthorizeAccessToSecurityGroup(
	svc ec2iface.EC2API,
	perm *ec2.IpPermission,
	destination string,
) *ec2.AuthorizeSecurityGroupIngressOutput {
	params := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       &destination,
		IpPermissions: []*ec2.IpPermission{perm},
	}
	out, error := svc.AuthorizeSecurityGroupIngress(params)
	if error != nil {
		log.Panic(error)
	}
	return out
}

// RevokeAccessToSecurityGroup adds the specified permissions to the Ingress
// list of the destination security group on protocol and port
func RevokeAccessToSecurityGroup(
	svc ec2iface.EC2API,
	perm *ec2.IpPermission,
	destination string,
) *ec2.RevokeSecurityGroupIngressOutput {
	params := &ec2.RevokeSecurityGroupIngressInput{
		GroupId:       &destination,
		IpPermissions: []*ec2.IpPermission{perm},
	}
	out, error := svc.RevokeSecurityGroupIngress(params)
	if error != nil {
		log.Panic(error)
	}
	return out
}

// Init initializes connection to AWS API
func Init() ec2iface.EC2API {
	region := "us-east-1"
	sess := session.New(&aws.Config{Region: aws.String(region)})
	return ec2.New(sess)
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

// FindSGByName gets an array of sgids for a name search
func FindSGByName(name string, vpc string, svc ec2iface.EC2API) (ret []string) {
	params := &ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("group-name"),
				Values: []*string{&name},
			},
		},
	}
	res, err := svc.DescribeSecurityGroups(params)
	if err != nil {
		log.Panic(err.Error())
	}
	for _, sg := range res.SecurityGroups {
		ret = append(ret, *sg.GroupId)
	}
	return ret
}
