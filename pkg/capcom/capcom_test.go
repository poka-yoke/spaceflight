package capcom

import (
	"errors"
	"fmt"
	"net"
	"testing"
)

func TestListSecurityGroups(t *testing.T) {
	svc := &mockEC2Client{}
	out := ListSecurityGroups(svc)
	expected := []string{
		fmt.Sprintf("* %10s %20s %s\n", "sg-1234", "", ""),
	}
	for index, line := range out {
		if line != expected[index] {
			t.Error("Unexpected output")
		}
	}
}

func TestBuildIPPermission(t *testing.T) {
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
		_, err := BuildIPPermission(
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

func TestCreateSG(t *testing.T) {
	data := []struct {
		name        string
		description string
		vpcid       string
		out         string
	}{
		{
			name:        "",
			description: "Non-VPC success",
			vpcid:       "",
			out:         "sg-12345678",
		},
		{
			name:        "",
			description: "VPC success",
			vpcid:       "vpc-12345678",
			out:         "sg-12345678",
		},
	}

	svc := &mockEC2Client{}
	for _, tc := range data {
		t.Run(
			tc.description,
			func(t *testing.T) {
				out := CreateSG(
					tc.name,
					tc.description,
					tc.vpcid,
					svc,
				)
				if out != tc.out {
					t.Error("Unexpected output")
				}
			},
		)
	}
}

func TestFindSGByName(t *testing.T) {
	data := []struct {
		name string
		vpc  string
		ret  []string
	}{
		{
			name: "",
			vpc:  "",
			ret: []string{
				"sg-1234",
			},
		},
	}
	svc := &mockEC2Client{}
	for _, tc := range data {
		ret := FindSGByName(tc.name, tc.vpc, svc)
		for index := range ret {
			if ret[index] != tc.ret[index] {
				t.Error("Unexpected output")
			}
		}
	}
}

func TestFindSecurityGroupsWithRange(t *testing.T) {
	data := []struct {
		cidr string
		err  error
		ret  []SearchResult
	}{
		{
			cidr: "1.2.3.4/32",
			err:  nil,
			ret: []SearchResult{
				SearchResult{
					GroupID:  "sg-1234",
					Protocol: "tcp",
					Port:     22,
					Source:   "1.2.3.4/32",
				}},
		},
	}

	svc := &mockEC2Client{}
	for _, tc := range data {
		ret, err := FindSecurityGroupsWithRange(svc, tc.cidr)
		if (err != nil && tc.err == nil) ||
			(err == nil && tc.err != nil) {
			t.Error("Unexpected/mismatched error")
		}
		if len(ret) != len(tc.ret) {
			t.Error("Mismatched results and expectations length")
		}
		for k, v := range ret {
			if v.String() != tc.ret[k].String() {
				t.Errorf(
					"Unexpected output %s != %s",
					v.String(),
					tc.ret[k].String(),
				)
			}
		}
	}
}

func TestNetworkContainsIPCheck(t *testing.T) {
	data := []struct {
		cidr string
		ip   string
		ret  bool
		err  error
	}{
		{
			cidr: "0.0.0.0/0",
			ip:   "1.2.3.4/32",
			ret:  true,
			err:  nil,
		},
		{
			cidr: "deadbeef",
			ip:   "1.2.3.4/32",
			ret:  false,
			err:  errors.New(""),
		},
		{
			cidr: "192.168.1.0/24",
			ip:   "192.168.1.1/32",
			ret:  true,
			err:  nil,
		},
		{
			cidr: "192.168.1.0/24",
			ip:   "192.168.3.1/32",
			ret:  false,
			err:  nil,
		},
		{
			cidr: "192.168.1.0/24",
			ip:   "192.168.1.0/24",
			ret:  true,
			err:  nil,
		},
		{
			cidr: "192.168.1.0/24",
			ip:   "192.168.3.0/24",
			ret:  false,
			err:  nil,
		},
		{
			cidr: "192.168.0.0/16",
			ip:   "192.168.10.0/24",
			ret:  true,
			err:  nil,
		},
	}
	for _, tc := range data {
		ip, _, _ := net.ParseCIDR(tc.ip)
		ret, err := NetworkContainsIPCheck(tc.cidr, ip)
		if (err != nil && tc.err == nil) ||
			(err == nil && tc.err != nil) {
			t.Error("Unexpected/mismatched error")
		}
		if ret != tc.ret {
			t.Error("Unexpected mismatch")
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
	svc := &mockEC2Client{}
	for _, tc := range data {
		perm, _ := BuildIPPermission(tc.origin, tc.proto, tc.port)
		out := AuthorizeAccessToSecurityGroup(
			svc,
			perm,
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
	svc := &mockEC2Client{}
	for _, tc := range data {
		perm, _ := BuildIPPermission(tc.origin, tc.proto, tc.port)
		out := RevokeAccessToSecurityGroup(
			svc,
			perm,
			tc.destination,
		)
		if out != tc.expected {
			t.Error("Unexpected mismatch")
		}
	}

}
