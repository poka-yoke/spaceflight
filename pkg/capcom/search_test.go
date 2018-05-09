package capcom

import "testing"

func TestString(t *testing.T) {
	sr := SearchResult{
		GroupID:  "sg-idsgtest",
		Protocol: "tcp",
		Port:     22,
		Source:   "0.0.0.0/0",
	}
	expected := "sg-idsgtest 22/tcp 0.0.0.0/0"
	result := sr.String()
	if expected != result {
		t.Errorf("%s is not %s", result, expected)
	}
}
