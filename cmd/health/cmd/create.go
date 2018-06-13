package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/cronitor"
	"github.com/Devex/spaceflight/pkg/health"
)

// Check is an interface to be implemented by all backend providers
type Check interface {
	SetAPIKey(string)
	Create(string, map[string]interface{}) (*http.Request, error)
}

var apikey, endpoint, schedule, name, tags, sep, email string
var check Check

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		message := make(map[string]interface{})
		out := []string{}
		// Select backend to prepare request
		switch {
		case strings.Contains(endpoint, "healthchecks.io"):
			check = health.NewCheck()
			message = health.SetMessage(schedule, name, tags)
		case strings.Contains(endpoint, "cronitor.io"):
			check = cronitor.NewCheck()
			message = cronitor.SetMessage(schedule, name, tags, email)
		default:
			log.Fatal("Unrecognized provider")
		}
		// backend independent processing
		check.SetAPIKey(apikey)
		req, err := check.Create(endpoint, message)
		if err != nil {
			log.Fatal(err)
		}
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		m := make(map[string]interface{})
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(body, &m)
		if err != nil {
			log.Fatal(err)
		}
		// Select backend to process answer
		switch {
		case strings.Contains(endpoint, "healthchecks.io"):
			v, ok := m["update_url"].(string)
			if !ok {
				log.Fatal("Can't access field update_url")
			}
			slug := health.GetSlugFromURL(v)
			out = []string{schedule, name, tags, slug}
		case strings.Contains(endpoint, "cronitor.io"):
			v, ok := m["code"].(string)
			if !ok {
				log.Fatal("Can't retrieve id")
			}
			out = []string{schedule, name, tags, v}
		}
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
	createCmd.PersistentFlags().StringVarP(
		&email,
		"email",
		"",
		"",
		"Contact email adresses (cronitor only)",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
