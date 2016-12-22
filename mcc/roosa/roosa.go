package roosa

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/poka-yoke/spaceflight/mcc/ttl"
)

// Node type represents the reference between data.
type Node struct {
	parent   *Node
	children []*Node
	content  *route53.ResourceRecordSet
}

func (n *Node) isRoot() bool {
	return n.parent == nil
}

// PrintNode prints a Node, and all its children.
func PrintNode(node *Node, indent int) {
	indents := ""
	for i := 0; i < indent; i++ {
		indents += "\t"
	}
	extra := ""
	l := len(node.content.ResourceRecords)
	for _, record := range node.content.ResourceRecords[:l-1] {
		extra += fmt.Sprintf("%v, ", *record.Value)
	}
	extra += *node.content.ResourceRecords[l-1].Value
	fmt.Printf("%v%v %v %v\n", indents, *node.content.Name, *node.content.Type, extra)
	indent++
	for _, child := range node.children {
		PrintNode(child, indent)
	}
}

// FilterReferrals returns 'A', 'AAAA', and 'CNAME' records from a []*route53.ResourceRecordSet.
func FilterReferrals(records []*route53.ResourceRecordSet) []*route53.ResourceRecordSet {
	return ttl.FilterResourceRecordSetType(
		records,
		[]string{
			"A",
			"AAAA",
			"CNAME",
		},
	)
}

// FillReferenceTrees fills a referral lookup table with the base records from a []*route53.ResourceRecordSet.
func FillReferenceTrees(records []*route53.ResourceRecordSet) (rootLookup map[string][]*Node) {
	rootLookup = map[string][]*Node{}
	for _, val := range records {
		node := &Node{
			parent:  nil,
			content: val,
		}
		name := strings.TrimSuffix(*val.Name, ".")
		rootLookup[name] = append(rootLookup[name], node)
	}
	log.Printf("Added %d records to Lookup\n", len(rootLookup))
	return
}

// CleanReferenceTrees removes non-root elements from a referral lookup table.
func CleanReferenceTrees(rootLookup map[string][]*Node) {
	for name, nodes := range rootLookup {
		for _, node := range nodes {
			if !node.isRoot() {
				delete(rootLookup, name)
			}
		}
	}
}

// CompactReferenceTrees modifies a referral lookup table finding children and roots, and relating these appropriately.
func CompactReferenceTrees(rootLookup map[string][]*Node) {
	for name, nodes := range rootLookup {
		for _, node := range nodes {
			if *node.content.Type == "CNAME" {
				if parents, ok := rootLookup[*node.content.ResourceRecords[0].Value]; ok {
					for _, parent := range parents {
						log.Printf("%v (%v) has Parent %v\n", name, *node.content.ResourceRecords[0].Value, *parent.content.Name)
						parent.children = append(parent.children, node)
						node.parent = parent
					}
				} else {
					log.Printf("%v (%v) is Out of domain\n", name, *node.content.ResourceRecords[0].Value)
				}
			} else {
				log.Printf("%v (%v) is Root\n", name, *node.content.ResourceRecords[0].Value)
			}
		}
	}
	CleanReferenceTrees(rootLookup)
	log.Printf("Cleared children, %d records left in Lookup\n", len(rootLookup))
}

// GetReferenceTrees returns the reference trees among the records defined in records parameter.
func GetReferenceTrees(records []*route53.ResourceRecordSet) (rootLookup map[string][]*Node) {
	filteredRecords := FilterReferrals(records)
	rootLookup = FillReferenceTrees(filteredRecords)
	CompactReferenceTrees(rootLookup)
	return
}

// PrintReferenceTrees prints reference trees contents.
func PrintReferenceTrees(referenceTrees map[string][]*Node) {
	for _, tree := range referenceTrees {
		for _, node := range tree {
			PrintNode(node, 0)
		}
	}
}
