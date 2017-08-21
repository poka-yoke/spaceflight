package odin_test

import (
	"testing"
	"time"

	"github.com/go-test/deep"

	"github.com/poka-yoke/spaceflight/mcc/odin/odin"
)

type createInstanceCase struct {
	name          string
	identifier    string
	instanceType  string
	password      string
	user          string
	size          int64
	expected      string
	expectedError string
}

func (t *createInstanceCase) expectingError(err error) bool {
	return t.expectedError != "" && err.Error() != t.expectedError
}

var createInstanceCases = []createInstanceCase{
	// Creating simple instance
	{
		name:          "Creating simple instance",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		user:          "master",
		password:      "master",
		size:          5,
		expected:      "test1.0.us-east-1.rds.amazonaws.com",
		expectedError: "",
	},
	// Fail because empty user
	{
		name:          "Fail because empty user",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		user:          "",
		password:      "master",
		size:          5,
		expected:      "",
		expectedError: "Specify Master User",
	},
	// Fail because empty password
	{
		name:          "Fail because empty password",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		user:          "master",
		password:      "",
		size:          5,
		expected:      "",
		expectedError: "Specify Master User Password",
	},
	// Fail because non-present size
	{
		name:          "Fail because non-present size",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		user:          "master",
		password:      "master",
		expected:      "",
		expectedError: "Specify size between 5 and 6144",
	},
	// Fail because too small size
	{
		name:          "Fail because too small size",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		user:          "master",
		password:      "master",
		expected:      "",
		expectedError: "Specify size between 5 and 6144",
	},
	// Fail because too big size
	{
		name:          "Fail because too big size",
		identifier:    "test1",
		instanceType:  "db.m1.small",
		user:          "master",
		password:      "master",
		size:          6145,
		expected:      "",
		expectedError: "Specify size between 5 and 6144",
	},
}

func TestCreateInstance(t *testing.T) {
	svc := newMockRDSClient()
	odin.Duration = time.Duration(0)
	for _, test := range createInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				params := odin.CreateParams{
					InstanceType: test.instanceType,
					User:         test.user,
					Password:     test.password,
					Size:         test.size,
				}
				actual, err := odin.CreateInstance(
					test.identifier,
					params,
					svc,
				)
				switch {
				case err != nil && test.expectingError(err):
					t.Errorf(
						"Unexpected error: %v",
						err,
					)
				case err == nil && test.expectedError != "":
					t.Errorf(
						"Expected error: %v missing",
						test.expectedError,
					)
				case err == nil:
					if diff := deep.Equal(
						actual,
						test.expected,
					); diff != nil {
						t.Errorf(
							"Unexpected output: %s",
							diff,
						)
					}
				}
			},
		)
	}
}
