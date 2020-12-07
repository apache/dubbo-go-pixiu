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
	"time"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
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

// Client client to generic invoke dubbo
type Client struct {
}

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

// Init init dubbo, config mapping can do here
func (dc *Client) Init() error {
	return nil
}

// Close close
func (dc *Client) Close() error {
	return nil
}

// Call invoke service
func (dc *Client) Call(req *client.Request) (resp interface{}, err error) {
	urlStr := req.GetURL()
	u, err := url.ParseRequestURI(urlStr)
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
	// url query add.
	urlStr = strings.TrimRight(u.String(), "/") + "?" + params.Query.Encode()

	newReq, err := http.NewRequest(req.IngressRequest.Method, urlStr, params.Body)

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
