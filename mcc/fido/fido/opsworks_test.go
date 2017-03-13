package fido

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/opsworks"
	"github.com/aws/aws-sdk-go/service/opsworks/opsworksiface"
)

var list = []*opsworks.Stack{
	{
		Name:       aws.String("test1"),
		StackId:    aws.String("test1"),
		CustomJson: aws.String("test1"),
	},
	{
		Name:       aws.String("test2"),
		StackId:    aws.String("test2"),
		CustomJson: aws.String("test2"),
	},
	{
		Name:       aws.String("test3"),
		StackId:    aws.String("test3"),
		CustomJson: aws.String("test3"),
	},
	{
		Name:       aws.String("test4"),
		StackId:    aws.String("test4"),
		CustomJson: aws.String("test4"),
	},
	{
		Name:       aws.String("test5"),
		StackId:    aws.String("test5"),
		CustomJson: aws.String("test5"),
	},
	{
		Name:       aws.String("test6"),
		StackId:    aws.String("test6"),
		CustomJson: aws.String("test6"),
	},
	{
		Name:       aws.String("test7"),
		StackId:    aws.String("test7"),
		CustomJson: aws.String("test7"),
	},
	{
		Name:       aws.String("test8"),
		StackId:    aws.String("test8"),
		CustomJson: aws.String("test8"),
	},
	{
		Name:       aws.String("test9"),
		StackId:    aws.String("test9"),
		CustomJson: aws.String("test9"),
	},
	{
		Name:       aws.String("test10"),
		StackId:    aws.String("test10"),
		CustomJson: aws.String("test10"),
	},
}

func getStackMap() map[string]*opsworks.Stack {
	m := make(map[string]*opsworks.Stack)
	for _, i := range list {
		m[*i.Name] = i
	}
	return m
}

type mockOpsWorksClient struct {
	opsworksiface.OpsWorksAPI
}

func (m *mockOpsWorksClient) DescribeStacks(
	params *opsworks.DescribeStacksInput,
) (
	out *opsworks.DescribeStacksOutput,
	err error,
) {
	if params != nil {
		stacks := getStackMap()
		for _, sID := range params.StackIds {
			stack := stacks[*sID]
			out = &opsworks.DescribeStacksOutput{
				Stacks: []*opsworks.Stack{
					stack,
				},
			}
			return
		}
	}
	out = &opsworks.DescribeStacksOutput{
		Stacks: list,
	}
	return
}

func (m *mockOpsWorksClient) UpdateStack(
	params *opsworks.UpdateStackInput,
) (
	out *opsworks.UpdateStackOutput,
	err error,
) {
	ma := getStackMap()
	s := ma[*params.StackId]
	s.CustomJson = params.CustomJson
	return
}

func TestGetStackID(t *testing.T) {
	svc := &mockOpsWorksClient{}
	for name, tt := range getStackMap() {
		out, err := GetStackID(svc, name)
		if err != nil {
			t.Errorf("Unexpected error %s\n", err.Error())
		}
		if out != *tt.StackId {
			t.Errorf(
				"Expected StackID %s but received %s\n",
				*tt.StackId,
				out,
			)
		}
	}
}

func TestGetCustomJSON(t *testing.T) {
	svc := &mockOpsWorksClient{}
	for _, tt := range list {
		out, err := GetCustomJSON(svc, *tt.StackId)
		if err != nil {
			t.Errorf("Unexpected error %s\n", err.Error())
		}
		if out != *tt.CustomJson {
			t.Errorf(
				"Expected StackID %s but received %s\n",
				*tt.CustomJson,
				out,
			)
		}
	}
}

func TestPushCustomJSON(t *testing.T) {
	svc := &mockOpsWorksClient{}
	cJSON := "NewCustomJSON"
	for _, tt := range list {
		err := PushCustomJSON(svc, *tt.StackId, cJSON)
		if err != nil {
			t.Errorf("Unexpected error %s\n", err.Error())
		}
		if *tt.CustomJson != cJSON {
			t.Errorf(
				"Expected %s but got %s\n",
				cJSON,
				*tt.CustomJson,
			)
		}
	}
}
