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

func applySnapshotParams(identifier string, in *rds.CreateDBInstanceInput, svc rdsiface.RDSAPI) (out *rds.CreateDBInstanceInput, err error) {
	var snapshot *rds.DBSnapshot
	out = in
	snapshot, err = GetLastSnapshot(identifier, svc)
	if err != nil {
		err = fmt.Errorf(
			"Couldn't find snapshot for %s instance",
			identifier,
		)
		return
	}
	out.AllocatedStorage = snapshot.AllocatedStorage
	out.MasterUsername = snapshot.MasterUsername
	return
}

// ModifiableParams is interface for params structs supporting DBInstance modification.
type ModifiableParams interface {
	GetModifyDBInstanceInput(string, rdsiface.RDSAPI) *rds.ModifyDBInstanceInput
}

func modifyInstance(instanceName string, params ModifiableParams, svc rdsiface.RDSAPI) (err error) {
	rdsParams := params.GetModifyDBInstanceInput(
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
	identifier string,
	svc rdsiface.RDSAPI,
) (result *rds.DBSnapshot, err error) {
	params := &rds.DescribeDBSnapshotsInput{
		DBInstanceIdentifier: &identifier,
	}
	results, err := svc.DescribeDBSnapshots(params)
	if err != nil || len(results.DBSnapshots) == 0 {
		err = fmt.Errorf("There are no Snapshots for %s", identifier)
		return
	}
	return results.DBSnapshots[0], nil
}
