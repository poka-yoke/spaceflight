package http

import (
	"log"
	"net/http"
)

// Get does http.Get, but catching and logging errors arising.
func Get(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to get URL %s: %s", url, err)
	}
	return resp
}
