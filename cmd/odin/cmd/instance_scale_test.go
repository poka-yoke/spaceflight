package cmd

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/internal/test_case"
	"github.com/poka-yoke/spaceflight/pkg/odin"
)

type scaleInstanceCase struct {
	testcase.TestCase
	name         string
	identifier   string
	instanceType string
	delayChange  bool
	instances    []*rds.DBInstance
}

var scaleInstanceCases = []scaleInstanceCase{
	// Scaling up instance
	{
		TestCase: testcase.TestCase{
			Expected:      "Instance test1 is db.m1.small",
			ExpectedError: "",
		},
		name:         "Scaling up instance",
		identifier:   "test1",
		instanceType: "db.m1.small",
		delayChange:  false,
		instances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: aws.String("test1"),
				DBInstanceClass:      aws.String("db.t1.micro"),
				DBInstanceStatus:     aws.String("available"),
			},
		},
	},
	// Fail to scale up non existing instance
	{
		TestCase: testcase.TestCase{
			Expected:      "",
			ExpectedError: "No such instance test1",
		},
		name:         "Fail to scale up non existing instance",
		identifier:   "test1",
		instanceType: "db.m1.small",
		delayChange:  false,
		instances:    []*rds.DBInstance{},
	},
	// Fail to scale up non available instance
	{
		TestCase: testcase.TestCase{
			Expected:      "",
			ExpectedError: "test1 instance state is not available",
		},
		name:         "Fail to scale up non available instance",
		identifier:   "test1",
		instanceType: "db.m1.small",
		delayChange:  false,
		instances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: aws.String("test1"),
				DBInstanceClass:      aws.String("db.t1.micro"),
				DBInstanceStatus:     aws.String("modifying"),
			},
		},
	},
	// Scaling down instance
	{
		TestCase: testcase.TestCase{
			Expected:      "Instance test1 is db.m1.small",
			ExpectedError: "",
		},
		name:         "Scaling down instance",
		identifier:   "test1",
		instanceType: "db.m1.small",
		delayChange:  false,
		instances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: aws.String("test1"),
				DBInstanceClass:      aws.String("db.m1.medium"),
				DBInstanceStatus:     aws.String("available"),
			},
		},
	},
	// Delayed scaling down instance
	{
		TestCase: testcase.TestCase{
			Expected:      "Instance test1 is db.m1.medium",
			ExpectedError: "",
		},
		name:         "Scaling down instance with delay",
		identifier:   "test1",
		instanceType: "db.m1.small",
		delayChange:  true,
		instances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: aws.String("test1"),
				DBInstanceClass:      aws.String("db.m1.medium"),
				DBInstanceStatus:     aws.String("available"),
			},
		},
	},
}

func TestScaleInstance(t *testing.T) {
	svc := newMockRDSClient()
	for _, test := range scaleInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc.addInstances(test.instances)
				actual, err := scaleInstance(
					odin.Instance{
						Identifier: test.identifier,
						Type:       test.instanceType,
					},
					test.delayChange,
					svc,
					time.Duration(0),
				)
				test.Check(actual, err, t)
			},
		)
	}
}
