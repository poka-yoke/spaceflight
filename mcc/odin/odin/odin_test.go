package odin_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

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
	endpoint, _ := m.dbInstancesEndpoints[*describeParams.DBInstanceIdentifier]
	result = &rds.DescribeDBInstancesOutput{
		DBInstances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: describeParams.DBInstanceIdentifier,
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
	endpoint := fmt.Sprintf(
		"%s.0.%s.rds.amazonaws.com",
		*inputParams.DBInstanceIdentifier,
		region,
	)
	port := int64(5432)
	m.dbInstancesEndpoints[*inputParams.DBInstanceIdentifier] = rds.Endpoint{
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
					inputParams.DBInstanceIdentifier,
				),
			),
			DBInstanceIdentifier: inputParams.DBInstanceIdentifier,
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
	endpoint := fmt.Sprintf(
		"%s.0.%s.rds.amazonaws.com",
		*inputParams.DBInstanceIdentifier,
		region,
	)
	port := int64(5432)
	m.dbInstancesEndpoints[*inputParams.DBInstanceIdentifier] = rds.Endpoint{
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
					inputParams.DBInstanceIdentifier,
				),
			),
			DBInstanceIdentifier: inputParams.DBInstanceIdentifier,
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

var exampleSnapshot1 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: aws.String("production"),
	DBSnapshotIdentifier: aws.String("rds:production-2015-06-11"),
	MasterUsername:       aws.String("owner"),
	Status:               aws.String("available"),
}

type createDBInstanceCase struct {
	name                 string
	identifier           string
	instanceType         string
	masterUserPassword   string
	masterUser           string
	size                 int64
	originalInstanceName string
	endpoint             string
	expectedError        string
	snapshot             *rds.DBSnapshot
}

var createDBInstanceCases = []createDBInstanceCase{
	// Creating simple instance
	{
		name:               "Creating simple instance",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		size:               5,
		endpoint:           "test1.0.us-east-1.rds.amazonaws.com",
		expectedError:      "",
		snapshot:           nil,
	},
	// Fail because empty user
	{
		name:               "Fail because empty user",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "",
		masterUserPassword: "master",
		size:               5,
		endpoint:           "",
		expectedError:      "Specify Master User",
		snapshot:           nil,
	},
	// Fail because empty password
	{
		name:               "Fail because empty password",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "",
		size:               5,
		endpoint:           "",
		expectedError:      "Specify Master User Password",
		snapshot:           nil,
	},
	// Fail because non-present size
	{
		name:               "Fail because non-present size",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		endpoint:           "",
		expectedError:      "Specify size between 5 and 6144",
		snapshot:           nil,
	},
	// Fail because too small size
	{
		name:               "Fail because too small size",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		endpoint:           "",
		expectedError:      "Specify size between 5 and 6144",
		snapshot:           nil,
	},
	// Fail because too big size
	{
		name:               "Fail because too big size",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		size:               6145,
		endpoint:           "",
		expectedError:      "Specify size between 5 and 6144",
		snapshot:           nil,
	},
	// Uses snapshot to copy from
	{
		name:                 "Uses snapshot to copy from",
		identifier:           "test1",
		instanceType:         "db.m1.small",
		masterUser:           "master",
		masterUserPassword:   "master",
		size:                 6144,
		originalInstanceName: "production",
		endpoint:             "test1.0.us-east-1.rds.amazonaws.com",
		expectedError:        "",
		snapshot:             exampleSnapshot1,
	},
	// Uses non existing snapshot to copy from
	{
		name:                 "Uses non existing snapshot to copy from",
		identifier:           "test1",
		instanceType:         "db.m1.small",
		masterUser:           "master",
		masterUserPassword:   "master",
		size:                 6144,
		originalInstanceName: "develop",
		endpoint:             "",
		expectedError:        "Couldn't find snapshot for develop instance",
		snapshot:             exampleSnapshot1,
	},
}

func TestCreateDB(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, useCase := range createDBInstanceCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				if useCase.originalInstanceName != "" {
					svc.dbSnapshots[*useCase.snapshot.DBInstanceIdentifier] = []*rds.DBSnapshot{
						useCase.snapshot,
					}
				}
				params := odin.CreateDBParams{
					DBInstanceType:       useCase.instanceType,
					DBUser:               useCase.masterUser,
					DBPassword:           useCase.masterUserPassword,
					Size:                 useCase.size,
					OriginalInstanceName: useCase.originalInstanceName,
				}
				endpoint, err := odin.CreateDBInstance(
					useCase.identifier,
					params,
					svc,
				)
				if err != nil {
					if err.Error() != useCase.expectedError {
						t.Errorf(
							"Unexpected error %s",
							err,
						)
					}
				} else if useCase.expectedError != "" {
					t.Errorf(
						"Expected error %s didn't happened",
						useCase.expectedError,
					)
				} else {
					if endpoint != useCase.endpoint {
						t.Errorf(
							"Unexpected output: %s should be %s",
							endpoint,
							useCase.endpoint,
						)
					}
				}
			},
		)
	}
}

type getLastSnapshotCase struct {
	name          string
	identifier    string
	snapshots     []*rds.DBSnapshot
	snapshot      *rds.DBSnapshot
	expectedError string
}

var getLastSnapshotCases = []getLastSnapshotCase{
	{
		name:       "Get snapshot id by instance id",
		identifier: "production",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot1,
		},
		snapshot:      exampleSnapshot1,
		expectedError: "",
	},
	{
		name:          "Get non-existant snapshot id by instance id",
		identifier:    "production",
		snapshots:     []*rds.DBSnapshot{},
		snapshot:      nil,
		expectedError: "There are no Snapshots for production",
	},
}

func TestGetLastSnapshot(t *testing.T) {
	svc := newMockRDSClient()
	for _, useCase := range getLastSnapshotCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				svc.dbSnapshots[useCase.identifier] = useCase.snapshots
				snapshot, err := odin.GetLastSnapshot(
					useCase.identifier,
					svc,
				)
				if useCase.expectedError == "" {
					if err != nil {
						t.Errorf(
							"Unexpected error %s",
							err,
						)
					}
					if *snapshot.DBInstanceIdentifier != *useCase.snapshot.DBInstanceIdentifier ||
						*snapshot.DBSnapshotIdentifier != *useCase.snapshot.DBSnapshotIdentifier {
						t.Errorf(
							"Unexpected output: %s should be %s",
							snapshot,
							useCase.snapshot,
						)
					}
				} else if err.Error() != useCase.expectedError {
					t.Errorf(
						"Unexpected error %s",
						err,
					)
				}
			},
		)
	}
}
