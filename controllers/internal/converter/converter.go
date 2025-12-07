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
	"fmt"
)

import (
	"controllers/internal/ir"

	"gopkg.in/yaml.v3"
)

type Converter struct{}

func NewConverter() *Converter {
	return &Converter{}
}

func (c *Converter) ConvertIRToPixiuConfig(xds *ir.Xds, clusters []*ir.Cluster) (*PixiuConfig, error) {
	config := &PixiuConfig{
		StaticResources: StaticResources{
			Listeners: []*Listener{},
			Clusters:  []*Cluster{},
		},
		ShutdownConfig: ShutdownConfig{
			Timeout:      "60s",
			StepTimeout:  "10s",
			RejectPolicy: "immediacy",
		},
		Log: LogConfig{
			Level: "info",
		},
	}

	for _, httpListener := range xds.HTTP {
		listener, err := c.convertHTTPListener(httpListener, clusters)
		if err != nil {
			return nil, fmt.Errorf("failed to convert HTTP listener %s: %w", httpListener.Name, err)
		}
		if listener != nil {
			config.StaticResources.Listeners = append(config.StaticResources.Listeners, listener)
		}
	}

	for _, tcpListener := range xds.TCP {
		listener, err := c.convertTCPListener(tcpListener, clusters)
		if err != nil {
			return nil, fmt.Errorf("failed to convert TCP listener %s: %w", tcpListener.Name, err)
		}
		if listener != nil {
			config.StaticResources.Listeners = append(config.StaticResources.Listeners, listener)
		}
	}

	for _, udpListener := range xds.UDP {
		listener, err := c.convertUDPListener(udpListener)
		if err != nil {
			return nil, fmt.Errorf("failed to convert UDP listener %s: %w", udpListener.Name, err)
		}
		if listener != nil {
			config.StaticResources.Listeners = append(config.StaticResources.Listeners, listener)
		}
	}

	for _, irCluster := range clusters {
		cluster := c.convertCluster(irCluster)
		config.StaticResources.Clusters = append(config.StaticResources.Clusters, cluster)
	}

	return config, nil
}

func (c *Converter) convertHTTPListener(httpListener *ir.HTTPListener, clusters []*ir.Cluster) (*Listener, error) {
	routes := []*Route{}
	for _, irRoute := range httpListener.Routes {
		route := c.convertHTTPRoute(irRoute)
		if route != nil {
			routes = append(routes, route)
		}
	}

	httpFilters := []HTTPFilter{
		{
			Name:   "dgp.filter.http.httpproxy",
			Config: map[string]interface{}{},
		},
	}

	for _, irRoute := range httpListener.Routes {
		for _, filter := range irRoute.Filters {
			httpFilters = append(httpFilters, HTTPFilter{
				Name:   filter.Name,
				Config: filter.Config,
			})
		}
	}

	hcmConfig := HTTPConnectionManagerConfig{
		RouteConfig: RouteConfiguration{
			Routes: routes,
		},
		HTTPFilters: httpFilters,
	}

	filterChain := FilterChain{
		Filters: []NetworkFilter{
			{
				Name:   "dgp.filter.httpconnectionmanager",
				Config: hcmConfig,
			},
		},
	}

	listenerConfig := map[string]interface{}{
		"idle_timeout":  "5s",
		"read_timeout":  "5s",
		"write_timeout": "5s",
	}

	return &Listener{
		Name:         httpListener.Name,
		ProtocolType: "HTTP",
		Address: Address{
			SocketAddress: SocketAddress{
				Address: httpListener.Address,
				Port:    int(httpListener.Port),
			},
		},
		FilterChain: filterChain,
		Config:      listenerConfig,
	}, nil
}

func (c *Converter) convertHTTPRoute(irRoute *ir.HTTPRoute) *Route {
	match := RouteMatch{}

	if irRoute.PathMatch != nil {
		if irRoute.PathMatch.Exact != nil {
			match.Path = *irRoute.PathMatch.Exact
		} else if irRoute.PathMatch.Prefix != nil {
			match.Prefix = *irRoute.PathMatch.Prefix
		} else if irRoute.PathMatch.SafeRegex != nil {
			match.Path = *irRoute.PathMatch.SafeRegex
		}
	} else {
		match.Prefix = "/"
	}

	if irRoute.Method != nil {
		match.Methods = []string{*irRoute.Method}
	}

	if len(irRoute.HeaderMatches) > 0 {
		match.Headers = []HeaderMatcher{}
		for _, headerMatch := range irRoute.HeaderMatches {
			header := HeaderMatcher{
				Name: *headerMatch.Name,
			}
			if headerMatch.Exact != nil {
				header.Values = []string{*headerMatch.Exact}
			} else if headerMatch.SafeRegex != nil {
				header.Values = []string{*headerMatch.SafeRegex}
				header.Regex = true
			}
			match.Headers = append(match.Headers, header)
		}
	}

	clusterName := ""
	if len(irRoute.Destinations) > 0 {
		clusterName = irRoute.Destinations[0].Name
	}

	routeAction := RouteAction{
		Cluster:                     clusterName,
		ClusterNotFoundResponseCode: 503,
	}

	return &Route{
		Match: match,
		Route: routeAction,
	}
}

func (c *Converter) convertTCPListener(tcpListener *ir.TCPListener, clusters []*ir.Cluster) (*Listener, error) {
	// TODO: TCP listener conversion
	return &Listener{
		Name:         tcpListener.Name,
		ProtocolType: "TCP",
		Address: Address{
			SocketAddress: SocketAddress{
				Address: tcpListener.Address,
				Port:    int(tcpListener.Port),
			},
		},
		FilterChain: FilterChain{
			Filters: []NetworkFilter{},
		},
	}, nil
}

// convertUDPListener converts an IR UDPListener to a Pixiu Listener.
func (c *Converter) convertUDPListener(udpListener *ir.UDPListener) (*Listener, error) {
	// TODO: UDP listener conversion
	return &Listener{
		Name:         udpListener.Name,
		ProtocolType: "UDP",
		Address: Address{
			SocketAddress: SocketAddress{
				Address: udpListener.Address,
				Port:    int(udpListener.Port),
			},
		},
		FilterChain: FilterChain{
			Filters: []NetworkFilter{},
		},
	}, nil
}

func (c *Converter) convertCluster(irCluster *ir.Cluster) *Cluster {
	endpoints := []*Endpoint{}
	for i, irEndpoint := range irCluster.Endpoints {
		endpoints = append(endpoints, &Endpoint{
			ID: i + 1,
			SocketAddress: SocketAddress{
				Address: irEndpoint.Address,
				Port:    int(irEndpoint.Port),
			},
		})
	}

	lbPolicy := "lb" // Default to round_robin
	if irCluster.LoadBalancerPolicy != nil {
		switch irCluster.LoadBalancerPolicy.Type {
		case "round_robin":
			lbPolicy = "lb"
		case "least_conn":
			lbPolicy = "least_conn"
		default:
			lbPolicy = "lb"
		}
	}

	return &Cluster{
		Name:      irCluster.Name,
		Type:      "Static",
		LbPolicy:  lbPolicy,
		Endpoints: endpoints,
	}
}

func ConvertToYAML(config *PixiuConfig) (string, error) {
	data, err := yaml.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal config to YAML: %w", err)
	}
	return string(data), nil
}
