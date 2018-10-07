package cmd

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/internal/test_case"
)

var exampleSnapshot1Out = fmt.Sprintf(
	"%v %v %v %v\n",
	*exampleSnapshot1ID,
	*exampleSnapshot1DBID,
	exampleSnapshot1Time,
	"available",
)

var exampleSnapshot2Out = fmt.Sprintf(
	"%v %v %v %v\n",
	*exampleSnapshot2ID,
	*exampleSnapshot2DBID,
	exampleSnapshot2Time,
	"available",
)

type listSnapshotsCase struct {
	testcase.TestCase
	name       string
	snapshots  []*rds.DBSnapshot
	instanceID string
}

var printSnapshotsCases = []listSnapshotsCase{
	// No snapshots
	{
		TestCase: testcase.TestCase{
			Expected:      "",
			ExpectedError: "",
		},
		name:      "No snapshots",
		snapshots: []*rds.DBSnapshot{},
	},
	// One snapshot
	{
		TestCase: testcase.TestCase{
			Expected:      exampleSnapshot1Out,
			ExpectedError: "",
		},
		name:      "One snapshot",
		snapshots: []*rds.DBSnapshot{exampleSnapshot1},
	},
	// Two snapshots
	{
		TestCase: testcase.TestCase{
			Expected: strings.Join(
				[]string{
					exampleSnapshot2Out,
					exampleSnapshot1Out,
				},
				"",
			),
			ExpectedError: "",
		},
		name: "Two snapshots",
		snapshots: []*rds.DBSnapshot{
			exampleSnapshot2,
			exampleSnapshot1,
		},
	},
}

func TestPrintSnapshot(t *testing.T) {
	for _, test := range printSnapshotsCases {
		t.Run(
			test.name,
			func(t *testing.T) {
				actual := printSnapshots(
					test.snapshots,
				)
				test.Check(actual, nil, t)
			},
		)
	}
}
