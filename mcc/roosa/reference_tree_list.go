package roosa

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/Devex/spaceflight/mcc/ttl"
)

// ReferenceTreeList is a type representing the reference trees for a list of
// DNS records, explicitly A, AAAA, and CNAME records.
type ReferenceTreeList struct {
	records []*route53.ResourceRecordSet
	lookup  map[string][]*Node
}

var recordTypes = []string{
	"A",
	"AAAA",
	"CNAME",
}

// NewReferenceTreeList is the constructor for ReferenceTreeList. It filters
// A, AAAA, and CNAME records from `records` argument.
func NewReferenceTreeList(
	records []*route53.ResourceRecordSet,
) *ReferenceTreeList {
	return &ReferenceTreeList{
		records: ttl.FilterResourceRecords(
			records,
			recordTypes,
			func(elem *route53.ResourceRecordSet, filter string) *route53.ResourceRecordSet {
				if *elem.Type == filter {
					return elem
				}
				return nil
			},
		),
		lookup: nil,
	}
}

// GetReferenceTrees builds and returns the reference trees among the
// ReferenceTreeList records attribute.
func (rtl *ReferenceTreeList) GetReferenceTrees() map[string][]*Node {
	rtl.fill()
	rtl.compact()
	return rtl.lookup
}

// String returns a string representing ReferenceTreeList contents.
func (rtl *ReferenceTreeList) String() (output string) {
	if rtl.lookup == nil {
		rtl.GetReferenceTrees()
	}
	output = ""
	for _, tree := range rtl.lookup {
		for _, node := range tree {
			output += fmt.Sprintf("%v\n", node)
		}
	}
	return
}

// fill fills the referral lookup table with the base records.
func (rtl *ReferenceTreeList) fill() {
	rtl.lookup = map[string][]*Node{}
	for _, val := range rtl.records {
		node := &Node{
			content: val,
		}
		name := strings.TrimSuffix(*val.Name, ".")
		rtl.lookup[name] = append(rtl.lookup[name], node)
	}
	log.Printf("Added %d records to Lookup\n", len(rtl.lookup))
	return
}

// clean removes non-root elements from the referral lookup table.
func (rtl *ReferenceTreeList) clean() {
	for name, nodes := range rtl.lookup {
		for _, node := range nodes {
			if !node.IsRoot() {
				delete(rtl.lookup, name)
			}
		}
	}
}

// compact modifies the referral lookup table finding children and
// roots, and relating these appropriately.
func (rtl *ReferenceTreeList) compact() {
	for name, nodes := range rtl.lookup {
		for _, node := range nodes {
			value := *node.content.ResourceRecords[0].Value
			if *node.content.Type == "CNAME" {
				if parents, ok := rtl.lookup[value]; ok {
					for _, parent := range parents {
						log.Printf(
							"%v (%v) has Parent %v\n",
							name,
							value,
							*parent.content.Name,
						)
						parent.children = append(
							parent.children,
							node,
						)
						node.parent = parent
					}
				} else {
					log.Printf(
						"%v (%v) is Out of domain\n",
						name,
						value,
					)
				}
			} else {
				log.Printf(
					"%v (%v) is Root\n",
					name,
					value,
				)
			}
		}
	}
	rtl.clean()
	log.Printf(
		"Cleared children, %d records left in Lookup\n",
		len(rtl.lookup),
	)
}
