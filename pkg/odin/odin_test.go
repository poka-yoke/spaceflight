package odin_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/go-test/deep"

	"github.com/Devex/spaceflight/pkg/odin"
)

type testCase struct {
	expected      interface{}
	expectedError string
}

func (tc *testCase) expectingError(err error) bool {
	return tc.expectedError != "" && err.Error() == tc.expectedError
}

func (tc *testCase) check(actual interface{}, err error, t *testing.T) {
	switch {
	case err != nil && !tc.expectingError(err):
		t.Errorf(
			"Unexpected error: %v",
			err,
		)
	case err != nil && tc.expectingError(err):
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

type getLastSnapshotCase struct {
	testCase
	name       string
	identifier string
	snapshots  []*rds.DBSnapshot
}

var getLastSnapshotCases = []getLastSnapshotCase{
	// Get snapshot id by instance id
	{
		testCase: testCase{
			expected:      exampleSnapshot1,
			expectedError: "",
		},
		name:       "Get snapshot id by instance id",
		identifier: "production-rds",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot1,
		},
	},
	// Get non-existing snapshot id by instance id
	{
		testCase: testCase{
			expected:      nil,
			expectedError: "No snapshot found for develop instance",
		},
		name:       "Get non-existing snapshot id by instance id",
		identifier: "develop",
		snapshots:  []*rds.DBSnapshot{},
	},
	// Get last snapshot id by instance id out of two
	{
		testCase: testCase{
			expected:      exampleSnapshot3,
			expectedError: "",
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
				test.check(actual, err, t)
			},
		)
	}
}
