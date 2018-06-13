package cmd

import (
	"github.com/spf13/cobra"
)

var goldenFile string

// RobotsCmd represents the robots command
var RobotsCmd = &cobra.Command{
	Use:   "robots",
	Short: "Monitor published robots.txt",
	Long: `Helps with robots.txt monitoring, detecting changes from a
golden file, that needs to be generated.`,
}

func init() {
	RootCmd.AddCommand(RobotsCmd)
}
