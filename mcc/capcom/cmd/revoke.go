package cmd

import (
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/mcc/capcom/capcom"
)

// revokeCmd represents the revoke command
var revokeCmd = &cobra.Command{
	Use:   "revoke [flags] <sgid1> [[sgid2] [[...]]]",
	Short: "Revoke a rule in the specified Security Group",
	Long: `
This option removes a rule allowing inbound access to AWS
machines pertaining to the selected security group (as
sgid) from the specified source (as either CIDR or sgid
string) to the specified port. E.g.:

    capcom revoke --source 1.2.3.4/32 sg-abc01234`,
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
			_ = capcom.RevokeAccessToSecurityGroup(
				svc,
				perm,
				sgid,
			)
		}
	},
}

func init() {
	RootCmd.AddCommand(revokeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// revokeCmd.PersistentFlags().String("foo", "", "A help for foo")
	revokeCmd.PersistentFlags().StringVarP(&source, "source", "s", "", "CIDR or sgid to be used as source of the Security Group Inbound rule")
	revokeCmd.PersistentFlags().StringVarP(&proto, "proto", "", "tcp", "Which protocol will the rule affect to")
	revokeCmd.PersistentFlags().Int64VarP(&port, "port", "p", 22, "Port for the rule")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// revokeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
