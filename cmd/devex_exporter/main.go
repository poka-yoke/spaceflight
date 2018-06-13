package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/Devex/spaceflight/pkg/dnsbl"
)

func main() {
	blacklist := flag.String("file", "Path to file containing black list addresses", "")
	ipAddress := flag.String("ip", "IP Address to look for in the BLs", "127.0.0.1")
	flag.Parse()

	blfile, err := os.Open(*blacklist)
	if err != nil {
		log.Fatal("Could't open file ", blacklist, err)
	}
	defer blfile.Close()

	providers := dnsbl.GetProviders(*ipAddress, blfile)
	prometheus.MustRegister(dnsbl.NewCollector(providers))

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
