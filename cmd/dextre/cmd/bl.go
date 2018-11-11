package cmd

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cobra"

	"github.com/poka-yoke/spaceflight/pkg/dnsbl"
)

var ipAddress, blacklist, pgaddress string
var warning, critical int

// blCmd represents the bl command
var blCmd = &cobra.Command{
	Use:   "bl",
	Short: "Check for blacklist presence",
	Long:  `Checks the supplied list of DNS-based blacklists for a specific IP presence and sends the results as metrics to a Prometheus' Push Gateway.`,
	Run: func(cmd *cobra.Command, args []string) {
		if blacklist == "" {
			log.Fatal("No file with blacklist addresses specified")
		}
		if pgaddress == "" {
			log.Fatal("Push Gateway address is required")
		}

		blfile, err := os.Open(blacklist)
		if err != nil {
			log.Fatal("Could't open file ", blacklist, err)
		}
		defer blfile.Close()

		providers := dnsbl.GetProviders(ipAddress, blfile)

		err = push.New(pgaddress, "dnsbl").
			Collector(dnsbl.NewCollector(providers)).Push()

		if err != nil {
			log.Fatal("Could not push completion time to Pushgateway: ", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(blCmd)

	blCmd.Flags().StringVarP(&ipAddress, "ip", "", "127.0.0.1", "IP Address to look for in the BLs")
	blCmd.Flags().StringVarP(&blacklist, "file", "f", "", "Path to file containing black list addresses")
	blCmd.Flags().StringVarP(
		&pgaddress,
		"push-gateway",
		"p",
		"",
		"Address of the Prometheus PushGateway to send results to.",
	)
}
