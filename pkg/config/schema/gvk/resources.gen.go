// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// GENERATED FILE -- DO NOT EDIT
//

package gvk

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
)

var (
	AuthorizationPolicy          = config.GroupVersionKind{Group: "security.istio.io", Version: "v1beta1", Kind: "AuthorizationPolicy"}
	ConfigMap                    = config.GroupVersionKind{Group: "", Version: "v1", Kind: "ConfigMap"}
	CustomResourceDefinition     = config.GroupVersionKind{Group: "apiextensions.k8s.io", Version: "v1", Kind: "CustomResourceDefinition"}
	Deployment                   = config.GroupVersionKind{Group: "apps", Version: "v1", Kind: "Deployment"}
	DestinationRule              = config.GroupVersionKind{Group: "networking.istio.io", Version: "v1alpha3", Kind: "DestinationRule"}
	Endpoints                    = config.GroupVersionKind{Group: "", Version: "v1", Kind: "Endpoints"}
	EnvoyFilter                  = config.GroupVersionKind{Group: "networking.istio.io", Version: "v1alpha3", Kind: "EnvoyFilter"}
	Gateway                      = config.GroupVersionKind{Group: "networking.istio.io", Version: "v1alpha3", Kind: "Gateway"}
	GatewayClass                 = config.GroupVersionKind{Group: "gateway.networking.k8s.io", Version: "v1alpha2", Kind: "GatewayClass"}
	HTTPRoute                    = config.GroupVersionKind{Group: "gateway.networking.k8s.io", Version: "v1alpha2", Kind: "HTTPRoute"}
	Ingress                      = config.GroupVersionKind{Group: "extensions", Version: "v1beta1", Kind: "Ingress"}
	KubernetesGateway            = config.GroupVersionKind{Group: "gateway.networking.k8s.io", Version: "v1alpha2", Kind: "Gateway"}
	MeshConfig                   = config.GroupVersionKind{Group: "", Version: "v1alpha1", Kind: "MeshConfig"}
	MeshNetworks                 = config.GroupVersionKind{Group: "", Version: "v1alpha1", Kind: "MeshNetworks"}
	MutatingWebhookConfiguration = config.GroupVersionKind{Group: "admissionregistration.k8s.io", Version: "v1", Kind: "MutatingWebhookConfiguration"}
	Namespace                    = config.GroupVersionKind{Group: "", Version: "v1", Kind: "Namespace"}
	Node                         = config.GroupVersionKind{Group: "", Version: "v1", Kind: "Node"}
	PeerAuthentication           = config.GroupVersionKind{Group: "security.istio.io", Version: "v1beta1", Kind: "PeerAuthentication"}
	Pod                          = config.GroupVersionKind{Group: "", Version: "v1", Kind: "Pod"}
	ProxyConfig                  = config.GroupVersionKind{Group: "networking.istio.io", Version: "v1beta1", Kind: "ProxyConfig"}
	ReferenceGrant               = config.GroupVersionKind{Group: "gateway.networking.k8s.io", Version: "v1alpha2", Kind: "ReferenceGrant"}
	ReferencePolicy              = config.GroupVersionKind{Group: "gateway.networking.k8s.io", Version: "v1alpha2", Kind: "ReferencePolicy"}
	RequestAuthentication        = config.GroupVersionKind{Group: "security.istio.io", Version: "v1beta1", Kind: "RequestAuthentication"}
	Secret                       = config.GroupVersionKind{Group: "", Version: "v1", Kind: "Secret"}
	Service                      = config.GroupVersionKind{Group: "", Version: "v1", Kind: "Service"}
	ServiceEntry                 = config.GroupVersionKind{Group: "networking.istio.io", Version: "v1alpha3", Kind: "ServiceEntry"}
	Sidecar                      = config.GroupVersionKind{Group: "networking.istio.io", Version: "v1alpha3", Kind: "Sidecar"}
	TCPRoute                     = config.GroupVersionKind{Group: "gateway.networking.k8s.io", Version: "v1alpha2", Kind: "TCPRoute"}
	TLSRoute                     = config.GroupVersionKind{Group: "gateway.networking.k8s.io", Version: "v1alpha2", Kind: "TLSRoute"}
	Telemetry                    = config.GroupVersionKind{Group: "telemetry.istio.io", Version: "v1alpha1", Kind: "Telemetry"}
	VirtualService               = config.GroupVersionKind{Group: "networking.istio.io", Version: "v1alpha3", Kind: "VirtualService"}
	WasmPlugin                   = config.GroupVersionKind{Group: "extensions.istio.io", Version: "v1alpha1", Kind: "WasmPlugin"}
	WorkloadEntry                = config.GroupVersionKind{Group: "networking.istio.io", Version: "v1alpha3", Kind: "WorkloadEntry"}
	WorkloadGroup                = config.GroupVersionKind{Group: "networking.istio.io", Version: "v1alpha3", Kind: "WorkloadGroup"}
)