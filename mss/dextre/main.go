package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"

	"github.com/olorin/nagiosplugin"
	"github.com/poka-yoke/spaceflight/mss/dextre/dnsbl"
)

func fromFile(path string) <-chan string {
	out := make(chan string)
	go func() {
		blfile, err := os.Open(path)
		if err != nil {
			log.Fatal("Could't open file ", path)
			log.Fatal(err)
		}
		defer blfile.Close()

		scanner := bufio.NewScanner(blfile)
		for scanner.Scan() {
			out <- scanner.Text()
		}
		close(out)
	}()
	return out
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

	blacklists := fromFile(*blacklist)
	responses := make(chan int)

	queried := 0
	positive := 0
	length := 0
	for list := range blacklists {
		go dnsbl.Query(*ipAddress, list, responses)
		length++
	}

	for i := 0; i < length; i++ {
		response := <-responses
		if response > 0 {
			positive += response
		}
		queried++
	}

	warningAmount := length * (*warning) / 100
	criticalAmount := length * (*critical) / 100

	check := nagiosplugin.NewCheck()
	defer check.Finish()
	check.AddPerfDatum("queried", "", float64(queried), 0.0, math.Inf(1))
	check.AddPerfDatum("positive", "", float64(positive), 0.0, math.Inf(1))
	check.AddResultf(
		nagiosplugin.OK,
		"%v present in %v(%v%%) out of %v BLs",
		*ipAddress,
		positive,
		positive*100/length,
		length,
	)
	switch {
	case positive > warningAmount:
		check.AddResult(nagiosplugin.WARNING, "")
	case positive > criticalAmount:
		check.AddResultf(nagiosplugin.CRITICAL, "")
	}
}
