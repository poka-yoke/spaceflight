package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/odin"
)

// snapshotCreateCmd represents the snapshot create command
var snapshotCreateCmd = &cobra.Command{
	Use:   "create instanceID snapshotID",
	Short: "Creates a new snapshot",
	Long: `Create a new snapshot of the specified instance,
using the specified snapshot name.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			log.Fatal(CreateParamsReq)
		}
		svc := odin.Init()
		snapshot, err := odin.CreateSnapshot(
			args[0],
			args[1],
			svc,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Printf(
			"%s snapshot is being created",
			*snapshot.DBSnapshotIdentifier,
		)
	},
}

func init() {
	SnapshotCmd.AddCommand(snapshotCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Toggle help message")

}
