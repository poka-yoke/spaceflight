package cmd

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/internal/test_case"
	"github.com/poka-yoke/spaceflight/pkg/odin"
)

type getRestoreDBInputCase struct {
	testcase.TestCase
	name       string
	identifier string
	params     odin.Instance
	snapshots  []*rds.DBSnapshot
}

var getRestoreDBInputCases = []getRestoreDBInputCase{
	// Params with Snapshot
	{
		TestCase: testcase.TestCase{
			Expected: &rds.RestoreDBInstanceFromDBSnapshotInput{
				DBInstanceClass:      exampleSnapshot1Type,
				DBInstanceIdentifier: aws.String("develop"),
				DBSnapshotIdentifier: exampleSnapshot1ID,
				DBSubnetGroupName:    aws.String(""),
				Engine:               aws.String("postgres"),
			},
			ExpectedError: "",
		},
		name: "Params with Snapshot",
		params: odin.Instance{
			Identifier:           "develop",
			Type:                 "db.m1.medium",
			OriginalInstanceName: "production-rds",
		},
		snapshots: []*rds.DBSnapshot{exampleSnapshot1},
	},
	// Params with Snapshot without OriginalInstanceName
	{
		TestCase: testcase.TestCase{
			Expected:      nil,
			ExpectedError: "Original Instance Name was empty",
		},
		name: "Params with Snapshot without OriginalInstanceName",
		params: odin.Instance{
			Identifier: "production-rds",
			Type:       "db.m1.medium",
		},
		snapshots: []*rds.DBSnapshot{exampleSnapshot1},
	},
	// Params with non existing Snapshot
	{
		TestCase: testcase.TestCase{
			Expected:      nil,
			ExpectedError: "No snapshot found for develop instance",
		},
		name: "Params with non existing Snapshot",
		params: odin.Instance{
			Identifier:           "production-rds",
			Type:                 "db.m1.medium",
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
				actual, err := params.RestoreDBInput(
					svc,
				)
				test.Check(actual, err, t)
			},
		)
	}
}

var restoreInstanceCases = []cloneInstanceCase{
	// Uses snapshot to restore from
	{
		TestCase: testcase.TestCase{
			Expected:      "test1.0.us-east-1.rds.amazonaws.com",
			ExpectedError: "",
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
		TestCase: testcase.TestCase{
			Expected:      "",
			ExpectedError: "No snapshot found for develop instance",
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
	for _, test := range restoreInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				if test.from != "" {
					svc.AddSnapshots(test.snapshots)
				}
				params := odin.Instance{
					Identifier:           test.identifier,
					Type:                 test.instanceType,
					OriginalInstanceName: test.from,
				}
				actual, err := restoreInstance(
					params,
					svc,
					time.Duration(0),
				)
				test.Check(actual, err, t)
			},
		)
	}
}
