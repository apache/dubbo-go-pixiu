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

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

// PolicyTargetReference identifies an API object to apply policy to.
type PolicyTargetReference struct {
	// Group is the group of the target resource.
	Group string `json:"group,omitempty"`

	// Kind is kind of the target resource.
	Kind string `json:"kind"`

	// Name is the name of the target resource.
	Name string `json:"name"`

	// Namespace is the namespace of the target resource.
	Namespace *gatewayv1.Namespace `json:"namespace,omitempty"`
}

// PixiuGatewayPolicySpec defines the desired state of PixiuGatewayPolicy
type PixiuGatewayPolicySpec struct {
	// TargetRef identifies an API object to apply policy to.
	TargetRef PolicyTargetReference `json:"targetRef"`

	// Listener configuration
	Listener *ListenerConfig `json:"listener,omitempty"`

	// Shutdown configuration
	Shutdown *ShutdownConfig `json:"shutdown,omitempty"`

	// Log configuration
	Log *LogConfig `json:"log,omitempty"`

	// Tracing configuration
	Tracing *TracingConfig `json:"tracing,omitempty"`

	// Global timeout configuration
	Timeout *TimeoutConfig `json:"timeout,omitempty"`
}

// ListenerConfig defines listener timeout settings
type ListenerConfig struct {
	Timeout *ListenerTimeout `json:"timeout,omitempty"`
}

// ListenerTimeout defines timeout values for listener
type ListenerTimeout struct {
	Idle  string `json:"idle,omitempty"`
	Read  string `json:"read,omitempty"`
	Write string `json:"write,omitempty"`
}

// ShutdownConfig defines graceful shutdown settings
type ShutdownConfig struct {
	Timeout      string `json:"timeout,omitempty"`
	StepTimeout  string `json:"stepTimeout,omitempty"`
	RejectPolicy string `json:"rejectPolicy,omitempty"`
}

// LogConfig defines logging settings
type LogConfig struct {
	Level             string   `json:"level,omitempty"`
	Development       *bool    `json:"development,omitempty"`
	DisableCaller     *bool    `json:"disableCaller,omitempty"`
	DisableStacktrace *bool    `json:"disableStacktrace,omitempty"`
	Encoding          string   `json:"encoding,omitempty"`
	OutputPaths       []string `json:"outputPaths,omitempty"`
	ErrorOutputPaths  []string `json:"errorOutputPaths,omitempty"`
}

// TracingConfig defines distributed tracing settings
type TracingConfig struct {
	Name        string             `json:"name,omitempty"`
	ServiceName string             `json:"serviceName,omitempty"`
	Sampler     *TracingSampler    `json:"sampler,omitempty"`
	Config      *TracingConfigData `json:"config,omitempty"`
}

// TracingSampler defines tracing sampler settings
type TracingSampler struct {
	Type string `json:"type,omitempty"`
	// Param expressed as string to avoid float in CRD schema
	Param string `json:"param,omitempty"`
}

// TracingConfigData defines tracing backend configuration
type TracingConfigData struct {
	URL      string            `json:"url,omitempty"`
	Endpoint string            `json:"endpoint,omitempty"`
	Headers  map[string]string `json:"headers,omitempty"`
}

// TimeoutConfig defines global timeout settings
type TimeoutConfig struct {
	Connect string `json:"connect,omitempty"`
	Request string `json:"request,omitempty"`
}

// PixiuGatewayPolicyStatus defines the observed state of PixiuGatewayPolicy
type PixiuGatewayPolicyStatus struct {
	// Conditions describe the current conditions of the PixiuGatewayPolicy.
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:shortName=pgp

// PixiuGatewayPolicy is the Schema for the pixiugatewaypolicies API
type PixiuGatewayPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PixiuGatewayPolicySpec   `json:"spec,omitempty"`
	Status PixiuGatewayPolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PixiuGatewayPolicyList contains a list of PixiuGatewayPolicy
type PixiuGatewayPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PixiuGatewayPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PixiuGatewayPolicy{}, &PixiuGatewayPolicyList{})
}
