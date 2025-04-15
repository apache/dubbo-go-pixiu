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
	"context"
	"fmt"
	"net"
	stdhttp "net/http"
	"net/http/httptest"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	_ "github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/trie"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	streamEventCh = make(chan string, 10)
)

// 测试各种普通流式HTTP响应类型
func TestStreamableHTTPResponse(t *testing.T) {
	// 定义要测试的内容类型
	contentTypes := []string{
		"text/plain",
		"application/json",
		"application/octet-stream",
		"application/x-ndjson",
	}

	for _, contentType := range contentTypes {
		t.Run(fmt.Sprintf("ContentType_%s", contentType), func(t *testing.T) {
			testStreamableResponse(t, contentType)
		})
	}
}

func testStreamableResponse(t *testing.T, contentType string) {
	hcmc := model.HttpConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{
			RouteTrie: trie.NewTrieWithDefault("GET/api/stream", model.RouteAction{
				Cluster: "mock_stream_cluster",
			}),
		},
		HTTPFilters: []*model.HTTPFilter{
			{
				Name: mock.Kind,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Clear any data that may have been left over from the previous test
	for len(streamEventCh) > 0 {
		<-streamEventCh
	}

	// mock server
	upstreamServer, _ := NewTestStreamServerWithURL("localhost:8080", stdhttp.HandlerFunc(func(w stdhttp.ResponseWriter, r *stdhttp.Request) {
		w.Header().Set("Content-Type", contentType)
		flusher := w.(stdhttp.Flusher)

		// Generate appropriate test data based on content type
		var data []byte
		for i := 1; i <= 5; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				time.Sleep(10 * time.Millisecond)

				switch contentType {
				case "application/json":
					data = []byte(fmt.Sprintf(`{"id": %d, "message": "test chunk %d"}\n`, i, i))
				case "application/x-ndjson":
					data = []byte(fmt.Sprintf(`{"id": %d, "message": "test chunk %d"}\n`, i, i))
				case "application/octet-stream":
					data = []byte(fmt.Sprintf("CHUNK-%d", i))
				default: // text/plain
					data = []byte(fmt.Sprintf("Chunk %d\n", i))
				}

				_, _ = w.Write(data)
				flusher.Flush()
				logger.Info("Upstream sent chunk ", i)
			}
		}
	}))
	defer upstreamServer.Close()

	req := httptest.NewRequest("GET", "http://localhost:8080/api/stream", nil).WithContext(ctx)
	done := make(chan struct{})

	httpCtx := &contexthttp.HttpContext{
		Request: req,
		Writer:  NewStreamHTTPRecorder(),
		Ctx:     ctx,
	}

	go func() {
		defer close(done)

		hcm := CreateHttpConnectionManager(&hcmc)

		if err := hcm.Handle(httpCtx); err != nil {
			t.Errorf("Handle failed: %v", err)
		}

		// verify that targetResp exists
		if httpCtx.TargetResp == nil {
			t.Error("TargetResp is nil")
			return
		}
	}()

	// collect and validate responses
	receivedChunks := 0
	for {
		receivedEvents := httpCtx.Writer.(*StreamHTTPRecorder).receivedBuf
		select {
		case event := <-streamEventCh:
			logger.Info("Received chunk: %s", event)
			receivedChunks++
		case <-done:
			assert.Equal(t, 5, len(receivedEvents), "Should receive 5 chunks")
			return
		case <-time.After(5 * time.Second):
			t.Fatal("Test timeout")
			return
		}
	}
}

// StreamHTTPRecorder
type StreamHTTPRecorder struct {
	stdhttp.ResponseWriter
	receivedBuf []string
	headers     stdhttp.Header
	status      int
	flushCount  int
}

func NewStreamHTTPRecorder() *StreamHTTPRecorder {
	return &StreamHTTPRecorder{
		receivedBuf: make([]string, 0),
		headers:     make(stdhttp.Header),
		flushCount:  0,
	}
}

func (r *StreamHTTPRecorder) Header() stdhttp.Header {
	return r.headers
}

func (r *StreamHTTPRecorder) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *StreamHTTPRecorder) Write(data []byte) (int, error) {
	streamEventCh <- string(data)
	r.receivedBuf = append(r.receivedBuf, string(data))
	return len(data), nil
}

func (r *StreamHTTPRecorder) Flush() {
	r.flushCount++
}

// NewTestStreamServerWithURL
func NewTestStreamServerWithURL(URL string, handler stdhttp.Handler) (*httptest.Server, error) {
	ts := httptest.NewUnstartedServer(handler)
	if URL != "" {
		l, err := net.Listen("tcp", URL)
		if err != nil {
			return nil, err
		}
		ts.Listener.Close()
		ts.Listener = l
	}
	ts.Start()
	return ts, nil
}
