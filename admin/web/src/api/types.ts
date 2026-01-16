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

// ==================== API Response ====================

// API Response codes (backend returns number)
export const API_CODE = {
  SUCCESS: 0,
  ERROR: -1,
  TOKEN_EXPIRED: 401,
} as const;

// Base API response structure
export interface ApiResponse<T = unknown> {
  code: number;
  message?: string;
  msg?: string;
  data: T;
}

// ==================== Auth Types ====================

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  username: string;
  userId?: number;
}

export interface UserInfo {
  userId: number;
  username: string;
  nickname?: string;
  email?: string;
  phone?: string;
  avatar?: string;
  role?: string;
}

export interface ChangePasswordRequest {
  oldPassword: string;
  newPassword: string;
}

// ==================== Cluster Types ====================

// Cluster endpoint (matches Go: pkg/model.Endpoint)
export interface ClusterEndpoint {
  ID?: string;
  name?: string;
  socket_address?: SocketAddress;
  meta?: Record<string, string>;
}

// Matches Go: pkg/model/cluster.go - ClusterConfig struct
// Note: Cluster uses 'name' as identifier, no 'id' field
// Address and port are nested in endpoints[].socket_address
export interface Cluster {
  name: string;
  type?: string;              // cluster discovery type (Static, StrictDNS, etc.)
  lb_policy?: string;         // load balancer policy (RoundRobin, Rand, etc.)
  endpoints?: ClusterEndpoint[];  // cluster endpoints with socket addresses
  health_checks?: unknown[];  // health check configurations
}

// ==================== Listener Types ====================

// Socket address structure
export interface SocketAddress {
  address?: string;
  host?: string;  // alias for address
  port?: number;
}

// Listener address structure
export interface ListenerAddress {
  socket_address?: SocketAddress;
  name?: string;
}

// Route match configuration
export interface RouteMatch {
  prefix: string;
}

// Route target configuration
export interface RouteTarget {
  cluster: string;
  cluster_not_found_response_code?: number;
}

// Route configuration item
export interface RouteItem {
  match: RouteMatch;
  route: RouteTarget;
}

// Route configuration wrapper (inside httpconnectionmanager filter config)
export interface RouteConfig {
  routes: RouteItem[];
}

// Network filter configuration (matches pkg/model.NetworkFilter)
export interface NetworkFilter {
  name: string;
  config?: Record<string, unknown>;
}

// Filter chain configuration (matches pkg/model.FilterChain)
export interface FilterChain {
  filters: NetworkFilter[];
}

// Matches Go: pkg/model/listener.go - Listener struct (Pixiu native format)
// Uses filter_chains for filters and routes configuration
export interface Listener {
  name: string;
  protocol_type?: string;  // HTTP, HTTPS, GRPC, TCP, etc.
  address: ListenerAddress;
  filter_chains: FilterChain;
  config?: Record<string, unknown>;  // Listener-level config (timeouts, etc.)
}

// HTTP filter configuration (used inside httpconnectionmanager)
export interface HttpFilter {
  name: string;
  config?: Record<string, unknown>;
}

// ==================== Resource Types (Mapping) ====================

// Filter configuration
export interface Filter {
  name?: string;
  config?: Record<string, unknown>;
}

// Params definition
export interface Param {
  name: string;
  type: string;
  required: boolean;
}

// Body definition
export interface BodyDefinition {
  definitionName: string;
}

// Inbound request configuration
export interface InboundRequest {
  requestType?: string;
  headers?: Param[];
  queryStrings?: Param[];
  requestBody?: BodyDefinition[];
}

// Dubbo backend configuration
export interface DubboBackendConfig {
  clusterName?: string;
  applicationName?: string;
  protocol?: string;
  group?: string;
  version?: string;
  interface?: string;
  method?: string;
  retries?: string;
}

// HTTP backend configuration
export interface HttpBackendConfig {
  url?: string;
  host?: string;
  path?: string;
  schema?: string;
}

// Mapping parameter
export interface MappingParam {
  name?: string;
  mapTo?: string;
  mapType?: string;
}

// Integration request configuration
export interface IntegrationRequest {
  requestType?: string;
  dubboBackendConfig?: DubboBackendConfig;
  httpBackendConfig?: HttpBackendConfig;
  mappingParams?: MappingParam[];
}

// Matches Go: pkg/config/api_config.go - Method struct
export interface Method {
  id?: number;
  resourcePath?: string;
  enable?: boolean;
  timeout?: number | string;
  mock?: boolean;
  filters?: Filter[];
  httpVerb: string;
  inboundRequest?: InboundRequest;
  integrationRequest?: IntegrationRequest;
}

// Matches Go: pkg/config/api_config.go - Resource struct
export interface Resource {
  id?: number;
  type: string;           // Restful, Dubbo
  path: string;
  timeout?: number | string;
  description?: string;
  filters?: Filter[];
  methods?: Method[];
  resources?: Resource[]; // nested resources
  headers?: Record<string, string>;
}

// ==================== Plugin Types ====================

// Matches Go: pkg/config/api_config.go - Plugin struct
export interface Plugin {
  name: string;
  version: string;
  priority: number;
  externalLookupName?: string;
  config?: Record<string, unknown>;
}

// Matches Go: pkg/config/api_config.go - PluginGroup struct
export interface PluginGroup {
  groupName: string;
  plugins: Plugin[];
}

// Plugin group detail with YAML config
export interface PluginGroupDetail {
  groupName: string;
  plugins: Plugin[];
  yamlConfig?: string;
}

// ==================== API Config Types ====================

// Definition for complex JSON schema
export interface Definition {
  name: string;
  schema: string;
}

// Full API configuration
export interface APIConfig {
  name: string;
  description?: string;
  resources: Resource[];
  definitions?: Definition[];
}

// ==================== Helper Types ====================

// Base config for YAML content
export interface BaseConfig {
  content: string;
}

// Type aliases for backward compatibility
export type ClusterDetail = Cluster;
export type ListenerDetail = Listener;
export type ResourceDetail = Resource;
export type MethodDetail = Method;
export type PluginConfig = Plugin;
export type PluginGroupItem = PluginGroup;

// ==================== User Management Types ====================

// User list item (safe representation without password)
export interface UserListItem {
  id: number;
  username: string;
  role: number;
  enabled: boolean;
  dateCreated: string;
  dateUpdated: string;
}

// User list response with pagination
export interface UserListResponse {
  items: UserListItem[];
  total: number;
  page: number;
  pageSize: number;
}

// Update user request
export interface UpdateUserRequest {
  role?: number;
  enabled?: boolean;
}

// Role definition
export interface Role {
  id: number;
  role_name: string;
  description: string;
}

// Role names mapping
export const ROLE_NAMES: Record<number, string> = {
  1: 'admin',
  2: 'user',
};

// Permission definition
export interface Permission {
  id: number;
  resource: string;
  action: string;
}

// ==================== Instance Types ====================

export interface Instance {
  nodeId: string;
  address: string;
  cluster: string;
  version: string;
  metadata?: Record<string, string>;
  status: string;
  lastSeen: string;
  connectedAt: string;
  uptime: string;
}

export interface InstanceStats {
  total: number;
  connected: number;
  disconnected: number;
}
