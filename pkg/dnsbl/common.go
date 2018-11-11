package dnsbl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// GetProviders returns a slice of provider addresses to check
func GetProviders(ipAddress string, lists io.Reader) (providers []string) {
	reverseAddress := reverseAddress(ipAddress)
	scanner := bufio.NewScanner(lists)
	for scanner.Scan() {
		reversedIPAddress := fmt.Sprintf(
			"%v.%v",
			reverseAddress,
			scanner.Text(),
		)
		providers = append(providers, reversedIPAddress)
	}
	return providers
}

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

// lookuper interface provides a method to do hostname lookups
type lookuper interface {
	lookup(string) ([]string, error)
}

// query queries a DNSBL and returns true if the argument gets a match
// in the BL.
func query(c lookuper, address string) int {
	// We ignore errors because the providers where we are not
	// flagged can't be resolved. We can not distinguish if we are
	// not on their list or their service is broken.
	result, _ := c.lookup(address)
	if len(result) > 0 {
		log.Printf("%v returned %v\n", address, result)
	}
	return len(result)
}
