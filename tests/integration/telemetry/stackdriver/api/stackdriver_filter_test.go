//go:build integ
// +build integ

// Copyright Istio Authors. All Rights Reserved.
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

package stackdriver

import (
	"context"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
)

import (
	"golang.org/x/sync/errgroup"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/test/env"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/istio"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/stackdriver"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/label"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/resource"
	"github.com/apache/dubbo-go-pixiu/pkg/test/scopes"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/retry"
	"github.com/apache/dubbo-go-pixiu/tests/integration/telemetry"
	stackdrivertest "github.com/apache/dubbo-go-pixiu/tests/integration/telemetry/stackdriver"
)

const (
	serverRequestCount = "tests/integration/telemetry/stackdriver/testdata/server_request_count.json.tmpl"
	clientRequestCount = "tests/integration/telemetry/stackdriver/testdata/client_request_count.json.tmpl"
	serverLogEntry     = "tests/integration/telemetry/stackdriver/testdata/unsampled_server_access_log.json.tmpl"

	fakeGCEMetadataServerValues = `
meshConfig:
  defaultConfig:
    proxyMetadata:
      GCE_METADATA_HOST: `
)

// TestStackdriverMonitoring verifies that stackdriver WASM filter exports metrics with expected labels.
func TestStackdriverMonitoring(t *testing.T) {
	framework.NewTest(t).
		Features("observability.telemetry.stackdriver.api").
		Run(func(t framework.TestContext) {
			g, _ := errgroup.WithContext(context.Background())
			for _, cltInstance := range stackdrivertest.Clt {
				cltInstance := cltInstance
				g.Go(func() error {
					err := retry.UntilSuccess(func() error {
						if err := stackdrivertest.SendTraffic(cltInstance, http.Header{}, false); err != nil {
							return err
						}
						clName := cltInstance.Config().Cluster.Name()
						trustDomain := telemetry.GetTrustDomain(cltInstance.Config().Cluster, stackdrivertest.Ist.Settings().SystemNamespace)
						scopes.Framework.Infof("Validating for cluster %s", clName)

						// Validate cluster names in telemetry below once https://github.com/istio/istio/issues/28125 is fixed.
						if err := stackdrivertest.ValidateMetrics(t, filepath.Join(env.IstioSrc, serverRequestCount), filepath.Join(env.IstioSrc, clientRequestCount),
							clName, trustDomain); err != nil {
							return err
						}
						t.Logf("Metrics validated")

						if err := stackdrivertest.ValidateLogs(t, filepath.Join(env.IstioSrc, serverLogEntry), clName, trustDomain, stackdriver.ServerAccessLog); err != nil {
							return err
						}
						t.Logf("logs validated")
						// TODO: add trace validation

						return nil
					}, retry.Delay(framework.TelemetryRetryDelay), retry.Timeout(framework.TelemetryRetryTimeout))
					if err != nil {
						return err
					}
					return nil
				})
			}
			if err := g.Wait(); err != nil {
				t.Fatalf("test failed: %v", err)
			}
		})
}

func TestMain(m *testing.M) {
	framework.NewSuite(m).
		Label(label.CustomSetup).
		Label(label.IPv4). // We get detected as on GCE, since our tests run there, but don't have connectivity
		Setup(stackdrivertest.ConditionallySetupMetadataServer).
		Setup(istio.Setup(&stackdrivertest.Ist, setupConfig)).
		Setup(func(ctx resource.Context) error {
			i, err := istio.Get(ctx)
			if err != nil {
				return err
			}
			return ctx.ConfigIstio().YAML(i.Settings().SystemNamespace, `
apiVersion: telemetry.istio.io/v1alpha1
kind: Telemetry
metadata:
  name: mesh-default
  namespace: istio-system
spec:
  accessLogging:
  - providers:
    - name: stackdriver
  metrics:
  - providers:
    - name: stackdriver
`).Apply()
		}).
		Setup(stackdrivertest.TestSetup).
		Run()
}

func setupConfig(_ resource.Context, cfg *istio.Config) {
	if cfg == nil {
		return
	}

	// conditionally use a fake metadata server for testing off of GCP
	if stackdrivertest.GCEInst != nil {
		cfg.ControlPlaneValues = strings.Join([]string{cfg.ControlPlaneValues, fakeGCEMetadataServerValues, stackdrivertest.GCEInst.Address()}, "")
	}
}
