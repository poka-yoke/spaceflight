package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/pkg/capcom"
)

var graph, search bool

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list [flags]",
	Short: "Show all Security Groups",
	Long: `
This option shows a information about the Security groups
present in your account. The information is shown as a list
but can also be presented in dot format for graphics processing.`,
	Run: func(cmd *cobra.Command, args []string) {
		svc := capcom.Init()
		if graph {
			fmt.Print(capcom.GraphSGRelations(svc))
		} else if search {
			list, err := capcom.FindSecurityGroupsWithRange(svc, args[0])
			if err != nil {
				log.Fatal(err.Error())
			}
			for _, l := range list {
				fmt.Println(l)
			}
		} else {
			fmt.Print(capcom.ListSecurityGroups(svc))
		}
	},
}

func init() {
	RootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	listCmd.Flags().BoolVarP(&graph, "graph", "g", false, "Output relations as a graph in DOT format")
	listCmd.Flags().BoolVarP(&search, "search", "s", false, "Search for this following CIDR in all SGs")

}
