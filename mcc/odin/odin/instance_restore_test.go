package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/go-test/deep"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

type getRestoreDBInputCase struct {
	name          string
	identifier    string
	params        odin.RestoreParams
	snapshot      *rds.DBSnapshot
	expected      *rds.RestoreDBInstanceFromDBSnapshotInput
	expectedError string
}

func (t *getRestoreDBInputCase) expectingError(err error) bool {
	return t.expectedError != "" && err.Error() != t.expectedError
}

var getRestoreDBInputCases = []getRestoreDBInputCase{
	// Params with Snapshot
	{
		name:       "Params with Snapshot",
		identifier: "production-rds",
		params: odin.RestoreParams{
			InstanceType:         "db.m1.medium",
			OriginalInstanceName: "production",
		},
		snapshot: exampleSnapshot1,
		expected: &rds.RestoreDBInstanceFromDBSnapshotInput{
			DBInstanceClass:      aws.String("db.m1.medium"),
			DBInstanceIdentifier: aws.String("production-rds"),
			DBSnapshotIdentifier: exampleSnapshot1Id,
			DBSubnetGroupName:    aws.String(""),
			Engine:               aws.String("postgres"),
		},
		expectedError: "",
	},
	// Params with Snapshot without OriginalInstanceName
	{
		name:       "Params with Snapshot without OriginalInstanceName",
		identifier: "production-rds",
		params: odin.RestoreParams{
			InstanceType: "db.m1.medium",
		},
		snapshot:      exampleSnapshot1,
		expected:      nil,
		expectedError: "Original Instance Name was empty",
	},
	// Params with non existing Snapshot
	{
		name:       "Params with non existing Snapshot",
		identifier: "production-rds",
		params: odin.RestoreParams{
			InstanceType:         "db.m1.medium",
			OriginalInstanceName: "develop",
		},
		snapshot:      exampleSnapshot1,
		expected:      nil,
		expectedError: "No snapshot found for develop instance",
	},
}

func TestGetRestoreDBInput(t *testing.T) {
	svc := newMockRDSClient()
	for _, test := range getRestoreDBInputCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				if test.snapshot != nil {
					snapshot := test.snapshot
					id := *snapshot.DBInstanceIdentifier
					snapshots := []*rds.DBSnapshot{snapshot}
					svc.dbSnapshots[id] = snapshots
				}
				params := test.params
				actual, err := params.GetRestoreDBInput(
					test.identifier,
					svc,
				)
				switch {
				case err != nil && test.expectingError(err):
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

type restoreInstanceCase struct {
	name          string
	identifier    string
	instanceType  string
	password      string
	user          string
	size          int64
	from          string
	expected      string
	expectedError string
	snapshot      *rds.DBSnapshot
}

func (t *restoreInstanceCase) expectingError(err error) bool {
	return t.expectedError != "" && err.Error() != t.expectedError
}

var restoreInstanceCases = []restoreInstanceCase{
	// Uses snapshot to restore from
	{
		name:          "Uses snapshot to restore from",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		user:          "master",
		password:      "master",
		size:          6144,
		from:          "production",
		expected:      "test1.0.us-east-1.rds.amazonaws.com",
		expectedError: "",
		snapshot:      exampleSnapshot1,
	},
	// Uses non existing snapshot to restore from
	{
		name:          "Uses non existing snapshot to restore from",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		user:          "master",
		password:      "master",
		size:          6144,
		from:          "develop",
		expected:      "",
		expectedError: "No snapshot found for develop instance",
		snapshot:      exampleSnapshot1,
	},
}

func TestRestoreInstance(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, test := range restoreInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				if test.from != "" {
					snapshot := test.snapshot
					id := *snapshot.DBInstanceIdentifier
					snapshots := []*rds.DBSnapshot{snapshot}
					svc.dbSnapshots[id] = snapshots
				}
				params := odin.RestoreParams{
					InstanceType:         test.instanceType,
					OriginalInstanceName: test.from,
				}
				actual, err := odin.RestoreInstance(
					test.identifier,
					params,
					svc,
				)
				switch {
				case err != nil && test.expectingError(err):
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
