package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/fido"
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Modify contents of supplied JSON file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		customJSON, err := fido.ReadFromPipe()
		if err != nil {
			log.Fatal(err.Error())
		}

		// Convert JSON into struct
		data := new(interface{})
		err = json.Unmarshal([]byte(customJSON), data)
		if err != nil {
			log.Panic(err)
		}

		// Convert struct into JSON
		out, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			log.Panic(err)
		}
		fmt.Print(string(out))

	},
}

func init() {
	RootCmd.AddCommand(jsonCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jsonCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jsonCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
