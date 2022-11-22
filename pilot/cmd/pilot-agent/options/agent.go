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

package options

import (
	"os"
	"path/filepath"
	"strings"
)

import (
	meshconfig "istio.io/api/mesh/v1alpha1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/bootstrap/platform"
	istioagent "github.com/apache/dubbo-go-pixiu/pkg/istio-agent"
)

// Similar with ISTIO_META_, which is used to customize the node metadata - this customizes extra header.
const xdsHeaderPrefix = "XDS_HEADER_"

func NewAgentOptions(proxy *model.Proxy, cfg *meshconfig.ProxyConfig) *istioagent.AgentOptions {
	o := &istioagent.AgentOptions{
		XDSRootCerts:                xdsRootCA,
		CARootCerts:                 caRootCA,
		XDSHeaders:                  map[string]string{},
		XdsUdsPath:                  filepath.Join(cfg.ConfigPath, "XDS"),
		IsIPv6:                      proxy.IsIPv6(),
		ProxyType:                   proxy.Type,
		EnableDynamicProxyConfig:    enableProxyConfigXdsEnv,
		EnableDynamicBootstrap:      enableBootstrapXdsEnv,
		WASMInsecureRegistries:      strings.Split(wasmInsecureRegistries, ","),
		ProxyIPAddresses:            proxy.IPAddresses,
		ServiceNode:                 proxy.ServiceNode(),
		EnvoyStatusPort:             envoyStatusPortEnv,
		EnvoyPrometheusPort:         envoyPrometheusPortEnv,
		MinimumDrainDuration:        minimumDrainDurationEnv,
		ExitOnZeroActiveConnections: exitOnZeroActiveConnectionsEnv,
		Platform:                    platform.Discover(proxy.SupportsIPv6()),
		GRPCBootstrapPath:           grpcBootstrapEnv,
		DisableEnvoy:                disableEnvoyEnv,
		ProxyXDSDebugViaAgent:       proxyXDSDebugViaAgent,
		ProxyXDSDebugViaAgentPort:   proxyXDSDebugViaAgentPort,
		DNSCapture:                  DNSCaptureByAgent.Get(),
		DNSAddr:                     DNSCaptureAddr.Get(),
		ProxyNamespace:              PodNamespaceVar.Get(),
		ProxyDomain:                 proxy.DNSDomain,
		IstiodSAN:                   istiodSAN.Get(),
	}
	extractXDSHeadersFromEnv(o)
	return o
}

// Simplified extraction of gRPC headers from environment.
// Unlike ISTIO_META, where we need JSON and advanced features - this is just for small string headers.
func extractXDSHeadersFromEnv(o *istioagent.AgentOptions) {
	envs := os.Environ()
	for _, e := range envs {
		if strings.HasPrefix(e, xdsHeaderPrefix) {
			parts := strings.SplitN(e, "=", 2)
			if len(parts) != 2 {
				continue
			}
			o.XDSHeaders[parts[0][len(xdsHeaderPrefix):]] = parts[1]
		}
	}
}
