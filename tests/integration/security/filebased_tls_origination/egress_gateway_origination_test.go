//go:build integ
// +build integ

//  Copyright Istio Authors
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package filebasedtlsorigination

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

import (
	envoyAdmin "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config/protocol"
	"github.com/apache/dubbo-go-pixiu/pkg/http/headers"
	"github.com/apache/dubbo-go-pixiu/pkg/test"
	echoClient "github.com/apache/dubbo-go-pixiu/pkg/test/echo"
	"github.com/apache/dubbo-go-pixiu/pkg/test/echo/common"
	"github.com/apache/dubbo-go-pixiu/pkg/test/env"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo/check"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo/deployment"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/istio"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/namespace"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/resource"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/retry"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/structpath"
)

func mustReadCert(t framework.TestContext, f string) string {
	b, err := os.ReadFile(path.Join(env.IstioSrc, "tests/testdata/certs/dns", f))
	if err != nil {
		t.Fatalf("failed to read %v: %v", f, err)
	}
	return string(b)
}

// TestEgressGatewayTls brings up an cluster and will ensure that the TLS origination at
// egress gateway allows secure communication between the egress gateway and external workload.
// This test brings up an egress gateway to originate TLS connection. The test will ensure that requests
// are routed securely through the egress gateway and that the TLS origination happens at the gateway.
func TestEgressGatewayTls(t *testing.T) {
	framework.NewTest(t).
		Features("security.egress.tls.filebased").
		Run(func(t framework.TestContext) {
			internalClient, externalServer, _, serviceNamespace := setupEcho(t, t)
			// Set up Host Name
			host := "server." + serviceNamespace.Name() + ".svc.cluster.local"

			testCases := map[string]struct {
				destinationRuleMode string
				code                int
				gateway             bool //  If gateway is true, request is expected to pass through the egress gateway
				fakeRootCert        bool // If Fake root cert is to be used to verify server's presented certificate
			}{
				// Mutual Connection is originated by our DR but server side drops the connection to
				// only use Simple TLS as it doesn't verify client side cert
				// TODO: mechanism to enforce mutual TLS(client cert) validation by the server
				// 1. Mutual TLS origination from egress gateway to https endpoint:
				//    internalClient ) ---HTTP request (Host: some-external-site.com----> Hits listener 0.0.0.0_80 ->
				//      VS Routing (add Egress Header) --> Egress Gateway(originates mTLS with client certs)
				//      --> externalServer(443 with only Simple TLS used and client cert is not verified)
				"Mutual TLS origination from egress gateway to https endpoint": {
					destinationRuleMode: "MUTUAL",
					code:                http.StatusOK,
					gateway:             true,
					fakeRootCert:        false,
				},
				// 2. Simple TLS case:
				//    internalClient ) ---HTTP request (Host: some-external-site.com----> Hits listener 0.0.0.0_80 ->
				//      VS Routing (add Egress Header) --> Egress Gateway(originates TLS)
				//      --> externalServer(443 with TLS enforced)

				"SIMPLE TLS origination from egress gateway to https endpoint": {
					destinationRuleMode: "SIMPLE",
					code:                http.StatusOK,
					gateway:             true,
					fakeRootCert:        false,
				},
				// 3. No TLS case:
				//    internalClient ) ---HTTP request (Host: some-external-site.com----> Hits listener 0.0.0.0_80 ->
				//      VS Routing (add Egress Header) --> Egress Gateway(does not originate TLS)
				//      --> externalServer(443 with TLS enforced) request fails as gateway tries plain text only
				"No TLS origination from egress gateway to https endpoint": {
					destinationRuleMode: "DISABLE",
					code:                http.StatusBadRequest,
					gateway:             false, // 400 response will not contain header
				},
				// 5. SIMPLE TLS origination with "fake" root cert::
				//   internalClient ) ---HTTP request (Host: some-external-site.com----> Hits listener 0.0.0.0_80 ->
				//     VS Routing (add Egress Header) --> Egress Gateway(originates simple TLS)
				//     --> externalServer(443 with TLS enforced)
				//    request fails as the server cert can't be validated using the fake root cert used during origination
				"SIMPLE TLS origination from egress gateway to https endpoint with fake root cert": {
					destinationRuleMode: "SIMPLE",
					code:                http.StatusServiceUnavailable,
					gateway:             false, // 503 response will not contain header
					fakeRootCert:        true,
				},
			}

			for name, tc := range testCases {
				t.NewSubTest(name).
					Run(func(t framework.TestContext) {
						bufDestinationRule := createDestinationRule(t, serviceNamespace, tc.destinationRuleMode, tc.fakeRootCert)

						istioCfg := istio.DefaultConfigOrFail(t, t)
						systemNamespace := namespace.ClaimOrFail(t, t, istioCfg.SystemNamespace)

						t.ConfigIstio().YAML(systemNamespace.Name(), bufDestinationRule.String()).ApplyOrFail(t)

						opts := echo.CallOptions{
							To: externalServer,
							Port: echo.Port{
								Name: "http",
							},
							HTTP: echo.HTTP{
								Headers: headers.New().WithHost(host).Build(),
							},
							Retry: echo.Retry{
								Options: []retry.Option{retry.Delay(1 * time.Second), retry.Timeout(2 * time.Minute)},
							},
							Check: check.And(
								check.NoError(),
								check.Status(tc.code),
								check.Each(func(r echoClient.Response) error {
									if _, f := r.RequestHeaders["Handled-By-Egress-Gateway"]; tc.gateway && !f {
										return fmt.Errorf("expected to be handled by gateway. response: %s", r)
									}
									return nil
								})),
						}

						internalClient.CallOrFail(t, opts)
					})
			}
		})
}

const (
	// Destination Rule configs
	DestinationRuleConfigSimple = `
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: originate-tls-for-server-filebased-simple
spec:
  host: "server.{{.AppNamespace}}.svc.cluster.local"
  trafficPolicy:
    portLevelSettings:
      - port:
          number: 443
        tls:
          mode: {{.Mode}}
          caCertificates: {{.RootCertPath}}
          sni: server.{{.AppNamespace}}.svc.cluster.local

`
	// Destination Rule configs
	DestinationRuleConfigDisabledOrIstioMutual = `
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: originate-tls-for-server-filebased-disabled
spec:
  host: "server.{{.AppNamespace}}.svc.cluster.local"
  trafficPolicy:
    portLevelSettings:
      - port:
          number: 443
        tls:
          mode: {{.Mode}}
          sni: server.{{.AppNamespace}}.svc.cluster.local

`
	DestinationRuleConfigMutual = `
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: originate-tls-for-server-filebased-mutual
spec:
  host: "server.{{.AppNamespace}}.svc.cluster.local"
  trafficPolicy:
    portLevelSettings:
      - port:
          number: 443
        tls:
          mode: {{.Mode}}
          clientCertificate: /etc/certs/custom/cert-chain.pem
          privateKey: /etc/certs/custom/key.pem
          caCertificates: {{.RootCertPath}}
          sni: server.{{.AppNamespace}}.svc.cluster.local
`
)

func createDestinationRule(t framework.TestContext, serviceNamespace namespace.Instance,
	destinationRuleMode string, fakeRootCert bool) bytes.Buffer {
	var destinationRuleToParse string
	var rootCertPathToUse string
	if destinationRuleMode == "MUTUAL" {
		destinationRuleToParse = DestinationRuleConfigMutual
	} else if destinationRuleMode == "SIMPLE" {
		destinationRuleToParse = DestinationRuleConfigSimple
	} else {
		destinationRuleToParse = DestinationRuleConfigDisabledOrIstioMutual
	}
	if fakeRootCert {
		rootCertPathToUse = "/etc/certs/custom/fake-root-cert.pem"
	} else {
		rootCertPathToUse = "/etc/certs/custom/root-cert.pem"
	}
	tmpl, err := template.New("DestinationRule").Parse(destinationRuleToParse)
	if err != nil {
		t.Errorf("failed to create template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, map[string]string{
		"AppNamespace": serviceNamespace.Name(),
		"Mode":         destinationRuleMode, "RootCertPath": rootCertPathToUse,
	}); err != nil {
		t.Fatalf("failed to create template: %v", err)
	}
	return buf
}

// setupEcho creates two namespaces app and service. It also brings up two echo instances server and
// client in app namespace. HTTP and HTTPS port on the server echo are set up. Sidecar scope config
// is applied to only allow egress traffic to service namespace such that when client to server calls are made
// we are able to simulate "external" traffic by going outside this namespace. Egress Gateway is set up in the
// service namespace to handle egress for "external" calls.
func setupEcho(t framework.TestContext, ctx resource.Context) (echo.Instance, echo.Instance, namespace.Instance, namespace.Instance) {
	appsNamespace := namespace.NewOrFail(t, ctx, namespace.Config{
		Prefix: "app",
		Inject: true,
	})
	serviceNamespace := namespace.NewOrFail(t, ctx, namespace.Config{
		Prefix: "service",
		Inject: true,
	})

	var internalClient, externalServer echo.Instance
	deployment.New(ctx).
		With(&internalClient, echo.Config{
			Service:   "client",
			Namespace: appsNamespace,
			Ports:     []echo.Port{},
			Subsets: []echo.SubsetConfig{{
				Version: "v1",
			}},
			Cluster: ctx.Clusters().Default(),
		}).
		With(&externalServer, echo.Config{
			Service:   "server",
			Namespace: serviceNamespace,
			Ports: []echo.Port{
				{
					// Plain HTTP port only used to route request to egress gateway
					Name:         "http",
					Protocol:     protocol.HTTP,
					ServicePort:  80,
					WorkloadPort: 8080,
				},
				{
					// HTTPS port
					Name:         "https",
					Protocol:     protocol.HTTPS,
					ServicePort:  443,
					WorkloadPort: 8443,
					TLS:          true,
				},
			},
			// Set up TLS certs on the server. This will make the server listen with these credentials.
			TLSSettings: &common.TLSSettings{
				// Echo has these test certs baked into the docker image
				RootCert:   mustReadCert(t, "root-cert.pem"),
				ClientCert: mustReadCert(t, "cert-chain.pem"),
				Key:        mustReadCert(t, "key.pem"),
				// Override hostname to match the SAN in the cert we are using
				Hostname: "server.default.svc",
			},
			Subsets: []echo.SubsetConfig{{
				Version:     "v1",
				Annotations: echo.NewAnnotations().SetBool(echo.SidecarInject, false),
			}},
			Cluster: ctx.Clusters().Default(),
		}).
		BuildOrFail(t)

	// Apply Egress Gateway for service namespace to originate external traffic
	createGateway(t, ctx, appsNamespace, serviceNamespace)

	if err := WaitUntilNotCallable(internalClient, externalServer); err != nil {
		t.Fatalf("failed to apply sidecar, %v", err)
	}

	return internalClient, externalServer, appsNamespace, serviceNamespace
}

const (
	Gateway = `
apiVersion: networking.istio.io/v1alpha3
kind: Gateway
metadata:
  name: istio-egressgateway-filebased
spec:
  selector:
    istio: egressgateway
  servers:
    - port:
        number: 443
        name: https-filebased
        protocol: HTTPS
      hosts:
        - server.{{.ServerNamespace}}.svc.cluster.local
      tls:
        mode: ISTIO_MUTUAL
---
apiVersion: networking.istio.io/v1alpha3
kind: DestinationRule
metadata:
  name: egressgateway-for-server-filebased
spec:
  host: istio-egressgateway.istio-system.svc.cluster.local
  subsets:
  - name: server
    trafficPolicy:
      portLevelSettings:
      - port:
          number: 443
        tls:
          mode: ISTIO_MUTUAL
          sni: server.{{.ServerNamespace}}.svc.cluster.local
`
	VirtualService = `
apiVersion: networking.istio.io/v1alpha3
kind: VirtualService
metadata:
  name: route-via-egressgateway-filebased
spec:
  hosts:
    - server.{{.ServerNamespace}}.svc.cluster.local
  gateways:
    - istio-egressgateway-filebased
    - mesh
  http:
    - match:
        - gateways:
            - mesh # from sidecars, route to egress gateway service
          port: 80
      route:
        - destination:
            host: istio-egressgateway.istio-system.svc.cluster.local
            subset: server
            port:
              number: 443
          weight: 100
    - match:
        - gateways:
            - istio-egressgateway-filebased
          port: 443
      route:
        - destination:
            host: server.{{.ServerNamespace}}.svc.cluster.local
            port:
              number: 443
          weight: 100
      headers:
        request:
          add:
            handled-by-egress-gateway: "true"
`
)

func createGateway(t test.Failer, ctx resource.Context, appsNamespace namespace.Instance, serviceNamespace namespace.Instance) {
	tmplGateway, err := template.New("Gateway").Parse(Gateway)
	if err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	var bufGateway bytes.Buffer
	if err := tmplGateway.Execute(&bufGateway, map[string]string{"ServerNamespace": serviceNamespace.Name()}); err != nil {
		t.Fatalf("failed to create template: %v", err)
	}
	if err := ctx.ConfigIstio().YAML(appsNamespace.Name(), bufGateway.String()).Apply(); err != nil {
		t.Fatalf("failed to apply gateway: %v. template: %v", err, bufGateway.String())
	}

	// Have to wait for DR to apply to all sidecars first!
	time.Sleep(5 * time.Second)

	tmplVS, err := template.New("Gateway").Parse(VirtualService)
	if err != nil {
		t.Fatalf("failed to create template: %v", err)
	}

	var bufVS bytes.Buffer

	if err := tmplVS.Execute(&bufVS, map[string]string{"ServerNamespace": serviceNamespace.Name()}); err != nil {
		t.Fatalf("failed to create template: %v", err)
	}
	if err := ctx.ConfigIstio().YAML(appsNamespace.Name(), bufVS.String()).Apply(); err != nil {
		t.Fatalf("failed to apply virtualservice: %v. template: %v", err, bufVS.String())
	}
}

func clusterName(target echo.Instance, port echo.Port) string {
	cfg := target.Config()
	return fmt.Sprintf("outbound|%d||%s.%s.svc.%s", port.ServicePort, cfg.Service, cfg.Namespace.Name(), cfg.Domain)
}

// Wait for the server to NOT be callable by the client. This allows us to simulate external traffic.
// This essentially just waits for the Sidecar to be applied, without sleeping.
func WaitUntilNotCallable(c echo.Instance, dest echo.Instance) error {
	accept := func(cfg *envoyAdmin.ConfigDump) (bool, error) {
		validator := structpath.ForProto(cfg)
		for _, port := range dest.Config().Ports {
			clusterName := clusterName(dest, port)
			// Ensure that we have an outbound configuration for the target port.
			err := validator.NotExists("{.configs[*].dynamicActiveClusters[?(@.cluster.Name == '%s')]}", clusterName).Check()
			if err != nil {
				return false, err
			}
		}

		return true, nil
	}

	workloads, _ := c.Workloads()
	// Wait for the outbound config to be received by each workload from Pilot.
	for _, w := range workloads {
		if w.Sidecar() != nil {
			if err := w.Sidecar().WaitForConfig(accept, retry.Timeout(time.Second*10)); err != nil {
				return err
			}
		}
	}

	return nil
}
