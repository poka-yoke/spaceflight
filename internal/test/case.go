package test

import (
	"testing"

	"github.com/go-test/deep"
)

// Case contains the basic elements of a test case:
//   * expected is what should be got.
//   * expectedError is the string of the error message in case of such error being expected.
type Case struct {
	Expected      interface{}
	ExpectedError string
}

// expectingError helps determine if the case was expecting an specific error.
func (c *Case) expectingError(err error) bool {
	return c.ExpectedError != "" && err.Error() == c.ExpectedError
}

// Check is to be used by the test executor to validate the case is passing.
func (c *Case) Check(actual interface{}, err error, t *testing.T) {
	switch {
	case err != nil && !c.expectingError(err):
		t.Errorf(
			"Unexpected error: %v",
			err,
		)
	case err != nil && c.expectingError(err):
	case err == nil && c.ExpectedError != "":
		t.Errorf(
			"Expected error: %v missing",
			c.ExpectedError,
		)
	case err == nil:
		if diff := deep.Equal(
			actual,
			c.Expected,
		); diff != nil {
			t.Errorf(
				"Unexpected output: %s",
				diff,
			)
		}
	}
}
