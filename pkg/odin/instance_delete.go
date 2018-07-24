package odin

import (
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// DeleteInstance deletes an existing RDS database instance.
func DeleteInstance(
	params Instance,
	svc rdsiface.RDSAPI,
) (err error) {
	rdsParams, err := params.DeleteDBInput(svc)
	if err != nil {
		return err
	}
	out, err := svc.DeleteDBInstance(
		rdsParams,
	)
	if err != nil {
		return err
	}
	err = WaitForInstance(out.DBInstance, svc, "deleted")
	return nil
}
