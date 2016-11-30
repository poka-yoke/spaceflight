package main

import (
	"flag"
	"fmt"

	"github.com/olorin/nagiosplugin"
	"github.com/poka-yoke/spaceflight/mss/dextre"
)

func main() {
	ipAddress := flag.String(
		"ip",
		"127.0.0.1",
		"IP Address to look for in the BLs",
	)
	warning := flag.Int("w", 90, "Warning threshold")
	critical := flag.Int("c", 95, "Critical threshold")

	flag.Parse()

	check := nagiosplugin.NewCheck()
	defer check.Finish()
	responses := make(chan int, len(dextre.Blacklists))

	queried := 0
	positive := 0
	for i := range dextre.Blacklists {
		go dextre.DNSBLQuery(*ipAddress, dextre.Blacklists[i], responses)
	}
	for i := 0; i < len(dextre.Blacklists); i++ {
		response := <-responses
		if response > 0 {
			positive += response
		}
		queried++
	}
	warningAmount := len(dextre.Blacklists) * (*warning) / 100
	criticalAmount := len(dextre.Blacklists) * (*critical) / 100
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
			positive*100/len(dextre.Blacklists),
			len(dextre.Blacklists),
			fmt.Sprintf(
				"queried=%v positive=%v",
				queried,
				positive,
			),
		),
	)
}
