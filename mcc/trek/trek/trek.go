package trek

import "fmt"

// Add adds a new redirect line for original to final, and returns it.
func Add(redirects, original, final string) (resultingRedirects string, err error) {
	if original == "" || final == "" {
		err = fmt.Errorf("You must specify both original and final")
		return
	}
	newRedirect := fmt.Sprintf("%s %s;\n", original, final)
	resultingRedirects = fmt.Sprintf("%s%s", redirects, newRedirect)
	return
}
