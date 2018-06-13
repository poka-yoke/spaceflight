package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/fido"
)

var name string

// rawCmd represents the raw command
var getCmd = &cobra.Command{
	Use:   "get [flags]",
	Short: "Obtain CustomJSON from OpsWorks stack",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		svc := fido.Init()
		sID, err := fido.GetStackID(svc, name)
		if err != nil {
			log.Fatal(err)
		}

		customJSON, err := fido.GetCustomJSON(svc, sID)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(customJSON)
	},
}

func init() {
	RootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rawCmd.PersistentFlags().String("foo", "", "A help for foo")
	getCmd.PersistentFlags().StringVarP(
		&name,
		"name",
		"",
		"",
		"Name of stack to retrieve",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rawCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
