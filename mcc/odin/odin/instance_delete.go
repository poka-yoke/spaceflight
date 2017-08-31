package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// DeleteInstance deletes an existing RDS database instance.
func DeleteInstance(
	identifier string,
	svc rdsiface.RDSAPI,
) error {
	snapshotID := fmt.Sprintf(
		"%s-final",
		identifier,
	)
	out, err := svc.DeleteDBInstance(
		&rds.DeleteDBInstanceInput{
			DBInstanceIdentifier:      aws.String(identifier),
			FinalDBSnapshotIdentifier: aws.String(snapshotID),
		},
	)
	if err != nil {
		return err
	}
	WaitForInstance(out.DBInstance, svc, "deleted")
	return nil
}
