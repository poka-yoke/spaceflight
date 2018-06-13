package cmd

import (
	"log"

	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/got"
)

var dryrun, exclude, wait, filterByName, filterByType bool
var zoneName string
var ttl int64

// ttlCmd represents the ttl command
var ttlCmd = &cobra.Command{
	Use:   "ttl [flags] [filters ...]",
	Short: "Modify Time To Live of a set of records in a DNS zone",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		svc := got.Init()
		if len(zoneName) <= 0 {
			log.Fatal("No zone name specified")
		}
		zoneID := got.GetZoneID(zoneName, svc)

		list := got.GetResourceRecordSet(zoneID, svc)
		// Filter list in between
		if exclude && filterByType {
			list = got.FilterResourceRecords(
				list,
				args,
				func(
					elem *route53.ResourceRecordSet,
					filter string,
				) *route53.ResourceRecordSet {
					if *elem.Type != filter {
						return elem
					}
					return nil
				},
			)
		} else if exclude && filterByName {
			list = got.FilterResourceRecords(
				list,
				args,
				func(
					elem *route53.ResourceRecordSet,
					filter string,
				) *route53.ResourceRecordSet {
					if *elem.Name != filter {
						return elem
					}
					return nil
				},
			)
		} else if filterByType {
			list = got.FilterResourceRecords(
				list,
				args,
				func(
					elem *route53.ResourceRecordSet,
					filter string,
				) *route53.ResourceRecordSet {
					if *elem.Type == filter {
						return elem
					}
					return nil
				},
			)
		} else if filterByName {
			list = got.FilterResourceRecords(
				list,
				args,
				func(
					elem *route53.ResourceRecordSet,
					filter string,
				) *route53.ResourceRecordSet {
					if *elem.Name == filter {
						return elem
					}
					return nil
				},
			)
		}
		changeResponse, err := got.UpsertResourceRecordSetTTL(
			list,
			ttl,
			&zoneID,
			svc,
		)
		if err != nil {
			log.Panic(err.Error())
		}
		if wait && !dryrun {
			got.WaitForChangeToComplete(
				changeResponse.ChangeInfo,
				svc,
			)
		}
	},
}

func init() {
	RootCmd.AddCommand(ttlCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// ttlCmd.PersistentFlags().String("foo", "", "A help for foo")
	ttlCmd.PersistentFlags().BoolVarP(
		&dryrun,
		"dryrun",
		"",
		false,
		"Don't really do anything",
	)
	ttlCmd.PersistentFlags().BoolVarP(
		&exclude,
		"exclude",
		"",
		false,
		"Exclude records matching list",
	)
	ttlCmd.PersistentFlags().BoolVarP(
		&wait,
		"wait",
		"",
		false,
		"Don't return until operation is completed",
	)
	ttlCmd.PersistentFlags().StringVarP(
		&zoneName,
		"zone",
		"",
		"",
		"Name of the zone to work on.",
	)
	ttlCmd.PersistentFlags().Int64VarP(
		&ttl,
		"ttl",
		"",
		30,
		"New TTL value",
	)
	ttlCmd.PersistentFlags().BoolVarP(
		&filterByName,
		"name",
		"n",
		false,
		"Filters are to be used as names",
	)
	ttlCmd.PersistentFlags().BoolVarP(
		&filterByType,
		"type",
		"t",
		false,
		"Filters are to be used as types",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// ttlCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
