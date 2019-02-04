package got

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"

	"github.com/poka-yoke/spaceflight/internal/http"
)

// GetResourceRecordSet returns a slice containing all responses for specified
// query. It may issue more than one request as each returns a fixed amount of
// entries at most.
func GetResourceRecordSet(
	zoneID string,
	svc route53iface.Route53API,
) (
	resourceRecordSet []*route53.ResourceRecordSet,
	err error,
) {
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneID),
	}
	for respIsTruncated := true; respIsTruncated; {
		var resp *route53.ListResourceRecordSetsOutput
		resp, err = svc.ListResourceRecordSets(params)
		if err != nil {
			return
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
		err = fmt.Errorf("no records to process")
		return
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
		err = fmt.Errorf("no records to process")
		return
	}
	// Create batch with all jobs
	changeBatch := &route53.ChangeBatch{
		Changes: changes,
	}
	if err = changeBatch.Validate(); err != nil {
		return
	}

	changeRRSInput := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch:  changeBatch,
		HostedZoneId: zoneID,
	}

	if err = changeRRSInput.Validate(); err != nil {
		return
	}
	// Submit batch changes
	changeResponse, err = svc.ChangeResourceRecordSets(changeRRSInput)
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
func GetZoneID(zoneName string, svc route53iface.Route53API) (zoneID string, err error) {
	params := &route53.ListHostedZonesByNameInput{
		DNSName:  aws.String(zoneName),
		MaxItems: aws.String("100"),
	}
	resp, err := svc.ListHostedZonesByName(params)
	if err != nil {
		return
	}
	if len(resp.HostedZones) == 0 {
		err = fmt.Errorf("no results for zone %s. Exiting", zoneName)
		return
	}
	zoneID = *resp.HostedZones[0].Id
	return
}
