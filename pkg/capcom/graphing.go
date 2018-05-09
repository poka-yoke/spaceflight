package capcom

import (
	"fmt"
	"log"
	"strconv"

	"github.com/awalterschulze/gographviz"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
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

func getInstances(svc ec2iface.EC2API) *ec2.DescribeInstancesOutput {
	// TODO: Check for need of pagination and handle it
	resp, err := svc.DescribeInstances(
		&ec2.DescribeInstancesInput{
			MaxResults: aws.Int64(1000),
		},
	)
	if err != nil {
		log.Panic(err.Error())
	}
	return resp
}

func getInstancesStates(instances []*ec2.Reservation) sGInstanceState {
	iState := make(sGInstanceState)
	for _, res := range instances {
		state := map[string]int{
			"pending":       0,
			"running":       0,
			"shutting-down": 0,
			"terminated":    0,
			"stopping":      0,
			"stopped":       0,
		}
		for _, instance := range res.Instances {
			state[*instance.State.Name]++
		}
		for _, group := range res.Groups {
			iState[*group.GroupId] = state
		}
	}
	return iState
}

func nodeAttrs(sg *ec2.SecurityGroup) (attrs map[string]string) {
	attrs = make(map[string]string)
	attrs["label"] = fmt.Sprintf("{{%s|}|%s}", *sg.GroupId, *sg.GroupName)
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
			attrs["color"] = "green"
		case nodesPresence[*sg.GroupId]["running"] == 0 &&
			nodesPresence[*sg.GroupId]["stopped"] > 0:
			attrs["color"] = "yellow"
		case nodesPresence[*sg.GroupId]["running"] == 0 &&
			nodesPresence[*sg.GroupId]["stopped"] == 0:
			attrs["color"] = "red"
		}
		if err := graph.AddNode("G", *sg.GroupId, attrs); err != nil {
			log.Println(err)
		}
		if nodesPresence[*sg.GroupId] == nil {
			nodesPresence[*sg.GroupId] = nil
		}
	}
}

func edgeAttrs(perm *ec2.IpPermission) (attrs map[string]string) {
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
		attrs = make(map[string]string)
		attrs["label"] = val
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
					if err := graph.AddEdge(
						*sg.GroupId,
						*pair.GroupId,
						true,
						edgeAttrs(perm),
					); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

// GraphSGRelations returns a string containing a graph representation in DOT
// format of the relations between Security Groups in the service.
func GraphSGRelations(svc ec2iface.EC2API) string {
	sglist := getSecurityGroups(svc).SecurityGroups

	g := gographviz.NewEscape()
	if err := g.SetName("G"); err != nil {
		log.Println(err)
	}
	if err := g.SetDir(true); err != nil {
		log.Println(err)
	}
	log.Println("Created graph")

	nodesPresence := getInstancesStates(getInstances(svc).Reservations)
	registerNodes(sglist, g, nodesPresence)
	registerEdges(sglist, g, nodesPresence)
	return g.String()
}
