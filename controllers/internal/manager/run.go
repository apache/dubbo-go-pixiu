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
	"crypto/tls"
	"os"
)

import (
	"controllers/api/v1alpha1"

	"controllers/internal/controller"
	"controllers/internal/controller/config"
	"controllers/internal/controller/status"

	"github.com/go-logr/logr"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	"k8s.io/utils/ptr"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/metrics/filters"
	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	gatewayv1 "sigs.k8s.io/gateway-api/apis/v1"
	"sigs.k8s.io/gateway-api/apis/v1beta1"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	if err := gatewayv1.Install(scheme); err != nil {
		panic(err)
	}
	if err := v1alpha1.AddToScheme(scheme); err != nil {
		panic(err)
	}
	if err := v1beta1.Install(scheme); err != nil {
		panic(err)
	}
	// +kubebuilder:scaffold:scheme
}

func Run(ctx context.Context, logger logr.Logger) error {
	cfg := config.ControllerConfig
	setupLog := ctrl.LoggerFrom(ctx).WithName("setup")
	disableHTTP2 := func(c *tls.Config) {
		setupLog.Info("disabling http/2")
		c.NextProtos = []string{"http/1.1"}
	}
	var tlsOpts []func(*tls.Config)
	if !cfg.EnableHTTP2 {
		tlsOpts = append(tlsOpts, disableHTTP2)
	}
	webhookServer := webhook.NewServer(webhook.Options{
		TLSOpts: tlsOpts,
	})
	metricsServerOptions := metricsserver.Options{
		BindAddress:   cfg.MetricsAddr,
		SecureServing: cfg.SecureMetrics,
		TLSOpts:       tlsOpts,
	}

	if cfg.SecureMetrics {
		metricsServerOptions.FilterProvider = filters.WithAuthenticationAndAuthorization
	}

	namespace := os.Getenv("POD_NAMESPACE")
	if namespace == "" {
		namespace = "default"
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                  scheme,
		Metrics:                 metricsServerOptions,
		WebhookServer:           webhookServer,
		HealthProbeBindAddress:  cfg.ProbeAddr,
		LeaderElection:          !config.ControllerConfig.LeaderElection.Disable,
		LeaderElectionID:        cfg.LeaderElectionID,
		LeaderElectionNamespace: namespace,
		LeaseDuration:           ptr.To(config.ControllerConfig.LeaderElection.LeaseDuration.Duration),
		RenewDeadline:           ptr.To(config.ControllerConfig.LeaderElection.RenewDeadline.Duration),
		RetryPeriod:             ptr.To(config.ControllerConfig.LeaderElection.RetryPeriod.Duration),
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		return err
	}

	updater := status.NewStatusUpdateHandler(ctrl.LoggerFrom(ctx).WithName("status").WithName("updater"), mgr.GetClient())
	if err := mgr.Add(updater); err != nil {
		setupLog.Error(err, "unable to add status updater")
		return err
	}

	setupLog.Info("check ReferenceGrants is enabled")
	_, err = mgr.GetRESTMapper().KindsFor(schema.GroupVersionResource{
		Group:    v1beta1.GroupVersion.Group,
		Version:  v1beta1.GroupVersion.Version,
		Resource: "referencegrants",
	})
	if err != nil {
		setupLog.Info("CRD ReferenceGrants is not installed", "err", err)
	}
	controller.SetEnableReferenceGrant(err == nil)
	setupLog.Info("setting up controllers")
	controllers, err := setupControllers(ctx, mgr, updater.Writer())
	if err != nil {
		setupLog.Error(err, "unable to set up controllers")
		return err
	}

	for _, c := range controllers {
		if err := c.SetupWithManager(mgr); err != nil {
			return err
		}
	}

	// +kubebuilder:scaffold:builder

	setupLog.Info("setting up health checks")
	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		return err
	}

	setupLog.Info("setting up ready checks")
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		return err
	}

	setupLog.Info("starting controller manager")
	return mgr.Start(ctrl.SetupSignalHandler())
}
