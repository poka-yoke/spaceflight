package surgeon

import (
	"database/sql/driver"
	"errors"
)

// Name represents a `name` type from PostgreSQL
type Name struct {
	Name string
}

// String returns the string value of the Name.
func (name Name) String() string {
	return string(name.Name)
}

// MarshalText marshals name into text.
func (name Name) MarshalText() ([]byte, error) {
	return []byte(name.String()), nil
}

// UnmarshalText unmarshals Name from text.
func (name *Name) UnmarshalText(text []byte) error {
	name.Name = string(text)
	return nil
}

// Value satisfies the sql/driver.Valuer interface for Name.
func (name Name) Value() (driver.Value, error) {
	return name.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Name.
func (name *Name) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid name")
	}

	return name.UnmarshalText(buf)
}
