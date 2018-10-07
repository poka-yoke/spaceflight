package cmd

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/internal/test_case"
)

type createSnapshotsCase struct {
	testcase.TestCase
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
		TestCase: testcase.TestCase{
			Expected:      exampleSnapshot4,
			ExpectedError: "",
		},
		name:       "Create simple snapshot",
		instances:  []*rds.DBInstance{instance1},
		snapshots:  []*rds.DBSnapshot{},
		snapshotID: *exampleSnapshot4ID,
		instanceID: *exampleSnapshot4DBID,
	},
	// Snapshot already exists
	{
		TestCase: testcase.TestCase{
			Expected: nil,
			ExpectedError: fmt.Sprintf(
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
		TestCase: testcase.TestCase{
			Expected: nil,
			ExpectedError: fmt.Sprintf(
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
		TestCase: testcase.TestCase{
			Expected: nil,
			ExpectedError: fmt.Sprintf(
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
				actual, err := createSnapshot(
					test.instanceID,
					test.snapshotID,
					svc,
				)
				test.Check(actual, err, t)
			},
		)
	}
}
