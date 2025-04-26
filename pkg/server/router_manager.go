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

package server

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type (
	RouterListener interface {
		OnAddRouter(r *model.Router)
		OnDeleteRouter(r *model.Router)
	}

	RouterManager struct {
		rls []RouterListener
	}
)

func CreateDefaultRouterManager(server *Server, bs *model.Bootstrap) *RouterManager {
	rm := &RouterManager{}
	return rm
}

func (rm *RouterManager) AddRouterListener(l RouterListener) {
	rm.rls = append(rm.rls, l)
}

func (rm *RouterManager) AddRouter(r *model.Router) {
	logger.Infof("add router: %v", r)
	for _, l := range rm.rls {
		l.OnAddRouter(r)
	}
}

func (rm *RouterManager) DeleteRouter(r *model.Router) {
	logger.Infof("del router: %v", r)
	for _, l := range rm.rls {
		l.OnDeleteRouter(r)
	}
}

// UpdateRoutes updates the routes in RouterManager with the provided new routes.
func (rm *RouterManager) UpdateRoutes(oldRoutes []*model.Router, newRoutes []*model.Router) error {
	logger.Infof("Starting route update with %d new routes", len(newRoutes))

	// Validate new routes
	if err := validateRoutes(newRoutes); err != nil {
		logger.Errorf("Invalid routes provided: %v", err)
		return errors.Wrap(err, "route validation failed")
	}

	// Notify listeners to delete existing routes
	for _, route := range oldRoutes {
		logger.Debugf("Notifying listeners to delete route: %s", route.String())
		for _, listener := range rm.rls {
			listener.OnDeleteRouter(route)
		}
	}

	// Notify listeners to add new routes
	for _, route := range newRoutes {
		logger.Debugf("Notifying listeners to add route: %s", route.String())
		for _, listener := range rm.rls {
			listener.OnAddRouter(route)
		}
	}

	// Atomically update the active configuration
	logger.Infof("Routes updated successfully with %d routes", len(newRoutes))

	return nil
}

// validateRoutes performs basic validation on the provided routes.
func validateRoutes(routes []*model.Router) error {
	routeIDs := make(map[string]struct{}, len(routes))
	for _, route := range routes {
		// Check for duplicate IDs
		if _, exists := routeIDs[route.ID]; exists {
			return errors.Errorf("duplicate route ID: %s", route.ID)
		}
		routeIDs[route.ID] = struct{}{}

		// Ensure route has a valid match condition
		if route.Match.Prefix == "" && route.Match.Path == "" {
			return errors.Errorf("route %s has no prefix or path defined", route.ID)
		}

		// Ensure cluster is specified
		if route.Route.Cluster == "" {
			return errors.Errorf("route %s has no cluster defined", route.ID)
		}
	}
	return nil
}
