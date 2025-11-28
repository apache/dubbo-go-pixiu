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

package types

import (
	"controllers/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"

	netv1 "k8s.io/api/networking/v1"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	KindGateway      = "Gateway"
	KindGatewayClass = "GatewayClass"
	KindIngress      = "Ingress"
	KindSecret       = "Secret"
	KindIngressClass = "IngressClass"
	KindGatewayProxy = "GatewayProxy"
)

func KindOf(obj any) string {
	switch obj.(type) {
	case *gatewayv1.Gateway:
		return KindGateway
	case *gatewayv1.GatewayClass:
		return KindGatewayClass
	case *netv1.Ingress:
		return KindIngress
	case *netv1.IngressClass:
		return KindIngressClass
	case *corev1.Secret:
		return KindSecret
	case *v1alpha1.GatewayProxy:
		return KindGatewayProxy
	default:
		return "Unknown"
	}
}
