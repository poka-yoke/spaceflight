package digest

import (
	"crypto/sha1"
	"encoding/base64"
	"io"
)

// ContentBase64 is providing a Base64 encoded SHA1 hash from a io.Reader.
func ContentBase64(input io.Reader) (encoded string, err error) {
	hasher := sha1.New()
	if _, err = io.Copy(hasher, input); err != nil {
		return "", err
	}
	encoded = base64.StdEncoding.EncodeToString(hasher.Sum(nil))
	return encoded, nil
}
