package testcase

import (
	"testing"

	"github.com/go-test/deep"
)

// TestCase contains the basic elements of a test case:
//   * expected is what should be got.
//   * expectedError is the string of the error message in case of such error being expected.
type TestCase struct {
	Expected      interface{}
	ExpectedError string
}

// expectingError helps determine if the case was expecting an specific error.
func (tc *TestCase) expectingError(err error) bool {
	return tc.ExpectedError != "" && err.Error() == tc.ExpectedError
}

// Check is to be used by the test executor to validate the case is passing.
func (tc *TestCase) Check(actual interface{}, err error, t *testing.T) {
	switch {
	case err != nil && !tc.expectingError(err):
		t.Errorf(
			"Unexpected error: %v",
			err,
		)
	case err != nil && tc.expectingError(err):
	case err == nil && tc.ExpectedError != "":
		t.Errorf(
			"Expected error: %v missing",
			tc.ExpectedError,
		)
	case err == nil:
		if diff := deep.Equal(
			actual,
			tc.Expected,
		); diff != nil {
			t.Errorf(
				"Unexpected output: %s",
				diff,
			)
		}
	}
}
