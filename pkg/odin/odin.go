package odin

import (
	"fmt"
	"sort"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

const (
	// RFC8601 is the date/time format used by AWS.
	RFC8601 = "2006-01-02T15:04:05-07:00"
)

// ModifiableParams is interface for params structs supporting
// DBInstance modification.
type ModifiableParams interface {
	ModifyDBInput(bool, rdsiface.RDSAPI) (*rds.ModifyDBInstanceInput, error)
}

// GetLastSnapshot queries AWS looking for a Snapshot ID, depending on
// an instance ID.
func GetLastSnapshot(
	id string,
	svc rdsiface.RDSAPI,
) (
	result *rds.DBSnapshot,
	err error,
) {
	results, err := ListSnapshots(id, svc)
	if err != nil || len(results) == 0 {
		err = fmt.Errorf("No snapshot found for %s instance", id)
		return
	}
	return results[0], nil
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
