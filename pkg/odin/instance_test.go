package odin_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/rds"

	"github.com/poka-yoke/spaceflight/internal/test/mocks"
	"github.com/poka-yoke/spaceflight/pkg/odin"
)

func TestCloneDBInput(t *testing.T) {
	tt := []struct {
		input            odin.Instance
		allocatedStorage int64
		masterUsername   string
		err              error
	}{
		// Underspecified options
		{
			input:            odin.Instance{OriginalInstanceName: ""},
			allocatedStorage: 0,
			masterUsername:   "",
			err:              fmt.Errorf("original Instance Name was empty"),
		},
		// Non-existing snapshot
		{
			input:            odin.Instance{OriginalInstanceName: "im-not-here"},
			allocatedStorage: 0,
			masterUsername:   "",
			err: fmt.Errorf(
				"no snapshot found for %s instance",
				"im-not-here",
			),
		},
		// Existing snapshot
		{
			input:            odin.Instance{OriginalInstanceName: "production-rds"},
			allocatedStorage: 10,
			masterUsername:   "owner",
			err:              nil,
		},
	}
	svc := mocks.NewRDSClient()
	svc.AddSnapshots([]*rds.DBSnapshot{exampleSnapshot1})
	for _, tc := range tt {
		res, err := tc.input.CloneDBInput(svc)
		if tc.err != nil &&
			err.Error() != tc.err.Error() {
			t.Errorf(
				"Expected: %s, but got %s",
				tc.err.Error(),
				err.Error(),
			)
		}
		if err == nil {
			if res.AllocatedStorage != nil &&
				*res.AllocatedStorage != tc.allocatedStorage {
				t.Errorf(
					"Expected: %v, got %v",
					tc.allocatedStorage,
					*res.AllocatedStorage,
				)
			}
			if res.MasterUsername != nil &&
				*res.MasterUsername != tc.masterUsername {
				t.Errorf(
					"Expected: %s, but got %s",
					tc.masterUsername,
					*res.MasterUsername,
				)
			}
		}
	}
}

func TestCreateDBInput(t *testing.T) {
	tt := []struct {
		input            odin.Instance
		allocatedStorage int64
		masterUsername   string
	}{
		// No Snapshot
		{
			input:            odin.Instance{LastSnapshot: nil},
			allocatedStorage: 0,
			masterUsername:   "",
		},
		// With snapshot
		{
			input:            odin.Instance{LastSnapshot: exampleSnapshot1},
			allocatedStorage: 10,
			masterUsername:   "owner",
		},
	}
	svc := mocks.NewRDSClient()
	svc.AddSnapshots([]*rds.DBSnapshot{exampleSnapshot1})
	for _, tc := range tt {
		res, err := tc.input.CreateDBInput()
		if err != nil {
			t.Errorf("It should not fail")
		}
		if res.AllocatedStorage != nil &&
			*res.AllocatedStorage != tc.allocatedStorage {
			t.Errorf(
				"Expected: %v, got %v",
				tc.allocatedStorage,
				*res.AllocatedStorage,
			)
		}
		if res.MasterUsername != nil &&
			*res.MasterUsername != tc.masterUsername {
			t.Errorf(
				"Expected: %s, but got %s",
				tc.masterUsername,
				*res.MasterUsername,
			)
		}
	}
}

func TestDeleteDBInput(t *testing.T) {
	tt := []struct {
		input             odin.Instance
		skipFinalSnapshot bool
		finalSnapshotID   string
	}{
		{
			input:             odin.Instance{FinalSnapshotID: ""},
			skipFinalSnapshot: true,
		},
		{
			input:           odin.Instance{FinalSnapshotID: "123"},
			finalSnapshotID: "123",
		},
	}

	for _, tc := range tt {
		res, err := tc.input.DeleteDBInput()
		if err != nil {
			t.Errorf("It should not fail")
		}
		if res.SkipFinalSnapshot != nil &&
			*res.SkipFinalSnapshot != tc.skipFinalSnapshot {
			t.Errorf(
				"Expected: %v. Received %v",
				tc.skipFinalSnapshot,
				*res.SkipFinalSnapshot,
			)
		}
		if !tc.skipFinalSnapshot &&
			*res.FinalDBSnapshotIdentifier != tc.finalSnapshotID {
			t.Errorf(
				"Expected: %s, but got %s",
				*res.FinalDBSnapshotIdentifier,
				tc.finalSnapshotID,
			)
		}
	}
}

func TestModifyDBInputGroups(t *testing.T) {
	tt := []struct {
		securityGroups []string
	}{
		// No groups
		{
			securityGroups: []string{},
		},
		// Just one group
		{
			securityGroups: []string{
				"sg-12345",
			},
		},
		// Several groups
		{
			securityGroups: []string{
				"sg-12345",
				"sg-23456",
				"sg-34567",
				"sg-45678",
			},
		},
	}
	for _, tc := range tt {
		res, err := odin.Instance{SecurityGroups: tc.securityGroups}.ModifyDBInput(false)
		if err != nil {
			t.Errorf("It should not fail")
		}
		if len(res.VpcSecurityGroupIds) != len(tc.securityGroups) {
			t.Errorf(
				"Expected: %d entries, but got %d",
				len(tc.securityGroups),
				len(res.VpcSecurityGroupIds),
			)
		}
		for i, sgid := range tc.securityGroups {
			if *res.VpcSecurityGroupIds[i] != sgid {
				t.Errorf(
					"Expected: %s, but got %s at position %d",
					sgid,
					*res.VpcSecurityGroupIds[i],
					i,
				)
			}
		}
	}
}

func TestModifyDBInputApplyNow(t *testing.T) {
	tt := []struct {
		applyNow bool
	}{
		{applyNow: true},
		{applyNow: false},
	}
	for _, tc := range tt {
		res, err := odin.Instance{}.ModifyDBInput(tc.applyNow)
		if err != nil {
			t.Errorf("It should not fail")
		}
		if *res.ApplyImmediately != tc.applyNow {
			t.Errorf(
				"Expected: %v, but got %v",
				tc.applyNow,
				*res.ApplyImmediately,
			)
		}
	}
}

func TestRestoreDBInput(t *testing.T) {
	tt := []struct {
		input                odin.Instance
		dbSnapshotIdentifier string
		err                  error
	}{
		// Underspecified options
		{
			input: odin.Instance{OriginalInstanceName: ""},
			err:   fmt.Errorf("original Instance Name was empty"),
		},
		// Non-existing snapshot
		{
			input: odin.Instance{OriginalInstanceName: "im-not-here"},
			err: fmt.Errorf(
				"no snapshot found for %s instance",
				"im-not-here",
			),
			dbSnapshotIdentifier: "not-found",
		},
		// Existing snapshot
		{
			input:                odin.Instance{OriginalInstanceName: "production-rds"},
			err:                  nil,
			dbSnapshotIdentifier: "rds:production-2015-06-11",
		},
	}
	svc := mocks.NewRDSClient()
	svc.AddSnapshots([]*rds.DBSnapshot{exampleSnapshot1})
	for _, tc := range tt {
		res, err := tc.input.RestoreDBInput(svc)
		if tc.err != nil &&
			err.Error() != tc.err.Error() {
			t.Errorf(
				"Expected: %s, but got %s",
				tc.err.Error(),
				err.Error(),
			)
		}
		if err == nil {
			switch {
			case res == nil:
				t.Errorf("Response should not be nil")
			case res.DBSnapshotIdentifier == nil:
				t.Errorf("Snapshot Identifier should not be nil")
			case *res.DBSnapshotIdentifier != tc.dbSnapshotIdentifier:
				t.Errorf(
					"Expected: %s, but got %s",
					*res.DBSnapshotIdentifier,
					tc.dbSnapshotIdentifier,
				)
			}
		}
	}
}
