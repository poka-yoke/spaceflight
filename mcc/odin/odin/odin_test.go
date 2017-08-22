package odin_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/go-test/deep"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

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

type getLastSnapshotCase struct {
	testCase
	name       string
	identifier string
	snapshots  []*rds.DBSnapshot
}

var getLastSnapshotCases = []getLastSnapshotCase{
	{
		testCase: testCase{
			expected:      exampleSnapshot1,
			expectedError: "",
		},
		name:       "Get snapshot id by instance id",
		identifier: "production",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot1,
		},
	},
	{
		testCase: testCase{
			expected:      nil,
			expectedError: "There are no Snapshots for production",
		},
		name:       "Get non-existant snapshot id by instance id",
		identifier: "production",
		snapshots:  []*rds.DBSnapshot{},
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
				test.check(actual, err, t)
			},
		)
	}
}
