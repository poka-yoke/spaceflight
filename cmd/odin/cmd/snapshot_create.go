package cmd

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/spf13/cobra"
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
		svc := rdsLogin("us-east-1")
		output, err := svc.CreateDBSnapshot(
			&rds.CreateDBSnapshotInput{
				DBInstanceIdentifier: aws.String(args[0]),
				DBSnapshotIdentifier: aws.String(args[1]),
			},
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Printf(
			"%s snapshot is being created",
			*output.DBSnapshot.DBSnapshotIdentifier,
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
