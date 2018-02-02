package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CreateParams represents CreateInstance parameters.
type CreateParams struct {
	InstanceType    string
	User            string
	Password        string
	SubnetGroupName string
	SecurityGroups  []string
	Size            int64
}

// GetCreateDBInput method creates a new CreateDBInstanceInput from
// provided CreateParams and rds.DBSnapshot.
func (params CreateParams) GetCreateDBInput(
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

// GetModifyDBInput method creates a new ModifyDBInstanceInput from
// provided ModifyDBParams and rds.DBSnapshot.
func (params CreateParams) GetModifyDBInput(
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
	err = WaitForInstance(instance, svc, "available")
	if err != nil {
		return
	}
	result = *instance.Endpoint.Address
	err = modifyInstance(instanceName, params, svc)
	return
}

func doCreate(
	instanceName string,
	params CreateParams,
	svc rdsiface.RDSAPI,
) (
	instance *rds.DBInstance,
	err error,
) {
	rdsParams := params.GetCreateDBInput(
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
