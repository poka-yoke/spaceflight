package dnsbl

import "testing"

func TestReverse(t *testing.T) {
	stringSlice := []string{
		"1",
		"2",
		"3",
	}
	reverse(stringSlice)
	if stringSlice[0] != "3" {
		t.Errorf(
			"stringSlice[0] should be 3 and not %s",
			stringSlice[0],
		)
	}
	if stringSlice[1] != "2" {
		t.Errorf(
			"stringSlice[1] should be 2 and not %s",
			stringSlice[1],
		)
	}
	if stringSlice[2] != "1" {
		t.Errorf(
			"stringSlice[2] should be 1 and not %s",
			stringSlice[2],
		)
	}
}

func TestReverseAddress(t *testing.T) {
	ipAddress := "127.0.0.1"
	reversedIPAddress := reverseAddress(ipAddress)
	if reversedIPAddress != "1.0.0.127" {
		t.Errorf(
			"reversedIPAddress should be 1.0.0.127 and not %s",
			reversedIPAddress,
		)
	}
}

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
