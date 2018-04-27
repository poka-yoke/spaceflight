package cloudinary

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// Collector implements a Prometheus Collector to collect
// Cloudinary Usage Report statistics.
type Collector struct {
	transformations_usage_amount prometheus.Gauge
	transformations_limit_amount prometheus.Gauge
	transformations_usage_ratio  prometheus.Gauge
	objects_usage_amount         prometheus.Gauge
	objects_limit_amount         prometheus.Gauge
	objects_usage_ratio          prometheus.Gauge
	bandwidth_usage_bytes        prometheus.Gauge
	bandwidth_limit_bytes        prometheus.Gauge
	bandwidth_usage_ratio        prometheus.Gauge
	storage_usage_bytes          prometheus.Gauge
	storage_limit_bytes          prometheus.Gauge
	storage_usage_ratio          prometheus.Gauge
	requests_amount              prometheus.Gauge
	resources_amount             prometheus.Gauge
	derived_resources_amount     prometheus.Gauge
}

// NewCollector creates a new, default configured Collector
func NewCollector() *Collector {
	return &Collector{
		transformations_usage_amount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "transformations_usage_amount",
				Help:        "Number of used transformations in the last 30 days",
				ConstLabels: nil,
			},
		),
		transformations_limit_amount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "transformations_limit_amount",
				Help:        "Limit of transformations allowed in the last 30 days",
				ConstLabels: nil,
			},
		),
		transformations_usage_ratio: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "transformations_usage_ratio",
				Help:        "Ratio of used transformations over corresponding limit",
				ConstLabels: nil,
			},
		),
		objects_usage_amount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "objects_usage_amount",
				Help:        "Number of used objects in the last 30 days",
				ConstLabels: nil,
			},
		),
		objects_limit_amount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "objects_limit_amount",
				Help:        "Limit of objects allowed in the last 30 days",
				ConstLabels: nil,
			},
		),
		objects_usage_ratio: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "objects_usage_ratio",
				Help:        "Ratio of used objects over corresponding limit",
				ConstLabels: nil,
			},
		),
		bandwidth_usage_bytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "bandwidth_usage_bytes",
				Help:        "Bytes used in bandwidth in the last 30 days",
				ConstLabels: nil,
			},
		),
		bandwidth_limit_bytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "bandwidth_limit_bytes",
				Help:        "Limit of bytes in bandwidth to use in the last 30 days",
				ConstLabels: nil,
			},
		),
		bandwidth_usage_ratio: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "bandwidth_usage_ratio",
				Help:        "Ratio of used bytes in bandwidth over corresponding limit",
				ConstLabels: nil,
			},
		),
		storage_usage_bytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "storage_usage_bytes",
				Help:        "Bytes of storage used in the last 30 days",
				ConstLabels: nil,
			},
		),
		storage_limit_bytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "storage_limit_bytes",
				Help:        "Limit of storage used allowed in the last 30 days",
				ConstLabels: nil,
			},
		),
		storage_usage_ratio: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "storage_usage_ratio",
				Help:        "Ratio of storage bytes used over corresponding limit",
				ConstLabels: nil,
			},
		),
		requests_amount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "requests_amount",
				Help:        "Number of requests done to Cloudinary",
				ConstLabels: nil,
			},
		),
		resources_amount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "resources_amount",
				Help:        "Number of resources in Cloudinary",
				ConstLabels: nil,
			},
		),
		derived_resources_amount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "derived_resources_amount",
				Help:        "Number of derived resources in Cloudinary",
				ConstLabels: nil,
			},
		),
	}
}

// Describe is a requirement for the Collector interface of Prometheus
// that returns each exported metric's description to the Prometheus
// middleware
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.transformations_usage_amount.Describe(ch)
	c.transformations_limit_amount.Describe(ch)
	c.transformations_usage_ratio.Describe(ch)
	c.objects_usage_amount.Describe(ch)
	c.objects_limit_amount.Describe(ch)
	c.objects_usage_ratio.Describe(ch)
	c.bandwidth_usage_bytes.Describe(ch)
	c.bandwidth_limit_bytes.Describe(ch)
	c.bandwidth_usage_ratio.Describe(ch)
	c.storage_usage_bytes.Describe(ch)
	c.storage_limit_bytes.Describe(ch)
	c.storage_usage_ratio.Describe(ch)
	c.requests_amount.Describe(ch)
	c.resources_amount.Describe(ch)
	c.derived_resources_amount.Describe(ch)
}

// Collect is a requirement for the Collector interface of Prometheus
// that runs the queries to set the metrics values to be exported
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.runCollection()

	c.transformations_usage_amount.Collect(ch)
	c.transformations_limit_amount.Collect(ch)
	c.transformations_usage_ratio.Collect(ch)
	c.objects_usage_amount.Collect(ch)
	c.objects_limit_amount.Collect(ch)
	c.objects_usage_ratio.Collect(ch)
	c.bandwidth_usage_bytes.Collect(ch)
	c.bandwidth_limit_bytes.Collect(ch)
	c.bandwidth_usage_ratio.Collect(ch)
	c.storage_usage_bytes.Collect(ch)
	c.storage_limit_bytes.Collect(ch)
	c.storage_usage_ratio.Collect(ch)
	c.requests_amount.Collect(ch)
	c.resources_amount.Collect(ch)
	c.derived_resources_amount.Collect(ch)
}

func (c *Collector) runCollection() {
	cloudName, key, secret := CloudinaryCredentials.get()
	report, err := getUsageReport(
		fmt.Sprintf(
			"https://%s:%s@api.cloudinary.com/v1_1/%s/usage",
			key,
			secret,
			cloudName,
		),
	)
	if err != nil {
		fmt.Println(err)
	}

	c.transformations_usage_amount.Set(TransformationUsage(*report))
	c.transformations_limit_amount.Set(TransformationLimit(*report))
	c.transformations_usage_ratio.Set(TransformationUsageRatio(*report))
	c.objects_usage_amount.Set(ObjectsUsage(*report))
	c.objects_limit_amount.Set(ObjectsLimit(*report))
	c.objects_usage_ratio.Set(ObjectsUsageRatio(*report))
	c.bandwidth_usage_bytes.Set(BandwidthUsage(*report))
	c.bandwidth_limit_bytes.Set(BandwidthLimit(*report))
	c.bandwidth_usage_ratio.Set(BandwidthUsageRatio(*report))
	c.storage_usage_bytes.Set(StorageUsage(*report))
	c.storage_limit_bytes.Set(StorageLimit(*report))
	c.storage_usage_ratio.Set(StorageUsageRatio(*report))
	c.requests_amount.Set(Requests(*report))
	c.resources_amount.Set(Resources(*report))
	c.derived_resources_amount.Set(DerivedResources(*report))
}
