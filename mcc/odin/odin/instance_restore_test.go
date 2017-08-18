package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/go-test/deep"

	"github.com/Devex/spaceflight/mcc/odin/odin"
)

type getRestoreDBInstanceFromDBSnapshotInputCase struct {
	name                         string
	identifier                   string
	restoreParams                odin.RestoreParams
	snapshot                     *rds.DBSnapshot
	expectedRestoreSnapshotInput *rds.RestoreDBInstanceFromDBSnapshotInput
	expectedError                string
}

var getRestoreDBInstanceFromDBSnapshotInputCases = []getRestoreDBInstanceFromDBSnapshotInputCase{
	// Params with Snapshot
	{
		name:       "Params with Snapshot",
		identifier: "production-rds",
		restoreParams: odin.RestoreParams{
			InstanceType:         "db.m1.medium",
			OriginalInstanceName: "production",
		},
		snapshot: exampleSnapshot1,
		expectedRestoreSnapshotInput: &rds.RestoreDBInstanceFromDBSnapshotInput{
			DBInstanceClass:      aws.String("db.m1.medium"),
			DBInstanceIdentifier: aws.String("production-rds"),
			DBSnapshotIdentifier: aws.String("rds:production-2015-06-11"),
			DBSubnetGroupName:    aws.String(""),
			Engine:               aws.String("postgres"),
		},
		expectedError: "",
	},
	// Params with Snapshot without OriginalInstanceName
	{
		name:       "Params with Snapshot without OriginalInstanceName",
		identifier: "production-rds",
		restoreParams: odin.RestoreParams{
			InstanceType: "db.m1.medium",
		},
		snapshot:                     exampleSnapshot1,
		expectedRestoreSnapshotInput: nil,
		expectedError:                "Original Instance Name was empty",
	},
	// Params with non existing Snapshot
	{
		name:       "Params with non existing Snapshot",
		identifier: "production-rds",
		restoreParams: odin.RestoreParams{
			InstanceType:         "db.m1.medium",
			OriginalInstanceName: "develop",
		},
		snapshot:                     exampleSnapshot1,
		expectedRestoreSnapshotInput: nil,
		expectedError:                "Couldn't find snapshot for develop instance",
	},
}

func TestGetRestoreDBInstanceFromDBSnapshotInput(t *testing.T) {
	svc := newMockRDSClient()
	for _, useCase := range getRestoreDBInstanceFromDBSnapshotInputCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				if useCase.snapshot != nil {
					svc.dbSnapshots[*useCase.snapshot.DBInstanceIdentifier] = []*rds.DBSnapshot{useCase.snapshot}
				}
				restoreSnapshotInput, err := useCase.restoreParams.GetRestoreDBInstanceFromDBSnapshotInput(
					useCase.identifier,
					svc,
				)
				if err != nil {
					if useCase.expectedError == "" ||
						err.Error() != useCase.expectedError {
						t.Errorf(
							"Unexpected error happened: %v",
							err,
						)
					}
				} else {
					if diff := deep.Equal(
						restoreSnapshotInput,
						useCase.expectedRestoreSnapshotInput,
					); diff != nil {
						t.Errorf(
							"Unexpected output: %s should be %s",
							restoreSnapshotInput,
							useCase.expectedRestoreSnapshotInput,
						)
					}
				}
			},
		)
	}
}

type restoreInstanceCase struct {
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

var restoreInstanceCases = []restoreInstanceCase{
	// Uses snapshot to restore from
	{
		name:                 "Uses snapshot to restore from",
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
	// Uses non existing snapshot to restore from
	{
		name:                 "Uses non existing snapshot to restore from",
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

func TestRestoreInstance(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, useCase := range restoreInstanceCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				if useCase.originalInstanceName != "" {
					svc.dbSnapshots[*useCase.snapshot.DBInstanceIdentifier] = []*rds.DBSnapshot{
						useCase.snapshot,
					}
				}
				params := odin.RestoreParams{
					InstanceType:         useCase.instanceType,
					OriginalInstanceName: useCase.originalInstanceName,
				}
				endpoint, err := odin.RestoreInstance(
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
