package odin

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// Duration specified time to wait for instance to be available.
var Duration = time.Duration(5) * time.Second

// Init initializes connection to AWS API
func Init() rdsiface.RDSAPI {
	region := "us-east-1"
	sess := session.New(&aws.Config{Region: aws.String(region)})
	return rds.New(sess)
}

// CreateDBParams represents CreateDBInstance parameters.
type CreateDBParams struct {
	DBInstanceType      string
	DBUser              string
	DBPassword          string
	DBSubnetGroupName   string
	VpcSecurityGroupIds []string
	Size                int64

	OriginalInstanceName string
	Restore              bool
}

// GetRestoreDBInstanceFromDBSnapshotInput method creates a new
// RestoreDBInstanceFromDBSnapshotInput from provided CreateDBParams and
// rds.DBSnapshot.
func (params CreateDBParams) GetRestoreDBInstanceFromDBSnapshotInput(
	identifier string,
	svc rdsiface.RDSAPI,
) (out *rds.RestoreDBInstanceFromDBSnapshotInput, err error) {
	if params.OriginalInstanceName == "" {
		err = fmt.Errorf("Original Instance Name was empty")
		return
	}
	snapshot, err := GetLastSnapshot(params.OriginalInstanceName, svc)
	if err != nil {
		err = fmt.Errorf(
			"Couldn't find snapshot for %s instance",
			params.OriginalInstanceName,
		)
		return
	}
	out = &rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceClass:      &params.DBInstanceType,
		DBInstanceIdentifier: &identifier,
		DBSnapshotIdentifier: snapshot.DBSnapshotIdentifier,
		DBSubnetGroupName:    &params.DBSubnetGroupName,
		Engine:               aws.String("postgres"),
	}
	return
}

// GetCreateDBInstanceInput method creates a new CreateDBInstanceInput from provided
// CreateDBParams and rds.DBSnapshot.
func (params CreateDBParams) GetCreateDBInstanceInput(
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
func (params CreateDBParams) GetModifyDBInstanceInput(
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

func applySnapshotParams(identifier string, in *rds.CreateDBInstanceInput, svc rdsiface.RDSAPI) (out *rds.CreateDBInstanceInput, err error) {
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

// CreateDBInstance creates a new RDS database instance. If a vpcid is
// specified the security group will be in that VPC.
func CreateDBInstance(
	instanceName string,
	params CreateDBParams,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	var instance *rds.DBInstance
	if params.Restore {
		instance, err = getInstanceRestore(instanceName, params, svc)
		if err != nil {
			return
		}
	} else {
		if params.OriginalInstanceName == "" {
			instance, err = getInstanceCreate(instanceName, params, svc)
			if err != nil {
				return
			}
		} else {
			instance, err = getInstanceClone(instanceName, params, svc)
			if err != nil {
				return
			}
		}
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

func modifyInstance(instanceName string, params CreateDBParams, svc rdsiface.RDSAPI) (err error) {
	rdsParams := params.GetModifyDBInstanceInput(
		instanceName,
		svc,
	)
	if err = rdsParams.Validate(); err != nil {
		err = fmt.Errorf(
			"DB instance parameters failed to validate: %s",
			err,
		)
		return
	}
	_, err = svc.ModifyDBInstance(rdsParams)
	return
}

func getInstanceRestore(instanceName string, params CreateDBParams, svc rdsiface.RDSAPI) (instance *rds.DBInstance, err error) {
	var res *rds.RestoreDBInstanceFromDBSnapshotOutput
	rdsParams, err := params.GetRestoreDBInstanceFromDBSnapshotInput(
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

func getInstanceCreate(instanceName string, params CreateDBParams, svc rdsiface.RDSAPI) (instance *rds.DBInstance, err error) {
	rdsParams := params.GetCreateDBInstanceInput(
		instanceName,
		svc,
	)
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

func getInstanceClone(instanceName string, params CreateDBParams, svc rdsiface.RDSAPI) (instance *rds.DBInstance, err error) {
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

// GetLastSnapshot queries AWS looking for a Snapshot ID, depending on
// an instance ID.
func GetLastSnapshot(
	identifier string,
	svc rdsiface.RDSAPI,
) (result *rds.DBSnapshot, err error) {
	params := &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: &identifier,
	}
	results, err := svc.DescribeDBSnapshots(params)
	if err != nil || len(results.DBSnapshots) == 0 {
		err = fmt.Errorf("There are no Snapshots for %s", identifier)
		return
	}
	return results.DBSnapshots[0], nil
}
