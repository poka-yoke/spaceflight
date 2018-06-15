package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/internal/digest"
	"github.com/Devex/spaceflight/internal/http"
)

var robotsState = prometheus.NewGauge(prometheus.GaugeOpts{
	Name: "robots_txt_state",
	Help: `The state of the robots.txt comparison to the golden
file. 0 means equal, 1, different`,
})

// robotsCheckCmd represents the robots check command
var robotsCheckCmd = &cobra.Command{
	Use:   "check",
	Short: "Checks differences between a URL and a golden file",
	Long: `Checks differences between a URL and a golden file.
The golden file can be generated using the gen_gc command`,
	Run: func(cmd *cobra.Command, args []string) {
		must(checkPGAddress())
		must(checkGoldenFile())
		url := checkURL(args)

		resp := http.Get(url)
		defer resp.Body.Close()

		body, err := digest.ContentBase64(resp.Body)
		if err != nil {
			log.Fatalf("Failed to hash downloaded content: %s", err)
		}

		in, err := os.Open(goldenFile)
		if err != nil {
			log.Fatalf(
				"Failed to read %s: %s",
				goldenFile,
				err,
			)
		}
		defer in.Close()

		golden, err := digest.ContentBase64(in)
		if err != nil {
			log.Fatalf("Failed to hash golden file: %s", err)
		}

		pusher := push.New(pgaddress, "robots").
			Collector(robotsState).
			Grouping("web", strings.Split(url, "/")[2])

		if body != golden {
			robotsState.Set(1)
			log.Printf(
				"Hashes are different! Expected %s but got %s",
				golden,
				body,
			)
		} else {
			robotsState.Set(0)
			log.Printf("Hashes are equal")
		}

		if err := pusher.Push(); err != nil {
			log.Fatalf("Could not push to Push Gateway: %v", err)
		}
	},
}

func init() {
	RobotsCmd.AddCommand(robotsCheckCmd)

	robotsCheckCmd.Flags().StringVarP(
		&goldenFile,
		"golden-file",
		"g",
		"",
		"File name for the Golden file.",
	)

	robotsCheckCmd.Flags().StringVarP(
		&pgaddress,
		"push-gateway",
		"p",
		"",
		"Address of the Prometheus PushGateway to send results to.",
	)
}
