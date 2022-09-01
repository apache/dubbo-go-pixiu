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
	StaticResources  StaticResources   `yaml:"static_resources" json:"static_resources" mapstructure:"static_resources"`
	DynamicResources *DynamicResources `yaml:"dynamic_resources" json:"dynamic_resources" mapstructure:"dynamic_resources"`
	Metric           Metric            `yaml:"metric" json:"metric" mapstructure:"metric"`
	Node             *Node             `yaml:"node" json:"node" mapstructure:"node"`
	Trace            *TracerConfig     `yaml:"tracing" json:"tracing" mapstructure:"tracing"`
	Wasm             *WasmConfig       `yaml:"wasm" json:"wasm" mapstructure:"wasm"`
}

// Node node info for dynamic identifier
type Node struct {
	Cluster string `yaml:"cluster" json:"cluster" mapstructure:"cluster"`
	Id      string `yaml:"id" json:"id" mapstructure:"id"`
}

// GetListeners
func (bs *Bootstrap) GetListeners() []*Listener {
	return bs.StaticResources.Listeners
}

func (bs *Bootstrap) GetStaticListeners() []*Listener {
	return bs.StaticResources.Listeners
}

// GetPprof
func (bs *Bootstrap) GetPprof() PprofConf {
	return bs.StaticResources.PprofConf
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
	Listeners      []*Listener      `yaml:"listeners" json:"listeners" mapstructure:"listeners"`
	Clusters       []*ClusterConfig `yaml:"clusters" json:"clusters" mapstructure:"clusters"`
	Adapters       []*Adapter       `yaml:"adapters" json:"adapters" mapstructure:"adapters"`
	ShutdownConfig *ShutdownConfig  `yaml:"shutdown_config" json:"shutdown_config" mapstructure:"shutdown_config"`
	PprofConf      PprofConf        `yaml:"pprofConf" json:"pprofConf" mapstructure:"pprofConf"`
}

// DynamicResources config the dynamic resource source
//
//		"lds_config": "{...}", # config lister load source
//		"cds_config": "{...}", # config cluster load source
//		"ads_config": "{...}"
//	 "ada_config": "{...}" # config adaptor load source
type DynamicResources struct {
	LdsConfig *ApiConfigSource `yaml:"lds_config" json:"lds_config" mapstructure:"lds_config"`
	CdsConfig *ApiConfigSource `yaml:"cds_config" json:"cds_config" mapstructure:"cds_config"`
	AdsConfig *ApiConfigSource `yaml:"ads_config" json:"ads_config" mapstructure:"ads_config"`
}

// ShutdownConfig how to shutdown server.
type ShutdownConfig struct {
	Timeout      string `default:"60s" yaml:"timeout" json:"timeout,omitempty"`
	StepTimeout  string `default:"10s" yaml:"step_timeout" json:"step_timeout,omitempty"`
	RejectPolicy string `default:"immediacy" yaml:"reject_policy" json:"reject_policy,omitempty"`
}

// APIMetaConfig how to find api config, file or etcd etc.
type APIMetaConfig struct {
	Address       string `yaml:"address" json:"address,omitempty"`
	APIConfigPath string `default:"/pixiu/config/api" yaml:"api_config_path" json:"api_config_path,omitempty" mapstructure:"api_config_path"`
}

// TimeoutConfig the config of ConnectTimeout and RequestTimeout
type TimeoutConfig struct {
	ConnectTimeoutStr string `default:"5s" yaml:"connect_timeout" json:"connect_timeout,omitempty"` // ConnectTimeout timeout for connect to cluster node
	RequestTimeoutStr string `default:"10s" yaml:"request_timeout" json:"request_timeout,omitempty"`
}
