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

package converter

type PixiuConfig struct {
	StaticResources  StaticResources   `yaml:"static_resources" json:"static_resources"`
	DynamicResources *DynamicResources `yaml:"dynamic_resources,omitempty" json:"dynamic_resources,omitempty"`
	ShutdownConfig   ShutdownConfig    `yaml:"shutdown_config" json:"shutdown_config"`
	Log              LogConfig         `yaml:"log" json:"log"`
}

type StaticResources struct {
	Listeners []*Listener `yaml:"listeners" json:"listeners"`
	Clusters  []*Cluster  `yaml:"clusters" json:"clusters"`
}

type Listener struct {
	Name         string      `yaml:"name" json:"name"`
	ProtocolType string      `yaml:"protocol_type" json:"protocol_type"`
	Address      Address     `yaml:"address" json:"address"`
	FilterChain  FilterChain `yaml:"filter_chains" json:"filter_chains"`
	Config       interface{} `yaml:"config,omitempty" json:"config,omitempty"`
}

type Address struct {
	SocketAddress SocketAddress `yaml:"socket_address" json:"socket_address"`
}

type SocketAddress struct {
	Address string `yaml:"address" json:"address"`
	Port    int    `yaml:"port" json:"port"`
}

type FilterChain struct {
	Filters []NetworkFilter `yaml:"filters" json:"filters"`
}

type NetworkFilter struct {
	Name   string      `yaml:"name" json:"name"`
	Config interface{} `yaml:"config" json:"config"`
}

type HTTPConnectionManagerConfig struct {
	RouteConfig RouteConfiguration `yaml:"route_config" json:"route_config"`
	HTTPFilters []HTTPFilter       `yaml:"http_filters" json:"http_filters"`
}

type RouteConfiguration struct {
	Routes []*Route `yaml:"routes" json:"routes"`
}

type Route struct {
	Match RouteMatch  `yaml:"match" json:"match"`
	Route RouteAction `yaml:"route" json:"route"`
}

type RouteMatch struct {
	Prefix  string          `yaml:"prefix,omitempty" json:"prefix,omitempty"`
	Path    string          `yaml:"path,omitempty" json:"path,omitempty"`
	Methods []string        `yaml:"methods,omitempty" json:"methods,omitempty"`
	Headers []HeaderMatcher `yaml:"headers,omitempty" json:"headers,omitempty"`
}

type HeaderMatcher struct {
	Name   string   `yaml:"name" json:"name"`
	Values []string `yaml:"values,omitempty" json:"values,omitempty"`
	Regex  bool     `yaml:"regex,omitempty" json:"regex,omitempty"`
}

type RouteAction struct {
	Cluster                     string `yaml:"cluster" json:"cluster"`
	ClusterNotFoundResponseCode int    `yaml:"cluster_not_found_response_code,omitempty" json:"cluster_not_found_response_code,omitempty"`
}

type HTTPFilter struct {
	Name   string                 `yaml:"name" json:"name"`
	Config map[string]interface{} `yaml:"config,omitempty" json:"config,omitempty"`
}

type Cluster struct {
	Name      string      `yaml:"name" json:"name"`
	Type      string      `yaml:"type" json:"type"`
	LbPolicy  string      `yaml:"lb_policy" json:"lb_policy"`
	Endpoints []*Endpoint `yaml:"endpoints" json:"endpoints"`
}

type Endpoint struct {
	ID            int           `yaml:"id" json:"id"`
	SocketAddress SocketAddress `yaml:"socket_address" json:"socket_address"`
}

type ShutdownConfig struct {
	Timeout      string `yaml:"timeout" json:"timeout"`
	StepTimeout  string `yaml:"step_timeout" json:"step_timeout"`
	RejectPolicy string `yaml:"reject_policy" json:"reject_policy"`
}

type LogConfig struct {
	Level string `yaml:"level" json:"level"`
}

type DynamicResources struct {
	LdsConfig *ApiConfigSource `yaml:"lds_config,omitempty" json:"lds_config,omitempty"`
	CdsConfig *ApiConfigSource `yaml:"cds_config,omitempty" json:"cds_config,omitempty"`
}

type ApiConfigSource struct {
	ClusterName    []string `yaml:"cluster_name" json:"cluster_name"`
	APIType        string   `yaml:"api_type" json:"api_type"`
	RefreshDelay   string   `yaml:"refresh_delay,omitempty" json:"refresh_delay,omitempty"`
	RequestTimeout string   `yaml:"request_timeout,omitempty" json:"request_timeout,omitempty"`
	GrpcServices   []struct {
		Timeout string `yaml:"timeout,omitempty" json:"timeout,omitempty"`
	} `yaml:"grpc_services,omitempty" json:"grpc_services,omitempty"`
}
