package odin_test

import (
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
