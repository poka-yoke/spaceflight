package surgeon

import (
	"database/sql/driver"
	"strconv"
)

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
