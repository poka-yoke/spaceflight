package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
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
	var instance *rds.DBInstance
	if params.OriginalInstanceName == "" {
		return "", fmt.Errorf("Original instance name not provided")
	}
	instance, err = doClone(instanceName, params, svc)
	if err != nil {
		return
	}
	err = WaitForInstance(instance, svc, "available")
	if err != nil {
		return
	}
	result = *instance.Endpoint.Address
	err = modifyInstance(instanceName, params, svc)
	return
}

func doClone(
	instanceName string,
	params Instance,
	svc rdsiface.RDSAPI,
) (
	instance *rds.DBInstance,
	err error,
) {
	pparams, err := params.AddLastSnapshot(params.OriginalInstanceName, svc)
	if err != nil {
		return nil, err
	}
	rdsParams, err := pparams.CreateDBInput(
		instanceName,
		svc,
	)
	if err != nil {
		return
	}
	res, err := svc.CreateDBInstance(rdsParams)
	if err != nil {
		return
	}
	return res.DBInstance, nil
}
