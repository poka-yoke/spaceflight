package roosa

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/poka-yoke/spaceflight/mcc/ttl"
)

// FilterReferrals returns 'A', 'AAAA', and 'CNAME' records from a
// []*route53.ResourceRecordSet.
func FilterReferrals(
	records []*route53.ResourceRecordSet,
) []*route53.ResourceRecordSet {
	return ttl.FilterResourceRecordSetType(
		records,
		[]string{
			"A",
			"AAAA",
			"CNAME",
		},
	)
}

// FillReferenceTrees fills a referral lookup table with the base
// records from a []*route53.ResourceRecordSet.
func FillReferenceTrees(
	records []*route53.ResourceRecordSet,
) (rootLookup map[string][]*Node) {
	rootLookup = map[string][]*Node{}
	for _, val := range records {
		node := &Node{
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
			if !node.IsRoot() {
				delete(rootLookup, name)
			}
		}
	}
}

// CompactReferenceTrees modifies a referral lookup table finding children and
// roots, and relating these appropriately.
func CompactReferenceTrees(rootLookup map[string][]*Node) {
	for name, nodes := range rootLookup {
		for _, node := range nodes {
			value := *node.content.ResourceRecords[0].Value
			if *node.content.Type == "CNAME" {
				if parents, ok := rootLookup[value]; ok {
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
	CleanReferenceTrees(rootLookup)
	log.Printf(
		"Cleared children, %d records left in Lookup\n",
		len(rootLookup),
	)
}

// GetReferenceTrees returns the reference trees among the records defined in
// records parameter.
func GetReferenceTrees(
	records []*route53.ResourceRecordSet,
) (rootLookup map[string][]*Node) {
	filteredRecords := FilterReferrals(records)
	rootLookup = FillReferenceTrees(filteredRecords)
	CompactReferenceTrees(rootLookup)
	return
}

// PrintReferenceTrees prints reference trees contents.
func PrintReferenceTrees(referenceTrees map[string][]*Node) {
	for _, tree := range referenceTrees {
		for _, node := range tree {
			fmt.Print(node)
		}
	}
}
