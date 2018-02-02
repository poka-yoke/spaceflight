package short

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// NewService returns an empty Service
func NewService() *Service {
	return &Service{}
}

// Service holds identification and back-end methods
type Service struct {
	aPIKey            string
	AddURL, UpdateURL string
}

// SetAPIKey Configures the API Key in the Service
func (s *Service) SetAPIKey(key string) *Service {
	s.aPIKey = key
	return s
}

// Add an entry to the shortener with handle path and destination url
func (s Service) Add(domain, path, url string) error {
	body := s.Body(domain, path, url)
	req, err := s.Request(http.MethodPost, s.AddURL, body)
	if err != nil {
		return err
	}
	res, err := s.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return &shortenError{
			fmt.Sprintf(
				"Unexpected status, expected 201, received %d",
				res.StatusCode,
			),
			res.StatusCode,
		}
	}
	return nil
}

// Update modifies an existing short entry and changes where it points to
func (s Service) Update(domain, path, url string) error {
	body := s.Body(domain, path, url)
	req, err := s.Request(http.MethodPut, s.UpdateURL, body)
	if err != nil {
		return fmt.Errorf("Failed building request: %s", err.Error())
	}
	res, err := s.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return &shortenError{
			fmt.Sprintf(
				"Unexpected status, expected 200, received %d",
				res.StatusCode,
			),
			res.StatusCode,
		}
	}
	return nil
}

// Body renders the Request's body
func (s Service) Body(domain, path, url string) (out io.Reader) {
	message := make(map[string]string)
	message["path"] = path
	message["originalURL"] = url
	message["domain"] = domain

	buf, err := json.Marshal(message)
	if err != nil {
		return
	}
	out = bytes.NewBuffer(buf)
	return
}

// Request to Service
func (s Service) Request(method, endpoint string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, endpoint, body)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.aPIKey)
	return
}

// Do sends the request to the Service and returns the response for validation
func (s Service) Do(req *http.Request) (res *http.Response, err error) {
	res, err = http.DefaultClient.Do(req)
	if err != nil {
		err = fmt.Errorf("Failed executing request: %s", err.Error())
		return
	}
	return
}

type shortenError struct {
	s    string
	code int
}

func (e *shortenError) Error() string {
	return e.s
}

func (e *shortenError) ErrorCode() int {
	return e.code
}
