package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// jsonCmd represents the json command
var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "Modify contents of supplied JSON file",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("json called")
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
