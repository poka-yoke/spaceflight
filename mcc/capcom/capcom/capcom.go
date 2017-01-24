package capcom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/awalterschulze/gographviz"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type sGInstanceState map[string]map[string]int

func (s sGInstanceState) getKeys() []string {
	i := 0
	keys := make([]string, len(s))
	for k := range s {
		keys[i] = k
		i++
	}
	return keys
}

func (s sGInstanceState) has(key string) bool {
	for _, k := range s.getKeys() {
		if k == key {
			return true
		}
	}
	return false
}

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
func AuthorizeIPToSecurityGroup(svc *ec2.EC2, ipRange string, port int64, sgid string) {
	ran := &ec2.IpRange{
		CidrIp: aws.String(ipRange),
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
func RevokeIPToSecurityGroup(svc *ec2.EC2, ipRange string, port int64, sgid string) {
	ran := &ec2.IpRange{
		CidrIp: aws.String(ipRange),
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

func nodeAttrs(sg *ec2.SecurityGroup) (attrs gographviz.Attrs) {
	attrs = gographviz.NewAttrs()
	attrs.Add("label", fmt.Sprintf("{{%s|}|%s}", *sg.GroupId, *sg.GroupName))
	return
}

func registerNodes(
	sglist []*ec2.SecurityGroup,
	graph *gographviz.Escape,
	nodesPresence sGInstanceState,
) {
	for _, sg := range sglist {
		log.Printf(
			"Adding node for %s (%s)\n",
			*sg.GroupName,
			*sg.GroupId,
		)
		attrs := nodeAttrs(sg)
		switch {
		case nodesPresence[*sg.GroupId]["running"] > 0:
			attrs.Add("color", "green")
		case nodesPresence[*sg.GroupId]["running"] == 0 &&
			nodesPresence[*sg.GroupId]["stopped"] > 0:
			attrs.Add("color", "yellow")
		case nodesPresence[*sg.GroupId]["running"] == 0 &&
			nodesPresence[*sg.GroupId]["stopped"] == 0:
			attrs.Add("color", "red")
		}
		graph.AddNode("G", *sg.GroupId, attrs)
		if nodesPresence[*sg.GroupId] == nil {
			nodesPresence[*sg.GroupId] = nil
		}
	}
}

func edgeAttrs(perm *ec2.IpPermission) (attrs gographviz.Attrs) {
	var val string
	if perm.FromPort != nil && perm.ToPort != nil {
		fromport := strconv.FormatInt(*perm.FromPort, 10)
		toport := strconv.FormatInt(*perm.ToPort, 10)
		if *perm.FromPort == *perm.ToPort {
			val = fmt.Sprintf(
				"%s: %s",
				*perm.IpProtocol,
				fromport,
			)
		} else {
			val = fmt.Sprintf(
				"%s: %s - %s",
				*perm.IpProtocol,
				fromport,
				toport,
			)
		}
		attrs = gographviz.NewAttrs()
		attrs.Add("label", val)
	}
	return attrs
}

func registerEdges(
	sglist []*ec2.SecurityGroup,
	graph *gographviz.Escape,
	nodesPresence sGInstanceState,
) {
	for _, sg := range sglist {
		log.Printf(
			"Processing entries for %s (%s)\n",
			*sg.GroupName,
			*sg.GroupId,
		)
		for _, perm := range sg.IpPermissions {
			for _, pair := range perm.UserIdGroupPairs {
				if nodesPresence.has(*pair.GroupId) {
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
						edgeAttrs(perm),
					)
				}
			}
		}
	}
}

func getInstancesPerSG(svc *ec2.EC2) sGInstanceState {
	iState := make(sGInstanceState)
	// TODO: Check for need of pagination and handle it
	resp, err := svc.DescribeInstances(
		&ec2.DescribeInstancesInput{
			MaxResults: aws.Int64(1000),
		},
	)
	if err != nil {
		log.Panic(err.Error())
	}
	for _, res := range resp.Reservations {
		groupID := []string{}
		state := map[string]int{
			"pending":       0,
			"running":       0,
			"shutting-down": 0,
			"terminated":    0,
			"stopping":      0,
			"stopped":       0,
		}
		for _, group := range res.Groups {
			groupID = append(groupID, *group.GroupId)
		}
		for _, instance := range res.Instances {
			state[*instance.State.Name]++
		}
		for _, gid := range groupID {
			iState[gid] = state
		}
	}
	return iState
}

// GraphSGRelations returns a string containing a graph representation in DOT
// format of the relations between Security Groups in the service.
func GraphSGRelations(svc *ec2.EC2) string {
	sglist := getSecurityGroups(svc).SecurityGroups
	nodesPresence := getInstancesPerSG(svc)

	g := gographviz.NewEscape()
	g.SetName("G")
	g.SetDir(true)
	log.Println("Created graph")

	registerNodes(sglist, g, nodesPresence)
	registerEdges(sglist, g, nodesPresence)
	return g.String()
}

// Init initializes connection to AWS API
func Init() *ec2.EC2 {
	region := "us-east-1"
	sess := session.New(&aws.Config{Region: aws.String(region)})
	return ec2.New(sess)
}
