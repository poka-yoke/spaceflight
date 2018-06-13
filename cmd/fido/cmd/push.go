package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/fido"
)

var file string

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Upload CustomJSON to stack",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		customJSON, err := fido.ReadFromPipe()
		if err != nil {
			log.Fatal(err.Error())
		}
		svc := fido.Init()
		sID, err := fido.GetStackID(svc, name)
		if err != nil {
			log.Fatal(err)
		}

		err = fido.PushCustomJSON(svc, sID, customJSON)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Successfully uploaded CustomJSON to %s\n", name)
	},
}

func init() {
	RootCmd.AddCommand(pushCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pushCmd.PersistentFlags().String("foo", "", "A help for foo")
	pushCmd.PersistentFlags().StringVarP(
		&name,
		"name",
		"",
		"",
		"Name of stack to retrieve",
	)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// pushCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
