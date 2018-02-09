package odin

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// DeleteInstance deletes an existing RDS database instance.
func DeleteInstance(
	identifier string,
	snapshotID string,
	svc rdsiface.RDSAPI,
) error {
	params := &rds.DeleteDBInstanceInput{
		DBInstanceIdentifier: aws.String(identifier),
	}
	if snapshotID == "" {
		params.SkipFinalSnapshot = aws.Bool(true)
	} else {
		params.FinalDBSnapshotIdentifier = aws.String(snapshotID)
	}
	out, err := svc.DeleteDBInstance(
		params,
	)
	if err != nil {
		return err
	}
	WaitForInstance(out.DBInstance, svc, "deleted")
	return nil
}
