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

package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/retry"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/util"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	Kind                = constant.LLMProxyFilter
	APIKeyPrefix        = "Bearer"
	LLMUnhealthyKey     = "LLMUnhealthy"
	HealthyCheckTimeKey = "HealthyCheckTime"
	// Context key to pass attempt data from proxy to downstream filters
	LLMUpstreamAttemptsKey = "llm_upstream_attempts"
)

// UpstreamAttempt holds details for a single request attempt to an endpoint.
type UpstreamAttempt struct {
	EndpointID      string
	EndpointAddress string
	ClusterName     string
	Success         bool
	StatusCode      int
	ErrorType       string // e.g., "network_error", "status_code_error"
}

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is the main plugin entrypoint.
	Plugin struct{}

	// FilterFactory creates filter instances.
	FilterFactory struct {
		cfg    *Config
		client http.Client
	}

	// Filter is the processing entity for each request.
	Filter struct {
		client         http.Client
		scheme         string
		strategy       *Strategy
		clusterManager *server.ClusterManager
	}

	// Config describes the top-level configuration for the filter.
	// Note: Strategy-specific configurations are now defined on the endpoints.
	Config struct {
		Timeout             time.Duration `yaml:"timeout" json:"timeout,omitempty"`
		MaxIdleConns        int           `yaml:"maxIdleConns" json:"maxIdleConns,omitempty"`
		MaxIdleConnsPerHost int           `yaml:"maxIdleConnsPerHost" json:"maxIdleConnsPerHost,omitempty"`
		MaxConnsPerHost     int           `yaml:"maxConnsPerHost" json:"maxConnsPerHost,omitempty"`
		Scheme              string        `yaml:"scheme" json:"scheme,omitempty" default:"http"`
	}

	RequestExecutor struct {
		hc             *contexthttp.HttpContext
		filter         *Filter
		clusterName    string
		clusterManager *server.ClusterManager
	}
)

// Kind returns the unique name of this filter.
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilterFactory creates a new factory instance for this filter.
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

// Config returns the configuration struct for the factory.
func (factory *FilterFactory) Config() any {
	return factory.cfg
}

// Apply initializes the factory from its configuration.
func (factory *FilterFactory) Apply() error {
	scheme := strings.TrimSpace(strings.ToLower(factory.cfg.Scheme))
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("%s: scheme must be http or https", Kind)
	}
	factory.cfg.Scheme = scheme

	cfg := factory.cfg
	factory.client = http.Client{
		Timeout: cfg.Timeout,
		Transport: &http.Transport{
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			MaxConnsPerHost:     cfg.MaxConnsPerHost,
		},
	}
	return nil
}

// PrepareFilterChain creates a new Filter instance for a request chain.
func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{
		client:         factory.client,
		scheme:         factory.cfg.Scheme,
		strategy:       &Strategy{},
		clusterManager: server.GetClusterManager(),
	}
	chain.AppendDecodeFilters(f)
	return nil
}

// Decode is the main entry point for processing an incoming request.
func (f *Filter) Decode(hc *contexthttp.HttpContext) filter.FilterStatus {
	rEntry := hc.GetRouteEntry()
	if rEntry == nil {
		sendJSONError(hc, http.StatusBadRequest, "no route entry found for request")
		return filter.Stop
	}
	logger.Debugf("[dubbo-go-pixiu] client choose endpoint from cluster: %v", rEntry.Cluster)

	// Ensure the request body can be re-read for retries
	if err := f.prepareRequestBody(hc); err != nil {
		sendJSONError(hc, http.StatusInternalServerError, fmt.Sprintf("failed to read request body: %v", err))
		return filter.Stop
	}
	defer hc.Request.Body.Close()

	// Set up the context for our strategy executor
	executor := &RequestExecutor{
		hc:             hc,
		filter:         f,
		clusterName:    rEntry.Cluster,
		clusterManager: f.clusterManager,
	}

	// Delegate the complex execution logic to the strategy
	resp, err := f.strategy.Execute(executor)

	// Handle the outcome
	if err != nil {
		logger.Infof("[dubbo-go-pixiu] request execution failed after all attempts: %v", err)
		var urlErr *url.Error
		if errors.As(err, &urlErr) && urlErr.Timeout() {
			sendJSONError(hc, http.StatusGatewayTimeout, err.Error())
		} else if resp == nil {
			// This handles errors where no response was ever received (e.g., DNS error, connection refused)
			sendJSONError(hc, http.StatusServiceUnavailable, err.Error())
		} else {
			// A response was received, but it was a failure. Pass it along.
			hc.SourceResp = resp
		}
		return filter.Continue // Let the response writer handle the failed response
	}

	logger.Debugf("[dubbo-go-pixiu] client call successful, resp status: %s", resp.Status)
	hc.SourceResp = resp
	return filter.Continue
}

// prepareRequestBody ensures the request body can be read multiple times.
func (f *Filter) prepareRequestBody(hc *contexthttp.HttpContext) error {
	if hc.Request.Body == nil || hc.Request.GetBody != nil {
		return nil // Nothing to do
	}

	bodyBytes, err := io.ReadAll(hc.Request.Body)
	if err != nil {
		return err
	}
	hc.Request.Body.Close() // Close the original body

	// Set the body to a new reader and provide a function to get a new reader for later reads
	hc.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
	hc.Request.GetBody = func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(bodyBytes)), nil
	}
	return nil
}

// assembleRequest creates a new http.Request for a specific endpoint.
func (f *Filter) assembleRequest(endpoint *model.Endpoint, r *http.Request) (*http.Request, error) {
	// Reset the body to the beginning for each new request attempt
	if r.GetBody != nil {
		var err error
		r.Body, err = r.GetBody()
		if err != nil {
			return nil, err
		}
	}

	parsedURL := url.URL{
		Host:     endpoint.Address.GetAddress(),
		Scheme:   f.scheme,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}

	req, err := http.NewRequest(r.Method, parsedURL.String(), r.Body)
	if err != nil {
		return nil, err
	}
	// Copy headers from original request
	req.Header = r.Header

	// replace the header value with the endpoint's api key
	if apiKey := endpoint.LLMMeta.APIKey; apiKey != "" {
		req.Header.Set(constant.HeaderValueAuthorization, fmt.Sprintf("%s %s", APIKeyPrefix, apiKey))
	}
	req.Header.Set(constant.HeaderKeyUserAgent, fmt.Sprintf("%s %s", constant.Name, constant.Version))

	return req, nil
}

// Strategy is a stateless executor that orchestrates the request lifecycle.
type Strategy struct{}

// Execute orchestrates the request lifecycle using dynamic policies from endpoints.
func (s *Strategy) Execute(executor *RequestExecutor) (*http.Response, error) {
	var (
		resp     *http.Response
		err      error
		attempts []UpstreamAttempt
	)

	// 1. Pick initial endpoint from the cluster based on load balancing.
	endpoint := executor.clusterManager.PickEndpoint(executor.clusterName, executor.hc)

	// 2. The main fallback loop. It continues as long as we have a valid endpoint to try.
	for endpoint != nil {
		logger.Debugf("[dubbo-go-pixiu] client attempting endpoint [%s: %v]", endpoint.ID, endpoint.Address.GetAddress())

		// 3. Check the health of current endpoint,
		if unhealthy, ok := endpoint.Metadata[LLMUnhealthyKey]; ok && unhealthy == "true" {
			// check the health cooldown time
			if t, ok := endpoint.Metadata[HealthyCheckTimeKey]; ok {
				lt, err := time.Parse(time.RFC3339, t)
				if err == nil && time.Since(lt) < time.Millisecond*time.Duration(endpoint.LLMMeta.HealthCheckInterval) {
					logger.Debugf("[dubbo-go-pixiu] endpoint [%s: %v] is still in unhealthy cooldown period. Skipping to next endpoint.", endpoint.ID, endpoint.Address.GetAddress())
					endpoint = getNextFallbackEndpoint(endpoint, executor)
					continue
				}
				// The Cooldown period has passed, ready for a new attempt
				delete(endpoint.Metadata, LLMUnhealthyKey)
				delete(endpoint.Metadata, HealthyCheckTimeKey)
				logger.Debugf("[dubbo-go-pixiu] endpoint [%s: %v] cooldown period passed. Retrying this endpoint.", endpoint.ID, endpoint.Address.GetAddress())
			}
		}

		// 4. Dynamically load the retry policy for the current endpoint
		var retryPolicy retry.RetryPolicy
		retryPolicy, err = retry.GetRetryPolicy(endpoint)
		if err != nil {
			logger.Errorf("could not load retry policy for endpoint [%s: %v]. Skipping to next endpoint.", endpoint.ID, err)
			endpoint = getNextFallbackEndpoint(endpoint, executor)
			continue
		}
		retryPolicy.Reset()

		// 5. The retry loop for the current endpoint.
		for retryPolicy.Attempt() {
			var req *http.Request
			req, err = executor.filter.assembleRequest(endpoint, executor.hc.Request)
			if err != nil {
				// Request assembly error is fatal for this endpoint, break retry loop to go to fallback
				logger.Warnf("[dubbo-go-pixiu] failed to assemble request for endpoint [%s: %v]: %v. Skipping to next endpoint.", endpoint.ID, endpoint.Address.GetAddress(), err)
				break
			}

			resp, err = executor.filter.client.Do(req)

			attempt := UpstreamAttempt{
				EndpointID:      endpoint.ID,
				EndpointAddress: endpoint.Address.GetAddress(),
				ClusterName:     executor.clusterName,
			}

			if err != nil {
				logger.Warnf("[dubbo-go-pixiu] request to endpoint [%s: %v] failed: %v", endpoint.ID, endpoint.Address.GetAddress(), err)
				attempt.Success = false
				attempt.ErrorType = "network_error"
				attempts = append(attempts, attempt)
				break
			}

			attempt.StatusCode = resp.StatusCode
			if util.IsHTTPRespSuccessful(resp.StatusCode) {
				attempt.Success = true
				attempts = append(attempts, attempt)
				executor.hc.Params[LLMUpstreamAttemptsKey] = attempts
				return resp, nil
			}

			attempt.Success = false
			attempt.ErrorType = "status_code_error"
			attempts = append(attempts, attempt)

			logger.Debugf("[dubbo-go-pixiu] attempt failed for endpoint [%s: %v]. Error: %v, Status: %s trying to retry",
				endpoint.ID, endpoint.Address.GetAddress(), err, resp.Status)
		}

		// 6. If we are here, all retries for the current endpoint are exhausted.
		// Get the next endpoint for fallback. The loop will terminate if it's nil.
		endpoint.Metadata[LLMUnhealthyKey] = "true"
		endpoint.Metadata[HealthyCheckTimeKey] = time.Now().Format(time.RFC3339)
		endpoint = getNextFallbackEndpoint(endpoint, executor)
	}

	// 7. If we've exited the loop, all attempts and fallbacks have failed.
	executor.hc.Params[LLMUpstreamAttemptsKey] = attempts

	// Return the last known error and response.
	if err == nil && resp != nil {
		err = fmt.Errorf("request failed with status code %d after all retries and fallbacks", resp.StatusCode)
	} else if err == nil {
		err = errors.New("all retries and fallbacks failed without a definitive error or response")
	}
	return resp, err
}

// getNextFallbackEndpoint checks if fallback is enabled and returns the next endpoint.
func getNextFallbackEndpoint(currentEndpoint *model.Endpoint, executor *RequestExecutor) *model.Endpoint {
	// Fallback is controlled by a boolean in metadata. Default to false.
	if !currentEndpoint.LLMMeta.Fallback {
		return nil // Fallback disabled, end the process.
	}

	nextEndpoint := executor.clusterManager.PickNextEndpoint(executor.clusterName, currentEndpoint.ID)
	if nextEndpoint != nil {
		logger.Debugf("[dubbo-go-pixiu] client fallback to endpoint [%s: %v]", nextEndpoint.ID, nextEndpoint.Address.GetAddress())
	} else {
		logger.Debugf("[dubbo-go-pixiu] no more fallback endpoints available.")
	}

	return nextEndpoint
}

// sendJSONError is a helper to send a structured JSON error message.
func sendJSONError(hc *contexthttp.HttpContext, code int, message string) {
	bt, _ := json.Marshal(contexthttp.ErrResponse{Message: message})
	hc.SendLocalReply(code, bt)
}
