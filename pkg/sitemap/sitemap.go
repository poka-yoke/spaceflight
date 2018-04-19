package sitemap

import (
	"compress/gzip"
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Collector implements a Prometheus Collector to report sitemap
// metrics
type Collector struct {
	urls           []string
	timestamp, loc *prometheus.GaugeVec
}

// NewCollector creates a new, default configured Collector
func NewCollector(urls []string) *Collector {
	return &Collector{
		urls: urls,
		timestamp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   "sitemap",
				Name:        "last_modified_timestamp",
				Help:        "Unix timestamp of Last-Modified header response",
				ConstLabels: nil,
			},
			[]string{"sitemap"}, // The labels supported
		),
		loc: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace:   "sitemap",
				Name:        "loc_count",
				Help:        "Number of loc entries in the parsed sitemap",
				ConstLabels: nil,
			},
			[]string{"sitemap"}, // The labels supported
		),
	}
}

// Describe is a requirement for the Collector interface of Prometheus
// that returns each exported metric's description to the Prometheus
// middleware
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.timestamp.Describe(ch)
	c.loc.Describe(ch)
}

// Collect is a requirement for the Collector interface of Prometheus
// that runs the queries to set the metrics values to be exported
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	var wg sync.WaitGroup
	wg.Add(len(c.urls))

	for _, url := range c.urls {
		go func(url string) {
			// Get file from URL
			body := c.fetchBody(url)
			defer body.Close()

			reader, err := gzip.NewReader(body)
			if err != nil {
				panic(err)
			}
			defer reader.Close()
			c.loc.With(
				prometheus.Labels{
					"sitemap": nameFromURL(url),
				},
			).Set(float64(c.parseXML(reader)))
			wg.Done()
		}(url)
	}
	wg.Wait()

	c.timestamp.Collect(ch)
	c.loc.Collect(ch)
}

// fetchBody calls a url for a sitemap file. It captures the value of
// the Last-Modified header converted to Unix timestamp and passes the
// file contents for further processing.
func (c *Collector) fetchBody(url string) io.ReadCloser {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		log.Fatalf("Error contacting server %s: %s\n", url, err)
	}
	if resp.StatusCode != 200 {
		log.Fatalf("Server responded unsuccessfully: %s\n", resp.Status)
	}

	timestamp, err := time.Parse(
		"Mon, 2 Jan 2006 15:04:05 MST",
		resp.Header.Get("Last-Modified"),
	)
	if err != nil {
		log.Fatalf(
			"Couldn't parse timestamp %s: %s\n",
			resp.Header.Get("Last-Modified"),
			err,
		)
	}
	c.timestamp.With(
		prometheus.Labels{
			"sitemap": nameFromURL(url),
		},
	).Set(float64(timestamp.Unix()))
	return resp.Body
}

// parseXML gets a reader to an XML document and reads it. It captures
// the amount of `loc` tags present in the document.
func (c *Collector) parseXML(reader io.Reader) int {
	var loc int
	decoder := xml.NewDecoder(reader)
	for {
		t, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if e, ok := t.(xml.StartElement); ok && e.Name.Local == "loc" {
			loc++
		}
	}
	return loc
}

// nameFromURL takes an URL and returns the filename before extensions
// of the last member of the URL
func nameFromURL(url string) string {
	tmp := strings.Split(url, "/")
	return strings.Split(tmp[len(tmp)-1], ".")[0]
}
