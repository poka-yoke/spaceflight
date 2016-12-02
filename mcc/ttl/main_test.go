package main

import (
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
)

var one = "one.example.com"
var two = "two.example.com"

var A = "A"
var AAAA = "AAAA"
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
	rrsl   []*route53.ResourceRecordSet
	filter []string
	out    []*route53.ResourceRecordSet
}{
	{
		ResourceRecordSetList,
		[]string{"A"},
		[]*route53.ResourceRecordSet{
			onerecordA,
		},
	},
	{
		ResourceRecordSetList,
		[]string{"AAAA"},
		[]*route53.ResourceRecordSet{
			tworecordAAAA,
		},
	},
	{
		ResourceRecordSetList,
		[]string{"MX"},
		[]*route53.ResourceRecordSet{},
	},
	{
		ResourceRecordSetList,
		[]string{"TXT"},
		[]*route53.ResourceRecordSet{},
	},
	{
		ResourceRecordSetList,
		[]string{
			"A",
			"AAAA",
		},
		[]*route53.ResourceRecordSet{
			onerecordA,
			tworecordAAAA,
		},
	},
}

func TestFilterResourceRecordSetType(t *testing.T) {
	for _, tt := range rrstest {
		t.Run(strings.Join(tt.filter, ","), func(t *testing.T) {
			r := FilterResourceRecordSetType(tt.rrsl, tt.filter)
			if len(r) != len(tt.out) {
				t.Error("Result has different length than expected")
			}
			for index, value := range r {
				if tt.out[index] != value {
					t.Error("Results don't match as expected")
				}
			}
		})
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