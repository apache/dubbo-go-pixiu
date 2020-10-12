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
package registry_load

import (
	"github.com/apache/dubbo-go/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"net/url"
	"strconv"
	"strings"
)

import (
	"github.com/apache/dubbo-go/common"
)

import (
	consul "github.com/hashicorp/consul/api"
	perrors "github.com/pkg/errors"
)

func init() {
	var _ RegistryLoad = new(ConsulRegistryLoad)
}

const (
	dubbo_api_filter = "dubbo in Tags"
)

type ConsulRegistryLoad struct {
	Address string
	// Consul client.
	client *consul.Client
}

func newConsulRegistryLoad(address string) (RegistryLoad, error) {
	config := &consul.Config{Address: address}
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	r := &ConsulRegistryLoad{
		Address: address,
		client:  client,
	}

	return r, nil
}

func (crl *ConsulRegistryLoad) transfer2Url(service consul.AgentService) (common.URL, error) {
	var err error

	var params url.Values

	for _, tag := range service.Tags {
		kv := strings.Split(tag, "=")
		if len(kv) != 2 {
			continue
		}
		params.Add(kv[0], kv[1])
	}

	url := common.NewURLWithOptions(common.WithPort(strconv.Itoa(service.Port)),
		common.WithIp(service.Address), common.WithParams(params))

	// tags

	// meta
	meta := make(map[string]string, 8)
	meta["url"] = url.String()

	// check
	check := &consul.AgentServiceCheck{
		TCP:                            tcp,
		Interval:                       url.GetParam("consul-check-interval", "10s"),
		Timeout:                        url.GetParam("consul-check-timeout", "1s"),
		DeregisterCriticalServiceAfter: url.GetParam("consul-deregister-critical-service-after", "20s"),
	}

	service := &consul.AgentServiceRegistration{
		Name:  url.Service(),
		Tags:  tags,
		Meta:  meta,
		Check: check,
	}

	return *url, nil
}

func (crl *ConsulRegistryLoad) LoadAllApis() ([]common.URL, error) {
	agentServices, err := crl.client.Agent().ServicesWithFilter(dubbo_api_filter)
	if err != nil {
		logger.Error("consul load all apis error:%v", err)
		return nil, perrors.Wrap(err, "consul load all apis")
	}
	var urls []common.URL
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
