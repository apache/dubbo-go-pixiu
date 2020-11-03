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

package httpclient

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go/common/constant"
	dg "github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
)

// TODO java class name elem
const (
	JavaStringClassName = "java.lang.String"
	JavaLangClassName   = "java.lang.Long"
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
	_httpClient *HTTPClient
	countDown   = sync.Once{}
	dgCfg       dg.ConsumerConfig
)

// HTTPClient client to generic invoke dubbo
type HTTPClient struct {
	mLock              sync.RWMutex
	GenericServicePool map[string]*dg.GenericService
}

// SingleHTTPClient singleton HTTP Client
func SingleHTTPClient() *HTTPClient {

	if _httpClient == nil {
		countDown.Do(func() {
			_httpClient = NewHTTPClient()
		})
	}
	return _httpClient
}

// NewHTTPClient create dubbo client
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		mLock:              sync.RWMutex{},
		GenericServicePool: make(map[string]*dg.GenericService),
	}
}

// Init init dubbo, config mapping can do here
func (dc *HTTPClient) Init() error {
	dgCfg = dg.GetConsumerConfig()
	dg.SetConsumerConfig(dgCfg)
	dg.Load()
	dc.GenericServicePool = make(map[string]*dg.GenericService)
	return nil
}

// Close close
func (dc *HTTPClient) Close() error {
	return nil
}

// Call invoke service
func (dc *HTTPClient) Call(r *client.Request) (resp client.Response, err error) {

	var urlStr = r.API.IntegrationRequest.HTTPBackendConfig.Protocol + "://" + r.API.IntegrationRequest.HTTPBackendConfig.TargetURL
	var httpClient = &http.Client{Timeout: 5 * time.Second}
	var request = r.IngressRequest.Clone(context.Background())
	request.URL, _ = url.ParseRequestURI(urlStr)
	//TODO header replace, url rewrite....

	var tmpRet, e = httpClient.Do(request)
	var ret = client.Response{Data: tmpRet}
	return ret, e
}

func (dc *HTTPClient) get(key string) *dg.GenericService {
	dc.mLock.RLock()
	defer dc.mLock.RUnlock()
	return dc.GenericServicePool[key]
}

func (dc *HTTPClient) check(key string) bool {
	dc.mLock.RLock()
	defer dc.mLock.RUnlock()
	if _, ok := dc.GenericServicePool[key]; ok {
		return true
	} else {
		return false
	}
}

func (dc *HTTPClient) create(key string, dm *RestMetadata) *dg.GenericService {
	referenceConfig := dg.NewReferenceConfig(dm.Interface, context.TODO())
	referenceConfig.InterfaceName = dm.Interface
	referenceConfig.Cluster = constant.DEFAULT_CLUSTER
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

// Get find a dubbo GenericService
func (dc *HTTPClient) Get(interfaceName, version, group string, dm *RestMetadata) *dg.GenericService {
	key := strings.Join([]string{dm.ApplicationName, interfaceName, version, group}, "_")
	if dc.check(key) {
		return dc.get(key)
	}
	return dc.create(key, dm)
}
