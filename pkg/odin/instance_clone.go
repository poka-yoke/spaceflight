package odin

import (
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CloneInstance creates a new RDS database instance, copying parameters
// from a snapshot. If a vpcid is specified the security group will be
// in that VPC.
func CloneInstance(
	params Instance,
	svc rdsiface.RDSAPI,
) (
	result string,
	err error,
) {
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
	err = ModifyInstance(params, svc)
	return
}
