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

package constant

import (
	"time"
)

// default timeout 1s.
const (
	DefaultTimeoutStr = "1s"
	DefaultTimeout    = time.Second
)

const (
	// DefaultBodyAll body passthrough.
	DefaultBodyAll = "_all"
)

// strategy for response.
const (
	ResponseStrategyNormal = "normal"
	ResponseStrategyHump   = "hump"
)

const (
	// DefaultDiscoveryType Set up default discovery type.
	DefaultDiscoveryType = "EDS"
	// DefaultLoadBalanceType Set up default load balance type.
	DefaultLoadBalanceType = "RoundRobin"
	// DefaultFilterType Set up default filter type.
	DefaultFilterType = "dgp.filter.httpconnectionmanager"
	// DefaultHTTPType Set up default HTTP Type.
	DefaultHTTPType = "net/http"
	// DefaultProtocolType Set up default protocol type.
	DefaultProtocolType = "HTTP"
)

const (
	// YAML .yaml
	YAML = ".yaml"
	//YML .yml
	YML = ".yml"
)

const (
	StringSeparator      = ","
	DefaultConfigPath    = "configs/conf.yaml"
	DefaultApiConfigPath = "configs/api_config.yaml"
	DefaultLogConfigPath = "configs/log.yml"
	DefaultLogLevel      = "info"
	DefaultLimitCpus     = "0"
	DefaultLogFormat     = ""
)
