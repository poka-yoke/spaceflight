package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/go-test/deep"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

type cloneInstanceCase struct {
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

func (t *cloneInstanceCase) expectingError(err error) bool {
	return t.expectedError != "" && err.Error() != t.expectedError
}

var cloneInstanceCases = []cloneInstanceCase{
	// Uses snapshot to copy from
	{
		name:          "Uses snapshot to copy from",
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
	// Uses non existing snapshot to copy from
	{
		name:          "Uses non existing snapshot to copy from",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		user:          "master",
		password:      "master",
		size:          6144,
		from:          "develop",
		expected:      "",
		expectedError: "Couldn't find snapshot for develop instance",
		snapshot:      exampleSnapshot1,
	},
}

func TestCloneInstance(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, test := range cloneInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				if test.from != "" {
					snapshot := test.snapshot
					id := test.snapshot.DBInstanceIdentifier
					snapshots := []*rds.DBSnapshot{snapshot}
					svc.dbSnapshots[*id] = snapshots
				}
				params := odin.CloneParams{
					InstanceType:         test.instanceType,
					User:                 test.user,
					Password:             test.password,
					Size:                 test.size,
					OriginalInstanceName: test.from,
				}
				actual, err := odin.CloneInstance(
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
