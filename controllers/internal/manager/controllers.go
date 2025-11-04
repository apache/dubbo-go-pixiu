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

package manager

import (
	"context"
)

import (
	"controllers/internal/controller"
	"controllers/internal/controller/indexer"
	"controllers/internal/controller/status"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// K8s
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch
// +kubebuilder:rbac:groups=coordination.k8s.io,resources=leases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="discovery.k8s.io",resources=endpointslices,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=services,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch
// GatewayAPI
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gatewayclasses/status,verbs=get;update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=gateways/status,verbs=get;update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes,verbs=get;list;watch
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=httproutes/status,verbs=get;update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=referencegrants,verbs=list;watch;update
// +kubebuilder:rbac:groups=gateway.networking.k8s.io,resources=referencegrants/status,verbs=get;update
// Networking
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses,verbs=get;list;watch;update
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingresses/status,verbs=get;update
// +kubebuilder:rbac:groups=networking.k8s.io,resources=ingressclasses,verbs=get;list;watch
type Controller interface {
	SetupWithManager(mgr manager.Manager) error
}

func setupControllers(ctx context.Context, mgr manager.Manager, updater status.Updater) ([]Controller, error) {
	if err := indexer.SetupIndexer(mgr); err != nil {
		return nil, err
	}
	return []Controller{
		&controller.GatewayClassReconciler{
			Client:  mgr.GetClient(),
			Scheme:  mgr.GetScheme(),
			Log:     ctrl.LoggerFrom(ctx).WithName("controllers").WithName("GatewayClass"),
			Updater: updater,
		},
		&controller.GatewayReconciler{
			Client:  mgr.GetClient(),
			Scheme:  mgr.GetScheme(),
			Log:     ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Gateway"),
			Updater: updater,
		},
		&controller.IngressReconciler{
			Client:  mgr.GetClient(),
			Scheme:  mgr.GetScheme(),
			Log:     ctrl.LoggerFrom(ctx).WithName("controllers").WithName("Ingress"),
			Updater: updater,
		},
		&controller.IngressClassReconciler{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
			Log:    ctrl.LoggerFrom(ctx).WithName("controllers").WithName("IngressClass"),
		},
		&controller.GatewayProxyController{
			Client: mgr.GetClient(),
			Scheme: mgr.GetScheme(),
			Log:    ctrl.LoggerFrom(ctx).WithName("controllers").WithName("GatewayProxy"),
		},
	}, nil
}
