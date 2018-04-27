package cloudinary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// UsageInfo represents the three different metrics for Transformations,
// Object, Bandwidth, and Storage.
type UsageInfo struct {
	Usage       int64   `json:"usage"`
	Limit       int64   `json:"limit"`
	UsedPercent float64 `json:"used_percent"`
}

// UsageReport represents the response from Cloudinary's Admin API on
// usage report.
type UsageReport struct {
	Plan             string    `json:"plan"`
	LastUpdate       string    `json:"last_updated"`
	Transformations  UsageInfo `json:"transformations"`
	Objects          UsageInfo `json:"objects"`
	Bandwidth        UsageInfo `json:"bandwidth"`
	Storage          UsageInfo `json:"storage"`
	Requests         int64     `json:"requests"`
	Resources        int64     `json:"resources"`
	DerivedResources int64     `json:"derived_resources"`
}

func getUsageReport(url string) (usageReport *UsageReport, err error) {
	rs, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf(
			"ERROR: Request failure: %v",
			err,
		)
	}
	defer rs.Body.Close()

	if rs.StatusCode != 200 && rs.StatusCode != 201 {
		return nil, fmt.Errorf(
			"ERROR: Cloudinary API complained: %v",
			rs.Header.Get("X-Cld-Error"),
		)
	}
	bodyBytes, err := ioutil.ReadAll(rs.Body)
	if err != nil {
		return nil, err
	}

	usageReport = new(UsageReport)
	err = json.Unmarshal(bodyBytes, &usageReport)
	return usageReport, err
}

func derivedResources(usageReport UsageReport) float64 {
	return float64(usageReport.DerivedResources)
}

func resources(usageReport UsageReport) float64 {
	return float64(usageReport.Resources)
}

func requests(usageReport UsageReport) float64 {
	return float64(usageReport.Requests)
}

func storageUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Storage.Usage)
}

func storageLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Storage.Limit)
}

func storageUsageRatio(usageReport UsageReport) float64 {
	return float64(usageReport.Storage.UsedPercent / 100)
}

func bandwidthUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Bandwidth.Usage)
}

func bandwidthLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Bandwidth.Limit)
}

func bandwidthUsageRatio(usageReport UsageReport) float64 {
	return float64(usageReport.Bandwidth.UsedPercent / 100)
}

func objectsUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Objects.Usage)
}

func objectsLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Objects.Limit)
}

func objectsUsageRatio(usageReport UsageReport) float64 {
	return float64(usageReport.Objects.UsedPercent / 100)
}

func transformationsUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Transformations.Usage)
}

func transformationsLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Transformations.Limit)
}

func transformationsUsageRatio(usageReport UsageReport) float64 {
	return float64(usageReport.Transformations.UsedPercent / 100)
}
