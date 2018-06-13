package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/odin"
)

// instanceRestoreCmd represents the instance restore command
var instanceRestoreCmd = &cobra.Command{
	Use:   "restore [flags] identifier",
	Short: "Restores from a Snapshot to a database",
	Long:  `Restores from a Snapshot to a database, in RDS.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal(NewInstanceIDReq)
		}
		svc := odin.Init()
		securityGroupsList := strings.Split(securityGroups, ",")
		params := odin.RestoreParams{
			InstanceType:         instanceType,
			SubnetGroupName:      subnetName,
			SecurityGroups:       securityGroupsList,
			OriginalInstanceName: from,
		}
		endpoint, err := odin.RestoreInstance(
			args[0],
			params,
			svc,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Println(endpoint)
	},
}

func init() {
	InstanceCmd.AddCommand(instanceRestoreCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	instanceRestoreCmd.PersistentFlags().StringVarP(
		&instanceType,
		"instance-type",
		"t",
		"db.m1.small",
		"Instance type to use when creating DB Instance",
	)
	instanceRestoreCmd.PersistentFlags().StringVarP(
		&from,
		"from",
		"f",
		"",
		"RDS Instance to look for snapshot",
	)
	instanceRestoreCmd.PersistentFlags().StringVarP(
		&subnetName,
		"subnet",
		"n",
		"",
		"DB Subnet Name to attach to (effectively VPC)",
	)
	instanceRestoreCmd.PersistentFlags().StringVarP(
		&securityGroups,
		"securityGroups",
		"g",
		"",
		"VPC SG IDs separated to attach to (effectively VPC)",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Toggle help message")

}
