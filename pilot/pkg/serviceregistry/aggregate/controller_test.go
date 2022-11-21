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

package aggregate

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

import (
	"github.com/google/go-cmp/cmp"
	"go.uber.org/atomic"
	meshconfig "istio.io/api/mesh/v1alpha1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/serviceregistry"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/serviceregistry/memory"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/serviceregistry/mock"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/serviceregistry/provider"
	"github.com/apache/dubbo-go-pixiu/pkg/cluster"
	"github.com/apache/dubbo-go-pixiu/pkg/config/host"
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/retry"
)

type mockMeshConfigHolder struct {
	trustDomainAliases []string
}

func (mh mockMeshConfigHolder) Mesh() *meshconfig.MeshConfig {
	return &meshconfig.MeshConfig{
		TrustDomainAliases: mh.trustDomainAliases,
	}
}

func buildMockController() *Controller {
	discovery1 := memory.NewServiceDiscovery(mock.ReplicatedFooServiceV1.DeepCopy(),
		mock.HelloService.DeepCopy(),
		mock.ExtHTTPService.DeepCopy(),
	)
	for _, port := range mock.HelloService.Ports {
		discovery1.AddInstance(mock.HelloService.Hostname, mock.MakeServiceInstance(mock.HelloService, port, 0, model.Locality{}))
		discovery1.AddInstance(mock.HelloService.Hostname, mock.MakeServiceInstance(mock.HelloService, port, 1, model.Locality{}))
	}

	discovery2 := memory.NewServiceDiscovery(mock.ReplicatedFooServiceV2.DeepCopy(),
		mock.WorldService.DeepCopy(),
		mock.ExtHTTPSService.DeepCopy(),
	)
	for _, port := range mock.WorldService.Ports {
		discovery2.AddInstance(mock.WorldService.Hostname, mock.MakeServiceInstance(mock.WorldService, port, 0, model.Locality{}))
		discovery2.AddInstance(mock.WorldService.Hostname, mock.MakeServiceInstance(mock.WorldService, port, 1, model.Locality{}))
	}
	registry1 := serviceregistry.Simple{
		ProviderID:       provider.ID("mockAdapter1"),
		ServiceDiscovery: discovery1,
		Controller:       &mock.Controller{},
	}

	registry2 := serviceregistry.Simple{
		ProviderID:       provider.ID("mockAdapter2"),
		ServiceDiscovery: discovery2,
		Controller:       &mock.Controller{},
	}

	ctls := NewController(Options{&mockMeshConfigHolder{}})
	ctls.AddRegistry(registry1)
	ctls.AddRegistry(registry2)

	return ctls
}

// return aggregator and cluster1 and cluster2 service discovery
func buildMockControllerForMultiCluster() (*Controller, *memory.ServiceDiscovery, *memory.ServiceDiscovery) {
	discovery1 := memory.NewServiceDiscovery(mock.HelloService)

	discovery2 := memory.NewServiceDiscovery(mock.MakeService(mock.ServiceArgs{
		Hostname:        mock.HelloService.Hostname,
		Address:         "10.1.2.0",
		ServiceAccounts: []string{},
		ClusterID:       "cluster-2",
	}), mock.WorldService)

	registry1 := serviceregistry.Simple{
		ProviderID:       provider.Kubernetes,
		ClusterID:        "cluster-1",
		ServiceDiscovery: discovery1,
		Controller:       &mock.Controller{},
	}

	registry2 := serviceregistry.Simple{
		ProviderID:       provider.Kubernetes,
		ClusterID:        "cluster-2",
		ServiceDiscovery: discovery2,
		Controller:       &mock.Controller{},
	}

	ctls := NewController(Options{})
	ctls.AddRegistry(registry1)
	ctls.AddRegistry(registry2)

	return ctls, discovery1, discovery2
}

func TestServicesForMultiCluster(t *testing.T) {
	originalHelloService := mock.HelloService.DeepCopy()
	aggregateCtl, _, registry2 := buildMockControllerForMultiCluster()
	// List Services from aggregate controller
	services := aggregateCtl.Services()

	// Set up ground truth hostname values
	hosts := map[host.Name]bool{
		mock.HelloService.Hostname: false,
		mock.WorldService.Hostname: false,
	}

	count := 0
	// Compare return value to ground truth
	for _, svc := range services {
		if counted, existed := hosts[svc.Hostname]; existed && !counted {
			count++
			hosts[svc.Hostname] = true
		}
	}

	if count != len(hosts) {
		t.Fatalf("Cluster local service map expected size %d, actual %v vs %v", count, hosts, services)
	}

	// Now verify ClusterVIPs for each service
	ClusterVIPs := map[host.Name]map[cluster.ID][]string{
		mock.HelloService.Hostname: {
			"cluster-1": []string{"10.1.0.0"},
			"cluster-2": []string{"10.1.2.0"},
		},
		mock.WorldService.Hostname: {
			"cluster-2": []string{"10.2.0.0"},
		},
	}
	for _, svc := range services {
		if !reflect.DeepEqual(svc.ClusterVIPs.Addresses, ClusterVIPs[svc.Hostname]) {
			t.Fatalf("Service %s ClusterVIPs actual %v, expected %v", svc.Hostname,
				svc.ClusterVIPs.Addresses, ClusterVIPs[svc.Hostname])
		}
	}

	registry2.RemoveService(mock.HelloService.Hostname)
	// List Services from aggregate controller
	services = aggregateCtl.Services()
	// Now verify ClusterVIPs for each service
	ClusterVIPs = map[host.Name]map[cluster.ID][]string{
		mock.HelloService.Hostname: {
			"cluster-1": []string{"10.1.0.0"},
		},
		mock.WorldService.Hostname: {
			"cluster-2": []string{"10.2.0.0"},
		},
	}
	for _, svc := range services {
		if !reflect.DeepEqual(svc.ClusterVIPs.Addresses, ClusterVIPs[svc.Hostname]) {
			t.Fatalf("Service %s ClusterVIPs actual %v, expected %v", svc.Hostname,
				svc.ClusterVIPs.Addresses, ClusterVIPs[svc.Hostname])
		}
	}

	// check HelloService is not mutated
	if !reflect.DeepEqual(originalHelloService, mock.HelloService) {
		t.Errorf("Original hello service is mutated")
	}
}

func TestServices(t *testing.T) {
	aggregateCtl := buildMockController()
	// List Services from aggregate controller
	services := aggregateCtl.Services()

	// Set up ground truth hostname values
	serviceMap := map[host.Name]bool{
		mock.HelloService.Hostname:    false,
		mock.ExtHTTPService.Hostname:  false,
		mock.WorldService.Hostname:    false,
		mock.ExtHTTPSService.Hostname: false,
	}

	svcCount := 0
	// Compare return value to ground truth
	for _, svc := range services {
		if counted, existed := serviceMap[svc.Hostname]; existed && !counted {
			svcCount++
			serviceMap[svc.Hostname] = true
		}
	}

	if svcCount != len(serviceMap) {
		t.Fatal("Return services does not match ground truth")
	}
}

func TestGetService(t *testing.T) {
	aggregateCtl := buildMockController()

	// Get service from mockAdapter1
	svc := aggregateCtl.GetService(mock.HelloService.Hostname)
	if svc == nil {
		t.Fatal("Fail to get service")
	}
	if svc.Hostname != mock.HelloService.Hostname {
		t.Fatal("Returned service is incorrect")
	}

	// Get service from mockAdapter2
	svc = aggregateCtl.GetService(mock.WorldService.Hostname)
	if svc == nil {
		t.Fatal("Fail to get service")
	}
	if svc.Hostname != mock.WorldService.Hostname {
		t.Fatal("Returned service is incorrect")
	}
}

func TestGetProxyServiceInstances(t *testing.T) {
	aggregateCtl := buildMockController()

	// Get Instances from mockAdapter1
	instances := aggregateCtl.GetProxyServiceInstances(&model.Proxy{IPAddresses: []string{mock.HelloInstanceV0}})
	if len(instances) != 6 {
		t.Fatalf("Returned GetProxyServiceInstances' amount %d is not correct", len(instances))
	}
	for _, inst := range instances {
		if inst.Service.Hostname != mock.HelloService.Hostname {
			t.Fatal("Returned Instance is incorrect")
		}
	}

	// Get Instances from mockAdapter2
	instances = aggregateCtl.GetProxyServiceInstances(&model.Proxy{IPAddresses: []string{mock.MakeIP(mock.WorldService, 1)}})
	if len(instances) != 6 {
		t.Fatalf("Returned GetProxyServiceInstances' amount %d is not correct", len(instances))
	}
	for _, inst := range instances {
		if inst.Service.Hostname != mock.WorldService.Hostname {
			t.Fatal("Returned Instance is incorrect")
		}
	}
}

func TestGetProxyWorkloadLabels(t *testing.T) {
	// If no registries return workload labels, we must return nil, rather than an empty list.
	// This ensures callers can distinguish between no labels, and labels not found.
	aggregateCtl := buildMockController()

	instances := aggregateCtl.GetProxyWorkloadLabels(&model.Proxy{IPAddresses: []string{mock.HelloInstanceV0}})
	if instances != nil {
		t.Fatalf("expected nil workload labels, got: %v", instances)
	}
}

func TestInstances(t *testing.T) {
	aggregateCtl := buildMockController()

	// Get Instances from mockAdapter1
	instances := aggregateCtl.InstancesByPort(mock.HelloService, 80, nil)
	if len(instances) != 2 {
		t.Fatal("Returned wrong number of instances from controller")
	}
	for _, instance := range instances {
		if instance.Service.Hostname != mock.HelloService.Hostname {
			t.Fatal("Returned instance's hostname does not match desired value")
		}
		if _, ok := instance.Service.Ports.Get(mock.PortHTTPName); !ok {
			t.Fatal("Returned instance does not contain desired port")
		}
	}

	// Get Instances from mockAdapter2
	instances = aggregateCtl.InstancesByPort(mock.WorldService, 80, nil)
	if len(instances) != 2 {
		t.Fatal("Returned wrong number of instances from controller")
	}
	for _, instance := range instances {
		if instance.Service.Hostname != mock.WorldService.Hostname {
			t.Fatal("Returned instance's hostname does not match desired value")
		}
		if _, ok := instance.Service.Ports.Get(mock.PortHTTPName); !ok {
			t.Fatal("Returned instance does not contain desired port")
		}
	}
}

func TestGetIstioServiceAccounts(t *testing.T) {
	aggregateCtl := buildMockController()
	testCases := []struct {
		name               string
		svc                *model.Service
		trustDomainAliases []string
		want               []string
	}{
		{
			name: "HelloEmpty",
			svc:  mock.HelloService,
			want: []string{},
		},
		{
			name: "World",
			svc:  mock.WorldService,
			want: []string{
				"spiffe://cluster.local/ns/default/sa/world1",
				"spiffe://cluster.local/ns/default/sa/world2",
			},
		},
		{
			name: "ReplicatedFoo",
			svc:  mock.ReplicatedFooServiceV1,
			want: []string{
				"spiffe://cluster.local/ns/default/sa/foo-share",
				"spiffe://cluster.local/ns/default/sa/foo1",
				"spiffe://cluster.local/ns/default/sa/foo2",
			},
		},
		{
			name:               "ExpansionByTrustDomainAliases",
			trustDomainAliases: []string{"cluster.local", "example.com"},
			svc:                mock.WorldService,
			want: []string{
				"spiffe://cluster.local/ns/default/sa/world1",
				"spiffe://cluster.local/ns/default/sa/world2",
				"spiffe://example.com/ns/default/sa/world1",
				"spiffe://example.com/ns/default/sa/world2",
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			aggregateCtl.meshHolder = &mockMeshConfigHolder{trustDomainAliases: tc.trustDomainAliases}
			accounts := aggregateCtl.GetIstioServiceAccounts(tc.svc, []int{})
			if diff := cmp.Diff(accounts, tc.want); diff != "" {
				t.Errorf("unexpected service account, diff %v, %v", diff, accounts)
			}
		})
	}
}

func TestAddRegistry(t *testing.T) {
	registries := []serviceregistry.Simple{
		{
			ProviderID:       "registry1",
			ClusterID:        "cluster1",
			Controller:       &mock.Controller{},
			ServiceDiscovery: memory.NewServiceDiscovery(),
		},
		{
			ProviderID:       "registry2",
			ClusterID:        "cluster2",
			Controller:       &mock.Controller{},
			ServiceDiscovery: memory.NewServiceDiscovery(),
		},
	}
	ctrl := NewController(Options{})

	registry1Counter := atomic.NewInt32(0)
	registry2Counter := atomic.NewInt32(0)

	for _, r := range registries {
		clusterID := r.Cluster()
		counter := registry1Counter
		if clusterID == "cluster2" {
			counter = registry2Counter
		}
		ctrl.AppendServiceHandlerForCluster(clusterID, func(service *model.Service, event model.Event) {
			counter.Add(1)
		})
		ctrl.AddRegistry(r)
	}
	if l := len(ctrl.registries); l != 2 {
		t.Fatalf("Expected length of the registries slice should be 2, got %d", l)
	}

	registries[0].Controller.(*mock.Controller).OnServiceEvent(mock.HelloService, model.EventAdd)
	registries[1].Controller.(*mock.Controller).OnServiceEvent(mock.WorldService, model.EventAdd)

	ctrl.DeleteRegistry(registries[1].Cluster(), registries[1].Provider())
	ctrl.UnRegisterHandlersForCluster(registries[1].Cluster())
	registries[0].Controller.(*mock.Controller).OnServiceEvent(mock.HelloService, model.EventAdd)

	if registry1Counter.Load() != 3 {
		t.Errorf("cluster1 expected 3 event, but got %d", registry1Counter.Load())
	}
	if registry2Counter.Load() != 2 {
		t.Errorf("cluster2 expected 2 event, but got %d", registry2Counter.Load())
	}
}

func TestGetDeleteRegistry(t *testing.T) {
	registries := []serviceregistry.Simple{
		{
			ProviderID:       "registry1",
			ClusterID:        "cluster1",
			Controller:       &mock.Controller{},
			ServiceDiscovery: memory.NewServiceDiscovery(),
		},
		{
			ProviderID:       "registry2",
			ClusterID:        "cluster2",
			Controller:       &mock.Controller{},
			ServiceDiscovery: memory.NewServiceDiscovery(),
		},
		{
			ProviderID:       "registry3",
			ClusterID:        "cluster3",
			Controller:       &mock.Controller{},
			ServiceDiscovery: memory.NewServiceDiscovery(),
		},
	}
	wrapRegistry := func(r serviceregistry.Instance) serviceregistry.Instance {
		return &registryEntry{Instance: r}
	}

	ctrl := NewController(Options{})
	for _, r := range registries {
		ctrl.AddRegistry(r)
	}

	// Test Get
	result := ctrl.GetRegistries()
	if l := len(result); l != 3 {
		t.Fatalf("Expected length of the registries slice should be 3, got %d", l)
	}

	// Test Delete cluster2
	ctrl.DeleteRegistry(registries[1].ClusterID, registries[1].ProviderID)
	result = ctrl.GetRegistries()
	if l := len(result); l != 2 {
		t.Fatalf("Expected length of the registries slice should be 2, got %d", l)
	}
	// check left registries are orders as before
	if !reflect.DeepEqual(result[0], wrapRegistry(registries[0])) || !reflect.DeepEqual(result[1], wrapRegistry(registries[2])) {
		t.Fatalf("Expected registries order has been changed")
	}
}

func TestSkipSearchingRegistryForProxy(t *testing.T) {
	cluster1 := serviceregistry.Simple{
		ClusterID:        "cluster-1",
		ProviderID:       provider.Kubernetes,
		Controller:       &mock.Controller{},
		ServiceDiscovery: memory.NewServiceDiscovery(),
	}
	cluster2 := serviceregistry.Simple{
		ClusterID:        "cluster-2",
		ProviderID:       provider.Kubernetes,
		Controller:       &mock.Controller{},
		ServiceDiscovery: memory.NewServiceDiscovery(),
	}
	// external registries may eventually be associated with a cluster
	external := serviceregistry.Simple{
		ClusterID:        "cluster-1",
		ProviderID:       provider.External,
		Controller:       &mock.Controller{},
		ServiceDiscovery: memory.NewServiceDiscovery(),
	}

	cases := []struct {
		nodeClusterID cluster.ID
		registry      serviceregistry.Instance
		want          bool
	}{
		// matching kube registry
		{"cluster-1", cluster1, false},
		// unmatching kube registry
		{"cluster-1", cluster2, true},
		// always search external
		{"cluster-1", external, false},
		{"cluster-2", external, false},
		{"", external, false},
		// always search for empty node cluster id
		{"", cluster1, false},
		{"", cluster2, false},
		{"", external, false},
	}

	for i, c := range cases {
		got := skipSearchingRegistryForProxy(c.nodeClusterID, c.registry)
		if got != c.want {
			t.Errorf("%s: got %v want %v",
				fmt.Sprintf("[%v] registry=%v node=%v", i, c.registry, c.nodeClusterID),
				got, c.want)
		}
	}
}

func runnableRegistry(name string) *RunnableRegistry {
	return &RunnableRegistry{
		Instance: serviceregistry.Simple{
			ClusterID: cluster.ID(name), ProviderID: "test",
			Controller:       &mock.Controller{},
			ServiceDiscovery: memory.NewServiceDiscovery(),
		},
		running: atomic.NewBool(false),
	}
}

type RunnableRegistry struct {
	serviceregistry.Instance
	running *atomic.Bool
}

func (rr *RunnableRegistry) Run(stop <-chan struct{}) {
	if rr.running.Load() {
		panic("--- registry has been run twice ---")
	}
	rr.running.Store(true)
	<-stop
}

func expectRunningOrFail(t *testing.T, ctrl *Controller, want bool) {
	// running gets flipped in a goroutine, retry to avoid race
	retry.UntilSuccessOrFail(t, func() error {
		for _, registry := range ctrl.registries {
			if running := registry.Instance.(*RunnableRegistry).running.Load(); running != want {
				return fmt.Errorf("%s running is %v but wanted %v", registry.Cluster(), running, want)
			}
		}
		return nil
	}, retry.Timeout(50*time.Millisecond), retry.Delay(0))
}

func TestDeferredRun(t *testing.T) {
	stop := make(chan struct{})
	defer close(stop)
	ctrl := NewController(Options{})

	t.Run("AddRegistry before aggregate Run does not run", func(t *testing.T) {
		ctrl.AddRegistry(runnableRegistry("earlyAdd"))
		ctrl.AddRegistryAndRun(runnableRegistry("earlyAddAndRun"), nil)
		expectRunningOrFail(t, ctrl, false)
	})
	t.Run("aggregate Run starts all registries", func(t *testing.T) {
		go ctrl.Run(stop)
		expectRunningOrFail(t, ctrl, true)
		ctrl.DeleteRegistry("earlyAdd", "test")
		ctrl.DeleteRegistry("earlyAddAndRun", "test")
	})
	t.Run("AddRegistry after aggregate Run does not start registry", func(t *testing.T) {
		ctrl.AddRegistry(runnableRegistry("missed"))
		expectRunningOrFail(t, ctrl, false)
		ctrl.DeleteRegistry("missed", "test")
		expectRunningOrFail(t, ctrl, true)
	})
	t.Run("AddRegistryAndRun after aggregate Run starts registry", func(t *testing.T) {
		ctrl.AddRegistryAndRun(runnableRegistry("late"), nil)
		expectRunningOrFail(t, ctrl, true)
	})
}
