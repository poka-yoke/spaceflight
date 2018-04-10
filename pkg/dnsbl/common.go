package dnsbl

import (
	"bufio"
	"fmt"
	"io"
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
