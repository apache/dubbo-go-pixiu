// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package capture

import (
	"fmt"
	"net"
	"strconv"
)

import (
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
)

import (
	"github.com/apache/dubbo-go-pixiu/tools/istio-iptables/pkg/config"
	"github.com/apache/dubbo-go-pixiu/tools/istio-iptables/pkg/constants"
)

// configureTProxyRoutes configures ip firewall rules to enable TPROXY support.
// See https://www.envoyproxy.io/docs/envoy/latest/configuration/http/http_filters/original_src_filter
func configureTProxyRoutes(cfg *config.Config) error {
	if cfg.InboundPortsInclude != "" {
		if cfg.InboundInterceptionMode == constants.TPROXY {
			link, err := netlink.LinkByName("lo")
			if err != nil {
				return fmt.Errorf("failed to find 'lo' link: %v", err)
			}
			tproxyTable, err := strconv.Atoi(cfg.InboundTProxyRouteTable)
			if err != nil {
				return fmt.Errorf("failed to parse InboundTProxyRouteTable: %v", err)
			}
			tproxyMark, err := strconv.Atoi(cfg.InboundTProxyMark)
			if err != nil {
				return fmt.Errorf("failed to parse InboundTProxyMark: %v", err)
			}
			// Route all packets marked in chain ISTIODIVERT using routing table ${INBOUND_TPROXY_ROUTE_TABLE}.
			// Equivalent to `ip rule add fwmark <tproxyMark> lookup <tproxyTable>`
			families := []int{unix.AF_INET}
			if cfg.EnableInboundIPv6 {
				families = append(families, unix.AF_INET6)
			}
			for _, family := range families {
				r := netlink.NewRule()
				r.Family = family
				r.Table = tproxyTable
				r.Mark = tproxyMark
				if err := netlink.RuleAdd(r); err != nil {
					return fmt.Errorf("failed to configure netlink rule: %v", err)
				}
			}
			// In routing table ${INBOUND_TPROXY_ROUTE_TABLE}, create a single default rule to route all traffic to
			// the loopback interface.
			// Equivalent to `ip route add local default dev lo table <table>`
			cidrs := []string{"0.0.0.0/0"}
			if cfg.EnableInboundIPv6 {
				cidrs = append(cidrs, "0::0/0")
			}
			for _, fullCIDR := range cidrs {
				_, dst, err := net.ParseCIDR(fullCIDR)
				if err != nil {
					return fmt.Errorf("parse CIDR: %v", err)
				}

				if err := netlink.RouteAdd(&netlink.Route{
					Dst:       dst,
					Scope:     netlink.SCOPE_HOST,
					Type:      unix.RTN_LOCAL,
					Table:     tproxyTable,
					LinkIndex: link.Attrs().Index,
				}); ignoreExists(err) != nil {
					return fmt.Errorf("failed to add route: %v", err)
				}
			}
		}
	}
	return nil
}
