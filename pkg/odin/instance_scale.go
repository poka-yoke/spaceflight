package odin

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// ScaleInstance scales an existing RDS database instance.
func ScaleInstance(
	params Instance,
	delayChange bool,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	rdsParams, err := params.ModifyDBInput(!delayChange, svc)
	if err != nil {
		return "", err
	}
	out, err := svc.ModifyDBInstance(rdsParams)
	if err != nil {
		return "", err
	}
	err = WaitForInstance(out.DBInstance, svc, "available")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"Instance %s is %s",
		*out.DBInstance.DBInstanceIdentifier,
		*out.DBInstance.DBInstanceClass,
	), nil
}
