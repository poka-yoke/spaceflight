package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
)

// awsLogin initializes connection to AWS API
func rdsLogin(region string) rdsiface.RDSAPI {
	return rds.New(
		session.New(
			&aws.Config{
				Region: aws.String(
					region,
				),
			},
		),
	)
}
