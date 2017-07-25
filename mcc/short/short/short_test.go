package short

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func ExampleServiceAdd() {
	service := NewService().SetAPIKey("my-api-key")
	service.AddURL = "https://example.com/shortener/add"
	err := service.Add("xample.cm", "this", "https://example.com/very/long/url")
	if err != nil {
		// Handle error
	}
}

func ExampleServiceUpdate() {
	service := NewService().SetAPIKey("my-api-key")
	service.UpdateURL = "https://example.com/shortener/add"
	err := service.Update("xample.cm", "this", "https://example.com/very/long/url")
	if err != nil {
		// Handle error
	}
}

func TestAddShortEntry(t *testing.T) {
	data := []struct {
		domain, originalURL, path, title string
		tags                             []string
		status                           int
	}{
		{
			domain:      "shrt.co",
			originalURL: "http://yourlongdomain.com/yourlonglink",
			path:        "correct",
			title:       "Some url title",
			tags: []string{
				"tag1",
				"tag2",
			},
			status: http.StatusOK, // API docs are wrong,
			// it retugrns a 200
		},
		{
			domain:      "shrt.co",
			originalURL: "http://yourlongdomain.com/yourlonglink",
			path:        "duplicated",
			title:       "Some url title",
			tags: []string{
				"tag1",
				"tag2",
			},
			status: 409,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		// Process body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error()
		}
		s := make(map[string]string)
		err = json.Unmarshal(body, &s)
		if err != nil {
			t.Error()
		}
		for _, d := range data {
			if d.path == s["path"] {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(d.status)
				fmt.Fprintln(w, "")
			}
		}
	}))
	defer server.Close()

	service := NewService()
	service.AddURL = server.URL
	for _, tt := range data {
		err := service.Add(tt.domain, tt.path, tt.originalURL)
		if err != nil {
			if e, ok := err.(*shortenError); !ok {
				t.Error(err)
			} else {
				if e.ErrorCode() != tt.status {
					t.Errorf(
						"Unexpected response status, expected %d but received %d",
						tt.status,
						e.ErrorCode(),
					)
				}
			}
		}
	}
}

func TestAddShortEntryReq(t *testing.T) {
	data := []struct {
		domain, path, url   string
		apikey, contentType string
		apiPass             bool
	}{
		{
			domain:      "dvx.cm",
			path:        "correct",
			url:         "",
			contentType: "application/json",
			apikey:      "your-api-key",
			apiPass:     true,
		},
		{
			domain:      "example.com",
			path:        "correct",
			url:         "",
			contentType: "application/json",
			apikey:      "your-api-key",
			apiPass:     true,
		},
		{
			domain:      "dvx.cm",
			path:        "noAPIKey",
			url:         "",
			contentType: "application/json",
			apikey:      "",
			apiPass:     false,
		},
		{
			domain:      "dvx.cm",
			path:        "wrongAPIKey",
			url:         "",
			contentType: "application/json",
			apikey:      "this-is-not-an-api-key",
			apiPass:     false,
		},
	}

	for _, tt := range data {
		service := NewService().
			SetAPIKey(tt.apikey)

		bod := service.Body(tt.domain, tt.path, tt.url)
		r, err := service.Request(http.MethodPost, service.AddURL, bod)
		if err != nil {
			t.Error(err)
		}

		// Verify request headers
		t.Run("Content-Type", func(t *testing.T) {
			if r.Header.Get("Content-Type") != tt.contentType {
				t.Errorf("Content-Type equals %s", r.Header.Get("Content-Type"))
			}
		})
		t.Run("Authorization", func(t *testing.T) {
			if r.Header.Get("Authorization") != tt.apikey && tt.apiPass {
				t.Errorf(
					"Wrong API key, expected %s got %s",
					tt.apikey,
					r.Header.Get("Authorization"),
				)
			}
		})

		// Verify request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error()
		}
		s := make(map[string]string)
		err = json.Unmarshal(body, &s)
		if err != nil {
			t.Error()
		}
		t.Run("Domain field", func(t *testing.T) {
			if s["domain"] != tt.domain {
				t.Errorf(
					"Expected domain %s, got %s",
					tt.domain,
					s["domain"],
				)
			}
		})
		t.Run("Path field", func(t *testing.T) {
			if s["path"] != tt.path {
				t.Errorf(
					"Expected path %s, got %s",
					tt.path,
					s["path"],
				)
			}
		})
		t.Run("URL field", func(t *testing.T) {
			if s["originalURL"] != tt.url {
				t.Errorf(
					"Expected originalURL %s, got %s",
					tt.url,
					s["originalURL"],
				)
			}
		})
	}
}

func TestUpdateShortEntry(t *testing.T) {
	data := []struct {
		domain, originalURL, path, title string
		tags                             []string
		status                           int
	}{
		{
			domain:      "shrt.co",
			originalURL: "http://yourlongdomain.com/yourlonglink",
			path:        "correct",
			title:       "Some url title",
			tags: []string{
				"tag1",
				"tag2",
			},
			status: http.StatusOK,
		},
		{
			domain:      "shrt.co",
			originalURL: "http://yourlongdomain.com/yourlonglink",
			path:        "duplicated",
			title:       "Some url title",
			tags: []string{
				"tag1",
				"tag2",
			},
			status: http.StatusOK,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		// Process body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			t.Error()
		}
		s := make(map[string]string)
		err = json.Unmarshal(body, &s)
		if err != nil {
			t.Error()
		}
		for _, d := range data {
			if d.path == s["path"] {
				w.Header().Set(
					"Content-Type",
					"application/json",
				)
				w.WriteHeader(d.status)
				fmt.Fprintln(w, "")
			}
		}
	}))
	defer server.Close()

	service := NewService()
	service.UpdateURL = server.URL
	for _, tt := range data {
		err := service.Update(tt.domain, tt.path, tt.originalURL)
		if err != nil {
			if e, ok := err.(*shortenError); !ok {
				t.Error(err)
			} else {
				if e.ErrorCode() != tt.status {
					t.Errorf(
						"Unexpected response status, expected %d but received %d",
						tt.status,
						e.ErrorCode(),
					)
				}
			}
		}
	}
}
