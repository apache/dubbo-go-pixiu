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

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

import (
	"controllers/api/v1alpha1"

	"controllers/internal/ir"

	corev1 "k8s.io/api/core/v1"

	discoveryv1 "k8s.io/api/discovery/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ApplyGatewayPolicy applies PixiuGatewayPolicy to PixiuConfig
func ApplyGatewayPolicy(config *PixiuConfig, policy *v1alpha1.PixiuGatewayPolicy) {
	if policy == nil {
		return
	}

	// Apply listener timeout
	if policy.Spec.Listener != nil && policy.Spec.Listener.Timeout != nil {
		for _, listener := range config.StaticResources.Listeners {
			if listener.Config == nil {
				listener.Config = make(map[string]interface{})
			}
			configMap := listener.Config.(map[string]interface{})
			if policy.Spec.Listener.Timeout.Idle != "" {
				configMap["idle_timeout"] = policy.Spec.Listener.Timeout.Idle
			}
			if policy.Spec.Listener.Timeout.Read != "" {
				configMap["read_timeout"] = policy.Spec.Listener.Timeout.Read
			}
			if policy.Spec.Listener.Timeout.Write != "" {
				configMap["write_timeout"] = policy.Spec.Listener.Timeout.Write
			}
		}
	}

	// Apply shutdown config
	if policy.Spec.Shutdown != nil {
		if policy.Spec.Shutdown.Timeout != "" {
			config.ShutdownConfig.Timeout = policy.Spec.Shutdown.Timeout
		}
		if policy.Spec.Shutdown.StepTimeout != "" {
			config.ShutdownConfig.StepTimeout = policy.Spec.Shutdown.StepTimeout
		}
		if policy.Spec.Shutdown.RejectPolicy != "" {
			config.ShutdownConfig.RejectPolicy = policy.Spec.Shutdown.RejectPolicy
		}
	}

	// Apply log config
	if policy.Spec.Log != nil {
		if policy.Spec.Log.Level != "" {
			config.Log.Level = policy.Spec.Log.Level
		}
		if policy.Spec.Log.Development != nil {
			// LogConfig needs to be extended to support more fields
		}
		if len(policy.Spec.Log.OutputPaths) > 0 {
			// LogConfig needs to be extended
		}
	}

	// Apply tracing config
	if policy.Spec.Tracing != nil {
		// Tracing config needs to be added to PixiuConfig
		// For now, we skip it as it's not in the current PixiuConfig structure
	}

	// Apply timeout config
	if policy.Spec.Timeout != nil {
		// Global timeout config needs to be added to PixiuConfig
	}
}

// ApplyClusterPolicy applies PixiuClusterPolicy service config to a Cluster
func ApplyClusterPolicy(cluster *Cluster, serviceConfig *v1alpha1.ServiceClusterConfig) {
	if serviceConfig == nil {
		return
	}

	// Apply load balancer policy
	if serviceConfig.LoadBalancer != nil && serviceConfig.LoadBalancer.Policy != "" {
		cluster.LbPolicy = serviceConfig.LoadBalancer.Policy
	}

	// Apply cluster type
	if serviceConfig.Type != "" {
		cluster.Type = serviceConfig.Type
	}

	// Apply static endpoints if specified
	if len(serviceConfig.Endpoints) > 0 {
		cluster.Endpoints = []*Endpoint{}
		for i, ep := range serviceConfig.Endpoints {
			endpoint := &Endpoint{
				SocketAddress: SocketAddress{
					Address: ep.Address,
					Port:    int(ep.Port),
				},
			}
			if ep.ID != nil {
				endpoint.ID = int(*ep.ID)
			} else {
				endpoint.ID = i + 1
			}
			cluster.Endpoints = append(cluster.Endpoints, endpoint)
		}
	}

	// Health check and registry config would need Cluster structure extension
}

// ApplyClusterConfig applies PixiuClusterPolicy cluster config to a Cluster
func ApplyClusterConfig(cluster *Cluster, clusterConfig *v1alpha1.ClusterConfig) {
	if clusterConfig == nil {
		return
	}

	// Apply load balancer policy
	if clusterConfig.LoadBalancerPolicy != "" {
		cluster.LbPolicy = clusterConfig.LoadBalancerPolicy
	}

	// Apply cluster type
	if clusterConfig.Type != "" {
		cluster.Type = clusterConfig.Type
	}

	if len(clusterConfig.Endpoints) > 0 {
		cluster.Endpoints = []*Endpoint{}
		for i, ep := range clusterConfig.Endpoints {
			endpoint := &Endpoint{
				SocketAddress: SocketAddress{
					Address: ep.Address,
					Port:    int(ep.Port),
				},
			}
			if ep.ID != nil {
				endpoint.ID = int(*ep.ID)
			} else {
				endpoint.ID = i + 1
			}
			cluster.Endpoints = append(cluster.Endpoints, endpoint)
		}
	}
}

// ApplyFilterPolicy applies PixiuFilterPolicy config to ir.HTTPFilter
func ApplyFilterPolicy(filter *ir.HTTPFilter, policy *v1alpha1.PixiuFilterPolicy) error {
	if policy == nil {
		return nil
	}

	// Set filter name
	if policy.Spec.FilterType != "" {
		filter.Name = policy.Spec.FilterType
	}

	// Parse the raw JSON config if available
	// runtime.RawExtension.Raw contains JSON bytes when loaded from Kubernetes API
	if len(policy.Spec.Config.Raw) > 0 {
		var configMap map[string]interface{}
		if err := json.Unmarshal(policy.Spec.Config.Raw, &configMap); err != nil {
			return err
		}

		// Merge with existing config
		if filter.Config == nil {
			filter.Config = make(map[string]interface{})
		}
		for k, v := range configMap {
			filter.Config[k] = v
		}
	} else if policy.Spec.Config.Object != nil {
		// If Raw is empty but Object is set, marshal Object to JSON first
		// This can happen when the policy is loaded directly from the API
		configBytes, err := json.Marshal(policy.Spec.Config.Object)
		if err != nil {
			return err
		}
		var configMap map[string]interface{}
		if err := json.Unmarshal(configBytes, &configMap); err != nil {
			return err
		}
		if filter.Config == nil {
			filter.Config = make(map[string]interface{})
		}
		for k, v := range configMap {
			filter.Config[k] = v
		}
	}

	return nil
}

// ApplyFilterPolicyToListener applies PixiuFilterPolicy config to a Listener's network filter
// This is used for TCP listeners with dubboconnectionmanager
func ApplyFilterPolicyToListener(listener *Listener, policy *v1alpha1.PixiuFilterPolicy) error {
	if policy == nil {
		return nil
	}

	// Only apply if this is a dubboconnectionmanager filter
	filterType := policy.Spec.FilterType
	if filterType == "dgp.filter.dubboconnectionmanager" {
		filterType = "dgp.filter.network.dubboconnectionmanager"
	}
	if filterType != "dgp.filter.network.dubboconnectionmanager" {
		return nil
	}

	// Find the dubboconnectionmanager filter in the listener's filter chain
	for i, filter := range listener.FilterChain.Filters {
		if filter.Name == "dgp.filter.network.dubboconnectionmanager" {
			// Parse the config from policy
			var configMap map[string]interface{}
			if len(policy.Spec.Config.Raw) > 0 {
				if err := json.Unmarshal(policy.Spec.Config.Raw, &configMap); err != nil {
					return err
				}
			} else if policy.Spec.Config.Object != nil {
				configBytes, err := json.Marshal(policy.Spec.Config.Object)
				if err != nil {
					return err
				}
				if err := json.Unmarshal(configBytes, &configMap); err != nil {
					return err
				}
			}

			// Merge the config into the existing filter config
			// Normalize snake_case to camelCase for consistent processing
			normalizedConfig := normalizeConfigKeys(configMap)
			if err := mergeFilterConfig(&listener.FilterChain.Filters[i], normalizedConfig, "dgp.filter.network.dubboconnectionmanager"); err != nil {
				// If merge fails, replace config
				listener.FilterChain.Filters[i].Config = configMap
			}
			break
		}
	}

	return nil
}

// ApplyListenersRefToConfig applies PixiuFilterPolicy listenersRef to PixiuConfig
func ApplyListenersRefToConfig(ctx context.Context, k8sClient client.Client, namespace string, config *PixiuConfig, policy *v1alpha1.PixiuFilterPolicy) error {
	if policy == nil || len(policy.Spec.ListenersRef) == 0 {
		return nil
	}

	for _, listenerRef := range policy.Spec.ListenersRef {
		// Find the listener by name
		var targetListener *Listener
		for _, listener := range config.StaticResources.Listeners {
			if listener.Name == listenerRef.Name {
				targetListener = listener
				break
			}
		}

		if targetListener == nil {
			// Listener not found, might need to create it from Gateway spec
			// For now, skip
			continue
		}

		// Parse config from listenerRef
		// In Kubernetes, runtime.RawExtension can be serialized in different ways:
		// 1. If Raw contains JSON bytes, use Raw directly
		// 2. If Raw is empty but Object is set, marshal Object to JSON then unmarshal
		// 3. If both are empty, the config might be embedded directly in the JSON
		// 4. When Kubernetes client deserializes, RawExtension might be directly serialized as the object itself
		var configMap map[string]interface{}

		// First, try to unmarshal Raw if it exists
		if len(listenerRef.Config.Raw) > 0 {
			if err := json.Unmarshal(listenerRef.Config.Raw, &configMap); err != nil {
				return fmt.Errorf("failed to unmarshal listener config from Raw: %w", err)
			}
		} else if listenerRef.Config.Object != nil {
			// Marshal Object to JSON bytes, then unmarshal to map
			configBytes, err := json.Marshal(listenerRef.Config.Object)
			if err != nil {
				return fmt.Errorf("failed to marshal listener config Object: %w", err)
			}
			if err := json.Unmarshal(configBytes, &configMap); err != nil {
				return fmt.Errorf("failed to unmarshal listener config from Object: %w", err)
			}
		} else {
			// Both Raw and Object are empty - try to unmarshal from the listenerRef itself
			// This can happen when Kubernetes client deserializes the resource
			// The config might be embedded directly in the listenerRef struct
			// Marshal the entire listenerRef to JSON, then extract config
			listenerRefBytes, err := json.Marshal(listenerRef)
			if err != nil {
				// If marshal fails, skip this listenerRef
				continue
			}
			var listenerRefMap map[string]interface{}
			if err := json.Unmarshal(listenerRefBytes, &listenerRefMap); err != nil {
				// If unmarshal fails, skip this listenerRef
				continue
			}
			if configVal, ok := listenerRefMap["config"]; ok && configVal != nil {
				if configMapVal, ok := configVal.(map[string]interface{}); ok {
					configMap = configMapVal
				} else {
					// Try to marshal/unmarshal to convert
					configBytes, _ := json.Marshal(configVal)
					if err := json.Unmarshal(configBytes, &configMap); err != nil {
						// If unmarshal fails, configMap will be empty and we'll skip
						configMap = nil
					}
				}
			}
			if len(configMap) == 0 {
				// Still empty, skip this listenerRef
				continue
			}
		}

		// Check if configMap is empty
		if len(configMap) == 0 {
			// Config map is empty, skip
			continue
		}

		// Normalize config keys (route_config -> routeConfig, etc.)
		configMap = normalizeConfigKeys(configMap)

		// Apply filter chain type
		filterType := listenerRef.FilterChains.Type
		if filterType != "" {
			// Normalize filter type: dgp.filter.dubboconnectionmanager -> dgp.filter.network.dubboconnectionmanager
			if filterType == "dgp.filter.dubboconnectionmanager" {
				filterType = "dgp.filter.network.dubboconnectionmanager"
			}

			// If filter type is httpconnectionmanager, update ProtocolType to HTTP
			if filterType == "dgp.filter.httpconnectionmanager" {
				targetListener.ProtocolType = "HTTP"
			} else if filterType == "dgp.filter.network.grpcconnectionmanager" {
				targetListener.ProtocolType = "GRPC"
			} else if filterType == "dgp.filter.network.dubboconnectionmanager" {
				targetListener.ProtocolType = "TCP"
			}

			// Update or create filter chain - always merge instead of replacing
			// Note: A listener should only have one network filter, so we replace if type doesn't match
			if len(targetListener.FilterChain.Filters) == 0 {
				// Create new filter - convert configMap to proper type
				var filterConfig interface{}
				if filterType == "dgp.filter.httpconnectionmanager" {
					// Convert configMap to HTTPConnectionManagerConfig
					configBytes, _ := json.Marshal(configMap)
					var hcmConfig HTTPConnectionManagerConfig
					if err := json.Unmarshal(configBytes, &hcmConfig); err == nil {
						// Ensure all routes have methods field
						for _, route := range hcmConfig.RouteConfig.Routes {
							if len(route.Match.Methods) == 0 {
								route.Match.Methods = []string{}
							}
						}
						filterConfig = hcmConfig
					} else {
						filterConfig = configMap
					}
				} else if filterType == "dgp.filter.network.dubboconnectionmanager" {
					// Convert configMap to DubboConnectionManagerConfig
					configBytes, _ := json.Marshal(configMap)
					var dcmConfig DubboConnectionManagerConfig
					if err := json.Unmarshal(configBytes, &dcmConfig); err == nil {
						filterConfig = dcmConfig
					} else {
						filterConfig = configMap
					}
				} else {
					filterConfig = configMap
				}
				targetListener.FilterChain.Filters = []NetworkFilter{
					{
						Name:   filterType,
						Config: filterConfig,
					},
				}
			} else {
				// Update existing filter or replace if type doesn't match
				found := false
				for i := range targetListener.FilterChain.Filters {
					if targetListener.FilterChain.Filters[i].Name == filterType {
						// Merge config - this will merge routes and filters instead of replacing
						if err := mergeFilterConfig(&targetListener.FilterChain.Filters[i], configMap, filterType); err != nil {
							return fmt.Errorf("failed to merge filter config: %w", err)
						}
						found = true
						break
					}
				}
				if !found {
					// Replace existing filter if type doesn't match (a listener should only have one network filter)
					// Convert configMap to proper type
					normalizeHTTPFilterNames(configMap)
					var filterConfig interface{}
					if filterType == "dgp.filter.httpconnectionmanager" {
						// Convert configMap to HTTPConnectionManagerConfig
						configBytes, _ := json.Marshal(configMap)
						var hcmConfig HTTPConnectionManagerConfig
						if err := json.Unmarshal(configBytes, &hcmConfig); err == nil {
							// Ensure all routes have methods field
							for _, route := range hcmConfig.RouteConfig.Routes {
								if len(route.Match.Methods) == 0 {
									route.Match.Methods = []string{}
								}
							}
							filterConfig = hcmConfig
						} else {
							filterConfig = configMap
						}
					} else if filterType == "dgp.filter.network.dubboconnectionmanager" {
						// Convert configMap to DubboConnectionManagerConfig
						configBytes, _ := json.Marshal(configMap)
						var dcmConfig DubboConnectionManagerConfig
						if err := json.Unmarshal(configBytes, &dcmConfig); err == nil {
							filterConfig = dcmConfig
						} else {
							filterConfig = configMap
						}
					} else {
						filterConfig = configMap
					}
					targetListener.FilterChain.Filters = []NetworkFilter{
						{
							Name:   filterType,
							Config: filterConfig,
						},
					}
				}
			}
		}

		// Resolve registry addresses from k8s endpoints if needed
		if err := resolveRegistryAddresses(ctx, k8sClient, namespace, configMap); err != nil {
			return fmt.Errorf("failed to resolve registry addresses: %w", err)
		}
	}

	return nil
}

// normalizeConfigKeys normalizes snake_case keys to camelCase for consistent processing
func normalizeConfigKeys(configMap map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{})
	for k, v := range configMap {
		switch k {
		case "route_config":
			normalized["routeConfig"] = v
		case "dubbo_filters":
			normalized["dubboFilters"] = v
		case "http_filters":
			normalized["httpFilters"] = v
		default:
			normalized[k] = v
		}
	}
	return normalized
}

// mergeFilterConfig merges config into an existing filter
func mergeFilterConfig(filter *NetworkFilter, configMap map[string]interface{}, filterType string) error {
	switch filterType {
	case "dgp.filter.httpconnectionmanager":
		// normalize legacy alias before merge
		normalizeHTTPFilterNames(configMap)

		// Try to convert existing config to HTTPConnectionManagerConfig
		var existingConfig HTTPConnectionManagerConfig
		var existingConfigOK bool

		if typedConfig, ok := filter.Config.(HTTPConnectionManagerConfig); ok {
			existingConfig = typedConfig
			existingConfigOK = true
		} else if mapConfig, ok := filter.Config.(map[string]interface{}); ok {
			// Try to convert map[string]interface{} to HTTPConnectionManagerConfig
			configBytes, err := json.Marshal(mapConfig)
			if err == nil {
				if err := json.Unmarshal(configBytes, &existingConfig); err == nil {
					existingConfigOK = true
				}
			}
		}

		if existingConfigOK {
			// Merge route_config - append new routes to existing ones instead of replacing
			if routeConfig, ok := configMap["routeConfig"].(map[string]interface{}); ok {
				if routes, ok := routeConfig["routes"].([]interface{}); ok {
					for _, routeInterface := range routes {
						routeBytes, _ := json.Marshal(routeInterface)
						var route Route
						if err := json.Unmarshal(routeBytes, &route); err == nil {
							// Check if route already exists (by prefix/path and cluster)
							exists := false
							for _, existingRoute := range existingConfig.RouteConfig.Routes {
								if (existingRoute.Match.Prefix != "" && existingRoute.Match.Prefix == route.Match.Prefix) ||
									(existingRoute.Match.Path != "" && existingRoute.Match.Path == route.Match.Path) {
									if existingRoute.Route.Cluster == route.Route.Cluster {
										exists = true
										break
									}
								}
							}
							if !exists {
								// Ensure route has methods field (empty methods will default to all HTTP methods in Pixiu)
								if len(route.Match.Methods) == 0 {
									route.Match.Methods = []string{}
								}
								existingConfig.RouteConfig.Routes = append(existingConfig.RouteConfig.Routes, &route)
							}
						}
					}
				}
			}
			// Merge httpFilters - append new filters to existing ones instead of replacing.
			// If any custom filters are provided and httpproxy is not explicitly requested,
			// drop the default httpproxy to avoid conflicts with protocol-specific filters.
			if httpFilters, ok := configMap["httpFilters"].([]interface{}); ok {
				hasDubboProxy := false
				requestsHttpProxy := false
				for _, filterInterface := range httpFilters {
					filterBytes, _ := json.Marshal(filterInterface)
					var httpFilter HTTPFilter
					if err := json.Unmarshal(filterBytes, &httpFilter); err == nil {
						// Normalize alias: policy may use dgp.filter.http.dubboproxy (legacy) which should map to directdubboproxy
						if httpFilter.Name == "dgp.filter.http.dubboproxy" {
							httpFilter.Name = "dgp.filter.http.directdubboproxy"
						}
						if httpFilter.Name == "dgp.filter.http.httpproxy" {
							requestsHttpProxy = true
						}
						if httpFilter.Name == "dgp.filter.http.dubboproxy" || httpFilter.Name == "dgp.filter.http.directdubboproxy" {
							hasDubboProxy = true
						}
						// Check if filter already exists (by name)
						exists := false
						for _, existingFilter := range existingConfig.HTTPFilters {
							if existingFilter.Name == httpFilter.Name {
								exists = true
								break
							}
						}
						if !exists {
							existingConfig.HTTPFilters = append(existingConfig.HTTPFilters, httpFilter)
						}
					}
				}
				// Remove default httpproxy filter if a protocol-specific filter is present
				if hasDubboProxy || (!requestsHttpProxy && len(httpFilters) > 0) {
					filtered := []HTTPFilter{}
					for _, f := range existingConfig.HTTPFilters {
						if f.Name != "dgp.filter.http.httpproxy" {
							filtered = append(filtered, f)
						}
					}
					existingConfig.HTTPFilters = filtered
				}
			}
			filter.Config = existingConfig
		} else {
			// If we can't convert existing config, replace it with new config
			// Convert configMap to HTTPConnectionManagerConfig
			configBytes, _ := json.Marshal(configMap)
			var hcmConfig HTTPConnectionManagerConfig
			if err := json.Unmarshal(configBytes, &hcmConfig); err == nil {
				// Ensure all routes have methods field
				for _, route := range hcmConfig.RouteConfig.Routes {
					if len(route.Match.Methods) == 0 {
						route.Match.Methods = []string{}
					}
				}
				filter.Config = hcmConfig
			} else {
				filter.Config = configMap
			}
		}
	case "dgp.filter.network.dubboconnectionmanager":
		if existingConfig, ok := filter.Config.(DubboConnectionManagerConfig); ok {
			// Merge route_config - append new routes to existing ones instead of replacing
			if routeConfig, ok := configMap["routeConfig"].(map[string]interface{}); ok {
				if routes, ok := routeConfig["routes"].([]interface{}); ok {
					for _, routeInterface := range routes {
						routeBytes, _ := json.Marshal(routeInterface)
						var route Route
						if err := json.Unmarshal(routeBytes, &route); err == nil {
							// Check if route already exists (by prefix/path and cluster)
							exists := false
							for _, existingRoute := range existingConfig.RouteConfig.Routes {
								if (existingRoute.Match.Prefix != "" && existingRoute.Match.Prefix == route.Match.Prefix) ||
									(existingRoute.Match.Path != "" && existingRoute.Match.Path == route.Match.Path) {
									if existingRoute.Route.Cluster == route.Route.Cluster {
										exists = true
										break
									}
								}
							}
							if !exists {
								existingConfig.RouteConfig.Routes = append(existingConfig.RouteConfig.Routes, &route)
							}
						}
					}
				}
			}
			// Merge dubboFilters - update existing filters or append new ones
			if dubboFilters, ok := configMap["dubboFilters"].([]interface{}); ok {
				for _, filterInterface := range dubboFilters {
					filterBytes, _ := json.Marshal(filterInterface)
					var dubboFilter DubboFilter
					if err := json.Unmarshal(filterBytes, &dubboFilter); err == nil {
						// Check if filter already exists (by name)
						found := false
						for i := range existingConfig.DubboFilters {
							if existingConfig.DubboFilters[i].Name == dubboFilter.Name {
								// Update existing filter config instead of skipping
								existingConfig.DubboFilters[i] = dubboFilter
								found = true
								break
							}
						}
						if !found {
							existingConfig.DubboFilters = append(existingConfig.DubboFilters, dubboFilter)
						}
					}
				}
			}
			filter.Config = existingConfig
		} else {
			filter.Config = configMap
		}
	default:
		// For other filter types, log a warning and replace config
		log.Printf("[WARN] mergeFilterConfig: unsupported filter type %q, replacing config instead of merging", filterType)
		filter.Config = configMap
	}
	return nil
}

// normalizeHTTPFilterNames rewrites legacy filter names to the actual plugin kind names
func normalizeHTTPFilterNames(configMap map[string]interface{}) {
	if httpFilters, ok := configMap["httpFilters"].([]interface{}); ok {
		for i, filterInterface := range httpFilters {
			if fMap, ok := filterInterface.(map[string]interface{}); ok {
				if name, ok := fMap["name"].(string); ok && name == "dgp.filter.http.dubboproxy" {
					fMap["name"] = "dgp.filter.http.directdubboproxy"
					httpFilters[i] = fMap
				}
			}
		}
		configMap["httpFilters"] = httpFilters
	}
}

// resolveRegistryAddresses resolves registry service names to IP addresses from k8s endpoints
func resolveRegistryAddresses(ctx context.Context, k8sClient client.Client, namespace string, configMap map[string]interface{}) error {
	// Check for dubboProxyConfig.registries
	if dubboProxyConfig, ok := configMap["dubboProxyConfig"].(map[string]interface{}); ok {
		if registries, ok := dubboProxyConfig["registries"].(map[string]interface{}); ok {
			for _, registryConfig := range registries {
				if registryMap, ok := registryConfig.(map[string]interface{}); ok {
					if address, ok := registryMap["address"].(string); ok {
						// Check if address is a service name (contains : but not IP format)
						if !isIPAddress(address) {
							// Parse service name and port
							parts := strings.Split(address, ":")
							if len(parts) == 2 {
								serviceName := parts[0]
								// Try to resolve service to endpoints
								var service corev1.Service
								if err := k8sClient.Get(ctx, client.ObjectKey{
									Namespace: namespace,
									Name:      serviceName,
								}, &service); err == nil {
									// Get endpoints
									var endpointSliceList discoveryv1.EndpointSliceList
									if err := k8sClient.List(ctx, &endpointSliceList, client.MatchingLabels{
										discoveryv1.LabelServiceName: serviceName,
									}, client.InNamespace(namespace)); err == nil {
										if len(endpointSliceList.Items) > 0 {
											endpointSlice := endpointSliceList.Items[0]
											if len(endpointSlice.Endpoints) > 0 {
												endpoint := endpointSlice.Endpoints[0]
												if len(endpoint.Addresses) > 0 {
													// Use first endpoint address
													resolvedAddress := fmt.Sprintf("%s:%s", endpoint.Addresses[0], parts[1])
													registryMap["address"] = resolvedAddress
												}
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

// isIPAddress checks if a string is an IP address
func isIPAddress(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		// Check for IPv6
		return strings.Contains(s, ":")
	}
	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return false
			}
		}
	}
	return true
}
