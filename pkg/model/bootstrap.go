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

package model

// Bootstrap the door
type Bootstrap struct {
	StaticResources  StaticResources  `yaml:"static_resources" json:"static_resources" mapstructure:"static_resources"`
	DynamicResources DynamicResources `yaml:"dynamic_resources" json:"dynamic_resources" mapstructure:"dynamic_resources"`
	Tracing          Tracing          `yaml:"tracing" json:"tracing" mapstructure:"tracing"`
}

// GetListeners
func (bs *Bootstrap) GetListeners() []Listener {
	return bs.StaticResources.Listeners
}

// GetPprof
func (bs *Bootstrap) GetPprof() PprofConf {
	return bs.StaticResources.PprofConf
}

// GetAPIMetaConfig get api meta config from bootstrap
func (bs *Bootstrap) GetAPIMetaConfig() *APIMetaConfig {
	return bs.StaticResources.APIMetaConfig
}

// ExistCluster
func (bs *Bootstrap) ExistCluster(name string) bool {
	if len(bs.StaticResources.Clusters) > 0 {
		for _, v := range bs.StaticResources.Clusters {
			if v.Name == name {
				return true
			}
		}
	}

	return false
}

// StaticResources
type StaticResources struct {
	Listeners      []Listener      `yaml:"listeners" json:"listeners" mapstructure:"listeners"`
	Clusters       []Cluster       `yaml:"clusters" json:"clusters" mapstructure:"clusters"`
	ShutdownConfig *ShutdownConfig `yaml:"shutdown_config" json:"shutdown_config" mapstructure:"shutdown_config"`
	PprofConf      PprofConf       `yaml:"pprofConf" json:"pprofConf" mapstructure:"pprofConf"`
	APIMetaConfig  *APIMetaConfig  `yaml:"api_meta_config" json:"api_meta_config,omitempty"`
}

// DynamicResources TODO
type DynamicResources struct {
}

// ShutdownConfig how to shutdown proxy.
type ShutdownConfig struct {
	Timeout      string `default:"60s" yaml:"timeout" json:"timeout,omitempty"`
	StepTimeout  string `default:"10s" yaml:"step_timeout" json:"step_timeout,omitempty"`
	RejectPolicy string `yaml:"reject_policy" json:"reject_policy,omitempty"`
}

// APIMetaConfig how to find api config, file or etcd etc.
type APIMetaConfig struct {
	Address string `yaml:"address" json:"address,omitempty"`
	APIConfigPath string `default:"/proxy/config/api" yaml:"api_config_path" json:"api_config_path,omitempty" mapstructure:"api_config_path"`
}
