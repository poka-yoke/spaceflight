package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

var deleteInstanceCases = []createInstanceCase{
	// Deleting simple instance
	{
		testCase: testCase{
			expected:      "",
			expectedError: "",
		},
		name:       "Deleting simple instance",
		identifier: "test1",
		instances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: aws.String("test1"),
				DBInstanceStatus:     aws.String("available"),
			},
		},
	},
	// Deleting non existing instance
	{
		testCase: testCase{
			expected:      "",
			expectedError: "No such instance test2",
		},
		name:       "Deleting non existing instance",
		identifier: "test2",
	},
}

func TestDeleteInstance(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, test := range deleteInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc.AddInstances(test.instances)
				err := odin.DeleteInstance(
					test.identifier,
					svc,
				)
				test.check("", err, t)
				_, instance := svc.FindInstance(test.identifier)
				if instance != nil {
					t.Errorf(
						"%s instance should be deleted",
						test.identifier,
					)
				}
			},
		)
	}
}
