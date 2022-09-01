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

package apiconfig

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

// ApiConfigConfig the config for api_config filter
type ApiConfigConfig struct {
	APIMetaConfig  *model.APIMetaConfig `yaml:"api_meta_config" json:"api_meta_config,omitempty"`
	Path           string               `yaml:"path" json:"path,omitempty"`
	Dynamic        bool                 `yaml:"dynamic" json:"dynamic,omitempty"`
	DynamicAdapter string               `yaml:"dynamic_adapter" json:"dynamic_adapter,omitempty"`
}
