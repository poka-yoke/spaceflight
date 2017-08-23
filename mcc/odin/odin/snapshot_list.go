package odin

import (
	"sort"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// ListSnapshots return a list of all DBSnapshots.
func ListSnapshots(
	instanceName string,
	svc rdsiface.RDSAPI,
) (
	result []*rds.DBSnapshot,
	err error,
) {
	input := &rds.DescribeDBSnapshotsInput{}
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
