package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// connect initializes connection to AWS API
func connect() *route53.Route53 {
	region := "us-east-1"
	return route53.New(
		session.New(
			&aws.Config{
				Region: aws.String(region),
			},
		),
	)
}
