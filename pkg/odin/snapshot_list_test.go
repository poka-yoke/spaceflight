package odin_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/Devex/spaceflight/pkg/odin"
)

var exampleSnapshot1Type = aws.String("db.m1.medium")
var exampleSnapshot1DBID = aws.String("production-rds")
var exampleSnapshot1ID = aws.String("rds:production-2015-06-11")
var exampleSnapshot1Time = "2015-06-11T22:00:00+00:00"
var exampleSnapshot2DBID = aws.String("develop-rds")
var exampleSnapshot2ID = aws.String("rds:develop-2016-06-11")
var exampleSnapshot2Time = "2016-06-11T22:00:00+00:00"
var exampleSnapshot3DBID = aws.String("develop-rds")
var exampleSnapshot3ID = aws.String("rds:develop-2017-06-11")
var exampleSnapshot3Time = "2017-06-11T22:00:00+00:00"
var exampleSnapshot4DBID = aws.String("develop-rds")
var exampleSnapshot4ID = aws.String("rds:develop-2017-07-11")

func getTime(original string) (parsed time.Time) {
	parsed, _ = time.Parse(
		odin.RFC8601,
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

var exampleSnapshot3 = &rds.DBSnapshot{
	DBInstanceIdentifier: exampleSnapshot3DBID,
	DBSnapshotIdentifier: exampleSnapshot3ID,
	SnapshotCreateTime:   aws.Time(getTime(exampleSnapshot3Time)),
}

var exampleSnapshot4 = &rds.DBSnapshot{
	DBInstanceIdentifier: exampleSnapshot4DBID,
	DBSnapshotIdentifier: exampleSnapshot4ID,
}

var exampleSnapshot1Out = fmt.Sprintf(
	"%v %v %v %v\n",
	*exampleSnapshot1ID,
	*exampleSnapshot1DBID,
	exampleSnapshot1Time,
	"available",
)

var exampleSnapshot2Out = fmt.Sprintf(
	"%v %v %v %v\n",
	*exampleSnapshot2ID,
	*exampleSnapshot2DBID,
	exampleSnapshot2Time,
	"available",
)

var printSnapshotsCases = []listSnapshotsCase{
	// No snapshots
	{
		testCase: testCase{
			expected:      "",
			expectedError: "",
		},
		name:      "No snapshots",
		snapshots: []*rds.DBSnapshot{},
	},
	// One snapshot
	{
		testCase: testCase{
			expected:      exampleSnapshot1Out,
			expectedError: "",
		},
		name:      "One snapshot",
		snapshots: []*rds.DBSnapshot{exampleSnapshot1},
	},
	// Two snapshots
	{
		testCase: testCase{
			expected: strings.Join(
				[]string{
					exampleSnapshot2Out,
					exampleSnapshot1Out,
				},
				"",
			),
			expectedError: "",
		},
		name: "Two snapshots",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot2,
			exampleSnapshot1,
		},
	},
}

func TestPrintSnapshot(t *testing.T) {
	for _, test := range printSnapshotsCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				actual := odin.PrintSnapshots(
					test.snapshots,
				)
				test.check(actual, nil, t)
			},
		)
	}
}

type listSnapshotsCase struct {
	testCase
	name       string
	snapshots  []*rds.DBSnapshot
	instanceID string
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
		name:       "One instance one snapshot",
		snapshots:  []*rds.DBSnapshot{exampleSnapshot1},
		instanceID: "",
	},
	// Two instances, two snapshots
	{
		testCase: testCase{
			expected: []*rds.DBSnapshot{
				exampleSnapshot2,
				exampleSnapshot1,
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
	// Instance selection
	{
		testCase: testCase{
			expected: []*rds.DBSnapshot{
				exampleSnapshot2,
			},
			expectedError: "",
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
				svc := newMockRDSClient()
				svc.AddSnapshots(test.snapshots)
				actual, err := odin.ListSnapshots(
					test.instanceID,
					svc,
				)
				test.check(actual, err, t)
			},
		)
	}
}
