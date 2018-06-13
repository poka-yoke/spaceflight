package got

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"

	"github.com/Devex/spaceflight/pkg/internal/http"
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

// Init initializes conections to Route53
func Init() *route53.Route53 {
	sess, err := session.NewSession()
	if err != nil {
		log.Panicf("Failed to create session: %s", err)
	}
	return route53.New(sess)
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

// NewResourceRecordList creates a list of ResourceRecords with a ResourceRecord
// per each string passed.
func NewResourceRecordList(values []string) (ret []*route53.ResourceRecord) {
	for _, val := range values {
		ret = append(ret,
			&route53.ResourceRecord{
				Value: aws.String(val),
			},
		)
	}
	return
}

// UpsertChangeList generates a list of changes for UPSERT the records in list
// according with ttl, name, and type
func UpsertChangeList(
	list []*route53.ResourceRecord,
	ttl int64,
	name string,
	typ string,
) (res []*route53.Change) {
	val := &route53.ResourceRecordSet{
		ResourceRecords: list,
		TTL:             &ttl,
		Type:            &typ,
		Name:            &name,
	}
	change := &route53.Change{
		Action:            aws.String("UPSERT"),
		ResourceRecordSet: val,
	}
	log.Printf(
		"Adding %s to change list for TTL %d\n",
		*val.Name,
		ttl,
	)
	res = append(res, change)
	return
}

// DeleteChangeList generates a list of changes for DELETEing the records in
// names, according to a common type
func DeleteChangeList(
	names []string,
	typ string,
	list []*route53.ResourceRecordSet,
) (res []*route53.Change) {
	var record *route53.ResourceRecordSet
	for _, name := range names {
		for _, i := range list {
			if *i.Name == name &&
				*i.Type == typ {
				record = i
			}
		}
		change := &route53.Change{
			Action:            aws.String("DELETE"),
			ResourceRecordSet: record,
		}
		res = append(res, change)
		log.Printf("Added %s to delete list\n", name)
	}
	return
}

// WaitForChangeToComplete waits until the ChangeInfo described by the argument is completed.
func WaitForChangeToComplete(
	changeInfo *route53.ChangeInfo,
	svc *route53.Route53,
) {
	getChangeInput := &route53.GetChangeInput{Id: changeInfo.Id}
	req, getChangeOutput := svc.GetChangeRequest(getChangeInput)

	http.NewLinearBackoff(
		func() bool {
			return *getChangeOutput.ChangeInfo.Status ==
				route53.ChangeStatusInsync
		},
		time.Second,
		30*time.Second,
	).Do(req)

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
	svc route53iface.Route53API,
) (
	changeResponse *route53.ChangeResourceRecordSetsOutput,
	err error,
) {
	if len(list) <= 0 {
		log.Fatal("No records to process.")
	}
	changeSlice := []*route53.Change{}
	for _, r := range list {
		partialChangeSlice := UpsertChangeList(r.ResourceRecords, ttl, *r.Name, *r.Type)
		changeSlice = append(changeSlice, partialChangeSlice...)
	}

	changeResponse, err = ApplyChanges(changeSlice, zoneID, svc)
	return
}

// ApplyChanges performs the request to change the list
// of records.
func ApplyChanges(
	changes []*route53.Change,
	zoneID *string,
	svc route53iface.Route53API,
) (
	changeResponse *route53.ChangeResourceRecordSetsOutput,
	err error,
) {
	if len(changes) <= 0 {
		log.Fatal("No records to process.")
	}
	// Create batch with all jobs
	changeBatch := &route53.ChangeBatch{
		Changes: changes,
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
		if err != nil {
			log.Panic(err)
		}
		if Verbose {
			fmt.Println(changeResponse.ChangeInfo)
		}
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
	p func(*route53.ResourceRecordSet, string,
	) *route53.ResourceRecordSet,
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
	if len(resp.HostedZones) == 0 {
		log.Fatalf("No results for zone %s. Exiting.\n", zoneName)
	}
	zoneID = *resp.HostedZones[0].Id
	if Verbose {
		// Pretty-print the response data.
		fmt.Println(zoneID)
	}
	return
}
