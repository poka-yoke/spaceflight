package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/Devex/spaceflight/pkg/odin"
)

type getRestoreDBInputCase struct {
	testCase
	name       string
	identifier string
	params     odin.RestoreParams
	snapshots  []*rds.DBSnapshot
}

var getRestoreDBInputCases = []getRestoreDBInputCase{
	// Params with Snapshot
	{
		testCase: testCase{
			expected: &rds.RestoreDBInstanceFromDBSnapshotInput{
				DBInstanceClass:      exampleSnapshot1Type,
				DBInstanceIdentifier: aws.String("develop"),
				DBSnapshotIdentifier: exampleSnapshot1ID,
				DBSubnetGroupName:    aws.String(""),
				Engine:               aws.String("postgres"),
			},
			expectedError: "",
		},
		name:       "Params with Snapshot",
		identifier: "develop",
		params: odin.RestoreParams{
			InstanceType:         "db.m1.medium",
			OriginalInstanceName: "production-rds",
		},
		snapshots: []*rds.DBSnapshot{exampleSnapshot1},
	},
	// Params with Snapshot without OriginalInstanceName
	{
		testCase: testCase{
			expected:      nil,
			expectedError: "Original Instance Name was empty",
		},
		name:       "Params with Snapshot without OriginalInstanceName",
		identifier: "production-rds",
		params: odin.RestoreParams{
			InstanceType: "db.m1.medium",
		},
		snapshots: []*rds.DBSnapshot{exampleSnapshot1},
	},
	// Params with non existing Snapshot
	{
		testCase: testCase{
			expected:      nil,
			expectedError: "No snapshot found for develop instance",
		},
		name:       "Params with non existing Snapshot",
		identifier: "production-rds",
		params: odin.RestoreParams{
			InstanceType:         "db.m1.medium",
			OriginalInstanceName: "develop",
		},
		snapshots: []*rds.DBSnapshot{exampleSnapshot1},
	},
}

func TestGetRestoreDBInput(t *testing.T) {
	svc := newMockRDSClient()
	for _, test := range getRestoreDBInputCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc.AddSnapshots(test.snapshots)
				params := test.params
				actual, err := params.GetRestoreDBInput(
					test.identifier,
					svc,
				)
				test.check(actual, err, t)
			},
		)
	}
}

var restoreInstanceCases = []cloneInstanceCase{
	// Uses snapshot to restore from
	{
		testCase: testCase{
			expected:      "test1.0.us-east-1.rds.amazonaws.com",
			expectedError: "",
		},
		name:         "Uses snapshot to restore from",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
		size:         6144,
		from:         "production-rds",
		snapshots:    []*rds.DBSnapshot{exampleSnapshot1},
	},
	// Uses non existing snapshot to restore from
	{
		testCase: testCase{
			expected:      "",
			expectedError: "No snapshot found for develop instance",
		},
		name:         "Uses non existing snapshot to restore from",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
		size:         6144,
		from:         "develop",
		snapshots:    []*rds.DBSnapshot{exampleSnapshot1},
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
					svc.AddSnapshots(test.snapshots)
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
				test.check(actual, err, t)
			},
		)
	}
}
