package odin

import (
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CreateInstance creates a new RDS database instance. If a vpcid is
// specified the security group will be in that VPC.
func CreateInstance(
	params Instance,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	rdsParams, err := params.CreateDBInput(
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
