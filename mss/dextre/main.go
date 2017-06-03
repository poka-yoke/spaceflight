package main

import (
	"bufio"
	"flag"
	"log"
	"math"
	"os"
	"sync"

	"github.com/olorin/nagiosplugin"
	"github.com/poka-yoke/spaceflight/mss/dextre/dnsbl"
)

var length = 0
var wg sync.WaitGroup

func fromFile(path string) <-chan string {
	out := make(chan string)
	go func() {
		blfile, err := os.Open(path)
		if err != nil {
			log.Fatal("Could't open file ", path, err)
		}
		defer blfile.Close()

		scanner := bufio.NewScanner(blfile)
		for scanner.Scan() {
			wg.Add(1)
			out <- scanner.Text()
			length++
		}
		close(out)
	}()
	return out
}

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

	list := fromFile(*blacklist)
	responses := dnsbl.Queries(*ipAddress, list)

	queried := 0
	positive := 0
	go func() {
		for response := range responses {
			if response > 0 {
				positive += response
			}
			queried++
			wg.Done()
		}
	}()
	wg.Wait()
	close(responses)

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
