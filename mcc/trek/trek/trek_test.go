package trek

import (
	"fmt"
	"testing"
)

// TestAdd tests adding a new redirect.
func TestAdd(t *testing.T) {
	redirect := "/en /;\n"
	output := Add(redirect, "/en/about", "/about")
	if output != fmt.Sprintf("%s/en/about /about;\n", redirect) {
		t.Errorf("Redirect for /en/about was not added:\n%s", output)
	}
}
