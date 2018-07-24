package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// Instance holds parameters needed for any operations related to
// instances and provides methods to obtain the AWS structures needed
// to perform them.
type Instance struct {
	Type                 string
	User                 string
	Password             string
	SubnetGroupName      string
	SecurityGroups       []string
	Size                 int64
	OriginalInstanceName string
	LastSnapshot         *rds.DBSnapshot
}

// CreateDBInput method creates a new CreateDBInstanceInput from
// rds.DBSnapshot.
func (params Instance) CreateDBInput(
	identifier string,
	svc rdsiface.RDSAPI,
) (result *rds.CreateDBInstanceInput, err error) {
	result = &rds.CreateDBInstanceInput{
		AllocatedStorage:     &params.Size,
		DBInstanceIdentifier: &identifier,
		DBSubnetGroupName:    &params.SubnetGroupName,
		DBInstanceClass:      &params.Type,
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
	if params.LastSnapshot != nil {
		result.AllocatedStorage = params.LastSnapshot.AllocatedStorage
		result.MasterUsername = params.LastSnapshot.MasterUsername
	}
	if err = result.Validate(); err != nil {
		err = fmt.Errorf(
			"DB instance parameters failed to validate: %s",
			err,
		)
		return nil, err
	}
	return result, nil
}

// ModifyDBInput method creates a new ModifyDBInstanceInput
// rds.DBSnapshot.
func (params Instance) ModifyDBInput(
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

// RestoreDBInput method creates a new
// RestoreDBInstanceFromDBSnapshotInput from provided CreateDBParams and
// rds.DBSnapshot.
func (p Instance) RestoreDBInput(
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
		DBInstanceClass:      &p.Type,
		DBInstanceIdentifier: &identifier,
		DBSnapshotIdentifier: snapshot.DBSnapshotIdentifier,
		DBSubnetGroupName:    &p.SubnetGroupName,
		Engine:               aws.String("postgres"),
	}
	return
}

// AddLastSnapshot adds the reference to the last available snapshot
// to the Instance structure, when it represents an already existing
// RDS instance.
func (params *Instance) AddLastSnapshot(
	identifier string,
	svc rdsiface.RDSAPI,
) (
	self *Instance,
	err error,
) {
	params.LastSnapshot, err = GetLastSnapshot(identifier, svc)
	if err != nil {
		err = fmt.Errorf(
			"No snapshot found for %s instance",
			identifier,
		)
		return params, err
	}
	return params, nil
}
