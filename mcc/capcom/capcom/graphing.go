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

func getInstancesPerSG(svc ec2iface.EC2API) sGInstanceState {
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
		graph.AddNode("G", *sg.GroupId, attrs)
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

// GraphSGRelations returns a string containing a graph representation in DOT
// format of the relations between Security Groups in the service.
func GraphSGRelations(svc ec2iface.EC2API) string {
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
