package odin

import(
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

func getTime(original string) (parsed time.Time) {
	parsed, _ = time.Parse(
		RFC8601,
		original,
	)
	return
}

type mockRDSClient struct {
	rdsiface.RDSAPI
	dbInstances []*rds.DBInstance
	dbSnapshots []*rds.DBSnapshot
}

// newMockRDSClient creates a mockRDSClient.
func newMockRDSClient() *mockRDSClient {
	return &mockRDSClient{
		dbInstances: []*rds.DBInstance{},
		dbSnapshots: []*rds.DBSnapshot{},
	}
}

// DescribeDBSnapshots mocks rds.DescribeDBSnapshots.
func (m mockRDSClient) DescribeDBSnapshots(
	describeParams *rds.DescribeDBSnapshotsInput,
) (
	result *rds.DescribeDBSnapshotsOutput,
	err error,
) {
	var snapshots []*rds.DBSnapshot
	if describeParams.DBInstanceIdentifier != nil {
		snapshots = []*rds.DBSnapshot{}
		id := describeParams.DBInstanceIdentifier
		for _, snapshot := range m.dbSnapshots {
			if *snapshot.DBInstanceIdentifier == *id {
				snapshots = append(
					snapshots,
					snapshot,
				)
			}
		}
	} else {
		snapshots = m.dbSnapshots
	}
	result = &rds.DescribeDBSnapshotsOutput{
		DBSnapshots: snapshots,
	}
	return
}

// AddSnapshots adds a list of snapshots to the mock
func (m *mockRDSClient) AddSnapshots(
	snapshots []*rds.DBSnapshot,
) {
	m.dbSnapshots = append(
		m.dbSnapshots,
		snapshots...,
	)
}

func TestDeleteDBInput(t *testing.T) {
	tt := []struct{
		input Instance
		skipFinalSnapshot bool
		finalSnapshotID string
	}{
		{
			input: Instance{FinalSnapshotID: ""},
			skipFinalSnapshot: true,
		},
		{
			input: Instance{FinalSnapshotID: "123"},
			finalSnapshotID: "123",
		},
	}

	for _, tc := range tt {
		res, err := tc.input.DeleteDBInput()
		if err != nil {
			t.Errorf("It should not fail")
		}
		if res.SkipFinalSnapshot != nil &&
			*res.SkipFinalSnapshot != tc.skipFinalSnapshot {
			t.Errorf(
				"SkipFinalSnapshot should be %v. Received %v",
				tc.skipFinalSnapshot,
				*res.SkipFinalSnapshot,
			)
		}
		if !tc.skipFinalSnapshot &&
			*res.FinalDBSnapshotIdentifier != tc.finalSnapshotID {
			t.Errorf(
				"Wrong identifier. Expected %s, but got %s",
				*res.FinalDBSnapshotIdentifier,
				tc.finalSnapshotID,
			)
		}
	}
}

var exampleSnapshot1 = &rds.DBSnapshot{
	AllocatedStorage:     aws.Int64(10),
	AvailabilityZone:     aws.String("us-east-1c"),
	DBInstanceIdentifier: aws.String("production-rds"),
	DBSnapshotIdentifier: aws.String("rds:production-2015-06-11"),
	MasterUsername:       aws.String("owner"),
	SnapshotCreateTime:   aws.Time(getTime("2015-06-11T22:00:00+00:00")),
	Status:               aws.String("available"),
}

func TestRestoreDBInput(t *testing.T) {
	tt := []struct{
		input Instance
		dbSnapshotIdentifier string
		err error
	}{
		// Underspecified options
		{
			input: Instance{OriginalInstanceName: ""},
			err: fmt.Errorf("Original Instance Name was empty"),
		},
		// Non-existing snapshot
		{
			input: Instance{OriginalInstanceName: "im-not-here"},
			err: fmt.Errorf(
				"No snapshot found for %s instance",
				"im-not-here",
			),
			dbSnapshotIdentifier: "not-found",
		},
		// Existing snapshot
		{
			input: Instance{OriginalInstanceName: "production-rds"},
			err: nil,
			dbSnapshotIdentifier: "rds:production-2015-06-11",
		},
	}
	svc := newMockRDSClient()
	svc.AddSnapshots([]*rds.DBSnapshot{exampleSnapshot1})
	for _, tc := range tt {
		res, err := tc.input.RestoreDBInput(svc)
		if tc.err != nil &&
			err.Error() != tc.err.Error() {
			t.Errorf(
				"Expected %s, but got %s",
				err.Error(),
				tc.err.Error(),
			)
		}
		if err == nil {
			switch {
			case res == nil:
				t.Errorf("Response should not be nil")
			case res.DBSnapshotIdentifier == nil:
				t.Errorf("Snapshot Identifier should not be nil")
			case *res.DBSnapshotIdentifier != tc.dbSnapshotIdentifier:
				t.Errorf(
					"Expected Snapshot identifier to be %s, but got %s",
					*res.DBSnapshotIdentifier,
					tc.dbSnapshotIdentifier,
				)
			}
		}
	}
}
