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

package istiocontrolplane

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
)

import (
	"istio.io/api/operator/v1alpha1"
	"istio.io/pkg/log"
	"istio.io/pkg/version"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	kubeversion "k8s.io/apimachinery/pkg/version"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"sigs.k8s.io/yaml"
)

import (
	"github.com/apache/dubbo-go-pixiu/operator/pkg/apis/istio"
	iopv1alpha1 "github.com/apache/dubbo-go-pixiu/operator/pkg/apis/istio/v1alpha1"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/cache"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/helm"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/helmreconciler"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/metrics"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/name"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/object"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/tpath"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/translate"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/util"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/util/clog"
	"github.com/apache/dubbo-go-pixiu/operator/pkg/util/progress"
	"github.com/apache/dubbo-go-pixiu/pkg/errdict"
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
	"github.com/apache/dubbo-go-pixiu/pkg/url"
)

const (
	finalizer = "istio-finalizer.install.istio.io"
	// finalizerMaxRetries defines the maximum number of attempts to remove the finalizer.
	finalizerMaxRetries = 1
	// IgnoreReconcileAnnotation is annotation of IstioOperator CR so it would be ignored during Reconcile loop.
	IgnoreReconcileAnnotation = "install.istio.io/ignoreReconcile"
)

var (
	scope      = log.RegisterScope("installer", "installer", 0)
	restConfig *rest.Config
)

type Options struct {
	Force bool
}

const (
	autoscalingV2MinK8SVersion = 23
	pdbV1MinK8SVersion         = 21
)

// watchedResources contains all resources we will watch and reconcile when changed
// Ideally this would also contain Istio CRDs, but there is a race condition here - we cannot watch
// a type that does not yet exist.
func watchedResources(version *kubeversion.Info) []schema.GroupVersionKind {
	res := []schema.GroupVersionKind{
		{Group: "apps", Version: "v1", Kind: name.DeploymentStr},
		{Group: "apps", Version: "v1", Kind: name.DaemonSetStr},
		{Group: "", Version: "v1", Kind: name.ServiceStr},
		// Endpoints should not be pruned because these are generated and not in the manifest.
		// {Group: "", Version: "v1", Kind: name.EndpointStr},
		{Group: "", Version: "v1", Kind: name.CMStr},
		{Group: "", Version: "v1", Kind: name.PVCStr},
		{Group: "", Version: "v1", Kind: name.PodStr},
		{Group: "", Version: "v1", Kind: name.SecretStr},
		{Group: "", Version: "v1", Kind: name.SAStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.RoleBindingStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.RoleStr},
		{Group: "admissionregistration.k8s.io", Version: "v1", Kind: name.MutatingWebhookConfigurationStr},
		{Group: "admissionregistration.k8s.io", Version: "v1", Kind: name.ValidatingWebhookConfigurationStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.ClusterRoleStr},
		{Group: "rbac.authorization.k8s.io", Version: "v1", Kind: name.ClusterRoleBindingStr},
		{Group: "apiextensions.k8s.io", Version: "v1", Kind: name.CRDStr},
	}
	// autoscaling v2 API is available on >=1.23
	if kube.IsKubeAtLeastOrLessThanVersion(version, autoscalingV2MinK8SVersion, true) {
		res = append(res, schema.GroupVersionKind{Group: "autoscaling", Version: "v2", Kind: name.HPAStr})
	} else {
		res = append(res, schema.GroupVersionKind{Group: "autoscaling", Version: "v2beta2", Kind: name.HPAStr})
	}
	// policy/v1 is available on >=1.21
	if kube.IsKubeAtLeastOrLessThanVersion(version, pdbV1MinK8SVersion, true) {
		res = append(res, schema.GroupVersionKind{Group: "policy", Version: "v1", Kind: name.PDBStr})
	} else {
		res = append(res, schema.GroupVersionKind{Group: "policy", Version: "v1beta1", Kind: name.PDBStr})
	}
	return res
}

var (
	ownedResourcePredicates = predicate.Funcs{
		CreateFunc: func(_ event.CreateEvent) bool {
			// no action
			return false
		},
		GenericFunc: func(_ event.GenericEvent) bool {
			// no action
			return false
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			obj, err := meta.Accessor(e.Object)
			scope.Debugf("got delete event for %s.%s", obj.GetName(), obj.GetNamespace())
			if err != nil {
				return false
			}
			unsObj, err := runtime.DefaultUnstructuredConverter.ToUnstructured(e.Object)
			if err != nil {
				return false
			}
			if isOperatorCreatedResource(obj) {
				crName := obj.GetLabels()[helmreconciler.OwningResourceName]
				crNamespace := obj.GetLabels()[helmreconciler.OwningResourceNamespace]
				componentName := obj.GetLabels()[helmreconciler.IstioComponentLabelStr]
				var host string
				if restConfig != nil {
					host = restConfig.Host
				}
				crHash := strings.Join([]string{crName, crNamespace, componentName, host}, "-")
				oh := object.NewK8sObject(&unstructured.Unstructured{Object: unsObj}, nil, nil).Hash()
				cache.RemoveObject(crHash, oh)
				return true
			}
			return false
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			// no action
			return false
		},
	}

	operatorPredicates = predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
		UpdateFunc: func(e event.UpdateEvent) bool {
			oldIOP, ok := e.ObjectOld.(*iopv1alpha1.IstioOperator)
			if !ok {
				scope.Error(errdict.OperatorFailedToGetObjectInCallback, "failed to get old IstioOperator")
				return false
			}
			newIOP := e.ObjectNew.(*iopv1alpha1.IstioOperator)
			if !ok {
				scope.Error(errdict.OperatorFailedToGetObjectInCallback, "failed to get new IstioOperator")
				return false
			}
			if !reflect.DeepEqual(oldIOP.Spec, newIOP.Spec) ||
				oldIOP.GetDeletionTimestamp() != newIOP.GetDeletionTimestamp() ||
				oldIOP.GetGeneration() != newIOP.GetGeneration() {
				return true
			}
			return false
		},
	}
)

// NewReconcileIstioOperator creates a new ReconcileIstioOperator and returns a ptr to it.
func NewReconcileIstioOperator(client client.Client, kubeClient kube.Client, scheme *runtime.Scheme) *ReconcileIstioOperator {
	return &ReconcileIstioOperator{
		client:     client,
		kubeClient: kubeClient,
		scheme:     scheme,
	}
}

// ReconcileIstioOperator reconciles a IstioOperator object
type ReconcileIstioOperator struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client     client.Client
	kubeClient kube.Client
	scheme     *runtime.Scheme
	options    *Options
}

// Reconcile reads that state of the cluster for a IstioOperator object and makes changes based on the state read
// and what is in the IstioOperator.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileIstioOperator) Reconcile(_ context.Context, request reconcile.Request) (reconcile.Result, error) {
	scope.Info("Reconciling IstioOperator")

	ns, iopName := request.Namespace, request.Name
	reqNamespacedName := types.NamespacedName{
		Name:      request.Name,
		Namespace: ns,
	}
	// declare read-only iop instance to create the reconciler
	iop := &iopv1alpha1.IstioOperator{}
	if err := r.client.Get(context.TODO(), reqNamespacedName, iop); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			metrics.CRDeletionTotal.Increment()
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		scope.Warnf(errdict.OperatorFailedToGetObjectFromAPIServer, "error getting IstioOperator %s: %s", iopName, err)
		metrics.CountCRFetchFail(errors.ReasonForError(err))
		return reconcile.Result{}, err
	}
	if iop.Spec == nil {
		iop.Spec = &v1alpha1.IstioOperatorSpec{Profile: name.DefaultProfileName}
	}
	operatorRevision, _ := os.LookupEnv("REVISION")
	if operatorRevision != "" && operatorRevision != iop.Spec.Revision {
		scope.Infof("Ignoring IstioOperator CR %s with revision %s, since operator revision is %s.", iopName, iop.Spec.Revision, operatorRevision)
		return reconcile.Result{}, nil
	}
	if iop.Annotations != nil {
		if ir := iop.Annotations[IgnoreReconcileAnnotation]; ir == "true" {
			scope.Infof("Ignoring the IstioOperator CR %s because it is annotated to be ignored for reconcile ", iopName)
			return reconcile.Result{}, nil
		}
	}

	// for backward compatibility, the previous applied installed-state CR does not have the ignore reconcile annotation
	// TODO(richardwxn): remove this check and rely on annotation check only
	if strings.HasPrefix(iop.Name, name.InstalledSpecCRPrefix) {
		scope.Infof("Ignoring the installed-state IstioOperator CR %s ", iopName)
		return reconcile.Result{}, nil
	}

	var err error
	iopMerged := &iopv1alpha1.IstioOperator{}
	*iopMerged = *iop
	// get the merged values in iop on top of the defaults for the profile given by iop.profile
	iopMerged.Spec, err = mergeIOPSWithProfile(iopMerged)
	if err != nil {
		scope.Errorf(errdict.OperatorFailedToMergeUserIOP, "failed to merge base profile with user IstioOperator CR %s, %s", iopName, err)
		return reconcile.Result{}, err
	}

	deleted := iop.GetDeletionTimestamp() != nil
	finalizers := sets.NewString(iop.GetFinalizers()...)
	if deleted {
		if !finalizers.Has(finalizer) {
			scope.Infof("IstioOperator %s deleted", iopName)
			metrics.CRDeletionTotal.Increment()
			return reconcile.Result{}, nil
		}
		scope.Infof("Deleting IstioOperator %s", iopName)

		reconciler, err := helmreconciler.NewHelmReconciler(r.client, r.kubeClient, iopMerged, nil)
		if err != nil {
			return reconcile.Result{}, err
		}
		if err := reconciler.Delete(); err != nil {
			scope.Errorf("Failed to delete resources with helm reconciler: %s.", err)
			return reconcile.Result{}, err
		}

		finalizers.Delete(finalizer)
		iop.SetFinalizers(finalizers.List())
		finalizerError := r.client.Update(context.TODO(), iop)
		for retryCount := 0; errors.IsConflict(finalizerError) && retryCount < finalizerMaxRetries; retryCount++ {
			scope.Info("API server conflict during finalizer removal, retrying.")
			_ = r.client.Get(context.TODO(), request.NamespacedName, iop)
			finalizers = sets.NewString(iop.GetFinalizers()...)
			finalizers.Delete(finalizer)
			iop.SetFinalizers(finalizers.List())
			finalizerError = r.client.Update(context.TODO(), iop)
		}
		if finalizerError != nil {
			if errors.IsNotFound(finalizerError) {
				scope.Infof("Did not remove finalizer from %s: the object was previously deleted.", iopName)
				metrics.CRDeletionTotal.Increment()
				return reconcile.Result{}, nil
			} else if errors.IsConflict(finalizerError) {
				scope.Infof("Could not remove finalizer from %s due to conflict. Operation will be retried in next reconcile attempt.", iopName)
				return reconcile.Result{}, nil
			}
			scope.Errorf(errdict.OperatorFailedToRemoveFinalizer, "error removing finalizer: %s", finalizerError)
			return reconcile.Result{}, finalizerError
		}
		return reconcile.Result{}, nil
	} else if !finalizers.Has(finalizer) {
		log.Infof("Adding finalizer %v to %v", finalizer, request)
		finalizers.Insert(finalizer)
		iop.SetFinalizers(finalizers.List())
		err := r.client.Update(context.TODO(), iop)
		if err != nil {
			if errors.IsNotFound(err) {
				scope.Infof("Could not add finalizer to %s: the object was deleted.", iopName)
				metrics.CRDeletionTotal.Increment()
				return reconcile.Result{}, nil
			} else if errors.IsConflict(err) {
				scope.Infof("Could not add finalizer to %s due to conflict. Operation will be retried in next reconcile attempt.", iopName)
				return reconcile.Result{}, nil
			}
			scope.Errorf(errdict.OperatorFailedToAddFinalizer, "Failed to add finalizer to IstioOperator CR %s: %s", iopName, err)
			return reconcile.Result{}, err
		}
	}

	scope.Info("Updating IstioOperator")
	val := iopMerged.Spec.Values.AsMap()
	if _, ok := val["global"]; !ok {
		val["global"] = make(map[string]interface{})
	}
	globalValues := val["global"].(map[string]interface{})
	scope.Info("Detecting third-party JWT support")
	var jwtPolicy util.JWTPolicy
	if jwtPolicy, err = util.DetectSupportedJWTPolicy(r.kubeClient); err != nil {
		// TODO(howardjohn): add to dictionary. When resolved, replace this sentence with Done or WontFix - if WontFix, add reason.
		scope.Warnf("Failed to detect third-party JWT support: %v", err)
	} else {
		if jwtPolicy == util.FirstPartyJWT {
			scope.Info("Detected that your cluster does not support third party JWT authentication. " +
				"Falling back to less secure first party JWT. " +
				"See " + url.ConfigureSAToken + " for details.")
		}
		globalValues["jwtPolicy"] = string(jwtPolicy)
	}
	err = util.ValidateIOPCAConfig(r.kubeClient, iopMerged)
	if err != nil {
		scope.Errorf(errdict.OperatorFailedToConfigure, "failed to apply IstioOperator resources. Error %s", err)
		return reconcile.Result{}, err
	}
	helmReconcilerOptions := &helmreconciler.Options{
		Log:         clog.NewDefaultLogger(),
		ProgressLog: progress.NewLog(),
	}
	if r.options != nil {
		helmReconcilerOptions.Force = r.options.Force
	}
	reconciler, err := helmreconciler.NewHelmReconciler(r.client, r.kubeClient, iopMerged, helmReconcilerOptions)
	if err != nil {
		return reconcile.Result{}, err
	}
	if err := reconciler.SetStatusBegin(); err != nil {
		return reconcile.Result{}, err
	}
	status, err := reconciler.Reconcile()
	if err != nil {
		scope.Errorf("Error during reconcile: %s", err)
	}
	if err := reconciler.SetStatusComplete(status); err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, err
}

// mergeIOPSWithProfile overlays the values in iop on top of the defaults for the profile given by iop.profile and
// returns the merged result.
func mergeIOPSWithProfile(iop *iopv1alpha1.IstioOperator) (*v1alpha1.IstioOperatorSpec, error) {
	profileYAML, err := helm.GetProfileYAML(iop.Spec.InstallPackagePath, iop.Spec.Profile)
	if err != nil {
		metrics.CountCRMergeFail(metrics.CannotFetchProfileError)
		return nil, err
	}

	// Due to the fact that base profile is compiled in before a tag can be created, we must allow an additional
	// override from variables that are set during release build time.
	hub := version.DockerInfo.Hub
	tag := version.DockerInfo.Tag
	if hub != "" && hub != "unknown" && tag != "" && tag != "unknown" {
		buildHubTagOverlayYAML, err := helm.GenerateHubTagOverlay(hub, tag)
		if err != nil {
			metrics.CountCRMergeFail(metrics.OverlayError)
			return nil, err
		}
		profileYAML, err = util.OverlayYAML(profileYAML, buildHubTagOverlayYAML)
		if err != nil {
			metrics.CountCRMergeFail(metrics.OverlayError)
			return nil, err
		}
	}

	overlayYAMLB, err := yaml.Marshal(iop)
	if err != nil {
		metrics.CountCRMergeFail(metrics.IOPFormatError)
		return nil, err
	}
	overlayYAML := string(overlayYAMLB)
	t := translate.NewReverseTranslator()
	overlayYAML, err = t.TranslateK8SfromValueToIOP(overlayYAML)
	if err != nil {
		metrics.CountCRMergeFail(metrics.TranslateValuesError)
		return nil, fmt.Errorf("could not overlay k8s settings from values to IOP: %s", err)
	}

	mergedYAML, err := util.OverlayIOP(profileYAML, overlayYAML)
	if err != nil {
		metrics.CountCRMergeFail(metrics.OverlayError)
		return nil, err
	}

	mergedYAML, err = translate.OverlayValuesEnablement(mergedYAML, overlayYAML, "")
	if err != nil {
		metrics.CountCRMergeFail(metrics.TranslateValuesError)
		return nil, err
	}

	mergedYAMLSpec, err := tpath.GetSpecSubtree(mergedYAML)
	if err != nil {
		metrics.CountCRMergeFail(metrics.InternalYAMLParseError)
		return nil, err
	}

	return istio.UnmarshalAndValidateIOPS(mergedYAMLSpec)
}

// Add creates a new IstioOperator Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started. It also provides additional options to modify internal reconciler behavior.
func Add(mgr manager.Manager, options *Options) error {
	restConfig = mgr.GetConfig()
	kubeClient, err := kube.NewExtendedClient(kube.NewClientConfigForRestConfig(restConfig), "")
	if err != nil {
		return fmt.Errorf("create Kubernetes client: %v", err)
	}
	return add(mgr, &ReconcileIstioOperator{client: mgr.GetClient(), scheme: mgr.GetScheme(), kubeClient: kubeClient, options: options})
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r *ReconcileIstioOperator) error {
	scope.Info("Adding controller for IstioOperator.")
	// Create a new controller
	c, err := controller.New("istiocontrolplane-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource IstioOperator
	err = c.Watch(&source.Kind{Type: &iopv1alpha1.IstioOperator{}}, &handler.EnqueueRequestForObject{}, operatorPredicates)
	if err != nil {
		return err
	}
	ver, err := r.kubeClient.GetKubernetesVersion()
	if err != nil {
		return err
	}
	// watch for changes to Istio resources
	err = watchIstioResources(c, ver)
	if err != nil {
		return err
	}
	scope.Info("Controller added")
	return nil
}

// Watch changes for Istio resources managed by the operator
func watchIstioResources(c controller.Controller, ver *kubeversion.Info) error {
	for _, t := range watchedResources(ver) {
		u := &unstructured.Unstructured{}
		u.SetGroupVersionKind(schema.GroupVersionKind{
			Kind:    t.Kind,
			Group:   t.Group,
			Version: t.Version,
		})
		err := c.Watch(&source.Kind{Type: u}, handler.EnqueueRequestsFromMapFunc(func(a client.Object) []reconcile.Request {
			scope.Infof("Watching a change for istio resource: %s/%s", a.GetNamespace(), a.GetName())
			return []reconcile.Request{
				{NamespacedName: types.NamespacedName{
					Name:      a.GetLabels()[helmreconciler.OwningResourceName],
					Namespace: a.GetLabels()[helmreconciler.OwningResourceNamespace],
				}},
			}
		}),
			ownedResourcePredicates)
		if err != nil {
			scope.Errorf("Could not create watch for %s/%s/%s: %s.", t.Kind, t.Group, t.Version, err)
		}
	}
	return nil
}

// Check if the specified object is created by operator
func isOperatorCreatedResource(obj metav1.Object) bool {
	return obj.GetLabels()[helmreconciler.OwningResourceName] != "" &&
		obj.GetLabels()[helmreconciler.OwningResourceNamespace] != "" &&
		obj.GetLabels()[helmreconciler.IstioComponentLabelStr] != ""
}
