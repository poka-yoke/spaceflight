package cloudinary

import (
	"errors"
)

// CloudinaryCredentials is the global variable for the Cloudinary's account's
// credentials.
var CloudinaryCredentials *Credentials

// Credentials represents a set of Cloudinary credentials.
type Credentials struct {
	cloudName string
	key       string
	secret    string
}

// NewCredentials creates the Credentials object.
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
