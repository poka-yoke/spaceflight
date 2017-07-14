package short

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func NewService() *Service {
	return &Service{}
}

type Service struct {
	aPIKey            string
	AddURL, UpdateURL string
}

func (s *Service) SetAPIKey(key string) *Service {
	s.aPIKey = key
	return s
}

// Add an entry to the shortener with handle path and destination url
func (s Service) Add(domain, path, url string) error {
	body := s.Body(domain, path, url)
	req, err := s.Request(http.MethodPost, body)
	if err != nil {
		return err
	}
	res, err := s.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 201 {
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

// Request to add entry
func (s Service) Request(method string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, s.AddURL, body)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", s.aPIKey)
	return
}

// Do executes the request to add an entry
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
