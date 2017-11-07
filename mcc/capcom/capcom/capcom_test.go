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
	for _, tt := range biptable {
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
		_, err := BuildIPPermission(
			tt.origin,
			tt.proto,
			tt.port,
		)
		if (err != nil && tt.err == nil) ||
			(err == nil && tt.err != nil) {
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
	for _, tt := range csgtable {
		t.Run(
			tt.description,
			func(t *testing.T) {
				out := CreateSG(
					tt.name,
					tt.description,
					tt.vpcid,
					svc,
				)
				if out != tt.out {
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
	for _, tt := range fsgbntable {
		ret := FindSGByName(tt.name, tt.vpc, svc)
		for index := range ret {
			if ret[index] != tt.ret[index] {
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
	for _, tt := range fsgwtable {
		ret, err := FindSecurityGroupsWithRange(svc, tt.cidr)
		if (err != nil && tt.err == nil) ||
			(err == nil && tt.err != nil) {
			t.Error("Unexpected/mismatched error")
		}
		if len(ret) != len(tt.ret) {
			t.Error("Mismatched results and expectations length")
		}
		for k, v := range ret {
			if v.String() != tt.ret[k].String() {
				t.Errorf(
					"Unexpected output %s != %s",
					v.String(),
					tt.ret[k].String(),
				)
			}
		}
	}
}

func TestNetworkContainsIPCheck(t *testing.T) {
	for _, tt := range ncictable {
		ip, _, _ := net.ParseCIDR(tt.ip)
		ret, err := NetworkContainsIPCheck(tt.cidr, ip)
		if (err != nil && tt.err == nil) ||
			(err == nil && tt.err != nil) {
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
			t.Error("Unexpected/mismatched error")
		}
		if ret != tt.ret {
			t.Error("Unexpected mismatch")
		}
	}
}
