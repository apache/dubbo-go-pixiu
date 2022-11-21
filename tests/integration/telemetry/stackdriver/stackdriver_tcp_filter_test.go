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
	"path/filepath"
	"testing"
)

import (
	"golang.org/x/sync/errgroup"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/test/env"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/stackdriver"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/retry"
	"github.com/apache/dubbo-go-pixiu/tests/integration/telemetry"
)

const (
	tcpServerConnectionCount = "tests/integration/telemetry/stackdriver/testdata/server_tcp_connection_count.json.tmpl"
	tcpClientConnectionCount = "tests/integration/telemetry/stackdriver/testdata/client_tcp_connection_count.json.tmpl"
	tcpServerLogEntry        = "tests/integration/telemetry/stackdriver/testdata/tcp_server_access_log.json.tmpl"
)

// TestTCPStackdriverMonitoring verifies that stackdriver TCP filter works.
func TestTCPStackdriverMonitoring(t *testing.T) {
	framework.NewTest(t).
		Features("observability.telemetry.stackdriver").
		Run(func(t framework.TestContext) {
			g, _ := errgroup.WithContext(context.Background())
			for _, cltInstance := range Clt {
				cltInstance := cltInstance
				g.Go(func() error {
					err := retry.UntilSuccess(func() error {
						_, err := cltInstance.Call(echo.CallOptions{
							To: Srv,
							Port: echo.Port{
								Name: "tcp",
							},
							Count: telemetry.RequestCountMultipler * Srv.WorkloadsOrFail(t).Len(),
							Retry: echo.Retry{
								NoRetry: true,
							},
						})
						if err != nil {
							return err
						}
						t.Logf("Validating Telemetry for Cluster %v", cltInstance.Config().Cluster.Name())
						clName := cltInstance.Config().Cluster.Name()
						trustDomain := telemetry.GetTrustDomain(cltInstance.Config().Cluster, Ist.Settings().SystemNamespace)
						if err := ValidateMetrics(t, filepath.Join(env.IstioSrc, tcpServerConnectionCount),
							filepath.Join(env.IstioSrc, tcpClientConnectionCount), clName, trustDomain); err != nil {
							return err
						}
						if err := ValidateLogs(t, filepath.Join(env.IstioSrc, tcpServerLogEntry), clName, trustDomain, stackdriver.ServerAccessLog); err != nil {
							return err
						}

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
