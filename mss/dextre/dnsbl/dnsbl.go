package dnsbl

import (
	"fmt"
	"log"
	"net"
	"strings"
)

// Lookup contains the lookup function used
var Lookup = net.LookupHost

// Reverse reverses slice of string elements.
func Reverse(original []string) {
	for i := len(original)/2 - 1; i >= 0; i-- {
		opp := len(original) - 1 - i
		original[i], original[opp] = original[opp], original[i]
	}
}

// ReverseAddress converts IP address in string to reversed address for query.
func ReverseAddress(ipAddress string) (reversedIPAddress string) {
	ipAddressValues := strings.Split(ipAddress, ".")
	Reverse(ipAddressValues)
	reversedIPAddress = strings.Join(ipAddressValues, ".")
	return
}

// Query queries a DNSBL and returns true if the argument gets a match
// in the BL.
func Query(ipAddress, bl string, addresses chan int) {
	reversedIPAddress := fmt.Sprintf(
		"%v.%v",
		ReverseAddress(ipAddress),
		bl,
	)
	result, _ := Lookup(reversedIPAddress)
	if len(result) > 0 {
		log.Printf("%v present in %v(%v)", reversedIPAddress, bl, result)
	}
	addresses <- len(result)
}
