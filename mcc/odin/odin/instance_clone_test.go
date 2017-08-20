package odin_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

type cloneInstanceCase struct {
	name                 string
	identifier           string
	instanceType         string
	masterUserPassword   string
	masterUser           string
	size                 int64
	originalInstanceName string
	endpoint             string
	expectedError        string
	snapshot             *rds.DBSnapshot
}

var cloneInstanceCases = []cloneInstanceCase{
	// Uses snapshot to copy from
	{
		name:                 "Uses snapshot to copy from",
		identifier:           "test1",
		instanceType:         "db.m1.small",
		masterUser:           "master",
		masterUserPassword:   "master",
		size:                 6144,
		originalInstanceName: "production",
		endpoint:             "test1.0.us-east-1.rds.amazonaws.com",
		expectedError:        "",
		snapshot:             exampleSnapshot1,
	},
	// Uses non existing snapshot to copy from
	{
		name:                 "Uses non existing snapshot to copy from",
		identifier:           "test1",
		instanceType:         "db.m1.small",
		masterUser:           "master",
		masterUserPassword:   "master",
		size:                 6144,
		originalInstanceName: "develop",
		endpoint:             "",
		expectedError:        "Couldn't find snapshot for develop instance",
		snapshot:             exampleSnapshot1,
	},
}

func TestCloneInstance(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, useCase := range cloneInstanceCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				if useCase.originalInstanceName != "" {
					svc.dbSnapshots[*useCase.snapshot.DBInstanceIdentifier] = []*rds.DBSnapshot{
						useCase.snapshot,
					}
				}
				params := odin.CloneParams{
					DBInstanceType:       useCase.instanceType,
					DBUser:               useCase.masterUser,
					DBPassword:           useCase.masterUserPassword,
					Size:                 useCase.size,
					OriginalInstanceName: useCase.originalInstanceName,
				}
				endpoint, err := odin.CloneInstance(
					useCase.identifier,
					params,
					svc,
				)
				if err != nil {
					if err.Error() != useCase.expectedError {
						t.Errorf(
							"Unexpected error %s",
							err,
						)
					}
				} else if useCase.expectedError != "" {
					t.Errorf(
						"Expected error %s didn't happened",
						useCase.expectedError,
					)
				} else {
					if endpoint != useCase.endpoint {
						t.Errorf(
							"Unexpected output: %s should be %s",
							endpoint,
							useCase.endpoint,
						)
					}
				}
			},
		)
	}
}
