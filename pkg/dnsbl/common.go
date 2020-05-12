package dnsbl

import (
	"context"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	// LookupHostFunc allows to override the function used for
	// Host lookup. Defaults to the standard library's
	// implementation
	LookupHostFunc = (*net.Resolver).LookupHost
	// Resolver allows to override the DNS resolver. Defaults to
	// the standard library's default resolver.
	Resolver = net.DefaultResolver
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

// Query queries a DNSBL and returns the number of matches present in
// the BL. If the number of matches in the lookup were 0, it means not
// present.
func Query(ctx context.Context, provider, address string) int {
	queryAddress := fmt.Sprintf(
		"%v.%v",
		reverseAddress(address),
		provider,
	)

	// We ignore errors because the providers where we are not
	// flagged can't be resolved. We can not distinguish if we are
	// not on their list or their service is broken.
	result, _ := LookupHostFunc(Resolver, ctx, queryAddress)
	if len(result) > 0 {
		log.Printf("%v returned %v\n", queryAddress, result)
	}
	return len(result)
}
