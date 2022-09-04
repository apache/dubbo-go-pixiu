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
	"github.com/stretchr/testify/assert"
)

func TestExporterApiMetric(t *testing.T) {

	rules := MetricCollectConfiguration{
		MetricCollectRule{
			Enable:              true,
			MeticPath:           "/metrics",
			PushGatewayURL:      "http://domain:port",
			PushIntervalSeconds: 3,
			PushJobName:         "prometheus",
		},
	}
	p := Plugin{
		Cfg: &rules,
	}

	metricFilter, _ := p.CreateFilterFactory()
	config := metricFilter.Config()
	mockYaml, err := yaml.MarshalYML(rules)
	assert.Nil(t, err)
	err = yaml.UnmarshalYML(mockYaml, config)
	assert.Nil(t, err)
	err = metricFilter.Apply()
	assert.Nil(t, err)
	f := &Filter{
		Cfg: &rules,
	}
	data := GetApiStatsResponse()
	body, _ := json.Marshal(&data)
	request, _ := http.NewRequest("POST", "/_api/health", bytes.NewBuffer(body))
	ctx := mock.GetMockHTTPContext(request)
	filterStatus := f.Decode(ctx)
	assert.Equal(t, filterStatus, filter.Continue)

}

func GetApiStatsResponse() ApiStatsResponse {
	return ApiStatsResponse{
		ApiStats: []ApiStat{
			{
				ApiName:     "api1",
				ApiRequests: 1000,
			},
		},
	}
}

type ApiStatsResponse struct {
	ApiStats []ApiStat `json:"api_stats"`
}

type ApiStat struct {
	ApiName     string `json:"api_name"`
	ApiRequests int64  `json:"api_requests"`
}
