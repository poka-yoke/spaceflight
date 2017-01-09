package roosa

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
)

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
