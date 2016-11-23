package surgeon

// Process represents a row from ProcessList.
type Process struct {
	Pid      int      // pid
	Usename  Name     // usename
	Datname  Name     // datname
	Client   Inet     // client_addr
	Duration Duration // duration
	Query    string   // query
	Waiting  string   // waiting
}

// IsActive returns true if the Process was active when extracted.
func (ps Process) IsActive() bool {
	return ps.Query != "Inactive"
}
