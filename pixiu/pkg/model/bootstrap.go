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

import (
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

// Bootstrap the door
type Bootstrap struct {
	StaticResources  StaticResources   `yaml:"static_resources" json:"static_resources" mapstructure:"static_resources"`
	DynamicResources *DynamicResources `yaml:"dynamic_resources" json:"dynamic_resources" mapstructure:"dynamic_resources"`
	Metric           Metric            `yaml:"metric" json:"metric" mapstructure:"metric"`
	Node             *Node             `yaml:"node" json:"node" mapstructure:"node"`
	Trace            *TracerConfig     `yaml:"tracing" json:"tracing" mapstructure:"tracing"`
	Wasm             *WasmConfig       `yaml:"wasm" json:"wasm" mapstructure:"wasm"`
	Config           *ConfigCenter     `yaml:"config-center" json:"config-center" mapstructure:"config-center"`
	// Third party dependency
	Nacos *Nacos `yaml:"nacos" json:"nacos" mapstructure:"nacos"`
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

// GetShutdownConfig
func (bs *Bootstrap) GetShutdownConfig() *ShutdownConfig {
	if bs.StaticResources.ShutdownConfig == nil {
		bs.StaticResources.ShutdownConfig = &ShutdownConfig{
			Timeout:      "0s",
			StepTimeout:  "0s",
			RejectPolicy: "immediacy",
		}
	}
	return bs.StaticResources.ShutdownConfig
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
	Timeout      string `default:"0s" yaml:"timeout" json:"timeout,omitempty"`
	StepTimeout  string `default:"0s" yaml:"step_timeout" json:"step_timeout,omitempty"`
	RejectPolicy string `default:"immediacy" yaml:"reject_policy" json:"reject_policy,omitempty"`
}

// GetTimeoutOfShutdown
func (sdc *ShutdownConfig) GetTimeout() time.Duration {
	result, err := time.ParseDuration(sdc.Timeout)
	if err != nil {
		defaultTimeout := 60 * time.Second
		logger.Errorf("The Timeout configuration is invalid: %s, and we will use the default value: %s, err: %v",
			sdc.Timeout, defaultTimeout.String(), err)
		return defaultTimeout
	}
	return result
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

type Config struct {
	Listeners      []*Listener      `yaml:"listeners" json:"listeners" mapstructure:"listeners"`
	Clusters       []*ClusterConfig `yaml:"clusters" json:"clusters" mapstructure:"clusters"`
	Adapters       []*Adapter       `yaml:"adapters" json:"adapters" mapstructure:"adapters"`
	ShutdownConfig *ShutdownConfig  `yaml:"shutdown_config" json:"shutdown_config" mapstructure:"shutdown_config"`
	PprofConf      PprofConf        `yaml:"pprofConf" json:"pprofConf" mapstructure:"pprofConf"`
}

type Nacos struct {
	ServerConfigs []*NacosServerConfig `yaml:"server_configs" json:"server-configs" mapstructure:"server_configs"`
	ClientConfig  *NacosClientConfig   `yaml:"client-config" json:"client-config" mapstructure:"client_config"`
	DataId        string               `yaml:"data-id" json:"data-id" mapstructure:"data_id" default:"pixiu.yaml"`
	Group         string               `yaml:"group" json:"group" mapstructure:"group" default:"DEFAULT_GROUP"`
}

type NacosServerConfig struct {
	IpAddr      string `json:"ip_addr,omitempty" yaml:"ip_addr" mapstructure:"ip_addr"`
	Port        uint64 `json:"port,omitempty" yaml:"port" mapstructure:"port"`
	Scheme      string `json:"scheme" yaml:"scheme" mapstructure:"scheme"`
	ContextPath string `json:"contextPath" yaml:"contextPath" mapstructure:"contextPath"`
}

type NacosClientConfig struct {
	TimeoutMs            uint64 `json:"timeout_ms,omitempty" yaml:"timeout_ms" mapstructure:"timeout_ms"`                                        // timeout for requesting Nacos server, default value is 10000ms
	ListenInterval       uint64 `json:"listen_interval,omitempty" yaml:"listen_interval" mapstructure:"listen_interval"`                         // Deprecated
	BeatInterval         int64  `json:"beat_interval,omitempty" yaml:"beat_interval" mapstructure:"beat_interval"`                               // the time interval for sending beat to server,default value is 5000ms
	NamespaceId          string `json:"namespace_id,omitempty" yaml:"namespace_id" mapstructure:"namespace_id"`                                  // the namespaceId of Nacos.When namespace is public, fill in the blank string here.
	AppName              string `json:"app_name,omitempty" yaml:"app_name" mapstructure:"app_name"`                                              // the appName
	Endpoint             string `json:"endpoint,omitempty" yaml:"endpoint" mapstructure:"endpoint"`                                              // the endpoint for get Nacos server addresses
	RegionId             string `json:"region_id,omitempty" yaml:"region_id" mapstructure:"region_id"`                                           // the regionId for kms
	AccessKey            string `json:"access_key,omitempty" yaml:"access_key" mapstructure:"access_key"`                                        // the AccessKey for kms
	SecretKey            string `json:"secret_key,omitempty" yaml:"secret_key" mapstructure:"secret_key"`                                        // the SecretKey for kms
	OpenKMS              bool   `json:"open_kms,omitempty" yaml:"open_kms" mapstructure:"open_kms"`                                              // it's to open kms,default is false. https://help.aliyun.com/product/28933.html
	CacheDir             string `json:"cache_dir,omitempty" yaml:"cache_dir" mapstructure:"cache_dir" default:"/tmp/nacos/cache"`                // the directory for persist nacos service info,default value is current path
	UpdateThreadNum      int    `json:"update_thread_num,omitempty" yaml:"update_thread_num" mapstructure:"update_thread_num"`                   // the number of gorutine for update nacos service info,default value is 20
	NotLoadCacheAtStart  bool   `json:"not_load_cache_at_start,omitempty" yaml:"not_load_cache_at_start" mapstructure:"not_load_cache_at_start"` // not to load persistent nacos service info in CacheDir at start time
	UpdateCacheWhenEmpty bool   `json:"update_cache_when_empty,omitempty" yaml:"update_cache_when_empty" mapstructure:"update_cache_when_empty"` // update cache when get empty service instance from server
	Username             string `json:"username,omitempty" yaml:"username" mapstructure:"username"`                                              // the username for nacos auth
	Password             string `json:"password,omitempty" yaml:"password" mapstructure:"password"`                                              // the password for nacos auth
	LogDir               string `json:"log_dir,omitempty" yaml:"log_dir" mapstructure:"log_dir" default:"/tmp/nacos/log"`                        // the directory for log, default is current path
	LogLevel             string `json:"log_level,omitempty" yaml:"log_level" mapstructure:"log_level"`                                           // the level of log, it's must be debug,info,warn,error, default value is info
	//LogSampling          *ClientLogSamplingConfig // the sampling config of log
	ContextPath string `json:"context_path,omitempty" yaml:"context_path" mapstructure:"context_path"` // the nacos server contextpath
	//LogRollingConfig     *ClientLogRollingConfig  // the log rolling config
}

type ConfigCenter struct {
	Type   string `json:"type,omitempty" yaml:"type"`
	Enable string `json:"enable" yaml:"enable"`
}
