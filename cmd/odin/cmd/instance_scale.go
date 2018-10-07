package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/poka-yoke/spaceflight/pkg/odin"
)

var delay bool

// instanceScaleCmd invokes the ScaleInstance function with defined
// parameter's from user's input.
var instanceScaleCmd = &cobra.Command{
	Use:   "scale [flags] identifier",
	Short: "Scales a database",
	Long:  `Scales a database, according to defined attributes, in RDS.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal(InstanceIDReq)
		}
		svc := rdsLogin("us-east-1")
		params := odin.Instance{
			Identifier: args[0],
			Type:       instanceType,
		}
		rdsParams, err := params.ModifyDBInput(!delay, svc)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		out, err := svc.ModifyDBInstance(rdsParams)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		err = waitForInstance(
			out.DBInstance,
			svc,
			"available",
			5*time.Second,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Printf(
			"Instance %s is %s\n",
			*out.DBInstance.DBInstanceIdentifier,
			*out.DBInstance.DBInstanceClass,
		)
	},
}

func init() {
	InstanceCmd.AddCommand(instanceScaleCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	instanceScaleCmd.PersistentFlags().StringVarP(
		&instanceType,
		"instance-type",
		"t",
		"db.m1.small",
		"Instance type to use when creating DB Instance",
	)
	instanceScaleCmd.PersistentFlags().BoolVarP(
		&delay,
		"delay",
		"d",
		false,
		"Scales on next reboot or during maintenance",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Toggle help message")

}
