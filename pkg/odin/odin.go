package odin

import (
	"fmt"
	"sort"
	"time"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

const (
	// RFC8601 is the date/time format used by AWS.
	RFC8601 = "2006-01-02T15:04:05-07:00"
)

// Duration specified time to wait for instance to be available.
var Duration = time.Duration(5) * time.Second

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

// WaitForInstance waits until instance's status is "available".
func WaitForInstance(
	instance *rds.DBInstance,
	svc rdsiface.RDSAPI,
	status string,
) (err error) {
	var res *rds.DescribeDBInstancesOutput
	for *instance.DBInstanceStatus != status {
		id := instance.DBInstanceIdentifier
		res, err = svc.DescribeDBInstances(
			&rds.DescribeDBInstancesInput{
				DBInstanceIdentifier: id,
			},
		)
		if err != nil {
			return
		}
		*instance = *res.DBInstances[0]
		// This is to avoid AWS API rate throttling.
		// Should use configurable exponential back-off
		time.Sleep(Duration)
	}
	return
}
