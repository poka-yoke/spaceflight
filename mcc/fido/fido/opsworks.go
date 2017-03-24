package fido

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/opsworks"
	"github.com/aws/aws-sdk-go/service/opsworks/opsworksiface"
)

// Init initializes connection to AWS API
func Init() opsworksiface.OpsWorksAPI {
	region := "us-east-1"
	sess := session.New(&aws.Config{Region: aws.String(region)})
	return opsworks.New(sess)
}

// GetStackID returns the ID of a stack whose name matches the passed string
func GetStackID(svc opsworksiface.OpsWorksAPI, name string) (string, error) {
	out, err := svc.DescribeStacks(nil)
	if err != nil {
		log.Panic(err)
	}
	for _, stack := range out.Stacks {
		if name == *stack.Name {
			return *stack.StackId, nil
		}
	}
	return "", fmt.Errorf("No stack matches %s", name)
}

// GetCustomJSON obtains the CustomJSON string from OpsWorks
func GetCustomJSON(
	svc opsworksiface.OpsWorksAPI,
	stackID string,
) (
	string,
	error,
) {
	params := &opsworks.DescribeStacksInput{
		StackIds: []*string{
			&stackID,
		},
	}
	out, err := svc.DescribeStacks(params)
	if err != nil {
		log.Panic(err)
	}
	if len(out.Stacks) < 1 {
		return "", fmt.Errorf("No Stack found for ID %s", stackID)
	}
	return *out.Stacks[0].CustomJson, nil
}

// PushCustomJSON uploads a CustomJSON to an OpsWorks stack
func PushCustomJSON(svc opsworksiface.OpsWorksAPI, stackid, customJSON string) (err error) {
	params := &opsworks.UpdateStackInput{
		StackId:    &stackid,
		CustomJson: &customJSON,
	}
	_, err = svc.UpdateStack(params)
	return
}
