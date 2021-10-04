package capcom

import (
	"errors"
	"fmt"
	"testing"

	"github.com/poka-yoke/spaceflight/internal/test/mocks"
)

func TestNewPermission(t *testing.T) {
	data := []struct {
		origin string
		proto  string
		port   int64
		err    error
	}{
		{
			origin: "1.2.3.4/32",
			proto:  "tcp",
			port:   int64(0),
			err:    nil,
		},
		{
			origin: "sg-",
			proto:  "tcp",
			port:   int64(0),
			err:    nil,
		},
		{
			origin: "1.2.3.4/32",
			proto:  "udp",
			port:   int64(0),
			err:    nil,
		},
		{
			origin: "sg-",
			proto:  "icmp",
			port:   int64(0),
			err:    nil,
		},
		{
			origin: "1.2.3./32",
			proto:  "udp",
			port:   int64(0),
			err:    errors.New(""),
		},
	}
	for _, tc := range data {
		_, err := NewPermission(
			tc.origin,
			tc.proto,
			tc.port,
		)
		if (err != nil && tc.err == nil) ||
			(err == nil && tc.err != nil) {
			t.Error(err)
		}
	}
}

func TestAuthorizeAccessToSecurityGroup(t *testing.T) {
	data := []struct {
		origin, proto string
		port          int64
		destination   string
		expected      bool
	}{
		{
			origin:   "1.2.3.4/32",
			proto:    "tcp",
			port:     int64(0),
			expected: true,
		},
		{
			origin:   "1.2.3.4/32",
			proto:    "tcp",
			port:     int64(0),
			expected: false,
		},
	}
	svc := &mocks.EC2Client{}
	for _, tc := range data {
		svc.FailAuthorizeSG = !tc.expected
		perm, _ := NewPermission(tc.origin, tc.proto, tc.port)
		out := perm.AddToSG(
			svc,
			tc.destination,
		)
		if out != tc.expected {
			t.Errorf(
				"Unexpected return. Expected %t but was %t",
				tc.expected,
				out,
			)
		}
		_, err := perm.Err()
		if err != nil && tc.expected {
			t.Error(err)
		}
	}
}

func TestRevokeAccessToSecurityGroup(t *testing.T) {
	data := []struct {
		origin, proto string
		port          int64
		destination   string
		expected      bool
	}{
		{
			origin:   "1.2.3.4/32",
			proto:    "tcp",
			port:     int64(0),
			expected: true,
		},
		{
			origin:   "1.2.3.4/32",
			proto:    "tcp",
			port:     int64(0),
			expected: false,
		},
	}
	svc := &mocks.EC2Client{}
	for _, tc := range data {
		svc.FailRevokeSG = !tc.expected
		perm, _ := NewPermission(tc.origin, tc.proto, tc.port)
		out := perm.RemoveToSG(
			svc,
			tc.destination,
		)
		if out != tc.expected {
			t.Errorf(
				"Unexpected return. Expected %t but was %t",
				tc.expected,
				out,
			)
		}
		_, err := perm.Err()
		if err != nil && tc.expected {
			t.Error(err)
		}
	}
}

// ExamplePermission_Err shows how to retrieve all the errors
// collected during a series of Permission operations.
func ExamplePermission_Err() {
	svc := &mocks.EC2Client{}
	perm, err := NewPermission("127.0.0.1/32", "tcp", 80)
	if err != nil {
		// Handle error
	}

	// Add the rule to some groups
	ok := true
	for _, sgid := range []string{"sg-1234", "sg-2345"} {
		// Force the first API call to fail
		svc.FailAuthorizeSG = ok

		if !perm.AddToSG(svc, sgid) {
			// We want to know if at least one of the
			// calls failed
			ok = false
		}
	}

	// Remove the rule from some groups
	ok2 := true
	for _, sgid := range []string{"sg-3456", "sg-5678"} {
		// Force the first API call to fail
		svc.FailRevokeSG = ok2

		if !perm.RemoveToSG(svc, sgid) {
			// We want to know if at least one of the
			// calls failed
			ok2 = false
		}
	}
	// If some call failed we check the errors
	if !ok || !ok2 {
		// At least one iteration
		cont := true
		for cont {
			// We don't know how many errors are there, so
			// loop until there are none left
			more, error := perm.Err()
			if error != nil {
				// Handle error and report it
				fmt.Println(error)
			}
			// `more` will be false once there are no more
			// errors, then `cont` will also be false and
			// we break out of the loop
			cont = more
		}
	}
	// Unordered output:
	// Error while adding on sg-1234: it had to fail
	// Error while revoking on sg-3456: it had to fail
}
