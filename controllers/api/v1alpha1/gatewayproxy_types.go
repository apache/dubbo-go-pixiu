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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.
// GatewayProxySpec defines the desired state of GatewayProxy
type GatewayProxySpec struct{}

// GatewayProxyStatus defines the observed state of GatewayProxy.

type GatewayProxyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// GatewayProxy is the Schema for the gatewayproxies API
type GatewayProxy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// GatewayProxySpec defines configuration of gateway proxy instances,
	// including networking settings, global plugins, and plugin metadata.
	Spec GatewayProxySpec `json:"spec,omitempty"`
}

// +kubebuilder:object:root=true

// GatewayProxyList contains a list of GatewayProxy
type GatewayProxyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GatewayProxy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GatewayProxy{}, &GatewayProxyList{})
}
