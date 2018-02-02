package odin

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// CreateSnapshot creates a Snapshot for specified instance.
func CreateSnapshot(
	instanceName string,
	snapshotName string,
	svc rdsiface.RDSAPI,
) (
	result *rds.DBSnapshot,
	err error,
) {
	output, err := svc.CreateDBSnapshot(
		&rds.CreateDBSnapshotInput{
			DBInstanceIdentifier: aws.String(instanceName),
			DBSnapshotIdentifier: aws.String(snapshotName),
		},
	)
	if err != nil {
		return
	}
	result = output.DBSnapshot
	return
}
