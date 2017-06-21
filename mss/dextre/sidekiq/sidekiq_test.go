package sidekiq

import (
	"encoding/json"
	"testing"
)

var gsstable = []struct {
	s      Info
	output []string
}{
	{
		Info{},
		[]string{},
	},
	{
		Info{
			Stats: sidekiqStats{
				Processed:      0,
				Failed:         0,
				Busy:           0,
				Processes:      0,
				Enqueued:       0,
				Scheduled:      0,
				Retries:        0,
				Dead:           0,
				DefaultLatency: 0.0,
			},
			Processes: []sidekiqProcess{
				{
					Attribs: sidekiqAttribs{
						Hostname: "node1",
						Queues: []string{
							"my_queue",
							"also_this_queue",
						},
						Pid:         1,
						Tag:         "init",
						Concurrency: 5,
						Labels: []string{
							"I",
							"Have",
							"Labels",
							"Too!",
						},
						Identity: "just-me",
						Busy:     3,
						Beat:     165.8,
					},
				},
				{
					Attribs: sidekiqAttribs{
						Hostname: "node2",
						Queues: []string{
							"my_other_queue",
						},
					},
				},
				{
					Attribs: sidekiqAttribs{
						Hostname: "node3",
					},
				},
			},
		},
		[]string{
			"node1",
			"node2",
			"node3",
		},
	},
	{
		Info{
			Stats: sidekiqStats{
				Processed:      1000,
				Failed:         500,
				Busy:           1,
				Processes:      1,
				Enqueued:       50,
				Scheduled:      0,
				Retries:        10,
				Dead:           300,
				DefaultLatency: 1.50,
			},
			Processes: []sidekiqProcess{
				{
					Attribs: sidekiqAttribs{
						Hostname: "node1",
						Queues: []string{
							"my_other_queue",
						},
					},
				},
				{
					Attribs: sidekiqAttribs{
						Hostname: "node2",
					},
				},
			},
		},
		[]string{
			"node1",
			"node2",
		},
	},
	{
		Info{
			Stats: sidekiqStats{
				Processed:      0,
				Failed:         0,
				Busy:           0,
				Processes:      0,
				Enqueued:       0,
				Scheduled:      0,
				Retries:        0,
				Dead:           0,
				DefaultLatency: 1,
			},
			Processes: []sidekiqProcess{
				{
					Attribs: sidekiqAttribs{
						Hostname: "node3",
						Queues: []string{
							"I_didnt_know_of_this_queue",
						},
					},
				},
			},
		},
		[]string{
			"node3",
		},
	},
}

func TestRunningHosts(t *testing.T) {
	for _, tt := range gsstable {
		for k, v := range tt.s.runningHosts() {
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

func TestReadSidekiqStats(t *testing.T) {
	for _, tt := range gsstable {
		body, err := json.Marshal(tt.s.Stats)
		if err != nil {
			panic(err)
		}
		info := Info{}
		info = info.readSidekiqStats(body)

		if tt.s.Stats != info.Stats {
			t.Error("Should have been equal")
		}
	}
}

func TestReadSidekiqProcessList(t *testing.T) {
	for _, tt := range gsstable {
		body, err := json.Marshal(tt.s.Processes)
		if err != nil {
			panic(err)
		}
		info := Info{}
		info = info.readSidekiqProcessList(body)

		if !compareSidekiqProcess(info.Processes, tt.s.Processes) {
			t.Error("Should have been equal")
		}
	}
}

func compareSidekiqProcess(a, b []sidekiqProcess) bool {
	if len(a) != len(b) {
		return false
	}
	res := true
	for k, v := range a {
		if !res || !compareSidekiqAttribs(v.Attribs, b[k].Attribs) {
			res = false
		}
	}
	return res
}

func compareSidekiqAttribs(a, b sidekiqAttribs) bool {
	if a.Hostname != b.Hostname ||
		a.StartedAt != b.StartedAt ||
		a.Pid != b.Pid ||
		a.Tag != b.Tag ||
		a.Concurrency != b.Concurrency ||
		a.Identity != b.Identity ||
		a.Busy != b.Busy ||
		a.Beat != b.Beat {
		return false
	}
	if len(a.Queues) != len(b.Queues) ||
		len(a.Labels) != len(b.Labels) {
		return false
	}

	res := true
	for k, v := range a.Labels {
		if !res || v != b.Labels[k] {
			res = false
		}
	}
	for k, v := range a.Queues {
		if !res || v != b.Queues[k] {
			res = false
		}
	}
	return res
}
