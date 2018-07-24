package odin

import (
	"fmt"

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
	var instance *rds.DBInstance
	instance, err = doRestore(instanceName, params, svc)
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

func doRestore(
	instanceName string,
	params Instance,
	svc rdsiface.RDSAPI,
) (
	instance *rds.DBInstance,
	err error,
) {
	var res *rds.RestoreDBInstanceFromDBSnapshotOutput
	rdsParams, err := params.RestoreDBInput(
		instanceName,
		svc,
	)
	if err != nil {
		return
	}
	res, err = svc.RestoreDBInstanceFromDBSnapshot(rdsParams)
	if err != nil {
		return
	}
	return res.DBInstance, nil
}
