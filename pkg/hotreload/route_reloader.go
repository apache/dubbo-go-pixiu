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

package hotreload

import (
	"encoding/json"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

// RouteReloader implements the HotReloader interface for reloading route configurations.
type RouteReloader struct{}

// CheckUpdate compares the old and new route configurations to determine if a reload is needed.
func (r *RouteReloader) CheckUpdate(oldConfig, newConfig *model.Bootstrap) bool {
	oldRoutes := extractRoutes(oldConfig)
	newRoutes := extractRoutes(newConfig)
	// Compare the number of routes
	if len(oldRoutes.Routes) != len(newRoutes.Routes) || oldRoutes.Dynamic != newRoutes.Dynamic {
		return true
	}

	// Compare each route
	for i := range newRoutes.Routes {
		if oldRoutes.Routes[i].Match.Prefix != newRoutes.Routes[i].Match.Prefix ||
			oldRoutes.Routes[i].Route.Cluster != newRoutes.Routes[i].Route.Cluster {
			return true
		}
	}
	return false
}

// HotReload applies the new route configuration.
func (r *RouteReloader) HotReload(oldConfig, newConfig *model.Bootstrap) error {
	oldRoutes := extractRoutes(oldConfig)
	newRoutes := extractRoutes(newConfig)

	// Update routes in the RouterManager
	err := server.GetRouterManager().UpdateRoutes(oldRoutes.Routes, newRoutes.Routes)
	if err != nil {
		logger.Infof("Failed to Routes reloaded.")
		return err
	}

	return nil
}

// extractRoutes extracts routes from the configuration by parsing the filters.
func extractRoutes(config *model.Bootstrap) model.RouteConfiguration {
	var (
		routeConfig     model.RouteConfiguration
		invalidRouteIDs []string
	)
	for _, listener := range config.StaticResources.Listeners {
		for _, filterChain := range listener.FilterChain.Filters {
			if filterChain.Name == constant.HTTPConnectManagerFilter {
				// Extract route_config
				rawRouteConfig, ok := filterChain.Config["route_config"]
				if !ok {
					logger.Debugf("No route_config found in filter chain: %+v", filterChain)
					continue
				}
				logger.Debugf("Raw route_config: %+v", rawRouteConfig)

				// Convert route_config to JSON bytes
				routeConfigBytes, err := json.Marshal(rawRouteConfig)
				if err != nil {
					logger.Errorf("Failed to marshal route_config: %v", err)
					continue
				}

				// Parse JSON bytes into model.RouteConfiguration
				if err := json.Unmarshal(routeConfigBytes, &routeConfig); err != nil {
					logger.Errorf("Failed to unmarshal route_config: %v", err)
					continue
				}

				logger.Debugf("Parsed route_config: %+v", routeConfig)

				// Validate and filter routes
				validRoutes := make([]*model.Router, 0, len(routeConfig.Routes))
				for _, route := range routeConfig.Routes {
					if err := validateRoute(route); err != nil {
						invalidRouteIDs = append(invalidRouteIDs, route.ID)
						logger.Warnf("Skipping invalid route %s: %v", route.ID, err)
						continue
					}
					validRoutes = append(validRoutes, route)
				}

				routeConfig.Routes = validRoutes
				logger.Debugf("Valid routes after filtering: %+v", validRoutes)

				// Return if we have valid routes
				if len(validRoutes) > 0 {
					return routeConfig
				}
			}
		}
	}

	if len(invalidRouteIDs) > 0 {
		logger.Warnf("No valid routes found in configuration: %v", invalidRouteIDs)
	}
	return routeConfig
}

// validateRoute validates a single route, returning an error if invalid.
func validateRoute(route *model.Router) error {
	// Ensure route has a valid match condition
	if route.Match.Prefix == "" && route.Match.Path == "" {
		return errors.Errorf("route %s has no prefix or path defined", route.ID)
	}

	// Ensure cluster is specified
	if route.Route.Cluster == "" {
		return errors.Errorf("route %s has no cluster defined", route.ID)
	}

	return nil
}
