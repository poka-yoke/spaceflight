package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

const (
	RFC8601 = "2006-01-02T15:04:05-07:00"
)

var exampleSnapshot1Type = aws.String("db.m1.medium")
var exampleSnapshot1DBID = aws.String("production-rds")
var exampleSnapshot1ID = aws.String("rds:production-2015-06-11")
var exampleSnapshot1Time = "2015-06-11T22:00:00+00:00"
var exampleSnapshot2DBID = aws.String("develop-rds")
var exampleSnapshot2ID = aws.String("rds:develop-2016-06-11")
var exampleSnapshot2Time = "2016-06-11T22:00:00+00:00"

func getTime(original string) (parsed time.Time) {
	parsed, _ = time.Parse(
		RFC8601,
		original,
	)
	return
}

var exampleSnapshot1 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: exampleSnapshot1DBID,
	DBSnapshotIdentifier: exampleSnapshot1ID,
	MasterUsername:       aws.String("owner"),
	SnapshotCreateTime:   aws.Time(getTime(exampleSnapshot1Time)),
	Status:               aws.String("available"),
}

var exampleSnapshot2 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: exampleSnapshot2DBID,
	DBSnapshotIdentifier: exampleSnapshot2ID,
	MasterUsername:       aws.String("owner"),
	SnapshotCreateTime:   aws.Time(getTime(exampleSnapshot2Time)),
	Status:               aws.String("available"),
}

type listSnapshotsCase struct {
	testCase
	name       string
	snapshots  []*rds.DBSnapshot
	instanceID string
}

// loadSnapshots method loads `snapshots` from current test case
// to the instance of the mocked RDSAPI, passed as argument.
// It returns nothing.
func (c *listSnapshotsCase) loadSnapshots(
	svc *mockRDSClient,
) {
	svc.dbSnapshots = map[string][]*rds.DBSnapshot{}
	// For each snapshot in the test case
	for _, snapshot := range c.snapshots {
		instanceID := snapshot.DBInstanceIdentifier
		var instanceSnapshots []*rds.DBSnapshot
		// if this snapshot's instance is not in the mock,
		if snapshotList, ok := svc.dbSnapshots[*instanceID]; !ok {
			// create the list of snapshots
			instanceSnapshots = make([]*rds.DBSnapshot, 0)
		} else {
			// or get the reference to the list
			instanceSnapshots = snapshotList
		}
		// append the snapshot to the instance's snapshot list
		instanceSnapshots = append(instanceSnapshots, snapshot)
		svc.dbSnapshots[*instanceID] = instanceSnapshots
	}
}

var listSnapshotsCases = []listSnapshotsCase{
	// No snapshots for any instance
	{
		testCase: testCase{
			expected:      []*rds.DBSnapshot{},
			expectedError: "",
		},
		name:       "No snapshots for any instance",
		snapshots:  []*rds.DBSnapshot{},
		instanceID: "",
	},
	// One instance, one snapshot
	{
		testCase: testCase{
			expected:      []*rds.DBSnapshot{exampleSnapshot1},
			expectedError: "",
		},
		name:       "One instance, one snapshot",
		snapshots:  []*rds.DBSnapshot{exampleSnapshot1},
		instanceID: "",
	},
	// Two instances, two snapshots
	{
		testCase: testCase{
			expected: []*rds.DBSnapshot{
				exampleSnapshot1,
				exampleSnapshot2,
			},
			expectedError: "",
		},
		name: "Two instances two snapshots",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot1,
			exampleSnapshot2,
		},
		instanceID: "",
	},
}

func TestListSnapshots(t *testing.T) {
	svc := newMockRDSClient()
	for _, test := range listSnapshotsCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				test.loadSnapshots(svc)
				actual, err := odin.ListSnapshots(
					test.instanceID,
					svc,
				)
				test.check(actual, err, t)
			},
		)
	}
}
