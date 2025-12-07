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
	"errors"
	"fmt"
	"reflect"
)

import (
	"controllers/internal/controller/status"

	"controllers/internal/converter"

	"controllers/internal/ir"

	"controllers/internal/translator"

	"controllers/internal/utils"

	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
)

// GatewayReconciler reconciles a Gateway object.
type GatewayReconciler struct { //nolint:revive
	client.Client
	Scheme  *runtime.Scheme
	Log     logr.Logger
	Updater status.Updater
}

// SetupWithManager sets up the controller with the Manager.
func (r *GatewayReconciler) SetupWithManager(mgr ctrl.Manager) error {
	bdr := ctrl.NewControllerManagedBy(mgr).
		For(
			&gatewayv1.Gateway{},
			builder.WithPredicates(
				predicate.NewPredicateFuncs(r.checkGatewayClass),
			),
		).
		WithEventFilter(
			predicate.Or(
				predicate.GenerationChangedPredicate{},
				predicate.NewPredicateFuncs(TypePredicate[*corev1.Secret]()),
			),
		).
		Watches(
			&gatewayv1.GatewayClass{},
			handler.EnqueueRequestsFromMapFunc(r.listGatewayForGatewayClass),
			builder.WithPredicates(
				predicate.NewPredicateFuncs(r.matchesGatewayClass),
			),
		).
		Watches(
			&gatewayv1.HTTPRoute{},
			handler.EnqueueRequestsFromMapFunc(r.listGatewaysForHTTPRoute),
		).
		Watches(
			&appsv1.Deployment{},
			handler.EnqueueRequestsFromMapFunc(r.listGatewaysForDeployment),
		).
		Watches(
			&corev1.Service{},
			handler.EnqueueRequestsFromMapFunc(r.listGatewaysForService),
		)

	if GetEnableReferenceGrant() {
		bdr.Watches(&v1beta1.ReferenceGrant{},
			handler.EnqueueRequestsFromMapFunc(r.listReferenceGrantsForGateway),
			builder.WithPredicates(referenceGrantPredicates(KindGateway)),
		)
	}

	return bdr.Complete(r)
}

func (r *GatewayReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	gateway := new(gatewayv1.Gateway)
	if err := r.Get(ctx, req.NamespacedName, gateway); err != nil {
		if client.IgnoreNotFound(err) == nil {
			gateway.Namespace = req.Namespace
			gateway.Name = req.Name

			gateway.TypeMeta = metav1.TypeMeta{
				Kind:       KindGateway,
				APIVersion: gatewayv1.GroupVersion.String(),
			}
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}
	conditionProgrammedStatus, conditionProgrammedMsg := true, "Programmed"

	r.Log.Info("gateway has been accepted", "gateway", gateway.GetName())
	type conditionStatus struct {
		status bool
		msg    string
	}
	acceptStatus := conditionStatus{
		status: true,
		msg:    acceptedMessage("gateway"),
	}

	r.processListenerConfig(gateway)

	if err := r.ensureGatewayConfigMap(ctx, gateway); err != nil {
		r.Log.Error(err, "failed to ensure gateway configmap", "gateway", gateway.GetName())
		conditionProgrammedStatus = false
		conditionProgrammedMsg = fmt.Sprintf("Failed to create configmap: %v", err)
	}

	if err := r.ensureDataPlane(ctx, gateway); err != nil {
		r.Log.Error(err, "failed to ensure data plane", "gateway", gateway.GetName())
		conditionProgrammedStatus = false
		conditionProgrammedMsg = fmt.Sprintf("Failed to create data plane: %v", err)
	}

	if err := r.updateGatewayAddresses(ctx, gateway); err != nil {
		r.Log.Error(err, "failed to update gateway addresses", "gateway", gateway.GetName())
	}

	listenerStatuses, err := getListenerStatus(ctx, r.Client, gateway)
	if err != nil {
		r.Log.Error(err, "failed to get listener status", "gateway", req.NamespacedName)
		return ctrl.Result{}, err
	}

	accepted := SetGatewayConditionAccepted(gateway, acceptStatus.status, acceptStatus.msg)
	programmed := SetGatewayConditionProgrammed(gateway, conditionProgrammedStatus, conditionProgrammedMsg)

	needsUpdate := accepted || programmed || len(listenerStatuses) > 0 || len(gateway.Status.Addresses) > 0
	if needsUpdate {
		if len(listenerStatuses) > 0 {
			gateway.Status.Listeners = listenerStatuses
		}

		r.Updater.Update(status.Update{
			NamespacedName: utils.NamespacedName(gateway),
			Resource:       &gatewayv1.Gateway{},
			Mutator: status.MutatorFunc(func(obj client.Object) client.Object {
				t, ok := obj.(*gatewayv1.Gateway)
				if !ok {
					err := fmt.Errorf("unsupported object type %T", obj)
					panic(err)
				}
				tCopy := t.DeepCopy()
				tCopy.Status = gateway.Status
				return tCopy
			}),
		})
	}

	return ctrl.Result{}, nil
}

func (r *GatewayReconciler) matchesGatewayClass(obj client.Object) bool {
	gateway, ok := obj.(*gatewayv1.GatewayClass)
	if !ok {
		r.Log.Error(fmt.Errorf("unexpected object type"), "failed to convert object to Gateway")
		return false
	}
	return matchesController(string(gateway.Spec.ControllerName))
}

func (r *GatewayReconciler) listGatewayForGatewayClass(ctx context.Context, gatewayClass client.Object) []reconcile.Request {
	gatewayList := &gatewayv1.GatewayList{}
	if err := r.List(context.Background(), gatewayList); err != nil {
		r.Log.Error(err, "failed to list gateways for gateway class",
			"gatewayclass", gatewayClass.GetName(),
		)
		return nil
	}

	return reconcileGatewaysMatchGatewayClass(gatewayClass, gatewayList.Items)
}

func (r *GatewayReconciler) checkGatewayClass(obj client.Object) bool {
	gateway := obj.(*gatewayv1.Gateway)
	gatewayClass := &gatewayv1.GatewayClass{}
	if err := r.Get(context.Background(), client.ObjectKey{Name: string(gateway.Spec.GatewayClassName)}, gatewayClass); err != nil {
		r.Log.Error(err, "failed to get gateway class", "gateway", gateway.GetName(), "gatewayclass", gateway.Spec.GatewayClassName)
		return false
	}

	return matchesController(string(gatewayClass.Spec.ControllerName))
}

func (r *GatewayReconciler) listGatewaysForHTTPRoute(ctx context.Context, obj client.Object) []reconcile.Request {
	httpRoute, ok := obj.(*gatewayv1.HTTPRoute)
	if !ok {
		r.Log.Error(
			fmt.Errorf("unexpected object type"),
			"HTTPRoute watch predicate received unexpected object type",
			"expected", "*gatewayapi.HTTPRoute", "found", reflect.TypeOf(obj),
		)
		return nil
	}
	recs := []reconcile.Request{}
	for _, parentRef := range httpRoute.Spec.ParentRefs {
		if parentRef.Group != nil && *parentRef.Group != gatewayv1.GroupName {
			continue
		}
		if parentRef.Kind != nil && *parentRef.Kind != "Gateway" {
			continue
		}

		gatewayNamespace := httpRoute.GetNamespace()
		if parentRef.Namespace != nil {
			gatewayNamespace = string(*parentRef.Namespace)
		}

		gateway := new(gatewayv1.Gateway)
		if err := r.Get(ctx, client.ObjectKey{
			Namespace: gatewayNamespace,
			Name:      string(parentRef.Name),
		}, gateway); err != nil {
			continue
		}

		if !r.checkGatewayClass(gateway) {
			continue
		}

		recs = append(recs, reconcile.Request{
			NamespacedName: client.ObjectKey{
				Namespace: gatewayNamespace,
				Name:      string(parentRef.Name),
			},
		})
	}
	return recs
}

func (r *GatewayReconciler) listReferenceGrantsForGateway(ctx context.Context, obj client.Object) (requests []reconcile.Request) {
	grant, ok := obj.(*v1beta1.ReferenceGrant)
	if !ok {
		r.Log.Error(
			errors.New("unexpected object type"),
			"ReferenceGrant watch predicate received unexpected object type",
			"expected", FullTypeName(new(v1beta1.ReferenceGrant)), "found", FullTypeName(obj),
		)
		return nil
	}

	var gatewayList gatewayv1.GatewayList
	if err := r.List(ctx, &gatewayList); err != nil {
		r.Log.Error(err, "failed to list gateways in watch predicate", "ReferenceGrant", grant.GetName())
		return nil
	}

	for _, gateway := range gatewayList.Items {
		gw := v1beta1.ReferenceGrantFrom{
			Group:     gatewayv1.GroupName,
			Kind:      KindGateway,
			Namespace: v1beta1.Namespace(gateway.GetNamespace()),
		}
		for _, from := range grant.Spec.From {
			if from == gw {
				requests = append(requests, reconcile.Request{
					NamespacedName: types.NamespacedName{
						Namespace: gateway.GetNamespace(),
						Name:      gateway.GetName(),
					},
				})
			}
		}
	}
	return requests
}

func (r *GatewayReconciler) processListenerConfig(gateway *gatewayv1.Gateway) {
	listeners := gateway.Spec.Listeners
	for _, listener := range listeners {
		if listener.TLS == nil || listener.TLS.CertificateRefs == nil {
			continue
		}
		secret := corev1.Secret{}
		for _, ref := range listener.TLS.CertificateRefs {
			ns := gateway.GetNamespace()
			if ref.Namespace != nil {
				ns = string(*ref.Namespace)
			}
			if ref.Kind != nil && *ref.Kind == KindSecret {
				if err := r.Get(context.Background(), client.ObjectKey{
					Namespace: ns,
					Name:      string(ref.Name),
				}, &secret); err != nil {
					r.Log.Error(err, "failed to get secret", "namespace", ns, "name", ref.Name)
					SetGatewayListenerConditionProgrammed(gateway, string(listener.Name), false, err.Error())
					SetGatewayListenerConditionResolvedRefs(gateway, string(listener.Name), false, err.Error())
					break
				}
			}
		}
	}
}

func (r *GatewayReconciler) ensureGatewayConfigMap(ctx context.Context, gateway *gatewayv1.Gateway) error {
	configMapName := fmt.Sprintf("%s-config", gateway.GetName())

	translator := translator.NewTranslator(r.Client, r.Log)
	conv := converter.NewConverter()

	xds, err := translator.TranslateGateway(ctx, gateway)
	if err != nil {
		return fmt.Errorf("failed to translate Gateway to IR: %w", err)
	}

	allClusters := []*ir.Cluster{}
	clusterMap := make(map[string]*ir.Cluster)

	var httpRouteList gatewayv1.HTTPRouteList
	if err := r.Client.List(ctx, &httpRouteList); err == nil {
		matchedHTTPRoutes := 0
		for _, httpRoute := range httpRouteList.Items {
			attached := false
			for _, parentRef := range httpRoute.Spec.ParentRefs {
				if parentRef.Group != nil && *parentRef.Group != gatewayv1.GroupName {
					continue
				}
				if parentRef.Kind != nil && *parentRef.Kind != "Gateway" {
					continue
				}
				if string(parentRef.Name) != gateway.Name {
					continue
				}
				namespaceMatch := false
				if parentRef.Namespace != nil {
					namespaceMatch = string(*parentRef.Namespace) == gateway.Namespace
				} else {
					namespaceMatch = gateway.Namespace == httpRoute.Namespace
				}
				if !namespaceMatch {
					continue
				}
				attached = true
				break
			}
			if !attached {
				continue
			}

			matchedHTTPRoutes++
			clusters, err := translator.TranslateBackendToCluster(ctx, &httpRoute)
			if err != nil {
				r.Log.Error(err, "failed to translate backend to cluster", "httproute", httpRoute.Name)
				continue
			}
			for _, cluster := range clusters {
				if existing, exists := clusterMap[cluster.Name]; exists {
					existing.Endpoints = append(existing.Endpoints, cluster.Endpoints...)
				} else {
					clusterMap[cluster.Name] = cluster
					allClusters = append(allClusters, cluster)
				}
			}
		}
	} else {
		r.Log.Error(err, "failed to list HTTPRoutes")
	}

	pixiuConfig, err := conv.ConvertIRToPixiuConfig(xds, allClusters)
	if err != nil {
		return fmt.Errorf("failed to convert IR to Pixiu config: %w", err)
	}

	configYAML, err := converter.ConvertToYAML(pixiuConfig)
	if err != nil {
		return fmt.Errorf("failed to convert config to YAML: %w", err)
	}

	existingConfigMap := &corev1.ConfigMap{}
	err = r.Get(ctx, client.ObjectKey{
		Name:      configMapName,
		Namespace: gateway.GetNamespace(),
	}, existingConfigMap)

	configMap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configMapName,
			Namespace: gateway.GetNamespace(),
			Labels: map[string]string{
				"app.kubernetes.io/managed-by":           "pg-controller",
				"gateway.networking.k8s.io/gateway-name": gateway.GetName(),
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: gateway.APIVersion,
					Kind:       gateway.Kind,
					Name:       gateway.GetName(),
					UID:        gateway.GetUID(),
					Controller: func() *bool { b := true; return &b }(),
				},
			},
		},
		Data: map[string]string{
			"conf.yaml": configYAML,
		},
	}

	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			if err := r.Create(ctx, configMap); err != nil {
				return fmt.Errorf("failed to create configmap: %w", err)
			}
			r.Log.Info("created gateway configmap", "gateway", gateway.GetName(), "configmap", configMapName)
		} else {
			return fmt.Errorf("failed to check configmap: %w", err)
		}
	} else {
		if existingConfigMap.Data["conf.yaml"] != configYAML {
			existingConfigMap.Data = configMap.Data
			existingConfigMap.Labels = configMap.Labels
			if err := r.Update(ctx, existingConfigMap); err != nil {
				return fmt.Errorf("failed to update configmap: %w", err)
			}
			r.Log.Info("updated gateway configmap", "gateway", gateway.GetName(), "configmap", configMapName)
		}
	}

	return nil
}

func (r *GatewayReconciler) ensureDataPlane(ctx context.Context, gateway *gatewayv1.Gateway) error {
	configMapName := fmt.Sprintf("%s-config", gateway.GetName())
	deploymentName := fmt.Sprintf("%s-%s", gateway.GetName(), string(gateway.GetUID())[:8])
	serviceName := deploymentName

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: gateway.GetNamespace(),
			Labels: map[string]string{
				"app.kubernetes.io/name":                 "pixiu-gateway",
				"app.kubernetes.io/managed-by":           "pg-controller",
				"gateway.networking.k8s.io/gateway-name": gateway.GetName(),
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: gateway.APIVersion,
					Kind:       gateway.Kind,
					Name:       gateway.GetName(),
					UID:        gateway.GetUID(),
					Controller: func() *bool { b := true; return &b }(),
				},
			},
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: func() *int32 { r := int32(1); return &r }(),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app.kubernetes.io/name":                 "pixiu-gateway",
					"gateway.networking.k8s.io/gateway-name": gateway.GetName(),
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app.kubernetes.io/name":                 "pixiu-gateway",
						"gateway.networking.k8s.io/gateway-name": gateway.GetName(),
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "pixiu",
							Image: "mfordjody/pixiugateway:debug",
							Args: []string{
								"gateway",
								"start",
								"-c",
								"/etc/pixiu/conf.yaml",
							},
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8888,
									Protocol:      corev1.ProtocolTCP,
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "config",
									MountPath: "/etc/pixiu",
									ReadOnly:  true,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "config",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: configMapName,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	existingDeployment := &appsv1.Deployment{}
	err := r.Get(ctx, client.ObjectKey{
		Name:      deploymentName,
		Namespace: gateway.GetNamespace(),
	}, existingDeployment)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			if err := r.Create(ctx, deployment); err != nil {
				return fmt.Errorf("failed to create deployment: %w", err)
			}
			r.Log.Info("created data plane deployment", "gateway", gateway.GetName(), "deployment", deploymentName)
		} else {
			return fmt.Errorf("failed to get deployment: %w", err)
		}
	} else {
		needsUpdate := false
		if !reflect.DeepEqual(existingDeployment.Spec, deployment.Spec) {
			needsUpdate = true
		}
		if !reflect.DeepEqual(existingDeployment.Labels, deployment.Labels) {
			needsUpdate = true
		}

		if needsUpdate {
			patch := client.MergeFrom(existingDeployment.DeepCopy())
			existingDeployment.Spec = deployment.Spec
			existingDeployment.Labels = deployment.Labels
			if err := r.Patch(ctx, existingDeployment, patch); err != nil {
				return fmt.Errorf("failed to patch deployment: %w", err)
			}
			r.Log.Info("updated data plane deployment", "gateway", gateway.GetName(), "deployment", deploymentName)
		}
	}

	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: gateway.GetNamespace(),
			Labels: map[string]string{
				"app.kubernetes.io/name":                 "pixiu-gateway",
				"app.kubernetes.io/managed-by":           "pg-controller",
				"gateway.networking.k8s.io/gateway-name": gateway.GetName(),
			},
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: gateway.APIVersion,
					Kind:       gateway.Kind,
					Name:       gateway.GetName(),
					UID:        gateway.GetUID(),
					Controller: func() *bool { b := true; return &b }(),
				},
			},
		},
		Spec: corev1.ServiceSpec{
			Type: corev1.ServiceTypeLoadBalancer,
			Selector: map[string]string{
				"app.kubernetes.io/name":                 "pixiu-gateway",
				"gateway.networking.k8s.io/gateway-name": gateway.GetName(),
			},
			Ports: []corev1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(8888),
					Protocol:   corev1.ProtocolTCP,
				},
			},
		},
	}

	existingService := &corev1.Service{}
	err = r.Get(ctx, client.ObjectKey{
		Name:      serviceName,
		Namespace: gateway.GetNamespace(),
	}, existingService)
	if err != nil {
		if client.IgnoreNotFound(err) == nil {
			if err := r.Create(ctx, service); err != nil {
				return fmt.Errorf("failed to create service: %w", err)
			}
			r.Log.Info("created data plane service", "gateway", gateway.GetName(), "service", serviceName)
		} else {
			return fmt.Errorf("failed to get service: %w", err)
		}
	} else {
		service.Spec.ClusterIP = existingService.Spec.ClusterIP
		service.Spec.ClusterIPs = existingService.Spec.ClusterIPs
		existingService.Spec = service.Spec
		existingService.Labels = service.Labels
		if err := r.Update(ctx, existingService); err != nil {
			return fmt.Errorf("failed to update service: %w", err)
		}
		r.Log.Info("updated data plane service", "gateway", gateway.GetName(), "service", serviceName)
	}

	return nil
}

func (r *GatewayReconciler) updateGatewayAddresses(ctx context.Context, gateway *gatewayv1.Gateway) error {
	deploymentName := fmt.Sprintf("%s-%s", gateway.GetName(), string(gateway.GetUID())[:8])
	serviceName := deploymentName

	service := &corev1.Service{}
	if err := r.Get(ctx, client.ObjectKey{
		Name:      serviceName,
		Namespace: gateway.GetNamespace(),
	}, service); err != nil {
		if client.IgnoreNotFound(err) == nil {
			gateway.Status.Addresses = nil
			return nil
		}
		return fmt.Errorf("failed to get service: %w", err)
	}

	var addresses []gatewayv1.GatewayStatusAddress
	if service.Spec.Type == corev1.ServiceTypeLoadBalancer {
		if len(service.Status.LoadBalancer.Ingress) > 0 {
			for _, ingress := range service.Status.LoadBalancer.Ingress {
				if ingress.IP != "" {
					addresses = append(addresses, gatewayv1.GatewayStatusAddress{
						Type:  func() *gatewayv1.AddressType { t := gatewayv1.IPAddressType; return &t }(),
						Value: ingress.IP,
					})
				} else if ingress.Hostname != "" {
					addresses = append(addresses, gatewayv1.GatewayStatusAddress{
						Type:  func() *gatewayv1.AddressType { t := gatewayv1.HostnameAddressType; return &t }(),
						Value: ingress.Hostname,
					})
				}
			}
		}
	}

	if len(addresses) == 0 && service.Spec.ClusterIP != "" && service.Spec.ClusterIP != "None" {
		addresses = append(addresses, gatewayv1.GatewayStatusAddress{
			Type:  func() *gatewayv1.AddressType { t := gatewayv1.IPAddressType; return &t }(),
			Value: service.Spec.ClusterIP,
		})
	}

	gateway.Status.Addresses = addresses
	return nil
}

func (r *GatewayReconciler) listGatewaysForDeployment(ctx context.Context, obj client.Object) []reconcile.Request {
	deployment, ok := obj.(*appsv1.Deployment)
	if !ok {
		return nil
	}

	if len(deployment.OwnerReferences) == 0 {
		return nil
	}

	for _, ownerRef := range deployment.OwnerReferences {
		if ownerRef.Kind == "Gateway" && ownerRef.APIVersion == gatewayv1.GroupVersion.String() {
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: deployment.GetNamespace(),
						Name:      ownerRef.Name,
					},
				},
			}
		}
	}

	return nil
}

func (r *GatewayReconciler) listGatewaysForService(ctx context.Context, obj client.Object) []reconcile.Request {
	service, ok := obj.(*corev1.Service)
	if !ok {
		return nil
	}

	if len(service.OwnerReferences) == 0 {
		return nil
	}

	for _, ownerRef := range service.OwnerReferences {
		if ownerRef.Kind == "Gateway" && ownerRef.APIVersion == gatewayv1.GroupVersion.String() {
			return []reconcile.Request{
				{
					NamespacedName: types.NamespacedName{
						Namespace: service.GetNamespace(),
						Name:      ownerRef.Name,
					},
				},
			}
		}
	}

	return nil
}
