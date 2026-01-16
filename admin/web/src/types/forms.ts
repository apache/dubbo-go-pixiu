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

export interface ClusterFormData {
  id?: string;
  name: string;
  type: string;
  address: string;
  port: number;
}

export interface ListenerFormData {
  id?: string;
  name: string;
  protocol: string;
  address: string;
  port: number;
}

export interface ResourceFormData {
  id?: string;
  path: string;
  type: string;
  description?: string;
  timeout?: string;
}

export type FormDataType = ClusterFormData | ListenerFormData | ResourceFormData;

export type ResourceType = 'cluster' | 'listener' | 'resource';

// Cluster discovery types supported by Pixiu (based on Envoy cluster types)
export const CLUSTER_TYPES = [
  'Static',
  'StrictDNS',
  'LogicalDns',
  'EDS',
  'OriginalDst',
] as const;

// Listener protocol types supported by Pixiu
export const LISTENER_PROTOCOLS = [
  'HTTP',
  'HTTPS',
  'HTTP2',
  'TCP',
  'UDP',
  'GRPC',
  'TRIPLE',
] as const;

// Resource/Mapping types supported by Pixiu
export const RESOURCE_TYPES = ['restful', 'dubbo'] as const;

// HTTP Filter plugins supported by Pixiu
export const HTTP_FILTER_PLUGINS = [
  'dgp.filter.httpconnectionmanager',
  'dgp.filter.http.accesslog',
  'dgp.filter.http.apiconfig',
  'dgp.filter.http.auth.jwt',
  'dgp.filter.http.auth.mcp',
  'dgp.filter.http.authority',
  'dgp.filter.http.circuitbreaker',
  'dgp.filter.http.cors',
  'dgp.filter.http.csrf',
  'dgp.filter.http.directdubboproxy',
  'dgp.filter.http.dubboproxy',
  'dgp.filter.http.event',
  'dgp.filter.http.faultinjection',
  'dgp.filter.http.grpcproxy',
  'dgp.filter.http.header',
  'dgp.filter.http.host',
  'dgp.filter.http.httpproxy',
  'dgp.filter.http.loadbalance',
  'dgp.filter.http.metric',
  'dgp.filter.http.opa',
  'dgp.filter.http.prometheusmetric',
  'dgp.filter.http.proxyrewrite',
  'dgp.filter.http.ratelimit',
  'dgp.filter.http.recovery',
  'dgp.filter.http.response',
  'dgp.filter.http.timeout',
  'dgp.filter.http.traffic',
  'dgp.filter.http.webassembly',
] as const;

// Network Filter plugins supported by Pixiu
export const NETWORK_FILTER_PLUGINS = [
  'dgp.filter.grpcconnectionmanager',
  'dgp.filter.network.grpcconnectionmanager',
  'dgp.filter.network.dubboconnectionmanager',
] as const;

// RPC Filter plugins supported by Pixiu
export const RPC_FILTER_PLUGINS = [
  'dgp.filter.dubbo.http',
  'dgp.filter.dubbo.proxy',
  'dgp.filter.grpc.proxy',
  'dgp.filters.tracing',
] as const;

// Other plugins
export const OTHER_PLUGINS = [
  'dgp.filter.llm.proxy',
  'dgp.filter.llm.tokenizer',
  'dgp.filter.mcp.mcpserver',
] as const;

// All plugin types
export const PLUGIN_TYPES = [
  ...HTTP_FILTER_PLUGINS,
  ...NETWORK_FILTER_PLUGINS,
  ...RPC_FILTER_PLUGINS,
  ...OTHER_PLUGINS,
] as const;

// Rate limit matching strategies (Sentinel-based)
// 0=EXACT, 1=REGEX, 2=AntPath
export const RATE_LIMIT_MATCH_STRATEGIES = [
  { value: 0, key: 'EXACT' },
  { value: 1, key: 'REGEX' },
  { value: 2, key: 'AntPath' },
] as const;

// Rate limit flow control behaviors (Sentinel-based)
// 0=Reject, 1=Throttle, 2=WarmUp, 3=WarmUpThrottling
export const RATE_LIMIT_CONTROL_BEHAVIORS = [
  { value: 0, key: 'Reject' },
  { value: 1, key: 'Throttle' },
  { value: 2, key: 'WarmUp' },
  { value: 3, key: 'WarmUpThrottling' },
] as const;

// Rate limit metric types
// 0=QPS, 1=Threads
export const RATE_LIMIT_METRIC_TYPES = [
  { value: 0, key: 'QPS' },
  { value: 1, key: 'Threads' },
] as const;

// Circuit breaker strategies (Sentinel-based)
// 0=ErrorCount, 1=ErrorRatio, 2=SlowRequestRatio
export const CIRCUIT_BREAKER_STRATEGIES = [
  { value: 0, key: 'ErrorCount' },
  { value: 1, key: 'ErrorRatio' },
  { value: 2, key: 'SlowRequestRatio' },
] as const;
