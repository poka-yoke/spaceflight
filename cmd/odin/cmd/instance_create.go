package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/poka-yoke/spaceflight/pkg/odin"
)

var instanceType, password, user, subnetName, securityGroups string
var size int64

// instanceCreateCmd represents the instance create command
var instanceCreateCmd = &cobra.Command{
	Use:   "create [flags] identifier",
	Short: "Creates a database",
	Long:  `Creates a database, according to defined attributes, in RDS.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal(NewInstanceIDReq)
		}
		svc := rdsLogin("us-east-1")
		params := odin.Instance{
			Identifier:      args[0],
			Type:            instanceType,
			User:            user,
			Password:        password,
			SubnetGroupName: subnetName,
			SecurityGroups:  strings.Split(securityGroups, ","),
			Size:            size,
		}
		rdsParams, err := params.CreateDBInput(
			svc,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		res, err := svc.CreateDBInstance(rdsParams)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		err = waitForInstance(
			res.DBInstance,
			svc,
			"available",
			5*time.Second,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		if err := modifyInstance(params, svc, false); err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Println(*res.DBInstance.Endpoint.Address)
	},
}

func init() {
	InstanceCmd.AddCommand(instanceCreateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	instanceCreateCmd.PersistentFlags().StringVarP(
		&instanceType,
		"instance-type",
		"t",
		"db.m1.small",
		"Instance type to use when creating DB Instance",
	)
	instanceCreateCmd.PersistentFlags().StringVarP(
		&user,
		"user",
		"u",
		"",
		"User to use when creating DB Instance",
	)
	instanceCreateCmd.PersistentFlags().StringVarP(
		&password,
		"password",
		"p",
		"",
		"Password to use when creating DB Instance",
	)
	instanceCreateCmd.PersistentFlags().Int64VarP(
		&size,
		"size",
		"s",
		5,
		"Size to use when creating DB Instance",
	)
	instanceCreateCmd.PersistentFlags().StringVarP(
		&subnetName,
		"subnet",
		"n",
		"",
		"DB Subnet Name to attach to (effectively VPC)",
	)
	instanceCreateCmd.PersistentFlags().StringVarP(
		&securityGroups,
		"securityGroups",
		"g",
		"",
		"VPC SG IDs separated with , to attach to (effectively VPC)",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Toggle help message")

}
