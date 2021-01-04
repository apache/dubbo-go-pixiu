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
package registry

import (
	"net/url"
	"strconv"
	"strings"
)

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	consul "github.com/hashicorp/consul/api"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

func init() {
	var _ Loader = new(ConsulRegistryLoad)
}

const (
	dubboAPIFilter = "dubbo in Tags"
)

// ConsulRegistryLoad load dubbo apis from consul registry
type ConsulRegistryLoad struct {
	Address string
	// Consul client.
	client  *consul.Client
	cluster string
}

func newConsulRegistryLoad(address, cluster string) (Loader, error) {
	config := &consul.Config{Address: address}
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	r := &ConsulRegistryLoad{
		Address: address,
		client:  client,
		cluster: cluster,
	}

	return r, nil
}

// nolint
func (crl *ConsulRegistryLoad) GetCluster() (string, error) {
	return crl.cluster, nil
}

func (crl *ConsulRegistryLoad) transfer2Url(service consul.AgentService) (*common.URL, error) {
	var params = url.Values{}
	var protocol string

	for _, tag := range service.Tags {
		kv := strings.Split(tag, "=")
		if len(kv) != 2 {
			continue
		}
		params.Add(kv[0], kv[1])
	}

	if url, ok := service.Meta["url"]; ok {
		protocol = strings.Split(url, ":")[0]
	}

	methodsParam := strings.Split(params.Get(constant.METHODS_KEY), ",")
	var methods = make([]string, len(methodsParam))
	for _, method := range methodsParam {
		if method != "" {
			methods = append(methods, method)
		}
	}
	url := common.NewURLWithOptions(common.WithPort(strconv.Itoa(service.Port)),
		common.WithProtocol(protocol), common.WithMethods(methods),
		common.WithIp(service.Address), common.WithParams(params))

	return url, nil
}

// LoadAllServices load all services from consul registry
func (crl *ConsulRegistryLoad) LoadAllServices() ([]*common.URL, error) {
	agentServices, err := crl.client.Agent().ServicesWithFilter(dubboAPIFilter)
	if err != nil {
		logger.Error("consul load all apis error:%v", err)
		return nil, perrors.Wrap(err, "consul load all apis")
	}
	var urls []*common.URL
	for _, service := range agentServices {
		url, err := crl.transfer2Url(*service)
		if err != nil {
			logger.Warnf("consul transfer service to url error:%v", err)
			continue
		}
		urls = append(urls, url)
	}
	return urls, nil
}
