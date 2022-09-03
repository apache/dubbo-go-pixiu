package scrapeImpl

import (
	"encoding/json"
	"strings"
)

import (
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/global"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Subsystem.
	apihealthscraper = "api_health_scraper"
)

var (
	apilabelValueFunc = func(apiStat ApiStat) []string {
		return []string{apiStat.ApiName}
	}
	apilabels  = []string{"api"}
	apimetrics = []*apiMetric{
		{
			Type: prometheus.CounterValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, "api_stats", "api_requests_total"),
				"Number of Api Requests",
				apilabels, nil,
			),
			Value: func(apiStats ApiStat) float64 {
				return float64(apiStats.ApiRequests)
			},
			Labels: apilabelValueFunc,
		},
	}
)

type apiMetric struct {
	Type   prometheus.ValueType
	Desc   *prometheus.Desc
	Value  func(dataStreamStats ApiStat) float64
	Labels func(dataStreamStats ApiStat) []string
}

type ApiStatsResponse struct {
	ApiStats []ApiStat `json:"api_stats"`
}

type ApiStat struct {
	ApiName     string `json:"api_name"`
	ApiRequests int64  `json:"api_requests"`
}

type ApiScraper struct{}

func (ApiScraper) Name() string {
	return apihealthscraper
}

func (ApiScraper) Help() string {
	return "The  Scraper for collect metric about API"
}

func (ApiScraper) Scrape(ctx *contextHttp.HttpContext, ch chan<- prometheus.Metric) error {
	var data ApiStatsResponse
	path := ctx.Request.URL.RequestURI()
	if strings.HasPrefix(path, "/_api/health") {
		_ = json.NewDecoder(ctx.Request.Body).Decode(&data)
		for _, metric := range apimetrics {
			for _, v := range data.ApiStats {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc,
					metric.Type,
					metric.Value(v),
					metric.Labels(v)...,
				)

			}

		}
	}
	return nil
}
