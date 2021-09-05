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
	"strings"
	"time"
)

import (
	"github.com/apache/dubbo-go/common"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type RegisteredType int8

const (
	RegisteredTypeApplication RegisteredType = iota
	RegisteredTypeInterface
)

var registryMap = make(map[string]func(model.Registry) (Registry, error), 8)

func (t *RegisteredType) String() string {
	return []string{"application", "interface"}[*t]
}

// Registry interface defines the basic features of a registry
type Registry interface {
	// Subscribe monitors the target registry.
	Subscribe() error
	// Unsubscribe stops monitoring the target registry.
	Unsubscribe() error
}

// SetRegistry will store the registry by name
func SetRegistry(name string, newRegFunc func(model.Registry) (Registry, error)) {
	registryMap[name] = newRegFunc
}

// GetRegistry will return the registry
// if not found, it will panic
func GetRegistry(name string, regConfig model.Registry) Registry {
	if registry, ok := registryMap[name]; ok {
		reg, err := registry(regConfig)
		if err != nil {
			panic("Initialize Registry" + name + "failed due to: " + err.Error())
		}
		return reg
	}
	panic("Registry " + name + " does not support yet")
}

// CreateAPIConfig returns router.API struct base on the input
func CreateAPIConfig(urlPattern string, dboBackendConfig config.DubboBackendConfig, methodString string, mappingParams []config.MappingParam) router.API {
	dboBackendConfig.Method = methodString
	url := strings.Join([]string{urlPattern, methodString}, constant.PathSlash)
	method := config.Method{
		Enable:    true,
		Timeout:  3 * time.Second,
		Mock:     false,
		HTTPVerb: config.MethodPost,
		InboundRequest: config.InboundRequest{
			RequestType: config.HTTPRequest,
		},
		IntegrationRequest: config.IntegrationRequest{
			RequestType:        config.DubboRequest,
			DubboBackendConfig: dboBackendConfig,
			MappingParams:      mappingParams,
		},
	}
	return router.API{
		URLPattern: url,
		Method:     method,
	}
}

// ParseDubboString parse the dubbo urls
// dubbo://192.168.3.46:20002/org.apache.dubbo.UserProvider2?anyhost=true&app.version=0.0.1&application=UserInfoServer&bean.name=UserProvider
// &cluster=failover&environment=dev&export=true&interface=org.apache.dubbo.UserProvider2&ip=192.168.3.46&loadbalance=random&message_size=4
// &methods=GetUser&methods.GetUser.loadbalance=random&methods.GetUser.retries=1&methods.GetUser.weight=0&module=dubbo-go user-info server
// &name=UserInfoServer&organization=dubbo.io&pid=11037&registry.role=3&release=dubbo-golang-1.5.6
// &service.filter=echo,token,accesslog,tps,generic_service,execute,pshutdown&side=provider&ssl-enabled=false&timestamp=1624716984&warmup=100
func ParseDubboString(urlString string) (config.DubboBackendConfig, []string, error) {
	url, err := common.NewURL(urlString)
	if err != nil {
		return config.DubboBackendConfig{}, nil, errors.WithStack(err)
	}
	return config.DubboBackendConfig{
		ClusterName:     url.GetParam(constant.ClusterKey, ""),
		ApplicationName: url.GetParam(constant.ApplicationKey, ""),
		Version:         url.GetParam(constant.AppVersionKey, ""),
		Protocol:        string(config.DubboRequest),
		Group:           url.GetParam(constant.GroupKey, ""),
		Interface:       url.GetParam(constant.InterfaceKey, ""),
		Retries:         url.GetParam(constant.RetriesKey, ""),
	}, strings.Split(url.GetParam(constant.MethodsKey, ""), constant.StringSeparator), nil
}

// GetAPIPattern generate the API path pattern. /application/interface/version
func GetAPIPattern(bkConfig config.DubboBackendConfig) string {
	return strings.Join([]string{"/" + bkConfig.ApplicationName, bkConfig.Interface, bkConfig.Version}, constant.PathSlash)
}

func GetRouter() model.Router {
	return model.Router{
		Match: model.RouterMatch{
			Prefix:  "",
			Path:    "",
			Regex:   "",
			Methods: nil,
			Headers: nil,
		},

	}
}