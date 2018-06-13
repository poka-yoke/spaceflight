package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/odin"
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
		svc := odin.Init()
		result, err := odin.ScaleInstance(
			args[0],
			instanceType,
			delay,
			svc,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Println(result)
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
