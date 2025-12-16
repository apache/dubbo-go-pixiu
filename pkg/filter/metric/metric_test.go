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

package metric

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"

	"go.opentelemetry.io/otel/sdk/metric"
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
		wantMode  string
	}{
		{
			name: "empty mode defaults to push",
			config: &Config{
				Mode: "",
			},
			wantError: false,
			wantMode:  "push",
		},
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
			wantMode:  "pull",
		},
		{
			name: "valid push mode",
			config: &Config{
				Mode: "push",
				Push: PushConfig{
					GatewayURL:   "http://localhost:9091",
					JobName:      "pixiu",
					PushInterval: 100,
					MetricPath:   "/metrics",
				},
			},
			wantError: false,
			wantMode:  "push",
		},
		{
			name: "push mode with empty fields applies defaults",
			config: &Config{
				Mode: "push",
				Push: PushConfig{
					GatewayURL:   "",
					JobName:      "",
					PushInterval: 0,
					MetricPath:   "",
				},
			},
			wantError: false,
			wantMode:  "push",
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
				if tt.wantMode != "" {
					assert.Equal(t, tt.wantMode, factory.cfg.Mode, "Mode should be set to %s", tt.wantMode)
				}
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
			Push: PushConfig{
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

	// Verify filter was created - push mode now has both decode and encode filters
	require.Len(t, chain.decodeFilters, 1)
	require.Len(t, chain.encodeFilters, 1, "Push mode should have both decode and encode filters")
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
			Push: PushConfig{
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

	// Push mode now has both decode and encode filters
	require.Len(t, chain.decodeFilters, 1)
	require.Len(t, chain.encodeFilters, 1)

	// Execute decode (only records start time)
	status := chain.decodeFilters[0].Decode(ctx)
	assert.Equal(t, 0, int(status))

	// Execute encode (reports metrics)
	status = chain.encodeFilters[0].Encode(ctx)
	assert.Equal(t, 0, int(status))
}

// TestPluginKind tests the plugin kind
func TestPluginKind(t *testing.T) {
	plugin := &Plugin{}
	assert.Equal(t, "dgp.filter.http.metric", plugin.Kind())
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
		name string
		mode string
	}{
		{"pull mode with nil instruments in encode", "pull"},
		{"push mode with nil collector in encode", "push"},
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

			// Decode should always succeed (just records time)
			decodeStatus := int(filter.Decode(ctx))
			assert.Equal(t, 0, decodeStatus) // filter.Continue = 0

			// Encode should fail when reporter is not initialized
			encodeStatus := int(filter.Encode(ctx))
			assert.Equal(t, 1, encodeStatus) // filter.Stop = 1

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

// TestMetricReporterPullMode tests pull mode with OpenTelemetry.
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
func TestMetricReporterPushMode(t *testing.T) {
	// Create factory with push mode
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "push",
			Push: PushConfig{
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

	// Push mode now has both decode and encode filters
	require.Len(t, chain.decodeFilters, 1)
	require.Len(t, chain.encodeFilters, 1, "Push mode should have both decode and encode filters")

	// Simulate multiple requests (to trigger push)
	for i := 0; i < 15; i++ {
		// Record some context metrics before decode
		ctx.RecordMetric("api_requests_total", "counter", 1.0, map[string]string{
			"api": "health",
		})

		// Execute decode (only records start time)
		decodeStatus := chain.decodeFilters[0].Decode(ctx)
		assert.Equal(t, 0, int(decodeStatus))

		// Execute encode (reports metrics in push mode)
		encodeStatus := chain.encodeFilters[0].Encode(ctx)
		assert.Equal(t, 0, int(encodeStatus))

		// Clear metrics for next iteration
		ctx.ClearMetrics()
	}

	t.Log("Push mode metric reporter test finished successfully")
}

// TestOTelInstrumentNoErrorOnDuplicateName tests that OpenTelemetry SDK
// allows creating instruments with the same name without errors.
// Although it returns different wrapper objects, it doesn't cause duplicate registration issues.
func TestOTelInstrumentNoErrorOnDuplicateName(t *testing.T) {
	// Initialize OTel instruments first
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

	// Get the meter
	meter := global.MeterProvider().Meter("pixiu")

	// Create the same counter multiple times with the same name
	// This should NOT cause errors even though it's the same name
	counter1, err1 := meter.SyncInt64().Counter("test_duplicate_counter",
		instrument.WithDescription("First call"))
	require.NoError(t, err1)
	require.NotNil(t, counter1)

	counter2, err2 := meter.SyncInt64().Counter("test_duplicate_counter",
		instrument.WithDescription("Second call"))
	require.NoError(t, err2)
	require.NotNil(t, counter2)

	counter3, err3 := meter.SyncInt64().Counter("test_duplicate_counter",
		instrument.WithDescription("Third call"))
	require.NoError(t, err3)
	require.NotNil(t, counter3)

	// Test with histogram
	hist1, err4 := meter.SyncFloat64().Histogram("test_duplicate_histogram",
		instrument.WithDescription("First histogram"))
	require.NoError(t, err4)
	require.NotNil(t, hist1)

	hist2, err5 := meter.SyncFloat64().Histogram("test_duplicate_histogram",
		instrument.WithDescription("Second histogram"))
	require.NoError(t, err5)
	require.NotNil(t, hist2)

	// Test with gauge (UpDownCounter)
	gauge1, err6 := meter.SyncInt64().UpDownCounter("test_duplicate_gauge",
		instrument.WithDescription("First gauge"))
	require.NoError(t, err6)
	require.NotNil(t, gauge1)

	gauge2, err7 := meter.SyncInt64().UpDownCounter("test_duplicate_gauge",
		instrument.WithDescription("Second gauge"))
	require.NoError(t, err7)
	require.NotNil(t, gauge2)

	// All instruments can be used without errors
	counter1.Add(ctx.Ctx, 1)
	counter2.Add(ctx.Ctx, 1)
	counter3.Add(ctx.Ctx, 1)

	hist1.Record(ctx.Ctx, 10.5)
	hist2.Record(ctx.Ctx, 20.5)

	gauge1.Add(ctx.Ctx, 1)
	gauge2.Add(ctx.Ctx, -1)

	t.Log("OpenTelemetry SDK allows same metric name without errors - no duplicate registration issues")
}

// TestDynamicMetricsMultipleRequests tests that dynamic metrics from context
// can be reported multiple times without issues (simulates real scenario).
func TestDynamicMetricsMultipleRequests(t *testing.T) {
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "pull",
		},
	}

	err := factory.Apply()
	require.NoError(t, err)

	// Simulate 100 requests with the same custom metric
	for i := 0; i < 100; i++ {
		ctx := newTestHTTPContext(t)

		// Record the same metric name in every request
		ctx.RecordMetric("api_requests", "counter", 1.0, map[string]string{
			"endpoint": "/test",
			"method":   "GET",
		})
		ctx.RecordMetric("api_latency", "histogram", float64(i*10), map[string]string{
			"endpoint": "/test",
		})

		// Create a new chain for each request
		chain := &mockFilterChain{}
		err = factory.PrepareFilterChain(ctx, chain)
		require.NoError(t, err)

		// Execute decode and encode
		require.Len(t, chain.decodeFilters, 1)
		require.Len(t, chain.encodeFilters, 1)

		decodeStatus := chain.decodeFilters[0].Decode(ctx)
		assert.Equal(t, 0, int(decodeStatus))

		encodeStatus := chain.encodeFilters[0].Encode(ctx)
		assert.Equal(t, 0, int(encodeStatus), "Request #%d: should handle repeated metric names without error", i)
	}

	t.Log("Successfully processed 100 requests with same metric names - no duplicate registration issues")
}

// TestSDKProviderRejectsRepeatedRegistration tests that when using SDK MeterProvider directly,
// repeated registration of the same metric name WILL cause an error.
// This is different from using global.MeterProvider().
func TestSDKProviderRejectsRepeatedRegistration(t *testing.T) {
	reader := metric.NewManualReader()
	provider := metric.NewMeterProvider(metric.WithReader(reader))
	meter := provider.Meter("pixiu")

	// First registration - should succeed
	counter1, err1 := meter.SyncInt64().Counter("test_counter",
		instrument.WithDescription("First"))
	require.NoError(t, err1)
	require.NotNil(t, counter1)

	// Second registration with SAME NAME - should FAIL with SDK provider
	_, err2 := meter.SyncInt64().Counter("test_counter",
		instrument.WithDescription("Second"))

	// SDK MeterProvider DOES reject duplicate registration
	assert.Error(t, err2, "SDK MeterProvider should reject duplicate instrument registration")
	assert.Contains(t, err2.Error(), "instrument already registered",
		"Error should indicate duplicate registration")

	t.Log("✓ Confirmed: SDK MeterProvider rejects duplicate instrument registration")
}

// TestGlobalProviderHandlesRepeatedCalls tests that when using global.MeterProvider(),
// which is what the actual code uses, repeated calls do NOT cause errors.
func TestGlobalProviderHandlesRepeatedCalls(t *testing.T) {
	// Use global meter provider (default noop or whatever is set)
	meter := global.MeterProvider().Meter("pixiu")

	// Create the same counter multiple times - this is what happens in actual code
	for i := 0; i < 10; i++ {
		counter, err := meter.SyncInt64().Counter("global_test_counter",
			instrument.WithDescription(fmt.Sprintf("Iteration %d", i)))

		// With global provider, this should NOT error
		assert.NoError(t, err, "global.MeterProvider() should handle repeated instrument creation")
		assert.NotNil(t, counter)

		// Use the counter
		counter.Add(context.Background(), int64(i+1))
	}

	t.Log("✓ Confirmed: global.MeterProvider() handles repeated instrument creation without errors")
	t.Log("✓ This explains why the actual code (which uses global.MeterProvider) works fine")
}
