package main

import (
	"flag"

	"github.com/Devex/spaceflight/mss/sidekiq"
)

func main() {
	baseURL := flag.String("url", "", "Base URL for sidekiq check API")
	flag.Parse()

	info := sidekiq.ProcessGetResponse(*baseURL)
	check := info.NagiosCheck()
	defer check.Finish()
}
