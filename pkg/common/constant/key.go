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

const (
	HTTPConnectManagerFilter = "dgp.filter.httpconnectionmanager"

	HTTPAuthorityFilter  = "dgp.filter.http.authority"
	HTTPProxyFilter      = "dgp.filter.http.httpproxy"
	HTTPHeaderFilter     = "dgp.filter.http.header"
	HTTPHostFilter       = "dgp.filter.http.host"
	HTTPMetricFilter     = "dgp.filter.http.metric"
	HTTPRecoveryFilter   = "dgp.filter.http.recovery"
	HTTPResponseFilter   = "dgp.filter.http.response"
	HTTPAccessLogFilter  = "dgp.filter.http.accesslog"
	HTTPRateLimitFilter  = "dgp.filter.http.ratelimit"
	HTTPGrpcProxyFilter  = "dgp.filter.http.grpcproxy"
	HTTPDubboProxyFilter = "dgp.filter.http.dubboproxy"
	HTTPApiConfigFilter  = "dgp.filter.http.apiconfig"
	HTTPTimeoutFilter    = "dgp.filter.http.timeout"
)

const (
	SpringCloudAdapter = "dgp.adapter.springcloud"
)

const (
	ConfigPathKey    = "config"
	ApiConfigPathKey = "api-config"
	LogConfigPathKey = "log-config"
	LogLevelKey      = "log-level"
	LimitCpusKey     = "limit-cpus"
	LogFormatKey     = "log-format"
)

const (
	ApplicationKey = "application"
	AppVersionKey  = "app.version"
	ClusterKey     = "cluster"
	GroupKey       = "group"
	VersionKey     = "version"
	InterfaceKey   = "interface"
	MethodsKey     = "methods"
	// NameKey name of interface
	NameKey = "name"
	// RetriesKey retry times
	RetriesKey = "retries"
)
