package main

import (
	"flag"
	"log"
	"math"

	"github.com/olorin/nagiosplugin"
	"github.com/poka-yoke/spaceflight/mss/dextre/dnsbl"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
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

	list := dnsbl.FromFile(*blacklist)
	dnsbl.Queries(*ipAddress, list)

	positive := dnsbl.Stats.Positive
	queried := dnsbl.Stats.Queried
	length := dnsbl.Stats.Length

	check := nagiosplugin.NewCheck()
	defer check.Finish()
	must(check.AddPerfDatum("queried", "", float64(queried), 0.0, math.Inf(1)))
	must(check.AddPerfDatum("positive", "", float64(positive), 0.0, math.Inf(1)))
	check.AddResultf(
		nagiosplugin.OK,
		"%v present in %v(%v%%) out of %v BLs",
		*ipAddress,
		positive,
		positive*100/length,
		length,
	)
	switch {
	case positive > length*(*warning)/100:
		check.AddResult(nagiosplugin.WARNING, "")
	case positive > length*(*critical)/100:
		check.AddResultf(nagiosplugin.CRITICAL, "")
	}
}
