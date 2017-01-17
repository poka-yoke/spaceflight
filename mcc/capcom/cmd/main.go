package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"

	"github.com/awalterschulze/gographviz"

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

func getSecurityGroups(svc *ec2.EC2) *ec2.DescribeSecurityGroupsOutput {
	res, err := svc.DescribeSecurityGroups(nil)
	if err != nil {
		log.Panic(err)
	}
	return res
}

// ListSecurityGroups prints all available Security groups accessible by the
// account on svc
func ListSecurityGroups(svc *ec2.EC2) {
	for _, sg := range getSecurityGroups(svc).SecurityGroups {
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

func registerNodes(
	sglist []*ec2.SecurityGroup,
	graph *gographviz.Escape,
	nodesPresence map[string]bool,
) {
	for _, sg := range sglist {
		val := fmt.Sprintf(
			"{{%s|%s}}",
			*sg.GroupId,
			*sg.GroupName,
		)
		attrs := gographviz.NewAttrs()
		attrs.Add("label", val)
		log.Printf(
			"Adding node for %s (%s)\n",
			*sg.GroupName,
			*sg.GroupId,
		)
		graph.AddNode("G", *sg.GroupId, attrs)
		nodesPresence[*sg.GroupId] = true
	}
}

func registerEdges(
	sglist []*ec2.SecurityGroup,
	graph *gographviz.Escape,
	nodesPresence map[string]bool,
) {
	for _, sg := range sglist {
		log.Printf(
			"Processing entries for %s (%s)\n",
			*sg.GroupName,
			*sg.GroupId,
		)
		for _, perm := range sg.IpPermissions {
			for _, pair := range perm.UserIdGroupPairs {
				if nodesPresence[*pair.GroupId] {
					attrs := gographviz.NewAttrs()
					if perm.FromPort != nil &&
						perm.ToPort != nil {
						fromport := strconv.FormatInt(
							*perm.FromPort,
							10,
						)
						toport := strconv.FormatInt(
							*perm.ToPort,
							10,
						)
						if *perm.FromPort ==
							*perm.ToPort {
							val := fmt.Sprintf(
								"%s: %s",
								*perm.IpProtocol,
								fromport,
							)
							attrs.Add("label", val)
						} else {
							val := fmt.Sprintf(
								"%s: %s - %s",
								*perm.IpProtocol,
								fromport,
								toport,
							)
							attrs.Add("label", val)
						}
					}
					groupName := ""
					if pair.GroupName != nil {
						groupName = *pair.GroupName
					}
					log.Printf(
						"Adding Edge for %s (%s) to %s (%s)\n",
						*sg.GroupName,
						*sg.GroupId,
						groupName,
						*pair.GroupId,
					)
					graph.AddEdge(
						*sg.GroupId,
						*pair.GroupId,
						true,
						attrs,
					)
				}
			}
		}
	}
}

// GraphSGRelations returns a string containing a graph representation in DOT
// format of the relations between Security Groups in the service.
func GraphSGRelations(svc *ec2.EC2) string {
	nodesPresence := make(map[string]bool)
	sglist := getSecurityGroups(svc).SecurityGroups

	g := gographviz.NewEscape()
	g.SetName("G")
	g.SetDir(true)
	log.Println("Created graph")

	registerNodes(sglist, g, nodesPresence)
	registerEdges(sglist, g, nodesPresence)
	return g.String()
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
