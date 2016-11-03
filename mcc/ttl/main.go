package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"log"
)

func GetResourceRecordSet(params *route53.ListResourceRecordSetsInput, svc *route53.Route53) (resource_record_set []*route53.ResourceRecordSet) {
	for resp_is_truncated := true; resp_is_truncated; {
		log.Printf("Query params: %s\n", params)
		resp, err := svc.ListResourceRecordSets(params)
		if err != nil {
			panic(err)
		}
		if *resp.IsTruncated {
			params.StartRecordName = resp.NextRecordName
			params.StartRecordType = resp.NextRecordType
		}
		resp_is_truncated = *resp.IsTruncated

		// Iterate over all entries and add changes to change_slice
		for _, val := range resp.ResourceRecordSets {
			resource_record_set = append(resource_record_set, val)
		}
	}
	return
}

func UpsertChangeList(list []*route53.ResourceRecordSet, ttl int64) (res []*route53.Change) {
	for _, val := range list {
		*val.TTL = ttl
		change := &route53.Change{
			Action:            aws.String("UPSERT"),
			ResourceRecordSet: val,
		}
		log.Printf("Adding %s to change list for TTL %d\n", *val.Name, ttl)
		res = append(res, change)
	}
	return
}

func UpsertResourceRecordSetTTL(list []*route53.ResourceRecordSet, ttl int64, zone_id *string, svc *route53.Route53) {
	change_slice := UpsertChangeList(list, ttl)

	// Create batch with all jobs
	change_batch := &route53.ChangeBatch{
		Changes: change_slice,
	}
	if err := change_batch.Validate(); err != nil {
		log.Panic(err)
	}

	change_resource_record_sets_input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  change_batch,
		HostedZoneId: zone_id,
	}

	if err := change_resource_record_sets_input.Validate(); err != nil {
		log.Panic(err.Error())
	}
	// Submit batch changes
	resp2, err := svc.ChangeResourceRecordSets(change_resource_record_sets_input)
	if err != nil {
		log.Panic(err)
	}
	fmt.Println(resp2)
}

func PrintRecords(p *route53.ListResourceRecordSetsOutput, last bool) (shouldContinue bool) {
	shouldContinue = *p.IsTruncated
	for idx, val := range p.ResourceRecordSets {
		fmt.Println(idx, *val)
	}
	return
}

func FilterResourceRecordSetType(l []*route53.ResourceRecordSet, f []string) (result []*route53.ResourceRecordSet) {
	for _, elem := range l {
		for _, filter := range f {
			if *elem.Type == filter {
				result = append(result, elem)
			}
		}
	}
	return
}

func SplitResourceRecordSetTypeOnNames(l []*route53.ResourceRecordSet, f []string) (result1 []*route53.ResourceRecordSet, result2 []*route53.ResourceRecordSet) {
	result1 = l
	for _, elem := range l {
		for _, filter := range f {
			if *elem.Name == filter {
				result2 = append(result2, elem)
				break
			}
		}
	}
	return
}

func main() {
	sess, err := session.NewSession()
	if err != nil {
		log.Panicf("Failed to create session,", err)
		return
	}

	svc := route53.New(sess)

	zone_name := flag.String("zonename", "", "Hosted Zone's name to traverse")
	flag.Parse()

	if *zone_name == "" {
		log.Fatal("Insufficient input parameters!")
	}

	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(*zone_name),
		MaxItems: aws.String("100"),
	}
	resp, err := svc.ListHostedZonesByName(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println(err.Error())
	}

	zone_id := resp.HostedZones[0].Id
	// Pretty-print the response data.
	fmt.Println(*zone_id)

	var params2 *route53.ListResourceRecordSetsInput = &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(*zone_id),
	}
	list := GetResourceRecordSet(params2, svc)
	// Filter list in between
	filter := []string{"A", "AAAA", "CNAME", "MX", "NAPTR", "PTR", "SPF", "SRV", "TXT"}
	list = FilterResourceRecordSetType(list, filter)
	UpsertResourceRecordSetTTL(list, 300, zone_id, svc)
}
