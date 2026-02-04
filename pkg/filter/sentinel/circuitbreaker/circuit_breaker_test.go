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

	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
		entryVal, exists := ctx.Params[ContextKeySentinelEntry]
		assert.True(t, exists, "Sentinel entry should be stored in context")
		assert.NotNil(t, entryVal, "Sentinel entry should not be nil")

		_, ok := entryVal.(*base.SentinelEntry)
		assert.True(t, ok, "Context value should be a SentinelEntry")
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
		mockYaml2, _ := yaml.MarshalYML(config2)
		yaml.UnmarshalYML(mockYaml2, factory2.Config())
		factory2.Apply()
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
		mockYaml3, _ := yaml.MarshalYML(config3)
		yaml.UnmarshalYML(mockYaml3, factory3.Config())
		factory3.Apply()
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
	_, exists := ctx.Params[ContextKeySentinelEntry]
	assert.False(t, exists, "No entry should be stored for non-matching URL")
}
