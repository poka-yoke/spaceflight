package cmd

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/rds/rdsiface"
	"github.com/spf13/cobra"

	"github.com/poka-yoke/spaceflight/pkg/odin"
)

var from string

// instanceCloneCmd represents the instance clone command
var instanceCloneCmd = &cobra.Command{
	Use:   "clone [flags] identifier",
	Short: "Clones a database using params from a Snapshot",
	Long:  `Clones a database using params from a Snapshot, in RDS.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			log.Fatal(NewInstanceIDReq)
		}
		svc := rdsLogin("us-east-1")
		if from == "" {
			log.Fatal("Original instance name not provided")
		}
		params := odin.Instance{
			Identifier:           args[0],
			Type:                 instanceType,
			User:                 user,
			Password:             password,
			SubnetGroupName:      subnetName,
			SecurityGroups:       strings.Split(securityGroups, ","),
			Size:                 size,
			OriginalInstanceName: from,
		}
		endpoint, err := cloneInstance(
			params,
			svc,
			5*time.Second,
		)
		if err != nil {
			log.Fatalf("Error: %s", err)
		}
		fmt.Println(endpoint)
	},
}

func init() {
	InstanceCmd.AddCommand(instanceCloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	instanceCloneCmd.PersistentFlags().StringVarP(
		&instanceType,
		"instance-type",
		"t",
		"db.m1.small",
		"Instance type to use when creating DB Instance",
	)
	instanceCloneCmd.PersistentFlags().StringVarP(
		&user,
		"user",
		"u",
		"",
		"User to use when creating DB Instance",
	)
	instanceCloneCmd.PersistentFlags().StringVarP(
		&password,
		"password",
		"p",
		"",
		"Password to use when creating DB Instance",
	)
	instanceCloneCmd.PersistentFlags().Int64VarP(
		&size,
		"size",
		"s",
		5,
		"Size to use when creating DB Instance",
	)
	instanceCloneCmd.PersistentFlags().StringVarP(
		&from,
		"from",
		"f",
		"",
		"RDS Instance to look for snapshot",
	)
	instanceCloneCmd.PersistentFlags().StringVarP(
		&subnetName,
		"subnet",
		"n",
		"",
		"DB Subnet Name to attach to (effectively VPC)",
	)
	instanceCloneCmd.PersistentFlags().StringVarP(
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

// cloneInstance creates a new RDS database instance, copying parameters
// from a snapshot. If a vpcid is specified the security group will be
// in that VPC.
func cloneInstance(
	params odin.Instance,
	svc rdsiface.RDSAPI,
	duration time.Duration,
) (
	result string,
	err error,
) {
	rdsParams, err := params.CloneDBInput(
		svc,
	)
	if err != nil {
		return "", err
	}
	res, err := svc.CreateDBInstance(rdsParams)
	if err != nil {
		return "", err
	}
	err = waitForInstance(res.DBInstance, svc, "available", duration)
	if err != nil {
		return
	}
	result = *res.DBInstance.Endpoint.Address
	err = odin.ModifyInstance(params, svc)
	return
}
