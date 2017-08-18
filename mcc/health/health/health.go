package health

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

// Check represents the checking service and holds service dependent
// data
type Check struct {
	APIKey string
}

// Create adds a new check throuh an API call
func (c *Check) Create(endpoint string, message map[string]interface{}) (res *http.Response, err error) {
	// Request body
	buf, err := json.Marshal(message)
	if err != nil {
		return
	}
	req, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(buf),
	)
	if err != nil {
		return
	}
	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", c.APIKey)
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	return
}

// GetSlugFromURL extracts the slug from any of the returned URLs
func GetSlugFromURL(m string) (s string) {
	s = strings.TrimPrefix(
		m,
		"https://healthchecks.io/api/v1/checks/",
	)
	s = strings.TrimPrefix(s, "https://hchk.io/")
	s = strings.TrimSuffix(s, "/pause")
	return
}

// ParseResponse converts the services' response body into a map
func ParseResponse(in io.ReadCloser) (m map[string]interface{}, err error) {
	m = make(map[string]interface{})
	body, err := ioutil.ReadAll(in)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return
	}
	return
}

// NewCheck creates an empty Check
func NewCheck() *Check {
	return &Check{}
}
