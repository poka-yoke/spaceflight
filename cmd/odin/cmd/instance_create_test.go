package cmd

import (
	"testing"
	"time"

	"github.com/poka-yoke/spaceflight/internal/test_case"
	"github.com/poka-yoke/spaceflight/pkg/odin"
)

type createInstanceCase struct {
	testcase.TestCase
	name         string
	identifier   string
	instanceType string
	password     string
	user         string
	size         int64
}

var createInstanceCases = []createInstanceCase{
	// Creating simple instance
	{
		TestCase: testcase.TestCase{
			Expected:      "test1.0.us-east-1.rds.amazonaws.com",
			ExpectedError: "",
		},
		name:         "Creating simple instance",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
		size:         5,
	},
	// Fail because empty user
	{
		TestCase: testcase.TestCase{
			Expected:      "",
			ExpectedError: "Specify Master User",
		},
		name:         "Fail because empty user",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "",
		password:     "master",
		size:         5,
	},
	// Fail because empty password
	{
		TestCase: testcase.TestCase{
			Expected:      "",
			ExpectedError: "Specify Master User Password",
		},
		name:         "Fail because empty password",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "",
		size:         5,
	},
	// Fail because non-present size
	{
		TestCase: testcase.TestCase{
			Expected:      "",
			ExpectedError: "Specify size between 5 and 6144",
		},
		name:         "Fail because non-present size",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
	},
	// Fail because too small size
	{
		TestCase: testcase.TestCase{
			Expected:      "",
			ExpectedError: "Specify size between 5 and 6144",
		},
		name:         "Fail because too small size",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
	},
	// Fail because too big size
	{
		TestCase: testcase.TestCase{
			Expected:      "",
			ExpectedError: "Specify size between 5 and 6144",
		},
		name:         "Fail because too big size",
		identifier:   "test1",
		instanceType: "db.m1.small",
		user:         "master",
		password:     "master",
		size:         6145,
	},
}

func TestCreateInstance(t *testing.T) {
	svc := newMockRDSClient()
	for _, test := range createInstanceCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				params := odin.Instance{
					Identifier: test.identifier,
					Type:       test.instanceType,
					User:       test.user,
					Password:   test.password,
					Size:       test.size,
				}
				actual, err := createInstance(
					params,
					svc,
					time.Duration(0),
				)
				test.Check(actual, err, t)
			},
		)
	}
}
