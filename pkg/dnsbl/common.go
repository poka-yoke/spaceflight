package dnsbl

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// reverseAddress converts IP address into reversed address for query.
func reverseAddress(ipAddress string) string {
	ipAddressValues := strings.Split(ipAddress, ".")
	var sb strings.Builder
	maxIndex := len(ipAddressValues) - 1
	for i := range ipAddressValues {
		sb.WriteString(ipAddressValues[maxIndex-i])
		if i != maxIndex {
			sb.WriteString(".")
		}
	}
	return sb.String()
}

// Query queries a DNSBL and returns true if the argument gets a match
// in the BL. Returns the number of results in the lookup where 0
// means not present.
func Query(provider, address string) int {
	queryAddress := fmt.Sprintf(
		"%v.%v",
		reverseAddress(address),
		provider,
	)

	// We ignore errors because the providers where we are not
	// flagged can't be resolved. We can not distinguish if we are
	// not on their list or their service is broken.
	result, _ := net.LookupHost(queryAddress)
	if len(result) > 0 {
		log.Printf("%v returned %v\n", queryAddress, result)
	}
	return len(result)
}
