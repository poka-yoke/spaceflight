package got

import (
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/aws/aws-sdk-go/service/route53/route53iface"
)

var one = "one.example.com."
var two = "two.example.com."
var awsCname = "ec2-1-2-3-4.compute-1.amazonaws.com"
var hundred = "100"
var zoneName = "test"
var A = "A"
var AAAA = "AAAA"
var fals = false
var duration1 int64 = 1
var duration5 int64 = 5
var now = time.Now()

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

func (m *mockRoute53Client) ChangeResourceRecordSets(
	params *route53.ChangeResourceRecordSetsInput,
) (out *route53.ChangeResourceRecordSetsOutput, err error) {
	out = &route53.ChangeResourceRecordSetsOutput{
		ChangeInfo: &route53.ChangeInfo{
			Id:     params.HostedZoneId,
			Status: pstr("PENDING"),
		},
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
		[]string{"one.example.com."},
		[]*route53.ResourceRecordSet{
			onerecordA,
		},
	},
	{
		ResourceRecordSetList,
		[]string{"AAAA"},
		[]string{"two.example.com."},
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
		[]string{"one.example.com.", "two.example.com."},
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

var rrltest = []struct {
	in  []string
	out []*route53.ResourceRecord
}{
	{
		[]string{
			one,
		},
		[]*route53.ResourceRecord{
			&route53.ResourceRecord{
				Value: &one,
			},
		},
	},
	{
		[]string{
			two,
		},
		[]*route53.ResourceRecord{
			&route53.ResourceRecord{
				Value: &two,
			},
		},
	},
	{
		[]string{
			one,
			two,
		},
		[]*route53.ResourceRecord{
			&route53.ResourceRecord{
				Value: &one,
			},
			&route53.ResourceRecord{
				Value: &two,
			},
		},
	},
	{
		[]string{
			awsCname,
		},
		[]*route53.ResourceRecord{
			&route53.ResourceRecord{
				Value: &awsCname,
			},
		},
	},
}

func TestNewResourceRecordList(t *testing.T) {
	for _, tt := range rrltest {
		t.Run(strings.Join(tt.in, ";"), func(t *testing.T) {
			out := NewResourceRecordList(tt.in)
			if len(tt.in) != len(out) {
				t.Error(
					"Erroneous amount of responses."+
						" Expected %d, received %d.",
					len(tt.in),
					len(out),
				)
			}
			for i, v := range out {
				if *v.Value != *tt.out[i].Value {
					t.Error(
						"Erroneous response."+
							" Expected %s, received %s.",
						*tt.out[i].Value,
						*v,
					)
				}
			}
		})
	}
}

func pstr(str string) *string {
	return &str
}

var actest = []struct {
	input struct {
		changes []*route53.Change
		zoneid  string
	}
	output struct {
		out *route53.ChangeResourceRecordSetsOutput
		err error
	}
}{
	{
		struct {
			changes []*route53.Change
			zoneid  string
		}{
			[]*route53.Change{
				&route53.Change{
					Action:            pstr("UPSERT"),
					ResourceRecordSet: onerecordA,
				},
			},
			"test1",
		},
		struct {
			out *route53.ChangeResourceRecordSetsOutput
			err error
		}{
			&route53.ChangeResourceRecordSetsOutput{
				ChangeInfo: &route53.ChangeInfo{
					Id: pstr("test1"),
					// Possible values: PENDING | INSYNC
					Status: pstr("PENDING"),
					// It's a *time.Time
					SubmittedAt: &now,
				},
			},
			nil,
		},
	},
}

func TestApplyChanges(t *testing.T) {
	mockSvc := &mockRoute53Client{}
	for _, tt := range actest {
		t.Run(tt.input.zoneid, func(t *testing.T) {
			out, err := ApplyChanges(
				tt.input.changes,
				&tt.input.zoneid,
				mockSvc,
			)
			if out == nil {
				t.Error("Returned nil")
			} else if *out.ChangeInfo.Id != *tt.output.out.ChangeInfo.Id ||
				*out.ChangeInfo.Status != *tt.output.out.ChangeInfo.Status ||
				err != tt.output.err {
				t.Error("Unexpected outcome")
			}
		})
	}
}

var dcltest = []struct {
	names  []string
	typ    string
	result []*route53.ResourceRecordSet
}{
	{
		typ: "CNAME",
		names: []string{
			"one.example.com.",
		},
		result: []*route53.ResourceRecordSet{ResourceRecordSetList[0]},
	},
}

func TestDeleteChangeList(t *testing.T) {
	for _, tt := range dcltest {
		res := DeleteChangeList(tt.names, tt.typ, ResourceRecordSetList)
		if len(res) != len(tt.names) {
			t.Errorf(
				"Unexpected length of results, expected %d and got %d\n",
				len(tt.names),
				len(res),
			)
		}
		for _, change := range res {
			if *change.Action != "DELETE" {
				t.Errorf(
					"Expected DELETE action, got %s\n",
					*change.Action,
				)
			}
		}
	}
}
