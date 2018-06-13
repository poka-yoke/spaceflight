package odin_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/Devex/spaceflight/pkg/odin"
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
			expected:      exampleSnapshot4,
			expectedError: "",
		},
		name:       "Create simple snapshot",
		instances:  []*rds.DBInstance{instance1},
		snapshots:  []*rds.DBSnapshot{},
		snapshotID: *exampleSnapshot4ID,
		instanceID: *exampleSnapshot4DBID,
	},
	// Snapshot already exists
	{
		testCase: testCase{
			expected: nil,
			expectedError: fmt.Sprintf(
				"Snapshot %s already exists",
				*exampleSnapshot4ID,
			),
		},
		name:       "Snapshot already exists",
		instances:  []*rds.DBInstance{instance1},
		snapshots:  []*rds.DBSnapshot{exampleSnapshot4},
		snapshotID: *exampleSnapshot4ID,
		instanceID: *exampleSnapshot4DBID,
	},
	// Invalid instance state
	{
		testCase: testCase{
			expected: nil,
			expectedError: fmt.Sprintf(
				"%s instance state is not available",
				*exampleSnapshot4DBID,
			),
		},
		name:       "Invalid instance state",
		instances:  []*rds.DBInstance{instance2},
		snapshots:  []*rds.DBSnapshot{},
		snapshotID: *exampleSnapshot4ID,
		instanceID: *exampleSnapshot4DBID,
	},
	// Non existing instance
	{
		testCase: testCase{
			expected: nil,
			expectedError: fmt.Sprintf(
				"No such instance %s",
				*exampleSnapshot4DBID,
			),
		},
		name:       "Non existing instance",
		instances:  []*rds.DBInstance{},
		snapshots:  []*rds.DBSnapshot{},
		snapshotID: *exampleSnapshot4ID,
		instanceID: *exampleSnapshot4DBID,
	},
}

func TestCreateSnapshot(t *testing.T) {
	for _, test := range createSnapshotCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc := newMockRDSClient()
				svc.addInstances(test.instances)
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
