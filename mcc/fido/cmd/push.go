package cmd

import (
	"bufio"
	"io"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/Devex/spaceflight/mcc/fido/fido"
)

var file string

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Upload CustomJSON to stack",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		nBytes, nChunks := int64(0), int64(0)
		r := bufio.NewReader(os.Stdin)
		buf := make([]byte, 0, 4*1024)
		customJSON := ""
		for {
			n, err := r.Read(buf[:cap(buf)])
			buf = buf[:n]
			if n == 0 {
				if err == nil {
					continue
				}
				if err == io.EOF {
					break
				}
				log.Fatal(err)
			}
			nChunks++
			nBytes += int64(len(buf))
			customJSON += string(buf)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			log.Println("Bytes:", nBytes, "Chunks:", nChunks)
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
