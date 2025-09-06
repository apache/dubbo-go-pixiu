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

package tokenizer

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	"go.opentelemetry.io/otel/attribute"

	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"go.opentelemetry.io/otel/metric/unit"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/llm/proxy"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// General constants
	Kind                = constant.LLMTokenizerFilter
	LoggerFmt           = "[Tokenizer] [DOWNSTREAM] "
	PromptTokensDetails = "prompt_tokens_details"
	StreamBufferSIze    = 4096

	// Metric meter name
	meterName = "dubbo-go-pixiu.ai-gateway"

	// Attribute Keys for metrics
	attrClusterName     = "cluster_name"
	attrEndpointID      = "endpoint_id"
	attrEndpointAddress = "endpoint_address"
	attrStatusCode      = "status_code"
	attrErrorType       = "error_type"
	attrModel           = "model"

	// JSON Keys from LLM response
	jsonKeyUsage            = "usage"
	jsonKeyModel            = "model"
	jsonKeyPromptTokens     = "prompt_tokens"
	jsonKeyCompletionTokens = "completion_tokens"
	jsonKeyTotalTokens      = "total_tokens"

	// Server-Sent Events (SSE) constants
	sseDone       = "[DONE]"
	sseDataPrefix = "data:"

	// Metric Names
	metricPromptTokens      = "pixiu_llm_prompt_tokens_total"
	metricCompletionTokens  = "pixiu_llm_completion_tokens_total"
	metricTotalTokens       = "pixiu_llm_total_tokens_total"
	metricUpstreamRequests  = "pixiu_llm_upstream_requests_total"
	metricUpstreamSuccess   = "pixiu_llm_upstream_requests_success_total"
	metricUpstreamFailure   = "pixiu_llm_upstream_requests_failure_total"
	metricTotalDurationSum  = "pixiu_llm_total_duration_microseconds_sum_total"
	metricTTLTSum           = "pixiu_llm_time_to_last_token_milliseconds_sum_total"
	metricStreamingRequests = "pixiu_llm_streaming_requests_total"
)

var (
	// Metric Instruments
	llmPromptTokens                 syncint64.Counter
	llmCompletionTokens             syncint64.Counter
	llmTotalTokens                  syncint64.Counter
	llmUpstreamRequestsTotal        syncint64.Counter
	llmUpstreamRequestsSuccessTotal syncint64.Counter
	llmUpstreamRequestsFailureTotal syncint64.Counter
	llmTotalDurationSum             syncfloat64.Counter
	llmTimeToLastTokenSum           syncfloat64.Counter
	llmStreamingRequestsTotal       syncint64.Counter

	registerOnce sync.Once
	registerErr  error
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin        struct{}
	FilterFactory struct{ cfg *Config }
	Filter        struct {
		cfg            *Config
		start          time.Time
		recordTTFTOnce sync.Once
	}
	Config struct {
		LogToConsole bool `yaml:"log_to_console" json:"log_to_console,omitempty" default:"false"`
	}
)

func (p *Plugin) Kind() string { return Kind }

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() any { return factory.cfg }

func (factory *FilterFactory) Apply() error { return registerLLMMetrics() }

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{cfg: factory.cfg}
	chain.AppendDecodeFilters(f)
	chain.AppendEncodeFilters(f)
	return nil
}

func (f *Filter) Decode(hc *contexthttp.HttpContext) filter.FilterStatus {
	f.start = time.Now()
	return filter.Continue
}

func (f *Filter) Encode(hc *contexthttp.HttpContext) filter.FilterStatus {
	// Report all metrics from a central place
	f.reportUpstreamMetrics(hc)
	f.reportTotalDurationMetric(hc)

	// Handle response body for token counting and TTFT
	encoding := hc.Writer.Header().Get(constant.HeaderKeyContentEncoding)
	switch res := hc.TargetResp.(type) {
	case *client.StreamResponse:
		pr, pw := io.Pipe()
		res.Stream = newTeeReadCloser(res.Stream, pw)
		go f.processStreamResponse(hc, pr, encoding)
	case *client.UnaryResponse:
		f.processUnaryResponse(hc, res.Data, encoding)
	default:
		logger.Warnf(LoggerFmt+"Response type not suitable for token calculation: %T", res)
	}
	return filter.Continue
}

// reportUpstreamMetrics reads attempt data from context and reports endpoint-level metrics.
func (f *Filter) reportUpstreamMetrics(hc *contexthttp.HttpContext) {
	attemptsVal, ok := hc.Params[proxy.LLMUpstreamAttemptsKey]
	if !ok {
		return // No attempt data found, likely not an LLM proxy request
	}
	attempts, ok := attemptsVal.([]proxy.UpstreamAttempt)
	if !ok {
		logger.Warnf(LoggerFmt+"Upstream attempt data in context has wrong type: %T", attemptsVal)
		return
	}

	for _, attempt := range attempts {
		attrs := []attribute.KeyValue{
			attribute.String(attrClusterName, attempt.ClusterName),
			attribute.String(attrEndpointID, attempt.EndpointID),
			attribute.String(attrEndpointAddress, attempt.EndpointAddress),
		}

		llmUpstreamRequestsTotal.Add(hc.Ctx, 1, attrs...)

		if attempt.Success {
			successAttrs := append(attrs, attribute.String(attrStatusCode, strconv.Itoa(attempt.StatusCode)))
			llmUpstreamRequestsSuccessTotal.Add(hc.Ctx, 1, successAttrs...)
		} else {
			failureAttrs := append(attrs,
				attribute.String(attrStatusCode, strconv.Itoa(attempt.StatusCode)),
				attribute.String(attrErrorType, attempt.ErrorType),
			)
			llmUpstreamRequestsFailureTotal.Add(hc.Ctx, 1, failureAttrs...)
		}
	}
}

// reportTotalDurationMetric calculates and reports the total request-response duration.
func (f *Filter) reportTotalDurationMetric(hc *contexthttp.HttpContext) {
	totalDuration := time.Since(f.start)
	clusterName := "unknown"
	if rEntry := hc.GetRouteEntry(); rEntry != nil {
		clusterName = rEntry.Cluster
	}
	durationAttrs := []attribute.KeyValue{
		attribute.String(attrClusterName, clusterName),
		attribute.String(attrStatusCode, strconv.Itoa(hc.GetStatusCode())),
	}
	llmTotalDurationSum.Add(hc.Ctx, float64(totalDuration.Microseconds()), durationAttrs...)
}

// processStreamResponse handles streaming responses to calculate TTFT and count tokens.
func (f *Filter) processStreamResponse(hc *contexthttp.HttpContext, body io.Reader, encoding string) {
	streamStartTime := time.Now()
	decompressedReader, err := getDecompressedReader(body, encoding)
	if err != nil {
		logger.Errorf(LoggerFmt+"could not create decompressing reader: %v", err)
		return
	}
	defer decompressedReader.Close()

	buf := make([]byte, StreamBufferSIze)
	var eventBuffer bytes.Buffer

	for {
		n, err := decompressedReader.Read(buf)
		if n > 0 {
			eventBuffer.Write(buf[:n])
			for {
				event, remaining, found := splitSSEEvent(eventBuffer.Bytes())
				if !found {
					break
				}
				jsonData := strings.TrimSpace(strings.TrimPrefix(string(event), sseDataPrefix))
				if len(jsonData) > 0 {
					f.parseAndReportTokens(hc, []byte(jsonData))
				}
				eventBuffer.Reset()
				eventBuffer.Write(remaining)
			}
		}

		if err == io.EOF {
			if eventBuffer.Len() > 0 {
				jsonData := strings.TrimSpace(strings.TrimPrefix(eventBuffer.String(), sseDataPrefix))
				if len(jsonData) > 0 {
					f.parseAndReportTokens(hc, []byte(jsonData))
				}
			}
			break
		}
		if err != nil {
			logger.Errorf(LoggerFmt+"error reading decompressed stream: %v", err)
			break
		}
	}

	// On the very first successful read, record Time to First Token.
	f.recordTTFTOnce.Do(func() {
		ttlt := time.Since(streamStartTime)
		clusterName := "unknown"
		if rEntry := hc.GetRouteEntry(); rEntry != nil {
			clusterName = rEntry.Cluster
		}
		attrs := attribute.String(attrClusterName, clusterName)

		llmTimeToLastTokenSum.Add(hc.Ctx, float64(ttlt.Milliseconds()), attrs)
		llmStreamingRequestsTotal.Add(hc.Ctx, 1, attrs)
	})
}

// splitSSEEvent finds the next complete SSE event (ending in \n\n) in a byte slice.
func splitSSEEvent(data []byte) (event, remaining []byte, found bool) {
	if i := bytes.Index(data, []byte("\n\n")); i >= 0 {
		return data[:i], data[i+2:], true
	}
	return nil, data, false
}

// processUnaryResponse handles non-streaming responses for token counting.
func (f *Filter) processUnaryResponse(hc *contexthttp.HttpContext, data []byte, encoding string) {
	processedData := data
	if encoding != "" {
		bodyReader := bytes.NewReader(data)
		if decompressedData, ok := decompress(bodyReader, encoding); ok {
			processedData = decompressedData
		}
	}
	if len(processedData) == 0 {
		return
	}
	f.parseAndReportTokens(hc, processedData)
}

// parseAndReportTokens parses a JSON data chunk and reports token usage metrics.
func (f *Filter) parseAndReportTokens(hc *contexthttp.HttpContext, data []byte) {
	if len(data) == 0 || string(data) == sseDone {
		return
	}
	var dataCont map[string]any
	if err := json.Unmarshal(data, &dataCont); err != nil {
		return
	}
	usage, ok := dataCont[jsonKeyUsage].(map[string]any)
	if !ok || usage == nil {
		return
	}
	if f.cfg.LogToConsole {
		f.logUsage(usage)
	}
	modelName := "unknown"
	if m, ok := dataCont[jsonKeyModel].(string); ok {
		modelName = m
	}
	clusterName := "unknown"
	if rEntry := hc.GetRouteEntry(); rEntry != nil {
		clusterName = rEntry.Cluster
	}
	tokenAttrs := []attribute.KeyValue{
		attribute.String(attrModel, modelName),
		attribute.String(attrClusterName, clusterName),
	}
	if pTokens, ok := usage[jsonKeyPromptTokens].(float64); ok {
		llmPromptTokens.Add(hc.Ctx, int64(pTokens), tokenAttrs...)
	}
	if cTokens, ok := usage[jsonKeyCompletionTokens].(float64); ok {
		llmCompletionTokens.Add(hc.Ctx, int64(cTokens), tokenAttrs...)
	}
	if tTokens, ok := usage[jsonKeyTotalTokens].(float64); ok {
		llmTotalTokens.Add(hc.Ctx, int64(tTokens), tokenAttrs...)
	}
}

// logUsage prints usage details to the log.
func (f *Filter) logUsage(usage map[string]any) {
	for key, value := range usage {
		if key == PromptTokensDetails {
			details, ok := value.(map[string]any)
			if !ok {
				logger.Warnf(LoggerFmt+PromptTokensDetails+" is not a map, value: %+v", value)
				continue
			}
			for detailKey, detailValue := range details {
				logger.Infof(LoggerFmt+"Usage | %s: %v", detailKey, detailValue)
			}
		} else {
			logger.Infof(LoggerFmt+"Usage | %s: %v", key, value)
		}
	}
}

// getDecompressedReader returns an io.ReadCloser that decompresses the body based on the encoding.
func getDecompressedReader(body io.Reader, encoding string) (io.ReadCloser, error) {
	switch encoding {
	case constant.HeaderValueGzip:
		return gzip.NewReader(body)
	case constant.HeaderValueDeflate:
		return flate.NewReader(body), nil
	default:
		return io.NopCloser(body), nil
	}
}

// decompress reads all data from a reader and returns the decompressed byte slice.
func decompress(body io.Reader, encoding string) ([]byte, bool) {
	decompressedReader, err := getDecompressedReader(body, encoding)
	if err != nil {
		logger.Errorf(LoggerFmt+"%v", err)
		return nil, false
	}
	defer decompressedReader.Close()

	decompressedData, err := io.ReadAll(decompressedReader)
	if err != nil {
		logger.Errorf(LoggerFmt+"Error reading decompressed stream: %v", err)
		return nil, false
	}
	return decompressedData, true
}

// registerLLMMetrics creates and registers all metrics for this plugin.
func registerLLMMetrics() error {
	registerOnce.Do(func() {
		meter := global.MeterProvider().Meter(meterName)
		var err error

		llmPromptTokens, err = meter.SyncInt64().Counter(metricPromptTokens,
			instrument.WithDescription("Total prompt tokens."),
			instrument.WithUnit("1"))
		if err != nil {
			registerErr = err
			return
		}
		llmCompletionTokens, err = meter.SyncInt64().Counter(metricCompletionTokens,
			instrument.WithDescription("Total completion tokens."),
			instrument.WithUnit("1"))
		if err != nil {
			registerErr = err
			return
		}
		llmTotalTokens, err = meter.SyncInt64().Counter(metricTotalTokens,
			instrument.WithDescription("Total tokens."),
			instrument.WithUnit("1"))
		if err != nil {
			registerErr = err
			return
		}

		llmUpstreamRequestsTotal, err = meter.SyncInt64().Counter(metricUpstreamRequests,
			instrument.WithDescription("Total requests to upstream endpoints."),
			instrument.WithUnit("1"))
		if err != nil {
			registerErr = err
			return
		}
		llmUpstreamRequestsSuccessTotal, err = meter.SyncInt64().Counter(metricUpstreamSuccess,
			instrument.WithDescription("Total successful requests to upstream endpoints."),
			instrument.WithUnit("1"))
		if err != nil {
			registerErr = err
			return
		}
		llmUpstreamRequestsFailureTotal, err = meter.SyncInt64().Counter(metricUpstreamFailure,
			instrument.WithDescription("Total failed requests to upstream endpoints."),
			instrument.WithUnit("1"))
		if err != nil {
			registerErr = err
			return
		}

		llmTotalDurationSum, err = meter.SyncFloat64().Counter(
			metricTotalDurationSum,
			instrument.WithDescription("Sum of total duration of LLM requests in microseconds."),
			instrument.WithUnit(unit.Unit("µs")),
		)
		if err != nil {
			registerErr = err
			return
		}

		llmTimeToLastTokenSum, err = meter.SyncFloat64().Counter(
			metricTTLTSum,
			instrument.WithDescription("Sum of Time to Last Token for streaming responses in milliseconds."),
			instrument.WithUnit("ms"),
		)
		if err != nil {
			registerErr = err
			return
		}

		llmStreamingRequestsTotal, err = meter.SyncInt64().Counter(
			metricStreamingRequests,
			instrument.WithDescription("Total number of streaming LLM requests."),
			instrument.WithUnit("1"),
		)
		if err != nil {
			registerErr = err
			return
		}

		logger.Info("LLM metrics registered successfully.")
	})
	return registerErr
}

// --- TeeReadCloser Implementation ---
// This is necessary to read the response body for tokenizing without consuming it,
// so the original response can still be sent to the client.

type teeReadCloser struct {
	reader   io.Reader
	closer   io.Closer
	writer   io.Writer
	once     sync.Once
	closeErr error
}

func newTeeReadCloser(r io.ReadCloser, w io.Writer) *teeReadCloser {
	return &teeReadCloser{
		reader: r,
		closer: r,
		writer: w,
	}
}

func (t *teeReadCloser) Read(p []byte) (n int, err error) {
	n, err = t.reader.Read(p)
	if n > 0 {
		nw, writeErr := t.writer.Write(p[:n])
		if writeErr != nil {
			logger.Errorf(LoggerFmt+"Error writing to tee writer: %v", writeErr)
			return n, err // Return original read error, but log the write error
		}
		if nw != n {
			logger.Errorf(LoggerFmt+"Short write to tee writer: %d/%d", nw, n)
		}
	}
	return n, err
}

func (t *teeReadCloser) Close() (err error) {
	var (
		closerErr error
		writerErr error
	)
	t.once.Do(func() {
		closerErr = t.closer.Close()
		if closerErr != nil {
			logger.Errorf(LoggerFmt+"Error closing closer: %v", closerErr)
		}

		if t.writer != nil {
			if writerCloser, ok := t.writer.(io.Closer); ok {
				writerErr = writerCloser.Close()
				if writerErr != nil {
					logger.Errorf(LoggerFmt+"Error closing writer: %v", writerErr)
				}
			}
		}
	})

	if closerErr != nil || writerErr != nil {
		err = fmt.Errorf("closing closer error: %w. closing writer error: %w", closerErr, writerErr)
	}
	return err
}
