package cmd

import (
	"bufio"
	"log"
	"os"

	"github.com/prometheus/client_golang/prometheus/push"
	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/sitemap"
)

var url, file string

// sitemapCmd represents the sitemap command
var sitemapCmd = &cobra.Command{
	Use:   "sitemap",
	Short: "Check for sitemap presence and timestamp",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		var urls []string

		if pgaddress == "" {
			log.Fatal("pgaddress is mandatory")
		}
		if url == "" && file == "" {
			log.Fatal("Either URL or file are mandatory")
		}

		if url != "" && file == "" {
			urls = []string{url}
		}
		if file != "" {
			fd, err := os.Open(file)
			if err != nil {
				log.Fatal("Couldn't open file ", file, err)
			}
			defer fd.Close()

			scanner := bufio.NewScanner(fd)
			for scanner.Scan() {
				urls = append(urls, scanner.Text())
			}
		}

		err := push.New(pgaddress, "sitemap").
			Collector(sitemap.NewCollector(urls)).Push()
		if err != nil {
			log.Fatal("An error occurred while processing: ", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(sitemapCmd)

	sitemapCmd.Flags().StringVarP(
		&url,
		"url",
		"",
		"",
		"URL to fetch sitemap file.",
	)
	sitemapCmd.Flags().StringVarP(
		&file,
		"file",
		"",
		"",
		"File contains the list of sitemap addresses to validate",
	)
	sitemapCmd.Flags().StringVarP(
		&pgaddress,
		"push-gateway",
		"p",
		"",
		"Address of the Prometheus PushGateway to send results to.",
	)
}
