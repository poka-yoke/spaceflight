package surgeon

import (
	"database/sql/driver"
	"errors"
)

// Inet represents an `inet` type from PostgreSQL
type Inet struct {
	IP string
}

// String returns the string value of the Inet.
func (inet Inet) String() string {
	return string(inet.IP)
}

// MarshalText marshals inet into text.
func (inet Inet) MarshalText() ([]byte, error) {
	return []byte(inet.String()), nil
}

// UnmarshalText unmarshals Inet from text.
func (inet *Inet) UnmarshalText(text []byte) error {
	inet.IP = string(text)
	return nil
}

// Value satisfies the sql/driver.Valuer interface for Inet.
func (inet Inet) Value() (driver.Value, error) {
	return inet.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Inet.
func (inet *Inet) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		if src != nil {
			return errors.New("invalid inet")
		}
		buf = []byte("")
	}
	return inet.UnmarshalText(buf)
}
