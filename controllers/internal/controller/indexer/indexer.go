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

package indexer

import (
	"context"
)

import (
	"controllers/api/v1alpha1"

	networkingv1 "k8s.io/api/networking/v1"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	ServiceIndexRef           = "serviceRefs"
	ParametersRef             = "parametersRef"
	IngressClass              = "ingressClass"
	IngressClassRef           = "ingressClassRef"
	IngressClassParametersRef = "ingressClassParametersRef"
	SecretIndexRef            = "secretRefs"
	GatewayClassIndexRef      = "gatewayClassRef"
	ControllerName            = "controllerName"
)

func SetupIndexer(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		setupGatewayClassIndexer,
		setupGatewayIndexer,
		setupIngressIndexer,
		setupIngressClassIndexer,
		setupGatewayProxyIndexer,
	} {
		if err := setup(mgr); err != nil {
			return err
		}
	}
	return nil
}

func GenIndexKey(namespace, name string) string {
	return client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}.String()
}

func setupGatewayIndexer(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&gatewayv1.Gateway{},
		ParametersRef,
		GatewayParametersRefIndexFunc,
	); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&gatewayv1.Gateway{},
		GatewayClassIndexRef,
		func(obj client.Object) (requests []string) {
			return []string{string(obj.(*gatewayv1.Gateway).Spec.GatewayClassName)}
		},
	); err != nil {
		return err
	}
	return nil
}

func setupIngressIndexer(mgr ctrl.Manager) error {
	// create IngressClass index
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&networkingv1.Ingress{},
		IngressClassRef,
		IngressClassRefIndexFunc,
	); err != nil {
		return err
	}

	// create Service index for quick lookup of Ingresses using specific services
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&networkingv1.Ingress{},
		ServiceIndexRef,
		IngressServiceIndexFunc,
	); err != nil {
		return err
	}

	// create secret index for TLS
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&networkingv1.Ingress{},
		SecretIndexRef,
		IngressSecretIndexFunc,
	); err != nil {
		return err
	}

	return nil
}

func setupIngressClassIndexer(mgr ctrl.Manager) error {
	// create IngressClass index
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&networkingv1.IngressClass{},
		IngressClass,
		IngressClassIndexFunc,
	); err != nil {
		return err
	}

	// create IngressClassParametersRef index
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&networkingv1.IngressClass{},
		IngressClassParametersRef,
		IngressClassParametersRefIndexFunc,
	); err != nil {
		return err
	}

	return nil
}

func setupGatewayProxyIndexer(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&v1alpha1.GatewayProxy{},
		ServiceIndexRef,
		GatewayProxyServiceIndexFunc,
	); err != nil {
		return err
	}

	if err := mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&v1alpha1.GatewayProxy{},
		SecretIndexRef,
		GatewayProxySecretIndexFunc,
	); err != nil {
		return err
	}
	return nil
}

func setupGatewayClassIndexer(mgr ctrl.Manager) error {
	return mgr.GetFieldIndexer().IndexField(
		context.Background(),
		&gatewayv1.GatewayClass{},
		ControllerName,
		func(obj client.Object) []string {
			return []string{string(obj.(*gatewayv1.GatewayClass).Spec.ControllerName)}
		},
	)
}

func GatewayProxyServiceIndexFunc(rawObj client.Object) []string {
	return nil
}

func GatewayProxySecretIndexFunc(rawObj client.Object) []string {
	secretKeys := make([]string, 0)
	return secretKeys
}

func GatewayParametersRefIndexFunc(rawObj client.Object) []string {
	gw := rawObj.(*gatewayv1.Gateway)
	if gw.Spec.Infrastructure != nil && gw.Spec.Infrastructure.ParametersRef != nil {
		// now we only care about kind: GatewayProxy
		if gw.Spec.Infrastructure.ParametersRef.Kind == "GatewayProxy" {
			name := gw.Spec.Infrastructure.ParametersRef.Name
			return []string{GenIndexKey(gw.GetNamespace(), name)}
		}
	}
	return nil
}

func IngressClassIndexFunc(rawObj client.Object) []string {
	ingressClass := rawObj.(*networkingv1.IngressClass)
	if ingressClass.Spec.Controller == "" {
		return nil
	}
	controllerName := ingressClass.Spec.Controller
	return []string{controllerName}
}

func IngressClassRefIndexFunc(rawObj client.Object) []string {
	ingress := rawObj.(*networkingv1.Ingress)
	if ingress.Spec.IngressClassName == nil {
		return nil
	}
	return []string{*ingress.Spec.IngressClassName}
}

func IngressClassParametersRefIndexFunc(rawObj client.Object) []string {
	ingressClass := rawObj.(*networkingv1.IngressClass)
	// check if the IngressClass references this gateway proxy
	if ingressClass.Spec.Parameters != nil &&
		ingressClass.Spec.Parameters.APIGroup != nil &&
		*ingressClass.Spec.Parameters.APIGroup == v1alpha1.GroupVersion.Group &&
		ingressClass.Spec.Parameters.Kind == "GatewayProxy" {
		ns := ingressClass.GetNamespace()
		if ingressClass.Spec.Parameters.Namespace != nil {
			ns = *ingressClass.Spec.Parameters.Namespace
		}
		return []string{GenIndexKey(ns, ingressClass.Spec.Parameters.Name)}
	}
	return nil
}

func IngressServiceIndexFunc(rawObj client.Object) []string {
	ingress := rawObj.(*networkingv1.Ingress)
	var services []string

	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}

		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service == nil {
				continue
			}
			key := GenIndexKey(ingress.Namespace, path.Backend.Service.Name)
			services = append(services, key)
		}
	}
	return services
}

func IngressSecretIndexFunc(rawObj client.Object) []string {
	ingress := rawObj.(*networkingv1.Ingress)
	secrets := make([]string, 0)

	for _, tls := range ingress.Spec.TLS {
		if tls.SecretName == "" {
			continue
		}

		key := GenIndexKey(ingress.Namespace, tls.SecretName)
		secrets = append(secrets, key)
	}
	return secrets
}
