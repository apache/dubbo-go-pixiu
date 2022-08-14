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
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

import (
	"github.com/blang/semver/v4"
	"github.com/prometheus/client_golang/prometheus"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/global"
)

func init() {
	RegisterCollector("cluster-info", global.DefaultEnabled, NewClusterInfo)
}

type ClusterInfoCollector struct {
	logger logger.Logger
	u      *url.URL
	hc     *http.Client
}

func NewClusterInfo(logger logger.Logger, u *url.URL, hc *http.Client) (Collector, error) {
	return &ClusterInfoCollector{
		logger: logger,
		u:      u,
		hc:     hc,
	}, nil
}

var clusterInfoDesc = map[string]*prometheus.Desc{
	"version": prometheus.NewDesc(
		prometheus.BuildFQName(global.Namespace, "", "version"),
		"Pixiu version information.",
		[]string{
			"cluster",
			"cluster_uuid",
			"build_date",
			"build_hash",
			"version",
			"lucene_version",
		},
		nil,
	),
}

type ClusterInfoResponse struct {
	Name        string      `json:"name"`
	ClusterName string      `json:"cluster_name"`
	ClusterUUID string      `json:"cluster_uuid"`
	Version     VersionInfo `json:"version"`
	Tagline     string      `json:"tagline"`
}

type VersionInfo struct {
	Number        semver.Version `json:"number"`
	BuildHash     string         `json:"build_hash"`
	BuildDate     string         `json:"build_date"`
	BuildSnapshot bool           `json:"build_snapshot"`
	LuceneVersion semver.Version `json:"lucene_version"`
}

func (c *ClusterInfoCollector) Update(ctx context.Context, ch chan<- prometheus.Metric) error {
	resp, err := c.hc.Get(c.u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var info ClusterInfoResponse
	err = json.Unmarshal(b, &info)
	if err != nil {
		return err
	}

	ch <- prometheus.MustNewConstMetric(
		clusterInfoDesc["version"],
		prometheus.GaugeValue,
		1,
		info.ClusterName,
		info.ClusterUUID,
		info.Version.BuildDate,
		info.Version.BuildHash,
		info.Version.Number.String(),
		info.Version.LuceneVersion.String(),
	)

	return nil
}
