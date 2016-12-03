package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

var list bool

func ListSecurityGroups(svc *ec2.EC2) {
	res, err := svc.DescribeSecurityGroups(nil)
	if err != nil {
		log.Panic(err)
	}

	for _, sg := range res.SecurityGroups {
		fmt.Printf("* %10s %20s %s\n",
			*sg.GroupId,
			*sg.GroupName,
			*sg.Description)
	}
}

func main() {
	flag.BoolVar(&list, "l", false, "List all Security groups with ID, name and description")
	flag.Parse()

	region := "us-east-1"
	sess := session.New(&aws.Config{Region: aws.String(region)})
	svc := ec2.New(sess)

	if list {
		ListSecurityGroups(svc)
	}
}
