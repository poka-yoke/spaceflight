package capcom

import (
	"fmt"
	"net"
)

// NetworkContainsIPCheck returns true if the subnet expresed in the
// CIDR in contains the IP object
func networkContainsIPCheck(cidr string, searchIP net.IP) (out bool, err error) {
	ip, sub, err := net.ParseCIDR(cidr)
	if err != nil {
		err = fmt.Errorf(
			"failed parsing CIDR %s",
			cidr,
		)
		return
	}
	out = ip.Equal(searchIP) || sub.Contains(searchIP)
	return
}

func isCIDR(origin string) bool {
	_, _, err := net.ParseCIDR(origin)
	return err == nil
}
