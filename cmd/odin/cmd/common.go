package cmd

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"

	"github.com/poka-yoke/spaceflight/pkg/odin"
)

const (
	// RFC8601 is the date/time format used by AWS.
	RFC8601 = "2006-01-02T15:04:05-07:00"
)

// awsLogin initializes connection to AWS API
func rdsLogin(region string) rdsiface.RDSAPI {
	session, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
		},
	)
	if err != nil {
		panic(err)
	}
	return rds.New(session)
}

func waitForInstance(
	instance *rds.DBInstance,
	svc rdsiface.RDSAPI,
	status string,
	duration time.Duration,
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
		time.Sleep(duration)
	}
	return
}

// modifyInstance enqueues a modify operation
func modifyInstance(
	params odin.ModifiableParams,
	svc rdsiface.RDSAPI,
	delay bool,
) (err error) {
	rdsParams, err := params.ModifyDBInput(delay)
	if err != nil {
		return err
	}
	_, err = svc.ModifyDBInstance(rdsParams)
	return err
}
