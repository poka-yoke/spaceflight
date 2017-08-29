package odin

import (
	"fmt"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

const (
	// RFC8601 is the date/time format used by AWS.
	RFC8601 = "2006-01-02T15:04:05-07:00"
)

// PrintSnapshots return a string with the snapshot list output.
func PrintSnapshots(snapshots []*rds.DBSnapshot) string {
	lines := []string{}
	for _, snapshot := range snapshots {
		line := fmt.Sprintf(
			"%v %v %v %v\n",
			*snapshot.DBSnapshotIdentifier,
			*snapshot.DBInstanceIdentifier,
			(*snapshot.SnapshotCreateTime).Format(
				RFC8601,
			),
			*snapshot.Status,
		)
		lines = append(lines, line)
	}
	return strings.Join(lines, "")
}

// ListSnapshots return a list of all DBSnapshots.
func ListSnapshots(
	instanceName string,
	svc rdsiface.RDSAPI,
) (
	result []*rds.DBSnapshot,
	err error,
) {
	input := &rds.DescribeDBSnapshotsInput{}
	if instanceName != "" {
		input.SetDBInstanceIdentifier(instanceName)
	}
	output, err := svc.DescribeDBSnapshots(input)
	if err != nil {
		return
	}
	// Lambda function here implements reverse ordering based
	// on snapshot's creation time, so first element is newer.
	sort.Slice(
		output.DBSnapshots,
		func(i, j int) bool {
			return output.DBSnapshots[i].SnapshotCreateTime.After(
				*output.DBSnapshots[j].SnapshotCreateTime,
			)
		},
	)
	result = output.DBSnapshots
	return
}
