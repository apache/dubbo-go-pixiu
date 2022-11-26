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

package customizemetrics

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config/protocol"
	"github.com/apache/dubbo-go-pixiu/pkg/test/env"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/containerregistry"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo/deployment"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo/match"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/istio"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/namespace"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/prometheus"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/label"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/resource"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/retry"
	util "github.com/apache/dubbo-go-pixiu/tests/integration/telemetry"
	common "github.com/apache/dubbo-go-pixiu/tests/integration/telemetry/stats/prometheus"
)

var (
	client, server echo.Instances
	appNsInst      namespace.Instance
	promInst       prometheus.Instance
	registry       containerregistry.Instance
)

const (
	removedTag            = "source_principal"
	requestCountMultipler = 3
	httpProtocol          = "http"
	grpcProtocol          = "grpc"

	// Same user name and password as specified at pkg/test/fakes/imageregistry
	registryUser   = "user"
	registryPasswd = "passwd"
)

func TestCustomizeMetrics(t *testing.T) {
	framework.NewTest(t).
		Label(label.IPv4). // https://github.com/istio/istio/issues/35835
		Features("observability.telemetry.stats.prometheus.customize-metric").
		Features("observability.telemetry.request-classification").
		Features("extensibility.wasm.remote-load").
		Run(func(t framework.TestContext) {
			t.Cleanup(func() {
				if t.Failed() {
					util.PromDump(t.Clusters().Default(), promInst, prometheus.Query{Metric: "istio_requests_total"})
				}
			})
			httpDestinationQuery := buildQuery(httpProtocol)
			grpcDestinationQuery := buildQuery(grpcProtocol)
			var httpMetricVal string
			cluster := t.Clusters().Default()
			httpChecked := false
			retry.UntilSuccessOrFail(t, func() error {
				if err := sendTraffic(); err != nil {
					t.Log("failed to send traffic")
					return err
				}
				var err error
				if !httpChecked {
					httpMetricVal, err = common.QueryPrometheus(t, t.Clusters().Default(), httpDestinationQuery, promInst)
					if err != nil {
						return err
					}
					httpChecked = true
				}
				_, err = common.QueryPrometheus(t, t.Clusters().Default(), grpcDestinationQuery, promInst)
				if err != nil {
					return err
				}
				return nil
			}, retry.Delay(1*time.Second), retry.Timeout(300*time.Second))
			// check tag removed
			if strings.Contains(httpMetricVal, removedTag) {
				t.Errorf("failed to remove tag: %v", removedTag)
			}
			common.ValidateMetric(t, cluster, promInst, httpDestinationQuery, 1)
			common.ValidateMetric(t, cluster, promInst, grpcDestinationQuery, 1)
		})
}

func TestMain(m *testing.M) {
	framework.NewSuite(m).
		Label(label.CustomSetup).
		Label(label.IPv4). // https://github.com/istio/istio/issues/35915
		Setup(istio.Setup(common.GetIstioInstance(), setupConfig)).
		Setup(testSetup).
		Setup(setupWasmExtension).
		Run()
}

func testSetup(ctx resource.Context) (err error) {
	// enable custom tag in the stats
	bootstrapPatch := `
apiVersion: networking.istio.io/v1alpha3
kind: EnvoyFilter
metadata:
  name: bootstrap-tag
  namespace: istio-system
spec:
  configPatches:
  - applyTo: BOOTSTRAP
    patch:
      operation: MERGE
      value:
        stats_config:
          stats_tags:
          - regex: "(custom_dimension=\\.=(.*?);\\.;)"
            tag_name: "custom_dimension"
`
	if err := ctx.ConfigIstio().YAML("istio-system", bootstrapPatch).Apply(resource.Wait); err != nil {
		return err
	}

	registry, err = containerregistry.New(ctx, containerregistry.Config{Cluster: ctx.AllClusters().Default()})
	if err != nil {
		return
	}

	var nsErr error
	appNsInst, nsErr = namespace.New(ctx, namespace.Config{
		Prefix: "echo",
		Inject: true,
	})
	if nsErr != nil {
		return nsErr
	}
	proxyMetadata := fmt.Sprintf(`
proxyMetadata:
  BOOTSTRAP_XDS_AGENT: "true"
  WASM_INSECURE_REGISTRIES: %q`, registry.Address())

	echos, err := deployment.New(ctx).
		WithClusters(ctx.Clusters()...).
		WithConfig(echo.Config{
			Service:   "client",
			Namespace: appNsInst,
			Ports:     nil,
			Subsets: []echo.SubsetConfig{
				{
					Annotations: map[echo.Annotation]*echo.AnnotationValue{
						echo.SidecarProxyConfig: {
							Value: proxyMetadata,
						},
					},
				},
			},
		}).
		WithConfig(echo.Config{
			Service:   "server",
			Namespace: appNsInst,
			Subsets: []echo.SubsetConfig{
				{
					Annotations: map[echo.Annotation]*echo.AnnotationValue{
						echo.SidecarProxyConfig: {
							Value: proxyMetadata,
						},
					},
				},
			},
			Ports: []echo.Port{
				{
					Name:         "http",
					Protocol:     protocol.HTTP,
					WorkloadPort: 8090,
				},
				{
					Name:         "grpc",
					Protocol:     protocol.GRPC,
					WorkloadPort: 7070,
				},
			},
		}).
		Build()
	if err != nil {
		return err
	}
	client = match.ServiceName(echo.NamespacedName{Name: "client", Namespace: appNsInst}).GetMatches(echos)
	server = match.ServiceName(echo.NamespacedName{Name: "server", Namespace: appNsInst}).GetMatches(echos)
	promInst, err = prometheus.New(ctx, prometheus.Config{})
	if err != nil {
		return
	}
	return nil
}

func setupConfig(_ resource.Context, cfg *istio.Config) {
	if cfg == nil {
		return
	}
	cfValue := `
values:
 telemetry:
   v2:
     prometheus:
       configOverride:
         inboundSidecar:
           debug: false
           stat_prefix: istio
           metrics:
           - name: requests_total
             dimensions:
               response_code: istio_responseClass
               request_operation: istio_operationId
               grpc_response_status: istio_grpcResponseStatus
               custom_dimension: "'test'"
             tags_to_remove:
             - %s
`
	cfg.ControlPlaneValues = fmt.Sprintf(cfValue, removedTag)
}

func setupWasmExtension(ctx resource.Context) error {
	proxySHA, err := env.ReadProxySHA()
	if err != nil {
		return err
	}
	attrGenURL := fmt.Sprintf("https://storage.googleapis.com/istio-build/proxy/attributegen-%v.wasm", proxySHA)
	attrGenImageURL := fmt.Sprintf("oci://%v/istio-testing/wasm/attributegen:%v", registry.Address(), proxySHA)
	useRemoteWasmModule := false
	resp, err := http.Get(attrGenURL)
	if err == nil && resp.StatusCode == http.StatusOK {
		useRemoteWasmModule = true
	}

	args := map[string]interface{}{
		"WasmRemoteLoad":  useRemoteWasmModule,
		"AttributeGenURL": attrGenImageURL,
		"DockerConfigJson": base64.StdEncoding.EncodeToString(
			[]byte(createDockerCredential(registryUser, registryPasswd, registry.Address()))),
	}
	if err := ctx.ConfigIstio().EvalFile(appNsInst.Name(), args, "testdata/attributegen_envoy_filter.yaml").Apply(); err != nil {
		return err
	}

	return nil
}

func sendTraffic() error {
	for _, cltInstance := range client {
		count := requestCountMultipler * server.MustWorkloads().Len()
		httpOpts := echo.CallOptions{
			To: server,
			Port: echo.Port{
				Name: "http",
			},
			HTTP: echo.HTTP{
				Path:   "/path",
				Method: "GET",
			},
			Count: count,
			Retry: echo.Retry{
				NoRetry: true,
			},
		}

		if _, err := cltInstance.Call(httpOpts); err != nil {
			return err
		}

		httpOpts.HTTP.Method = "POST"
		if _, err := cltInstance.Call(httpOpts); err != nil {
			return err
		}

		grpcOpts := echo.CallOptions{
			To: server,
			Port: echo.Port{
				Name: "grpc",
			},
			Count: count,
		}
		if _, err := cltInstance.Call(grpcOpts); err != nil {
			return err
		}
	}

	return nil
}

func buildQuery(protocol string) (destinationQuery prometheus.Query) {
	labels := map[string]string{
		"request_protocol":               "http",
		"response_code":                  "2xx",
		"destination_app":                "server",
		"destination_version":            "v1",
		"destination_service":            "server." + appNsInst.Name() + ".svc.cluster.local",
		"destination_service_name":       "server",
		"destination_workload_namespace": appNsInst.Name(),
		"destination_service_namespace":  appNsInst.Name(),
		"source_app":                     "client",
		"source_version":                 "v1",
		"source_workload":                "client-v1",
		"source_workload_namespace":      appNsInst.Name(),
		"custom_dimension":               "test",
	}
	if protocol == httpProtocol {
		labels["request_operation"] = "getoperation"
	}
	if protocol == grpcProtocol {
		labels["grpc_response_status"] = "OK"
		labels["request_protocol"] = "grpc"
	}

	_, destinationQuery, _ = common.BuildQueryCommon(labels, appNsInst.Name())
	return destinationQuery
}

func createDockerCredential(user, passwd, registry string) string {
	credentials := `{
	"auths":{
		"%v":{
			"username": "%v",
			"password": "%v",
			"email": "test@abc.com",
			"auth": "%v"
		}
	}
}`
	auth := base64.StdEncoding.EncodeToString([]byte(user + ":" + passwd))
	return fmt.Sprintf(credentials, registry, user, passwd, auth)
}
