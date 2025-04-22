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

package http

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	commonmock "github.com/apache/dubbo-go-pixiu/pkg/common/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/trie"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	DEMO = "dgp.filters.http.demo"
	// Kind is the kind of plugin.
	Kind = DEMO
)

var (
	eventCh = make(chan string, 3)
)

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// HeaderFilter is http filter instance
	DemoFilterFactory struct {
		conf *Config
	}
	DemoFilter struct {
		str string
	}
	// Config describe the config of ResponseFilter
	Config struct {
		Foo string `json:"foo,omitempty" yaml:"foo,omitempty"`
		Bar string `json:"bar,omitempty" yaml:"bar,omitempty"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &DemoFilterFactory{conf: &Config{Foo: "default foo", Bar: "default bar"}}, nil
}

func (f *DemoFilter) Decode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	logger.Info("decode phase: ", f.str)

	runes := []rune(f.str)
	for i := 0; i < len(runes)/2; i += 1 {
		runes[i], runes[len(runes)-1-i] = runes[len(runes)-1-i], runes[i]
	}
	f.str = string(runes)

	return filter.Continue
}

func (f *DemoFilter) Encode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	logger.Info("encode phase: ", f.str)
	return filter.Continue
}

func (f *DemoFilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	c := f.conf
	str := fmt.Sprintf("%s is drinking in the %s", c.Foo, c.Bar)
	demoFilter := &DemoFilter{str: str}

	chain.AppendDecodeFilters(demoFilter)
	chain.AppendEncodeFilters(demoFilter)
	return nil
}

func (f *DemoFilterFactory) Config() interface{} {
	return f.conf
}

func (f *DemoFilterFactory) Apply() error {
	return nil
}

func TestCreateHttpConnectionManager(t *testing.T) {
	filter.RegisterHttpFilter(&Plugin{})

	hcmc := model.HttpConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{
			RouteTrie: trie.NewTrieWithDefault("POST/api/v1/**", model.RouteAction{
				Cluster:                     "test_dubbo",
				ClusterNotFoundResponseCode: 505,
			}),
			Dynamic: false,
		},
		HTTPFilters: []*model.HTTPFilter{
			{
				Name:   DEMO,
				Config: nil,
			},
		},
		ServerName:        "test_http_dubbo",
		GenerateRequestID: false,
		IdleTimeoutStr:    "100",
	}

	hcm := CreateHttpConnectionManager(&hcmc)
	assert.Equal(t, len(hcm.filterManager.GetFactory()), 1)
	request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/api/v1?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	assert.NoError(t, err)
	request.Header = map[string][]string{
		"X-Dgp-Way": {"Dubbo"},
	}
	assert.NoError(t, err)
	c := mock.GetMockHTTPContext(request)
	err = hcm.findRoute(c)
	assert.NoError(t, err)
	err = hcm.Handle(c)
	assert.NoError(t, err)
}

// test SSE case
func TestStreamingResponse(t *testing.T) {
	hcmc := model.HttpConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{
			RouteTrie: trie.NewTrieWithDefault("GET/api/sse", model.RouteAction{
				Cluster: "mock_stream_cluster",
			}),
		},
		HTTPFilters: []*model.HTTPFilter{
			{
				Name: commonmock.Kind,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// mock server
	upstreamServer, _ := NewTestServerWithURL("localhost:8080", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		flusher := w.(http.Flusher)

		for i := 1; i <= 3; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(10 * time.Millisecond)
				event := fmt.Sprintf("data: %d\nevent: %d\nid: %d\n\n", i, i, i)
				_, _ = w.Write([]byte(event))
				flusher.Flush()
				logger.Info("Upstream sent event ", i)
			}
		}
	}))
	defer upstreamServer.Close()

	req := httptest.NewRequest("GET", "http://localhost:8080/api/sse", nil).WithContext(ctx)

	done := make(chan struct{})

	httpCtx := &contexthttp.HttpContext{
		Request: req,
		Writer:  NewStreamRecorder(),
		Ctx:     ctx,
	}
	go func() {
		defer close(done)

		hcm := CreateHttpConnectionManager(&hcmc)

		if err := hcm.Handle(httpCtx); err != nil {
			t.Errorf("Handle failed: %v", err)
		}

		// test targetResp
		if httpCtx.TargetResp == nil {
			t.Error("TargetResp is nil")
			return
		}
	}()

	// event waiting test
	for {
		receivedEvents := httpCtx.Writer.(*StreamRecorder).receivedBuf
		select {
		case event := <-eventCh:
			logger.Info("Received event: %s", strings.ReplaceAll(event, "\n", "\\n"))
		case <-done:
			assert.Equal(t, 3, len(receivedEvents), "Should receive 3 events")
			return
		case <-time.After(5 * time.Second):
			t.Fatal("Test timeout")
			return
		}
	}

}

// mock recorder
type StreamRecorder struct {
	http.ResponseWriter
	http.Flusher
	receivedBuf []string
	headers     http.Header
	status      int
}

func NewStreamRecorder() *StreamRecorder {
	return &StreamRecorder{
		receivedBuf: make([]string, 0),
		headers:     make(http.Header),
	}
}

func (r *StreamRecorder) Header() http.Header {
	return r.headers
}

func (r *StreamRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *StreamRecorder) Write(data []byte) (int, error) {
	eventCh <- string(data)
	r.receivedBuf = append(r.receivedBuf, string(data))
	return len(data), nil
}

func NewTestServerWithURL(URL string, handler http.Handler) (*httptest.Server, error) {
	ts := httptest.NewUnstartedServer(handler)
	if URL != "" {
		l, err := net.Listen("tcp", URL)
		if err != nil {
			return nil, err
		}
		err = ts.Listener.Close()
		if err != nil {
			return nil, err
		}
		ts.Listener = l
	}
	ts.Start()
	return ts, nil
}
