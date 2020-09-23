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
package api_load

import "github.com/dubbogo/dubbo-go-proxy/pkg/model"

// TODO
func init() {
	var _ ApiLoader = new(NacosApiLoader)
}

type NacosApiLoader struct {
	NacosAddress string
	ApiConfigs   []model.Api
	Prior        int
}

type NacosApiLoaderOption func(*NacosApiLoader)

func WithNacosAddress(nacosAddress string) NacosApiLoaderOption {
	return func(opt *NacosApiLoader) {
		opt.NacosAddress = nacosAddress
	}
}

func NewNacosApiLoader(opts ...NacosApiLoaderOption) *NacosApiLoader {
	var NacosApiLoader = &NacosApiLoader{}
	for _, opt := range opts {
		opt(NacosApiLoader)
	}
	return NacosApiLoader
}

func (f *NacosApiLoader) GetPrior() int {
	return f.Prior
}

func (f *NacosApiLoader) GetLoadedApiConfigs() ([]model.Api, error) {
	return f.ApiConfigs, nil
}

func (f *NacosApiLoader) InitLoad() (err error) {
	panic("")
}

func (f *NacosApiLoader) HotLoad() (chan struct{}, error) {
	panic("")
}

func (f *NacosApiLoader) Clear() error {
	panic("")
}
