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

package server

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config/xds"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

// DynamicResourceManager help to management the dynamic resource
type DynamicResourceManager interface {
	GetLds() *model.ApiConfigSource
	GetNode() *model.Node
	GetCds() *model.ApiConfigSource
}

type DynamicResourceManagerImpl struct {
	config *model.DynamicResources
	node   *model.Node
}

func (d DynamicResourceManagerImpl) GetCds() *model.ApiConfigSource {
	return d.config.CdsConfig
}

func (d DynamicResourceManagerImpl) GetNode() *model.Node {
	return d.node
}

func (d DynamicResourceManagerImpl) GetLds() *model.ApiConfigSource {
	return d.config.LdsConfig
}

// createDynamicResourceManger create dynamic resource manager or nil if not config
func createDynamicResourceManger(bs *model.Bootstrap) DynamicResourceManager {
	if err := validate(bs); err != nil {
		panic(err) // settings error panic
	}
	if bs.DynamicResources == nil {
		return nil
	}
	m := &DynamicResourceManagerImpl{
		config: bs.DynamicResources, // todo deep copy as immutable value
		node:   bs.Node,
	}

	if bs.DynamicResources.AdsConfig != nil {
		logger.Warnf("un-support ada_config.")
	}
	_ = xds.StartXdsClient(GetServer().GetListenerManager(), GetServer().GetClusterManager(), m) //todo graceful shutdown
	return m
}

// validate validate and make apiType
func validate(bs *model.Bootstrap) error {
	if bs.DynamicResources == nil {
		return nil
	}

	if err := convertApiType(bs.DynamicResources.LdsConfig); err != nil {
		return err
	}
	if err := convertApiType(bs.DynamicResources.CdsConfig); err != nil {
		return err
	}

	if bs.DynamicResources.AdsConfig != nil {
		return errors.Errorf("un-support ads config")
	}

	return nil
}

func convertApiType(config *model.ApiConfigSource) error {
	if config != nil {
		apiType, ok := model.ApiTypeValue[config.APITypeStr]
		if !ok {
			return errors.Errorf("unknown apiType %s", config.APITypeStr)
		}
		config.APIType = api.ApiType(apiType)
		if config.APIType != model.ApiTypeGRPC && config.APIType != model.ApiTypeIstioGRPC {
			return errors.Errorf("APIType support GRPC/ISTIO only but get %s", config.APITypeStr)
		}
	}
	return nil
}
