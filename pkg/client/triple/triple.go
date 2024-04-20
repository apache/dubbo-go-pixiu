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
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

import (
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/client/proxy"
)

// InitDefaultTripleClient init default dubbo client
func InitDefaultTripleClient(protoset []string) {
	tripleClient = NewTripleClient(protoset)
}

var (
	clientOnce   sync.Once
	tripleClient *Client
)

// NewTripleClient create dubbo client
func NewTripleClient(protoset []string) *Client {
	clientOnce.Do(func() {
		tripleClient = &Client{}
		proxy.InitProtosetSource(protoset)
	})
	return tripleClient
}

// Client client to generic invoke dubbo
type Client struct {
}

func (tc *Client) Apply() error {
	panic("implement me")
}

func (tc *Client) MapParams(req *client.Request) (reqData interface{}, err error) {
	panic("implement me")
}

// Close clear GenericServicePool.
func (tc *Client) Close() error {
	return nil
}

// SingletonTripleClient singleton dubbo clent
func SingletonTripleClient(protoset []string) *Client {
	return NewTripleClient(protoset)
}

// Call invoke service
func (tc *Client) Call(req *client.Request) (res interface{}, err error) {
	address := strings.Split(req.API.IntegrationRequest.HTTPBackendConfig.URL, ":")
	p := proxy.NewProxy()
	targetURL := &url.URL{
		Scheme: address[0],
		Opaque: address[1],
	}
	if err = p.Connect(context.Background(), targetURL); err != nil {
		return nil, errors.Errorf("connect triple server error = %s", err)
	}

	header := tc.forwardHeaders(req.IngressRequest.Header)
	ctx := metadata.NewOutgoingContext(context.Background(), header)
	meta := make(map[string][]string)
	reqData, err := io.ReadAll(req.IngressRequest.Body)
	if err != nil {
		return nil, errors.Errorf("read request body error = %s", err)
	}

	ctx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()
	call, err := p.Call(ctx, req.API.Method.IntegrationRequest.Interface, req.API.Method.IntegrationRequest.Method, reqData, (*metadata.MD)(&meta))
	if err != nil {
		return "", errors.Errorf("call triple server error = %s", err)
	}
	return call, nil
}

// forwardHeaders specific what headers should be forwarded
func (tc *Client) forwardHeaders(header http.Header) metadata.MD {
	md := metadata.MD{}
	for k, vals := range header {
		if s := strings.ToLower(k); strings.HasPrefix(s, "tri-") {
			md.Set(k, vals...)
		}
	}
	return md
}
