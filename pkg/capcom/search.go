package capcom

import (
	"fmt"
	"strconv"
)

// SearchResult defines a result for a rule
type SearchResult struct {
	GroupID  string
	Protocol string
	Port     int64
	Source   string
}

// String method for SearchResult gets a String to be printed.
func (sr SearchResult) String() string {
	return fmt.Sprintf(
		"%s %s/%s %s",
		sr.GroupID,
		strconv.FormatInt(sr.Port, 10),
		sr.Protocol,
		sr.Source,
	)
}
