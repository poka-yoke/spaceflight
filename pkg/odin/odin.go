package odin

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// Duration specified time to wait for instance to be available.
var Duration = time.Duration(5) * time.Second

// Init initializes connection to AWS API
func Init() rdsiface.RDSAPI {
	region := "us-east-1"
	sess := session.New(&aws.Config{Region: aws.String(region)})
	return rds.New(sess)
}

// ModifiableParams is interface for params structs supporting
// DBInstance modification.
type ModifiableParams interface {
	GetModifyDBInput(string, rdsiface.RDSAPI) *rds.ModifyDBInstanceInput
}

func modifyInstance(
	instanceName string,
	params ModifiableParams,
	svc rdsiface.RDSAPI,
) (err error) {
	rdsParams := params.GetModifyDBInput(
		instanceName,
		svc,
	)
	if err = rdsParams.Validate(); err != nil {
		err = fmt.Errorf(
			"DB instance parameters failed to validate: %s",
			err,
		)
		return
	}
	_, err = svc.ModifyDBInstance(rdsParams)
	return
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
