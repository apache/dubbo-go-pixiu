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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"time"
)

import (
	"controllers/api/v1alpha1"

	"controllers/internal/controller/config"
	"controllers/internal/controller/status"

	"controllers/internal/converter"

	"controllers/internal/ir"

	"controllers/internal/translator"

	"controllers/internal/utils"

	"github.com/go-logr/logr"

	appsv1 "k8s.io/api/apps/v1"

	corev1 "k8s.io/api/core/v1"

	discoveryv1 "k8s.io/api/discovery/v1"

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
		).
		Watches(
			&v1alpha1.PixiuFilterPolicy{},
			handler.EnqueueRequestsFromMapFunc(r.listGatewaysForFilterPolicy),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
		).
		Watches(
			&v1alpha1.PixiuClusterPolicy{},
			handler.EnqueueRequestsFromMapFunc(r.listGatewaysForClusterPolicy),
			builder.WithPredicates(predicate.GenerationChangedPredicate{}),
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

	configHash, err := r.ensureGatewayConfigMap(ctx, gateway)
	if err != nil {
		r.Log.Error(err, "failed to ensure gateway configmap", "gateway", gateway.GetName())
		conditionProgrammedStatus = false
		conditionProgrammedMsg = fmt.Sprintf("Failed to create configmap: %v", err)
	}

	if err := r.ensureDataPlane(ctx, gateway, configHash); err != nil {
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

func (r *GatewayReconciler) listGatewaysForFilterPolicy(ctx context.Context, obj client.Object) []reconcile.Request {
	policy, ok := obj.(*v1alpha1.PixiuFilterPolicy)
	if !ok {
		r.Log.Error(fmt.Errorf("unexpected object type"), "PixiuFilterPolicy watch received unexpected object", "found", reflect.TypeOf(obj))
		return nil
	}
	group := policy.Spec.TargetRef.Group
	if group != "" && !strings.EqualFold(group, gatewayv1.GroupName) {
		return nil
	}
	if !strings.EqualFold(policy.Spec.TargetRef.Kind, KindGateway) {
		return nil
	}
	ns := policy.Namespace
	if policy.Spec.TargetRef.Namespace != nil {
		ns = string(*policy.Spec.TargetRef.Namespace)
	}
	gwName := string(policy.Spec.TargetRef.Name)
	gw := new(gatewayv1.Gateway)
	if err := r.Get(ctx, client.ObjectKey{Namespace: ns, Name: gwName}, gw); err != nil {
		return nil
	}
	if !r.checkGatewayClass(gw) {
		return nil
	}
	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Namespace: ns,
			Name:      gwName,
		},
	}}
}

func (r *GatewayReconciler) listGatewaysForClusterPolicy(ctx context.Context, obj client.Object) []reconcile.Request {
	policy, ok := obj.(*v1alpha1.PixiuClusterPolicy)
	if !ok {
		r.Log.Error(fmt.Errorf("unexpected object type"), "PixiuClusterPolicy watch received unexpected object", "found", reflect.TypeOf(obj))
		return nil
	}
	group := policy.Spec.TargetRef.Group
	if group != "" && !strings.EqualFold(group, gatewayv1.GroupName) {
		r.Log.Info(
			"filter policy targetRef group does not match Gateway API group",
			"policy", client.ObjectKeyFromObject(policy),
			"group", group,
		)
		return nil
	}
	if !strings.EqualFold(policy.Spec.TargetRef.Kind, KindGateway) {
		r.Log.Info(
			"filter policy targetRef kind is not Gateway",
			"policy", client.ObjectKeyFromObject(policy),
			"kind", policy.Spec.TargetRef.Kind,
		)
		return nil
	}
	ns := policy.Namespace
	if policy.Spec.TargetRef.Namespace != nil {
		ns = string(*policy.Spec.TargetRef.Namespace)
	}
	gwName := string(policy.Spec.TargetRef.Name)
	gw := new(gatewayv1.Gateway)
	if err := r.Get(ctx, client.ObjectKey{Namespace: ns, Name: gwName}, gw); err != nil {
		r.Log.Info(
			"referenced Gateway not found for filter policy",
			"policy", client.ObjectKeyFromObject(policy),
			"gateway", client.ObjectKey{Namespace: ns, Name: gwName},
			"error", err,
		)
		return nil
	}
	if !r.checkGatewayClass(gw) {
		r.Log.Info(
			"gateway class not supported by this controller",
			"policy", client.ObjectKeyFromObject(policy),
			"gateway", client.ObjectKeyFromObject(gw),
			"gatewayClass", gw.Spec.GatewayClassName,
		)
		return nil
	}
	return []reconcile.Request{{
		NamespacedName: types.NamespacedName{
			Namespace: ns,
			Name:      gwName,
		},
	}}
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

func (r *GatewayReconciler) ensureGatewayConfigMap(ctx context.Context, gateway *gatewayv1.Gateway) (string, error) {
	configMapName := fmt.Sprintf("%s-config", gateway.GetName())

	translator := translator.NewTranslator(r.Client, r.Log)
	conv := converter.NewConverter()

	xds, err := translator.TranslateGateway(ctx, gateway)
	if err != nil {
		return "", fmt.Errorf("failed to translate Gateway to IR: %w", err)
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
		return "", fmt.Errorf("failed to convert IR to Pixiu config: %w", err)
	}

	policyLoader := NewPolicyLoader(r.Client, r.Log)
	gatewayPolicy, err := policyLoader.LoadGatewayPolicy(ctx, gateway)
	if err != nil {
		r.Log.Error(err, "failed to load gateway policy", "gateway", gateway.Name)
	} else if gatewayPolicy != nil {
		converter.ApplyGatewayPolicy(pixiuConfig, gatewayPolicy)
	}

	clusterPolicies, err := policyLoader.LoadAllClusterPolicies(ctx, gateway.Namespace)
	if err != nil {
		r.Log.Error(err, "failed to load cluster policies")
	} else {
		clusterConfigMap := make(map[string]*v1alpha1.ClusterConfig)
		serviceConfigMap := make(map[string]*v1alpha1.ServiceClusterConfig)

		for i := range clusterPolicies {
			policy := &clusterPolicies[i]

			isGatewayPolicy := policy.Spec.TargetRef.Kind == "Gateway" &&
				string(policy.Spec.TargetRef.Name) == gateway.Name

			r.Log.Info("processing cluster policy", "policy", policy.Name, "isGatewayPolicy", isGatewayPolicy)

			if isGatewayPolicy {
				for j := range policy.Spec.ClusterRef {
					clusterConfig := &policy.Spec.ClusterRef[j]
					clusterConfigMap[clusterConfig.Name] = clusterConfig
					r.Log.Info("added cluster config", "clusterName", clusterConfig.Name, "endpointCount", len(clusterConfig.Endpoints))
				}
			} else {
				for j := range policy.Spec.ServiceRef {
					serviceConfig := &policy.Spec.ServiceRef[j]
					serviceConfigMap[serviceConfig.Name] = serviceConfig
					r.Log.Info("added service config", "serviceName", serviceConfig.Name, "endpointCount", len(serviceConfig.Endpoints))
				}
			}
		}

		r.Log.Info("cluster policy maps", "clusterConfigCount", len(clusterConfigMap), "serviceConfigCount", len(serviceConfigMap))

		for _, cluster := range pixiuConfig.StaticResources.Clusters {
			r.Log.Info("generated cluster", "name", cluster.Name, "endpointCount", len(cluster.Endpoints))
		}

		for _, cluster := range pixiuConfig.StaticResources.Clusters {
			originalEndpoints := cluster.Endpoints
			originalEndpointCount := len(originalEndpoints)

			r.Log.Info("cluster before policy", "cluster", cluster.Name, "endpointCount", originalEndpointCount)

			if clusterConfig, ok := clusterConfigMap[cluster.Name]; ok {
				r.Log.Info("applying cluster config (exact match)", "cluster", cluster.Name, "hasEndpoints", len(clusterConfig.Endpoints) > 0)
				if err := r.resolveClusterEndpoints(ctx, gateway.Namespace, cluster, clusterConfig); err != nil {
					r.Log.Error(err, "failed to resolve cluster endpoints", "cluster", cluster.Name)
				}
				converter.ApplyClusterConfig(cluster, clusterConfig)
			} else {
				// match without namespace prefix (e.g., "default-service" -> "service")
				parts := strings.SplitN(cluster.Name, "-", 2)
				if len(parts) == 2 {
					shortName := parts[1]
					if clusterConfig, ok := clusterConfigMap[shortName]; ok {
						r.Log.Info("applying cluster config (short name match)", "cluster", cluster.Name, "shortName", shortName, "hasEndpoints", len(clusterConfig.Endpoints) > 0)
						if err := r.resolveClusterEndpoints(ctx, gateway.Namespace, cluster, clusterConfig); err != nil {
							r.Log.Error(err, "failed to resolve cluster endpoints", "cluster", cluster.Name)
						}
						converter.ApplyClusterConfig(cluster, clusterConfig)
					}
				}
			}

			if serviceConfig, ok := serviceConfigMap[cluster.Name]; ok {
				r.Log.Info("applying service cluster config (exact match)", "cluster", cluster.Name, "hasEndpoints", len(serviceConfig.Endpoints) > 0)
				r.resolveServiceClusterEndpoints(ctx, gateway.Namespace, cluster, serviceConfig)
				converter.ApplyClusterPolicy(cluster, serviceConfig)
			} else {
				parts := strings.SplitN(cluster.Name, "-", 2)
				if len(parts) == 2 {
					shortName := parts[1]
					if serviceConfig, ok := serviceConfigMap[shortName]; ok {
						r.Log.Info("applying service cluster config (short name match)", "cluster", cluster.Name, "shortName", shortName, "hasEndpoints", len(serviceConfig.Endpoints) > 0)
						r.resolveServiceClusterEndpoints(ctx, gateway.Namespace, cluster, serviceConfig)
						converter.ApplyClusterPolicy(cluster, serviceConfig)
					}
				}
			}

			if len(cluster.Endpoints) == 0 && originalEndpointCount > 0 {
				r.Log.Info("policy resulted in no endpoints, restoring original endpoints",
					"cluster", cluster.Name, "originalCount", originalEndpointCount)
				cluster.Endpoints = originalEndpoints
			}

			r.Log.Info("cluster after policy", "cluster", cluster.Name, "endpointCount", len(cluster.Endpoints))
		}
	}

	var filterPolicyList v1alpha1.PixiuFilterPolicyList
	if err := r.Client.List(ctx, &filterPolicyList, client.InNamespace(gateway.Namespace)); err == nil {
		for _, policy := range filterPolicyList.Items {
			if policy.Spec.TargetRef.Kind == "Gateway" && string(policy.Spec.TargetRef.Name) == gateway.Name {
				if len(policy.Spec.ListenersRef) > 0 {
					if err := converter.ApplyListenersRefToConfig(ctx, r.Client, gateway.Namespace, pixiuConfig, &policy); err != nil {
						r.Log.Error(err, "failed to apply listenersRef", "policy", policy.Name)
					}
				} else {
					for i := range pixiuConfig.StaticResources.Listeners {
						listener := pixiuConfig.StaticResources.Listeners[i]
						if err := converter.ApplyFilterPolicyToListener(listener, &policy); err != nil {
							r.Log.Error(err, "failed to apply filter policy to listener", "listener", listener.Name, "policy", policy.Name)
						}
					}
				}
			}
		}
	} else {
		r.Log.Error(err, "failed to list filter policies")
	}

	normalizeHTTPMethods(pixiuConfig)
	configYAML, err := converter.ConvertToYAML(pixiuConfig)
	if err != nil {
		return "", fmt.Errorf("failed to convert config to YAML: %w", err)
	}
	configHash := hashString(configYAML)

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
				return "", fmt.Errorf("failed to create configmap: %w", err)
			}
			r.Log.Info("created gateway configmap", "gateway", gateway.GetName(), "configmap", configMapName)
		} else {
			return "", fmt.Errorf("failed to check configmap: %w", err)
		}
	} else {
		if existingConfigMap.Data["conf.yaml"] != configYAML {
			existingConfigMap.Data = configMap.Data
			existingConfigMap.Labels = configMap.Labels
			if err := r.Update(ctx, existingConfigMap); err != nil {
				return "", fmt.Errorf("failed to update configmap: %w", err)
			}
			r.Log.Info("updated gateway configmap", "gateway", gateway.GetName(), "configmap", configMapName)

			deploymentName := fmt.Sprintf("%s-%s", gateway.GetName(), string(gateway.GetUID())[:8])
			if err := r.triggerHotReload(ctx, gateway, deploymentName); err != nil {
				return "", fmt.Errorf("failed to trigger hot reload after configmap update: %w", err)
			}
			r.Log.Info("hot reload triggered successfully after configmap update", "gateway", gateway.GetName())
		}
	}

	return configHash, nil
}

func (r *GatewayReconciler) ensureDataPlane(ctx context.Context, gateway *gatewayv1.Gateway, configHash string) error {
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
					Annotations: map[string]string{
						"pixiu.apache.org/config-hash": configHash,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:            "pixiu",
							Image:           config.ControllerConfig.Gateway.Image,
							ImagePullPolicy: corev1.PullPolicy(config.ValidateImagePullPolicy(config.ControllerConfig.Gateway.ImagePullPolicy)),
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
								{
									Name:          "reload",
									ContainerPort: 18380,
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
		existingHash := existingDeployment.Spec.Template.Annotations["pixiu.apache.org/config-hash"]
		if existingHash != configHash {
			r.Log.Info("config hash changed, triggering hot reload",
				"gateway", gateway.GetName(),
				"deployment", deploymentName,
				"oldHash", existingHash,
				"newHash", configHash)

			if err := r.triggerHotReload(ctx, gateway, deploymentName); err != nil {
				r.Log.Error(err, "failed to trigger hot reload, will update deployment", "gateway", gateway.GetName())
				needsUpdate = true
			} else {
				// Update annotation only after successful hot reload
				existingDeployment.Spec.Template.Annotations["pixiu.apache.org/config-hash"] = configHash
				needsUpdate = true
			}
		}
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
			r.Log.Info("updated data plane deployment", "gateway", gateway.GetName(), "deployment", deploymentName, "configHash", configHash)
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

func hashString(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}

func normalizeHTTPMethods(cfg *converter.PixiuConfig) {
	for i := range cfg.StaticResources.Listeners {
		l := cfg.StaticResources.Listeners[i]
		for fi := range l.FilterChain.Filters {
			f := l.FilterChain.Filters[fi]
			if f.Name != "dgp.filter.httpconnectionmanager" {
				continue
			}
			switch c := f.Config.(type) {
			case converter.HTTPConnectionManagerConfig:
				for _, r := range c.RouteConfig.Routes {
					if r.Match.Methods == nil {
						r.Match.Methods = []string{}
					}
				}
				f.Config = c
			case map[string]any:
				b, _ := json.Marshal(c)
				var hcm converter.HTTPConnectionManagerConfig
				if err := json.Unmarshal(b, &hcm); err == nil {
					for _, r := range hcm.RouteConfig.Routes {
						if r.Match.Methods == nil {
							r.Match.Methods = []string{}
						}
					}
					f.Config = hcm
				}
			default:
				fmt.Println("unsupported http connection manager config type.")
			}
			l.FilterChain.Filters[fi] = f
		}
	}
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

func (r *GatewayReconciler) resolveClusterEndpoints(ctx context.Context, namespace string, cluster *converter.Cluster, clusterConfig *v1alpha1.ClusterConfig) error {
	if len(clusterConfig.Endpoints) == 0 {
		return nil
	}

	resolvedEndpoints := []*converter.Endpoint{}
	resolveFailed := false
	for _, epConfig := range clusterConfig.Endpoints {
		if !isIPAddress(epConfig.Address) {
			serviceName, serviceNamespace := parseServiceAddress(epConfig.Address, namespace)
			r.Log.Info("resolving service DNS to Pod IPs", "address", epConfig.Address, "serviceName", serviceName, "namespace", serviceNamespace)

			podEndpoints, err := r.resolveServiceEndpoints(ctx, serviceNamespace, serviceName, epConfig.Port)
			if err == nil && len(podEndpoints) > 0 {
				r.Log.Info("resolved service to Pod endpoints", "service", serviceName, "endpointCount", len(podEndpoints))
				for _, podEp := range podEndpoints {
					resolvedEndpoints = append(resolvedEndpoints, &converter.Endpoint{
						ID: len(resolvedEndpoints) + 1,
						SocketAddress: converter.SocketAddress{
							Address: podEp.Address,
							Port:    int(podEp.Port),
						},
					})
				}
				continue
			}

			r.Log.Info("failed to resolve to Pod IPs, trying Service ClusterIP", "service", serviceName, "error", err)
			clusterIP, svcPort, err := r.resolveServiceClusterIP(ctx, serviceNamespace, serviceName, epConfig.Port)
			if err == nil && clusterIP != "" {
				resolvedEndpoints = append(resolvedEndpoints, &converter.Endpoint{
					ID: func() int {
						if epConfig.ID != nil {
							return int(*epConfig.ID)
						}
						return len(resolvedEndpoints) + 1
					}(),
					SocketAddress: converter.SocketAddress{
						Address: clusterIP,
						Port:    int(svcPort),
					},
				})
				continue
			}

			resolveFailed = true
			r.Log.Info("failed to resolve service, skipping this endpoint", "address", epConfig.Address, "error", err)
			continue
		}

		resolvedEndpoints = append(resolvedEndpoints, &converter.Endpoint{
			ID: func() int {
				if epConfig.ID != nil {
					return int(*epConfig.ID)
				}
				return len(resolvedEndpoints) + 1
			}(),
			SocketAddress: converter.SocketAddress{
				Address: epConfig.Address,
				Port:    int(epConfig.Port),
			},
		})
	}

	if len(resolvedEndpoints) == 0 {
		if resolveFailed {
			return fmt.Errorf("failed to resolve configured endpoints from ClusterPolicy")
		}
		return fmt.Errorf("ClusterPolicy configured no usable endpoints")
	}

	cluster.Endpoints = resolvedEndpoints
	return nil
}

func (r *GatewayReconciler) resolveServiceEndpoints(ctx context.Context, namespace, serviceName string, port int32) ([]*ir.Endpoint, error) {
	var service corev1.Service
	if err := r.Client.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      serviceName,
	}, &service); err != nil {
		return nil, fmt.Errorf("failed to get service %s/%s: %w", namespace, serviceName, err)
	}

	var endpointSliceList discoveryv1.EndpointSliceList
	if err := r.Client.List(ctx, &endpointSliceList, client.MatchingLabels{
		discoveryv1.LabelServiceName: serviceName,
	}, client.InNamespace(namespace)); err != nil {
		return nil, fmt.Errorf("failed to list endpoint slices: %w", err)
	}

	endpoints := []*ir.Endpoint{}
	for _, endpointSlice := range endpointSliceList.Items {
		for _, endpoint := range endpointSlice.Endpoints {
			if endpoint.Conditions.Ready != nil && !*endpoint.Conditions.Ready {
				continue
			}
			for _, address := range endpoint.Addresses {
				targetPort := port
				if len(endpointSlice.Ports) > 0 {
					for _, endpointPort := range endpointSlice.Ports {
						if endpointPort.Port != nil && endpointPort.Name != nil {
							targetPort = int32(*endpointPort.Port)
							break // Found port, break to use it
						}
					}

					if targetPort == port {
						for _, endpointPort := range endpointSlice.Ports {
							if endpointPort.Port != nil {
								targetPort = int32(*endpointPort.Port)
								break // Found port, break to use it
							}
						}
					}
				}
				endpoints = append(endpoints, &ir.Endpoint{
					Address: address,
					Port:    targetPort,
				})
			}
		}
	}

	if len(endpoints) == 0 && service.Spec.ClusterIP != "" && service.Spec.ClusterIP != "None" {
		endpoints = append(endpoints, &ir.Endpoint{
			Address: service.Spec.ClusterIP,
			Port:    port,
		})
	}

	return endpoints, nil
}

func (r *GatewayReconciler) resolveServiceClusterIP(ctx context.Context, namespace, serviceName string, desiredPort int32) (string, int32, error) {
	var service corev1.Service
	if err := r.Client.Get(ctx, client.ObjectKey{
		Namespace: namespace,
		Name:      serviceName,
	}, &service); err != nil {
		return "", 0, fmt.Errorf("failed to get service %s/%s: %w", namespace, serviceName, err)
	}

	if service.Spec.ClusterIP == "" || service.Spec.ClusterIP == "None" {
		return "", 0, fmt.Errorf("service %s/%s is headless or has no ClusterIP", namespace, serviceName)
	}

	resolvedPort := desiredPort
	if resolvedPort == 0 && len(service.Spec.Ports) > 0 {
		resolvedPort = service.Spec.Ports[0].Port
	} else if resolvedPort != 0 {
		for _, p := range service.Spec.Ports {
			if p.Port == desiredPort {
				resolvedPort = p.Port
				break
			}
		}
	}

	if resolvedPort == 0 {
		return "", 0, fmt.Errorf("no valid port resolved for service %s/%s", namespace, serviceName)
	}

	return service.Spec.ClusterIP, resolvedPort, nil
}

func isIPAddress(s string) bool {
	parts := strings.Split(s, ".")
	if len(parts) != 4 {
		return strings.Contains(s, ":")
	}
	for _, part := range parts {
		if len(part) == 0 || len(part) > 3 {
			return false
		}
		for _, c := range part {
			if c < '0' || c > '9' {
				return false
			}
		}
	}
	return true
}

func parseServiceAddress(address, defaultNamespace string) (serviceName, namespace string) {
	if strings.HasSuffix(address, ".svc.cluster.local") {
		withoutSuffix := strings.TrimSuffix(address, ".svc.cluster.local")
		parts := strings.Split(withoutSuffix, ".")
		if len(parts) >= 2 {
			serviceName = parts[0]
			namespace = strings.Join(parts[1:], ".")
			return serviceName, namespace
		}
		serviceName = parts[0]
		namespace = defaultNamespace
		return serviceName, namespace
	}

	parts := strings.Split(address, ".")
	if len(parts) >= 2 {
		if !isNumeric(parts[1]) {
			serviceName = parts[0]
			namespace = strings.Join(parts[1:], ".")
			return serviceName, namespace
		}
	}

	return address, defaultNamespace
}

func isNumeric(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return len(s) > 0
}

// resolveServiceClusterEndpoints resolves Service DNS names in ServiceClusterConfig endpoints to actual IPs
func (r *GatewayReconciler) resolveServiceClusterEndpoints(ctx context.Context, namespace string, cluster *converter.Cluster, serviceConfig *v1alpha1.ServiceClusterConfig) {
	if len(serviceConfig.Endpoints) == 0 {
		return
	}

	resolved := make([]v1alpha1.EndpointConfig, 0, len(serviceConfig.Endpoints))
	idCounter := 1

	for i := range serviceConfig.Endpoints {
		ep := &serviceConfig.Endpoints[i]
		if isIPAddress(ep.Address) {
			// Keep IP addresses as-is
			resolved = append(resolved, *ep)
			idCounter++
			continue
		}

		serviceName, serviceNamespace := parseServiceAddress(ep.Address, namespace)
		r.Log.Info("resolving service DNS", "address", ep.Address, "serviceName", serviceName, "namespace", serviceNamespace)

		endpoints, err := r.resolveServiceEndpoints(ctx, serviceNamespace, serviceName, ep.Port)
		if err == nil && len(endpoints) > 0 {
			r.Log.Info("resolved service to endpoints", "service", serviceName, "endpointCount", len(endpoints))
			for _, resolvedEp := range endpoints {
				id := int32(idCounter)
				resolved = append(resolved, v1alpha1.EndpointConfig{
					ID:      &id,
					Address: resolvedEp.Address,
					Port:    resolvedEp.Port,
				})
				idCounter++
			}
		} else {
			r.Log.Info("failed to resolve service, keeping original address", "service", serviceName, "error", err)
			// Keep the original unresolved entry
			resolved = append(resolved, *ep)
			idCounter++
		}
	}

	serviceConfig.Endpoints = resolved
}

// triggerHotReload triggers hot reload on all pods in the deployment
func (r *GatewayReconciler) triggerHotReload(ctx context.Context, gateway *gatewayv1.Gateway, _ string) error {
	configMapName := fmt.Sprintf("%s-config", gateway.GetName())
	configMap := &corev1.ConfigMap{}
	if err := r.Get(ctx, client.ObjectKey{
		Name:      configMapName,
		Namespace: gateway.GetNamespace(),
	}, configMap); err != nil {
		return fmt.Errorf("failed to get configmap: %w", err)
	}

	configYAML, ok := configMap.Data["conf.yaml"]
	if !ok {
		return fmt.Errorf("conf.yaml not found in configmap")
	}

	httpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	podList := &corev1.PodList{}
	if err := r.List(ctx, podList,
		client.InNamespace(gateway.GetNamespace()),
		client.MatchingLabels{
			"app.kubernetes.io/name":                 "pixiu-gateway",
			"gateway.networking.k8s.io/gateway-name": gateway.GetName(),
		}); err != nil {
		return fmt.Errorf("failed to list pods: %w", err)
	}

	if len(podList.Items) == 0 {
		r.Log.Info("no pods found for hot reload", "gateway", gateway.GetName())
		return nil
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	var successCount int
	var failureCount int

	for _, pod := range podList.Items {
		pod := pod
		wg.Add(1)
		go func() {
			defer wg.Done()

			if pod.Status.Phase != corev1.PodRunning {
				r.Log.Info("skipping non-running pod", "pod", pod.Name, "phase", pod.Status.Phase)
				return
			}

			podIP := pod.Status.PodIP
			if podIP == "" {
				r.Log.Info("pod has no IP address", "pod", pod.Name)
				return
			}

			reloadURL := fmt.Sprintf("http://%s:18380/-/reload", podIP)
			r.Log.Info("triggering hot reload", "pod", pod.Name, "url", reloadURL)

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, reloadURL, strings.NewReader(configYAML))
			if err != nil {
				r.Log.Error(err, "failed to create reload request", "pod", pod.Name)
				mu.Lock()
				failureCount++
				mu.Unlock()
				return
			}
			req.Header.Set("Content-Type", "application/x-yaml")

			resp, err := httpClient.Do(req)
			if err != nil {
				r.Log.Error(err, "failed to trigger hot reload", "pod", pod.Name, "url", reloadURL)
				mu.Lock()
				failureCount++
				mu.Unlock()
				return
			}

			_, _ = io.Copy(io.Discard, resp.Body)
			if closeErr := resp.Body.Close(); closeErr != nil {
				r.Log.Error(closeErr, "failed to close reload response body", "pod", pod.Name)
			}

			mu.Lock()
			defer mu.Unlock()
			if resp.StatusCode == http.StatusOK {
				r.Log.Info("hot reload successful", "pod", pod.Name)
				successCount++
				return
			}

			failureCount++
			r.Log.Error(fmt.Errorf("unexpected status code: %d", resp.StatusCode),
				"hot reload failed", "pod", pod.Name)
		}()
	}

	wg.Wait()

	if successCount == 0 {
		return fmt.Errorf("hot reload failed on all pods")
	}

	r.Log.Info("hot reload completed", "gateway", gateway.GetName(),
		"successCount", successCount, "failureCount", failureCount, "totalPods", len(podList.Items))
	return nil
}
