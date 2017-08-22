package odin

import (
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// ListSnapshots return a list of all DBSnapshots.
func ListSnapshots(
	instanceName string,
	svc rdsiface.RDSAPI,
) (
	result []*rds.DBSnapshot,
	err error,
) {
	result = make([]*rds.DBSnapshot, 0)
	return
}
