package dnsbl

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strings"
)

// Reverse returns a slice of string elements in reverse order than
// the one supplied.
func reverse(original []string) []string {
	copy := original
	for i := len(original)/2 - 1; i >= 0; i-- {
		opp := len(original) - 1 - i
		copy[i], copy[opp] = original[opp], original[i]
	}
	return copy
}

// ReverseAddress converts IP address into reversed address for query.
func reverseAddress(ipAddress string) string {
	ipAddressValues := strings.Split(ipAddress, ".")
	return strings.Join(reverse(ipAddressValues), ".")
}

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
