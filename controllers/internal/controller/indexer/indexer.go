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
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
)

const (
	ServiceIndexRef      = "serviceRefs"
	ParametersRef        = "parametersRef"
	SecretIndexRef       = "secretRefs"
	GatewayClassIndexRef = "gatewayClassRef"
	ControllerName       = "controllerName"
)

func SetupIndexer(mgr ctrl.Manager) error {
	for _, setup := range []func(ctrl.Manager) error{
		setupGatewayClassIndexer,
		setupGatewayIndexer,
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

func GatewayParametersRefIndexFunc(rawObj client.Object) []string {
	// Infrastructure parameters are no longer used
	return nil
}
