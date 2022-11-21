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
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/router"
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
func (dc *Client) Call(req *client.Request) (resp interface{}, err error) {
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

	tr := otel.Tracer(traceNameHTTPClient)
	_, span := tr.Start(req.Context, "HTTP "+newReq.Method, trace.WithSpanKind(trace.SpanKindClient))
	trace.SpanFromContext(req.Context).SpanContext()
	span.SetAttributes(semconv.HTTPMethodKey.String(newReq.Method))
	span.SetAttributes(semconv.HTTPTargetKey.String(targetURL))
	span.SetAttributes(semconv.HTTPFlavorKey.String(newReq.Proto))
	newReq.Header.Set(jaegerTraceIDInHeader, span.SpanContext().TraceID().String())
	defer span.End()

	tmpRet, err := httpClient.Do(newReq)
	if tmpRet != nil {
		span.SetAttributes(semconv.HTTPStatusCodeKey.Int(tmpRet.StatusCode))
	}
	if err != nil {
		span.AddEvent(semconv.ExceptionEventName, trace.WithAttributes(semconv.ExceptionMessageKey.String(err.Error())))
	}

	return tmpRet, err
}

// MapParams param mapping to api.
func (dc *Client) MapParams(req *client.Request) (reqData interface{}, err error) {
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
