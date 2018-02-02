package roosa

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/route53"
)

var A = "A"
var CNAME = "CNAME"

var parentName = "parent.example.com"
var parentValue = "10.10.10.10"
var parentRRS = &route53.ResourceRecordSet{
	Name: &parentName,
	Type: &A,
	ResourceRecords: []*route53.ResourceRecord{
		&route53.ResourceRecord{
			Value: &parentValue,
		},
	},
}
var parent = &Node{
	content: parentRRS,
}
var child1Name = "child1.example.com"
var child1RRS = &route53.ResourceRecordSet{
	Name: &child1Name,
	Type: &CNAME,
	ResourceRecords: []*route53.ResourceRecord{
		&route53.ResourceRecord{
			Value: &parentName,
		},
	},
}
var child1 = &Node{
	content: child1RRS,
	parent:  parent,
}
var child2Name = "child2.example.com"
var child2RRS = &route53.ResourceRecordSet{
	Name: &child2Name,
	Type: &CNAME,
	ResourceRecords: []*route53.ResourceRecord{
		&route53.ResourceRecord{
			Value: &child1Name,
		},
	},
}
var child2 = &Node{
	content: child2RRS,
	parent:  child1,
}

func TestIsRoot(t *testing.T) {
	if !parent.IsRoot() {
		t.Errorf(
			"%v is supposed to be root",
			parent,
		)
	}
	if child1.IsRoot() {
		t.Errorf(
			"%v is not supposed to be root",
			child1,
		)
	}
	if child2.IsRoot() {
		t.Errorf(
			"%v is not supposed to be root",
			child2,
		)
	}
}

func TestString(t *testing.T) {
	parent.children = append(parent.children, child1)
	child1.children = append(child1.children, child2)
	parentOutput := parent.String()
	parentExpectedOutput := "parent.example.com A 10.10.10.10\n\tchild1.example.com CNAME parent.example.com\n\t\tchild2.example.com CNAME child1.example.com"
	if parentOutput != parentExpectedOutput {
		t.Errorf("Output doesn't match Expected: \n%v\n-----\n%v", parentOutput, parentExpectedOutput)
	}
}
