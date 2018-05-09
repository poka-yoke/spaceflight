package cronitor

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

// SetAPIKey configures the backend with proper API Key
func (c *Check) SetAPIKey(key string) {
	c.APIKey = key
}

// Create adds a new check throuh an API call
func (c *Check) Create(endpoint string, message map[string]interface{}) (req *http.Request, err error) {
	// Request body
	buf, err := json.Marshal(message)
	if err != nil {
		return
	}
	req, err = http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewBuffer(buf),
	)
	if err != nil {
		return
	}
	// Set headers
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(c.APIKey, "")
	return
}

// SetMessage builds the request body for cronitor.io
func SetMessage(schedule, name, tags, email string) map[string]interface{} {
	message := make(map[string]interface{})
	message["type"] = "heartbeat"
	if schedule != "" {
		message["rules"] = []map[string]interface{}{
			{
				"value":     schedule,
				"rule_type": "not_on_schedule",
			},
		}
	}
	if name != "" {
		message["name"] = name
	}
	if tags != "" {
		message["tags"] = strings.Split(tags, " ")
	}
	if email != "" {
		message["notifications"] = map[string][]string{
			"emails": []string{
				email,
			},
		}
	}
	return message
}

// ParseResponse converts the services' response body into a map
func ParseResponse(in io.ReadCloser) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	body, err := ioutil.ReadAll(in)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// NewCheck creates an empty Check
func NewCheck() *Check {
	return &Check{}
}
