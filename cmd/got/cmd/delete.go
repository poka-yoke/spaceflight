package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/got"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [flags] [record] [record] ...",
	Short: "Remove DNS records",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		svc := got.Init()
		if len(zoneName) <= 0 {
			log.Fatal("No zone name specified")
		}
		if len(typ) <= 0 {
			log.Fatal("No record type specified")
		}
		if len(args) <= 0 {
			log.Fatal("No record names specified")
		}
		zoneid := got.GetZoneID(zoneName, svc)
		list := got.GetResourceRecordSet(zoneid, svc)
		changes := got.DeleteChangeList(args, typ, list)
		if !dryrun {
			res, err := got.ApplyChanges(changes, &zoneid, svc)
			if err != nil {
				log.Fatal(err.Error())
			}
			log.Println(res)
		}
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")
	deleteCmd.PersistentFlags().BoolVarP(
		&dryrun,
		"dryrun",
		"",
		false,
		"Don't really do anything",
	)
	deleteCmd.PersistentFlags().BoolVarP(
		&wait,
		"wait",
		"",
		false,
		"Don't return until operation is completed",
	)
	deleteCmd.PersistentFlags().StringVarP(
		&zoneName,
		"zone",
		"",
		"",
		"Name of the zone to work on.",
	)
	deleteCmd.PersistentFlags().StringVarP(
		&typ,
		"type",
		"",
		"",
		"Type of the record to upsert.",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
