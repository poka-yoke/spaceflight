package cmd

import (
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/internal/http"
)

// robotsGenGCmd represents the robots gen_gc command
var robotsGenGCmd = &cobra.Command{
	Use:   "generate-golden",
	Short: "Get golden file for robots check",
	Long: `Gets the specified URL contents and stores it as a golden file
for subsequent robots checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		must(checkPGAddress())
		must(checkGoldenFile())
		url := checkURL(args)

		if _, err := os.Stat(goldenFile); err == nil {
			log.Fatalf("Golden file %s already exists", goldenFile)
		}

		resp := http.Get(url)
		defer resp.Body.Close()

		out, err := os.Create(goldenFile)
		if err != nil {
			log.Fatalf(
				"Failed to create golden file %s: %s",
				goldenFile,
				err,
			)
		}
		defer out.Close()

		_, err = io.Copy(out, resp.Body)
		if err != nil {
			log.Fatalf(
				"Failed to copy data from %s to %s: %s",
				url,
				goldenFile,
				err,
			)
		}
	},
}

func init() {
	RobotsCmd.AddCommand(robotsGenGCmd)

	robotsGenGCmd.Flags().StringVarP(
		&goldenFile,
		"golden-file",
		"g",
		"",
		"File name for the Golden File.",
	)
}
