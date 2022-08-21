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
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func TestApiHealth(t *testing.T) {

	tcs := []string{
		` {"share":{"total":1000,"successful":998,"failed":2},"api_stats":[{"api_name":"Api_1","api_requests":10,"api_requests_latency":1},{"api_name":"Api_12","api_requests":20,"api_requests_latency":2}]}`,
	}
	for _, out := range tcs {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintln(w, out)
		}))
		defer ts.Close()
		u, err := url.Parse(ts.URL)
		if err != nil {
			t.Fatalf("Failed to parse URL: %s", err)
		}
		log := logger.GetLogger()
		var chr ApiStatsResponse
		c := NewApiHealth(log, http.DefaultClient, u)
		if v, ok := (c).(*ApiHealth); ok {
			chr, err = v.FetchAndDecodeStats()
		}
		if err != nil {
			t.Fatalf("Failed to fetch or decode Api Stats: %s", err)
		}
		if chr.ApiStats[0].ApiName == "Api_1" {
			fmt.Println("ok")
		}

	}
}
