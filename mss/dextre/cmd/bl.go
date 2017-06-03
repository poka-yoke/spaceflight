package cmd

import (
	"log"
	"math"

	"github.com/olorin/nagiosplugin"
	"github.com/poka-yoke/spaceflight/mss/dextre/dnsbl"
	"github.com/spf13/cobra"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var ipAddress, blacklist string
var warning, critical int

// blCmd represents the bl command
var blCmd = &cobra.Command{
	Use:   "bl",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		list := dnsbl.FromFile(blacklist)
		dnsbl.Queries(ipAddress, list)

		positive := dnsbl.Stats.Positive
		queried := dnsbl.Stats.Queried
		length := dnsbl.Stats.Length

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
			check.AddResult(nagiosplugin.WARNING, "")
		case positive > length*critical/100:
			check.AddResultf(nagiosplugin.CRITICAL, "")
		}
	},
}

func init() {
	RootCmd.AddCommand(blCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// blCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// blCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	blCmd.Flags().StringVarP(&ipAddress, "ip", "", "127.0.0.1", "IP Address to look for in the BLs")
	blCmd.Flags().StringVarP(&blacklist, "file", "f", "", "Path to file containing black list addresses")
	blCmd.Flags().IntVarP(&warning, "warning", "w", 90, "IP Address to look for in the BLs")
	blCmd.Flags().IntVarP(&critical, "critical", "c", 95, "Path to file containing black list addresses")

}
