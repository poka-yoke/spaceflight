package dnsbl

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Reverse reverses slice of string elements.
func reverse(original []string) {
	for i := len(original)/2 - 1; i >= 0; i-- {
		opp := len(original) - 1 - i
		original[i], original[opp] = original[opp], original[i]
	}
}

// ReverseAddress converts IP address in string to reversed address for query.
func reverseAddress(ipAddress string) (reversedIPAddress string) {
	ipAddressValues := strings.Split(ipAddress, ".")
	reverse(ipAddressValues)
	reversedIPAddress = strings.Join(ipAddressValues, ".")
	return
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
