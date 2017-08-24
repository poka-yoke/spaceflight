package cmd

import (
	"github.com/spf13/cobra"
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
