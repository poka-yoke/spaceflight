package cronitor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type healthCase struct {
	name, schedule, channels, tags string
	slug, email, note              string
	status                         int
}

var healthCases = []healthCase{
	{
		name:     "New check",
		schedule: "* * * * *",
		channels: "*",
		email:    "test@example.com",
		tags:     "new minute alert",
		slug:     "f618072a-7bde-4eee-af63-71a77c5723bc",
		note:     "new minute alert",
		status:   http.StatusCreated,
	},
	{
		name:     "New check",
		schedule: "0 * * * *",
		channels: "*",
		email:    "test@example.com",
		tags:     "new hourly alert",
		slug:     "f618072a-7bde-4eee-af63-71a77c5723bc",
		note:     "new hourly alert",
		status:   http.StatusCreated,
	},
	{
		name:     "New check",
		schedule: "0 * * * *",
		email:    "test@example.com",
		tags:     "new hourly noalert",
		slug:     "f618072a-7bde-4eee-af63-71a77c5723bc",
		note:     "new hourly noalert",
		status:   http.StatusCreated,
	},
}

func healthCheckServer(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		// Process headers
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf(
				"Expected header to have Content-Type %s, received %s",
				"application/json",
				r.Header.Get("Content-Type"),
			)
		}
		if r.Header.Get("Authorization") == "" {
			t.Errorf("Expected header to have Authorization")
		}
		// Process body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		s := make(map[string]interface{})
		err = json.Unmarshal(body, &s)
		if err != nil {
			t.Fatal(err)
		}
		for _, d := range healthCases {
			if d.note == s["note"] {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(d.status)
				b := map[string]interface{}{
					"code":        d.slug,
					"key":         nil,
					"type":        "heartbeat",
					"dev":         false,
					"initialized": false,
					"disabled":    false,
					"paused":      false,
					"passing":     true,
					"running":     false,
					"status":      "Waiting for the first ping",
					"name":        d.name,
					"notifications": map[string][]string{
						"templates": []string{},
						"pagerduty": []string{},
						"slack":     []string{},
						"phones":    []string{},
						"emails":    []string{d.email},
						"webhooks":  []string{},
						"hipchat":   []string{},
					},
					"tags":     strings.Split(d.tags, " "),
					"timezone": nil,
					"rules": []map[string]interface{}{
						{
							"rule_type":               "not_on_schedule",
							"hours_to_followup_alert": 8,
							"value":                   d.schedule,
						},
					},
					"request":                  nil,
					"request_interval_seconds": nil,
					"note":                 d.note,
					"has_duration_history": false,
					"created":              time.Now().Format("2006-01-02T15:04:05-07:00"),
				}
				buf, err := json.Marshal(b)
				if err != nil {
					t.Error(err)
				}
				fmt.Fprintln(w, bytes.NewBuffer(buf))
				return
			}
		}
		http.Error(w, "Bad Request", http.StatusBadRequest)
	}
}

func TestCheckCreate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(healthCheckServer(t)))
	defer server.Close()

	for _, tc := range healthCases {
		check := NewCheck()
		check.APIKey = "your-api-key"
		message := map[string]interface{}{
			"tags": strings.Split(tc.tags, " "),
			"name": tc.name,
			"type": "heartbeat",
			"notifications": map[string][]string{
				"emails": []string{
					tc.email,
				},
			},
			"rules": []map[string]interface{}{
				{
					"rule_type": "not_on_schedule",
					"value":     tc.schedule,
				},
			},
			"note": tc.note,
		}
		res, err := check.Create(server.URL, message)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != tc.status {
			t.Errorf(
				"Expected status %d but received %d: %s",
				tc.status,
				res.StatusCode,
				res.Status,
			)
		}
		m, err := ParseResponse(res.Body)
		if err != nil {
			t.Error(err)
		}
		if m["name"] != tc.name {
			t.Errorf(
				"Wrong name. Expected %s, found %s",
				tc.name,
				m["name"],
			)
		}
		tmp, ok := m["tags"].([]interface{})
		if !ok {
			t.Error("Tags field is not slice of interfaces")
		}
		var tags []string
		for _, v := range tmp {
			s, ok := v.(string)
			if !ok {
				t.Error("Tag contents are not strings")
			}
			tags = append(tags, s)
		}
		if strings.Join(tags, " ") != tc.tags {
			t.Errorf(
				"Wrong tags. Expected %s, found %v",
				tc.tags,
				strings.Join(tags, " "),
			)
		}
	}
}
