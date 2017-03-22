package trek

import (
	"fmt"
	"testing"
)

// TestIsURL tests checking for a url.
func TestIsURL(t *testing.T) {
	values := map[string]bool{
		"hostname.domain.tld":        false,
		"domain.tld":                 false,
		"*.domain.tld":               false,
		"/path/to/file":              false,
		"http://hostname.domain.tld": true,
	}
	for k, v := range values {
		if IsURL(k) != v {
			should := "should "
			if !v {
				should = "should not "
			}
			t.Errorf("%v %sbe considered a URL", k, should)
		}
	}
}

// TestIsHostname tests checking for a hostname.
func TestIsHostname(t *testing.T) {
	values := map[string]bool{
		"hostname.domain.tld":        true,
		"domain.tld":                 true,
		"*.domain.tld":               true,
		"/path/to/file":              false,
		"http://hostname.domain.tld": false,
	}
	for k, v := range values {
		if IsHostname(k) != v {
			should := "should "
			if !v {
				should = "should not "
			}
			t.Errorf("%v %sbe considered a hostname", k, should)
		}
	}
}

// TestAdd tests adding a new redirect.
func TestAdd(t *testing.T) {
	redirect := "/en /;\n"
	output, err := Add(redirect, "/en/about", "/about")
	if output != fmt.Sprintf("%s/en/about /about;\n", redirect) && err == nil {
		t.Errorf("Redirect for /en/about was not added:\n%s", output)
	}
	_, err = Add(redirect, "", "/")
	if err == nil {
		t.Errorf("Original lacking call should fail")
	}
	_, err = Add(redirect, "/en/about", "")
	if err == nil {
		t.Errorf("Final lacking call should fail")
	}
	_, err = Add(redirect, "", "")
	if err == nil {
		t.Errorf("Both original and final should be present")
	}
}
