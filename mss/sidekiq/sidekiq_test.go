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
	{
		Info{
			Processes: []sidekiqProcess{
				{
					Attribs: sidekiqAttribs{
						Hostname: "grape1",
					},
				},
				{
					Attribs: sidekiqAttribs{
						Hostname: "grape2",
					},
				},
				{
					Attribs: sidekiqAttribs{
						Hostname: "grape3",
					},
				},
			},
		},
		[]string{
			"grape1",
			"grape2",
			"grape3",
		},
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
