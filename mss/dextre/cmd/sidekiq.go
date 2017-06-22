package cmd

import (
	"log"
	"strings"

	"github.com/poka-yoke/spaceflight/mss/dextre/sidekiq"
	"github.com/spf13/cobra"
)

var baseURL string

// sidekiqCmd represents the sidekiq command
var sidekiqCmd = &cobra.Command{
	Use:   "sidekiq [opts]",
	Short: "Gather information from exposed sidekiq API",
	Long: `This command connects to an exposed API of sidekiq, scrapes the
endpoints and offers the information gathered in the form of a Nagios check.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !strings.HasPrefix(baseURL, "http://") ||
			!strings.HasPrefix(baseURL, "https://") {
			log.Fatal("Unknown scheme: ", baseURL)
		}
		info := sidekiq.ProcessGetResponse(baseURL)
		check := info.NagiosCheck()
		defer check.Finish()
	},
}

func init() {
	RootCmd.AddCommand(sidekiqCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sidekiqCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sidekiqCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	sidekiqCmd.Flags().StringVarP(&baseURL, "url", "", "", "Base URL for sidekiq check API")

}
