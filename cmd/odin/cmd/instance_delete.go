package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/odin"
)

var finalSnapshotID string

// instanceDeleteCmd represents the instance delete command
var instanceDeleteCmd = &cobra.Command{
	Use:   "delete [flags] identifier",
	Short: "Deletes a database",
	Long:  `Deletes a database in RDS.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal(InstanceIDReq)
		}
		svc := odin.Init()
		err := odin.DeleteInstance(
			args[0],
			finalSnapshotID,
			svc,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Printf("%s instance was deleted\n", args[0])
	},
}

func init() {
	InstanceCmd.AddCommand(instanceDeleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Toggle help message")
	instanceDeleteCmd.PersistentFlags().StringVarP(
		&finalSnapshotID,
		"final-snapshot-id",
		"s",
		"",
		"Final snapshot ID, if desired. If missing, no snapshot",
	)

}
