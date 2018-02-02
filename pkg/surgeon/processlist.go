package surgeon

import (
	"fmt"
	"sort"
)

// ProcessList represents a slice of Process.
type ProcessList struct {
	Processes []Process
	Inactive  bool
}

// Len returns the length of a ProcessList.
func (pl ProcessList) Len() int {
	return len(pl.Processes)
}

// Less returns whether a Process is less than another in a ProcessList.
func (pl ProcessList) Less(i, j int) bool {
	firstDuration := pl.Processes[i].Duration.Seconds()
	secondDuration := pl.Processes[j].Duration.Seconds()
	return firstDuration < secondDuration
}

// Swap exchanges two Processes in the ProcessList.
func (pl ProcessList) Swap(i, j int) {
	pl.Processes[i], pl.Processes[j] = pl.Processes[j], pl.Processes[i]
}

// String provides a string formatting to ProcessList, which depends on
// Inactive field value.
func (pl ProcessList) String() (output string) {
	output = fmt.Sprintf(
		" %5s %10s %19s %17s %13s %9s %s \n",
		"PID",
		"USER",
		"DB",
		"CLIENT",
		"DURATION",
		"WAITING",
		"QUERY",
	)
	sort.Sort(pl)
	for _, process := range pl.Processes {
		if pl.Inactive || process.IsActive() {
			output += fmt.Sprintf(
				" %5d %12s %20s %17s %10.2f %9s %s \n",
				process.Pid,
				process.Usename,
				process.Datname,
				process.Client.String(),
				process.Duration.Duration,
				process.Waiting,
				process.Query,
			)
		}
	}
	return
}
