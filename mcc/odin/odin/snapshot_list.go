package odin

import (
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
	result = output.DBSnapshots
	return
}
