package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/capcom"
)

var source, proto string
var port int64

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add [flags] <sgid1> [[sgid2] [[...]]]",
	Short: "Add rule to specified Security Group",
	Long: `
This option adds a rule allowing inbound access to AWS
machines pertaining to the selected security group (as
sgid) from the specified source (as either CIDR or sgid
string) to the specified port. E.g.:

    capcom add --source 1.2.3.4/32 sg-abc01234`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := capcom.Init()
		for _, sgid := range args {
			if !strings.HasPrefix(sgid, "sg-") {
				log.Fatalf("%s is invalid SG id\n", sgid)
			}
			perm, err := capcom.BuildIPPermission(source, proto, port)
			if err != nil {
				log.Fatal(err)
			}
			if !capcom.AuthorizeAccessToSecurityGroup(
				svc,
				perm,
				sgid,
			) {
				log.Fatalf("Failed to add rule to %s: %s %s %d\n",
					sgid,
					source,
					proto,
					port,
				)
			}
			log.Printf("Rule added successfully to %s: %s %s %d\n",
				sgid,
				source,
				proto,
				port,
			)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")
	addCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "CIDR or sgid to be used as source of the Security Group Inbound rule")
	addCmd.PersistentFlags().StringVarP(&proto, "proto", "", "tcp", "Which protocol will the rule affect to")
	addCmd.PersistentFlags().Int64VarP(&port, "port", "p", 22, "Port for the rule")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
