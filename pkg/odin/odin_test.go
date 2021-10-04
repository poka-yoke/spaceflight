package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/internal/test"
	"github.com/poka-yoke/spaceflight/internal/test/mocks"
	"github.com/poka-yoke/spaceflight/pkg/odin"
)

type getLastSnapshotCase struct {
	test.Case
	name       string
	identifier string
	snapshots  []*rds.DBSnapshot
}

var getLastSnapshotCases = []getLastSnapshotCase{
	// Get snapshot id by instance id
	{
		Case: test.Case{
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
		Case: test.Case{
			Expected:      nil,
			ExpectedError: "no snapshot found for develop instance",
		},
		name:       "Get non-existing snapshot id by instance id",
		identifier: "develop",
		snapshots:  []*rds.DBSnapshot{},
	},
	// Get last snapshot id by instance id out of two
	{
		Case: test.Case{
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
	svc := mocks.NewRDSClient()
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

func getTime(original string) (parsed time.Time) {
	parsed, _ = time.Parse(
		odin.RFC8601,
		original,
	)
	return
}

var exampleSnapshot1DBID = aws.String("production-rds")
var exampleSnapshot1ID = aws.String("rds:production-2015-06-11")
var exampleSnapshot1Time = "2015-06-11T22:00:00+00:00"
var exampleSnapshot1 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: exampleSnapshot1DBID,
	DBSnapshotIdentifier: exampleSnapshot1ID,
	MasterUsername:       aws.String("owner"),
	SnapshotCreateTime:   aws.Time(getTime(exampleSnapshot1Time)),
	Status:               aws.String("available"),
}

var exampleSnapshot2DBID = aws.String("develop-rds")
var exampleSnapshot2ID = aws.String("rds:develop-2016-06-11")
var exampleSnapshot2Time = "2016-06-11T22:00:00+00:00"
var exampleSnapshot2 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: exampleSnapshot2DBID,
	DBSnapshotIdentifier: exampleSnapshot2ID,
	MasterUsername:       aws.String("owner"),
	SnapshotCreateTime:   aws.Time(getTime(exampleSnapshot2Time)),
	Status:               aws.String("available"),
}

var exampleSnapshot3DBID = aws.String("develop-rds")
var exampleSnapshot3ID = aws.String("rds:develop-2017-06-11")
var exampleSnapshot3Time = "2017-06-11T22:00:00+00:00"
var exampleSnapshot3 = &rds.DBSnapshot{
	DBInstanceIdentifier: exampleSnapshot3DBID,
	DBSnapshotIdentifier: exampleSnapshot3ID,
	SnapshotCreateTime:   aws.Time(getTime(exampleSnapshot3Time)),
}

type listSnapshotsCase struct {
	test.Case
	name       string
	snapshots  []*rds.DBSnapshot
	instanceID string
}

var listSnapshotsCases = []listSnapshotsCase{
	// No snapshots for any instance
	{
		Case: test.Case{
			Expected:      []*rds.DBSnapshot{},
			ExpectedError: "",
		},
		name:       "No snapshots for any instance",
		snapshots:  []*rds.DBSnapshot{},
		instanceID: "",
	},
	// One instance, one snapshot
	{
		Case: test.Case{
			Expected:      []*rds.DBSnapshot{exampleSnapshot1},
			ExpectedError: "",
		},
		name:       "One instance one snapshot",
		snapshots:  []*rds.DBSnapshot{exampleSnapshot1},
		instanceID: "",
	},
	// Two instances, two snapshots
	{
		Case: test.Case{
			Expected: []*rds.DBSnapshot{
				exampleSnapshot2,
				exampleSnapshot1,
			},
			ExpectedError: "",
		},
		name: "Two instances two snapshots",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot1,
			exampleSnapshot2,
		},
		instanceID: "",
	},
	// Instance selection
	{
		Case: test.Case{
			Expected: []*rds.DBSnapshot{
				exampleSnapshot2,
			},
			ExpectedError: "",
		},
		name: "Two instances two snapshots, one selected",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot1,
			exampleSnapshot2,
		},
		instanceID: "develop-rds",
	},
}

func TestListSnapshots(t *testing.T) {
	for _, test := range listSnapshotsCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc := mocks.NewRDSClient()
				svc.AddSnapshots(test.snapshots)
				actual, err := odin.ListSnapshots(
					test.instanceID,
					svc,
				)
				test.Check(actual, err, t)
			},
		)
	}
}
