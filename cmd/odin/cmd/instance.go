package cmd

import (
	"github.com/spf13/cobra"
)

const (
	// NewInstanceIDReq is the message to show when a new instance is to be
	// created and no Id was specified.
	NewInstanceIDReq = "You must specify a identifier for the new instance"
	// InstanceIDReq is the message to show when an instance is to be
	// operated and no Id was specified.
	InstanceIDReq = "You must specify an instance identifier"
)

// InstanceCmd represents the instance super command
var InstanceCmd = &cobra.Command{
	Use:   "instance",
	Short: "It is a database instance management tool",
	Long:  `Currently, it is mostly usable for AWS RDS.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	RootCmd.AddCommand(InstanceCmd)
}
