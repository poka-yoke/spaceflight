package dnsbl_test

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/poka-yoke/spaceflight/pkg/dnsbl"
)

func mockLookupHost(resolv *net.Resolver, ctx context.Context, host string) ([]string, error) {
	var addresses []string
	if host == "1.0.0.127.positive.dnsbl.com" {
		addresses = []string{"127.0.0.1"}
	} else {
		addresses = []string{}
	}
	return addresses, nil
}


func TestQuery(t *testing.T) {
	t.Parallel()

	var table = []struct {
		address, provider string
		result  int
	}{
		{
			address: "127.0.0.1",
			provider: "positive.dnsbl.com",
			result:  1,
		},
		{
			address: "127.0.0.1",
			provider: "negative.dnsbl.com",
			result:  0,
		},
	}
	dnsbl.Resolver = nil // Do not resolve external addresses
	dnsbl.LookupHostFunc = mockLookupHost
	for _, tc := range table {
		t.Run(
			fmt.Sprintf("%s %s - %d", tc.provider, tc.address, tc.result),
			func(t *testing.T) {
				result := dnsbl.Query(context.Background(), tc.provider, tc.address)
				if result != tc.result {
					t.Errorf(
						"Result for %s should be %d instead of %d",
						tc.address,
						tc.result,
						result,
					)
				}
			},
		)
	}
}
