package odin

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CreateInstance creates a new RDS database instance. If a vpcid is
// specified the security group will be in that VPC.
func CreateInstance(
	instanceName string,
	params Instance,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	var instance *rds.DBInstance
	instance, err = doCreate(instanceName, params, svc)
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

func doCreate(
	instanceName string,
	params Instance,
	svc rdsiface.RDSAPI,
) (
	instance *rds.DBInstance,
	err error,
) {
	rdsParams, err := params.CreateDBInput(
		instanceName,
		svc,
	)
	if err != nil {
		return nil, err
	}
	res, err := svc.CreateDBInstance(rdsParams)
	if err != nil {
		return
	}
	return res.DBInstance, nil
}
