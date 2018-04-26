package cloudinary

import (
	"errors"
)

var CloudinaryCredentials *Credentials

type Credentials struct {
	cloudName string
	key       string
	secret    string
}

func NewCredentials(cloudName, key, secret string) error {
	if key == "" || secret == "" || cloudName == "" {
		return errors.New("No credentials defined")
	}
	CloudinaryCredentials = &Credentials{
		cloudName: cloudName,
		key:       key,
		secret:    secret,
	}
	return nil
}

func (c *Credentials) get() (string, string, string) {
	return c.cloudName, c.key, c.secret
}
