package odin_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/go-test/deep"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

type mockRDSClient struct {
	rdsiface.RDSAPI
	dbInstancesEndpoints map[string]rds.Endpoint
	dbSnapshots          map[string][]*rds.DBSnapshot
}

// DescribeDBSnapshots mocks rds.DescribeDBSnapshots.
func (m mockRDSClient) DescribeDBSnapshots(
	describeParams *rds.DescribeDBSnapshotsInput,
) (
	result *rds.DescribeDBSnapshotsOutput,
	err error,
) {
	dbSnapshots := m.dbSnapshots[*describeParams.DBInstanceIdentifier]
	result = &rds.DescribeDBSnapshotsOutput{
		DBSnapshots: dbSnapshots,
	}
	return
}

// DescribeDBInstances mocks rds.DescribeDBInstances.
func (m mockRDSClient) DescribeDBInstances(
	describeParams *rds.DescribeDBInstancesInput,
) (
	result *rds.DescribeDBInstancesOutput,
	err error,
) {
	status := "available"
	id := describeParams.DBInstanceIdentifier
	endpoint, _ := m.dbInstancesEndpoints[*id]
	result = &rds.DescribeDBInstancesOutput{
		DBInstances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: id,
				DBInstanceStatus:     &status,
				Endpoint:             &endpoint,
			},
		},
	}
	return
}

// CreateDBInstance mocks rds.CreateDBInstance.
func (m mockRDSClient) CreateDBInstance(
	inputParams *rds.CreateDBInstanceInput,
) (
	result *rds.CreateDBInstanceOutput,
	err error,
) {
	if err = inputParams.Validate(); err != nil {
		return
	}
	az := "us-east-1c"
	if inputParams.AvailabilityZone != nil {
		az = *inputParams.AvailabilityZone
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
	region := az[:len(az)-1]
	id := inputParams.DBInstanceIdentifier
	endpoint := fmt.Sprintf(
		"%s.0.%s.rds.amazonaws.com",
		*id,
		region,
	)
	port := int64(5432)
	m.dbInstancesEndpoints[*id] = rds.Endpoint{
		Address: &endpoint,
		Port:    &port,
	}
	status := "creating"
	result = &rds.CreateDBInstanceOutput{
		DBInstance: &rds.DBInstance{
			AllocatedStorage: inputParams.AllocatedStorage,
			DBInstanceArn: aws.String(
				fmt.Sprintf(
					"arn:aws:rds:%s:0:db:%s",
					region,
					id,
				),
			),
			DBInstanceIdentifier: id,
			DBInstanceStatus:     &status,
			Engine:               inputParams.Engine,
		},
	}
	return
}

// RestoreDBInstanceFromDBSnapshot mocks rds.RestoreDBInstanceFromDBSnapshot.
func (m mockRDSClient) RestoreDBInstanceFromDBSnapshot(
	inputParams *rds.RestoreDBInstanceFromDBSnapshotInput,
) (
	result *rds.RestoreDBInstanceFromDBSnapshotOutput,
	err error,
) {
	if err = inputParams.Validate(); err != nil {
		return
	}
	az := "us-east-1c"
	if inputParams.AvailabilityZone != nil {
		az = *inputParams.AvailabilityZone
	}
	region := az[:len(az)-1]
	id := inputParams.DBInstanceIdentifier
	endpoint := fmt.Sprintf(
		"%s.0.%s.rds.amazonaws.com",
		*id,
		region,
	)
	port := int64(5432)
	m.dbInstancesEndpoints[*id] = rds.Endpoint{
		Address: &endpoint,
		Port:    &port,
	}
	status := "creating"
	result = &rds.RestoreDBInstanceFromDBSnapshotOutput{
		DBInstance: &rds.DBInstance{
			DBInstanceArn: aws.String(
				fmt.Sprintf(
					"arn:aws:rds:%s:0:db:%s",
					region,
					id,
				),
			),
			DBInstanceIdentifier: id,
			DBInstanceStatus:     &status,
			Engine:               inputParams.Engine,
		},
	}
	return
}

// ModifyDBInstance mocks rds.ModifyDBInstance.
func (m mockRDSClient) ModifyDBInstance(
	inputParams *rds.ModifyDBInstanceInput,
) (
	result *rds.ModifyDBInstanceOutput,
	err error,
) {
	if err = inputParams.Validate(); err != nil {
		return
	}
	result = &rds.ModifyDBInstanceOutput{
		DBInstance: &rds.DBInstance{
			DBInstanceIdentifier: inputParams.DBInstanceIdentifier,
		},
	}
	return
}

// newMockRDSClient creates a mockRDSClient.
func newMockRDSClient() *mockRDSClient {
	return &mockRDSClient{
		dbInstancesEndpoints: map[string]rds.Endpoint{},
		dbSnapshots:          map[string][]*rds.DBSnapshot{},
	}
}

type testCase struct {
	expected      interface{}
	expectedError string
}

func (tc *testCase) expectingError(err error) bool {
	return tc.expectedError != "" && err.Error() != tc.expectedError
}

func (tc *testCase) check(actual interface{}, err error, t *testing.T) {
	switch {
	case err != nil && tc.expectingError(err):
		t.Errorf(
			"Unexpected error: %v",
			err,
		)
	case err == nil && tc.expectedError != "":
		t.Errorf(
			"Expected error: %v missing",
			tc.expectedError,
		)
	case err == nil:
		if diff := deep.Equal(
			actual,
			tc.expected,
		); diff != nil {
			t.Errorf(
				"Unexpected output: %s",
				diff,
			)
		}
	}
}

var exampleSnapshot1Type = aws.String("db.m1.medium")
var exampleSnapshot1DBID = aws.String("production-rds")
var exampleSnapshot1ID = aws.String("rds:production-2015-06-11")

var exampleSnapshot1 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: exampleSnapshot1DBID,
	DBSnapshotIdentifier: exampleSnapshot1ID,
	MasterUsername:       aws.String("owner"),
	Status:               aws.String("available"),
}

type getLastSnapshotCase struct {
	name          string
	identifier    string
	snapshots     []*rds.DBSnapshot
	expected      *rds.DBSnapshot
	expectedError string
}

func (t *getLastSnapshotCase) expectingError(err error) bool {
	return t.expectedError != "" && err.Error() == t.expectedError
}

var getLastSnapshotCases = []getLastSnapshotCase{
	{
		name:       "Get snapshot id by instance id",
		identifier: "production",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot1,
		},
		expected:      exampleSnapshot1,
		expectedError: "",
	},
	{
		name:          "Get non-existant snapshot id by instance id",
		identifier:    "production",
		snapshots:     []*rds.DBSnapshot{},
		expected:      nil,
		expectedError: "There are no Snapshots for production",
	},
}

func TestGetLastSnapshot(t *testing.T) {
	svc := newMockRDSClient()
	for _, test := range getLastSnapshotCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				id := test.identifier
				svc.dbSnapshots[id] = test.snapshots
				actual, err := odin.GetLastSnapshot(
					id,
					svc,
				)
				switch {
				case err != nil && !test.expectingError(err):
					t.Errorf(
						"Unexpected error: %v",
						err,
					)
				case err == nil && test.expectedError != "":
					t.Errorf(
						"Expected error: %v missing",
						test.expectedError,
					)
				case err == nil:
					if diff := deep.Equal(
						actual,
						test.expected,
					); diff != nil {
						t.Errorf(
							"Unexpected output: %s",
							diff,
						)
					}
				}
			},
		)
	}
}
