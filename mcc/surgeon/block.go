package surgeon

// Block represents a row from '[custom block]'.
type Block struct {
	BlockedPid    int    // bdpid
	BlockedUser   Name   // bduser
	BlockingPid   int    // bgpid
	BlockingUser  Name   // bguser
	BlockedQuery  string // bdquery
	BlockingQuery string // bgquery
}
