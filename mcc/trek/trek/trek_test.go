package trek

import (
	"fmt"
	"testing"
)

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
