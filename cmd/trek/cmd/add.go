package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/trek"
)

var original, final string

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds a redirect",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		redirects, err := trek.ReadFromPipe()
		if err != nil {
			log.Fatal(err.Error())
		}
		output, err := trek.Add(redirects, original, final)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Printf("%s", output)
	},
}

func init() {
	RootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	addCmd.PersistentFlags().StringVarP(
		&original,
		"original",
		"",
		"",
		"Path or URI to redirect",
	)
	addCmd.PersistentFlags().StringVarP(
		&final,
		"final",
		"",
		"",
		"Path or URI to redirect to",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
