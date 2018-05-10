package odin_test

import (
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// trimLast trims last character from a given string and returns the trimmed
// one.
func trimLast(in string) string {
	if in == "" {
		return in
	}
	return in[:len(in)-1]
}

type mockRDSClient struct {
	rdsiface.RDSAPI
	dbInstances []*rds.DBInstance
	dbSnapshots []*rds.DBSnapshot
}

// TakeFinalSnapshot emulates taking the final snapshot, if specified,
// validating params.
func (m *mockRDSClient) TakeFinalSnapshot(
	params *rds.DeleteDBInstanceInput,
) (err error) {
	if params.SkipFinalSnapshot == nil {
		params.SkipFinalSnapshot = aws.Bool(false)
	}
	if !*params.SkipFinalSnapshot {
		if params.FinalDBSnapshotIdentifier == nil {
			err = fmt.Errorf("Final Snapshot ID not specified")
			return
		}
		if *params.FinalDBSnapshotIdentifier == "" {
			err = fmt.Errorf("Final Snapshot ID not specified")
			return
		}
		identifier := params.DBInstanceIdentifier
		snapshotID := params.FinalDBSnapshotIdentifier
		m.AddSnapshots([]*rds.DBSnapshot{
			{
				DBInstanceIdentifier: identifier,
				DBSnapshotIdentifier: snapshotID,
				Status:               aws.String("available"),
			},
		})
	}
	return
}

// DeleteDBInstance mocks rds.DeleteDBInstance.
func (m *mockRDSClient) DeleteDBInstance(
	params *rds.DeleteDBInstanceInput,
) (
	result *rds.DeleteDBInstanceOutput,
	err error,
) {
	if err = params.Validate(); err != nil {
		return
	}
	_, instance, err := m.findInstance(*params.DBInstanceIdentifier)
	if err != nil {
		return
	}
	if err = m.TakeFinalSnapshot(params); err != nil {
		return
	}
	instance.DBInstanceStatus = aws.String("deleting")
	result = &rds.DeleteDBInstanceOutput{
		DBInstance: instance,
	}
	return
}

// FindSnapshotInstance return index and snapshot in mockRDSClient.dbSnapshots
// for a specific instance id.
func (m mockRDSClient) FindSnapshotInstance(instanceID string) (
	index int,
	snapshot *rds.DBSnapshot,
	err error,
) {
	found := false
	for i, obj := range m.dbSnapshots {
		if *obj.DBInstanceIdentifier == instanceID {
			snapshot = obj
			index = i
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf(
			"No snapshot for instance %s",
			instanceID,
		)
	}
	return
}

// FindSnapshot return index and snapshot in mockRDSClient.dbSnapshots
// for a specific id.
func (m mockRDSClient) FindSnapshot(id string) (
	index int,
	snapshot *rds.DBSnapshot,
	err error,
) {
	found := false
	for i, obj := range m.dbSnapshots {
		if *obj.DBSnapshotIdentifier == id {
			snapshot = obj
			index = i
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf(
			"No such snapshot %s",
			id,
		)
	}
	return
}

// findInstance return index and instance in mockRDSClient.dbInstances
// for a specific id.
func (m mockRDSClient) findInstance(id string) (
	index int,
	instance *rds.DBInstance,
	err error,
) {
	found := false
	for i, obj := range m.dbInstances {
		if *obj.DBInstanceIdentifier == id {
			instance = obj
			index = i
			found = true
			break
		}
	}
	if !found {
		err = fmt.Errorf(
			"No such instance %s",
			id,
		)
	}
	return
}

// addInstances add a list of instances to the mock
func (m *mockRDSClient) addInstances(
	instances []*rds.DBInstance,
) {
	m.dbInstances = []*rds.DBInstance{}
	m.dbInstances = append(
		m.dbInstances,
		instances...,
	)
}

// AddSnapshots add a list of snapshots to the mock
func (m *mockRDSClient) AddSnapshots(
	snapshots []*rds.DBSnapshot,
) {
	m.dbSnapshots = []*rds.DBSnapshot{}
	m.dbSnapshots = append(
		m.dbSnapshots,
		snapshots...,
	)
}

// CreateDBSnapshot mocks rds.CreateDBSnapshot.
func (m *mockRDSClient) CreateDBSnapshot(
	params *rds.CreateDBSnapshotInput,
) (
	output *rds.CreateDBSnapshotOutput,
	err error,
) {
	instanceID := params.DBInstanceIdentifier
	_, instance, err := m.findInstance(*instanceID)
	if err != nil {
		return
	}
	if *instance.DBInstanceStatus != "available" {
		err = fmt.Errorf(
			"%s instance state is not available",
			*instanceID,
		)
	}
	id := params.DBSnapshotIdentifier
	_, snapshot, _ := m.FindSnapshot(*id)
	if snapshot != nil {
		err = fmt.Errorf(
			"Snapshot %s already exists",
			*id,
		)
		return
	}
	createdSnapshot := &rds.DBSnapshot{
		DBInstanceIdentifier: instanceID,
		DBSnapshotIdentifier: id,
	}
	output = &rds.CreateDBSnapshotOutput{
		DBSnapshot: createdSnapshot,
	}
	return
}

// DescribeDBSnapshots mocks rds.DescribeDBSnapshots.
func (m mockRDSClient) DescribeDBSnapshots(
	describeParams *rds.DescribeDBSnapshotsInput,
) (
	result *rds.DescribeDBSnapshotsOutput,
	err error,
) {
	var snapshots []*rds.DBSnapshot
	if describeParams.DBInstanceIdentifier != nil {
		snapshots = []*rds.DBSnapshot{}
		id := describeParams.DBInstanceIdentifier
		for _, snapshot := range m.dbSnapshots {
			if *snapshot.DBInstanceIdentifier == *id {
				snapshots = append(
					snapshots,
					snapshot,
				)
			}
		}
	} else {
		snapshots = m.dbSnapshots
	}
	result = &rds.DescribeDBSnapshotsOutput{
		DBSnapshots: snapshots,
	}
	return
}

// DescribeDBInstances mocks rds.DescribeDBInstances.
func (m *mockRDSClient) DescribeDBInstances(
	describeParams *rds.DescribeDBInstancesInput,
) (
	result *rds.DescribeDBInstancesOutput,
	err error,
) {
	id := describeParams.DBInstanceIdentifier
	index, instance, err := m.findInstance(*id)
	if err != nil {
		return
	}
	if *instance.DBInstanceStatus == "deleting" {
		m.dbInstances = append(
			m.dbInstances[:index],
			m.dbInstances[index+1:]...,
		)
	}
	if *instance.DBInstanceStatus == "creating" {
		az := *instance.AvailabilityZone
		region := trimLast(az)
		endpoint := fmt.Sprintf(
			"%s.0.%s.rds.amazonaws.com",
			*id,
			region,
		)
		port := int64(5432)
		instance.Endpoint = &rds.Endpoint{
			Address: &endpoint,
			Port:    &port,
		}
	}
	if *instance.DBInstanceStatus == "creating" ||
		*instance.DBInstanceStatus == "modifying" {
		*instance.DBInstanceStatus = "available"
		m.dbInstances[index] = instance
	}
	result = &rds.DescribeDBInstancesOutput{
		DBInstances: []*rds.DBInstance{
			instance,
		},
	}
	return
}

// CreateDBInstance mocks rds.CreateDBInstance.
func (m *mockRDSClient) CreateDBInstance(
	inputParams *rds.CreateDBInstanceInput,
) (
	result *rds.CreateDBInstanceOutput,
	err error,
) {
	if err = inputParams.Validate(); err != nil {
		return
	}
	az := aws.String("us-east-1c")
	if inputParams.AvailabilityZone != nil {
		az = inputParams.AvailabilityZone
	}
	if inputParams.MasterUsername == nil ||
		*inputParams.MasterUsername == "" {
		err = errors.New("Specify Master User")
		return
	}
	if inputParams.MasterUserPassword == nil ||
		*inputParams.MasterUserPassword == "" {
		err = errors.New("Specify Master User Password")
		return
	}
	if inputParams.AllocatedStorage == nil ||
		*inputParams.AllocatedStorage < 5 ||
		*inputParams.AllocatedStorage > 6144 {
		err = errors.New("Specify size between 5 and 6144")
		return
	}
	region := trimLast(*az)
	id := inputParams.DBInstanceIdentifier
	instance := rds.DBInstance{
		AllocatedStorage: inputParams.AllocatedStorage,
		AvailabilityZone: az,
		DBInstanceArn: aws.String(
			fmt.Sprintf(
				"arn:aws:rds:%s:0:db:%s",
				region,
				*id,
			),
		),
		DBInstanceIdentifier: id,
		DBInstanceStatus:     aws.String("creating"),
		Engine:               inputParams.Engine,
	}
	m.dbInstances = append(
		m.dbInstances,
		&instance,
	)
	result = &rds.CreateDBInstanceOutput{
		DBInstance: &instance,
	}
	return
}

// RestoreDBInstanceFromDBSnapshot mocks rds.RestoreDBInstanceFromDBSnapshot.
func (m *mockRDSClient) RestoreDBInstanceFromDBSnapshot(
	inputParams *rds.RestoreDBInstanceFromDBSnapshotInput,
) (
	result *rds.RestoreDBInstanceFromDBSnapshotOutput,
	err error,
) {
	if err = inputParams.Validate(); err != nil {
		return
	}
	az := aws.String("us-east-1c")
	if inputParams.AvailabilityZone != nil {
		az = inputParams.AvailabilityZone
	}
	region := trimLast(*az)
	id := inputParams.DBInstanceIdentifier
	instance := rds.DBInstance{
		AvailabilityZone: az,
		DBInstanceArn: aws.String(
			fmt.Sprintf(
				"arn:aws:rds:%s:0:db:%s",
				region,
				*id,
			),
		),
		DBInstanceIdentifier: id,
		DBInstanceStatus:     aws.String("creating"),
		Engine:               inputParams.Engine,
	}
	m.dbInstances = append(
		m.dbInstances,
		&instance,
	)
	result = &rds.RestoreDBInstanceFromDBSnapshotOutput{
		DBInstance: &instance,
	}
	return
}

// ModifyDBInstance mocks rds.ModifyDBInstance.
func (m *mockRDSClient) ModifyDBInstance(
	params *rds.ModifyDBInstanceInput,
) (
	out *rds.ModifyDBInstanceOutput,
	err error,
) {
	if err = params.Validate(); err != nil {
		return
	}
	index, instance, err := m.findInstance(*params.DBInstanceIdentifier)
	if err != nil {
		return
	}
	if *instance.DBInstanceStatus != "available" {
		err = fmt.Errorf(
			"%s instance state is not available",
			*params.DBInstanceIdentifier,
		)
	}
	if params.ApplyImmediately != nil &&
		*params.ApplyImmediately {
		instance.DBInstanceStatus = aws.String("modifying")
		if params.DBInstanceClass != nil &&
			params.DBInstanceClass != instance.DBInstanceClass {
			instance.DBInstanceClass = params.DBInstanceClass
		}
	}
	m.dbInstances[index] = instance
	out = &rds.ModifyDBInstanceOutput{
		DBInstance: instance,
	}
	return
}

// newMockRDSClient creates a mockRDSClient.
func newMockRDSClient() *mockRDSClient {
	return &mockRDSClient{
		dbInstances: []*rds.DBInstance{},
		dbSnapshots: []*rds.DBSnapshot{},
	}
}
