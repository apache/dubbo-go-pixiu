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
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

import (
	dgConstant "github.com/apache/dubbo-go/common/constant"
	dg "github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
	"github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
)

// RestMetadata dubbo metadata, api config
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
	_httpClient *Client
	countDown   = sync.Once{}
	dgCfg       dg.ConsumerConfig
)

// Client client to generic invoke dubbo
type Client struct {
	mLock              sync.RWMutex
	GenericServicePool map[string]*dg.GenericService
}

// SingletonHTTPClient singleton HTTP Client
func SingletonHTTPClient() *Client {
	if _httpClient == nil {
		countDown.Do(func() {
			_httpClient = NewHTTPClient()
		})
	}
	return _httpClient
}

// NewHTTPClient create dubbo client
func NewHTTPClient() *Client {
	return &Client{
		mLock:              sync.RWMutex{},
		GenericServicePool: make(map[string]*dg.GenericService, 4),
	}
}

// Init init dubbo, config mapping can do here
func (dc *Client) Init() error {
	dgCfg = dg.GetConsumerConfig()
	dg.SetConsumerConfig(dgCfg)
	dg.Load()
	dc.GenericServicePool = make(map[string]*dg.GenericService)
	return nil
}

// Close close
func (dc *Client) Close() error {
	return nil
}

// Call invoke service
func (dc *Client) Call(req *client.Request) (resp interface{}, err error) {
	if err != nil {
		return nil, err
	}

	request := req.IngressRequest.Clone(req.Context)
	//Map the origin paramters to backend parameters according to the API configure
	transformedParams, err := dc.MapParams(req)
	if err != nil {
		return nil, err
	}
	params, _ := transformedParams.(*requestParams)
	request.Body = params.Body
	request.Header = params.Header

	targetURL, err := dc.parseURL(req, *params)
	if err != nil {
		return nil, err
	}

	newReq, err := http.NewRequest(req.IngressRequest.Method, targetURL, params.Body)

	httpClient := &http.Client{Timeout: 5 * time.Second}
	tmpRet, err := httpClient.Do(newReq)

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
		if req.API.IsWildCardBackendPath() {
			r.URIParams = req.API.GetURIParams(*req.IngressRequest.URL)
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
	if req.API.IsWildCardBackendPath() {
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

func (dc *Client) get(key string) *dg.GenericService {
	dc.mLock.RLock()
	defer dc.mLock.RUnlock()
	return dc.GenericServicePool[key]
}

func (dc *Client) create(key string, dm *RestMetadata) *dg.GenericService {
	referenceConfig := dg.NewReferenceConfig(dm.Interface, context.TODO())
	referenceConfig.InterfaceName = dm.Interface
	referenceConfig.Cluster = dgConstant.DEFAULT_CLUSTER
	var registers []string
	for k := range dgCfg.Registries {
		registers = append(registers, k)
	}
	referenceConfig.Registry = strings.Join(registers, ",")

	if dm.ProtocolTypeStr == "" {
		referenceConfig.Protocol = dubbo.DUBBO
	} else {
		referenceConfig.Protocol = dm.ProtocolTypeStr
	}

	referenceConfig.Version = dm.Version
	referenceConfig.Group = dm.Group
	referenceConfig.Generic = true
	if dm.Retries == "" {
		referenceConfig.Retries = "3"
	} else {
		referenceConfig.Retries = dm.Retries
	}
	dc.mLock.Lock()
	defer dc.mLock.Unlock()
	referenceConfig.GenericLoad(key)
	time.Sleep(200 * time.Millisecond) //sleep to wait invoker create
	clientService := referenceConfig.GetRPCService().(*dg.GenericService)

	dc.GenericServicePool[key] = clientService
	return clientService
}
