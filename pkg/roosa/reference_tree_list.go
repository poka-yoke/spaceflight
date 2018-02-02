package roosa

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
)

// GetResourceRecordSet returns a slice containing all responses for specified
// query. It may issue more than one request as each returns a fixed amount of
// entries at most.
func GetResourceRecordSet(
	zoneID string,
	svc route53iface.Route53API,
) (resourceRecordSet []*route53.ResourceRecordSet) {
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
	}
	for respIsTruncated := true; respIsTruncated; {
		resp, err := svc.ListResourceRecordSets(params)
		if err != nil {
			panic(err)
		}
		if *resp.IsTruncated {
			params.StartRecordName = resp.NextRecordName
			params.StartRecordType = resp.NextRecordType
		}
		respIsTruncated = *resp.IsTruncated

		// Iterate over all entries and add changes to change_slice
		for _, val := range resp.ResourceRecordSets {
			resourceRecordSet = append(resourceRecordSet, val)
		}
	}
	return
}

// GetZoneID returns a string containing the ZoneID for use in further API
// actions
func GetZoneID(zoneName string, svc route53iface.Route53API) (zoneID string) {
	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(zoneName),
		MaxItems: aws.String("100"),
	}
	resp, err := svc.ListHostedZonesByName(params)
	if err != nil {
		log.Println(err.Error())
	}
	zoneID = *resp.HostedZones[0].Id
	return
}

// FilterResourceRecords returns a slice containing only the entries that
// pass the check performed by the function argument
func FilterResourceRecords(
	l []*route53.ResourceRecordSet,
	f []string,
	p func(*route53.ResourceRecordSet, string) *route53.ResourceRecordSet,
) (result []*route53.ResourceRecordSet) {
	for _, elem := range l {
		for _, filter := range f {
			res := p(elem, filter)
			if res != nil {
				result = append(result, res)
			}
		}
	}
	return
}

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
		records: FilterResourceRecords(
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
