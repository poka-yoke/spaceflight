package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CloneParams represents CreateDBInstance parameters for cloning.
type CloneParams struct {
	CreateParams

	OriginalInstanceName string
}

// CloneInstance creates a new RDS database instance, copying parameters
// from a snapshot. If a vpcid is specified the security group will be
// in that VPC.
func CloneInstance(
	instanceName string,
	params CloneParams,
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

func applySnapshotParams(
	identifier string,
	in *rds.CreateDBInstanceInput,
	svc rdsiface.RDSAPI,
) (
	out *rds.CreateDBInstanceInput,
	err error,
) {
	var snapshot *rds.DBSnapshot
	out = in
	snapshot, err = GetLastSnapshot(identifier, svc)
	if err != nil {
		err = fmt.Errorf(
			"No snapshot found for %s instance",
			identifier,
		)
		return
	}
	out.AllocatedStorage = snapshot.AllocatedStorage
	out.MasterUsername = snapshot.MasterUsername
	return
}

func doClone(
	instanceName string,
	params CloneParams,
	svc rdsiface.RDSAPI,
) (
	instance *rds.DBInstance,
	err error,
) {
	rdsParams := params.GetCreateDBInput(
		instanceName,
		svc,
	)
	rdsParams, err = applySnapshotParams(
		params.OriginalInstanceName,
		rdsParams,
		svc,
	)
	if err != nil {
		return
	}
	if err = rdsParams.Validate(); err != nil {
		err = fmt.Errorf(
			"DB instance parameters failed to validate: %s",
			err,
		)
		return
	}
	res, err := svc.CreateDBInstance(rdsParams)
	if err != nil {
		return
	}
	return res.DBInstance, nil
}
