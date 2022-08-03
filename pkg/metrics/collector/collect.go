package collector

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
)

import (
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/alecthomas/kingpin.v2"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/global"
)

// Collector is the interface a collector has to implement.
type Collector interface {
	// Get new metrics and expose them via prometheus registry.
	Update(context.Context, chan<- prometheus.Metric) error
}

type FactoryFunc func(logger log.Logger, u *url.URL, hc *http.Client) (Collector, error)

var (
	factories              = make(map[string]FactoryFunc)
	initiatedCollectorsMtx = sync.Mutex{}
	initiatedCollectors    = make(map[string]Collector)
	collectorState         = make(map[string]*bool)
	//collectors which have been explicitly enabled or disabled 明确启用或禁用的收集器
	forcedCollectors = map[string]bool{}
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(global.Namespace, "scrape", "duration_seconds"),
		"pixiu_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(global.Namespace, "scrape", "success"),
		"pixiu_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)

//Provide collector call
func RegisterCollector(name string, isDefaultEnabled bool, createFunc FactoryFunc) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	// Create flag for this collector
	flagName := fmt.Sprintf("collector.%s", name)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", name, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", isDefaultEnabled)
	// Dispatched to the given function after parsing and validating the flags.
	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Action(CollectorFlagAction(name)).Bool()
	collectorState[name] = flag

	// Register the create function for this collector
	factories[name] = createFunc
}

// CollectorFlagAction generates a new action function for the given collector
// to track whether it has been explicitly enabled or disabled from the command line.
// A new action function is needed for each collector flag because the ParseContext
// does not contain information about which flag called the action.
func CollectorFlagAction(collector string) func(ctx *kingpin.ParseContext) error {
	return func(ctx *kingpin.ParseContext) error {

		forcedCollectors[collector] = true
		return nil
	}
}
