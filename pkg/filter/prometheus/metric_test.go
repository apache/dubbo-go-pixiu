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
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
)

func TestCounterExporterApiMetric(t *testing.T) {
	rules := MetricCollectConfiguration{
		MetricCollectRule{
			MetricPath:            "/metrics",
			PushGatewayURL:        "http://127.0.0.1:9091",
			PushJobName:           "pixiu",
			CounterPush:           true,
			PushIntervalThreshold: 100,
		},
	}
	_, err := yaml.MarshalYML(rules)
	assert.Nil(t, err)
	config := &rules
	p := Plugin{}
	msg := "this is test msg"
	metricFilterFactory, _ := p.CreateFilterFactory()
	if factory, ok := metricFilterFactory.(*FilterFactory); ok {
		factory.Cfg = config
		err = factory.Apply()
		assert.Nil(t, err)
		chain := filter.NewDefaultFilterChain()
		for i := 0; i < 1000; i++ {
			data := GetApiStatsResponse()
			body, _ := json.Marshal(&data)
			request, _ := http.NewRequest("POST", "/_api/health", bytes.NewBuffer(body))
			ctx := mock.GetMockHTTPContext(request)
			ctx.TargetResp = client.NewResponse([]byte(msg))
			err := factory.PrepareFilterChain(ctx, chain)
			assert.Nil(t, err)
			chain.OnDecode(ctx)
		}
	}
	time.Sleep(20 * time.Second)
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
