/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	clusterHealthSubsystem = "cluster_health_subsystem"
)

var (
	colors                     = []string{"green", "yellow", "red"}
	defaultClusterHealthLabels = []string{"cluster"}
)

type ClusterHealthResponse struct {
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
type ClusterResponses struct {
	ClusterStat []ClusterHealthResponse `json:"cluster_stats"`
}

type clusterHealthMetric struct {
	Type   prometheus.ValueType
	Desc   *prometheus.Desc
	Value  func(clusterHealth ClusterHealthResponse) float64
	Labels func(clusterName ClusterHealthResponse) []string
}

type clusterHealthStatusMetric struct {
	Type   prometheus.ValueType
	Desc   *prometheus.Desc
	Value  func(clusterHealth ClusterHealthResponse) float64
	Labels func(clusterName, color string) []string
}

type ClusterHealth struct {
	logger logger.Logger
	client *http.Client
	url    *url.URL

	up                prometheus.Gauge
	totalScrapes      prometheus.Counter
	jsonParseFailures prometheus.Counter

	metrics      []*clusterHealthMetric
	statusMetric *clusterHealthStatusMetric
}

func NewClusterHealth(logger logger.Logger, client *http.Client, url *url.URL) prometheus.Collector {
	return &ClusterHealth{
		logger: logger,
		client: client,
		url:    url,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "up"),
			Help: "Was the last scrape of the Pixiu cluster health  successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "total_scrapes"),
			Help: "Current total Pixiu cluster health scrapes.",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "json_parse_failures"),
			Help: "Number of errors while parsing JSON.",
		}),

		metrics: []*clusterHealthMetric{
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "active_primary_shards"),
					"The number of primary shards in your cluster. This is an aggregate total across all indices.",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.ActivePrimaryShards)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "active_shards"),
					"Aggregate total of all shards across all indices, which includes replica shards.",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.ActiveShards)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "delayed_unassigned_shards"),
					"Shards delayed to reduce reallocation overhead",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.DelayedUnassignedShards)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "initializing_shards"),
					"Count of shards that are being freshly created.",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.InitializingShards)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "number_of_data_nodes"),
					"Number of data nodes in the cluster.",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.NumberOfDataNodes)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "number_of_in_flight_fetch"),
					"The number of ongoing shard info requests.",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.NumberOfInFlightFetch)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "task_max_waiting_in_queue_millis"),
					"Tasks max time waiting in queue.",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.TaskMaxWaitingInQueueMillis)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "number_of_nodes"),
					"Number of nodes in the cluster.",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.NumberOfNodes)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "number_of_pending_tasks"),
					"Cluster level changes which have not yet been executed",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.NumberOfPendingTasks)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "relocating_shards"),
					"The number of shards that are currently moving from one node to another node.",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.RelocatingShards)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
			{
				Type: prometheus.GaugeValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "unassigned_shards"),
					"The number of shards that exist in the cluster state, but cannot be found in the cluster itself.",
					defaultClusterHealthLabels, nil,
				),
				Value: func(clusterHealth ClusterHealthResponse) float64 {
					return float64(clusterHealth.UnassignedShards)
				},
				Labels: func(clusterHealth ClusterHealthResponse) []string {
					return []string{
						clusterHealth.ClusterName,
					}
				},
			},
		},
		statusMetric: &clusterHealthStatusMetric{
			Type: prometheus.GaugeValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(config.Namespace, clusterHealthSubsystem, "different_cluster_stats_total"),
				"Number of Different Cluster Stat",
				[]string{"cluster", "color"}, nil,
			),
			Value: func(clusterHealth ClusterHealthResponse) float64 {
				return 1
			},
			Labels: func(clusterName string, color string) []string {
				return []string{
					clusterName,
					color,
				}
			},
		},
	}
}

func (c *ClusterHealth) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.metrics {
		ch <- metric.Desc
	}
	ch <- c.statusMetric.Desc
	ch <- c.up.Desc()
	ch <- c.totalScrapes.Desc()
	ch <- c.jsonParseFailures.Desc()
}

func (c *ClusterHealth) Collect(ch chan<- prometheus.Metric) {
	var err error
	c.totalScrapes.Inc()
	defer func() {
		ch <- c.up
		ch <- c.totalScrapes
		ch <- c.jsonParseFailures
	}()

	statResps, err := c.FetchAndDecodeStats()
	if err != nil {
		c.up.Set(0)
		c.logger.Infof("watching (msg:{%s}) = error{%v}", "Failed to fetch and decode Cluster  Status", err.Error())
		return
	}

	c.up.Set(1)

	for _, metric := range c.metrics {
		for _, v := range statResps.ClusterStat {

			ch <- prometheus.MustNewConstMetric(
				metric.Desc,
				metric.Type,
				metric.Value(v),
				metric.Labels(v)...,
			)
		}
	}

	for _, color := range colors {
		for _, v := range statResps.ClusterStat {
			if v.Status == color {
				ch <- prometheus.MustNewConstMetric(
					c.statusMetric.Desc,
					c.statusMetric.Type,
					c.statusMetric.Value(v),
					c.statusMetric.Labels(v.ClusterName, color)...,
				)
			}
		}
	}
}

func (c *ClusterHealth) FetchAndDecodeStats() (ClusterResponses, error) {
	var chr ClusterResponses

	u := *c.url
	u.Path = path.Join(u.Path, "/_cluster/health")
	res, err := c.client.Get(u.String())
	if err != nil {
		return chr, fmt.Errorf("failed to get cluster health from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}

	defer func() {
		err = res.Body.Close()
		if err != nil {
			c.logger.Infof("watching (msg:{%s}) = error{%v}", "failed to close http.Client", err.Error())

		}
	}()

	if res.StatusCode != http.StatusOK {
		return chr, fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		c.jsonParseFailures.Inc()
		return chr, err
	}

	if err := json.Unmarshal(bts, &chr); err != nil {
		c.jsonParseFailures.Inc()
		return chr, err
	}

	return chr, nil
}

func NewHandler(logger logger.Logger, client *http.Client, url *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		registry.MustRegister(NewClusterHealth(logger, client, url))
		gatherers := prometheus.Gatherers{
			registry,
		}
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}
