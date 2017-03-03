package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/poka-yoke/spaceflight/mcc/fido/fido"
)

// rawCmd represents the raw command
var rawCmd = &cobra.Command{
	Use:   "raw",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var root, sn fido.Section
		en := fido.NewStringEntry("key1", "val1")
		en2 := fido.NewStringEntry("key2", "val2")
		root = root.AddEntry(en)
		root = root.AddEntry(en2)
		sn = sn.AddEntry(en)
		root = root.AddSection(sn)

		fmt.Println(root)
	},
}

func init() {
	RootCmd.AddCommand(rawCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// rawCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// rawCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

}
