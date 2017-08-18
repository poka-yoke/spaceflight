package odin_test

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/Devex/spaceflight/mcc/odin/odin"
)

type getRestoreDBInstanceFromDBSnapshotInputCase struct {
	name                         string
	identifier                   string
	createDBParams               odin.CreateDBParams
	snapshot                     *rds.DBSnapshot
	expectedRestoreSnapshotInput *rds.RestoreDBInstanceFromDBSnapshotInput
	expectedError                string
}

var getRestoreDBInstanceFromDBSnapshotInputCases = []getRestoreDBInstanceFromDBSnapshotInputCase{
	// Params with Snapshot
	{
		name:       "Params with Snapshot",
		identifier: "production",
		createDBParams: odin.CreateDBParams{
			DBInstanceType:       "db.m1.medium",
			DBUser:               "owner",
			DBPassword:           "password",
			Size:                 5,
			OriginalInstanceName: "production",
		},
		snapshot: exampleSnapshot1,
		expectedRestoreSnapshotInput: &rds.RestoreDBInstanceFromDBSnapshotInput{
			DBInstanceClass:      aws.String("db.m1.medium"),
			DBInstanceIdentifier: aws.String("production"),
			DBSnapshotIdentifier: aws.String("rds:production-2015-06-11"),
			Engine:               aws.String("postgres"),
		},
		expectedError: "",
	},
	// Params with Snapshot without OriginalInstanceName
	{
		name:       "Params with Snapshot without OriginalInstanceName",
		identifier: "production",
		createDBParams: odin.CreateDBParams{
			DBInstanceType: "db.m1.medium",
			DBUser:         "owner",
			DBPassword:     "password",
			Size:           5,
		},
		snapshot:                     exampleSnapshot1,
		expectedRestoreSnapshotInput: nil,
		expectedError:                "Original Instance Name was empty",
	},
	// Params with non existing Snapshot
	{
		name:       "Params with non existing Snapshot",
		identifier: "production",
		createDBParams: odin.CreateDBParams{
			DBInstanceType:       "db.m1.medium",
			DBUser:               "owner",
			DBPassword:           "password",
			Size:                 5,
			OriginalInstanceName: "develop",
		},
		snapshot:                     exampleSnapshot1,
		expectedRestoreSnapshotInput: nil,
		expectedError:                "Couldn't find snapshot for develop instance",
	},
}

func equalsRestoreDBInstanceFromDBSnapshotInput(input1, input2 *rds.RestoreDBInstanceFromDBSnapshotInput) bool {
	switch {
	case *input1.DBInstanceIdentifier != *input2.DBInstanceIdentifier:
		return false
	case *input1.DBSnapshotIdentifier != *input2.DBSnapshotIdentifier:
		return false
	case *input1.DBInstanceClass != *input2.DBInstanceClass:
		return false
	case *input1.Engine != *input2.Engine:
		return false
	}
	return true
}

func TestGetRestoreDBInstanceFromDBSnapshotInput(t *testing.T) {
	svc := newMockRDSClient()
	for _, useCase := range getRestoreDBInstanceFromDBSnapshotInputCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				if useCase.snapshot != nil {
					svc.dbSnapshots[*useCase.snapshot.DBInstanceIdentifier] = []*rds.DBSnapshot{useCase.snapshot}
				}
				restoreSnapshotInput, err := useCase.createDBParams.GetRestoreDBInstanceFromDBSnapshotInput(
					useCase.identifier,
					svc,
				)
				if err != nil {
					if useCase.expectedError == "" ||
						err.Error() != useCase.expectedError {
						t.Errorf(
							"Unexpected error happened: %v",
							err,
						)
					}
				} else {
					if !equalsRestoreDBInstanceFromDBSnapshotInput(
						restoreSnapshotInput,
						useCase.expectedRestoreSnapshotInput,
					) {
						t.Errorf(
							"Unexpected output: %s should be %s",
							restoreSnapshotInput,
							useCase.expectedRestoreSnapshotInput,
						)
					}
				}
			},
		)
	}
}
