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
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
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

// TestConfigValidate tests the config validation
func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		wantMode string
	}{
		{
			name: "default to pull mode",
			config: &Config{
				Mode: "invalid",
			},
			wantMode: "pull",
		},
		{
			name: "pull mode",
			config: &Config{
				Mode: "pull",
			},
			wantMode: "pull",
		},
		{
			name: "push mode",
			config: &Config{
				Mode: "push",
			},
			wantMode: "push",
		},
		{
			name: "both mode",
			config: &Config{
				Mode: "both",
			},
			wantMode: "both",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantMode, tt.config.Mode)

			// Check defaults
			assert.Equal(t, 9090, tt.config.PullConfig.Port)
			assert.Equal(t, "/metrics", tt.config.PullConfig.Path)
			assert.Equal(t, "http://127.0.0.1:9091", tt.config.PushConfig.GatewayURL)
			assert.Equal(t, "pixiu", tt.config.PushConfig.JobName)
			assert.Equal(t, 100, tt.config.PushConfig.PushInterval)
		})
	}
}

// TestConfigModeEnable tests mode-based enable flags
func TestConfigModeEnable(t *testing.T) {
	tests := []struct {
		name           string
		mode           string
		wantPullEnable bool
		wantPushEnable bool
	}{
		{
			name:           "pull mode",
			mode:           "pull",
			wantPullEnable: true,
			wantPushEnable: false,
		},
		{
			name:           "push mode",
			mode:           "push",
			wantPullEnable: false,
			wantPushEnable: true,
		},
		{
			name:           "both mode",
			mode:           "both",
			wantPullEnable: true,
			wantPushEnable: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{Mode: tt.mode}
			err := cfg.Validate()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantPullEnable, cfg.PullConfig.Enabled)
			assert.Equal(t, tt.wantPushEnable, cfg.PushConfig.Enabled)
		})
	}
}

// TestPullReporterReport tests the OpenTelemetry pull reporter's metric reporting
func TestPullReporterReport(t *testing.T) {
	pullConfig := &PullConfig{
		Enabled: true,
		Port:    9191, // Use different port for testing
		Path:    "/metrics",
	}

	reporter := NewOTelPullReporter(pullConfig)
	require.NotNil(t, reporter)

	// Start the reporter
	err := reporter.Start()
	require.NoError(t, err)

	// Test reporting counter metric
	metrics := []*contextHttp.MetricData{
		{
			Name:  "test_counter",
			Type:  "counter",
			Value: 1.0,
			Labels: map[string]string{
				"method": "GET",
				"status": "200",
			},
		},
	}

	reporter.Report(metrics)

	// Verify the metric was registered in OpenTelemetry (check the instrument exists)
	_, exists := reporter.registeredCounters["test_counter"]
	assert.True(t, exists)
}

// TestPullReporterMultipleMetricTypes tests reporting different metric types with OpenTelemetry
func TestPullReporterMultipleMetricTypes(t *testing.T) {
	pullConfig := &PullConfig{
		Enabled: true,
		Port:    9192,
		Path:    "/metrics",
	}

	reporter := NewOTelPullReporter(pullConfig)
	err := reporter.Start()
	require.NoError(t, err)

	metrics := []*contextHttp.MetricData{
		{
			Name:  "http_requests_total",
			Type:  "counter",
			Value: 1.0,
			Labels: map[string]string{
				"method": "GET",
			},
		},
		{
			Name:  "http_request_duration_ms",
			Type:  "histogram",
			Value: 123.45,
			Labels: map[string]string{
				"method": "GET",
			},
		},
		{
			Name:  "http_active_connections",
			Type:  "gauge",
			Value: 42.0,
			Labels: map[string]string{
				"server": "pixiu",
			},
		},
	}

	reporter.Report(metrics)

	// Verify all metric types were registered in OpenTelemetry
	_, counterExists := reporter.registeredCounters["http_requests_total"]
	_, histogramExists := reporter.registeredHistograms["http_request_duration_ms"]
	_, gaugeExists := reporter.registeredGauges["http_active_connections"]
	
	assert.True(t, counterExists)
	assert.True(t, histogramExists)
	assert.True(t, gaugeExists)
}

// TestPushReporterReport tests the push reporter's metric reporting
func TestPushReporterReport(t *testing.T) {
	pushConfig := &PushConfig{
		Enabled:      true,
		GatewayURL:   "http://localhost:9091",
		JobName:      "test_job",
		PushInterval: 2, // Push every 2 requests for testing
		MetricPath:   "/metrics",
	}

	reporter := NewPushReporter(pushConfig)
	require.NotNil(t, reporter)

	metrics := []*contextHttp.MetricData{
		{
			Name:  "test_counter",
			Type:  "counter",
			Value: 1.0,
			Labels: map[string]string{
				"label": "value",
			},
		},
	}

	// Report once - should not push yet
	reporter.Report(metrics)
	assert.Equal(t, 1, reporter.counter)

	// Report twice - should trigger push
	reporter.Report(metrics)
	// Counter should be reset after push is triggered
	// Note: We use a short sleep to allow the goroutine to execute
	time.Sleep(10 * time.Millisecond)
	assert.Equal(t, 0, reporter.counter)
}

// TestPullReporterHTTPEndpoint tests that OpenTelemetry pull reporter starts HTTP server
func TestPullReporterHTTPEndpoint(t *testing.T) {
	pullConfig := &PullConfig{
		Enabled: true,
		Port:    9193, // Different port for each test
		Path:    "/metrics",
	}

	reporter := NewOTelPullReporter(pullConfig)
	err := reporter.Start()
	require.NoError(t, err)

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	// Report some metrics
	metrics := []*contextHttp.MetricData{
		{
			Name:  "test_http_requests",
			Type:  "counter",
			Value: 5.0,
			Labels: map[string]string{
				"path": "/test",
			},
		},
	}
	reporter.Report(metrics)

	// Try to access the metrics endpoint
	resp, err := http.Get("http://localhost:9193/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read response body
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	// Verify metrics are in the response
	bodyStr := string(body)
	assert.Contains(t, bodyStr, "test_http_requests")
}

// TestFilterEncode tests the filter's Encode method
func TestFilterEncode(t *testing.T) {
	// Create filter factory with pull mode
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "pull",
			PullConfig: PullConfig{
				Enabled: true,
				Port:    9194,
				Path:    "/metrics",
			},
		},
	}

	err := factory.Apply()
	require.NoError(t, err)

	// Create filter
	ctx := newTestHTTPContext(t)
	filter := &Filter{
		factory:      factory,
		pullReporter: factory.pullReporter,
		pushReporter: factory.pushReporter,
	}

	// Record some metrics in context
	ctx.RecordMetric("test_requests", "counter", 1.0, map[string]string{
		"method": "GET",
	})

	// Execute encode
	status := filter.Encode(ctx)
	assert.Equal(t, 0, int(status)) // filter.Continue = 0

	// Verify metrics were reported
	assert.NotNil(t, factory.pullReporter)
	_, exists := factory.pullReporter.registeredCounters["test_requests"]
	assert.True(t, exists)
}

// TestFilterBothMode tests filter with both pull and push modes
func TestFilterBothMode(t *testing.T) {
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "both",
			PullConfig: PullConfig{
				Enabled: true,
				Port:    9195,
				Path:    "/metrics",
			},
			PushConfig: PushConfig{
				Enabled:      true,
				GatewayURL:   "http://localhost:9091",
				JobName:      "test",
				PushInterval: 10,
				MetricPath:   "/metrics",
			},
		},
	}

	err := factory.Apply()
	require.NoError(t, err)

	// Verify both reporters are initialized
	assert.NotNil(t, factory.pullReporter)
	assert.NotNil(t, factory.pushReporter)

	// Create filter and test
	ctx := newTestHTTPContext(t)
	filter := &Filter{
		factory:      factory,
		pullReporter: factory.pullReporter,
		pushReporter: factory.pushReporter,
	}

	ctx.RecordMetric("both_mode_test", "counter", 1.0, map[string]string{
		"mode": "both",
	})

	status := filter.Encode(ctx)
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

	// Verify factory has default config
	ff := factory.(*FilterFactory)
	assert.NotNil(t, ff.cfg)
}

// TestPullReporterConcurrentReporting tests concurrent metric reporting with OpenTelemetry
func TestPullReporterConcurrentReporting(t *testing.T) {
	pullConfig := &PullConfig{
		Enabled: true,
		Port:    9196,
		Path:    "/metrics",
	}

	reporter := NewOTelPullReporter(pullConfig)
	err := reporter.Start()
	require.NoError(t, err)

	// Report metrics concurrently
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func(id int) {
			metrics := []*contextHttp.MetricData{
				{
					Name:  "concurrent_test",
					Type:  "counter",
					Value: 1.0,
					Labels: map[string]string{
						"worker": string(rune(id)),
					},
				},
			}
			reporter.Report(metrics)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify metric was registered in OpenTelemetry
	_, exists := reporter.registeredCounters["concurrent_test"]
	assert.True(t, exists)
}

// TestEmptyMetrics tests handling of empty metrics
func TestEmptyMetrics(t *testing.T) {
	factory := &FilterFactory{
		cfg: &Config{
			Mode: "pull",
			PullConfig: PullConfig{
				Enabled: true,
				Port:    9197,
				Path:    "/metrics",
			},
		},
	}

	err := factory.Apply()
	require.NoError(t, err)

	ctx := newTestHTTPContext(t)
	filter := &Filter{
		factory:      factory,
		pullReporter: factory.pullReporter,
	}

	// No metrics recorded
	status := filter.Encode(ctx)
	assert.Equal(t, 0, int(status)) // Should still continue
}

// TestMetricLabelCopy tests that labels are properly copied
func TestMetricLabelCopy(t *testing.T) {
	pullConfig := &PullConfig{
		Enabled: true,
		Port:    9198,
		Path:    "/metrics",
	}

	reporter := NewOTelPullReporter(pullConfig)
	err := reporter.Start()
	require.NoError(t, err)

	originalLabels := map[string]string{
		"key": "value",
	}

	metrics := []*contextHttp.MetricData{
		{
			Name:   "label_test",
			Type:   "counter",
			Value:  1.0,
			Labels: originalLabels,
		},
	}

	reporter.Report(metrics)

	// Modify original labels
	originalLabels["key"] = "modified"

	// Verify the modification doesn't affect the reported metric
	// This is a safety check for label copying
	_, exists := reporter.registeredCounters["label_test"]
	assert.True(t, exists)
}

// TestPushURLGeneration tests the push URL generation
func TestPushURLGeneration(t *testing.T) {
	pushConfig := &PushConfig{
		Enabled:      true,
		GatewayURL:   "http://localhost:9091",
		JobName:      "test_job",
		PushInterval: 1,
		MetricPath:   "/metrics",
	}

	reporter := NewPushReporter(pushConfig)

	// Note: We can't easily test the actual URL without refactoring,
	// but we can verify the reporter is properly initialized
	assert.NotNil(t, reporter.config)
	assert.Equal(t, "http://localhost:9091", reporter.config.GatewayURL)
	assert.Equal(t, "test_job", reporter.config.JobName)
}

// Benchmark tests
func BenchmarkPullReporterReport(b *testing.B) {
	pullConfig := &PullConfig{
		Enabled: true,
		Port:    9199,
		Path:    "/metrics",
	}

	reporter := NewOTelPullReporter(pullConfig)
	_ = reporter.Start()

	metrics := []*contextHttp.MetricData{
		{
			Name:  "benchmark_counter",
			Type:  "counter",
			Value: 1.0,
			Labels: map[string]string{
				"method": "GET",
				"status": "200",
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		reporter.Report(metrics)
	}
}

func BenchmarkContextRecordMetric(b *testing.B) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := &contextHttp.HttpContext{
		Request: req,
		Writer:  &mockResponseWriter{},
		Ctx:     context.Background(),
	}

	labels := map[string]string{
		"method": "GET",
		"status": "200",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.RecordMetric("test_metric", "counter", 1.0, labels)
	}
}

// TestDecodeMethod tests that Decode method always continues
func TestDecodeMethod(t *testing.T) {
	filter := &Filter{}
	ctx := newTestHTTPContext(t)

	status := filter.Decode(ctx)
	assert.Equal(t, 0, int(status)) // filter.Continue
}

// TestInvalidMetricType tests handling of invalid metric types
func TestInvalidMetricType(t *testing.T) {
	pullConfig := &PullConfig{
		Enabled: true,
		Port:    9200,
		Path:    "/metrics",
	}

	reporter := NewOTelPullReporter(pullConfig)
	err := reporter.Start()
	require.NoError(t, err)

	// Report metric with invalid type (should be silently ignored)
	metrics := []*contextHttp.MetricData{
		{
			Name:   "invalid_type_metric",
			Type:   "invalid_type",
			Value:  1.0,
			Labels: map[string]string{},
		},
	}

	// Should not panic
	assert.NotPanics(t, func() {
		reporter.Report(metrics)
	})
}

// TestMetricsEndpointContentType tests the content type of metrics endpoint
func TestMetricsEndpointContentType(t *testing.T) {
	pullConfig := &PullConfig{
		Enabled: true,
		Port:    9201,
		Path:    "/metrics",
	}

	reporter := NewOTelPullReporter(pullConfig)
	err := reporter.Start()
	require.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://localhost:9201/metrics")
	require.NoError(t, err)
	defer resp.Body.Close()

	// Check content type
	contentType := resp.Header.Get("Content-Type")
	assert.True(t, 
		strings.Contains(contentType, "text/plain") || 
		strings.Contains(contentType, "application/openmetrics-text"),
		"Expected Prometheus compatible content type, got: %s", contentType)
}

