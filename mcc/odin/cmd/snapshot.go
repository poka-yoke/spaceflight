package cmd

import (
	"github.com/spf13/cobra"
)

const (
	// CreateParamsReq is the message to show when arguments for
	// snapshot create are not satisfied.
	CreateParamsReq = `An existing available instance ID and
a non existing snapshot ID are expected`
)

// SnapshotCmd represents the snapshot super command
var SnapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Odin's Snapshot management",
	Long:  `Currently, it is mostly usable for AWS RDS.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	RootCmd.AddCommand(SnapshotCmd)
}
