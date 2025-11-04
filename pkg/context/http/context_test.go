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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
)

// mockResponseWriter is a test implementation of http.ResponseWriter
type mockResponseWriter struct {
	header http.Header
}

func (w *mockResponseWriter) Header() http.Header {
	if w.header == nil {
		w.header = make(http.Header)
	}
	return w.header
}

func (w *mockResponseWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (w *mockResponseWriter) WriteHeader(statusCode int) {
}

// newTestHTTPContext creates a mock HttpContext for testing
func newTestHTTPContext(r *http.Request) *HttpContext {
	ctx := &HttpContext{
		Index:   -1,
		Request: r,
		Writer:  &mockResponseWriter{},
		Ctx:     context.Background(),
	}
	ctx.Reset()
	return ctx
}

// TestErrorBuilder tests the ErrorBuilder methods
func TestErrorBuilder(t *testing.T) {
	tests := []struct {
		name           string
		builder        *ErrorBuilder
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "BadRequest",
			builder:        BadRequest,
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "Bad request",
		},
		{
			name:           "NotFound",
			builder:        RouteNotFound,
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "Route not found",
		},
		{
			name:           "InternalError",
			builder:        InternalError,
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    "Internal server error",
		},
		{
			name:           "ServiceUnavailable",
			builder:        ServiceUnavailable,
			expectedStatus: http.StatusServiceUnavailable,
			expectedMsg:    "Service unavailable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test New()
			errResp := tt.builder.New()
			if errResp.Status != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, errResp.Status)
			}
			if errResp.Message != tt.expectedMsg {
				t.Errorf("expected message %q, got %q", tt.expectedMsg, errResp.Message)
			}
			if errResp.Err != nil {
				t.Errorf("expected nil error, got %v", errResp.Err)
			}

			// Test GetStatus()
			if tt.builder.GetStatus() != tt.expectedStatus {
				t.Errorf("GetStatus() = %d, want %d", tt.builder.GetStatus(), tt.expectedStatus)
			}
		})
	}
}

// TestErrorBuilderWithError tests the WithError method
func TestErrorBuilderWithError(t *testing.T) {
	testErr := errors.New("test error")
	wrappedErr := fmt.Errorf("wrapped: %w", testErr)

	tests := []struct {
		name        string
		builder     *ErrorBuilder
		err         error
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "BadRequest with simple error",
			builder:     BadRequest,
			err:         testErr,
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Bad request",
		},
		{
			name:        "InternalError with wrapped error",
			builder:     InternalError,
			err:         wrappedErr,
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "Internal server error",
		},
		{
			name:        "GatewayTimeout with formatted error",
			builder:     GatewayTimeout,
			err:         fmt.Errorf("timeout after 30s: %w", testErr),
			wantStatus:  http.StatusGatewayTimeout,
			wantMessage: "Gateway timeout",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errResp := tt.builder.WithError(tt.err)
			if errResp.Status != tt.wantStatus {
				t.Errorf("Status = %d, want %d", errResp.Status, tt.wantStatus)
			}
			if errResp.Message != tt.wantMessage {
				t.Errorf("Message = %q, want %q", errResp.Message, tt.wantMessage)
			}
			if errResp.Err != tt.err {
				t.Errorf("Err = %v, want %v", errResp.Err, tt.err)
			}
		})
	}
}

// TestErrorResponseToJSON tests JSON serialization
func TestErrorResponseToJSON(t *testing.T) {
	tests := []struct {
		name     string
		errResp  *ErrorResponse
		wantJSON string
	}{
		{
			name:     "without error",
			errResp:  BadRequest.New(),
			wantJSON: `{"status":400,"message":"Bad request"}`,
		},
		{
			name:     "with simple error",
			errResp:  BadRequest.WithError(errors.New("invalid parameter")),
			wantJSON: `{"status":400,"message":"Bad request","error":"invalid parameter"}`,
		},
		{
			name:     "with wrapped error",
			errResp:  InternalError.WithError(fmt.Errorf("failed to process: %w", errors.New("connection refused"))),
			wantJSON: `{"status":500,"message":"Internal server error","error":"failed to process: connection refused"}`,
		},
		{
			name:     "NotFound without error",
			errResp:  RouteNotFound.New(),
			wantJSON: `{"status":404,"message":"Route not found"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJSON := tt.errResp.ToJSON()

			// Compare JSON structure
			var got, want map[string]any
			if err := json.Unmarshal(gotJSON, &got); err != nil {
				t.Fatalf("failed to unmarshal got JSON: %v", err)
			}
			if err := json.Unmarshal([]byte(tt.wantJSON), &want); err != nil {
				t.Fatalf("failed to unmarshal want JSON: %v", err)
			}

			if got["status"] != want["status"] {
				t.Errorf("status = %v, want %v", got["status"], want["status"])
			}
			if got["message"] != want["message"] {
				t.Errorf("message = %v, want %v", got["message"], want["message"])
			}
			if want["error"] != nil && got["error"] != want["error"] {
				t.Errorf("error = %v, want %v", got["error"], want["error"])
			}
		})
	}
}

// TestErrorResponseError tests the Error() method
func TestErrorResponseError(t *testing.T) {
	tests := []struct {
		name    string
		errResp *ErrorResponse
		wantStr string
	}{
		{
			name:    "without error",
			errResp: BadRequest.New(),
			wantStr: "[400] Bad request",
		},
		{
			name:    "with simple error",
			errResp: BadRequest.WithError(errors.New("invalid id")),
			wantStr: "[400] Bad request: invalid id",
		},
		{
			name:    "with wrapped error",
			errResp: InternalError.WithError(fmt.Errorf("database error: %w", errors.New("connection lost"))),
			wantStr: "[500] Internal server error: database error: connection lost",
		},
		{
			name:    "503 with context",
			errResp: ServiceUnavailable.WithError(fmt.Errorf("endpoint not found: %w", errors.New("no healthy hosts"))),
			wantStr: "[503] Service unavailable: endpoint not found: no healthy hosts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.errResp.Error()
			if got != tt.wantStr {
				t.Errorf("Error() = %q, want %q", got, tt.wantStr)
			}
		})
	}
}

// TestAllErrorTypes verifies all predefined error types
func TestAllErrorTypes(t *testing.T) {
	errorTypes := []struct {
		name    string
		builder *ErrorBuilder
		status  int
	}{
		// 404 errors
		{"RouteNotFound", RouteNotFound, http.StatusNotFound},
		{"ServiceNotFound", ServiceNotFound, http.StatusNotFound},
		{"APINotFound", APINotFound, http.StatusNotFound},
		// 400 errors
		{"BadRequest", BadRequest, http.StatusBadRequest},
		// 401 errors
		{"Unauthorized", Unauthorized, http.StatusUnauthorized},
		// 403 errors
		{"Forbidden", Forbidden, http.StatusForbidden},
		// 405 errors
		{"MethodNotAllowed", MethodNotAllowed, http.StatusMethodNotAllowed},
		// 406 errors
		{"NotAcceptable", NotAcceptable, http.StatusNotAcceptable},
		// 429 errors
		{"RateLimited", RateLimited, http.StatusTooManyRequests},
		// 500 errors
		{"InternalError", InternalError, http.StatusInternalServerError},
		{"ConfigurationError", ConfigurationError, http.StatusInternalServerError},
		// 502 errors
		{"BadGateway", BadGateway, http.StatusBadGateway},
		// 503 errors
		{"ServiceUnavailable", ServiceUnavailable, http.StatusServiceUnavailable},
		// 504 errors
		{"GatewayTimeout", GatewayTimeout, http.StatusGatewayTimeout},
	}

	for _, tt := range errorTypes {
		t.Run(tt.name, func(t *testing.T) {
			if tt.builder == nil {
				t.Fatal("error builder is nil")
			}
			if tt.builder.GetStatus() != tt.status {
				t.Errorf("status = %d, want %d", tt.builder.GetStatus(), tt.status)
			}

			// Test that New() works
			errResp := tt.builder.New()
			if errResp.Status != tt.status {
				t.Errorf("New().Status = %d, want %d", errResp.Status, tt.status)
			}

			// Test that WithError() works
			testErr := errors.New("test")
			errResp = tt.builder.WithError(testErr)
			if errResp.Err != testErr {
				t.Errorf("WithError().Err = %v, want %v", errResp.Err, testErr)
			}
		})
	}
}

// TestErrorResponseJSONMarshaling tests edge cases in JSON marshaling
func TestErrorResponseJSONMarshaling(t *testing.T) {
	t.Run("error with special characters", func(t *testing.T) {
		errResp := BadRequest.WithError(errors.New(`error with "quotes" and \backslash`))
		jsonBytes := errResp.ToJSON()

		var result map[string]any
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if result["error"] != `error with "quotes" and \backslash` {
			t.Errorf("error field not properly escaped: %v", result["error"])
		}
	})

	t.Run("nil error", func(t *testing.T) {
		errResp := &ErrorResponse{
			Status:  http.StatusBadRequest,
			Message: "Bad request",
			Err:     nil,
		}
		jsonBytes := errResp.ToJSON()

		var result map[string]any
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if _, hasError := result["error"]; hasError {
			t.Error("error field should be omitted when Err is nil")
		}
	})

	t.Run("WithError(nil) behavior", func(t *testing.T) {
		errResp := BadRequest.WithError(nil)
		jsonBytes := errResp.ToJSON()

		var result map[string]any
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		// Verify no error field when error is nil
		if _, hasError := result["error"]; hasError {
			t.Error("WithError(nil) should not include error field in JSON")
		}

		// Verify basic fields are present
		if result["status"] != float64(http.StatusBadRequest) {
			t.Errorf("status = %v, want %v", result["status"], http.StatusBadRequest)
		}
		if result["message"] != "Bad request" {
			t.Errorf("message = %v, want 'Bad request'", result["message"])
		}
	})

	t.Run("empty message and zero status", func(t *testing.T) {
		errResp := &ErrorResponse{
			Status:  0,
			Message: "",
			Err:     nil,
		}

		// Should not panic
		jsonBytes := errResp.ToJSON()

		var result map[string]any
		if err := json.Unmarshal(jsonBytes, &result); err != nil {
			t.Fatalf("failed to unmarshal JSON: %v", err)
		}

		if result["status"] != float64(0) {
			t.Errorf("status = %v, want 0", result["status"])
		}
		if result["message"] != "" {
			t.Errorf("message = %v, want empty string", result["message"])
		}

		// No error field expected
		if _, hasError := result["error"]; hasError {
			t.Error("error field should be omitted when Err is nil")
		}
	})

	t.Run("Error() with empty message and zero status", func(t *testing.T) {
		errResp := &ErrorResponse{
			Status:  0,
			Message: "",
			Err:     nil,
		}

		// Should not panic
		got := errResp.Error()
		want := "[0] "
		if got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("Error() with empty message but has error", func(t *testing.T) {
		errResp := &ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: "",
			Err:     errors.New("internal error"),
		}

		// Should not panic
		got := errResp.Error()
		want := "[500] : internal error"
		if got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("Error() with zero status and has error", func(t *testing.T) {
		errResp := &ErrorResponse{
			Status:  0,
			Message: "Unknown error",
			Err:     errors.New("something went wrong"),
		}

		// Should not panic
		got := errResp.Error()
		want := "[0] Unknown error: something went wrong"
		if got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})
}

// TestRecordMetric tests recording metrics to context
func TestRecordMetric(t *testing.T) {
	req, err := http.NewRequest("GET", "http://example.com/test", nil)
	if err != nil {
		t.Fatalf("failed to create request: %v", err)
	}

	ctx := newTestHTTPContext(req)

	// Record a metric
	ctx.RecordMetric("test_metric", "counter", 1.0, map[string]string{
		"label1": "value1",
		"label2": "value2",
	})

	// Get all metrics
	metrics := ctx.GetAllMetrics()

	// Verify
	if len(metrics) != 1 {
		t.Errorf("expected 1 metric, got %d", len(metrics))
	}

	metric := metrics[0]
	if metric.Name != "test_metric" {
		t.Errorf("expected metric name 'test_metric', got %s", metric.Name)
	}
	if metric.Type != "counter" {
		t.Errorf("expected metric type 'counter', got %s", metric.Type)
	}
	if metric.Value != 1.0 {
		t.Errorf("expected metric value 1.0, got %f", metric.Value)
	}
	if len(metric.Labels) != 2 {
		t.Errorf("expected 2 labels, got %d", len(metric.Labels))
	}
	if metric.Labels["label1"] != "value1" {
		t.Errorf("expected label1=value1, got %s", metric.Labels["label1"])
	}
}

// TestRecordMultipleMetrics tests recording multiple metrics
func TestRecordMultipleMetrics(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	// Record multiple metrics
	ctx.RecordMetric("metric1", "counter", 1.0, map[string]string{"type": "counter"})
	ctx.RecordMetric("metric2", "histogram", 123.45, map[string]string{"type": "histogram"})
	ctx.RecordMetric("metric3", "gauge", 42.0, map[string]string{"type": "gauge"})

	// Get all metrics
	metrics := ctx.GetAllMetrics()

	// Verify
	if len(metrics) != 3 {
		t.Errorf("expected 3 metrics, got %d", len(metrics))
	}

	// Verify each metric
	metricNames := make(map[string]bool)
	for _, m := range metrics {
		metricNames[m.Name] = true
	}

	if !metricNames["metric1"] || !metricNames["metric2"] || !metricNames["metric3"] {
		t.Error("not all metrics were recorded")
	}
}

// TestClearMetrics tests clearing metrics
func TestClearMetrics(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	// Record metrics
	ctx.RecordMetric("test1", "counter", 1.0, nil)
	ctx.RecordMetric("test2", "counter", 2.0, nil)

	// Verify metrics exist
	metrics := ctx.GetAllMetrics()
	if len(metrics) != 2 {
		t.Errorf("expected 2 metrics before clear, got %d", len(metrics))
	}

	// Clear metrics
	ctx.ClearMetrics()

	// Verify metrics are cleared
	metrics = ctx.GetAllMetrics()
	if len(metrics) != 0 {
		t.Errorf("expected 0 metrics after clear, got %d", len(metrics))
	}
}

// TestMetricLabelIsolation tests that modifying labels doesn't affect stored metrics
func TestMetricLabelIsolation(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	// Create labels
	labels := map[string]string{
		"key": "original",
	}

	// Record metric
	ctx.RecordMetric("test_metric", "counter", 1.0, labels)

	// Modify original labels
	labels["key"] = "modified"
	labels["new_key"] = "new_value"

	// Get metrics
	metrics := ctx.GetAllMetrics()

	// Verify the stored metric labels are not affected
	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	storedLabels := metrics[0].Labels
	if storedLabels["key"] != "original" {
		t.Errorf("expected label to be 'original', got %s", storedLabels["key"])
	}
	if _, exists := storedLabels["new_key"]; exists {
		t.Error("new_key should not exist in stored labels")
	}
}

// TestRecordMetricWithNilLabels tests recording metrics with nil labels
func TestRecordMetricWithNilLabels(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	// Record metric with nil labels
	ctx.RecordMetric("test_metric", "counter", 1.0, nil)

	// Get metrics
	metrics := ctx.GetAllMetrics()

	// Verify
	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if metrics[0].Labels == nil {
		t.Error("labels should not be nil, should be empty map")
	}
	if len(metrics[0].Labels) != 0 {
		t.Errorf("expected 0 labels, got %d", len(metrics[0].Labels))
	}
}

// TestRecordMetricWithEmptyLabels tests recording metrics with empty labels
func TestRecordMetricWithEmptyLabels(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	// Record metric with empty labels
	ctx.RecordMetric("test_metric", "counter", 1.0, map[string]string{})

	// Get metrics
	metrics := ctx.GetAllMetrics()

	// Verify
	if len(metrics) != 1 {
		t.Fatalf("expected 1 metric, got %d", len(metrics))
	}

	if len(metrics[0].Labels) != 0 {
		t.Errorf("expected 0 labels, got %d", len(metrics[0].Labels))
	}
}

// TestGetAllMetricsReturnsCopy tests that GetAllMetrics returns a copy
func TestGetAllMetricsReturnsCopy(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	// Record a metric
	ctx.RecordMetric("test1", "counter", 1.0, map[string]string{"key": "value"})

	// Get metrics twice
	metrics1 := ctx.GetAllMetrics()
	metrics2 := ctx.GetAllMetrics()

	// Verify they are different slices (not the same reference)
	if &metrics1[0] == &metrics2[0] {
		t.Error("GetAllMetrics should return a copy, not the same reference")
	}

	// Verify contents are the same
	if metrics1[0].Name != metrics2[0].Name {
		t.Error("metric contents should be the same")
	}
}

// TestConcurrentMetricRecording tests concurrent metric recording
func TestConcurrentMetricRecording(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	// Record metrics concurrently
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func(id int) {
			ctx.RecordMetric(
				fmt.Sprintf("metric_%d", id),
				"counter",
				float64(id),
				map[string]string{"id": fmt.Sprintf("%d", id)},
			)
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}

	// Verify all metrics were recorded
	metrics := ctx.GetAllMetrics()
	if len(metrics) != 100 {
		t.Errorf("expected 100 metrics, got %d", len(metrics))
	}
}

// TestConcurrentGetAllMetrics tests concurrent calls to GetAllMetrics
func TestConcurrentGetAllMetrics(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	// Pre-populate with metrics
	for i := 0; i < 50; i++ {
		ctx.RecordMetric(fmt.Sprintf("metric_%d", i), "counter", float64(i), nil)
	}

	// Concurrently call GetAllMetrics
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func() {
			metrics := ctx.GetAllMetrics()
			if len(metrics) != 50 {
				t.Errorf("expected 50 metrics, got %d", len(metrics))
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}
}

// TestMetricDataTypes tests different metric types
func TestMetricDataTypes(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	tests := []struct {
		name       string
		metricName string
		metricType string
		value      float64
	}{
		{"counter", "test_counter", "counter", 1.0},
		{"histogram", "test_histogram", "histogram", 123.45},
		{"gauge", "test_gauge", "gauge", 42.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx.ClearMetrics() // Clear before each test

			ctx.RecordMetric(tt.metricName, tt.metricType, tt.value, nil)

			metrics := ctx.GetAllMetrics()
			if len(metrics) != 1 {
				t.Fatalf("expected 1 metric, got %d", len(metrics))
			}

			if metrics[0].Name != tt.metricName {
				t.Errorf("expected name %s, got %s", tt.metricName, metrics[0].Name)
			}
			if metrics[0].Type != tt.metricType {
				t.Errorf("expected type %s, got %s", tt.metricType, metrics[0].Type)
			}
			if metrics[0].Value != tt.value {
				t.Errorf("expected value %f, got %f", tt.value, metrics[0].Value)
			}
		})
	}
}

// BenchmarkRecordMetric benchmarks recording metrics
func BenchmarkRecordMetric(b *testing.B) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	labels := map[string]string{
		"method": "GET",
		"status": "200",
		"path":   "/api/test",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx.RecordMetric("benchmark_metric", "counter", 1.0, labels)
	}
}

// BenchmarkGetAllMetrics benchmarks getting all metrics
func BenchmarkGetAllMetrics(b *testing.B) {
	req, _ := http.NewRequest("GET", "http://example.com/test", nil)
	ctx := newTestHTTPContext(req)

	// Pre-populate with metrics
	for i := 0; i < 10; i++ {
		ctx.RecordMetric(
			fmt.Sprintf("metric_%d", i),
			"counter",
			float64(i),
			map[string]string{"id": fmt.Sprintf("%d", i)},
		)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = ctx.GetAllMetrics()
	}
}
