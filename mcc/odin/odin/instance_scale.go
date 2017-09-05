package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// ScaleInstance scales an existing RDS database instance.
func ScaleInstance(
	instanceName string,
	instanceType string,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	out, err := svc.ModifyDBInstance(
		&rds.ModifyDBInstanceInput{
			DBInstanceIdentifier: aws.String(instanceName),
			DBInstanceClass:      aws.String(instanceType),
		},
	)
	if err != nil {
		return
	}
	err = WaitForInstance(out.DBInstance, svc, "available")
	if err != nil {
		return
	}
	result = fmt.Sprintf(
		"Instance %s is %s",
		*out.DBInstance.DBInstanceIdentifier,
		*out.DBInstance.DBInstanceClass,
	)
	return
}
