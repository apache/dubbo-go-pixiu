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

import type { FieldConfig } from '../components/DualModeEditor';
import {
  RATE_LIMIT_CONTROL_BEHAVIORS,
  RATE_LIMIT_METRIC_TYPES,
  CIRCUIT_BREAKER_STRATEGIES,
} from '../types/forms';

export type PluginCategory = 'http' | 'network' | 'rpc' | 'other';

export interface PluginConfigDefinition {
  type: string;
  category: PluginCategory;
  labelKey: string;
  basicFields: FieldConfig[];
  advancedFields?: FieldConfig[];
  defaultValues: Record<string, unknown>;
}

// Helper to create control behavior options
const controlBehaviorOptions = RATE_LIMIT_CONTROL_BEHAVIORS.map((item) => ({
  label: `rateLimitControlBehaviors.${item.key}`,
  value: item.value,
}));

// Helper to create metric type options
const metricTypeOptions = RATE_LIMIT_METRIC_TYPES.map((item) => ({
  label: `rateLimitMetricTypes.${item.key}`,
  value: item.value,
}));

// Helper to create circuit breaker strategy options
const circuitBreakerStrategyOptions = CIRCUIT_BREAKER_STRATEGIES.map((item) => ({
  label: `circuitBreakerStrategies.${item.key}`,
  value: item.value,
}));

// Plugin configurations for common plugins with wizard support
// Field names use snake_case to match Pixiu's YAML config format
export const PLUGIN_CONFIGS: Record<string, PluginConfigDefinition> = {
  // Rate Limit Plugin (Sentinel-based)
  'dgp.filter.http.ratelimit': {
    type: 'dgp.filter.http.ratelimit',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.ratelimit',
    basicFields: [
      {
        name: 'threshold',
        label: 'pluginConfig.ratelimit.threshold',
        type: 'number',
        required: true,
        min: 1,
        placeholder: 'pluginConfig.ratelimit.thresholdPlaceholder',
      },
      {
        name: 'stat_interval_in_ms',
        label: 'pluginConfig.ratelimit.statInterval',
        type: 'number',
        required: false,
        min: 100,
        placeholder: 'pluginConfig.ratelimit.statIntervalPlaceholder',
      },
      {
        name: 'metric_type',
        label: 'pluginConfig.ratelimit.metricType',
        type: 'select',
        required: true,
        options: metricTypeOptions,
      },
      {
        name: 'control_behavior',
        label: 'pluginConfig.ratelimit.controlBehavior',
        type: 'select',
        required: true,
        options: controlBehaviorOptions,
      },
    ],
    advancedFields: [
      {
        name: 'warm_up_period_sec',
        label: 'pluginConfig.ratelimit.warmUpPeriod',
        type: 'number',
        min: 0,
        placeholder: 'pluginConfig.ratelimit.warmUpPeriodPlaceholder',
      },
      {
        name: 'max_queuing_timeout_ms',
        label: 'pluginConfig.ratelimit.maxQueueingTimeout',
        type: 'number',
        min: 0,
        placeholder: 'pluginConfig.ratelimit.maxQueueingTimeoutPlaceholder',
      },
    ],
    defaultValues: {
      threshold: 100,
      stat_interval_in_ms: 1000,
      metric_type: 0,
      control_behavior: 0,
    },
  },

  // Circuit Breaker Plugin (Sentinel-based)
  'dgp.filter.http.circuitbreaker': {
    type: 'dgp.filter.http.circuitbreaker',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.circuitbreaker',
    basicFields: [
      {
        name: 'strategy',
        label: 'pluginConfig.circuitbreaker.strategy',
        type: 'select',
        required: true,
        options: circuitBreakerStrategyOptions,
      },
      {
        name: 'threshold',
        label: 'pluginConfig.circuitbreaker.threshold',
        type: 'number',
        required: true,
        min: 1,
        placeholder: 'pluginConfig.circuitbreaker.thresholdPlaceholder',
      },
      {
        name: 'retry_timeout_ms',
        label: 'pluginConfig.circuitbreaker.retryTimeout',
        type: 'number',
        required: true,
        min: 1000,
        placeholder: 'pluginConfig.circuitbreaker.retryTimeoutPlaceholder',
      },
      {
        name: 'stat_interval_ms',
        label: 'pluginConfig.circuitbreaker.statInterval',
        type: 'number',
        required: false,
        min: 100,
        placeholder: 'pluginConfig.circuitbreaker.statIntervalPlaceholder',
      },
    ],
    advancedFields: [
      {
        name: 'min_request_amount',
        label: 'pluginConfig.circuitbreaker.minRequestAmount',
        type: 'number',
        min: 1,
        placeholder: 'pluginConfig.circuitbreaker.minRequestAmountPlaceholder',
      },
      {
        name: 'slow_request_ratio_threshold',
        label: 'pluginConfig.circuitbreaker.slowRatioThreshold',
        type: 'number',
        min: 0,
        max: 1,
        placeholder: 'pluginConfig.circuitbreaker.slowRatioThresholdPlaceholder',
      },
      {
        name: 'max_allowed_rt_ms',
        label: 'pluginConfig.circuitbreaker.maxAllowedRt',
        type: 'number',
        min: 1,
        placeholder: 'pluginConfig.circuitbreaker.maxAllowedRtPlaceholder',
      },
    ],
    defaultValues: {
      strategy: 0,
      threshold: 5,
      retry_timeout_ms: 3000,
      stat_interval_ms: 10000,
      min_request_amount: 5,
    },
  },

  // JWT Authentication Plugin
  // Note: Pixiu JWT uses complex nested config with rules/providers
  // This provides simplified form fields; use YAML editor for advanced config
  'dgp.filter.http.auth.jwt': {
    type: 'dgp.filter.http.auth.jwt',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.auth.jwt',
    basicFields: [
      {
        name: 'issuer',
        label: 'pluginConfig.jwt.issuer',
        type: 'input',
        required: true,
        placeholder: 'pluginConfig.jwt.issuerPlaceholder',
      },
      {
        name: 'jwks_uri',
        label: 'pluginConfig.jwt.jwksUri',
        type: 'input',
        required: false,
        placeholder: 'pluginConfig.jwt.jwksUriPlaceholder',
      },
      {
        name: 'from_headers',
        label: 'pluginConfig.jwt.fromHeaders',
        type: 'input',
        required: false,
        placeholder: 'pluginConfig.jwt.fromHeadersPlaceholder',
      },
    ],
    advancedFields: [
      {
        name: 'forward',
        label: 'pluginConfig.jwt.forward',
        type: 'select',
        options: [
          { label: 'common.yes', value: true },
          { label: 'common.no', value: false },
        ],
      },
      {
        name: 'from_params',
        label: 'pluginConfig.jwt.fromParams',
        type: 'input',
        placeholder: 'pluginConfig.jwt.fromParamsPlaceholder',
      },
    ],
    defaultValues: {
      from_headers: 'Authorization',
      forward: true,
    },
  },

  // CORS Plugin
  // Pixiu uses: allow_origin ([]string), allow_methods (string), etc.
  'dgp.filter.http.cors': {
    type: 'dgp.filter.http.cors',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.cors',
    basicFields: [
      {
        name: 'allow_origin',
        label: 'pluginConfig.cors.allowOrigin',
        type: 'textarea',
        required: true,
        placeholder: 'pluginConfig.cors.allowOriginPlaceholder',
      },
      {
        name: 'allow_methods',
        label: 'pluginConfig.cors.allowMethods',
        type: 'input',
        required: false,
        placeholder: 'pluginConfig.cors.allowMethodsPlaceholder',
      },
      {
        name: 'allow_headers',
        label: 'pluginConfig.cors.allowHeaders',
        type: 'input',
        required: false,
        placeholder: 'pluginConfig.cors.allowHeadersPlaceholder',
      },
    ],
    advancedFields: [
      {
        name: 'expose_headers',
        label: 'pluginConfig.cors.exposeHeaders',
        type: 'input',
        placeholder: 'pluginConfig.cors.exposeHeadersPlaceholder',
      },
      {
        name: 'allow_credentials',
        label: 'pluginConfig.cors.allowCredentials',
        type: 'select',
        options: [
          { label: 'common.yes', value: true },
          { label: 'common.no', value: false },
        ],
      },
      {
        name: 'max_age',
        label: 'pluginConfig.cors.maxAge',
        type: 'number',
        min: 0,
        placeholder: 'pluginConfig.cors.maxAgePlaceholder',
      },
    ],
    defaultValues: {
      allow_origin: '*',
      allow_methods: 'GET,POST,PUT,DELETE,OPTIONS',
      allow_headers: 'Content-Type,Authorization',
      allow_credentials: false,
      max_age: 86400,
    },
  },

  // Timeout Plugin
  'dgp.filter.http.timeout': {
    type: 'dgp.filter.http.timeout',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.timeout',
    basicFields: [
      {
        name: 'timeout',
        label: 'pluginConfig.timeout.timeout',
        type: 'input',
        required: true,
        placeholder: 'pluginConfig.timeout.timeoutPlaceholder',
      },
    ],
    defaultValues: {
      timeout: '30s',
    },
  },

  // HTTP Proxy Plugin
  // Pixiu uses: timeout, max_idle_conns, scheme
  'dgp.filter.http.httpproxy': {
    type: 'dgp.filter.http.httpproxy',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.httpproxy',
    basicFields: [
      {
        name: 'timeout',
        label: 'pluginConfig.proxy.timeout',
        type: 'input',
        required: false,
        placeholder: 'pluginConfig.proxy.timeoutPlaceholder',
      },
      {
        name: 'scheme',
        label: 'pluginConfig.proxy.scheme',
        type: 'select',
        required: false,
        options: [
          { label: 'HTTP', value: 'http' },
          { label: 'HTTPS', value: 'https' },
        ],
      },
    ],
    advancedFields: [
      {
        name: 'max_idle_conns',
        label: 'pluginConfig.proxy.maxIdleConns',
        type: 'number',
        min: 1,
        placeholder: 'pluginConfig.proxy.maxIdleConnsPlaceholder',
      },
    ],
    defaultValues: {
      timeout: '30s',
      scheme: 'http',
    },
  },

  // Dubbo Proxy Plugin
  // Note: Pixiu Dubbo config uses complex nested structure
  'dgp.filter.http.dubboproxy': {
    type: 'dgp.filter.http.dubboproxy',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.dubboproxy',
    basicFields: [
      {
        name: 'registry_protocol',
        label: 'pluginConfig.dubbo.registryProtocol',
        type: 'select',
        required: true,
        options: [
          { label: 'Zookeeper', value: 'zookeeper' },
          { label: 'Nacos', value: 'nacos' },
          { label: 'Consul', value: 'consul' },
        ],
      },
      {
        name: 'registry_address',
        label: 'pluginConfig.dubbo.registryAddress',
        type: 'input',
        required: true,
        placeholder: 'pluginConfig.dubbo.registryAddressPlaceholder',
      },
    ],
    advancedFields: [
      {
        name: 'timeout',
        label: 'pluginConfig.timeout.timeout',
        type: 'input',
        placeholder: 'pluginConfig.timeout.timeoutPlaceholder',
      },
      {
        name: 'retries',
        label: 'pluginConfig.dubbo.retries',
        type: 'number',
        min: 0,
        placeholder: 'pluginConfig.dubbo.retriesPlaceholder',
      },
    ],
    defaultValues: {
      registry_protocol: 'zookeeper',
      timeout: '30s',
      retries: 2,
    },
  },

  // gRPC Proxy Plugin
  // Pixiu uses: timeout
  'dgp.filter.http.grpcproxy': {
    type: 'dgp.filter.http.grpcproxy',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.grpcproxy',
    basicFields: [
      {
        name: 'timeout',
        label: 'pluginConfig.timeout.timeout',
        type: 'input',
        required: false,
        placeholder: 'pluginConfig.timeout.timeoutPlaceholder',
      },
    ],
    advancedFields: [],
    defaultValues: {
      timeout: '30s',
    },
  },

  // Access Log Plugin
  'dgp.filter.http.accesslog': {
    type: 'dgp.filter.http.accesslog',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.accesslog',
    basicFields: [
      {
        name: 'path',
        label: 'pluginConfig.accesslog.path',
        type: 'input',
        required: true,
        placeholder: 'pluginConfig.accesslog.pathPlaceholder',
      },
    ],
    advancedFields: [
      {
        name: 'format',
        label: 'pluginConfig.accesslog.format',
        type: 'textarea',
        placeholder: 'pluginConfig.accesslog.formatPlaceholder',
      },
    ],
    defaultValues: {
      path: '/var/log/pixiu/access.log',
    },
  },

  // Metrics Plugin
  // Pixiu uses: mode (pull/push), push_config (nested)
  'dgp.filter.http.metric': {
    type: 'dgp.filter.http.metric',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.metric',
    basicFields: [
      {
        name: 'mode',
        label: 'pluginConfig.metric.mode',
        type: 'select',
        required: true,
        options: [
          { label: 'Pull', value: 'pull' },
          { label: 'Push', value: 'push' },
        ],
      },
    ],
    advancedFields: [
      {
        name: 'push_gateway_url',
        label: 'pluginConfig.metric.pushGatewayUrl',
        type: 'input',
        placeholder: 'pluginConfig.metric.pushGatewayUrlPlaceholder',
      },
      {
        name: 'push_interval_seconds',
        label: 'pluginConfig.metric.pushIntervalSeconds',
        type: 'number',
        min: 1,
        placeholder: 'pluginConfig.metric.pushIntervalSecondsPlaceholder',
      },
    ],
    defaultValues: {
      mode: 'pull',
    },
  },

  // Prometheus Metrics Plugin
  'dgp.filter.http.prometheusmetric': {
    type: 'dgp.filter.http.prometheusmetric',
    category: 'http',
    labelKey: 'plugins.dgp.filter.http.prometheusmetric',
    basicFields: [
      {
        name: 'path',
        label: 'pluginConfig.prometheus.path',
        type: 'input',
        required: true,
        placeholder: 'pluginConfig.prometheus.pathPlaceholder',
      },
    ],
    defaultValues: {
      path: '/metrics',
    },
  },
};

// Get plugins grouped by category
export function getPluginsByCategory(): Record<PluginCategory, PluginConfigDefinition[]> {
  const result: Record<PluginCategory, PluginConfigDefinition[]> = {
    http: [],
    network: [],
    rpc: [],
    other: [],
  };

  Object.values(PLUGIN_CONFIGS).forEach((config) => {
    result[config.category].push(config);
  });

  return result;
}

// Get plugin config by type
export function getPluginConfig(type: string): PluginConfigDefinition | undefined {
  return PLUGIN_CONFIGS[type];
}

// Check if a plugin has wizard support
export function hasWizardSupport(type: string): boolean {
  return type in PLUGIN_CONFIGS;
}
