package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/short"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update [flags]",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		svc := short.NewService().SetAPIKey(apiKey)
		svc.AddURL = addURL

		err := svc.Update(domain, path, url)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")
	updateCmd.PersistentFlags().StringVarP(
		&domain,
		"domain",
		"",
		"",
		"Domain for the shortener URL",
	)
	updateCmd.PersistentFlags().StringVarP(
		&path,
		"path",
		"",
		"",
		"Path for the shortener URL",
	)
	updateCmd.PersistentFlags().StringVarP(
		&url,
		"url",
		"",
		"",
		"Destination for the shortened URL",
	)
	updateCmd.PersistentFlags().StringVarP(
		&apiKey,
		"apikey",
		"",
		"",
		"API key for the shortener service",
	)
	updateCmd.PersistentFlags().StringVarP(
		&addURL,
		"addURL",
		"",
		"",
		"Base URL for the Add endpoint",
	)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
