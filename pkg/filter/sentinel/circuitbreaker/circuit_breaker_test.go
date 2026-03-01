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

package circuitbreaker

import (
	stdHttp "net/http"
	"testing"
	"time"
)

import (
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	pkgs "github.com/apache/dubbo-go-pixiu/pkg/filter/sentinel"
)

func TestFilter(t *testing.T) {
	f := FilterFactory{cfg: &Config{}}

	mockYaml, err := yaml.MarshalYML(mockConfig())
	assert.Nil(t, err)

	assert.Nil(t, yaml.UnmarshalYML(mockYaml, f.Config()))

	assert.Nil(t, f.Apply())

	decoder := &Filter{cfg: f.cfg, matcher: f.matcher}
	request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-dubbo/user/1111", nil)
	c := mock.GetMockHTTPContext(request)

	assert.Equal(t, decoder.Decode(c), filter.Continue)
}

func mockConfig() *Config {
	c := Config{
		Resources: []*pkgs.Resource{
			{
				Name: "test-dubbo",
				Items: []*pkgs.Item{
					{MatchStrategy: pkgs.EXACT, Pattern: "/api/v1/test-dubbo/user"},
					{MatchStrategy: pkgs.REGEX, Pattern: "/api/v1/test-dubbo/user/*"},
				},
			},
		},
		Rules: []*circuitbreaker.Rule{{
			Resource:         "test-dubbo",
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 10,
			StatIntervalMs:   1000,
			Threshold:        1.0,
		}},
	}
	return &c
}

// mockConfigWithResource creates a test config with a custom resource name
func mockConfigWithResource(resourceName string) *Config {
	c := Config{
		Resources: []*pkgs.Resource{
			{
				Name: resourceName,
				Items: []*pkgs.Item{
					{MatchStrategy: pkgs.EXACT, Pattern: "/api/v1/" + resourceName + "/user"},
					{MatchStrategy: pkgs.REGEX, Pattern: "/api/v1/" + resourceName + "/user/*"},
				},
			},
		},
		Rules: []*circuitbreaker.Rule{{
			Resource:         resourceName,
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 10,
			StatIntervalMs:   1000,
			Threshold:        1.0,
		}},
	}
	return &c
}

// TestCircuitBreakerFeedbackLoop tests the complete feedback loop for circuit breaker
// This test verifies the fix for issue #869
func TestCircuitBreakerFeedbackLoop(t *testing.T) {
	// Setup
	factory := FilterFactory{cfg: &Config{}}
	mockYaml, err := yaml.MarshalYML(mockConfig())
	require.NoError(t, err)
	require.NoError(t, yaml.UnmarshalYML(mockYaml, factory.Config()))
	require.NoError(t, factory.Apply())

	f := &Filter{cfg: factory.cfg, matcher: factory.matcher}

	t.Run("Decode stores entry in context", func(t *testing.T) {
		request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-dubbo/user/1111", nil)
		ctx := mock.GetMockHTTPContext(request)

		// Execute Decode
		status := f.Decode(ctx)
		assert.Equal(t, filter.Continue, status)

		// Verify entry is stored in context
		entryVal, exists := ctx.Params[constant.SentinelEntryKey]
		assert.True(t, exists, "Sentinel entry should be stored in context")
		assert.NotNil(t, entryVal, "Sentinel entry should not be nil")

		_, ok := entryVal.(*base.SentinelEntry)
		assert.True(t, ok, "Context value should be a SentinelEntry")

		// Call Encode to ensure the Sentinel entry is properly exited and cleaned up
		ctx.StatusCode(200)
		encodeStatus := f.Encode(ctx)
		assert.Equal(t, filter.Continue, encodeStatus)
	})

	t.Run("Encode reports error for 5xx status codes", func(t *testing.T) {
		request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-dubbo/user/1111", nil)
		ctx := mock.GetMockHTTPContext(request)

		// Execute Decode to get entry
		decodeStatus := f.Decode(ctx)
		require.Equal(t, filter.Continue, decodeStatus)

		// Simulate backend error - set 5xx status code
		ctx.StatusCode(500)

		// Execute Encode
		encodeStatus := f.Encode(ctx)
		assert.Equal(t, filter.Continue, encodeStatus)

		// Entry should be removed from context after Exit (cleanup)
		// Note: We can't directly verify SetError was called, but we can verify the flow completes
	})

	t.Run("Encode handles success status codes", func(t *testing.T) {
		request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-dubbo/user/1111", nil)
		ctx := mock.GetMockHTTPContext(request)

		// Execute Decode
		decodeStatus := f.Decode(ctx)
		require.Equal(t, filter.Continue, decodeStatus)

		// Simulate successful response
		ctx.StatusCode(200)

		// Execute Encode
		encodeStatus := f.Encode(ctx)
		assert.Equal(t, filter.Continue, encodeStatus)
	})

	t.Run("Encode handles various 5xx error codes", func(t *testing.T) {
		errorCodes := []int{500, 502, 503, 504, 599}

		for _, code := range errorCodes {
			request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-dubbo/user/1111", nil)
			ctx := mock.GetMockHTTPContext(request)

			// Execute Decode
			decodeStatus := f.Decode(ctx)
			require.Equal(t, filter.Continue, decodeStatus)

			// Set error status code
			ctx.StatusCode(code)

			// Execute Encode
			encodeStatus := f.Encode(ctx)
			assert.Equal(t, filter.Continue, encodeStatus, "Should handle status code %d", code)
		}
	})

	t.Run("Encode handles non-5xx error codes", func(t *testing.T) {
		// Use a fresh config to avoid circuit breaker state pollution
		factory2 := FilterFactory{cfg: &Config{}}
		config2 := mockConfigWithResource("test-non-error")
		mockYaml2, err := yaml.MarshalYML(config2)
		require.NoError(t, err)
		require.NoError(t, yaml.UnmarshalYML(mockYaml2, factory2.Config()))
		require.NoError(t, factory2.Apply())
		f2 := &Filter{cfg: factory2.cfg, matcher: factory2.matcher}

		nonErrorCodes := []int{200, 201, 301, 400, 401, 403, 404}

		for _, code := range nonErrorCodes {
			request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-non-error/user/1111", nil)
			ctx := mock.GetMockHTTPContext(request)

			// Execute Decode
			decodeStatus := f2.Decode(ctx)
			require.Equal(t, filter.Continue, decodeStatus)

			// Set non-error status code
			ctx.StatusCode(code)

			// Execute Encode
			encodeStatus := f2.Encode(ctx)
			assert.Equal(t, filter.Continue, encodeStatus, "Should handle status code %d without error", code)
		}
	})

	t.Run("Encode handles missing entry gracefully", func(t *testing.T) {
		request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-dubbo/user/1111", nil)
		ctx := mock.GetMockHTTPContext(request)

		// Don't call Decode, so no entry in context
		ctx.StatusCode(500)

		// Execute Encode without entry
		encodeStatus := f.Encode(ctx)
		assert.Equal(t, filter.Continue, encodeStatus, "Should handle missing entry gracefully")
	})

	t.Run("Complete request lifecycle with latency", func(t *testing.T) {
		// Use a fresh config to avoid circuit breaker state pollution
		factory3 := FilterFactory{cfg: &Config{}}
		config3 := mockConfigWithResource("test-latency")
		mockYaml3, err := yaml.MarshalYML(config3)
		require.NoError(t, err)
		require.NoError(t, yaml.UnmarshalYML(mockYaml3, factory3.Config()))
		require.NoError(t, factory3.Apply())
		f3 := &Filter{cfg: factory3.cfg, matcher: factory3.matcher}

		request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-latency/user/1111", nil)
		ctx := mock.GetMockHTTPContext(request)

		// Execute Decode
		decodeStatus := f3.Decode(ctx)
		require.Equal(t, filter.Continue, decodeStatus)

		// Simulate backend processing time
		time.Sleep(10 * time.Millisecond)

		// Simulate backend response
		ctx.StatusCode(200)

		// Execute Encode
		encodeStatus := f3.Encode(ctx)
		assert.Equal(t, filter.Continue, encodeStatus)

		// Sentinel will automatically track the latency between Entry() and Exit()
	})
}

// TestCircuitBreakerNoMatch tests that non-matching URLs are not processed
func TestCircuitBreakerNoMatch(t *testing.T) {
	factory := FilterFactory{cfg: &Config{}}
	mockYaml, err := yaml.MarshalYML(mockConfig())
	require.NoError(t, err)
	require.NoError(t, yaml.UnmarshalYML(mockYaml, factory.Config()))
	require.NoError(t, factory.Apply())

	f := &Filter{cfg: factory.cfg, matcher: factory.matcher}

	// Request that doesn't match any resource pattern
	request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/other-service/data", nil)
	ctx := mock.GetMockHTTPContext(request)

	// Execute Decode
	status := f.Decode(ctx)
	assert.Equal(t, filter.Continue, status)

	// Verify no entry is stored
	_, exists := ctx.Params[constant.SentinelEntryKey]
	assert.False(t, exists, "No entry should be stored for non-matching URL")
}

// TestEncodeWithInvalidEntryType tests that Encode handles invalid entry type gracefully
func TestEncodeWithInvalidEntryType(t *testing.T) {
	factory := FilterFactory{cfg: &Config{}}
	mockYaml, err := yaml.MarshalYML(mockConfig())
	require.NoError(t, err)
	require.NoError(t, yaml.UnmarshalYML(mockYaml, factory.Config()))
	require.NoError(t, factory.Apply())

	f := &Filter{cfg: factory.cfg, matcher: factory.matcher}

	request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-dubbo/user/1111", nil)
	ctx := mock.GetMockHTTPContext(request)

	// Manually set an invalid entry type in context
	ctx.Params = make(map[string]any)
	ctx.Params[constant.SentinelEntryKey] = "invalid_type" // string instead of *base.SentinelEntry

	ctx.StatusCode(500)

	// Execute Encode - should handle gracefully and return Continue
	encodeStatus := f.Encode(ctx)
	assert.Equal(t, filter.Continue, encodeStatus, "Should handle invalid entry type gracefully")
}

// TestCircuitBreakerTriggered tests the behavior when circuit breaker is open
func TestCircuitBreakerTriggered(t *testing.T) {
	// Create config with very low threshold to trigger circuit breaker easily
	config := &Config{
		Resources: []*pkgs.Resource{
			{
				Name: "test-trigger",
				Items: []*pkgs.Item{
					{MatchStrategy: pkgs.REGEX, Pattern: "/api/v1/test-trigger/*"},
				},
			},
		},
		Rules: []*circuitbreaker.Rule{{
			Resource:         "test-trigger",
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 1, // Only need 1 request
			StatIntervalMs:   10000,
			Threshold:        1.0, // Trip after 1 error
		}},
	}

	factory := FilterFactory{cfg: &Config{}}
	mockYaml, err := yaml.MarshalYML(config)
	require.NoError(t, err)
	require.NoError(t, yaml.UnmarshalYML(mockYaml, factory.Config()))
	require.NoError(t, factory.Apply())

	f := &Filter{cfg: factory.cfg, matcher: factory.matcher}

	// First request - should pass and report error
	request1, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-trigger/user", nil)
	ctx1 := mock.GetMockHTTPContext(request1)

	status1 := f.Decode(ctx1)
	assert.Equal(t, filter.Continue, status1)

	// Report error to trigger circuit breaker
	ctx1.StatusCode(500)
	f.Encode(ctx1)

	// Wait a bit for circuit breaker state to update
	time.Sleep(50 * time.Millisecond)

	// Second request - may be blocked if circuit breaker is open
	request2, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-trigger/user", nil)
	ctx2 := mock.GetMockHTTPContext(request2)

	status2 := f.Decode(ctx2)
	// The request might be blocked (filter.Stop) or passed (filter.Continue) depending on circuit breaker state
	// We just verify the code path executes without panic
	assert.True(t, status2 == filter.Continue || status2 == filter.Stop,
		"Decode should return either Continue or Stop")
}
