package odin

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CloneParams represents CreateDBInstance parameters for cloning.
type CloneParams struct {
	DBInstanceType      string
	DBUser              string
	DBPassword          string
	DBSubnetGroupName   string
	VpcSecurityGroupIds []string
	Size                int64

	OriginalInstanceName string
}

// GetCreateDBInstanceInput method creates a new CreateDBInstanceInput from provided
// CreateDBParams and rds.DBSnapshot.
func (params CloneParams) GetCreateDBInstanceInput(
	identifier string,
	svc rdsiface.RDSAPI,
) *rds.CreateDBInstanceInput {
	return &rds.CreateDBInstanceInput{
		AllocatedStorage:     &params.Size,
		DBInstanceIdentifier: &identifier,
		DBSubnetGroupName:    &params.DBSubnetGroupName,
		DBInstanceClass:      &params.DBInstanceType,
		DBSecurityGroups: []*string{
			aws.String("default"),
		},
		Engine:             aws.String("postgres"),
		EngineVersion:      aws.String("9.4.11"),
		MasterUsername:     &params.DBUser,
		MasterUserPassword: &params.DBPassword,
		Tags: []*rds.Tag{
			{
				Key:   aws.String("Name"),
				Value: &identifier,
			},
		},
	}
}

// GetModifyDBInstanceInput method creates a new ModifyDBInstanceInput from provided
// ModifyDBParams and rds.DBSnapshot.
func (params CloneParams) GetModifyDBInstanceInput(
	identifier string,
	svc rdsiface.RDSAPI,
) *rds.ModifyDBInstanceInput {
	VpcSecurityGroupIds := []*string{}
	for _, sgid := range params.VpcSecurityGroupIds {
		VpcSecurityGroupIds = append(VpcSecurityGroupIds, aws.String(sgid))
	}
	return &rds.ModifyDBInstanceInput{
		DBInstanceIdentifier: &identifier,
		VpcSecurityGroupIds:  VpcSecurityGroupIds,
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
	var res *rds.DescribeDBInstancesOutput
	for *instance.DBInstanceStatus != "available" {
		res, err = svc.DescribeDBInstances(&rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		})
		if err != nil {
			return
		}
		instance = res.DBInstances[0]
		// This is to avoid AWS API rate throttling.
		// Should use configurable exponential back-off
		time.Sleep(Duration)
	}
	result = *instance.Endpoint.Address
	err = modifyInstance(instanceName, params, svc)
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
	rdsParams := params.GetCreateDBInstanceInput(
		instanceName,
		svc,
	)
	rdsParams, err = applySnapshotParams(params.OriginalInstanceName, rdsParams, svc)
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
