package odin_test

import (
	"testing"
	"time"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

type createInstanceCase struct {
	name               string
	identifier         string
	instanceType       string
	masterUserPassword string
	masterUser         string
	size               int64
	endpoint           string
	expectedError      string
}

var createInstanceCases = []createInstanceCase{
	// Creating simple instance
	{
		name:               "Creating simple instance",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		size:               5,
		endpoint:           "test1.0.us-east-1.rds.amazonaws.com",
		expectedError:      "",
	},
	// Fail because empty user
	{
		name:               "Fail because empty user",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "",
		masterUserPassword: "master",
		size:               5,
		endpoint:           "",
		expectedError:      "Specify Master User",
	},
	// Fail because empty password
	{
		name:               "Fail because empty password",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "",
		size:               5,
		endpoint:           "",
		expectedError:      "Specify Master User Password",
	},
	// Fail because non-present size
	{
		name:               "Fail because non-present size",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		endpoint:           "",
		expectedError:      "Specify size between 5 and 6144",
	},
	// Fail because too small size
	{
		name:               "Fail because too small size",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		endpoint:           "",
		expectedError:      "Specify size between 5 and 6144",
	},
	// Fail because too big size
	{
		name:               "Fail because too big size",
		identifier:         "test1",
		instanceType:       "db.m1.small",
		masterUser:         "master",
		masterUserPassword: "master",
		size:               6145,
		endpoint:           "",
		expectedError:      "Specify size between 5 and 6144",
	},
}

func TestCreateInstance(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, useCase := range createInstanceCases {
		t.Run(
			useCase.name,
			func(t *testing.T) {
				params := odin.CreateParams{
					InstanceType: useCase.instanceType,
					User:         useCase.masterUser,
					Password:     useCase.masterUserPassword,
					Size:         useCase.size,
				}
				endpoint, err := odin.CreateInstance(
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
