package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/Devex/spaceflight/pkg/odin"
)

type scaleInstanceCase struct {
	testCase
	name         string
	identifier   string
	instanceType string
	delayChange  bool
	instances    []*rds.DBInstance
}

var scaleInstanceCases = []scaleInstanceCase{
	// Scaling up instance
	{
		testCase: testCase{
			expected:      "Instance test1 is db.m1.small",
			expectedError: "",
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
		testCase: testCase{
			expected:      "",
			expectedError: "No such instance test1",
		},
		name:         "Fail to scale up non existing instance",
		identifier:   "test1",
		instanceType: "db.m1.small",
		delayChange:  false,
		instances:    []*rds.DBInstance{},
	},
	// Fail to scale up non available instance
	{
		testCase: testCase{
			expected:      "",
			expectedError: "test1 instance state is not available",
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
		testCase: testCase{
			expected:      "Instance test1 is db.m1.small",
			expectedError: "",
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
		testCase: testCase{
			expected:      "Instance test1 is db.m1.medium",
			expectedError: "",
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
	odin.Duration = time.Duration(0)
	for _, test := range scaleInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc.addInstances(test.instances)
				actual, err := odin.ScaleInstance(
					test.identifier,
					test.instanceType,
					test.delayChange,
					svc,
				)
				test.check(actual, err, t)
			},
		)
	}
}
