package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CloneInstance creates a new RDS database instance, copying parameters
// from a snapshot. If a vpcid is specified the security group will be
// in that VPC.
func CloneInstance(
	instanceName string,
	params Instance,
	svc rdsiface.RDSAPI,
) (
	result string,
	err error,
) {
	if params.OriginalInstanceName == "" {
		return "", fmt.Errorf("Original instance name not provided")
	}
	rdsParams, err := params.CloneDBInput(
		svc,
	)
	if err != nil {
		return "", err
	}
	res, err := svc.CreateDBInstance(rdsParams)
	if err != nil {
		return "", err
	}
	err = WaitForInstance(res.DBInstance, svc, "available")
	if err != nil {
		return
	}
	result = *res.DBInstance.Endpoint.Address
	err = modifyInstance(instanceName, params, svc)
	return
}
