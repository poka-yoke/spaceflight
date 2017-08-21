package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CloneParams represents CreateDBInstance parameters for cloning.
type CloneParams struct {
	InstanceType    string
	User            string
	Password        string
	SubnetGroupName string
	SecurityGroups  []string
	Size            int64

	OriginalInstanceName string
}

// GetCreateDBInput method creates a new CreateDBInstanceInput
// from provided CreateDBParams and rds.DBSnapshot.
func (params CloneParams) GetCreateDBInput(
	identifier string,
	svc rdsiface.RDSAPI,
) *rds.CreateDBInstanceInput {
	return &rds.CreateDBInstanceInput{
		AllocatedStorage:     &params.Size,
		DBInstanceIdentifier: &identifier,
		DBSubnetGroupName:    &params.SubnetGroupName,
		DBInstanceClass:      &params.InstanceType,
		DBSecurityGroups: []*string{
			aws.String("default"),
		},
		Engine:             aws.String("postgres"),
		EngineVersion:      aws.String("9.4.11"),
		MasterUsername:     &params.User,
		MasterUserPassword: &params.Password,
		Tags: []*rds.Tag{
			{
				Key:   aws.String("Name"),
				Value: &identifier,
			},
		},
	}
}

// GetModifyDBInput method creates a new ModifyDBInstanceInput
// from provided ModifyDBParams and rds.DBSnapshot.
func (params CloneParams) GetModifyDBInput(
	identifier string,
	svc rdsiface.RDSAPI,
) *rds.ModifyDBInstanceInput {
	SecurityGroups := []*string{}
	for _, sgid := range params.SecurityGroups {
		SecurityGroups = append(SecurityGroups, aws.String(sgid))
	}
	return &rds.ModifyDBInstanceInput{
		DBInstanceIdentifier: &identifier,
		VpcSecurityGroupIds:  SecurityGroups,
	}
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
	err = WaitForInstance(instance, svc)
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
			"Couldn't find snapshot for %s instance",
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
