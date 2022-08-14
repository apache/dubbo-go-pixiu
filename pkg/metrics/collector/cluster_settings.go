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
	"strconv"
)

import (
	"github.com/imdario/mergo"
	"github.com/prometheus/client_golang/prometheus"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/global"
)

type ClusterSettings struct {
	logger logger.Logger
	client *http.Client
	url    *url.URL

	up                prometheus.Gauge
	maxShardsPerNode  prometheus.Gauge
	totalScrapes      prometheus.Counter
	jsonParseFailures prometheus.Counter
}

func NewClusterSettings(logger logger.Logger, client *http.Client, url *url.URL) *ClusterSettings {
	return &ClusterSettings{
		logger: logger,
		client: client,
		url:    url,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(global.Namespace, "clustersettings_stats", "up"),
			Help: "Was the last scrape of the Pixiu cluster settings endpoint successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(global.Namespace, "clustersettings_stats", "total_scrapes"),
			Help: "Current total Pixiu cluster settings scrapes.",
		}),
		maxShardsPerNode: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(global.Namespace, "clustersettings_stats", "max_shards_per_node"),
			Help: "Current maximum number of shards per node setting.",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(global.Namespace, "clustersettings_stats", "json_parse_failures"),
			Help: "Number of errors while parsing JSON.",
		}),
	}
}

func (cs *ClusterSettings) Describe(ch chan<- *prometheus.Desc) {
	ch <- cs.up.Desc()
	ch <- cs.totalScrapes.Desc()
	ch <- cs.maxShardsPerNode.Desc()
	ch <- cs.jsonParseFailures.Desc()
}

func (cs *ClusterSettings) getAndParseURL(u *url.URL, data interface{}) error {
	res, err := cs.client.Get(u.String())
	if err != nil {
		return fmt.Errorf("failed to get from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}

	defer func() {
		err = res.Body.Close()
		if err != nil {
			cs.logger.Infof("watching (msg:{%s}) = error{%v}", "failed to close http.Client", err.Error())
		}
	}()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		cs.jsonParseFailures.Inc()
		return err
	}

	if err := json.Unmarshal(bts, data); err != nil {
		cs.jsonParseFailures.Inc()
		return err
	}

	return nil
}

type ClusterSettingsResponse struct {
	MaxShardsPerNode interface{} `json:"max_shards_per_node"`
}

type ClusterSettingsFullResponse struct {
	Defaults   ClusterSettingsResponse `json:"defaults"`
	Persistent ClusterSettingsResponse `json:"persistent"`
	Transient  ClusterSettingsResponse `json:"transient"`
}

func (cs *ClusterSettings) fetchAndDecodeClusterSettingsStats() (ClusterSettingsResponse, error) {

	u := *cs.url
	u.Path = path.Join(u.Path, "/_cluster/settings")
	q := u.Query()
	q.Set("include_defaults", "true")
	u.RawQuery = q.Encode()
	u.RawPath = q.Encode()
	var csfr ClusterSettingsFullResponse
	var csr ClusterSettingsResponse
	err := cs.getAndParseURL(&u, &csfr)
	if err != nil {
		return csr, err
	}
	err = mergo.Merge(&csr, csfr.Defaults, mergo.WithOverride)
	if err != nil {
		return csr, err
	}
	err = mergo.Merge(&csr, csfr.Persistent, mergo.WithOverride)
	if err != nil {
		return csr, err
	}
	err = mergo.Merge(&csr, csfr.Transient, mergo.WithOverride)

	return csr, err
}

func (cs *ClusterSettings) Collect(ch chan<- prometheus.Metric) {

	cs.totalScrapes.Inc()
	defer func() {
		ch <- cs.up
		ch <- cs.totalScrapes
		ch <- cs.jsonParseFailures
		ch <- cs.maxShardsPerNode
	}()

	csr, err := cs.fetchAndDecodeClusterSettingsStats()
	if err != nil {
		cs.up.Set(0)
		cs.logger.Infof("watching (msg:{%s}) = error{%v}", "failed to fetch And Decode Cluster Settings Stats", err.Error())
		return
	}
	cs.up.Set(1)
	if maxShardsPerNodeString, ok := csr.MaxShardsPerNode.(string); ok {
		maxShardsPerNode, err := strconv.ParseInt(maxShardsPerNodeString, 10, 64)
		if err == nil {
			cs.maxShardsPerNode.Set(float64(maxShardsPerNode))
		}
	}
}
