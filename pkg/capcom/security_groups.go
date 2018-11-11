package capcom

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
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
		err = fmt.Errorf("%s is not a valid CIDR", cidr)
		return nil, err
	}
	// Obtain and traverse AWS's Security Group structure
	for _, sg := range getSecurityGroups(svc).SecurityGroups {
		for _, perm := range sg.IpPermissions {
			for _, ipRange := range perm.IpRanges {
				cont, err := networkContainsIPCheck(
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
					out = append(out,
						fmt.Sprintf(
							"%s %s/%s %s",
							*sg.GroupId,
							strconv.FormatInt(*perm.ToPort, 10),
							*perm.IpProtocol,
							*ipRange.CidrIp,
						),
					)
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
