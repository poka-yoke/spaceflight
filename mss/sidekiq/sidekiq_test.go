package sidekiq

import "testing"

var rhtable = []struct {
	Info
	output []string
}{
	{
		Info{},
		[]string{},
	},
}

func TestRunningHosts(t *testing.T) {
	for _, tt := range rhtable {
		for k, v := range tt.runningHosts() {
			if v != tt.output[k] {
				t.Error(
					"Should be equal. Got",
					v,
					"expected",
					tt.output[k],
				)
			}
		}
	}
}
