package dnsbl

import (
	"fmt"
	"testing"
)

func TestReverseAddress(t *testing.T) {
	t.Parallel()

	table := []struct {
		original, reverse string
	}{
		{
			original: "127.0.0.1",
			reverse:  "1.0.0.127",
		},
		{
			original: "1.2.3.4",
			reverse:  "4.3.2.1",
		},
	}
	for _, tc := range table {
		t.Run(
			fmt.Sprintf("%s - %s", tc.original, tc.reverse),
			func(t *testing.T) {
				reversedIPAddress := reverseAddress(tc.original)
				if reversedIPAddress != tc.reverse {
					t.Errorf(
						"reversedIPAddress should be %s and not %s",
						tc.reverse,
						reversedIPAddress,
					)
				}
			},
		)
	}
}

type mockLookup struct{}

func (m mockLookup) lookup(address string) ([]string, error) {
	var addresses []string
	if address == "1.0.0.127.positive.dnsbl.com" {
		addresses = []string{"127.0.0.1"}
	} else {
		addresses = []string{}
	}
	return addresses, nil
}

func TestQuery(t *testing.T) {
	t.Parallel()

	var table = []struct {
		address string
		result  int
	}{
		{
			address: "1.0.0.127.positive.dnsbl.com",
			result:  1,
		},
		{
			address: "1.0.0.127.negative.dnsbl.com",
			result:  0,
		},
	}
	checker := mockLookup{}
	for _, tc := range table {
		t.Run(
			fmt.Sprintf("%s - %d", tc.address, tc.result),
			func(t *testing.T) {
				result := query(checker, tc.address)
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
