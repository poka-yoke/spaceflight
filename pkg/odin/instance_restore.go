package odin

import (
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// RestoreInstance creates a new RDS database instance restoring
// from a snapshot.
func RestoreInstance(
	params Instance,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	rdsParams, err := params.RestoreDBInput(
		svc,
	)
	if err != nil {
		return "", err
	}
	res, err := svc.RestoreDBInstanceFromDBSnapshot(rdsParams)
	if err != nil {
		return "", err
	}
	err = WaitForInstance(res.DBInstance, svc, "available")
	if err != nil {
		return "", err
	}
	result = *res.DBInstance.Endpoint.Address
	err = ModifyInstance(params, svc)
	return
}
