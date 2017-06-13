package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/mcc/odin/odin"
)

var instanceType, password, user, from string
var size int64
var restore bool

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [flags] identifier",
	Short: "Creates a database",
	Long:  `Creates a database, into RDS.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal("You must specify the instance identifier for the new instance")
		}
		if user == "" && !restore {
			log.Fatal("User should be provided and not be blank")
		}
		if password == "" && !restore {
			log.Fatal("Password should be provided and not be blank")
		}
		svc := odin.Init()
		params := odin.CreateDBParams{
			DBInstanceType:       instanceType,
			DBUser:               user,
			DBPassword:           password,
			Size:                 size,
			OriginalInstanceName: from,
			Restore:              restore,
		}
		endpoint, err := odin.CreateDBInstance(
			args[0],
			params,
			svc,
		)
		if err != nil {
			log.Fatal("Error: %s", err)
		}
		fmt.Println(endpoint)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	createCmd.PersistentFlags().StringVarP(
		&instanceType,
		"instance-type",
		"t",
		"db.m1.small",
		"Instance type to use when creating DB Instance",
	)
	createCmd.PersistentFlags().StringVarP(
		&user,
		"user",
		"u",
		"",
		"User to use when creating DB Instance",
	)
	createCmd.PersistentFlags().StringVarP(
		&password,
		"password",
		"p",
		"",
		"Password to use when creating DB Instance",
	)
	createCmd.PersistentFlags().Int64VarP(
		&size,
		"size",
		"s",
		5,
		"Size to use when creating DB Instance",
	)
	createCmd.PersistentFlags().StringVarP(
		&from,
		"from",
		"f",
		"",
		"RDS Instance to look for snapshot",
	)
	createCmd.PersistentFlags().BoolVarP(
		&restore,
		"restore",
		"r",
		false,
		"If true, restores; else, it just clones the parameters",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
