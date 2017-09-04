package odin_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

type createSnapshotsCase struct {
	testCase
	name       string
	instances  []*rds.DBInstance
	snapshots  []*rds.DBSnapshot
	snapshotID string
	instanceID string
}

var instance1 = &rds.DBInstance{
	DBInstanceIdentifier: exampleSnapshot3DBID,
	DBInstanceStatus:     aws.String("available"),
}

var instance2 = &rds.DBInstance{
	DBInstanceIdentifier: exampleSnapshot3DBID,
	DBInstanceStatus:     aws.String("creating"),
}

var createSnapshotCases = []createSnapshotsCase{
	// Create simple snapshot
	{
		testCase: testCase{
			expected:      exampleSnapshot3,
			expectedError: "",
		},
		name:       "Create simple snapshot",
		instances:  []*rds.DBInstance{instance1},
		snapshots:  []*rds.DBSnapshot{},
		snapshotID: *exampleSnapshot3ID,
		instanceID: *exampleSnapshot3DBID,
	},
	// Snapshot already exists
	{
		testCase: testCase{
			expected: nil,
			expectedError: fmt.Sprintf(
				"Snapshot %s already exists",
				*exampleSnapshot3ID,
			),
		},
		name:       "Snapshot already exists",
		instances:  []*rds.DBInstance{instance1},
		snapshots:  []*rds.DBSnapshot{exampleSnapshot3},
		snapshotID: *exampleSnapshot3ID,
		instanceID: *exampleSnapshot3DBID,
	},
	// Invalid instance state
	{
		testCase: testCase{
			expected: nil,
			expectedError: fmt.Sprintf(
				"%s instance state is not available",
				*exampleSnapshot3DBID,
			),
		},
		name:       "Invalid instance state",
		instances:  []*rds.DBInstance{instance2},
		snapshots:  []*rds.DBSnapshot{},
		snapshotID: *exampleSnapshot3ID,
		instanceID: *exampleSnapshot3DBID,
	},
	// Non existing instance
	{
		testCase: testCase{
			expected: nil,
			expectedError: fmt.Sprintf(
				"No such instance %s",
				*exampleSnapshot3DBID,
			),
		},
		name:       "Non existing instance",
		instances:  []*rds.DBInstance{},
		snapshots:  []*rds.DBSnapshot{},
		snapshotID: *exampleSnapshot3ID,
		instanceID: *exampleSnapshot3DBID,
	},
}

func TestCreateSnapshot(t *testing.T) {
	for _, test := range createSnapshotCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc := newMockRDSClient()
				svc.AddInstances(test.instances)
				svc.AddSnapshots(test.snapshots)
				actual, err := odin.CreateSnapshot(
					test.instanceID,
					test.snapshotID,
					svc,
				)
				test.check(actual, err, t)
			},
		)
	}
}
