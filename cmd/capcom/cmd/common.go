package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// connect initializes connection to AWS API
func connect() ec2iface.EC2API {
	region := "us-east-1"
	return ec2.New(
		session.New(
			&aws.Config{
				Region: aws.String(region),
			},
		),
	)
}
