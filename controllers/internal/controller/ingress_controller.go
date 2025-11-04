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
	"fmt"
	"reflect"
)

import (
	"controllers/api/v1alpha1"

	"controllers/internal/controller/indexer"
	"controllers/internal/controller/status"

	"controllers/internal/utils"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"

	discoveryv1 "k8s.io/api/discovery/v1"

	networkingv1 "k8s.io/api/networking/v1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

const defaultIngressClassAnnotation = "ingressclass.kubernetes.io/is-default-class"

// IngressReconciler reconciles a Ingress object.
type IngressReconciler struct { //nolint:revive
	client.Client
	Scheme       *runtime.Scheme
	Log          logr.Logger
	genericEvent chan event.GenericEvent
	Updater      status.Updater
}

// SetupWithManager sets up the controller with the Manager.
func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.genericEvent = make(chan event.GenericEvent, 100)

	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.Ingress{},
			builder.WithPredicates(
				predicate.NewPredicateFuncs(r.checkIngressClass),
			),
		).
		WithEventFilter(
			predicate.Or(
				predicate.GenerationChangedPredicate{},
				predicate.AnnotationChangedPredicate{},
				predicate.NewPredicateFuncs(TypePredicate[*corev1.Secret]()),
			),
		).
		Watches(
			&networkingv1.IngressClass{},
			handler.EnqueueRequestsFromMapFunc(r.listIngressForIngressClass),
			builder.WithPredicates(
				predicate.NewPredicateFuncs(r.matchesIngressController),
			),
		).
		Watches(
			&discoveryv1.EndpointSlice{},
			handler.EnqueueRequestsFromMapFunc(r.listIngressesByService),
		).
		Watches(
			&corev1.Secret{},
			handler.EnqueueRequestsFromMapFunc(r.listIngressesBySecret),
		).
		Watches(&v1alpha1.GatewayProxy{},
			handler.EnqueueRequestsFromMapFunc(r.listIngressesForGatewayProxy),
		).
		WatchesRawSource(
			source.Channel(
				r.genericEvent,
				handler.EnqueueRequestsFromMapFunc(r.listIngressForGenericEvent),
			),
		).
		Complete(r)
}

// Reconcile handles the reconciliation of Ingress resources
func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	ingress := new(networkingv1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, ingress); err != nil {
		if client.IgnoreNotFound(err) == nil {
			ingress.Namespace = req.Namespace
			ingress.Name = req.Name

			ingress.TypeMeta = metav1.TypeMeta{
				Kind:       KindIngress,
				APIVersion: networkingv1.SchemeGroupVersion.String(),
			}
			r.Log.Info("deleted ingress resources", "ingress", ingress.Name)
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	r.Log.Info("reconciling ingress", "ingress", ingress.Name)

	ingressClass, err := r.getIngressClass(ctx, ingress)
	if err != nil {
		r.Log.Error(err, "failed to get IngressClass")
		return ctrl.Result{}, err
	}

	// process IngressClass parameters if they reference GatewayProxy
	if err := ProcessIngressClassParameters(r.Client, r.Log, ingress, ingressClass); err != nil {
		r.Log.Error(err, "failed to process IngressClass parameters", "ingressClass", ingressClass.Name)
		return ctrl.Result{}, err
	}

	// process backend services
	if err := r.processBackends(ingress); err != nil {
		r.Log.Error(err, "failed to process backend services", "ingress", ingress.Name)
		return ctrl.Result{}, err
	}

	// update the ingress status
	if err := r.updateStatus(ctx, ingress, ingressClass); err != nil {
		r.Log.Error(err, "failed to update ingress status", "ingress", ingress.Name)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// getIngressClass get the ingress class for the ingress
func (r *IngressReconciler) getIngressClass(ctx context.Context, obj client.Object) (*networkingv1.IngressClass, error) {
	ingress := obj.(*networkingv1.Ingress)
	var ingressClassName string
	if ingress.Spec.IngressClassName != nil {
		ingressClassName = *ingress.Spec.IngressClassName
	}
	return GetIngressClass(ctx, r.Client, r.Log, ingressClassName)
}

// checkIngressClass check if the ingress uses the ingress class that we control
func (r *IngressReconciler) checkIngressClass(obj client.Object) bool {
	_, err := r.getIngressClass(context.Background(), obj)
	return err == nil
}

// matchesIngressController check if the ingress class is controlled by us
func (r *IngressReconciler) matchesIngressController(obj client.Object) bool {
	ingressClass, ok := obj.(*networkingv1.IngressClass)
	if !ok {
		r.Log.Error(fmt.Errorf("unexpected object type"), "failed to convert object to IngressClass")
		return false
	}

	return matchesController(ingressClass.Spec.Controller)
}

// listIngressForIngressClass list all ingresses that use a specific ingress class
func (r *IngressReconciler) listIngressForIngressClass(ctx context.Context, obj client.Object) []reconcile.Request {
	ingressClass, ok := obj.(*networkingv1.IngressClass)
	if !ok {
		r.Log.Error(fmt.Errorf("unexpected object type"), "failed to convert object to IngressClass")
		return nil
	}

	// check if the ingress class is the default ingress class
	if IsDefaultIngressClass(ingressClass) {
		ingressList := &networkingv1.IngressList{}
		if err := r.List(ctx, ingressList); err != nil {
			r.Log.Error(err, "failed to list ingresses for ingress class", "ingressclass", ingressClass.GetName())
			return nil
		}

		requests := make([]reconcile.Request, 0, len(ingressList.Items))
		for _, ingress := range ingressList.Items {
			if ingress.Spec.IngressClassName == nil || *ingress.Spec.IngressClassName == "" ||
				*ingress.Spec.IngressClassName == ingressClass.GetName() {
				requests = append(requests, reconcile.Request{
					NamespacedName: client.ObjectKey{
						Namespace: ingress.Namespace,
						Name:      ingress.Name,
					},
				})
			}
		}
		return requests
	} else {
		ingressList := &networkingv1.IngressList{}
		if err := r.List(ctx, ingressList, client.MatchingFields{
			indexer.IngressClassRef: ingressClass.GetName(),
		}); err != nil {
			r.Log.Error(err, "failed to list ingresses for ingress class", "ingressclass", ingressClass.GetName())
			return nil
		}

		requests := make([]reconcile.Request, 0, len(ingressList.Items))
		for _, ingress := range ingressList.Items {
			requests = append(requests, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Namespace: ingress.Namespace,
					Name:      ingress.Name,
				},
			})
		}

		return requests
	}
}

func IsDefaultIngressClass(obj client.Object) bool {
	if ingressClass, ok := obj.(*networkingv1.IngressClass); ok {
		return ingressClass.Annotations[defaultIngressClassAnnotation] == "true"
	}
	return false
}

// listIngressesByService list all ingresses that use a specific service
func (r *IngressReconciler) listIngressesByService(ctx context.Context, obj client.Object) []reconcile.Request {
	endpointSlice, ok := obj.(*discoveryv1.EndpointSlice)
	if !ok {
		r.Log.Error(fmt.Errorf("unexpected object type"), "failed to convert object to EndpointSlice")
		return nil
	}

	namespace := endpointSlice.GetNamespace()
	serviceName := endpointSlice.Labels[discoveryv1.LabelServiceName]

	ingressList := &networkingv1.IngressList{}
	if err := r.List(ctx, ingressList, client.MatchingFields{
		indexer.ServiceIndexRef: indexer.GenIndexKey(namespace, serviceName),
	}); err != nil {
		r.Log.Error(err, "failed to list ingresses by service", "service", serviceName)
		return nil
	}

	requests := make([]reconcile.Request, 0, len(ingressList.Items))
	for _, ingress := range ingressList.Items {
		if r.checkIngressClass(&ingress) {
			requests = append(requests, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Namespace: ingress.Namespace,
					Name:      ingress.Name,
				},
			})
		}
	}
	return requests
}

// listIngressesBySecret list all ingresses that use a specific secret
func (r *IngressReconciler) listIngressesBySecret(ctx context.Context, obj client.Object) []reconcile.Request {
	secret, ok := obj.(*corev1.Secret)
	if !ok {
		r.Log.Error(fmt.Errorf("unexpected object type"), "failed to convert object to Secret")
		return nil
	}

	namespace := secret.GetNamespace()
	name := secret.GetName()

	ingressList := &networkingv1.IngressList{}
	if err := r.List(ctx, ingressList, client.MatchingFields{
		indexer.SecretIndexRef: indexer.GenIndexKey(namespace, name),
	}); err != nil {
		r.Log.Error(err, "failed to list ingresses by secret", "secret", name)
		return nil
	}

	requests := make([]reconcile.Request, 0, len(ingressList.Items))
	for _, ingress := range ingressList.Items {
		if r.checkIngressClass(&ingress) {
			requests = append(requests, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Namespace: ingress.Namespace,
					Name:      ingress.Name,
				},
			})
		}
	}

	gatewayProxyList := &v1alpha1.GatewayProxyList{}
	if err := r.List(ctx, gatewayProxyList, client.MatchingFields{
		indexer.SecretIndexRef: indexer.GenIndexKey(secret.GetNamespace(), secret.GetName()),
	}); err != nil {
		r.Log.Error(err, "failed to list gateway proxies by secret", "secret", secret.GetName())
		return nil
	}

	for _, gatewayProxy := range gatewayProxyList.Items {
		var (
			ingressClassList networkingv1.IngressClassList
			indexKey         = indexer.GenIndexKey(gatewayProxy.GetNamespace(), gatewayProxy.GetName())
			matchingFields   = client.MatchingFields{indexer.IngressClassParametersRef: indexKey}
		)
		if err := r.List(ctx, &ingressClassList, matchingFields); err != nil {
			r.Log.Error(err, "failed to list ingress classes for gateway proxy", "gatewayproxy", indexKey)
			continue
		}
		for _, ingressClass := range ingressClassList.Items {
			requests = append(requests, r.listIngressForIngressClass(ctx, &ingressClass)...)
		}
	}

	return distinctRequests(requests)
}

func (r *IngressReconciler) listIngressForGenericEvent(ctx context.Context, obj client.Object) (requests []reconcile.Request) {
	switch obj.(type) {
	default:
		r.Log.Error(fmt.Errorf("unexpected object type"), "failed to convert object to BackendTrafficPolicy")
		return nil
	}
}

// processTLS process the TLS configuration of the ingress
func (r *IngressReconciler) processTLS(ingress *networkingv1.Ingress) error {
	for _, tls := range ingress.Spec.TLS {
		if tls.SecretName == "" {
			continue
		}

		secret := corev1.Secret{}
		if err := r.Get(context.TODO(), client.ObjectKey{
			Namespace: ingress.Namespace,
			Name:      tls.SecretName,
		}, &secret); err != nil {
			r.Log.Error(err, "failed to get secret", "namespace", ingress.Namespace, "name", tls.SecretName)
			return err
		}

		if secret.Data == nil {
			continue
		}
	}

	return nil
}

// processBackends process the backend services of the ingress
func (r *IngressReconciler) processBackends(ingress *networkingv1.Ingress) error {
	var terr error

	// process all the backend services in the rules
	for _, rule := range ingress.Spec.Rules {
		if rule.HTTP == nil {
			continue
		}
		for _, path := range rule.HTTP.Paths {
			if path.Backend.Service == nil {
				continue
			}
			service := path.Backend.Service
			if err := r.processBackendService(ingress.Namespace, service); err != nil {
				terr = err
			}
		}
	}
	return terr
}

// processBackendService process a single backend service
func (r *IngressReconciler) processBackendService(namespace string, backendService *networkingv1.IngressServiceBackend) error {
	// get the service
	var service corev1.Service
	serviceNS := types.NamespacedName{
		Namespace: namespace,
		Name:      backendService.Name,
	}
	if err := r.Get(context.TODO(), serviceNS, &service); err != nil {
		if client.IgnoreNotFound(err) == nil {
			r.Log.Info("service not found", "namespace", namespace, "name", backendService.Name)
			return nil
		}
		return err
	}

	if service.Spec.Type == corev1.ServiceTypeExternalName {
		return nil
	}

	// verify if the port exists
	var portExists bool
	if backendService.Port.Number != 0 {
		for _, servicePort := range service.Spec.Ports {
			if servicePort.Port == backendService.Port.Number {
				portExists = true
				break
			}
		}
	} else if backendService.Port.Name != "" {
		for _, servicePort := range service.Spec.Ports {
			if servicePort.Name == backendService.Port.Name {
				portExists = true
				break
			}
		}
	}

	if !portExists {
		err := fmt.Errorf("port(name: %s, number: %d) not found in service %s/%s", backendService.Port.Name, backendService.Port.Number, namespace, backendService.Name)
		r.Log.Error(err, "service port not found")
		return err
	}

	// get the endpoint slices
	endpointSliceList := &discoveryv1.EndpointSliceList{}
	if err := r.List(context.TODO(), endpointSliceList,
		client.InNamespace(namespace),
		client.MatchingLabels{
			discoveryv1.LabelServiceName: backendService.Name,
		},
	); err != nil {
		r.Log.Error(err, "failed to list endpoint slices", "namespace", namespace, "name", backendService.Name)
		return err
	}

	return nil
}

// updateStatus update the status of the ingress
func (r *IngressReconciler) updateStatus(ctx context.Context, ingress *networkingv1.Ingress, ingressClass *networkingv1.IngressClass) error {
	var loadBalancerStatus networkingv1.IngressLoadBalancerStatus

	// update the load balancer status
	if len(loadBalancerStatus.Ingress) > 0 && !reflect.DeepEqual(ingress.Status.LoadBalancer, loadBalancerStatus) {
		ingress.Status.LoadBalancer = loadBalancerStatus
		r.Updater.Update(status.Update{
			NamespacedName: utils.NamespacedName(ingress),
			Resource:       ingress.DeepCopy(),
			Mutator: status.MutatorFunc(func(obj client.Object) client.Object {
				cp := obj.(*networkingv1.Ingress).DeepCopy()
				cp.Status = ingress.Status
				return cp
			}),
		})
		return nil
	}

	return nil
}

// listIngressesForGatewayProxy list all ingresses that use a specific gateway proxy
func (r *IngressReconciler) listIngressesForGatewayProxy(ctx context.Context, obj client.Object) []reconcile.Request {
	return listIngressClassRequestsForGatewayProxy(ctx, r.Client, obj, r.Log, r.listIngressForIngressClass)
}

func listIngressClassRequestsForGatewayProxy(
	ctx context.Context,
	c client.Client,
	obj client.Object,
	logger logr.Logger,
	listFunc func(context.Context, client.Object) []reconcile.Request,
) []reconcile.Request {
	gatewayProxy, ok := obj.(*v1alpha1.GatewayProxy)
	if !ok {
		return nil
	}

	ingressClassList := &networkingv1.IngressClassList{}
	if err := c.List(ctx, ingressClassList, client.MatchingFields{
		indexer.IngressClassParametersRef: indexer.GenIndexKey(gatewayProxy.GetNamespace(), gatewayProxy.GetName()),
	}); err != nil {
		logger.Error(err, "failed to list ingress classes for gateway proxy", "gatewayproxy", gatewayProxy.GetName())
		return nil
	}

	requestSet := make(map[string]reconcile.Request)
	for _, ingressClass := range ingressClassList.Items {
		for _, req := range listFunc(ctx, &ingressClass) {
			requestSet[req.String()] = req
		}
	}

	requests := make([]reconcile.Request, 0, len(requestSet))
	for _, req := range requestSet {
		requests = append(requests, req)
	}
	return requests
}

func distinctRequests(requests []reconcile.Request) []reconcile.Request {
	uniqueRequests := make(map[string]reconcile.Request)
	for _, request := range requests {
		uniqueRequests[request.String()] = request
	}

	distinctRequests := make([]reconcile.Request, 0, len(uniqueRequests))
	for _, request := range uniqueRequests {
		distinctRequests = append(distinctRequests, request)
	}
	return distinctRequests
}
