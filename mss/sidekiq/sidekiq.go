package sidekiq

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/olorin/nagiosplugin"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

type sidekiqStats struct {
	Processed      int64 `json:"processed"`
	Failed         int64 `json:"failed"`
	Busy           int64 `json:"busy"`
	Processes      int64 `json:"processes"`
	Enqueued       int64 `json:"enqueued"`
	Scheduled      int64 `json:"scheduled"`
	Retries        int64 `json:"retries"`
	Dead           int64 `json:"dead"`
	DefaultLatency int64 `json:"default_latency"`
}

type sidekiqProcess struct {
	Attribs struct {
		Hostname    string   `json:"hostname"`
		StartedAt   float64  `json:"started_at"`
		Pid         int64    `json:"pid"`
		Tag         string   `json:"tag"`
		Concurrency int64    `json:"concurrency"`
		Queues      []string `json:"queues"`
		Labels      []string `json:"labels"`
		Identity    string   `json:"identity"`
		Busy        int64    `jsno:"busy"`
		Beat        float64  `json:"beat"`
	} `json:"attribs"`
}

// Info holds all the information obtained from the sidekiq
// process
type Info struct {
	Stats     sidekiqStats
	Processes []sidekiqProcess
}

func (s Info) runningHosts() []string {
	hostnames := []string{}
	for _, v := range s.Processes {
		hostnames = append(hostnames, v.Attribs.Hostname)
	}
	return hostnames
}

// NagiosCheck returns a Nagios check populated with Sidekiq's
// information
func (s Info) NagiosCheck() *nagiosplugin.Check {
	check := nagiosplugin.NewCheck()
	must(check.AddPerfDatum("processed", "", float64(s.Stats.Processed)))
	must(check.AddPerfDatum("failed", "", float64(s.Stats.Failed)))
	must(check.AddPerfDatum("busy", "", float64(s.Stats.Busy)))
	must(check.AddPerfDatum("num_processes", "", float64(s.Stats.Processes), 1, 1, 1, 1))
	must(check.AddPerfDatum("enqueued", "", float64(s.Stats.Enqueued)))
	must(check.AddPerfDatum("scheduled", "", float64(s.Stats.Scheduled)))
	must(check.AddPerfDatum("retries", "", float64(s.Stats.Retries)))
	must(check.AddPerfDatum("dead", "", float64(s.Stats.Dead)))

	hostString := strings.Join(s.runningHosts()[:], ", ")
	if s.Stats.Processes > 1 {
		check.AddResult(nagiosplugin.CRITICAL, "Running in too many nodes. Nodes: "+hostString)
	}
	if s.Stats.Processes == 0 {
		check.AddResult(nagiosplugin.CRITICAL, "No sidekiq running")
	}
	check.AddResult(nagiosplugin.OK, "Running in node "+hostString)
	return check
}

func (s Info) getSidekiqProcessList(url string) Info {
	resp := getSidekiqData(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	processList := []sidekiqProcess{}
	if err := json.Unmarshal(body, &processList); err != nil {
		panic(err)
	}
	s.Processes = processList
	return s
}

func (s Info) getSidekiqStats(url string) Info {
	resp := getSidekiqData(url)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	stats := sidekiqStats{}
	if err := json.Unmarshal(body, &stats); err != nil {
		panic(err)
	}
	s.Stats = stats
	return s
}

func getSidekiqData(url string) *http.Response {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	return resp
}

// ProcessGetResponse creates a sidekiqInfo object contining sidekiq's
// exposed API information available
func ProcessGetResponse(baseURL string) (info Info) {
	info = Info{}
	info = info.getSidekiqStats(baseURL + "/system/sidekiq")
	info = info.getSidekiqProcessList(baseURL + "/system/sidekiq/processes")
	return
}
