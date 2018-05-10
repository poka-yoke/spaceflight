package cloudinary

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// Collector implements a Prometheus Collector to collect
// Cloudinary Usage Report statistics.
type Collector struct {
	transformationsUsageAmount prometheus.Gauge
	transformationsLimitAmount prometheus.Gauge
	transformationsUsageRatio  prometheus.Gauge
	objectsUsageAmount         prometheus.Gauge
	objectsLimitAmount         prometheus.Gauge
	objectsUsageRatio          prometheus.Gauge
	bandwidthUsageBytes        prometheus.Gauge
	bandwidthLimitBytes        prometheus.Gauge
	bandwidthUsageRatio        prometheus.Gauge
	storageUsageBytes          prometheus.Gauge
	storageLimitBytes          prometheus.Gauge
	storageUsageRatio          prometheus.Gauge
	requestsAmount             prometheus.Gauge
	resourcesAmount            prometheus.Gauge
	derivedResourcesAmount     prometheus.Gauge
}

// NewCollector creates a new, default configured Collector
func NewCollector() *Collector {
	return &Collector{
		transformationsUsageAmount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "transformations_usage_amount",
				Help:        "Number of used transformations in the last 30 days",
				ConstLabels: nil,
			},
		),
		transformationsLimitAmount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "transformations_limit_amount",
				Help:        "Limit of transformations allowed in the last 30 days",
				ConstLabels: nil,
			},
		),
		transformationsUsageRatio: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "transformations_usage_ratio",
				Help:        "Ratio of used transformations over corresponding limit",
				ConstLabels: nil,
			},
		),
		objectsUsageAmount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "objects_usage_amount",
				Help:        "Number of used objects in the last 30 days",
				ConstLabels: nil,
			},
		),
		objectsLimitAmount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "objects_limit_amount",
				Help:        "Limit of objects allowed in the last 30 days",
				ConstLabels: nil,
			},
		),
		objectsUsageRatio: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "objects_usage_ratio",
				Help:        "Ratio of used objects over corresponding limit",
				ConstLabels: nil,
			},
		),
		bandwidthUsageBytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "bandwidth_usage_bytes",
				Help:        "Bytes used in bandwidth in the last 30 days",
				ConstLabels: nil,
			},
		),
		bandwidthLimitBytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "bandwidth_limit_bytes",
				Help:        "Limit of bytes in bandwidth to use in the last 30 days",
				ConstLabels: nil,
			},
		),
		bandwidthUsageRatio: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "bandwidth_usage_ratio",
				Help:        "Ratio of used bytes in bandwidth over corresponding limit",
				ConstLabels: nil,
			},
		),
		storageUsageBytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "storage_usage_bytes",
				Help:        "Bytes of storage used in the last 30 days",
				ConstLabels: nil,
			},
		),
		storageLimitBytes: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "storage_limit_bytes",
				Help:        "Limit of storage used allowed in the last 30 days",
				ConstLabels: nil,
			},
		),
		storageUsageRatio: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "storage_usage_ratio",
				Help:        "Ratio of storage bytes used over corresponding limit",
				ConstLabels: nil,
			},
		),
		requestsAmount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "requests_amount",
				Help:        "Number of requests done to Cloudinary",
				ConstLabels: nil,
			},
		),
		resourcesAmount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace:   "cloudinary",
				Name:        "resources_amount",
				Help:        "Number of resources in Cloudinary",
				ConstLabels: nil,
			},
		),
		derivedResourcesAmount: prometheus.NewGauge(
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
	c.transformationsUsageAmount.Describe(ch)
	c.transformationsLimitAmount.Describe(ch)
	c.transformationsUsageRatio.Describe(ch)
	c.objectsUsageAmount.Describe(ch)
	c.objectsLimitAmount.Describe(ch)
	c.objectsUsageRatio.Describe(ch)
	c.bandwidthUsageBytes.Describe(ch)
	c.bandwidthLimitBytes.Describe(ch)
	c.bandwidthUsageRatio.Describe(ch)
	c.storageUsageBytes.Describe(ch)
	c.storageLimitBytes.Describe(ch)
	c.storageUsageRatio.Describe(ch)
	c.requestsAmount.Describe(ch)
	c.resourcesAmount.Describe(ch)
	c.derivedResourcesAmount.Describe(ch)
}

// Collect is a requirement for the Collector interface of Prometheus
// that runs the queries to set the metrics values to be exported
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.runCollection()

	c.transformationsUsageAmount.Collect(ch)
	c.transformationsLimitAmount.Collect(ch)
	c.transformationsUsageRatio.Collect(ch)
	c.objectsUsageAmount.Collect(ch)
	c.objectsLimitAmount.Collect(ch)
	c.objectsUsageRatio.Collect(ch)
	c.bandwidthUsageBytes.Collect(ch)
	c.bandwidthLimitBytes.Collect(ch)
	c.bandwidthUsageRatio.Collect(ch)
	c.storageUsageBytes.Collect(ch)
	c.storageLimitBytes.Collect(ch)
	c.storageUsageRatio.Collect(ch)
	c.requestsAmount.Collect(ch)
	c.resourcesAmount.Collect(ch)
	c.derivedResourcesAmount.Collect(ch)
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

	c.transformationsUsageAmount.Set(transformationsUsage(*report))
	c.transformationsLimitAmount.Set(transformationsLimit(*report))
	c.transformationsUsageRatio.Set(transformationsUsageRatio(*report))
	c.objectsUsageAmount.Set(objectsUsage(*report))
	c.objectsLimitAmount.Set(objectsLimit(*report))
	c.objectsUsageRatio.Set(objectsUsageRatio(*report))
	c.bandwidthUsageBytes.Set(bandwidthUsage(*report))
	c.bandwidthLimitBytes.Set(bandwidthLimit(*report))
	c.bandwidthUsageRatio.Set(bandwidthUsageRatio(*report))
	c.storageUsageBytes.Set(storageUsage(*report))
	c.storageLimitBytes.Set(storageLimit(*report))
	c.storageUsageRatio.Set(storageUsageRatio(*report))
	c.requestsAmount.Set(requests(*report))
	c.resourcesAmount.Set(resources(*report))
	c.derivedResourcesAmount.Set(derivedResources(*report))
}
