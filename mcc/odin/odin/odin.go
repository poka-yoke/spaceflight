package odin

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// Init initializes connection to AWS API
func Init() rdsiface.RDSAPI {
	region := "us-east-1"
	sess := session.New(&aws.Config{Region: aws.String(region)})
	return rds.New(sess)
}

var duration = time.Duration(5) * time.Second

// CreateDBInstance creates a new RDS database instance. If a vpcid is specified the security
// group will be in that VPC
func CreateDBInstance(
	instanceName string,
	instanceType string,
	masterUser string,
	masterUserPassword string,
	size int64,
	svc rdsiface.RDSAPI,
) (result string, err error) {
	params := &rds.CreateDBInstanceInput{
		AllocatedStorage:     &size,
		DBInstanceClass:      &instanceType,
		DBInstanceIdentifier: &instanceName,
		Engine:               aws.String("postgres"),
		EngineVersion:        aws.String("9.4.11"),
		DBSecurityGroups: []*string{
			aws.String("default"),
		},
		MasterUserPassword: &masterUserPassword,
		MasterUsername:     &masterUser,
		Tags: []*rds.Tag{
			{
				Key:   aws.String("Name"),
				Value: &instanceName,
			},
		},
	}
	if err = params.Validate(); err != nil {
		return
	}
	res, err := svc.CreateDBInstance(params)
	if err != nil {
		return
	}
	instance := *res.DBInstance
	for *instance.DBInstanceStatus != "available" {
		res2, err2 := svc.DescribeDBInstances(&rds.DescribeDBInstancesInput{
			DBInstanceIdentifier: instance.DBInstanceIdentifier,
		})
		if err2 != nil {
			err = err2
			return
		}
		instance = *res2.DBInstances[0]
		// This is to avoid AWS API rate throttling.
		time.Sleep(duration)
	}
	result = *instance.Endpoint.Address
	return
}
