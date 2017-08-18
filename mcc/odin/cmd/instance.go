package cmd

import (
	"github.com/spf13/cobra"
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
