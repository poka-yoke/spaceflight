package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var list bool
var add bool
var revoke bool
var iprange string
var port int64
var sgid string

// ListSecurityGroups prints all available Security groups accessible by the
// account on svc
func ListSecurityGroups(svc *ec2.EC2) {
	res, err := svc.DescribeSecurityGroups(nil)
	if err != nil {
		log.Panic(err)
	}

	for _, sg := range res.SecurityGroups {
		fmt.Printf("* %10s %20s %s\n",
			*sg.GroupId,
			*sg.GroupName,
			*sg.Description)
	}
}

// AuthorizeIPToSecurityGroup adds the IP to the Ingress list of the target
// security group at the specified port
func AuthorizeIPToSecurityGroup(svc *ec2.EC2) {
	ran := &ec2.IpRange{
		CidrIp: aws.String(iprange),
	}
	perm := &ec2.IpPermission{
		FromPort:   &port,
		ToPort:     &port,
		IpProtocol: aws.String("tcp"),
		IpRanges:   []*ec2.IpRange{ran},
	}
	params := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupId:       aws.String(sgid),
		IpPermissions: []*ec2.IpPermission{perm},
	}
	_, err := svc.AuthorizeSecurityGroupIngress(params)
	if err != nil {
		log.Panic(err)
	}
}

// RevokeIPToSecurityGroup removes the IP from the Ingress list of the target
// security group at the specified port
func RevokeIPToSecurityGroup(svc *ec2.EC2) {
	ran := &ec2.IpRange{
		CidrIp: aws.String(iprange),
	}
	perm := &ec2.IpPermission{
		FromPort:   &port,
		ToPort:     &port,
		IpProtocol: aws.String("tcp"),
		IpRanges:   []*ec2.IpRange{ran},
	}
	params := &ec2.RevokeSecurityGroupIngressInput{
		GroupId:       aws.String(sgid),
		IpPermissions: []*ec2.IpPermission{perm},
	}
	_, err := svc.RevokeSecurityGroupIngress(params)
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	flag.BoolVar(&list, "l", false,
		"List all Security groups with ID, name and description")
	flag.BoolVar(&add, "a", false, "Add a rule to a security group")
	flag.BoolVar(&revoke, "r", false, "Revoke a rule to a security group")
	flag.StringVar(&iprange, "ip", "127.0.0.1/32", "IPv4 CIDR: 127.0.0.1/32")
	flag.Int64Var(&port, "p", 22, "Port for connections (default: 22)")
	flag.StringVar(&sgid, "sgid", "", "Security Group ID, sg-XXXXXXX")
	flag.Parse()

	region := "us-east-1"
	sess := session.New(&aws.Config{Region: aws.String(region)})
	svc := ec2.New(sess)

	if list {
		ListSecurityGroups(svc)
	}
	if add {
		AuthorizeIPToSecurityGroup(svc)
	}
	if revoke {
		RevokeIPToSecurityGroup(svc)
	}
}
