package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"log"
)

var verbose bool

// GetResourceRecordSet returns a slice containing all responses for specified
// query. It may issue more than one request as each returns a fixed amount of
// entries at most.
func GetResourceRecordSet(
	params *route53.ListResourceRecordSetsInput,
	svc *route53.Route53,
) (resourceRecordSet []*route53.ResourceRecordSet) {
	for respIsTruncated := true; respIsTruncated; {
		if verbose {
			log.Printf("Query params: %s\n", params)
		}
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

// upsertChangeList iterates over a list of records and returns a list of
// Change objects of type Upsert with the specified TTL
func upsertChangeList(
	list []*route53.ResourceRecordSet,
	ttl int64,
) (res []*route53.Change) {
	for _, val := range list {
		*val.TTL = ttl
		change := &route53.Change{
			Action:            aws.String("UPSERT"),
			ResourceRecordSet: val,
		}
		log.Printf(
			"Adding %s to change list for TTL %d\n",
			*val.Name,
			ttl)
		res = append(res, change)
	}
	return
}

// UpsertResourceRecordSetTTL performs the request to change the TTL of the list
// of records.
func UpsertResourceRecordSetTTL(
	list []*route53.ResourceRecordSet,
	ttl int64,
	zoneID *string,
	svc *route53.Route53,
) {
	changeSlice := upsertChangeList(list, ttl)

	// Create batch with all jobs
	changeBatch := &route53.ChangeBatch{
		Changes: changeSlice,
	}
	if err := changeBatch.Validate(); err != nil {
		log.Panic(err)
	}

	changeRRSInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  changeBatch,
		HostedZoneId: zoneID,
	}

	if err := changeRRSInput.Validate(); err != nil {
		log.Panic(err.Error())
	}
	// Submit batch changes
	resp2, err := svc.ChangeResourceRecordSets(changeRRSInput)
	if err != nil {
		log.Panic(err)
	}
	if verbose {
		fmt.Println(resp2)
	}
}

// PrintRecords prints all records in a zone using the API's built-in method
// ListResourceRecordSets.
func PrintRecords(
	p *route53.ListResourceRecordSetsOutput,
	last bool,
) (shouldContinue bool) {
	shouldContinue = *p.IsTruncated
	for idx, val := range p.ResourceRecordSets {
		fmt.Println(idx, *val)
	}
	return
}

// FilterResourceRecordSetType returns a slice containing only the entries with
// specified types of the original record slice
func FilterResourceRecordSetType(
	l []*route53.ResourceRecordSet,
	f []string,
) (result []*route53.ResourceRecordSet) {
	for _, elem := range l {
		for _, filter := range f {
			if *elem.Type == filter {
				result = append(result, elem)
			}
		}
	}
	return
}

// SplitResourceRecordSetTypeOnNames returns two slices: one containing all the
// entries from the l argument, and another with the results of excluding the
// matches from the f argument.
func SplitResourceRecordSetTypeOnNames(
	l []*route53.ResourceRecordSet,
	f []string,
) (result1 []*route53.ResourceRecordSet, result2 []*route53.ResourceRecordSet) {
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
	zoneName := flag.String(
		"zonename",
		"",
		"Hosted Zone's name to traverse")
	ttl := flag.Int64("ttl", 300, "Desired TTL value")
	flag.BoolVar(&verbose, "v", false, "Increments output")

	flag.Parse()

	if *zoneName == "" {
		log.Fatal("Insufficient input parameters!")
	}

	sess, err := session.NewSession()
	if err != nil {
		log.Panicf("Failed to create session: %s", err)
		return
	}

	svc := route53.New(sess)

	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(*zoneName),
		MaxItems: aws.String("100"),
	}
	resp, err := svc.ListHostedZonesByName(params)
	if err != nil {
		log.Println(err.Error())
	}

	zoneID := resp.HostedZones[0].Id
	if verbose {
		// Pretty-print the response data.
		fmt.Println(*zoneID)
	}

	params2 := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(*zoneID),
	}
	list := GetResourceRecordSet(params2, svc)
	// Filter list in between
	filter := []string{
		"A",
		"AAAA",
		"CNAME",
		"MX",
		"NAPTR",
		"PTR",
		"SPF",
		"SRV",
		"TXT",
	}
	list = FilterResourceRecordSetType(list, filter)
	UpsertResourceRecordSetTTL(list, *ttl, zoneID, svc)
}
