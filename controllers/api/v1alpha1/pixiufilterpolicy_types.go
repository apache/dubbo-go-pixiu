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
	"k8s.io/apimachinery/pkg/runtime"
)

// PixiuFilterPolicySpec defines the desired state of PixiuFilterPolicy
type PixiuFilterPolicySpec struct {
	// TargetRef identifies an API object to apply policy to.
	TargetRef PolicyTargetReference `json:"targetRef"`

	// FilterType specifies the type of filter to configure
	// +kubebuilder:validation:Enum=dgp.filter.http.auth.jwt;dgp.filter.http.opa;dgp.filter.http.cors;dgp.filter.http.csrf;dgp.filter.http.traffic;dgp.filter.http.ratelimit;dgp.filter.http.httpproxy;dgp.filter.http.dubboproxy;dgp.filter.http.directdubboproxy;dgp.filter.http.apiconfig;dgp.filter.http.loadbalance;dgp.filter.http.prometheusmetric;dgp.filter.http.seata;dgp.filter.mcp.mcpserver;dgp.filter.network.dubboconnectionmanager;dgp.filter.network.grpcconnectionmanager
	FilterType string `json:"filterType,omitempty"`

	// Config contains the filter-specific configuration
	// Using runtime.RawExtension to support arbitrary JSON configuration
	Config runtime.RawExtension `json:"config,omitempty"`

	// ListenersRef defines listener-specific configurations
	// This is used when targeting a Gateway to configure multiple listeners
	ListenersRef []ListenerRefConfig `json:"listenersRef,omitempty"`
}

// ListenerRefConfig defines configuration for a specific listener
type ListenerRefConfig struct {
	// Name specifies the listener name
	Name string `json:"name"`

	// FilterChains defines the filter chain configuration
	FilterChains FilterChainConfig `json:"filterChains"`

	// Config contains the listener-specific configuration
	Config runtime.RawExtension `json:"config,omitempty"`
}

// FilterChainConfig defines filter chain configuration
type FilterChainConfig struct {
	// Type specifies the filter chain type
	Type string `json:"type"`
}

// PixiuFilterPolicyStatus defines the observed state of PixiuFilterPolicy
type PixiuFilterPolicyStatus struct {
	// Conditions describe the current conditions of the PixiuFilterPolicy.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=pfp

// PixiuFilterPolicy is the Schema for the pixiufilterpolicies API
type PixiuFilterPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PixiuFilterPolicySpec   `json:"spec,omitempty"`
	Status PixiuFilterPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PixiuFilterPolicyList contains a list of PixiuFilterPolicy
type PixiuFilterPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PixiuFilterPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PixiuFilterPolicy{}, &PixiuFilterPolicyList{})
}
