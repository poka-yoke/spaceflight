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

// CloneDBInput returns CreateDBInstanceInput for the instance with
// Snapshot belonging to a target instance.
func (i Instance) CloneDBInput(
	identifier string,
	svc rdsiface.RDSAPI,
) (
	result *rds.CreateDBInstanceInput,
	err error,
) {
	instance, err := i.addLastSnapshot(svc)
	if err != nil {
		return nil, err
	}
	i = *instance
	return i.CreateDBInput(identifier, svc)
}

// CreateDBInput returns CreateDBInstanceInput for the instance.
func (i Instance) CreateDBInput(
	identifier string,
	svc rdsiface.RDSAPI,
) (result *rds.CreateDBInstanceInput, err error) {
	result = &rds.CreateDBInstanceInput{
		AllocatedStorage:     &i.Size,
		DBInstanceIdentifier: &identifier,
		DBSubnetGroupName:    &i.SubnetGroupName,
		DBInstanceClass:      &i.Type,
		DBSecurityGroups: []*string{
			aws.String("default"),
		},
		Engine:             aws.String("postgres"),
		EngineVersion:      aws.String("9.4.11"),
		MasterUsername:     &i.User,
		MasterUserPassword: &i.Password,
		Tags: []*rds.Tag{
			{
				Key:   aws.String("Name"),
				Value: &identifier,
			},
		},
	}
	if i.LastSnapshot != nil {
		result.AllocatedStorage = i.LastSnapshot.AllocatedStorage
		result.MasterUsername = i.LastSnapshot.MasterUsername
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

// ModifyDBInput returns ModifyDBInstanceInput for the instance.
func (i Instance) ModifyDBInput(
	identifier string,
	svc rdsiface.RDSAPI,
) (result *rds.ModifyDBInstanceInput, err error) {
	SecurityGroups := []*string{}
	for _, sgid := range i.SecurityGroups {
		SecurityGroups = append(SecurityGroups, aws.String(sgid))
	}
	result = &rds.ModifyDBInstanceInput{
		DBInstanceIdentifier: &identifier,
		VpcSecurityGroupIds:  SecurityGroups,
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

// RestoreDBInput returns RestoreDBInstanceFromDBSnapshotInput for the
// instance.
func (i Instance) RestoreDBInput(
	identifier string,
	svc rdsiface.RDSAPI,
) (
	result *rds.RestoreDBInstanceFromDBSnapshotInput,
	err error,
) {
	if i.OriginalInstanceName == "" {
		err = fmt.Errorf("Original Instance Name was empty")
		return nil, err
	}
	instance, err := i.addLastSnapshot(svc)
	if err != nil {
		return nil, err
	}
	i = *instance
	result = &rds.RestoreDBInstanceFromDBSnapshotInput{
		DBInstanceClass:      &i.Type,
		DBInstanceIdentifier: &identifier,
		DBSnapshotIdentifier: i.LastSnapshot.DBSnapshotIdentifier,
		DBSubnetGroupName:    &i.SubnetGroupName,
		Engine:               aws.String("postgres"),
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

// addLastSnapshot adds the reference to the last available snapshot
// of the target instance to this instance.
func (i *Instance) addLastSnapshot(
	svc rdsiface.RDSAPI,
) (_ *Instance, err error) {
	i.LastSnapshot, err = GetLastSnapshot(i.OriginalInstanceName, svc)
	if err != nil {
		err = fmt.Errorf(
			"No snapshot found for %s instance",
			i.OriginalInstanceName,
		)
		return i, err
	}
	return i, nil
}
