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
	delayChange bool,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	in := &rds.ModifyDBInstanceInput{
		DBInstanceIdentifier: aws.String(instanceName),
		DBInstanceClass:      aws.String(instanceType),
	}
	if !delayChange {
		in.ApplyImmediately = aws.Bool(true)
	}
	out, err := svc.ModifyDBInstance(in)
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
