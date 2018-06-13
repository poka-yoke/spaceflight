package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/odin"
)

// snapshotListCmd represents the snapshot list command
var snapshotListCmd = &cobra.Command{
	Use:   "list identifier",
	Short: "Lists all snapshots",
	Long:  `Lists all snapshots, ordered from newer to older.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			args = append(args, "")
		}
		svc := odin.Init()
		snapshots, err := odin.ListSnapshots(
			args[0],
			svc,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Println(odin.PrintSnapshots(snapshots))
	},
}

func init() {
	SnapshotCmd.AddCommand(snapshotListCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Toggle help message")

}
