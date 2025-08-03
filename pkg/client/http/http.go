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
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

import (
	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

// RestMetadata http metadata, api config
type RestMetadata struct {
	ApplicationName      string   `yaml:"application_name" json:"application_name" mapstructure:"application_name"`
	Group                string   `yaml:"group" json:"group" mapstructure:"group"`
	Version              string   `yaml:"version" json:"version" mapstructure:"version"`
	Interface            string   `yaml:"interface" json:"interface" mapstructure:"interface"`
	Method               string   `yaml:"method" json:"method" mapstructure:"method"`
	Types                []string `yaml:"type" json:"types" mapstructure:"types"`
	Retries              string   `yaml:"retries"  json:"retries,omitempty" property:"retries"`
	ClusterName          string   `yaml:"cluster_name"  json:"cluster_name,omitempty" property:"cluster_name"`
	ProtocolTypeStr      string   `yaml:"protocol_type"  json:"protocol_type,omitempty" property:"protocol_type"`
	SerializationTypeStr string   `yaml:"serialization_type"  json:"serialization_type,omitempty" property:"serialization_type"`
}

var (
	httpClient *Client
	countDown  = sync.Once{}
)

const (
	traceNameHTTPClient   = "http-client"
	jaegerTraceIDInHeader = "uber-trace-id"
)

// Client client to generic invoke dubbo
type Client struct{}

// SingletonHTTPClient singleton HTTP Client
func SingletonHTTPClient() *Client {
	if httpClient == nil {
		countDown.Do(func() {
			httpClient = NewHTTPClient()
		})
	}
	return httpClient
}

// NewHTTPClient create dubbo client
func NewHTTPClient() *Client {
	return &Client{}
}

// Apply only init dubbo, config mapping can do here
func (dc *Client) Apply() error {
	return nil
}

// Close close
func (dc *Client) Close() error {
	return nil
}

// Call invoke service
func (dc *Client) Call(req *client.Request) (resp any, err error) {
	// Map the origin parameters to backend parameters according to the API configure
	transformedParams, err := dc.MapParams(req)
	if err != nil {
		return nil, err
	}
	params, _ := transformedParams.(*requestParams)

	targetURL, err := dc.parseURL(req, *params)
	if err != nil {
		return nil, err
	}

	newReq, _ := http.NewRequest(req.IngressRequest.Method, targetURL, params.Body)
	newReq.Header = params.Header
	httpClient := &http.Client{Timeout: req.Timeout}

	// Observability
	tr := otel.Tracer(traceNameHTTPClient)
	_, span := tr.Start(req.Context, "HTTP "+newReq.Method, trace.WithSpanKind(trace.SpanKindClient))
	trace.SpanFromContext(req.Context).SpanContext()
	span.SetAttributes(semconv.HTTPMethodKey.String(newReq.Method))
	span.SetAttributes(semconv.HTTPTargetKey.String(targetURL))
	span.SetAttributes(semconv.HTTPFlavorKey.String(newReq.Proto))
	newReq.Header.Set(jaegerTraceIDInHeader, span.SpanContext().TraceID().String())
	defer span.End()

	// Real request
	tmpRet, err := httpClient.Do(newReq)
	if tmpRet != nil {
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(tmpRet.StatusCode))
	}
	if err != nil {
		span.AddEvent(semconv.ExceptionEventName, trace.WithAttributes(semconv.ExceptionMessageKey.String(err.Error())))
		urlErr, ok := err.(*url.Error)
		if ok && urlErr.Timeout() {
			err = errors.Errorf("http req call timeout err: %s", err.Error())
		}
	}

	return tmpRet, err
}

// MapParams param mapping to api.
func (dc *Client) MapParams(req *client.Request) (reqData any, err error) {
	mp := req.API.IntegrationRequest.MappingParams
	r := newRequestParams()
	if len(mp) == 0 {
		r.Body = req.IngressRequest.Body
		r.Header = req.IngressRequest.Header.Clone()
		queryValues, err := url.ParseQuery(req.IngressRequest.URL.RawQuery)
		if err != nil {
			return nil, errors.New("Retrieve request query parameters failed")
		}
		r.Query = queryValues
		if router.IsWildCardBackendPath(&req.API) {
			r.URIParams = router.GetURIParams(&req.API, *req.IngressRequest.URL)
		}
		return r, nil
	}
	for i := 0; i < len(mp); i++ {
		source, _, err := client.ParseMapSource(mp[i].Name)
		if err != nil {
			return nil, err
		}
		if mapper, ok := mappers[source]; ok {
			if err := mapper.Map(mp[i], req, r, nil); err != nil {
				return nil, err
			}
		}
	}
	return r, nil
}

// ParseURL returns the actual target url. Supports wildcard target path value mapping.
func (dc *Client) parseURL(req *client.Request, params requestParams) (string, error) {
	var schema string
	if len(req.API.IntegrationRequest.HTTPBackendConfig.Schema) == 0 {
		schema = "http"
	} else {
		schema = req.API.IntegrationRequest.HTTPBackendConfig.Schema
	}

	rawPath := req.API.IntegrationRequest.HTTPBackendConfig.Path
	if router.IsWildCardBackendPath(&req.API) {
		paths := strings.Split(
			strings.TrimLeft(req.API.IntegrationRequest.HTTPBackendConfig.Path, constant.PathSlash),
			constant.PathSlash)
		for i := 0; i < len(paths); i++ {
			if strings.HasPrefix(paths[i], constant.PathParamIdentifier) {
				uriParam := string(paths[i][1:len(paths[i])])
				uriValue := params.URIParams.Get(uriParam)
				if len(uriValue) == 0 {
					return "", errors.New("No value for target URI")
				}
				paths[i] = uriValue
			}
		}
		rawPath = strings.Join(paths, constant.PathSlash)
	}

	parsedURL := url.URL{
		Host:     req.API.IntegrationRequest.HTTPBackendConfig.Host,
		Scheme:   schema,
		Path:     rawPath,
		RawQuery: params.Query.Encode(),
	}
	return parsedURL.String(), nil
}

// IsSSEStream check if the response is a SSE stream
func IsSSEStream(resp *http.Response) bool {
	contentType := resp.Header.Get(constant.HeaderKeyContextType)
	return strings.Contains(contentType, constant.HeaderValueTextEventStream)
}

// IsStreamableResponse check if the response is streamable
// Determine whether it is a streaming response based on the following conditions:
// 1. Using Transfer-Encoding: chunked
// 2. The content type indicates that this is a streamable response (e.g. text/event-stream, application/json, etc.)
func IsStreamableResponse(resp *http.Response) bool {
	if IsSSEStream(resp) {
		return true
	}

	// check if the block encoded transfer is used
	transferEncoding := resp.Header.Get(constant.HeaderKeyTransferEncoding)
	if strings.Contains(strings.ToLower(transferEncoding), constant.HeaderValueChunked) {
		return true
	}

	// check the content type
	contentType := resp.Header.Get(constant.HeaderKeyContextType)

	// check if it s a streamable content type
	streamableTypes := []string{
		constant.HeaderValueTextPrefix,
		constant.HeaderValueApplicationJson,
		constant.HeaderValueApplicationNDJson,
		constant.HeaderValueApplicationOctetStream,
	}

	for _, streamableType := range streamableTypes {
		if strings.HasPrefix(contentType, streamableType) {
			// For these content types, if you don't have Content-Length set or if Content-Length is large,
			// may be a good candidate for streaming
			contentLength := resp.Header.Get(constant.HeaderKeyContentLength)

			// If Content-Length is not specified, it is possible that the server is not aware of the content length
			if contentLength == "" {
				return true
			}

			// If the Content-Length is large (> 1MB), streaming is also suitable
			length, err := strconv.ParseInt(contentLength, 10, 64)
			if err == nil && length > 1024*1024 {
				return true
			}
		}
	}

	return false
}
