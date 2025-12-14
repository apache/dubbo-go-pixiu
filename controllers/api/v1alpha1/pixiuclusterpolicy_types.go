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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PixiuClusterPolicySpec defines the desired state of PixiuClusterPolicy
type PixiuClusterPolicySpec struct {
	// TargetRef identifies an API object to apply policy to.
	TargetRef PolicyTargetReference `json:"targetRef,omitempty"`

	// ServiceRef defines a list of services and their cluster configurations
	ServiceRef []ServiceClusterConfig `json:"serviceRef,omitempty"`

	// ClusterRef defines a list of clusters and their configurations
	ClusterRef []ClusterConfig `json:"clusterRef,omitempty"`
}

// ClusterConfig defines cluster configuration
type ClusterConfig struct {
	// Name specifies the cluster name
	Name string `json:"name"`

	// Type specifies the cluster type (static, dynamic, etc.)
	Type string `json:"type,omitempty"`

	// LoadBalancerPolicy specifies the load balancer policy
	LoadBalancerPolicy string `json:"loadBalancerPolicy,omitempty"`

	// Endpoints defines static endpoints (service names will be resolved from k8s endpoints)
	Endpoints []EndpointConfig `json:"endpoints,omitempty"`
}

// ServiceClusterConfig defines cluster configuration for a service
type ServiceClusterConfig struct {
	// Name specifies the service name (also used as cluster name in conf.yaml)
	Name string `json:"name"`

	// Type specifies the cluster type
	Type string `json:"type,omitempty"`

	// LoadBalancer defines load balancing settings
	LoadBalancer *LoadBalancerConfig `json:"loadBalancer,omitempty"`

	// Registries defines registry configurations
	Registries map[string]RegistryConfig `json:"registries,omitempty"`

	// HealthCheck defines health check settings
	HealthCheck *HealthCheckConfig `json:"healthCheck,omitempty"`

	// Endpoints defines static endpoints (overrides Service discovery if specified)
	Endpoints []EndpointConfig `json:"endpoints,omitempty"`

	// Timeout defines timeout settings for the cluster
	Timeout *ClusterTimeoutConfig `json:"timeout,omitempty"`
}

// LoadBalancerConfig defines load balancing settings
type LoadBalancerConfig struct {
	Policy string `json:"policy,omitempty"`
}

// RegistryConfig defines registry configuration
type RegistryConfig struct {
	Protocol  string `json:"protocol,omitempty"`
	Timeout   string `json:"timeout,omitempty"`
	Address   string `json:"address,omitempty"`
	Username  string `json:"username,omitempty"`
	Password  string `json:"password,omitempty"`
	Group     string `json:"group,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}

// HealthCheckConfig defines health check settings
type HealthCheckConfig struct {
	Protocol           string                 `json:"protocol,omitempty"`
	Timeout            string                 `json:"timeout,omitempty"`
	Interval           string                 `json:"interval,omitempty"`
	HealthyThreshold   *int32                 `json:"healthyThreshold,omitempty"`
	UnhealthyThreshold *int32                 `json:"unhealthyThreshold,omitempty"`
	HTTPHealthCheck    *HTTPHealthCheckConfig `json:"httpHealthCheck,omitempty"`
}

// HTTPHealthCheckConfig defines HTTP health check settings
type HTTPHealthCheckConfig struct {
	Path             string  `json:"path,omitempty"`
	ExpectedStatuses []int32 `json:"expectedStatuses,omitempty"`
}

// EndpointConfig defines a static endpoint
type EndpointConfig struct {
	ID           *int32 `json:"id,omitempty"`
	Address      string `json:"address"`
	Port         int32  `json:"port"`
	ProtocolType string `json:"protocolType,omitempty"`
}

// ClusterTimeoutConfig defines timeout settings for cluster
type ClusterTimeoutConfig struct {
	Connect string `json:"connect,omitempty"`
	Request string `json:"request,omitempty"`
}

// PixiuClusterPolicyStatus defines the observed state of PixiuClusterPolicy
type PixiuClusterPolicyStatus struct {
	// Conditions describe the current conditions of the PixiuClusterPolicy.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=pcp

// PixiuClusterPolicy is the Schema for the pixiuclusterpolicies API
type PixiuClusterPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PixiuClusterPolicySpec   `json:"spec,omitempty"`
	Status PixiuClusterPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PixiuClusterPolicyList contains a list of PixiuClusterPolicy
type PixiuClusterPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PixiuClusterPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PixiuClusterPolicy{}, &PixiuClusterPolicyList{})
}
