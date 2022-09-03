package scrape

import (
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/prometheus/client_golang/prometheus"
)

type Scraper interface {
	// Name of the Scraper.
	Name() string

	// Help describes the role of the Scraper.
	Help() string

	// Scrape collects data from database connection and sends it over channel as prometheus metric.
	Scrape(ctx *contextHttp.HttpContext, ch chan<- prometheus.Metric) error
}
