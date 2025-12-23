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

package controller

import (
	"context"
	"fmt"
)

import (
	"controllers/api/v1alpha1"

	"github.com/go-logr/logr"

	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// PolicyLoader loads and applies Pixiu policies
type PolicyLoader struct {
	client client.Client
	logger logr.Logger
}

// NewPolicyLoader creates a new PolicyLoader
func NewPolicyLoader(client client.Client, logger logr.Logger) *PolicyLoader {
	return &PolicyLoader{
		client: client,
		logger: logger,
	}
}

// LoadGatewayPolicy loads PixiuGatewayPolicy for a Gateway
func (pl *PolicyLoader) LoadGatewayPolicy(ctx context.Context, gateway *gatewayv1.Gateway) (*v1alpha1.PixiuGatewayPolicy, error) {
	var policyList v1alpha1.PixiuGatewayPolicyList
	if err := pl.client.List(ctx, &policyList, client.InNamespace(gateway.Namespace)); err != nil {
		return nil, fmt.Errorf("failed to list gateway policies: %w", err)
	}

	for _, policy := range policyList.Items {
		if pl.matchesGateway(&policy, gateway) {
			return &policy, nil
		}
	}

	return nil, nil
}

// LoadFilterPolicy loads a specific PixiuFilterPolicy by name
func (pl *PolicyLoader) LoadFilterPolicy(ctx context.Context, namespace, name string) (*v1alpha1.PixiuFilterPolicy, error) {
	var policy v1alpha1.PixiuFilterPolicy
	if err := pl.client.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}, &policy); err != nil {
		return nil, fmt.Errorf("failed to get filter policy: %w", err)
	}
	return &policy, nil
}

// LoadFilterPolicies loads PixiuFilterPolicy for an HTTPRoute
func (pl *PolicyLoader) LoadFilterPolicies(ctx context.Context, httpRoute *gatewayv1.HTTPRoute) (map[string]*v1alpha1.PixiuFilterPolicy, error) {
	policies := make(map[string]*v1alpha1.PixiuFilterPolicy)

	var policyList v1alpha1.PixiuFilterPolicyList
	if err := pl.client.List(ctx, &policyList, client.InNamespace(httpRoute.Namespace)); err != nil {
		return nil, fmt.Errorf("failed to list filter policies: %w", err)
	}

	for i := range policyList.Items {
		policy := &policyList.Items[i]
		if pl.matchesHTTPRoute(policy, httpRoute) {
			policies[policy.Name] = policy
		}
	}

	return policies, nil
}

// LoadClusterPolicy loads PixiuClusterPolicy for a Service
// This is kept for backward compatibility but may not be used with the new serviceRef format
func (pl *PolicyLoader) LoadClusterPolicy(ctx context.Context, namespace, serviceName string) (*v1alpha1.PixiuClusterPolicy, error) {
	var policyList v1alpha1.PixiuClusterPolicyList
	if err := pl.client.List(ctx, &policyList, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list cluster policies: %w", err)
	}

	// Search through all policies and their serviceRef entries
	for _, policy := range policyList.Items {
		for _, serviceConfig := range policy.Spec.ServiceRef {
			if serviceConfig.Name == serviceName {
				return &policy, nil
			}
		}
	}

	return nil, nil
}

// LoadClusterPolicyByClusterName loads PixiuClusterPolicy by cluster name
// This searches through serviceRef entries to find matching service name
func (pl *PolicyLoader) LoadClusterPolicyByClusterName(ctx context.Context, namespace, clusterName string) (*v1alpha1.PixiuClusterPolicy, error) {
	var policyList v1alpha1.PixiuClusterPolicyList
	if err := pl.client.List(ctx, &policyList, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list cluster policies: %w", err)
	}

	// Search through all policies and their serviceRef entries
	for _, policy := range policyList.Items {
		for _, serviceConfig := range policy.Spec.ServiceRef {
			if serviceConfig.Name == clusterName {
				return &policy, nil
			}
		}
	}

	return nil, nil
}

// LoadAllClusterPolicies loads all PixiuClusterPolicy in a namespace
func (pl *PolicyLoader) LoadAllClusterPolicies(ctx context.Context, namespace string) ([]v1alpha1.PixiuClusterPolicy, error) {
	var policyList v1alpha1.PixiuClusterPolicyList
	if err := pl.client.List(ctx, &policyList, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list cluster policies: %w", err)
	}
	return policyList.Items, nil
}

// matchesGateway checks if a policy matches a Gateway
func (pl *PolicyLoader) matchesGateway(policy *v1alpha1.PixiuGatewayPolicy, gateway *gatewayv1.Gateway) bool {
	ref := policy.Spec.TargetRef
	if ref.Kind != "Gateway" {
		return false
	}
	if ref.Name != gateway.Name {
		return false
	}
	if ref.Namespace != nil && string(*ref.Namespace) != gateway.Namespace {
		return false
	}
	if ref.Group != "" && ref.Group != "gateway.networking.k8s.io" {
		return false
	}
	return true
}

// matchesHTTPRoute checks if a policy matches an HTTPRoute
func (pl *PolicyLoader) matchesHTTPRoute(policy *v1alpha1.PixiuFilterPolicy, httpRoute *gatewayv1.HTTPRoute) bool {
	ref := policy.Spec.TargetRef
	if ref.Kind != "HTTPRoute" {
		return false
	}
	if ref.Name != httpRoute.Name {
		return false
	}
	if ref.Namespace != nil && string(*ref.Namespace) != httpRoute.Namespace {
		return false
	}
	if ref.Group != "" && ref.Group != "gateway.networking.k8s.io" {
		return false
	}
	return true
}

// matchesService checks if a policy matches a Service
// This checks if any serviceRef entry matches the service name
func (pl *PolicyLoader) matchesService(policy *v1alpha1.PixiuClusterPolicy, namespace, serviceName string) bool {
	for _, serviceConfig := range policy.Spec.ServiceRef {
		if serviceConfig.Name == serviceName {
			return true
		}
	}
	return false
}
