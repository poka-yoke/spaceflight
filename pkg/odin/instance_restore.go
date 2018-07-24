package odin

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// RestoreInstance creates a new RDS database instance restoring
// from a snapshot.
func RestoreInstance(
	instanceName string,
	params Instance,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	var res *rds.RestoreDBInstanceFromDBSnapshotOutput
	rdsParams, err := params.RestoreDBInput(
		svc,
	)
	if err != nil {
		return "", err
	}
	res, err = svc.RestoreDBInstanceFromDBSnapshot(rdsParams)
	if err != nil {
		return "", err
	}
	err = WaitForInstance(res.DBInstance, svc, "available")
	if err != nil {
		return "", err
	}
	result = *res.DBInstance.Endpoint.Address
	err = modifyInstance(instanceName, params, svc)
	return
}
