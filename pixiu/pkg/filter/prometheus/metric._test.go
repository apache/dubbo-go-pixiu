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

package prometheus

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)
import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/scrape/scrapeImpl"
	"github.com/stretchr/testify/assert"
)

func TestExporterApiMetric(t *testing.T) {

	rules := MetricCollectConfiguration{
		MetricCollectRule{Enable: true, MeticPath: "/metrics", ExporterPort: 9914},
	}
	p := Plugin{}
	metricFilter, _ := p.CreateFilterFactory()
	config := metricFilter.Config()
	mockYaml, err := yaml.MarshalYML(rules)
	assert.Nil(t, err)
	err = yaml.UnmarshalYML(mockYaml, config)
	assert.Nil(t, err)

	err = metricFilter.Apply()
	assert.Nil(t, err)
	f := &Filter{Cfg: &rules}
	data := GetApiStatsResponse()
	body, _ := json.Marshal(&data)
	request, _ := http.NewRequest("POST", "/_api/health", bytes.NewBuffer(body))
	ctx := mock.GetMockHTTPContext(request)
	filterStatus := f.Decode(ctx)
	assert.Equal(t, filterStatus, filter.Continue)

}

func TestExporterClusterMetric(t *testing.T) {

	rules := MetricCollectConfiguration{
		MetricCollectRule{Enable: true, MeticPath: "/metrics", ExporterPort: 9914},
	}
	p := Plugin{}
	metricFilter, _ := p.CreateFilterFactory()
	config := metricFilter.Config()
	mockYaml, err := yaml.MarshalYML(rules)
	assert.Nil(t, err)
	err = yaml.UnmarshalYML(mockYaml, config)
	assert.Nil(t, err)

	err = metricFilter.Apply()
	assert.Nil(t, err)
	f := &Filter{Cfg: &rules}
	data := GetClusterStatsResponses()
	body, _ := json.Marshal(&data)
	request, _ := http.NewRequest("POST", "/_cluster/health", bytes.NewBuffer(body))
	ctx := mock.GetMockHTTPContext(request)
	filterStatus := f.Decode(ctx)
	assert.Equal(t, filterStatus, filter.Continue)

}

func TestExporterFrontendMetric(t *testing.T) {

	rules := MetricCollectConfiguration{
		MetricCollectRule{Enable: true, MeticPath: "/metrics", ExporterPort: 9914},
	}
	p := Plugin{}
	metricFilter, _ := p.CreateFilterFactory()
	config := metricFilter.Config()
	mockYaml, err := yaml.MarshalYML(rules)
	assert.Nil(t, err)
	err = yaml.UnmarshalYML(mockYaml, config)
	assert.Nil(t, err)

	err = metricFilter.Apply()
	assert.Nil(t, err)
	f := &Filter{Cfg: &rules}
	data := GetFrontendStatsResponses()
	body, _ := json.Marshal(&data)
	request, _ := http.NewRequest("POST", "/_frontend/health", bytes.NewBuffer(body))
	ctx := mock.GetMockHTTPContext(request)
	filterStatus := f.Decode(ctx)
	assert.Equal(t, filterStatus, filter.Continue)

}

func GetApiStatsResponse() scrapeImpl.ApiStatsResponse {
	return scrapeImpl.ApiStatsResponse{
		ApiStats: []scrapeImpl.ApiStat{
			{
				ApiName:     "api1",
				ApiRequests: 1000,
			},
		},
	}
}

func GetClusterStatsResponses() scrapeImpl.ClusterStatsResponses {
	return scrapeImpl.ClusterStatsResponses{
		ClusterStat: []scrapeImpl.ClusterStat{
			{
				ClusterName:         "cluster1",
				Status:              "yellow",
				TimedOut:            true,
				NumberOfNodes:       100,
				NumberOfDataNodes:   1000,
				ActivePrimaryShards: 200,
				ActiveShards:        300,
				RelocatingShards:    800,
				InitializingShards:  77,
				UnassignedShards:    99,
			},
		},
	}
}

func GetFrontendStatsResponses() scrapeImpl.FrontendStatsResponse {
	return scrapeImpl.FrontendStatsResponse{
		FrontedStats: scrapeImpl.FrontendStat{
			Name:                "front2",
			RequestsDeniedTotal: 1001,
			RequestErrorsTotal:  2002,
			Responses1XXTotal:   2389,
			Responses2XXTotal:   4444,
			Responses3XXTotal:   9990,
			Responses4XXTotal:   7654,
			Responses5XXTotal:   9087,
			ConnectionsTotal:    5622,
		},
	}
}
