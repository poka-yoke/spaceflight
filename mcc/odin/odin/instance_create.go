package odin

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CreateParams represents CreateDBInstance parameters.
type CreateParams struct {
	InstanceType    string
	User            string
	Password        string
	SubnetGroupName string
	SecurityGroups  []string
	Size            int64
}

// GetCreateDBInstanceInput method creates a new CreateDBInstanceInput from provided
// CreateParams and rds.DBSnapshot.
func (params CreateParams) GetCreateDBInstanceInput(
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

// GetModifyDBInstanceInput method creates a new ModifyDBInstanceInput from provided
// ModifyDBParams and rds.DBSnapshot.
func (params CreateParams) GetModifyDBInstanceInput(
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

// CreateInstance creates a new RDS database instance. If a vpcid is
// specified the security group will be in that VPC.
func CreateInstance(
	instanceName string,
	params CreateParams,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	var instance *rds.DBInstance
	instance, err = doCreate(instanceName, params, svc)
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

func doCreate(instanceName string, params CreateParams, svc rdsiface.RDSAPI) (instance *rds.DBInstance, err error) {
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
