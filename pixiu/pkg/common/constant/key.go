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
	HTTPConnectManagerFilter  = "dgp.filter.httpconnectionmanager"
	GRPCConnectManagerFilter  = "dgp.filter.grpcconnectionmanager"
	DubboConnectManagerFilter = "dgp.filter.network.dubboconnectionmanager"

	HTTPAuthorityFilter        = "dgp.filter.http.authority"
	HTTPProxyFilter            = "dgp.filter.http.httpproxy"
	HTTPHeaderFilter           = "dgp.filter.http.header"
	HTTPHostFilter             = "dgp.filter.http.host"
	HTTPMetricFilter           = "dgp.filter.http.metric"
	HTTPRecoveryFilter         = "dgp.filter.http.recovery"
	HTTPResponseFilter         = "dgp.filter.http.response"
	HTTPAccessLogFilter        = "dgp.filter.http.accesslog"
	HTTPRateLimitFilter        = "dgp.filter.http.ratelimit"
	HTTPGrpcProxyFilter        = "dgp.filter.http.grpcproxy"
	HTTPDubboProxyFilter       = "dgp.filter.http.dubboproxy"
	HTTPDirectDubboProxyFilter = "dgp.filter.http.directdubboproxy"
	HTTPApiConfigFilter        = "dgp.filter.http.apiconfig"
	HTTPTimeoutFilter          = "dgp.filter.http.timeout"
	TracingFilter              = "dgp.filters.tracing"
	HTTPWasmFilter             = "dgp.filter.http.webassembly"
	HTTPCircuitBreakerFilter   = "dgp.filter.http.circuitbreaker"
	HTTPAuthJwtFilter          = "dgp.filter.http.auth.jwt"
	HTTPCorsFilter             = "dgp.filter.http.cors"
	HTTPCsrfFilter             = "dgp.filter.http.csrf"
	HTTPProxyRewriteFilter     = "dgp.filter.http.proxyrewrite"
	HTTPLoadBalanceFilter      = "dgp.filter.http.loadbalance"
	HTTPEventFilter            = "dgp.filter.http.event"
	HTTPTrafficFilter          = "dgp.filter.http.traffic"
	HTTPPrometheusMetricFilter = "dgp.filter.http.prometheusmetric"

	DubboHttpFilter  = "dgp.filter.dubbo.http"
	DubboProxyFilter = "dgp.filter.dubbo.proxy"
)

const (
	SpringCloudAdapter         = "dgp.adapter.springcloud"
	DubboRegistryCenterAdapter = "dgp.adapter.dubboregistrycenter"
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
	MQTypeKafka    = "kafka"
	MQTypeRocketMQ = "rocketmq"
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
	// MetadataStorageTypeKey the storage type of metadata
	MetadataStorageTypeKey     = "dubbo.metadata.storage-type"
	DefaultMetadataStorageType = "local"
)
