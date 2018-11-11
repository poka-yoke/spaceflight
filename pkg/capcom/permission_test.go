package capcom

import (
	"errors"
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
	}
	svc := &mocks.EC2Client{}
	for _, tc := range data {
		perm, _ := NewPermission(tc.origin, tc.proto, tc.port)
		out := perm.AddToSG(
			svc,
			tc.destination,
		)
		if out != tc.expected {
			t.Error("Unexpected mismatch")
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
	}
	svc := &mocks.EC2Client{}
	for _, tc := range data {
		perm, _ := NewPermission(tc.origin, tc.proto, tc.port)
		out := perm.RemoveToSG(
			svc,
			tc.destination,
		)
		if out != tc.expected {
			t.Error("Unexpected mismatch")
		}
	}

}
