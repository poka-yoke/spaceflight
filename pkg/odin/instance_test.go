package odin

import(
	"testing"
)

func TestDeleteDBInput(t *testing.T) {
	tt := []struct{
		input Instance
		skipFinalSnapshot bool
		finalSnapshotID string
	}{
		{
			input: Instance{FinalSnapshotID: ""},
			skipFinalSnapshot: true,
		},
		{
			input: Instance{FinalSnapshotID: "123"},
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
				"SkipFinalSnapshot should be %v. Received %v",
				tc.skipFinalSnapshot,
				*res.SkipFinalSnapshot,
			)
		}
		if !tc.skipFinalSnapshot &&
			*res.FinalDBSnapshotIdentifier != tc.finalSnapshotID {
			t.Errorf(
				"Wrong identifier. Expected %s, but got %s",
				*res.FinalDBSnapshotIdentifier,
				tc.finalSnapshotID,
			)
		}
	}
}
