package odin

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// DeleteInstance deletes an existing RDS database instance.
func DeleteInstance(
	identifier string,
	svc rdsiface.RDSAPI,
) error {
	_, err := svc.DeleteDBInstance(
		&rds.DeleteDBInstanceInput{
			DBInstanceIdentifier: aws.String(identifier),
		},
	)
	return err
}
