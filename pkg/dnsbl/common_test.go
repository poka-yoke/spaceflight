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
