package odin

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
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
	if len(dbSnapshots) == 0 {
		err = errors.New("Snapshot not found")
		return
	}
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

type createDBInstanceCase struct {
	name               string
	identifier         string
	instanceType       string
	masterUserPassword string
	masterUser         string
	size               int64
	endpoint           string
	expectedError      string
}

var createDBInstanceCases = []createDBInstanceCase{
	{
		name:               "Creating simple instance",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		size:               5,
		endpoint:           "test1.0.us-east-1.rds.amazonaws.com",
		expectedError:      "",
	},
	{
		name:               "Fail because non-present user",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUserPassword: "master",
		size:               5,
		endpoint:           "",
		expectedError:      "Specify Master User",
	},
	{
		name:               "Fail because empty user",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "",
		masterUserPassword: "master",
		size:               5,
		endpoint:           "",
		expectedError:      "Specify Master User",
	},
	{
		name:          "Fail because non-present password",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		masterUser:    "master",
		size:          5,
		endpoint:      "",
		expectedError: "Specify Master User Password",
	},
	{
		name:               "Fail because empty password",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "",
		size:               5,
		endpoint:           "",
		expectedError:      "Specify Master User Password",
	},
	{
		name:               "Fail because non-present size",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		endpoint:           "",
		expectedError:      "Specify size between 5 and 6144",
	},
	{
		name:               "Fail because too small size",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		size:               4,
		endpoint:           "",
		expectedError:      "Specify size between 5 and 6144",
	},
	{
		name:               "Fail because too big size",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		size:               6145,
		endpoint:           "",
		expectedError:      "Specify size between 5 and 6144",
	},
}

func TestCreateDB(t *testing.T) {
	svc := &mockRDSClient{
		dbInstancesEndpoints: map[string]rds.Endpoint{},
		dbSnapshots:          map[string][]*rds.DBSnapshot{},
	}
	duration = time.Duration(0)
	for _, useCase := range createDBInstanceCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				params := CreateDBParams{
					DBInstanceType: useCase.instanceType,
					DBUser:         useCase.masterUser,
					DBPassword:     useCase.masterUserPassword,
					Size:           useCase.size,
				}
				endpoint, err := CreateDBInstance(
					useCase.identifier,
					params,
					svc,
				)
				if err != nil {
					if fmt.Sprintf("%s", err) != useCase.expectedError {
						t.Errorf(
							"Unexpected error %s",
							err,
						)
					}
				}
				if endpoint != useCase.endpoint {
					t.Errorf(
						"Unexpected output: %s should be %s",
						endpoint,
						useCase.endpoint,
					)
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

var exampleSnapshot1 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: aws.String("production"),
	DBSnapshotIdentifier: aws.String("rds:production-2015-06-11"),
	MasterUsername:       aws.String("owner"),
	Status:               aws.String("available"),
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
		expectedError: "Snapshot not found",
	},
}

func TestGetLastSnapshot(t *testing.T) {
	svc := &mockRDSClient{
		dbInstancesEndpoints: map[string]rds.Endpoint{},
		dbSnapshots:          map[string][]*rds.DBSnapshot{},
	}
	for _, useCase := range getLastSnapshotCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				svc.dbSnapshots[useCase.identifier] = useCase.snapshots
				snapshot, err := GetLastSnapshot(
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
				} else if fmt.Sprintf("%s", err) != useCase.expectedError {
					t.Errorf(
						"Unexpected error %s",
						err,
					)
				}
			},
		)
	}
}

type getCreateDBInstanceInputCase struct {
	name                          string
	identifier                    string
	createDBParams                CreateDBParams
	snapshot                      *rds.DBSnapshot
	expectedCreateDBInstanceInput *rds.CreateDBInstanceInput
	expectedError                 string
}

var getCreateDBInstanceInputCases = []getCreateDBInstanceInputCase{
	{
		name:       "Params without Snapshot",
		identifier: "production",
		createDBParams: CreateDBParams{
			DBInstanceType: "db.m1.medium",
			DBUser:         "owner",
			DBPassword:     "password",
			Size:           5,
		},
		snapshot: nil,
		expectedCreateDBInstanceInput: &rds.CreateDBInstanceInput{
			AllocatedStorage:     aws.Int64(5),
			DBInstanceIdentifier: aws.String("production"),
			DBInstanceClass:      aws.String("db.m1.medium"),
			MasterUsername:       aws.String("owner"),
			MasterUserPassword:   aws.String("password"),
			Engine:               aws.String("postgres"),
			EngineVersion:        aws.String("9.4.11"),
			DBSecurityGroups: []*string{
				aws.String("default"),
			},
			Tags: []*rds.Tag{
				{
					Key:   aws.String("Name"),
					Value: aws.String("production"),
				},
			},
		},
		expectedError: "",
	},
}

func TestGetCreateDBInstanceInput(t *testing.T) {
	svc := &mockRDSClient{
		dbInstancesEndpoints: map[string]rds.Endpoint{},
		dbSnapshots:          map[string][]*rds.DBSnapshot{},
	}
	for _, useCase := range getCreateDBInstanceInputCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				createDBInstanceInput, err := GetCreateDBInstanceInput(
					useCase.identifier,
					useCase.createDBParams,
					useCase.snapshot,
					svc,
				)
				if useCase.expectedError == "" {
					if err != nil {
						t.Errorf(
							"Unexpected error %s",
							err,
						)
					}
					if *createDBInstanceInput.DBInstanceIdentifier != *useCase.expectedCreateDBInstanceInput.DBInstanceIdentifier ||
						*createDBInstanceInput.MasterUsername != *useCase.expectedCreateDBInstanceInput.MasterUsername ||
						*createDBInstanceInput.MasterUserPassword != *useCase.expectedCreateDBInstanceInput.MasterUserPassword ||
						*createDBInstanceInput.AllocatedStorage != *useCase.expectedCreateDBInstanceInput.AllocatedStorage ||
						*createDBInstanceInput.DBInstanceClass != *useCase.expectedCreateDBInstanceInput.DBInstanceClass {
						t.Errorf(
							"Unexpected output: %s should be %s",
							createDBInstanceInput,
							useCase.expectedCreateDBInstanceInput,
						)
					}
				} else if fmt.Sprintf("%s", err) != useCase.expectedError {
					t.Errorf(
						"Unexpected error %s",
						err,
					)
				}
			},
		)
	}
}
