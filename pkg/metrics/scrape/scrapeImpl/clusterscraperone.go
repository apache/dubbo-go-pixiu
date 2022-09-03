package scrapeImpl

import (
	"encoding/json"
	"strings"
)

import (
	contextHttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/global"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	clusterhealthscraper = "cluster_health_scraper"
)

var (
	defaultclusterlabels  = []string{"cluster"}
	clusterlabelValueFunc = func(clusterHealth ClusterStat) []string {
		return []string{
			clusterHealth.ClusterName,
		}
	}
	clustermetrics = []*clusterHealthMetric{
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "active_primary_shards"),
				"The number of primary shards in your cluster. This is an aggregate total across all indices.",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.ActivePrimaryShards)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "active_shards"),
				"Aggregate total of all shards across all indices, which includes replica shards.",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.ActiveShards)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "delayed_unassigned_shards"),
				"Shards delayed to reduce reallocation overhead",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.DelayedUnassignedShards)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "initializing_shards"),
				"Count of shards that are being freshly created.",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.InitializingShards)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "number_of_data_nodes"),
				"Number of data nodes in the cluster.",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.NumberOfDataNodes)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "number_of_in_flight_fetch"),
				"The number of ongoing shard info requests.",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.NumberOfInFlightFetch)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "task_max_waiting_in_queue_millis"),
				"Tasks max time waiting in queue.",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.TaskMaxWaitingInQueueMillis)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "number_of_nodes"),
				"Number of nodes in the cluster.",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.NumberOfNodes)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "number_of_pending_tasks"),
				"Cluster level changes which have not yet been executed",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.NumberOfPendingTasks)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "relocating_shards"),
				"The number of shards that are currently moving from one node to another node.",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.RelocatingShards)
			},
			Labels: clusterlabelValueFunc,
		},
		{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, clusterhealthscraper, "unassigned_shards"),
				"The number of shards that exist in the cluster state, but cannot be found in the cluster itself.",
				defaultclusterlabels, nil,
			),
			Value: func(clusterHealth ClusterStat) float64 {
				return float64(clusterHealth.UnassignedShards)
			},
			Labels: clusterlabelValueFunc,
		},
	}
)

type clusterHealthMetric struct {
	Type   prometheus.ValueType
	Desc   *prometheus.Desc
	Value  func(clusterHealth ClusterStat) float64
	Labels func(clusterName ClusterStat) []string
}

type ClusterStat struct {
	ClusterName                 string  `json:"cluster_name"`
	Status                      string  `json:"status"`
	TimedOut                    bool    `json:"timed_out"`
	NumberOfNodes               int     `json:"number_of_nodes"`
	NumberOfDataNodes           int     `json:"number_of_data_nodes"`
	ActivePrimaryShards         int     `json:"active_primary_shards"`
	ActiveShards                int     `json:"active_shards"`
	RelocatingShards            int     `json:"relocating_shards"`
	InitializingShards          int     `json:"initializing_shards"`
	UnassignedShards            int     `json:"unassigned_shards"`
	DelayedUnassignedShards     int     `json:"delayed_unassigned_shards"`
	NumberOfPendingTasks        int     `json:"number_of_pending_tasks"`
	NumberOfInFlightFetch       int     `json:"number_of_in_flight_fetch"`
	TaskMaxWaitingInQueueMillis int     `json:"task_max_waiting_in_queue_millis"`
	ActiveShardsPercentAsNumber float64 `json:"active_shards_percent_as_number"`
}

type ClusterStatsResponses struct {
	ClusterStat []ClusterStat `json:"cluster_stats"`
}

type ClusterScraper struct {
}

func (ClusterScraper) Name() string {
	return apihealthscraper
}

func (ClusterScraper) Help() string {
	return "The  Scraper for collect metric about Cluster"
}

func (ClusterScraper) Scrape(ctx *contextHttp.HttpContext, ch chan<- prometheus.Metric) error {
	var data ClusterStatsResponses
	path := ctx.Request.URL.RequestURI()
	if strings.HasPrefix(path, "/_cluster/health") {
		_ = json.NewDecoder(ctx.Request.Body).Decode(&data)
		for _, metric := range clustermetrics {
			for _, v := range data.ClusterStat {
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
