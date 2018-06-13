package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/short"
)

var apiKey, addURL string
var domain, path, url string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "add [flags]",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		svc := short.NewService().SetAPIKey(apiKey)
		svc.AddURL = addURL

		err := svc.Add(domain, path, url)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")
	addCmd.PersistentFlags().StringVarP(
		&domain,
		"domain",
		"",
		"",
		"Domain for the shortener URL",
	)
	addCmd.PersistentFlags().StringVarP(
		&path,
		"path",
		"",
		"",
		"Path for the shortener URL",
	)
	addCmd.PersistentFlags().StringVarP(
		&url,
		"url",
		"",
		"",
		"Destination for the shortened URL",
	)
	addCmd.PersistentFlags().StringVarP(
		&apiKey,
		"apikey",
		"",
		"",
		"API key for the shortener service",
	)
	addCmd.PersistentFlags().StringVarP(
		&addURL,
		"addURL",
		"",
		"",
		"Base URL for the Add endpoint",
	)
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
