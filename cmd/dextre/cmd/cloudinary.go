package cmd

import (
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cobra"

	"github.com/poka-yoke/spaceflight/pkg/cloudinary"
)

// cloudinaryCmd represents the cloudinary command
var cloudinaryCmd = &cobra.Command{
	Use:   "cloudinary",
	Short: "Get cloudinary metrics",
	Long:  `Gets the cloudinary usage metrics from the Admin API and sends the results as metrics to a Prometheus' Push Gateway.`,
	Run: func(cmd *cobra.Command, args []string) {
		must(cloudinary.NewCredentials(
			os.Getenv("CLOUDINARY_CLOUD_NAME"),
			os.Getenv("CLOUDINARY_KEY"),
			os.Getenv("CLOUDINARY_SECRET"),
		))

		must(checkPGAddress())
		err := push.New(pgaddress, "cloudinary").
			Collector(cloudinary.NewCollector()).Push()

		if err != nil {
			log.Fatal("Could not push completion time to Pushgateway: ", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(cloudinaryCmd)

	cloudinaryCmd.Flags().StringVarP(
		&pgaddress,
		"push-gateway",
		"p",
		"",
		"Address of the Prometheus PushGateway to send results to.",
	)
}
