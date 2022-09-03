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
	frontendhealthscraper = "frontend_health_scraper"
)

var (
	frontendLabelNames = []string{"frontend"}
	labelValueFunc     = func(stat FrontendStat) []string {
		return []string{
			stat.Name,
		}
	}

	frontendmetrics = []*frontendHealthMetric{
		{
			Type: prometheus.CounterValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, "frontend", "requests_denied_total"),
				"Total of requests denied for security.",
				frontendLabelNames,
				nil,
			),
			Value: func(stat FrontendStat) float64 {
				return float64(stat.RequestsDeniedTotal)
			},
			Labels: labelValueFunc,
		},
		{
			Type: prometheus.CounterValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, "frontend", "request_errors_total"),
				"Total of request errors.",
				frontendLabelNames,
				nil,
			),
			Value: func(stat FrontendStat) float64 {
				return float64(stat.RequestErrorsTotal)
			},
			Labels: labelValueFunc,
		},
	}
)

type frontendHealthMetric struct {
	Type   prometheus.ValueType
	Desc   *prometheus.Desc
	Value  func(stat FrontendStat) float64
	Labels func(stat FrontendStat) []string
}

type FrontendStatsResponse struct {
	FrontedStats FrontendStat
}

type FrontendStat struct {
	Name                string
	RequestsDeniedTotal int
	RequestErrorsTotal  int
	Responses1XXTotal   int
	Responses2XXTotal   int
	Responses3XXTotal   int
	Responses4XXTotal   int
	Responses5XXTotal   int
	ConnectionsTotal    int
}

type FtontendScraper struct{}

func (FtontendScraper) Name() string {
	return frontendhealthscraper
}

func (FtontendScraper) Help() string {
	return "The  Scraper for collect metric about Frontend"
}

func (FtontendScraper) Scrape(ctx *contextHttp.HttpContext, ch chan<- prometheus.Metric) error {
	var data FrontendStatsResponse
	path := ctx.Request.URL.RequestURI()
	if strings.HasPrefix(path, "/_frontend/health") {
		_ = json.NewDecoder(ctx.Request.Body).Decode(&data)
		for _, metric := range frontendmetrics {
			ch <- prometheus.MustNewConstMetric(
				metric.Desc,
				metric.Type,
				metric.Value(data.FrontedStats),
				metric.Labels(data.FrontedStats)...,
			)
		}
	}
	return nil
}
