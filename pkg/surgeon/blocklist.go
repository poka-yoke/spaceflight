package surgeon

import "fmt"

// BlockList represents a list of Blocks.
type BlockList struct {
	Blocks []Block
}

// String provides a string formatting to BlockList
func (bl BlockList) String() (output string) {
	output = fmt.Sprintf(
		" %11s %13s %12s %14s %20s %20s \n",
		"BLOCKED PID",
		"BLOCKED USER",
		"BLOCKING PID",
		"BLOCKING USER",
		"BLOCKED QUERY",
		"BLOCKING QUERY",
	)
	for _, block := range bl.Blocks {
		output += fmt.Sprintf(
			" %11d %13s %12d %14s %20s %20s \n",
			block.BlockedPid,
			block.BlockedUser,
			block.BlockingPid,
			block.BlockingUser,
			block.BlockedQuery,
			block.BlockingQuery,
		)
	}
	return
}
