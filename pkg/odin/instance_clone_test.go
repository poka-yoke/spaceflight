package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/Devex/spaceflight/pkg/odin"
)

type cloneInstanceCase struct {
	testCase
	name         string
	identifier   string
	instanceType string
	password     string
	user         string
	size         int64
	from         string
	snapshots    []*rds.DBSnapshot
}

var cloneInstanceCases = []cloneInstanceCase{
	// Uses snapshot to copy from
	{
		testCase: testCase{
			expected:      "test1.0.us-east-1.rds.amazonaws.com",
			expectedError: "",
		},
		name:         "Uses snapshot to copy from",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
		size:         6144,
		from:         "production-rds",
		snapshots:    []*rds.DBSnapshot{exampleSnapshot1},
	},
	// Uses non existing snapshot to copy from
	{
		testCase: testCase{
			expected:      "",
			expectedError: "No snapshot found for develop instance",
		},
		name:         "Uses non existing snapshot to copy from",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
		size:         6144,
		from:         "develop",
		snapshots:    []*rds.DBSnapshot{exampleSnapshot1},
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
					svc.AddSnapshots(test.snapshots)
				}
				createParams := odin.CreateParams{
					InstanceType: test.instanceType,
					User:         test.user,
					Password:     test.password,
					Size:         test.size,
				}
				params := odin.CloneParams{
					CreateParams:         createParams,
					OriginalInstanceName: test.from,
				}
				actual, err := odin.CloneInstance(
					test.identifier,
					params,
					svc,
				)
				test.check(actual, err, t)
			},
		)
	}
}
