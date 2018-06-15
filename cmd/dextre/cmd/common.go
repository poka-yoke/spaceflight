package cmd

import (
	"fmt"
	"log"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func checkGoldenFile() error {
	if goldenFile == "" {
		return fmt.Errorf("Golden file is mandatory")
	}
	return nil
}

func checkURL(args []string) string {
	if len(args) < 1 {
		log.Fatalf("No URL specified")
	}
	return args[0]
}

func checkPGAddress() error {
	if pgaddress == "" {
		return fmt.Errorf("Push Gateway address is mandatory")
	}
	return nil
}
