package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/olorin/nagiosplugin"
	"github.com/poka-yoke/spaceflight/mss/dextre"
)

func fromFile(path *string) []string {
	var blacklists []string

	blfile, err := os.Open(*path)
	if err != nil {
		log.Fatal("Could't open file ", *path)
		log.Fatal(err)
	}
	defer blfile.Close()

	scanner := bufio.NewScanner(blfile)
	for scanner.Scan() {
		blacklists = append(blacklists, scanner.Text())
	}
	return blacklists
}

func main() {
	ipAddress := flag.String(
		"ip",
		"127.0.0.1",
		"IP Address to look for in the BLs",
	)
	warning := flag.Int("w", 90, "Warning threshold")
	critical := flag.Int("c", 95, "Critical threshold")
	blacklist := flag.String(
		"f",
		"",
		"Path to file containing black list addresses",
	)

	flag.Parse()

	blacklists := fromFile(blacklist)

	check := nagiosplugin.NewCheck()
	defer check.Finish()
	responses := make(chan int, len(blacklists))

	queried := 0
	positive := 0
	for _, list := range blacklists {
		go dextre.DNSBLQuery(*ipAddress, list, responses)
	}
	for i := 0; i < len(blacklists); i++ {
		response := <-responses
		if response > 0 {
			positive += response
		}
		queried++
	}
	warningAmount := len(blacklists) * (*warning) / 100
	criticalAmount := len(blacklists) * (*critical) / 100
	checkLevel := nagiosplugin.OK
	if positive > warningAmount {
		checkLevel = nagiosplugin.WARNING
		if positive > criticalAmount {
			checkLevel = nagiosplugin.CRITICAL
		}
	}
	check.AddResult(
		checkLevel,
		fmt.Sprintf(
			"%v present in %v(%v%%) out of %v BLs | %v",
			*ipAddress,
			positive,
			positive*100/len(blacklists),
			len(blacklists),
			fmt.Sprintf(
				"queried=%v positive=%v",
				queried,
				positive,
			),
		),
	)
}
