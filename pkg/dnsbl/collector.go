package dnsbl

import (
	"net"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Collector implements a Prometheus Collector to report DNS blacklist
// values
type Collector struct {
	providers               []string
	positive, query, length prometheus.Gauge

	// Functional control
	wg sync.WaitGroup
}

// NewCollector creates a new, default configured Collector
func NewCollector(providers []string) *Collector {
	return &Collector{
		providers: providers,
		positive: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "dnsbl",
				Name:        "positive_count",
				Help:        "Number of positive results of DNS blacklists",
				ConstLabels: nil,
			},
		),
		query: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "dnsbl",
				Name:        "query_count",
				Help:        "Number of DNS blacklists providers that answered",
				ConstLabels: nil,
			},
		),
		length: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "dnsbl",
				Name:        "length_count",
				Help:        "Number of DNS blacklists contacted",
				ConstLabels: nil,
			},
		),
	}
}

// Describe is a requirement for the Collector interface of Prometheus
// that returns each exported metric's description to the Prometheus
// middleware
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.positive.Describe(ch)
	c.query.Describe(ch)
	c.length.Describe(ch)
}

// Collect is a requirement for the Collector interface of Prometheus
// that runs the queries to set the metrics values to be exported
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.runCollection()

	c.positive.Collect(ch)
	c.query.Collect(ch)
	c.length.Collect(ch)
}

func (c *Collector) runCollection() {
	length := len(c.providers)
	responses := make(chan int, length)
	c.wg.Add(length)
	go func() {
		c.positive.Set(0)
		c.query.Set(0)
		for response := range responses {
			c.positive.Add(float64(response))
			c.query.Inc()
			c.wg.Done()
		}
	}()
	for _, provider := range c.providers {
		go func(provider string) {
			responses <- query(c, provider)
		}(provider)
	}
	c.length.Set(float64(length))
	c.wg.Wait()
	close(responses)
}

func (c *Collector) lookup(address string) ([]string, error) {
	return net.LookupHost(address)
}
