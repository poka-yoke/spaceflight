package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/capcom"
)

var name, vpcid string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [flags] <description>",
	Short: "Create a new Security Group",
	Long: `Example:
    capcom create --name test This is a test SG
    capcom create --name test --vpcid vpc-12345678 This is a test SG in a vpc`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := capcom.Init()
		sgid := capcom.CreateSG(name, strings.Join(args, " "), vpcid, svc)
		fmt.Println(sgid)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	createCmd.PersistentFlags().StringVarP(
		&vpcid,
		"vpcid",
		"",
		"",
		"VPC ID where the SG should be created",
	)
	createCmd.PersistentFlags().StringVarP(
		&name,
		"name",
		"",
		"",
		"Name for the Security Group",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
