package ttl

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
)

// Verbose flag
var Verbose bool

// Dryrun flag
var Dryrun bool

// Filter type for generic filtering
type Filter []string

// String interface for Filter
func (f *Filter) String() string {
	return fmt.Sprint(*f)
}

// Set so flag can initialize flags of type Filter
func (f *Filter) Set(value string) error {
	for _, val := range strings.Split(value, ",") {
		*f = append(*f, val)
	}
	return nil
}

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
		if Verbose {
			fmt.Printf("Query params: %s\n", params)
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

// WaitForChangeToComplete waits until the ChangeInfo described by the argument is completed.
func WaitForChangeToComplete(
	changeInfo *route53.ChangeInfo,
	svc *route53.Route53,
) {
	getChangeInput := route53.GetChangeInput{Id: changeInfo.Id}
	getChangeOutput, err := svc.GetChange(&getChangeInput)
	if err != nil {
		log.Panic(err.Error())
	}
	for *getChangeOutput.ChangeInfo.Status != route53.ChangeStatusInsync {
		time.Sleep(1)
		getChangeOutput, err = svc.GetChange(&getChangeInput)
		if err != nil {
			log.Panic(err.Error())
		}
	}
	if Verbose {
		fmt.Println(getChangeOutput.ChangeInfo)
	}
	log.Println("All changes applied")
}

// UpsertResourceRecordSetTTL performs the request to change the TTL of the list
// of records.
func UpsertResourceRecordSetTTL(
	list []*route53.ResourceRecordSet,
	ttl int64,
	zoneID *string,
	svc *route53.Route53,
) (
	changeResponse *route53.ChangeResourceRecordSetsOutput,
	err error,
) {
	changeSlice := upsertChangeList(list, ttl)

	// Create batch with all jobs
	changeBatch := &route53.ChangeBatch{
		Changes: changeSlice,
	}
	if err := changeBatch.Validate(); err != nil {
		log.Panic(err.Error())
	}

	changeRRSInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  changeBatch,
		HostedZoneId: zoneID,
	}

	if err := changeRRSInput.Validate(); err != nil {
		log.Panic(err.Error())
	}
	// Submit batch changes
	if !Dryrun {
		changeResponse, err = svc.ChangeResourceRecordSets(changeRRSInput)
	}
	if err != nil {
		log.Panic(err)
	}
	if Verbose && !Dryrun {
		fmt.Println(changeResponse.ChangeInfo)
	}
	return
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
	if Verbose {
		// Pretty-print the response data.
		fmt.Println(zoneID)
	}
	return
}
