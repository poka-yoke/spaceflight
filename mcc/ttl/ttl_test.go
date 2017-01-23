package ttl

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
)

var one = "one.example.com"
var two = "two.example.com"
var hundred = "100"
var zoneName = "test"
var A = "A"
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

var hostedZone = route53.HostedZone{
	CallerReference: &hundred,
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

var ucltest = []struct {
	original []*route53.ResourceRecordSet
	ttl      int64
}{
	{
		ResourceRecordSetList,
		60,
	},
	{
		ResourceRecordSetList,
		300,
	},
}

func TestUpsertChangeList(t *testing.T) {
	for _, tt := range ucltest {
		for _, change := range upsertChangeList(tt.original, tt.ttl) {
			if *change.ResourceRecordSet.TTL != tt.ttl {
				t.Error("TTL doesn't match")
			}
		}
	}
}

var fstest = []string{
	"",
	"a,b",
	"a,,a,,c,,",
}

func TestFilterSet(t *testing.T) {
	for _, s := range fstest {
		t.Run(s, func(t *testing.T) {
			var f Filter
			_ = f.Set(s)
			j := strings.Join(f, ",")
			if s != j {
				t.Error(j, " doesn't match")
			}
		})
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
