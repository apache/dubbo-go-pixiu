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

package model_test

import (
	"fmt"
	"net"
	"reflect"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/miekg/dns"
	meshconfig "istio.io/api/mesh/v1alpha1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/serviceregistry/memory"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/xds"
	"github.com/apache/dubbo-go-pixiu/pkg/config/mesh"
	"github.com/apache/dubbo-go-pixiu/pkg/test/scopes"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/retry"
	"github.com/apache/dubbo-go-pixiu/pkg/util/sets"
)

func TestGatewayHostnames(t *testing.T) {
	origMinGatewayTTL := model.MinGatewayTTL
	model.MinGatewayTTL = 30 * time.Millisecond
	t.Cleanup(func() {
		model.MinGatewayTTL = origMinGatewayTTL
	})

	gwHost := "test.gw.istio.io"
	dnsServer := newFakeDNSServer(":0", 1, sets.New(gwHost))
	model.NetworkGatewayTestDNSServers = []string{dnsServer.Server.PacketConn.LocalAddr().String()}
	t.Cleanup(func() {
		if err := dnsServer.Shutdown(); err != nil {
			t.Logf("failed shutting down fake dns server")
		}
	})

	meshNetworks := mesh.NewFixedNetworksWatcher(nil)
	xdsUpdater := &xds.FakeXdsUpdater{Events: make(chan xds.FakeXdsEvent, 10)}
	env := &model.Environment{NetworksWatcher: meshNetworks, ServiceDiscovery: memory.NewServiceDiscovery()}
	if err := env.InitNetworksManager(xdsUpdater); err != nil {
		t.Fatal(err)
	}

	var firstIP string
	t.Run("initial resolution", func(t *testing.T) {
		meshNetworks.SetNetworks(&meshconfig.MeshNetworks{Networks: map[string]*meshconfig.Network{
			"nw0": {Gateways: []*meshconfig.Network_IstioNetworkGateway{{
				Gw: &meshconfig.Network_IstioNetworkGateway_Address{
					Address: gwHost,
				},
				Port: 15443,
			}}},
		}})
		xdsUpdater.WaitDurationOrFail(t, model.MinGatewayTTL+5*time.Second, "xds")
		gws := env.NetworkManager.AllGateways()
		if len(gws) != 1 {
			t.Fatalf("expected 1 IP")
		}
		if gws[0].Network != "nw0" {
			t.Fatalf("unexpected network: %v", gws)
		}
		firstIP = gws[0].Addr
	})
	t.Run("re-resolve after TTL", func(t *testing.T) {
		if testing.Short() {
			t.Skip()
		}
		// after the update, we should see the next gateway. Since TTL is low we don't know the exact IP, but we know it should change from
		// the original
		retry.UntilOrFail(t, func() bool {
			return !reflect.DeepEqual(env.NetworkManager.AllGateways(), []model.NetworkGateway{{Network: "nw0", Addr: firstIP, Port: 15443}})
		})
	})
	t.Run("forget", func(t *testing.T) {
		meshNetworks.SetNetworks(nil)
		xdsUpdater.WaitDurationOrFail(t, 5*time.Second, "xds")
		if len(env.NetworkManager.AllGateways()) > 0 {
			t.Fatalf("expected no gateways")
		}
	})
}

type fakeDNSServer struct {
	*dns.Server
	ttl uint32

	mu sync.Mutex
	// map fqdn hostname -> query count
	hosts map[string]int
}

func newFakeDNSServer(addr string, ttl uint32, hosts sets.Set) *fakeDNSServer {
	waitLock := sync.Mutex{}
	waitLock.Lock()
	s := &fakeDNSServer{
		Server: &dns.Server{Addr: addr, Net: "udp", NotifyStartedFunc: waitLock.Unlock},
		ttl:    ttl,
		hosts:  make(map[string]int, len(hosts)),
	}
	s.Handler = s

	for k := range hosts {
		s.hosts[dns.Fqdn(k)] = 0
	}

	go func() {
		if err := s.ListenAndServe(); err != nil {
			scopes.Framework.Errorf("fake dns server error: %v", err)
		}
	}()
	waitLock.Lock()
	return s
}

func (s *fakeDNSServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	s.mu.Lock()
	defer s.mu.Unlock()

	msg := (&dns.Msg{}).SetReply(r)
	switch r.Question[0].Qtype {
	case dns.TypeA, dns.TypeANY:
		domain := msg.Question[0].Name
		c, ok := s.hosts[domain]
		if ok {
			s.hosts[domain]++
			msg.Answer = append(msg.Answer, &dns.A{
				Hdr: dns.RR_Header{Name: domain, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: s.ttl},
				A:   net.ParseIP(fmt.Sprintf("10.0.0.%d", c)),
			})
		}
	}
	if err := w.WriteMsg(msg); err != nil {
		scopes.Framework.Errorf("failed writing fake DNS response: %v", err)
	}
}
