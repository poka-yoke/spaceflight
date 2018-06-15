package cmd

import (
	"log"
	"math"
	"os"

	"github.com/olorin/nagiosplugin"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/dnsbl"
)

var ipAddress, blacklist, pgaddress string
var warning, critical int

// blCmd represents the bl command
var blCmd = &cobra.Command{
	Use:   "bl",
	Short: "Check for blacklist presence",
	Long:  `Checks the supplied list of DNS-based blacklists for a specific IP presence. It returns the results in a Nagios compliat string. Thresholds for different Nagios states can be supplied as well.`,
	Run: func(cmd *cobra.Command, args []string) {
		blfile, err := os.Open(blacklist)
		if err != nil {
			log.Fatal("Could't open file ", blacklist, err)
		}
		defer blfile.Close()

		providers := dnsbl.GetProviders(ipAddress, blfile)
		if pgaddress == "" {
			positive, queried, length := dnsbl.NewChecker(providers).Query().Stats()

			check := nagiosplugin.NewCheck()
			defer check.Finish()
			must(check.AddPerfDatum("queried", "", float64(queried), 0.0, math.Inf(1)))
			must(check.AddPerfDatum("positive", "", float64(positive), 0.0, math.Inf(1)))
			check.AddResultf(
				nagiosplugin.OK,
				"%v present in %v(%v%%) out of %v BLs",
				ipAddress,
				positive,
				positive*100/length,
				length,
			)
			switch {
			case positive > length*warning/100:
				check.AddResultf(
					nagiosplugin.WARNING,
					"%v present in %v(%v%%) out of %v BLs",
					ipAddress,
					positive,
					positive*100/length,
					length,
				)
			case positive > length*critical/100:
				check.AddResultf(
					nagiosplugin.CRITICAL,
					"%v present in %v(%v%%) out of %v BLs",
					ipAddress,
					positive,
					positive*100/length,
					length,
				)
			}
		} else {
			err := push.New(pgaddress, "dnsbl").
				Collector(dnsbl.NewCollector(providers)).Push()

			if err != nil {
				log.Fatal("Could not push completion time to Pushgateway: ", err)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(blCmd)

	blCmd.Flags().StringVarP(&ipAddress, "ip", "", "127.0.0.1", "IP Address to look for in the BLs")
	blCmd.Flags().StringVarP(&blacklist, "file", "f", "", "Path to file containing black list addresses")
	blCmd.Flags().IntVarP(&warning, "warning", "w", 90, "IP Address to look for in the BLs")
	blCmd.Flags().IntVarP(&critical, "critical", "c", 95, "Path to file containing black list addresses")
	blCmd.Flags().StringVarP(
		&pgaddress,
		"push-gateway",
		"p",
		"",
		"Address of the Prometheus PushGateway to send results to.",
	)
}
