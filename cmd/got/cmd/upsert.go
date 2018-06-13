package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/got"
)

var name, typ string

// upsertCmd represents the upsert command
var upsertCmd = &cobra.Command{
	Use:   "upsert [flags] <destination>",
	Short: "Upsert a DNS record",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		svc := got.Init()
		if len(zoneName) <= 0 {
			log.Fatal("No zone name specified")
		}
		if len(name) <= 0 {
			log.Fatal("No record name specified")
		}
		if len(typ) <= 0 {
			log.Fatal("No record type specified")
		}
		if len(args) <= 0 {
			log.Fatal("No destination specified")
		}
		zoneid := got.GetZoneID(zoneName, svc)
		list := got.NewResourceRecordList(args)
		changes := got.UpsertChangeList(list, ttl, name, typ)
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
	RootCmd.AddCommand(upsertCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// upsertCmd.PersistentFlags().String("foo", "", "A help for foo")
	upsertCmd.PersistentFlags().BoolVarP(
		&dryrun,
		"dryrun",
		"",
		false,
		"Don't really do anything",
	)
	upsertCmd.PersistentFlags().BoolVarP(
		&exclude,
		"exclude",
		"",
		false,
		"Exclude records matching list",
	)
	upsertCmd.PersistentFlags().BoolVarP(
		&wait,
		"wait",
		"",
		false,
		"Don't return until operation is completed",
	)
	upsertCmd.PersistentFlags().StringVarP(
		&zoneName,
		"zone",
		"",
		"",
		"Name of the zone to work on.",
	)
	upsertCmd.PersistentFlags().Int64VarP(
		&ttl,
		"ttl",
		"",
		30,
		"New UPSERT value",
	)
	upsertCmd.PersistentFlags().StringVarP(
		&name,
		"name",
		"",
		"",
		"Name of the record to upsert.",
	)
	upsertCmd.PersistentFlags().StringVarP(
		&typ,
		"type",
		"",
		"",
		"Type of the record to upsert.",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// upsertCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
