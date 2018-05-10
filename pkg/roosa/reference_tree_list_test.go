package roosa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
)

var one = "one.example.com"
var two = "two.example.com"
var hundred = "100"
var zoneName = "test"
var AAAA = "AAAA"
var fals = false
var duration1 int64 = 1
var duration5 int64 = 5

var onerecordA = &route53.ResourceRecordSet{
	Name: &one,
	Type: &A,
	TTL:  &duration1,
}

var tworecordAAAA = &route53.ResourceRecordSet{
	Name: &two,
	Type: &AAAA,
	TTL:  &duration5,
}

var ResourceRecordSetList = []*route53.ResourceRecordSet{
	onerecordA,
	tworecordAAAA,
}

var rrstest = []struct {
	rrsl       []*route53.ResourceRecordSet
	typeFilter []string
	nameFilter []string
	out        []*route53.ResourceRecordSet
}{
	{
		ResourceRecordSetList,
		[]string{"A"},
		[]string{"one.example.com"},
		[]*route53.ResourceRecordSet{
			onerecordA,
		},
	},
	{
		ResourceRecordSetList,
		[]string{"AAAA"},
		[]string{"two.example.com"},
		[]*route53.ResourceRecordSet{
			tworecordAAAA,
		},
	},
	{
		ResourceRecordSetList,
		[]string{"MX"},
		[]string{"non-existent"},
		[]*route53.ResourceRecordSet{},
	},
	{
		ResourceRecordSetList,
		[]string{"TXT"},
		[]string{"non-existent"},
		[]*route53.ResourceRecordSet{},
	},
	{
		ResourceRecordSetList,
		[]string{
			"A",
			"AAAA",
		},
		[]string{"one.example.com", "two.example.com"},
		[]*route53.ResourceRecordSet{
			onerecordA,
			tworecordAAAA,
		},
	},
}

// Define a mock struct to be used in unit tests.
type mockRoute53Client struct {
	route53iface.Route53API
}

func (m *mockRoute53Client) ListResourceRecordSets(
	params *route53.ListResourceRecordSetsInput,
) (out *route53.ListResourceRecordSetsOutput, err error) {
	out = &route53.ListResourceRecordSetsOutput{
		IsTruncated:        &fals,
		MaxItems:           &hundred,
		ResourceRecordSets: ResourceRecordSetList,
	}
	return
}

var hostedZone = route53.HostedZone{
	CallerReference: &hundred,
}

func (m *mockRoute53Client) ListHostedZonesByName(
	params *route53.ListHostedZonesByNameInput,
) (out *route53.ListHostedZonesByNameOutput, err error) {
	hostedZone.Id = &zoneName
	hostedZone.Name = &zoneName
	out = &route53.ListHostedZonesByNameOutput{
		IsTruncated: &fals,
		MaxItems:    &hundred,
		HostedZones: []*route53.HostedZone{&hostedZone},
	}
	return
}

func TestFilterResourceRecords(t *testing.T) {
	for _, tt := range rrstest {
		t.Run(
			strings.Join(tt.typeFilter, ","),
			func(t *testing.T) {
				fn := func(elem *route53.ResourceRecordSet, filter string) *route53.ResourceRecordSet {
					if *elem.Type == filter {
						return elem
					}
					return nil
				}
				r := FilterResourceRecords(
					tt.rrsl,
					tt.typeFilter,
					fn,
				)
				if len(r) != len(tt.out) {
					t.Error("Result has different length than expected")
				}
				for index, value := range r {
					if tt.out[index] != value {
						t.Error("Results don't match as expected")
					}
				}
			},
		)
		t.Run(
			strings.Join(tt.nameFilter, ","),
			func(t *testing.T) {
				fn := func(elem *route53.ResourceRecordSet, filter string) *route53.ResourceRecordSet {
					if *elem.Name == filter {
						return elem
					}
					return nil
				}
				r := FilterResourceRecords(
					tt.rrsl,
					tt.nameFilter,
					fn,
				)
				if len(r) != len(tt.out) {
					t.Error("Result has different length than expected")
				}
				for index, value := range r {
					if tt.out[index] != value {
						t.Error("Results don't match as expected")
					}
				}
			},
		)
	}
}

var grrstest = []string{
	"test",
}

func TestGetResourceRecordSet(t *testing.T) {
	mockSvc := &mockRoute53Client{}
	for _, s := range grrstest {
		t.Run(s, func(t *testing.T) {
			out := GetResourceRecordSet(s, mockSvc)
			if len(out) != len(ResourceRecordSetList) {
				t.Error("Response doesn't match")
			}
			for i, val := range ResourceRecordSetList {
				if out[i] != val {
					t.Error("Response doesn't match")
				}
			}
		})
	}
}

var gzitest = []string{
	"test",
}

func TestGetZoneID(t *testing.T) {
	mockSvc := &mockRoute53Client{}
	for _, s := range grrstest {
		t.Run(s, func(t *testing.T) {
			out := GetZoneID(s, mockSvc)
			if out != s {
				t.Error("Response doesn't match")
			}
		})
	}
}

var recordsContents = []string{
	"example.com. SOA admin.example.com 2017010101 86400 7200 604800 300",
	"example.com. A 10.10.10.10",
	"root.example.com. A 127.0.0.1",
	"root-son.example.com. CNAME root.example.com",
	"root-son-sibling.example.com. CNAME root.example.com",
	"service1.example.com. CNAME root.example.com",
	"service1.example.com. CNAME root2.example.com",
	"root-grandson.example.com. CNAME root-son.example.com",
	"root2.example.com. A 127.0.0.2",
	"test.example.com. CNAME test.example2.com",
	"multiple-a.example.com. A 127.0.0.1,127.0.0.2,127.0.0.3",
}

var expectedOutputContents = []string{
	"root2.example.com. A 127.0.0.2",
	"\tservice1.example.com. CNAME root2.example.com",
	"test.example.com. CNAME test.example2.com",
	"multiple-a.example.com. A 127.0.0.1, 127.0.0.2, 127.0.0.3",
	"example.com. A 10.10.10.10",
	"root.example.com. A 127.0.0.1",
	"\troot-son-sibling.example.com. CNAME root.example.com",
	"\troot-son.example.com. CNAME root.example.com",
	"\t\troot-grandson.example.com. CNAME root-son.example.com",
	"\tservice1.example.com. CNAME root.example.com",
	"",
}

func generateRoute53RRS() (records []*route53.ResourceRecordSet) {
	records = []*route53.ResourceRecordSet{}
	for _, record := range recordsContents {
		recordValues := strings.Fields(record)
		values := strings.Split(recordValues[2], ",")
		rrs := []*route53.ResourceRecord{}
		for _, value := range values {
			newValue := value
			rr := &route53.ResourceRecord{Value: &newValue}
			rrs = append(rrs, rr)
		}
		route53RRS := &route53.ResourceRecordSet{
			Name:            &recordValues[0],
			Type:            &recordValues[1],
			ResourceRecords: rrs,
		}
		records = append(records, route53RRS)
	}
	return
}

func TestReferenceTreeListString(t *testing.T) {
	rtl := NewReferenceTreeList(generateRoute53RRS())
	output := fmt.Sprintf("%v", rtl)
	for _, line := range expectedOutputContents {
		if !strings.Contains(output, line) {
			t.Errorf("'%v' should be in output:\n%v", line, output)
		}
	}
	outputLines := strings.Split(output, "\n")
	if len(outputLines) != len(expectedOutputContents) {
		t.Errorf(
			"Output line count (%d) doesn't match with expected (%d)",
			len(outputLines),
			len(expectedOutputContents),
		)
	}
}
