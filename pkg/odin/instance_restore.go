package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// RestoreParams represents instance restoring parameters.
type RestoreParams struct {
	InstanceType         string
	SubnetGroupName      string
	SecurityGroups       []string
	OriginalInstanceName string
}

// GetRestoreDBInput method creates a new
// RestoreDBInstanceFromDBSnapshotInput from provided CreateDBParams and
// rds.DBSnapshot.
func (p RestoreParams) GetRestoreDBInput(
	identifier string,
	svc rdsiface.RDSAPI,
) (
	out *rds.RestoreDBInstanceFromDBSnapshotInput,
	err error,
) {
	if p.OriginalInstanceName == "" {
		err = fmt.Errorf("Original Instance Name was empty")
		return
	}
	snapshot, err := GetLastSnapshot(p.OriginalInstanceName, svc)
	if err != nil {
		err = fmt.Errorf(
			"No snapshot found for %s instance",
			p.OriginalInstanceName,
		)
		return
	}
	out = &rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceClass:      &p.InstanceType,
		DBInstanceIdentifier: &identifier,
		DBSnapshotIdentifier: snapshot.DBSnapshotIdentifier,
		DBSubnetGroupName:    &p.SubnetGroupName,
		Engine:               aws.String("postgres"),
	}
	return
}

// GetModifyDBInput method creates a new ModifyDBInstanceInput
// from provided ModifyDBParams and rds.DBSnapshot.
func (p RestoreParams) GetModifyDBInput(
	identifier string,
	svc rdsiface.RDSAPI,
) *rds.ModifyDBInstanceInput {
	SecurityGroups := []*string{}
	for _, sgid := range p.SecurityGroups {
		SecurityGroups = append(SecurityGroups, aws.String(sgid))
	}
	return &rds.ModifyDBInstanceInput{
		DBInstanceIdentifier: &identifier,
		VpcSecurityGroupIds:  SecurityGroups,
	}
}

// RestoreInstance creates a new RDS database instance restoring
// from a snapshot.
func RestoreInstance(
	instanceName string,
	params RestoreParams,
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
	params RestoreParams,
	svc rdsiface.RDSAPI,
) (
	instance *rds.DBInstance,
	err error,
) {
	var res *rds.RestoreDBInstanceFromDBSnapshotOutput
	rdsParams, err := params.GetRestoreDBInput(
		instanceName,
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
	res, err = svc.RestoreDBInstanceFromDBSnapshot(rdsParams)
	if err != nil {
		return
	}
	instance = res.DBInstance
	return
}
