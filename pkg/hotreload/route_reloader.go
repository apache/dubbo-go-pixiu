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
	if len(oldRoutes.Routes) != len(newRoutes.Routes) || oldRoutes.Dynamic != newRoutes.Dynamic {
		return true
	}

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
	logger.Info("Starting route hot reload")

	srv := server.GetServer()
	if srv == nil {
		logger.Error("Server instance is nil")
		return errors.New("server instance is nil")
	}
	logger.Info("Got server instance")

	logger.Info("Reinitializing server components...")

	listenerManager := srv.GetListenerManager()
	if listenerManager == nil {
		logger.Error("Listener manager is nil")
		return errors.New("listener manager is nil")
	}

	srv.GetRouterManager().ClearRouterListeners()

	refreshed := 0
	for _, listener := range newConfig.StaticResources.Listeners {
		logger.Infof("Refreshing listener: name=%s, protocol=%s", listener.Name, listener.ProtocolStr)

		if err := listenerManager.UpdateListener(listener); err != nil {
			logger.Errorf("Failed to refresh listener %s: %v", listener.Name, err)
			return errors.Wrapf(err, "failed to refresh listener %s", listener.Name)
		}
		logger.Infof("Successfully refreshed listener: %s", listener.Name)
		refreshed++
	}

	if refreshed == 0 {
		logger.Warn("No listeners were refreshed")
	} else {
		logger.Infof("Successfully refreshed %d listener(s)", refreshed)
	}

	logger.Info("Route hot reload completed successfully")
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
				rawRouteConfig, ok := filterChain.Config["route_config"]
				if !ok {
					logger.Debugf("No route_config found in filter chain: %+v", filterChain)
					continue
				}
				logger.Debugf("Raw route_config: %+v", rawRouteConfig)

				routeConfigBytes, err := json.Marshal(rawRouteConfig)
				if err != nil {
					logger.Errorf("Failed to marshal route_config: %v", err)
					continue
				}

				if err := json.Unmarshal(routeConfigBytes, &routeConfig); err != nil {
					logger.Errorf("Failed to unmarshal route_config: %v", err)
					continue
				}

				logger.Debugf("Parsed route_config: %+v", routeConfig)

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
	if route.Match.Prefix == "" && route.Match.Path == "" {
		return errors.Errorf("route %s has no prefix or path defined", route.ID)
	}

	if route.Route.Cluster == "" {
		return errors.Errorf("route %s has no cluster defined", route.ID)
	}

	return nil
}
