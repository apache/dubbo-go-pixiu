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

package translator

import (
	"context"
	"fmt"
)

import (
	"controllers/internal/ir"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"

	discoveryv1 "k8s.io/api/discovery/v1"

	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

type Translator struct {
	client client.Client
	logger logr.Logger
}

func NewTranslator(client client.Client, logger logr.Logger) *Translator {
	return &Translator{
		client: client,
		logger: logger,
	}
}

func (t *Translator) TranslateGateway(ctx context.Context, gateway *gatewayv1.Gateway) (*ir.Xds, error) {
	xds := &ir.Xds{
		HTTP: []*ir.HTTPListener{},
		TCP:  []*ir.TCPListener{},
		UDP:  []*ir.UDPListener{},
	}

	var httpRouteList gatewayv1.HTTPRouteList
	if err := t.client.List(ctx, &httpRouteList); err != nil {
		return nil, fmt.Errorf("failed to list HTTPRoutes: %w", err)
	}

	for _, listener := range gateway.Spec.Listeners {
		if listener.Protocol == gatewayv1.HTTPProtocolType || listener.Protocol == gatewayv1.HTTPSProtocolType {
			httpListener := t.translateHTTPListener(&listener, gateway, &httpRouteList)
			if httpListener != nil {
				xds.HTTP = append(xds.HTTP, httpListener)
			}
		} else if listener.Protocol == gatewayv1.TCPProtocolType {
			tcpListener := t.translateTCPListener(&listener, gateway)
			if tcpListener != nil {
				xds.TCP = append(xds.TCP, tcpListener)
			}
		} else if listener.Protocol == gatewayv1.UDPProtocolType {
			udpListener := t.translateUDPListener(&listener, gateway)
			if udpListener != nil {
				xds.UDP = append(xds.UDP, udpListener)
			}
		}
	}

	return xds, nil
}

func (t *Translator) translateHTTPListener(listener *gatewayv1.Listener, gateway *gatewayv1.Gateway, httpRouteList *gatewayv1.HTTPRouteList) *ir.HTTPListener {
	address := "0.0.0.0"
	if len(gateway.Spec.Addresses) > 0 {
		address = gateway.Spec.Addresses[0].Value
	}

	// Pixiu gateway listens on port 8888 internally
	// Gateway listener port is mapped via Service
	httpListener := &ir.HTTPListener{
		Name:    string(listener.Name),
		Address: address,
		Port:    8888, // Pixiu gateway default port
		Routes:  []*ir.HTTPRoute{},
	}

	matchedRoutes := 0
	for _, httpRoute := range httpRouteList.Items {
		attached := false
		for _, parentRef := range httpRoute.Spec.ParentRefs {
			if parentRef.Group != nil && *parentRef.Group != gatewayv1.GroupName {
				continue
			}
			if parentRef.Kind != nil && *parentRef.Kind != "Gateway" {
				continue
			}
			if string(parentRef.Name) != gateway.Name {
				continue
			}

			namespaceMatch := false
			if parentRef.Namespace != nil {
				namespaceMatch = string(*parentRef.Namespace) == gateway.Namespace
			} else {
				namespaceMatch = gateway.Namespace == httpRoute.Namespace
			}
			if !namespaceMatch {
				continue
			}

			if parentRef.SectionName != nil && *parentRef.SectionName != listener.Name {
				continue
			}

			attached = true
			break
		}
		if !attached {
			continue
		}

		matchedRoutes++
		for _, rule := range httpRoute.Spec.Rules {
			for _, match := range rule.Matches {
				irRoute := t.translateHTTPRouteMatch(&match, &httpRoute, rule.BackendRefs)
				if irRoute != nil {
					httpListener.Routes = append(httpListener.Routes, irRoute)
				}
			}
		}
	}

	return httpListener
}

func (t *Translator) translateHTTPRouteMatch(match *gatewayv1.HTTPRouteMatch, httpRoute *gatewayv1.HTTPRoute, backendRefs []gatewayv1.HTTPBackendRef) *ir.HTTPRoute {
	pathTypeStr := "PathPrefix"
	if match.Path != nil && match.Path.Type != nil {
		pathTypeStr = string(*match.Path.Type)
	}
	irRoute := &ir.HTTPRoute{
		Name: fmt.Sprintf("%s-%s", httpRoute.Name, pathTypeStr),
	}

	if match.Path != nil {
		if match.Path.Type != nil && match.Path.Value != nil {
			pathType := string(*match.Path.Type)
			pathValue := *match.Path.Value
			switch pathType {
			case string(gatewayv1.PathMatchExact):
				irRoute.PathMatch = &ir.StringMatch{
					Exact: &pathValue,
				}
			case string(gatewayv1.PathMatchPathPrefix):
				irRoute.PathMatch = &ir.StringMatch{
					Prefix: &pathValue,
				}
			case string(gatewayv1.PathMatchRegularExpression):
				irRoute.PathMatch = &ir.StringMatch{
					SafeRegex: &pathValue,
				}
			}
		}
	} else {
		// Default to prefix "/"
		prefix := "/"
		irRoute.PathMatch = &ir.StringMatch{
			Prefix: &prefix,
		}
	}

	if match.Method != nil {
		method := string(*match.Method)
		irRoute.Method = &method
	}

	if len(match.Headers) > 0 {
		irRoute.HeaderMatches = []*ir.StringMatch{}
		for _, headerMatch := range match.Headers {
			headerType := string(*headerMatch.Type)
			headerName := string(headerMatch.Name)
			headerValue := headerMatch.Value
			var stringMatch *ir.StringMatch
			switch headerType {
			case string(gatewayv1.HeaderMatchExact):
				stringMatch = &ir.StringMatch{
					Name:  &headerName,
					Exact: &headerValue,
				}
			case string(gatewayv1.HeaderMatchRegularExpression):
				stringMatch = &ir.StringMatch{
					Name:      &headerName,
					SafeRegex: &headerValue,
				}
			}
			if stringMatch != nil {
				irRoute.HeaderMatches = append(irRoute.HeaderMatches, stringMatch)
			}
		}
	}

	irRoute.Destinations = []*ir.RouteDestination{}
	for _, backendRef := range backendRefs {
		if backendRef.Kind != nil && *backendRef.Kind != "Service" {
			continue
		}
		if backendRef.Group != nil && *backendRef.Group != "" {
			continue
		}
		namespace := httpRoute.Namespace
		if backendRef.Namespace != nil {
			namespace = string(*backendRef.Namespace)
		}
		clusterName := fmt.Sprintf("%s-%s", namespace, string(backendRef.Name))
		irRoute.Destinations = append(irRoute.Destinations, &ir.RouteDestination{
			Name: clusterName,
		})
	}

	return irRoute
}

func (t *Translator) translateTCPListener(listener *gatewayv1.Listener, gateway *gatewayv1.Gateway) *ir.TCPListener {
	address := "0.0.0.0"
	if len(gateway.Spec.Addresses) > 0 {
		address = gateway.Spec.Addresses[0].Value
	}

	return &ir.TCPListener{
		Name:    string(listener.Name),
		Address: address,
		Port:    int32(listener.Port),
	}
}

func (t *Translator) translateUDPListener(listener *gatewayv1.Listener, gateway *gatewayv1.Gateway) *ir.UDPListener {
	address := "0.0.0.0"
	if len(gateway.Spec.Addresses) > 0 {
		address = gateway.Spec.Addresses[0].Value
	}

	return &ir.UDPListener{
		Name:    string(listener.Name),
		Address: address,
		Port:    int32(listener.Port),
	}
}

func (t *Translator) TranslateBackendToCluster(ctx context.Context, httpRoute *gatewayv1.HTTPRoute) ([]*ir.Cluster, error) {
	clusters := []*ir.Cluster{}
	clusterMap := make(map[string]*ir.Cluster)

	for _, rule := range httpRoute.Spec.Rules {
		for _, backendRef := range rule.BackendRefs {
			namespace := httpRoute.Namespace
			if backendRef.Namespace != nil {
				namespace = string(*backendRef.Namespace)
			}

			clusterName := fmt.Sprintf("%s-%s", namespace, string(backendRef.Name))
			port := int32(80)
			if backendRef.Port != nil {
				port = int32(*backendRef.Port)
			}

			endpoints, err := t.resolveEndpoints(ctx, namespace, string(backendRef.Name), port)
			if err != nil {
				t.logger.Error(err, "failed to resolve endpoints", "namespace", namespace, "service", backendRef.Name)
				continue
			}

			if existing, exists := clusterMap[clusterName]; exists {
				// Merge endpoints
				existing.Endpoints = append(existing.Endpoints, endpoints...)
			} else {
				cluster := &ir.Cluster{
					Name:      clusterName,
					Endpoints: endpoints,
				}
				clusterMap[clusterName] = cluster
				clusters = append(clusters, cluster)
			}
		}
	}

	return clusters, nil
}

func (t *Translator) resolveEndpoints(ctx context.Context, namespace, serviceName string, port int32) ([]*ir.Endpoint, error) {
	var service corev1.Service
	if err := t.client.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      serviceName,
	}, &service); err != nil {
		return nil, fmt.Errorf("failed to get service: %w", err)
	}

	var endpointSliceList discoveryv1.EndpointSliceList
	if err := t.client.List(ctx, &endpointSliceList, client.MatchingLabels{
		discoveryv1.LabelServiceName: serviceName,
	}, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list endpoint slices: %w", err)
	}

	endpoints := []*ir.Endpoint{}
	for _, endpointSlice := range endpointSliceList.Items {
		for _, endpoint := range endpointSlice.Endpoints {
			if endpoint.Conditions.Ready != nil && !*endpoint.Conditions.Ready {
				continue
			}
			for _, address := range endpoint.Addresses {
				targetPort := port
				if len(endpointSlice.Ports) > 0 {
					for _, endpointPort := range endpointSlice.Ports {
						if endpointPort.Port != nil {
							targetPort = int32(*endpointPort.Port)
							break
						}
					}
				}
				endpoints = append(endpoints, &ir.Endpoint{
					Address: address,
					Port:    targetPort,
				})
			}
		}
	}

	if len(endpoints) == 0 && service.Spec.ClusterIP != "" && service.Spec.ClusterIP != "None" {
		endpoints = append(endpoints, &ir.Endpoint{
			Address: service.Spec.ClusterIP,
			Port:    port,
		})
	}

	return endpoints, nil
}
