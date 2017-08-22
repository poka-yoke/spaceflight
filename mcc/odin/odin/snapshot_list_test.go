package odin_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

var exampleSnapshot1Type = aws.String("db.m1.medium")
var exampleSnapshot1DBID = aws.String("production-rds")
var exampleSnapshot1ID = aws.String("rds:production-2015-06-11")

var exampleSnapshot1 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: exampleSnapshot1DBID,
	DBSnapshotIdentifier: exampleSnapshot1ID,
	MasterUsername:       aws.String("owner"),
	Status:               aws.String("available"),
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
}

func TestListSnapshots(t *testing.T) {
	svc := newMockRDSClient()
	for _, test := range listSnapshotsCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				actual, err := odin.ListSnapshots(
					test.instanceID,
					svc,
				)
				test.check(actual, err, t)
			},
		)
	}
}
