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

import "github.com/apache/dubbo-go-pixiu/pkg/logger"

func TestNewCommonMetricExporter(t *testing.T) {
	tcs := []string{

		`{"backend_stats":[{"name":"http://service1.xxx.cn/"}],"fronted_stats":[{"name":"http://client1.xxx.cn/"},{"name":"http://client2.xxx.cn/"},{"name":"http://client3.xxx.cn/"},{"name":"http://client4.xxx.cn/"},{"name":"http://client5.xxx.cn/"},{"name":"http://client6.xxx.cn/"},{"name":"http://client7.xxx.cn/"},{"name":"http://client8.xxx.cn/"},{"name":"http://client9.xxx.cn/"},{"name":"http://client10.xxx.cn/"}],"server_stats":[{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"}]}`,
		`{"backend_stats":[{"name":"http://service1.xxx.cn/"}],"fronted_stats":[{"name":"http://client11.xxx.cn/"},{"name":"http://client12.xxx.cn/"},{"name":"http://client13.xxx.cn/"},{"name":"http://client14.xxx.cn/"},{"name":"http://client15.xxx.cn/"},{"name":"http://client16.xxx.cn/"},{"name":"http://client17.xxx.cn/"},{"name":"http://client18.xxx.cn/"},{"name":"http://client19.xxx.cn/"},{"name":"http://client20.xxx.cn/"}],"server_stats":[{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"},{"server":"http://0.0.0.0.0"}]}`}

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
		c := NewCommonMetricExporter(log, http.DefaultClient, u)
		_, err = c.fetchAndDecodeStats()
		if err != nil {
			t.Fatalf("Failed to fetch or decode backendMetrics、 frontendMetrics、 and serverMetrics Stat : %s", err)
		}

	}
}
