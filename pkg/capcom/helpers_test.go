package capcom

import (
	"errors"
	"net"
	"testing"
)

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
		ret, err := networkContainsIPCheck(tc.cidr, ip)
		if (err != nil && tc.err == nil) ||
			(err == nil && tc.err != nil) {
			t.Error("Unexpected/mismatched error")
		}
		if ret != tc.ret {
			t.Error("Unexpected mismatch")
		}
	}
}
