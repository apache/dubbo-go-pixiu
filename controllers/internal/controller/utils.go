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
	"path"
	"reflect"
	"strings"
)

import (
	"controllers/internal/utils"

	"github.com/go-logr/logr"

	"github.com/samber/lo"

	corev1 "k8s.io/api/core/v1"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
)

const (
	KindGateway = "Gateway"
	KindSecret  = "Secret"
)

var (
	enableReferenceGrant bool
)

func SetEnableReferenceGrant(enable bool) {
	enableReferenceGrant = enable
}

func GetEnableReferenceGrant() bool {
	return enableReferenceGrant
}

func TypePredicate[T client.Object]() func(obj client.Object) bool {
	return func(obj client.Object) bool {
		_, ok := obj.(T)
		return ok
	}
}

func IsConditionPresentAndEqual(conditions []metav1.Condition, condition metav1.Condition) bool {
	for _, cond := range conditions {
		if cond.Type == condition.Type &&
			cond.Reason == condition.Reason &&
			cond.Status == condition.Status &&
			cond.ObservedGeneration == condition.ObservedGeneration {
			return true
		}
	}
	return false
}

func ListRequests(ctx context.Context, c client.Client, logger logr.Logger, listObj client.ObjectList, opts ...client.ListOption) []reconcile.Request {
	return ListMatchingRequests(ctx, c, logger, listObj, func(obj client.Object) bool { return true }, opts...)
}

func ListMatchingRequests(ctx context.Context, c client.Client, logger logr.Logger, listObj client.ObjectList, matchFunc func(obj client.Object) bool, opts ...client.ListOption) []reconcile.Request {
	if err := c.List(ctx, listObj, opts...); err != nil {
		logger.Error(err, "failed to list resource")
		return nil
	}

	items, err := meta.ExtractList(listObj)
	if err != nil {
		logger.Error(err, "failed to extract list items")
		return nil
	}

	var requests []reconcile.Request
	for _, item := range items {
		obj, ok := item.(client.Object)
		if !ok {
			continue
		}

		if matchFunc(obj) {
			requests = append(requests, reconcile.Request{
				NamespacedName: utils.NamespacedName(obj),
			})
		}
	}
	return requests
}

func getListenerStatus(ctx context.Context, mrgc client.Client, gateway *gatewayv1.Gateway) ([]gatewayv1.ListenerStatus, error) {
	statuses := make(map[gatewayv1.SectionName]gatewayv1.ListenerStatus, len(gateway.Spec.Listeners))

	for i, listener := range gateway.Spec.Listeners {
		attachedRoutes, err := getAttachedRoutesForListener(ctx, mrgc, *gateway, listener)
		if err != nil {
			return nil, err
		}
		var (
			now                 = metav1.Now()
			conditionProgrammed = metav1.Condition{
				Type:               string(gatewayv1.ListenerConditionProgrammed),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: gateway.GetGeneration(),
				LastTransitionTime: now,
				Reason:             string(gatewayv1.ListenerReasonProgrammed),
			}
			conditionAccepted = metav1.Condition{
				Type:               string(gatewayv1.ListenerConditionAccepted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: gateway.GetGeneration(),
				LastTransitionTime: now,
				Reason:             string(gatewayv1.ListenerReasonAccepted),
			}
			conditionConflicted = metav1.Condition{
				Type:               string(gatewayv1.ListenerConditionConflicted),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: gateway.GetGeneration(),
				LastTransitionTime: now,
				Reason:             string(gatewayv1.ListenerReasonNoConflicts),
			}
			conditionResolvedRefs = metav1.Condition{
				Type:               string(gatewayv1.ListenerConditionResolvedRefs),
				Status:             metav1.ConditionTrue,
				ObservedGeneration: gateway.GetGeneration(),
				LastTransitionTime: now,
				Reason:             string(gatewayv1.ListenerReasonResolvedRefs),
			}

			supportedKinds = []gatewayv1.RouteGroupKind{}
		)

		status := gatewayv1.ListenerStatus{
			Name: listener.Name,
			Conditions: []metav1.Condition{
				conditionProgrammed,
				conditionAccepted,
				conditionConflicted,
				conditionResolvedRefs,
			},
			SupportedKinds: supportedKinds,
			AttachedRoutes: attachedRoutes,
		}

		changed := false
		if len(gateway.Status.Listeners) > i {
			if gateway.Status.Listeners[i].AttachedRoutes != attachedRoutes {
				changed = true
			}
			for _, condition := range status.Conditions {
				if !IsConditionPresentAndEqual(gateway.Status.Listeners[i].Conditions, condition) {
					changed = true
					break
				}
			}
		} else {
			changed = true
		}

		if changed {
			statuses[listener.Name] = status
		} else {
			statuses[listener.Name] = gateway.Status.Listeners[i]
		}
	}

	// check for conflicts

	statusArray := []gatewayv1.ListenerStatus{}
	for _, status := range statuses {
		statusArray = append(statusArray, status)
	}

	return statusArray, nil
}

func getAttachedRoutesForListener(ctx context.Context, mgrc client.Client, gateway gatewayv1.Gateway, listener gatewayv1.Listener) (int32, error) {
	httpRouteList := gatewayv1.HTTPRouteList{}
	if err := mgrc.List(ctx, &httpRouteList); err != nil {
		return 0, err
	}
	var attachedRoutes int32
	for _, route := range httpRouteList.Items {
		route := route
		acceptedByGateway := lo.ContainsBy(route.Status.Parents, func(parentStatus gatewayv1.RouteParentStatus) bool {
			parentRef := parentStatus.ParentRef
			if parentRef.Group != nil && *parentRef.Group != gatewayv1.GroupName {
				return false
			}
			if parentRef.Kind != nil && *parentRef.Kind != KindGateway {
				return false
			}
			gatewayNamespace := route.Namespace
			if parentRef.Namespace != nil {
				gatewayNamespace = string(*parentRef.Namespace)
			}
			return gateway.Namespace == gatewayNamespace && gateway.Name == string(parentRef.Name)
		})
		if !acceptedByGateway {
			continue
		}

		for _, parentRef := range route.Spec.ParentRefs {
			ok, _, err := checkRouteAcceptedByListener(
				ctx,
				mgrc,
				&route,
				gateway,
				listener,
				parentRef,
			)
			if err != nil {
				return 0, err
			}
			if ok {
				attachedRoutes++
			}
		}
	}
	return attachedRoutes, nil
}

func SetGatewayConditionAccepted(gw *gatewayv1.Gateway, status bool, message string) (ok bool) {
	condition := metav1.Condition{
		Type:               string(gatewayv1.GatewayConditionAccepted),
		Status:             ConditionStatus(status),
		Reason:             string(gatewayv1.GatewayReasonAccepted),
		ObservedGeneration: gw.GetGeneration(),
		Message:            message,
		LastTransitionTime: metav1.Now(),
	}

	if !IsConditionPresentAndEqual(gw.Status.Conditions, condition) {
		setGatewayCondition(gw, condition)
		ok = true
	}
	return
}

func ConditionStatus(status bool) metav1.ConditionStatus {
	if status {
		return metav1.ConditionTrue
	}
	return metav1.ConditionFalse
}

func setGatewayCondition(gw *gatewayv1.Gateway, newCondition metav1.Condition) {
	gw.Status.Conditions = MergeCondition(gw.Status.Conditions, newCondition)
}

func acceptedMessage(kind string) string {
	return fmt.Sprintf("the %s has been accepted by the pixiu-gateway-controller", kind)
}

func referenceGrantPredicates(kind gatewayv1.Kind) predicate.Funcs {
	var filter = func(obj client.Object) bool {
		grant, ok := obj.(*v1beta1.ReferenceGrant)
		if !ok {
			return false
		}
		for _, from := range grant.Spec.From {
			if from.Kind == kind && string(from.Group) == gatewayv1.GroupName {
				return true
			}
		}
		return false
	}
	predicates := predicate.NewPredicateFuncs(filter)
	predicates.UpdateFunc = func(e event.UpdateEvent) bool {
		return filter(e.ObjectOld) || filter(e.ObjectNew)
	}
	return predicates
}

func checkRouteAcceptedByListener(ctx context.Context, mgrc client.Client, route client.Object, gateway gatewayv1.Gateway, listener gatewayv1.Listener, parentRef gatewayv1.ParentReference) (bool, gatewayv1.RouteConditionReason, error) {
	if parentRef.SectionName != nil {
		if *parentRef.SectionName != "" && *parentRef.SectionName != listener.Name {
			return false, gatewayv1.RouteReasonNoMatchingParent, nil
		}
	}
	if parentRef.Port != nil {
		if *parentRef.Port != listener.Port {
			return false, gatewayv1.RouteReasonNoMatchingParent, nil
		}
	}
	if !routeMatchesListenerType(route, listener) {
		return false, gatewayv1.RouteReasonNoMatchingParent, nil
	}
	if !routeHostnamesIntersectsWithListenerHostname(route, listener) {
		return false, gatewayv1.RouteReasonNoMatchingListenerHostname, nil
	}
	if ok, err := routeMatchesListenerAllowedRoutes(ctx, mgrc, route, listener.AllowedRoutes, gateway.Namespace, parentRef.Namespace); err != nil {
		return false, gatewayv1.RouteReasonNotAllowedByListeners, fmt.Errorf("failed matching listener %s to a route %s for gateway %s: %w",
			listener.Name, route.GetName(), gateway.Name, err,
		)
	} else if !ok {
		return false, gatewayv1.RouteReasonNotAllowedByListeners, nil
	}
	return true, gatewayv1.RouteReasonAccepted, nil
}

func routeMatchesListenerAllowedRoutes(ctx context.Context, mgrc client.Client, route client.Object, allowedRoutes *gatewayv1.AllowedRoutes, gatewayNamespace string, parentRefNamespace *gatewayv1.Namespace) (bool, error) {
	if allowedRoutes == nil {
		return true, nil
	}

	if !isRouteKindAllowed(route, allowedRoutes.Kinds) {
		return false, fmt.Errorf("route %s/%s is not allowed in the kind", route.GetNamespace(), route.GetName())
	}

	if !isRouteNamespaceAllowed(ctx, route, mgrc, gatewayNamespace, parentRefNamespace, allowedRoutes.Namespaces) {
		return false, fmt.Errorf("route %s/%s is not allowed in the namespace", route.GetNamespace(), route.GetName())
	}

	return true, nil
}

func isRouteKindAllowed(route client.Object, kinds []gatewayv1.RouteGroupKind) (ok bool) {
	ok = true
	if len(kinds) > 0 {
		_, ok = lo.Find(kinds, func(rgk gatewayv1.RouteGroupKind) bool {
			gvk := route.GetObjectKind().GroupVersionKind()
			return (rgk.Group != nil && string(*rgk.Group) == gvk.Group) && string(rgk.Kind) == gvk.Kind
		})
	}
	return
}

func isRouteNamespaceAllowed(ctx context.Context, route client.Object, mgrc client.Client, gatewayNamespace string, parentRefNamespace *gatewayv1.Namespace, routeNamespaces *gatewayv1.RouteNamespaces) bool {
	if routeNamespaces == nil || routeNamespaces.From == nil {
		return true
	}

	switch *routeNamespaces.From {
	case gatewayv1.NamespacesFromAll:
		return true

	case gatewayv1.NamespacesFromSame:
		if parentRefNamespace == nil {
			return gatewayNamespace == route.GetNamespace()
		}
		return route.GetNamespace() == string(*parentRefNamespace)

	case gatewayv1.NamespacesFromSelector:
		namespace := corev1.Namespace{}
		if err := mgrc.Get(ctx, client.ObjectKey{Name: route.GetNamespace()}, &namespace); err != nil {
			return false
		}

		s, err := metav1.LabelSelectorAsSelector(routeNamespaces.Selector)
		if err != nil {
			return false
		}
		return s.Matches(labels.Set(namespace.Labels))
	default:
		return true
	}
}

func routeMatchesListenerType(route client.Object, listener gatewayv1.Listener) bool {
	switch route.(type) {
	case *gatewayv1.HTTPRoute:
		if listener.Protocol != gatewayv1.HTTPProtocolType && listener.Protocol != gatewayv1.HTTPSProtocolType {
			return false
		}

		if listener.Protocol == gatewayv1.HTTPSProtocolType {
			if listener.TLS == nil {
				return false
			}

			if listener.TLS.Mode != nil && *listener.TLS.Mode != gatewayv1.TLSModeTerminate {
				return false
			}
		}
	default:
		return false
	}
	return true
}

func routeHostnamesIntersectsWithListenerHostname(route client.Object, listener gatewayv1.Listener) bool {
	switch r := route.(type) {
	case *gatewayv1.HTTPRoute:
		return listenerHostnameIntersectWithRouteHostnames(listener, r.Spec.Hostnames)
	default:
		return false
	}
}

func reconcileGatewaysMatchGatewayClass(gatewayClass client.Object, gateways []gatewayv1.Gateway) (recs []reconcile.Request) {
	for _, gateway := range gateways {
		if string(gateway.Spec.GatewayClassName) == gatewayClass.GetName() {
			recs = append(recs, reconcile.Request{
				NamespacedName: client.ObjectKey{
					Name:      gateway.GetName(),
					Namespace: gateway.GetNamespace(),
				},
			})
		}
	}
	return
}

func listenerHostnameIntersectWithRouteHostnames(listener gatewayv1.Listener, hostnames []gatewayv1.Hostname) bool {
	if len(hostnames) == 0 {
		return true
	}

	if listener.Hostname == nil || *listener.Hostname == "" {
		return true
	}

	for _, hostname := range hostnames {
		if HostnamesIntersect(string(*listener.Hostname), string(hostname)) {
			return true
		}
	}

	return false
}

func FullTypeName(a any) string {
	typeOf := reflect.TypeOf(a)
	pkgPath := typeOf.PkgPath()
	name := typeOf.String()
	if typeOf.Kind() == reflect.Ptr {
		pkgPath = typeOf.Elem().PkgPath()
	}
	return path.Join(path.Dir(pkgPath), name)
}

func SetGatewayConditionProgrammed(gw *gatewayv1.Gateway, status bool, message string) (ok bool) {
	condition := metav1.Condition{
		Type:               string(gatewayv1.GatewayConditionProgrammed),
		Status:             ConditionStatus(status),
		Reason:             string(gatewayv1.GatewayReasonProgrammed),
		ObservedGeneration: gw.GetGeneration(),
		Message:            message,
		LastTransitionTime: metav1.Now(),
	}

	if !IsConditionPresentAndEqual(gw.Status.Conditions, condition) {
		setGatewayCondition(gw, condition)
		ok = true
	}
	return
}

func SetGatewayListenerConditionProgrammed(gw *gatewayv1.Gateway, listenerName string, status bool, message string) (ok bool) {
	condition := metav1.Condition{
		Type:               string(gatewayv1.ListenerConditionProgrammed),
		Status:             ConditionStatus(status),
		Reason:             string(gatewayv1.ListenerReasonProgrammed),
		ObservedGeneration: gw.GetGeneration(),
		Message:            message,
		LastTransitionTime: metav1.Now(),
	}

	if !IsConditionPresentAndEqual(gw.Status.Conditions, condition) {
		setListenerCondition(gw, listenerName, condition)
		ok = true
	}
	return
}

func SetGatewayListenerConditionResolvedRefs(gw *gatewayv1.Gateway, listenerName string, status bool, message string) (ok bool) {
	condition := metav1.Condition{
		Type:               string(gatewayv1.ListenerConditionResolvedRefs),
		Status:             ConditionStatus(status),
		Reason:             string(gatewayv1.ListenerReasonResolvedRefs),
		ObservedGeneration: gw.GetGeneration(),
		Message:            message,
		LastTransitionTime: metav1.Now(),
	}

	if !IsConditionPresentAndEqual(gw.Status.Conditions, condition) {
		setListenerCondition(gw, listenerName, condition)
		ok = true
	}
	return
}

func setListenerCondition(gw *gatewayv1.Gateway, listenerName string, newCondition metav1.Condition) {
	for i, listener := range gw.Status.Listeners {
		if listener.Name == gatewayv1.SectionName(listenerName) {
			gw.Status.Listeners[i].Conditions = MergeCondition(listener.Conditions, newCondition)
			return
		}
	}
}

func MergeCondition(conditions []metav1.Condition, newCondition metav1.Condition) []metav1.Condition {
	if newCondition.LastTransitionTime.IsZero() {
		newCondition.LastTransitionTime = metav1.Now()
	}
	newConditions := []metav1.Condition{}
	for _, condition := range conditions {
		if condition.Type != newCondition.Type {
			newConditions = append(newConditions, condition)
		}
	}
	newConditions = append(newConditions, newCondition)
	return newConditions
}

func HostnamesIntersect(a, b string) bool {
	return HostnamesMatch(a, b) || HostnamesMatch(b, a)
}

// HostnamesMatch checks that the hostnameB matches the hostnameA. HostnameA is treated as mask
// to be checked against the hostnameB.
func HostnamesMatch(hostnameA, hostnameB string) bool {
	// the hostnames are in the form of "foo.bar.com"; split them
	// in a slice of substrings
	hostnameALabels := strings.Split(hostnameA, ".")
	hostnameBLabels := strings.Split(hostnameB, ".")

	var a, b int
	var wildcard bool

	// iterate over the parts of both the hostnames
	for a, b = 0, 0; a < len(hostnameALabels) && b < len(hostnameBLabels); a, b = a+1, b+1 {
		var matchFound bool

		// if the current part of B is a wildcard, we need to find the first
		// A part that matches with the following B part
		if wildcard {
			for ; b < len(hostnameBLabels); b++ {
				if hostnameALabels[a] == hostnameBLabels[b] {
					matchFound = true
					break
				}
			}
		}

		// if no match was found, the hostnames don't match
		if wildcard && !matchFound {
			return false
		}

		// check if at least on of the current parts are a wildcard; if so, continue
		if hostnameALabels[a] == "*" {
			wildcard = true
			continue
		}
		// reset the wildcard  variables
		wildcard = false

		// if the current a part is different from the b part, the hostnames are incompatible
		if hostnameALabels[a] != hostnameBLabels[b] {
			return false
		}
	}
	return len(hostnameBLabels)-b == len(hostnameALabels)-a
}
