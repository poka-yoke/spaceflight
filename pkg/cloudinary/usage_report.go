package cloudinary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type UsageInfo struct {
	Usage       int64   `json:"usage"`
	Limit       int64   `json:"limit"`
	UsedPercent float64 `json:"used_percent"`
}

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

func DerivedResources(usageReport UsageReport) float64 {
	fmt.Println(usageReport.DerivedResources)
	return float64(usageReport.DerivedResources)
}

func Resources(usageReport UsageReport) float64 {
	return float64(usageReport.Resources)
}

func Requests(usageReport UsageReport) float64 {
	return float64(usageReport.Requests)
}

func StorageUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Storage.Usage)
}

func StorageLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Storage.Limit)
}

func StorageUsageRatio(usageReport UsageReport) float64 {
	return float64(usageReport.Storage.UsedPercent / 100)
}

func BandwidthUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Bandwidth.Usage)
}

func BandwidthLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Bandwidth.Limit)
}

func BandwidthUsageRatio(usageReport UsageReport) float64 {
	return float64(usageReport.Bandwidth.UsedPercent / 100)
}

func ObjectsUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Objects.Usage)
}

func ObjectsLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Objects.Limit)
}

func ObjectsUsageRatio(usageReport UsageReport) float64 {
	return float64(usageReport.Objects.UsedPercent / 100)
}

func TransformationUsage(usageReport UsageReport) float64 {
	return float64(usageReport.Transformations.Usage)
}

func TransformationLimit(usageReport UsageReport) float64 {
	return float64(usageReport.Transformations.Limit)
}

func TransformationUsageRatio(usageReport UsageReport) float64 {
	return float64(usageReport.Transformations.UsedPercent / 100)
}
