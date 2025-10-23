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

package metricreporter

import (
	"context"
	"net/http"
	"net/url"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

// mockResponseWriter is a test implementation of http.ResponseWriter
type mockResponseWriter struct {
	header http.Header
	body   []byte
	status int
}

func (w *mockResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *mockResponseWriter) Write(b []byte) (int, error) {
	w.body = append(w.body, b...)
	return len(b), nil
}

func (w *mockResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
}

// newTestHTTPContext creates a test HTTP context
func newTestHTTPContext(t *testing.T) *contextHttp.HttpContext {
	req, err := http.NewRequest("GET", "http://example.com/test", nil)
	require.NoError(t, err)

	return &contextHttp.HttpContext{
		Request: req,
		Writer:  &mockResponseWriter{},
		Ctx:     context.Background(),
	}
}

// mockFilterChain for testing
type mockFilterChain struct {
	decodeFilters []filter.HttpDecodeFilter
	encodeFilters []filter.HttpEncodeFilter
}

func (m *mockFilterChain) AppendDecodeFilters(f ...filter.HttpDecodeFilter) {
	m.decodeFilters = append(m.decodeFilters, f...)
}

func (m *mockFilterChain) AppendEncodeFilters(f ...filter.HttpEncodeFilter) {
	m.encodeFilters = append(m.encodeFilters, f...)
}

func (m *mockFilterChain) OnDecode(ctx *contextHttp.HttpContext) {
	// Not used in tests
}

func (m *mockFilterChain) OnEncode(ctx *contextHttp.HttpContext) {
	// Not used in tests
}

// TestConfigValidate tests the config validation
func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name      string
		config    *Config
		wantError bool
	}{
		{
			name: "invalid mode",
			config: &Config{
				Mode: "invalid",
			},
			wantError: true,
		},
		{
			name: "valid pull mode",
			config: &Config{
				Mode: "pull",
			},
			wantError: false,
		},
		{
			name: "valid push mode",
			config: &Config{
				Mode: "push",
				PushConfig: PushConfig{
					GatewayURL:   "http://localhost:9091",
					JobName:      "pixiu",
					PushInterval: 100,
					MetricPath:   "/metrics",
				},
			},
			wantError: false,
		},
		{
			name: "push mode with empty gateway_url",
			config: &Config{
				Mode: "push",
				PushConfig: PushConfig{
					GatewayURL:   "",
					JobName:      "pixiu",
					PushInterval: 100,
					MetricPath:   "/metrics",
				},
			},
			wantError: true,
		},
		{
			name: "push mode with invalid interval",
			config: &Config{
				Mode: "push",
				PushConfig: PushConfig{
					GatewayURL:   "http://localhost:9091",
					JobName:      "pixiu",
					PushInterval: 0,
					MetricPath:   "/metrics",
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory := &FilterFactory{cfg: tt.config}
			err := factory.Apply()

			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestPullModeInitialization tests pull mode initialization
func TestPullModeInitialization(t *testing.T) {
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "pull",
		},
	}

	err := factory.Apply()
	require.NoError(t, err)

	// Initialization happens in PrepareFilterChain, so we need to test that
	ctx := newTestHTTPContext(t)
	chain := &mockFilterChain{}
	err = factory.PrepareFilterChain(ctx, chain)
	require.NoError(t, err)

	// Verify filter was created
	require.Len(t, chain.decodeFilters, 1)
}

// TestPushModeInitialization tests push mode initialization
func TestPushModeInitialization(t *testing.T) {
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "push",
			PushConfig: PushConfig{
				GatewayURL:   "http://localhost:9091",
				JobName:      "test_job",
				PushInterval: 100,
				MetricPath:   "/metrics",
			},
		},
	}

	err := factory.Apply()
	require.NoError(t, err)

	// Initialization happens in PrepareFilterChain
	ctx := newTestHTTPContext(t)
	chain := &mockFilterChain{}
	err = factory.PrepareFilterChain(ctx, chain)
	require.NoError(t, err)

	// Verify filter was created (push mode only has decode filter)
	require.Len(t, chain.decodeFilters, 1)
	require.Len(t, chain.encodeFilters, 0)
}

// TestFilterWithPullMode tests filter encode with pull mode
func TestFilterWithPullMode(t *testing.T) {
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "pull",
		},
	}

	err := factory.Apply()
	require.NoError(t, err)

	ctx := newTestHTTPContext(t)
	chain := &mockFilterChain{}
	err = factory.PrepareFilterChain(ctx, chain)
	require.NoError(t, err)

	// Record metrics in context
	ctx.RecordMetric("custom_metric", "counter", 1.0, map[string]string{
		"key": "value",
	})

	// Pull mode should have both decode and encode filters
	require.Len(t, chain.decodeFilters, 1)
	require.Len(t, chain.encodeFilters, 1)

	// Execute decode (records start time)
	chain.decodeFilters[0].Decode(ctx)

	// Execute encode (reports metrics)
	status := chain.encodeFilters[0].Encode(ctx)
	assert.Equal(t, 0, int(status))
}

// TestFilterWithPushMode tests filter encode with push mode
func TestFilterWithPushMode(t *testing.T) {
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "push",
			PushConfig: PushConfig{
				GatewayURL:   "http://localhost:9091",
				JobName:      "test",
				PushInterval: 100,
				MetricPath:   "/metrics",
			},
		},
	}

	err := factory.Apply()
	require.NoError(t, err)

	ctx := newTestHTTPContext(t)
	chain := &mockFilterChain{}
	err = factory.PrepareFilterChain(ctx, chain)
	require.NoError(t, err)

	// Record metrics in context
	ctx.RecordMetric("custom_metric", "counter", 1.0, nil)

	// Push mode only has decode filter
	require.Len(t, chain.decodeFilters, 1)
	require.Len(t, chain.encodeFilters, 0)

	// Execute decode (reports metrics immediately)
	status := chain.decodeFilters[0].Decode(ctx)
	assert.Equal(t, 0, int(status))
}

// TestPluginKind tests the plugin kind
func TestPluginKind(t *testing.T) {
	plugin := &Plugin{}
	assert.Equal(t, "dgp.filter.http.metricreporter", plugin.Kind())
}

// TestCreateFilterFactory tests creating a filter factory
func TestCreateFilterFactory(t *testing.T) {
	plugin := &Plugin{}
	factory, err := plugin.CreateFilterFactory()
	require.NoError(t, err)
	require.NotNil(t, factory)

	ff := factory.(*FilterFactory)
	assert.NotNil(t, ff.cfg)
}

// TestFilterWithUninitializedReporter tests that filter stops when reporter is not initialized
func TestFilterWithUninitializedReporter(t *testing.T) {
	tests := []struct {
		name  string
		mode  string
		phase string // "decode" or "encode"
	}{
		{"pull mode with nil instruments in encode", "pull", "encode"},
		{"push mode with nil collector in decode", "push", "decode"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create filter with config but nil reporter
			filter := &Filter{
				cfg:             &Config{Mode: tt.mode},
				otelInstruments: nil,
				promCollector:   nil,
			}

			ctx := newTestHTTPContext(t)
			ctx.RecordMetric("test", "counter", 1.0, nil)

			var status int
			if tt.phase == "decode" {
				status = int(filter.Decode(ctx))
			} else {
				status = int(filter.Encode(ctx))
			}

			// Should stop when reporter is not initialized
			assert.Equal(t, 1, status) // filter.Stop = 1

			// Should have sent local reply
			assert.True(t, ctx.LocalReply())
			assert.Equal(t, 500, ctx.GetStatusCode())
		})
	}
}

// TestDecodeMethod tests that Decode method records start time
func TestDecodeMethod(t *testing.T) {
	filter := &Filter{
		cfg: &Config{Mode: "pull"},
	}
	ctx := newTestHTTPContext(t)

	status := filter.Decode(ctx)
	assert.Equal(t, 0, int(status)) // filter.Continue

	// Verify start time was recorded
	assert.False(t, filter.start.IsZero())
}

// ============================================================================
// Integration Tests - Prove new filter can replace old filters
// ============================================================================

// TestMetricReporterPullMode tests pull mode with OpenTelemetry.
// This test proves the new filter can replace dgp.filter.http.metric
// Inspired by TestMetric in pkg/filter/metric/metric_test.go
func TestMetricReporterPullMode(t *testing.T) {
	// Create factory with pull mode
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "pull",
		},
	}

	// Validate configuration
	err := factory.Apply()
	require.NoError(t, err)

	// Create HTTP request
	req, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/mock/test?name=tc", nil)
	require.NoError(t, err)

	ctx := &contextHttp.HttpContext{
		Request: req,
		Writer:  &mockResponseWriter{},
		Ctx:     context.Background(),
	}

	// Prepare filter chain
	chain := &mockFilterChain{}
	err = factory.PrepareFilterChain(ctx, chain)
	require.NoError(t, err)

	// Get filters
	// Pull mode should have both decode and encode filters
	require.Len(t, chain.decodeFilters, 1)
	require.Len(t, chain.encodeFilters, 1)

	// Execute decode (records start time)
	decodeStatus := chain.decodeFilters[0].Decode(ctx)
	assert.Equal(t, 0, int(decodeStatus))

	// Execute encode (reports metrics)
	encodeStatus := chain.encodeFilters[0].Encode(ctx)
	assert.Equal(t, 0, int(encodeStatus))

	t.Log("Pull mode metric reporter test finished successfully")
}

// TestMetricReporterPushMode tests push mode with Prometheus Push Gateway.
// This test proves the new filter can replace dgp.filter.http.prometheusmetric
// Inspired by TestCounterExporterApiMetric in pkg/filter/prometheus/metric_test.go
func TestMetricReporterPushMode(t *testing.T) {
	// Create factory with push mode
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "push",
			PushConfig: PushConfig{
				GatewayURL:   "http://127.0.0.1:9091",
				JobName:      "pixiu-test",
				PushInterval: 10, // Push every 10 requests for faster testing
				MetricPath:   "/metrics",
			},
		},
	}

	// Validate configuration
	err := factory.Apply()
	require.NoError(t, err)

	// Prepare filter chain
	testURL, _ := url.Parse("http://localhost/_api/health")
	ctx := &contextHttp.HttpContext{
		Request: &http.Request{
			Method: "POST",
			URL:    testURL,
			Host:   "localhost",
		},
		Writer: &mockResponseWriter{},
		Ctx:    context.Background(),
	}

	chain := &mockFilterChain{}
	err = factory.PrepareFilterChain(ctx, chain)
	require.NoError(t, err)

	// Push mode should only have decode filter, no encode filter
	require.Len(t, chain.decodeFilters, 1)
	require.Len(t, chain.encodeFilters, 0, "Push mode should not have encode filter")

	// Simulate multiple requests (to trigger push)
	for i := 0; i < 15; i++ {
		// Record some context metrics before decode
		ctx.RecordMetric("api_requests_total", "counter", 1.0, map[string]string{
			"api": "health",
		})

		// Execute decode (reports metrics immediately in push mode)
		decodeStatus := chain.decodeFilters[0].Decode(ctx)
		assert.Equal(t, 0, int(decodeStatus))

		// Clear metrics for next iteration
		ctx.ClearMetrics()
	}

	t.Log("Push mode metric reporter test finished successfully")
}
