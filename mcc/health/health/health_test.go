package health

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckCreate(t *testing.T) {
	data := []struct {
		name, schedule, channels, tags string
		slug                           string
		status                         int
	}{
		{
			name:     "New check",
			schedule: "* * * * *",
			channels: "*",
			tags:     "new minute alert",
			slug:     "f618072a-7bde-4eee-af63-71a77c5723bc",
			status:   http.StatusCreated,
		},
		{
			name:     "New check",
			schedule: "0 * * * *",
			channels: "*",
			tags:     "new hourly alert",
			slug:     "f618072a-7bde-4eee-af63-71a77c5723bc",
			status:   http.StatusCreated,
		},
		{
			name:     "New check",
			schedule: "0 * * * *",
			tags:     "new hourly noalert",
			slug:     "f618072a-7bde-4eee-af63-71a77c5723bc",
			status:   http.StatusCreated,
		},
	}
	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		// Process headers
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(
				w,
				"Undefined Content-Type",
				http.StatusBadRequest,
			)
			return
		}
		if r.Header.Get("X-Api-Key") == "" {
			http.Error(
				w,
				"Expected non-empty API Key header",
				http.StatusUnauthorized,
			)
			return
		}
		// Process body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(
				w,
				"Unexpected error reading request body",
				http.StatusInternalServerError,
			)
			return
		}
		s := make(map[string]interface{})
		err = json.Unmarshal(body, &s)
		if err != nil {
			http.Error(
				w,
				"Unexpected error parsing request body",
				http.StatusInternalServerError,
			)
			return
		}
		for _, d := range data {
			if d.tags == s["tags"] {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(d.status)
				b := map[string]interface{}{
					"grace":      60,
					"last_ping":  nil,
					"n_pings":    0,
					"name":       d.name,
					"next_ping":  nil,
					"pause_url":  "https://healthchecks.io/api/v1/checks/" + d.slug + "/pause",
					"ping_url":   "https://hchk.io/" + d.slug,
					"status":     "new",
					"tags":       d.tags,
					"timeout":    3600,
					"update_url": "https://healthchecks.io/api/v1/checks/" + d.slug,
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
	}))
	defer server.Close()

	for _, tt := range data {
		check := NewCheck()
		check.APIKey = "your-api-key"
		message := map[string]interface{}{
			"tags":     tt.tags,
			"schedule": tt.schedule,
			"name":     tt.name,
		}
		res, err := check.Create(server.URL, message)
		if err != nil {
			t.Error(err)
		}
		if res.StatusCode != tt.status {
			t.Errorf(
				"Expected status %d but received %d: %s",
				tt.status,
				res.StatusCode,
				res.Status,
			)
		}
		m, err := ParseResponse(res.Body)
		if err != nil {
			t.Error(err)
		}
		if m["name"] != tt.name {
			t.Errorf(
				"Wrong name. Expected %s, found %s",
				tt.name,
				m["name"],
			)
		}
		if m["tags"] != tt.tags {
			t.Errorf(
				"Wrong tags. Expected %s, found %s",
				tt.tags,
				m["tags"],
			)
		}
		v, ok := m["update_url"].(string)
		if !ok {
			t.Errorf(
				"Invalid field. Expected a string, received %v",
				m["udate_url"],
			)
		}
		if GetSlugFromURL(v) != tt.slug {
			t.Errorf(
				"Wrong slug. Expected %s, found %s",
				tt.slug,
				GetSlugFromURL(v),
			)
		}
		if GetSlugFromURL(tt.slug) != tt.slug {
			t.Errorf(
				"String was modified. Expected %s but is %s",
				tt.slug,
				GetSlugFromURL(tt.slug),
			)
		}
	}
}

func TestGetSlugFromURL(t *testing.T) {
	data := []struct {
		in, out string
	}{
		{
			in:  "",
			out: "",
		},
		{
			in:  "https://healthchecks.io/api/v1/checks/f618072a-7bde-4eee-af63-71a77c5723bc",
			out: "f618072a-7bde-4eee-af63-71a77c5723bc",
		},
		{
			in:  "https://healthchecks.io/api/v1/checks/f618072a-7bde-4eee-af63-71a77c5723bc/pause",
			out: "f618072a-7bde-4eee-af63-71a77c5723bc",
		},
		{
			in:  "https://healthchecks.io/api/v1/checks/f618072a-7bde-4eee-af63-71a77c5723bc/pauses",
			out: "f618072a-7bde-4eee-af63-71a77c5723bc/pauses",
		},
		{
			in:  "https://hchk.io/f618072a-7bde-4eee-af63-71a77c5723bc",
			out: "f618072a-7bde-4eee-af63-71a77c5723bc",
		},
	}
	for _, tt := range data {
		if GetSlugFromURL(tt.in) != tt.out {
			t.Errorf(
				"Wrong slug. Expected %s, found %s",
				tt.out,
				GetSlugFromURL(tt.in),
			)
		}
	}
}
