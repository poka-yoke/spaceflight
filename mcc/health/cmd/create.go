package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/mcc/health/health"
)

var apikey, endpoint, schedule, name, tags, sep string

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		message := make(map[string]interface{})
		out := []string{}
		check := health.NewCheck()
		check.APIKey = apikey
		if schedule != "" {
			message["schedule"] = schedule
			out = append(out, schedule)
		}
		if name != "" {
			message["name"] = name
			out = append(out, name)
		}
		if tags != "" {
			message["tags"] = tags
			out = append(out, tags)
		}
		res, err := check.Create(endpoint, message)
		if err != nil {
			log.Fatal(err)
		}
		m, err := health.ParseResponse(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		v, ok := m["update_url"].(string)
		if !ok {
			log.Fatal("Can't access field update_url")
		}
		slug := health.GetSlugFromURL(v)
		out = append(out, slug)
		fmt.Println(strings.Join(out, sep))
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")
	createCmd.PersistentFlags().StringVarP(
		&apikey,
		"apikey",
		"",
		"",
		"Healthchecks.io API key",
	)
	createCmd.PersistentFlags().StringVarP(
		&endpoint,
		"api",
		"",
		"https://healthchecks.io/api/v1/checks/",
		"Healthchecks.io API address",
	)
	createCmd.PersistentFlags().StringVarP(
		&schedule,
		"schedule",
		"",
		"",
		"Cron-like schedule",
	)
	createCmd.PersistentFlags().StringVarP(
		&name,
		"name",
		"",
		"",
		"Name identifier for the entry",
	)
	createCmd.PersistentFlags().StringVarP(
		&tags,
		"tags",
		"",
		"",
		"Space separated list of tags",
	)
	createCmd.PersistentFlags().StringVarP(
		&sep,
		"separator",
		"",
		" ",
		"Field separator string for output",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
