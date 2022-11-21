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

package authn

import (
	http_conn "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/network/http_connection_manager/v3"
	"istio.io/api/security/v1beta1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/networking/plugin"
)

// PolicyApplier is the interface provides essential functionalities to help config Envoy (xDS) to enforce
// authentication policy. Each version of authentication policy will implement this interface.
type PolicyApplier interface {
	// InboundMTLSSettings returns inbound mTLS settings for a given workload port
	InboundMTLSSettings(endpointPort uint32, node *model.Proxy, trustDomainAliases []string) plugin.MTLSSettings

	// JwtFilter returns the JWT HTTP filter to enforce the underlying authentication policy.
	// It may return nil, if no JWT validation is needed.
	JwtFilter() *http_conn.HttpFilter

	// AuthNFilter returns the (authn) HTTP filter to enforce the underlying authentication policy.
	// It may return nil, if no authentication is needed.
	AuthNFilter(forSidecar bool) *http_conn.HttpFilter

	// PortLevelSetting returns port level mTLS settings.
	PortLevelSetting() map[uint32]*v1beta1.PeerAuthentication_MutualTLS

	// GetMutualTLSModeForPort gets the mTLS mode for the given port. If there is no port level setting, it
	// returns the inherited namespace/mesh level setting.
	GetMutualTLSModeForPort(endpointPort uint32) model.MutualTLSMode
}
