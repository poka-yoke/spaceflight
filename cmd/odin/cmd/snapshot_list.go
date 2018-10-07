package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"

	"github.com/poka-yoke/spaceflight/pkg/odin"
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
		svc := rdsLogin("us-east-1")
		snapshots, err := odin.ListSnapshots(
			args[0],
			svc,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Println(printSnapshots(snapshots))
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

// printSnapshots return a string with the snapshot list output.
func printSnapshots(snapshots []*rds.DBSnapshot) string {
	lines := []string{}
	for _, snapshot := range snapshots {
		line := fmt.Sprintf(
			"%v %v %v %v\n",
			*snapshot.DBSnapshotIdentifier,
			*snapshot.DBInstanceIdentifier,
			(*snapshot.SnapshotCreateTime).Format(
				RFC8601,
			),
			*snapshot.Status,
		)
		lines = append(lines, line)
	}
	return strings.Join(lines, "")
}
