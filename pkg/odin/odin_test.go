package odin_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/internal/test_case"
	"github.com/poka-yoke/spaceflight/pkg/odin"
)

type getLastSnapshotCase struct {
	testcase.TestCase
	name       string
	identifier string
	snapshots  []*rds.DBSnapshot
}

var getLastSnapshotCases = []getLastSnapshotCase{
	// Get snapshot id by instance id
	{
		TestCase: testcase.TestCase{
			Expected:      exampleSnapshot1,
			ExpectedError: "",
		},
		name:       "Get snapshot id by instance id",
		identifier: "production-rds",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot1,
		},
	},
	// Get non-existing snapshot id by instance id
	{
		TestCase: testcase.TestCase{
			Expected:      nil,
			ExpectedError: "No snapshot found for develop instance",
		},
		name:       "Get non-existing snapshot id by instance id",
		identifier: "develop",
		snapshots:  []*rds.DBSnapshot{},
	},
	// Get last snapshot id by instance id out of two
	{
		TestCase: testcase.TestCase{
			Expected:      exampleSnapshot3,
			ExpectedError: "",
		},
		name:       "Get last snapshot id by instance id",
		identifier: "develop-rds",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot2,
			exampleSnapshot3,
		},
	},
}

func TestGetLastSnapshot(t *testing.T) {
	svc := newMockRDSClient()
	for _, test := range getLastSnapshotCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc.AddSnapshots(test.snapshots)
				actual, err := odin.GetLastSnapshot(
					test.identifier,
					svc,
				)
				test.Check(actual, err, t)
			},
		)
	}
}
