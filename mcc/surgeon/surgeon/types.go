package surgeon

import (
	"database/sql/driver"
	"errors"
	"strconv"
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

// Duration represents a `duration` type
type Duration struct {
	Duration float64
}

// String returns the string value of the Duration.
func (duration *Duration) String() string {
	return strconv.FormatFloat(duration.Duration, 'f', 3, 64)
}

// MarshalText marshals duration into text.
func (duration *Duration) MarshalText() ([]byte, error) {
	return []byte(duration.String()), nil
}

// UnmarshalText unmarshals Duration from text.
func (duration *Duration) UnmarshalText(text []byte) error {
	convertedDuration, _ := strconv.ParseFloat(string(text), 64)
	duration.Duration = convertedDuration
	return nil
}

// Value satisfies the sql/driver.Valuer interface for Duration.
func (duration *Duration) Value() (driver.Value, error) {
	return duration.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Duration.
func (duration *Duration) Scan(src interface{}) error {
	buf := ""
	floatDuration, ok := src.(float64)
	if ok {
		buf = strconv.FormatFloat(floatDuration, 'E', -1, 64)
	}
	return duration.UnmarshalText(([]byte)(buf))
}

// Seconds return the quantity of seconds this duration represents.
func (duration *Duration) Seconds() (amount float64) {
	if duration.String() == "" {
		return 0.0
	}
	return duration.Duration
}
