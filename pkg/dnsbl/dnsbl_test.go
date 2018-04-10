package dnsbl

import "testing"

func mockLookup(blacklist string) (addresses []string, err error) {
	if blacklist == "1.0.0.127.positive.dnsbl.com" {
		addresses = []string{"127.0.0.1"}
	} else {
		addresses = []string{}
	}
	return
}

func TestQuery(t *testing.T) {
	var table = []struct {
		query  string
		result int
	}{
		{
			query:  "1.0.0.127.positive.dnsbl.com",
			result: 1,
		},
		{
			query:  "1.0.0.127.negative.dnsbl.com",
			result: 0,
		},
	}
	checker := &Checker{lookup: mockLookup}
	for _, tc := range table {
		result := checker.query(tc.query)
		if result != tc.result {
			t.Errorf(
				"Result for %s should be %d instead of %d",
				tc.query,
				tc.result,
				result,
			)
		}
	}
}
