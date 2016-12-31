package roosa

import (
	"testing"
)

func TestNewReferenceTreeList(t *testing.T) {
	// 	records []*route53.ResourceRecordSet,
	// ) *ReferenceTreeList {
	// 	return &ReferenceTreeList{
	// 		records: ttl.FilterResourceRecordSetType(
	// 			records,
	// 			recordTypes,
	// 		),
	// 		lookup: nil,
	// 	}
}

func TestReferenceTreeListString(t *testing.T) {
	// 	if rtl.lookup == nil {
	// 		rtl.GetReferenceTrees()
	// 	}
	// 	output = ""
	// 	for _, tree := range rtl.lookup {
	// 		for _, node := range tree {
	// 			output += fmt.Sprintf("%v\n", node)
	// 		}
	// 	}
	// 	return
}

func TestReferenceTreeListFill(t *testing.T) {
	// 	rtl.lookup = map[string][]*Node{}
	// 	for _, val := range rtl.records {
	// 		node := &Node{
	// 			content: val,
	// 		}
	// 		name := strings.TrimSuffix(*val.Name, ".")
	// 		rtl.lookup[name] = append(rtl.lookup[name], node)
	// 	}
	// 	log.Printf("Added %d records to Lookup\n", len(rtl.lookup))
	// 	return
}

func TestReferenceTreeListClean(t *testing.T) {
	// 	for name, nodes := range rtl.lookup {
	// 		for _, node := range nodes {
	// 			if !node.IsRoot() {
	// 				delete(rtl.lookup, name)
	// 			}
	// 		}
	// 	}
}

func TestReferenceTreeListCompact(t *testing.T) {
	// 	for name, nodes := range rtl.lookup {
	// 		for _, node := range nodes {
	// 			value := *node.content.ResourceRecords[0].Value
	// 			if *node.content.Type == "CNAME" {
	// 				if parents, ok := rtl.lookup[value]; ok {
	// 					for _, parent := range parents {
	// 						log.Printf(
	// 							"%v (%v) has Parent %v\n",
	// 							name,
	// 							value,
	// 							*parent.content.Name,
	// 						)
	// 						parent.children = append(
	// 							parent.children,
	// 							node,
	// 						)
	// 						node.parent = parent
	// 					}
	// 				} else {
	// 					log.Printf(
	// 						"%v (%v) is Out of domain\n",
	// 						name,
	// 						value,
	// 					)
	// 				}
	// 			} else {
	// 				log.Printf(
	// 					"%v (%v) is Root\n",
	// 					name,
	// 					value,
	// 				)
	// 			}
	// 		}
	// 	}
	// 	rtl.clean()
	// 	log.Printf(
	// 		"Cleared children, %d records left in Lookup\n",
	// 		len(rtl.lookup),
	// 	)
}
