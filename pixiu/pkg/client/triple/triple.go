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

package triple

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/router"
	"io"
	"net/url"
	"strings"
	"sync"
)

import (
	proxymeta "github.com/mercari/grpc-http-proxy/metadata"
	"github.com/mercari/grpc-http-proxy/proxy"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
)

// InitDefaultTripleClient init default dubbo client
func InitDefaultTripleClient() {
	tripleClient = NewTripleClient()
}

var (
	clientOnce   sync.Once
	tripleClient *Client
)

// NewTripleClient create dubbo client
func NewTripleClient() *Client {
	clientOnce.Do(func() {
		tripleClient = &Client{}
	})
	return tripleClient
}

// Client client to generic invoke dubbo
type Client struct {
}

func (tc *Client) Apply() error {
	panic("implement me")
}

// suppose proto file and pixiu files as follows:

//test.proto
//	message C {
//		string f = 1;
//	}
//
//	message TestRequest {
//		int64 a = 1;
//		string b = 2;
//		C c = 3;
//		repeated string e = 4;
//	}
//
//	message TestResponse {
//		int64 result = 1;
//	}
//
//	service TestInterface {
//		rpc Test(TestRequest) returns (TestResponse) {}
//	}

// api-config.yaml
//		name: test
//		description: http -> dubbo for test service
//		resources:
//		- path: '/test/:a'
//		type: restful
//		description: test
//		methods:
//		- httpVerb: GET
//		enable: true
//		inboundRequest:
//			requestType: http
//		integrationRequest:
//			requestType: triple
//			interface: "test.TestInterface"
//			method: "Test"

// =================================SUPPORT===========================
//  1. support uri, body, query
//    HTTP:  GET /test/1?b=3&a=2&c.f=fvalue
//	  DUBBO: get TestRequest
//		{
//			"a": 2,
//			"b": "3",
//			"c": {
//				"f": "fvalue"
//			}
//		}

//    HTTP:  POST /test/1?c.f=fvalue
//                body:
//					{
//						"a": 0,
//						"b": "tony",
//						"c": {
//							"f": "f"
//						}
//					}
//	  DUBBO: get TestRequest
//		{
//			"a": 1,
//			"b": "tony",
//			"c": {
//				"f": "fvalue"
//			}
//		}

// 2. same for other http verbs; values is filled with the following priority
//    QUERY > URI > BODY

// ==============================NOT SUPPORT===========================
// 1. do not support array; which has repeated filed in proto files
//    for example:
// 	  	  HTTP:  GET /test/1?b=3&c.f=f&e=3&e=4
//		  DUBBO: get TestRequest(e will be ignored)
//			{
//				"a": 1,
//				"b": "3",
//				"c": {
//					"f": "f"
//				}
//			}

//    	  HTTP:  GET /test/1?b=3&e=3 ===> c
//		  DUBBO: get error: message type mismatch
// 2. do not support header

func (tc *Client) MapParams(req *client.Request) (reqData interface{}, err error) {
	//===== body =====
	rawBody, err := io.ReadAll(req.IngressRequest.Body)
	defer func() {
		req.IngressRequest.Body = io.NopCloser(bytes.NewReader(rawBody))
	}()
	if err != nil {
		return nil, err
	}
	mapBody := map[string]interface{}{}
	if err := json.Unmarshal(rawBody, &mapBody); err != nil {
		return nil, err
	}

	// ===== uri =====
	uriValues := router.GetURIParams(&req.API, *req.IngressRequest.URL)
	for k, v := range uriValues {
		mapBody[k] = v
	}

	// ===== query =====
	queryParams := req.IngressRequest.URL.Query()
	for k, v := range queryParams {
		if len(v) == 1 {
			keys := strings.Split(k, ".")
			if err := client.SetMapValue(mapBody, keys, queryParams.Get(k), ""); err != nil {
				return nil, err
			}
		}
	}
	return mapBody, err
}

// Close clear GenericServicePool.
func (dc *Client) Close() error {
	return nil
}

// SingletonTripleClient singleton dubbo clent
func SingletonTripleClient() *Client {
	return NewTripleClient()
}

// Call invoke service
func (dc *Client) Call(req *client.Request) (res interface{}, err error) {
	address := strings.Split(req.API.IntegrationRequest.HTTPBackendConfig.URL, ":")
	p := proxy.NewProxy()
	targetURL := &url.URL{
		Scheme: address[0],
		Opaque: address[1],
	}
	if err := p.Connect(context.Background(), targetURL); err != nil {
		return "", errors.Errorf("connect triple server error = %s", err)
	}
	meta := make(map[string][]string)
	reqBody, err := dc.MapParams(req)
	if err != nil {
		return "", err
	}
	reqData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()
	call, err := p.Call(ctx, req.API.Method.IntegrationRequest.Interface, req.API.Method.IntegrationRequest.Method, reqData, (*proxymeta.Metadata)(&meta))
	if err != nil {
		return "", errors.Errorf("call triple server error = %s", err)
	}
	return call, nil
}
