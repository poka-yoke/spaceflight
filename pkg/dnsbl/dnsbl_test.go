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
	Lookup = mockLookup
	checker := &Checker{}
	results := make(chan int)
	t.Log("Hola")
	go checker.query("127.0.0.1", "positive.dnsbl.com", results)
	result := <-results
	if result != 1 {
		t.Errorf(
			"Result for 1.0.0.127.positive.dnsbl.com "+
				"should be 1 instead of %v",
			result,
		)
	}
	go checker.query("127.0.0.1", "negative.dnsbl.com", results)
	result = <-results
	if result != 0 {
		t.Errorf(
			"Result for 1.0.0.127.negative.dnsbl.com "+
				"should be 0 instead of %v",
			result,
		)
	}
	close(results)
}
