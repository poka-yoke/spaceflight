package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/Devex/spaceflight/pkg/odin"
)

type deleteInstanceCase struct {
	testCase
	name       string
	identifier string
	snapshotID string
	instances  []*rds.DBInstance
}

var deleteInstanceCases = []deleteInstanceCase{
	// Deleting simple instance
	{
		testCase: testCase{
			expected:      "",
			expectedError: "",
		},
		name:       "Deleting simple instance",
		identifier: "test1",
		snapshotID: "",
		instances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: aws.String("test1"),
				DBInstanceStatus:     aws.String("available"),
			},
		},
	},
	// Deleting simple instance with final snapshot
	{
		testCase: testCase{
			expected:      "",
			expectedError: "",
		},
		name:       "Deleting simple instance with final snapshot",
		identifier: "test3",
		snapshotID: "test3-final",
		instances: []*rds.DBInstance{
			{
				DBInstanceIdentifier: aws.String("test3"),
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
		snapshotID: "",
	},
}

func TestDeleteInstance(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, test := range deleteInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				svc.addInstances(test.instances)
				err := odin.DeleteInstance(
					test.identifier,
					test.snapshotID,
					svc,
				)
				test.check("", err, t)
				_, instance, _ := svc.findInstance(
					test.identifier,
				)
				if instance != nil {
					t.Errorf(
						"%s instance should be deleted",
						test.identifier,
					)
				}
				if test.snapshotID == "" {
					_, s, _ := svc.FindSnapshotInstance(
						test.identifier,
					)
					if s != nil {
						t.Errorf(
							"%s should not exist",
							*s.DBSnapshotIdentifier,
						)
					}
				} else {
					_, _, err = svc.FindSnapshot(
						test.snapshotID,
					)
					if err != nil {
						t.Errorf(
							"%s should be created",
							test.snapshotID,
						)
					}
				}
			},
		)
	}
}
