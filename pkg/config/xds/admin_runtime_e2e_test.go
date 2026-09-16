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

package xds_test

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	_ "dubbo.apache.org/dubbo-go/v3/imports"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	rpcserver "dubbo.apache.org/dubbo-go/v3/server"

	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	envoyserver "github.com/envoyproxy/go-control-plane/pkg/server/v3"

	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"

	"github.com/stretchr/testify/require"

	clientv3 "go.etcd.io/etcd/client/v3"

	"go.etcd.io/etcd/server/v3/embed"

	"google.golang.org/grpc"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/controller/auth"
	"github.com/apache/dubbo-go-pixiu/admin/initialize"
	adminxds "github.com/apache/dubbo-go-pixiu/admin/xds"
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/rand"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	pixiuxds "github.com/apache/dubbo-go-pixiu/pkg/config/xds"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/http/apiconfig"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/http/remote"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/network/httpconnectionmanager"
	_ "github.com/apache/dubbo-go-pixiu/pkg/listener/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	runtimeE2EHTTPFilter         = "dgp.test.xds.runtime.proxy"
	runtimeE2EProviderProcessEnv = "PIXIU_RUNTIME_E2E_DUBBO_PROVIDER"
	runtimeE2EProviderPortEnv    = "PIXIU_RUNTIME_E2E_DUBBO_PORT"
	runtimeE2EProviderMarkerEnv  = "PIXIU_RUNTIME_E2E_DUBBO_MARKER"
)

type runtimeE2EHTTPPlugin struct {
	clusters *server.ClusterManager
}

type runtimeE2EHTTPConfig struct {
	Marker string `yaml:"marker" json:"marker"`
}

type runtimeE2EHTTPFactory struct {
	config   *runtimeE2EHTTPConfig
	clusters *server.ClusterManager
}

type runtimeE2EHTTPProxy struct {
	marker   string
	clusters *server.ClusterManager
}

type RuntimeE2EDubboProvider struct {
	marker string
}

type runtimeE2EDynamicResources struct {
	config *model.ApiConfigSource
	node   *model.Node
}

type runtimeE2EXDSServer struct {
	t             *testing.T
	ctx           context.Context
	snapshotCache cache.SnapshotCache
	address       string
	nacks         chan string

	mu         sync.Mutex
	grpcServer *grpc.Server
	listener   net.Listener
}

func (r runtimeE2EDynamicResources) GetLds() *model.ApiConfigSource { return r.config }
func (r runtimeE2EDynamicResources) GetCds() *model.ApiConfigSource { return r.config }
func (r runtimeE2EDynamicResources) GetNode() *model.Node           { return r.node }

func (p *runtimeE2EHTTPPlugin) Kind() string { return runtimeE2EHTTPFilter }

func (p *runtimeE2EHTTPPlugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &runtimeE2EHTTPFactory{config: &runtimeE2EHTTPConfig{}, clusters: p.clusters}, nil
}

func (f *runtimeE2EHTTPFactory) Config() any { return f.config }
func (f *runtimeE2EHTTPFactory) Apply() error {
	if strings.TrimSpace(f.config.Marker) == "" {
		return fmt.Errorf("marker is required")
	}
	return nil
}

func (f *runtimeE2EHTTPFactory) PrepareFilterChain(_ *contexthttp.HttpContext, chain filter.FilterChain) error {
	chain.AppendDecodeFilters(&runtimeE2EHTTPProxy{marker: f.config.Marker, clusters: f.clusters})
	return nil
}

func (f *runtimeE2EHTTPProxy) Decode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	route := ctx.GetRouteEntry()
	if route == nil {
		ctx.SendLocalReply(http.StatusNotFound, []byte("route not found"))
		return filter.Stop
	}
	endpoint := f.clusters.PickEndpoint(route.Cluster, ctx)
	if endpoint == nil {
		ctx.SendLocalReply(http.StatusServiceUnavailable, []byte("endpoint not found"))
		return filter.Stop
	}
	request, err := http.NewRequestWithContext(ctx.Ctx, ctx.Request.Method,
		"http://"+endpoint.Address.GetAddress()+ctx.Request.URL.RequestURI(), ctx.Request.Body)
	if err != nil {
		ctx.SendLocalReply(http.StatusInternalServerError, []byte(err.Error()))
		return filter.Stop
	}
	request.Header = ctx.Request.Header.Clone()
	request.Header.Set("X-XDS-E2E-Marker", f.marker)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		ctx.SendLocalReply(http.StatusBadGateway, []byte(err.Error()))
		return filter.Stop
	}
	ctx.SourceResp = response
	return filter.Continue
}

func (p *RuntimeE2EDubboProvider) Reference() string {
	return "pixiu.e2e.Greeter"
}

func (p *RuntimeE2EDubboProvider) SayHello(_ context.Context, name string) (string, error) {
	return p.marker + ":" + name, nil
}

var runtimeE2EDubboServiceInfo = rpcserver.ServiceInfo{
	InterfaceName: "pixiu.e2e.Greeter",
	ServiceType:   (*RuntimeE2EDubboProvider)(nil),
	Methods: []rpcserver.MethodInfo{{
		Name: "SayHello",
		Type: constant.CallUnary,
		ReqInitFunc: func() any {
			return new(string)
		},
		MethodFunc: func(ctx context.Context, args []any, handler any) (any, error) {
			return handler.(*RuntimeE2EDubboProvider).SayHello(ctx, args[0].(string))
		},
	}},
}

// TestRuntimeE2EDubboProviderProcess runs only as a child process of the
// runtime E2E test. Keeping dubbo-go's protocol globals in another process
// prevents this acceptance path from changing the state of sibling tests.
func TestRuntimeE2EDubboProviderProcess(t *testing.T) {
	if os.Getenv(runtimeE2EProviderProcessEnv) != "1" {
		return
	}

	port, err := strconv.Atoi(os.Getenv(runtimeE2EProviderPortEnv))
	require.NoError(t, err)
	// The protocol IP is intentionally left empty. dubbo-go always exports its
	// internal MetadataService on the same port with an empty IP, and
	// DubboProtocol.openServer caches listeners per url.Location. A concrete IP
	// would give this service "127.0.0.1:<port>", i.e. a second location whose
	// listen overlaps the wildcard MetadataService listener; Linux rejects that
	// with EADDRINUSE while BSD/macOS tolerates it. Sharing the wildcard
	// location makes both services reuse one listener.
	srv, err := rpcserver.NewServer(rpcserver.WithServerProtocol(
		protocol.WithDubbo(),
		protocol.WithPort(port),
	))
	require.NoError(t, err)
	require.NoError(t, srv.Register(
		&RuntimeE2EDubboProvider{marker: os.Getenv(runtimeE2EProviderMarkerEnv)},
		&runtimeE2EDubboServiceInfo,
		rpcserver.WithInterface(runtimeE2EDubboServiceInfo.InterfaceName),
		rpcserver.WithNotRegister(),
	))
	require.NoError(t, srv.Serve())
}

// TestAdminHTTPToRunningPixiuEndToEnd crosses the actual control-plane and
// runtime boundaries, including HTTP and Dubbo upstream calls.
func TestAdminHTTPToRunningPixiuEndToEnd(t *testing.T) {
	etcdEndpoint := runtimeE2EEtcdEndpoint(t)
	gin.SetMode(gin.TestMode)
	previousBootstrap := adminconfig.Bootstrap
	previousClient := adminconfig.Client
	previousStatus := adminxds.DefaultStatusStore
	t.Cleanup(func() {
		gin.SetMode(gin.DebugMode)
		adminconfig.CloseEtcdClient()
		adminconfig.Bootstrap = previousBootstrap
		adminconfig.Client = previousClient
		adminxds.DefaultStatusStore = previousStatus
	})

	etcdPath := fmt.Sprintf("/pixiu-e2e/%d", time.Now().UnixNano())
	adminconfig.Bootstrap = &adminconfig.AdminBootstrap{EtcdConfig: adminconfig.EtcdConfig{
		Address: etcdEndpoint,
		Path:    etcdPath,
	}}
	adminconfig.InitEtcdClient()
	require.NotNil(t, adminconfig.Client)
	t.Cleanup(func() {
		_, _ = adminconfig.Client.GetRawClient().Delete(
			adminconfig.Client.GetCtx(), etcdPath, clientv3.WithPrefix(),
		)
	})

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	nodeID := "pixiu-runtime-e2e"
	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	adminxds.DefaultStatusStore = adminxds.NewStatusStore(nodeID)
	publisher := adminxds.NewSnapshotPublisher(
		nodeID,
		adminxds.NewSnapshotBuilder(adminxds.LogicResourceLoader{}),
		snapshotCache,
		adminxds.DefaultStatusStore,
	)
	watchConfigForRuntimeE2E(t, ctx, publisher)

	xdsServer := newRuntimeE2EXDSServer(t, ctx, snapshotCache)
	t.Cleanup(xdsServer.Stop)
	managementHost, managementPort := splitAddress(t, xdsServer.address)
	bootstrap := &model.Bootstrap{StaticResources: model.StaticResources{
		Clusters: []*model.ClusterConfig{{
			Name:    "xds-management",
			TypeStr: "Static",
			Type:    model.Static,
			LbStr:   model.LoadBalancerRand,
			Endpoints: []*model.Endpoint{{
				ID: "xds-management-1",
				Address: model.SocketAddress{
					Address: managementHost,
					Port:    managementPort,
				},
			}},
		}},
		ShutdownConfig: &model.ShutdownConfig{Timeout: "0s"},
	}}
	clusterManager := server.CreateDefaultClusterManager(bootstrap)
	listenerManager := server.CreateDefaultListenerManager(bootstrap)
	filter.RegisterHttpFilter(&runtimeE2EHTTPPlugin{clusters: clusterManager})

	apiConfig := &model.ApiConfigSource{
		APIType:     model.ApiTypeGRPC,
		APITypeStr:  "GRPC",
		ClusterName: []string{"xds-management"},
	}
	node := &model.Node{Id: nodeID}
	xdsClient := pixiuxds.StartXdsClient(
		listenerManager,
		clusterManager,
		runtimeE2EDynamicResources{config: apiConfig, node: node},
	)
	require.NotNil(t, xdsClient)
	t.Cleanup(func() {
		xdsClient.Stop()
		listenerManager.RemoveXDSListeners(listenerManager.XDSListenerNames())
	})

	upstreamOne := httptest.NewServer(runtimeE2EUpstream("one"))
	t.Cleanup(upstreamOne.Close)
	upstreamTwo := httptest.NewServer(runtimeE2EUpstream("two"))
	t.Cleanup(upstreamTwo.Close)
	upstreamOneHost, upstreamOnePort := splitAddress(t, strings.TrimPrefix(upstreamOne.URL, "http://"))
	upstreamTwoHost, upstreamTwoPort := splitAddress(t, strings.TrimPrefix(upstreamTwo.URL, "http://"))
	listenerPort := reserveRuntimeE2EPort(t)

	router := initialize.Routers()
	token := runtimeE2EToken(t)
	clusterYAML := runtimeE2EClusterYAML(0, upstreamOneHost, upstreamOnePort)
	runtimeE2EAdminForm(t, router, token, http.MethodPut, "/config/api/cluster", clusterYAML)
	listenerYAML := runtimeE2EListenerYAML(listenerPort, "v1")
	runtimeE2EAdminForm(t, router, token, http.MethodPut, "/config/api/listener", listenerYAML)

	proxyURL := fmt.Sprintf("http://127.0.0.1:%d/orders", listenerPort)
	requireRuntimeE2EBody(t, proxyURL, "one-v1", 10*time.Second)

	runtimeE2EAdminForm(t, router, token, http.MethodPost, "/config/api/listener",
		runtimeE2EListenerYAML(listenerPort, "v2"))
	requireRuntimeE2EBody(t, proxyURL, "one-v2", 10*time.Second)

	// A syntactically valid Listener with an invalid nested plugin config must
	// be NACKed, and the v2 chain must continue serving throughout.
	runtimeE2EAdminForm(t, router, token, http.MethodPost, "/config/api/listener",
		runtimeE2EListenerYAML(listenerPort, ""))
	require.Eventually(t, func() bool {
		select {
		case nack := <-xdsServer.nacks:
			return strings.Contains(nack, "marker is required")
		default:
			return false
		}
	}, 10*time.Second, 25*time.Millisecond)
	require.Equal(t, "one-v2", runtimeE2EGet(proxyURL), "NACK replaced the last-good listener")

	runtimeE2EAdminForm(t, router, token, http.MethodPost, "/config/api/listener",
		runtimeE2EListenerYAML(listenerPort, "v3"))
	requireRuntimeE2EBody(t, proxyURL, "one-v3", 10*time.Second)

	// The data plane keeps the last-good chain while xDS is unavailable. A
	// config published during the outage is applied after the same endpoint
	// returns and the Delta client reconnects.
	xdsServer.Stop()
	runtimeE2EAdminForm(t, router, token, http.MethodPost, "/config/api/listener",
		runtimeE2EListenerYAML(listenerPort, "v4"))
	require.Equal(t, "one-v3", runtimeE2EGet(proxyURL))
	xdsServer.Start()
	requireRuntimeE2EBody(t, proxyURL, "one-v4", 15*time.Second)

	runtimeE2EAdminForm(t, router, token, http.MethodPost, "/config/api/cluster",
		runtimeE2EClusterYAML(1, upstreamTwoHost, upstreamTwoPort))
	requireRuntimeE2EBody(t, proxyURL, "two-v4", 10*time.Second)

	runtimeE2EAdminRequest(t, router, token, http.MethodDelete, "/config/api/listener?listener=gateway")
	runtimeE2EAdminRequest(t, router, token, http.MethodDelete, "/config/api/cluster?clusterId=1")
	requireRuntimeE2EPortClosed(t, listenerPort)

	// Exercise Pixiu's production HTTP-to-Dubbo path, rather than the test
	// HTTP forwarding filter above. The xDS Listener owns apiconfig and
	// dubboproxy filters; create/update/delete must affect real RPC traffic.
	providerOnePort := startRuntimeE2EDubboProvider(t, "rpc-one")
	providerTwoPort := startRuntimeE2EDubboProvider(t, "rpc-two")
	apiDir := t.TempDir()
	apiOne := filepath.Join(apiDir, "api-one.yaml")
	apiTwo := filepath.Join(apiDir, "api-two.yaml")
	writeRuntimeE2EDubboAPI(t, apiOne, providerOnePort)
	writeRuntimeE2EDubboAPI(t, apiTwo, providerTwoPort)

	rpcListenerPort := reserveRuntimeE2EPort(t)
	runtimeE2EAdminForm(t, router, token, http.MethodPut, "/config/api/cluster",
		runtimeE2EClusterYAML(2, "127.0.0.1", providerOnePort))
	runtimeE2EAdminForm(t, router, token, http.MethodPut, "/config/api/listener",
		runtimeE2EDubboListenerYAML(rpcListenerPort, apiOne))
	rpcURL := fmt.Sprintf("http://127.0.0.1:%d/rpc", rpcListenerPort)
	requireRuntimeE2EPostContains(t, rpcURL, `{"name":"alice"}`, `rpc-one:alice`, 15*time.Second)

	// The production dubboproxy resolves this direct API URL, not the Pixiu
	// ClusterManager. Prove that a cluster-only update does not masquerade as an
	// RPC routing update before changing the listener-owned API definition.
	runtimeE2EAdminForm(t, router, token, http.MethodPost, "/config/api/cluster",
		runtimeE2EClusterYAML(2, "127.0.0.1", providerTwoPort))
	require.Eventually(t, func() bool {
		endpoint := clusterManager.PickEndpoint("backend", nil)
		return endpoint != nil && endpoint.Address.Port == providerTwoPort
	}, 10*time.Second, 25*time.Millisecond, "cluster-only update was not applied before the RPC assertion")
	requireRuntimeE2EPostContains(t, rpcURL, `{"name":"cluster-only"}`, `rpc-one:cluster-only`, 15*time.Second)

	runtimeE2EAdminForm(t, router, token, http.MethodPost, "/config/api/listener",
		runtimeE2EDubboListenerYAML(rpcListenerPort, apiTwo))
	requireRuntimeE2EPostContains(t, rpcURL, `{"name":"bob"}`, `rpc-two:bob`, 15*time.Second)

	runtimeE2EAdminRequest(t, router, token, http.MethodDelete, "/config/api/listener?listener=gateway")
	runtimeE2EAdminRequest(t, router, token, http.MethodDelete, "/config/api/cluster?clusterId=2")
	requireRuntimeE2EPortClosed(t, rpcListenerPort)
}

func runtimeE2EEtcdEndpoint(t *testing.T) string {
	t.Helper()
	if endpoint := strings.TrimSpace(os.Getenv("PIXIU_E2E_ETCD_ENDPOINT")); endpoint != "" {
		return endpoint
	}

	clientURL := runtimeE2EReservedURL(t)
	peerURL := runtimeE2EReservedURL(t)
	config := embed.NewConfig()
	config.Name = fmt.Sprintf("pixiu-e2e-%d", time.Now().UnixNano())
	config.Dir = t.TempDir()
	config.LCUrls = []url.URL{clientURL}
	config.ACUrls = []url.URL{clientURL}
	config.LPUrls = []url.URL{peerURL}
	config.APUrls = []url.URL{peerURL}
	config.InitialCluster = config.InitialClusterFromName(config.Name)

	etcd, err := embed.StartEtcd(config)
	require.NoError(t, err)
	t.Cleanup(etcd.Close)
	select {
	case <-etcd.Server.ReadyNotify():
	case <-time.After(10 * time.Second):
		etcd.Server.Stop()
		t.Fatal("embedded etcd did not become ready")
	}
	return clientURL.String()
}

// runtimeE2EReservedPorts records every probe port this test binary has handed
// out. Each probe closes its listener immediately, so the kernel can hand the
// same port to a later probe; two consumers configured with it would then
// fight over one bind.
var runtimeE2EReservedPorts sync.Map

func runtimeE2EReservedURL(t *testing.T) url.URL {
	t.Helper()
	parsed, err := url.Parse(fmt.Sprintf("http://127.0.0.1:%d", reserveRuntimeE2EPort(t)))
	require.NoError(t, err)
	return *parsed
}

func startRuntimeE2EDubboProvider(t *testing.T, marker string) int {
	t.Helper()
	port := reserveRuntimeE2EPort(t)
	executable, err := os.Executable()
	require.NoError(t, err)
	logPath := filepath.Join(t.TempDir(), marker+".log")
	logFile, err := os.Create(logPath)
	require.NoError(t, err)

	cmd := exec.Command(executable, "-test.run=^TestRuntimeE2EDubboProviderProcess$", "-test.v")
	cmd.Env = append(os.Environ(),
		runtimeE2EProviderProcessEnv+"=1",
		runtimeE2EProviderPortEnv+"="+strconv.Itoa(port),
		runtimeE2EProviderMarkerEnv+"="+marker,
	)
	cmd.Stdout = logFile
	cmd.Stderr = logFile
	require.NoError(t, cmd.Start())

	var stopOnce sync.Once
	stop := func() {
		stopOnce.Do(func() {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
			_ = logFile.Close()
		})
	}
	t.Cleanup(stop)

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		connection, dialErr := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 100*time.Millisecond)
		if dialErr == nil {
			_ = connection.Close()
			return port
		}
		time.Sleep(25 * time.Millisecond)
	}
	stop()
	output, _ := os.ReadFile(logPath)
	t.Fatalf("Dubbo provider %q did not start on port %d:\n%s", marker, port, output)
	return 0
}

func writeRuntimeE2EDubboAPI(t *testing.T, path string, providerPort int) {
	t.Helper()
	content := fmt.Sprintf(`name: runtime-e2e
resources:
  - path: /rpc
    type: restful
    methods:
      - httpVerb: POST
        enable: true
        timeout: 5s
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          url: dubbo://127.0.0.1:%d
          protocol: dubbo
          interface: pixiu.e2e.Greeter
          method: SayHello
          parameterTypes:
            - java.lang.String
          serialization: hessian2
          mappingParams:
            - name: requestBody.name
              mapTo: "0"
              mapType: java.lang.String
`, providerPort)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
}

func watchConfigForRuntimeE2E(t *testing.T, ctx context.Context, publisher *adminxds.SnapshotPublisher) {
	t.Helper()
	watch, err := adminconfig.Client.WatchWithPrefix(adminconfig.Bootstrap.GetPath())
	require.NoError(t, err)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case response, ok := <-watch:
				if !ok {
					return
				}
				if response.Err() == nil && len(response.Events) > 0 {
					_ = publisher.Publish(ctx)
				}
			}
		}
	}()
}

func newRuntimeE2EXDSServer(t *testing.T, ctx context.Context, snapshotCache cache.SnapshotCache) *runtimeE2EXDSServer {
	t.Helper()
	server := &runtimeE2EXDSServer{
		t:             t,
		ctx:           ctx,
		snapshotCache: snapshotCache,
		nacks:         make(chan string, 16),
	}
	server.Start()
	return server
}

func (s *runtimeE2EXDSServer) Start() {
	s.t.Helper()
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.grpcServer != nil {
		return
	}
	address := s.address
	if address == "" {
		address = "127.0.0.1:0"
	}
	listener, err := net.Listen("tcp", address)
	require.NoError(s.t, err)
	callbacks := envoyserver.CallbackFuncs{StreamDeltaRequestFunc: func(_ int64, request *discoverypb.DeltaDiscoveryRequest) error {
		if request.GetErrorDetail() != nil {
			select {
			case s.nacks <- request.GetErrorDetail().GetMessage():
			default:
			}
		}
		return nil
	}}
	grpcServer := grpc.NewServer()
	extensionpb.RegisterExtensionConfigDiscoveryServiceServer(
		grpcServer,
		envoyserver.NewServer(s.ctx, s.snapshotCache, callbacks),
	)
	go func() { _ = grpcServer.Serve(listener) }()
	s.address = listener.Addr().String()
	s.listener = listener
	s.grpcServer = grpcServer
}

func (s *runtimeE2EXDSServer) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.grpcServer == nil {
		return
	}
	s.grpcServer.Stop()
	_ = s.listener.Close()
	s.grpcServer = nil
	s.listener = nil
}

func splitAddress(t *testing.T, address string) (string, int) {
	t.Helper()
	host, portText, err := net.SplitHostPort(address)
	require.NoError(t, err)
	port, err := strconv.Atoi(portText)
	require.NoError(t, err)
	return host, port
}

// reserveRuntimeE2EPort returns a loopback port that no earlier reservation in
// this process returned. The probe listener is closed before returning, so the
// port is only preflighted: a consumer binding it later can still race the
// kernel's transient port pool, which is why the port must not be shared.
func reserveRuntimeE2EPort(t *testing.T) int {
	t.Helper()
	for attempt := 0; attempt < 100; attempt++ {
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		require.NoError(t, err)
		port := listener.Addr().(*net.TCPAddr).Port
		require.NoError(t, listener.Close())
		if _, used := runtimeE2EReservedPorts.LoadOrStore(port, struct{}{}); !used {
			return port
		}
	}
	t.Fatal("could not reserve an unused loopback port")
	return 0
}

func runtimeE2EUpstream(name string) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		_, _ = io.WriteString(writer, name+"-"+request.Header.Get("X-XDS-E2E-Marker"))
	})
}

func runtimeE2EToken(t *testing.T) string {
	t.Helper()
	claims := auth.CustomClaims{
		Username: "xds-e2e",
		// SA1019: the upstream auth.CustomClaims embeds the deprecated
		// jwt.StandardClaims; migrate it together with the auth package.
		StandardClaims: jwt.StandardClaims{ //nolint:staticcheck // upstream CustomClaims requires it
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
		},
	}
	token, err := auth.NewJWT().CreateToken(claims)
	require.NoError(t, err)
	return token
}

func runtimeE2EAdminForm(t *testing.T, router http.Handler, token, method, path, content string) {
	t.Helper()
	body := url.Values{"content": []string{content}}.Encode()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	runtimeE2EServeAdmin(t, router, token, request)
}

func runtimeE2EAdminRequest(t *testing.T, router http.Handler, token, method, path string) {
	t.Helper()
	runtimeE2EServeAdmin(t, router, token, httptest.NewRequest(method, path, nil))
}

func runtimeE2EServeAdmin(t *testing.T, router http.Handler, token string, request *http.Request) {
	t.Helper()
	request.Header.Set("token", token)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	require.Equal(t, http.StatusOK, recorder.Code)
	require.Contains(t, recorder.Body.String(), fmt.Sprintf(`"code":"%s"`, adminconfig.OK))
}

func runtimeE2EClusterYAML(id int, address string, port int) string {
	return fmt.Sprintf("id: %d\nname: backend\ntype: Static\naddress: %s\nport: %d\n", id, address, port)
}

func runtimeE2EListenerYAML(port int, marker string) string {
	return fmt.Sprintf(`name: gateway
address:
  socket-address:
    address: 127.0.0.1
    port: %d
route_config:
  routes:
    - match:
        prefix: /
      route:
        cluster: backend
        cluster_not_found_response_code: 503
http_filters:
  - name: %s
    config:
      marker: %q
`, port, runtimeE2EHTTPFilter, marker)
}

func runtimeE2EDubboListenerYAML(port int, apiPath string) string {
	return fmt.Sprintf(`name: gateway
address:
  socket-address:
    address: 127.0.0.1
    port: %d
route_config:
  routes:
    - match:
        prefix: /
      route:
        cluster: backend
        cluster_not_found_response_code: 503
http_filters:
  - name: dgp.filter.http.apiconfig
    config:
      path: %q
  - name: dgp.filter.http.dubboproxy
    config:
      dubboProxyConfig:
        registries: {}
`, port, apiPath)
}

func runtimeE2EGet(target string) string {
	body, status, err := runtimeE2EGetResult(target)
	if err != nil || status != http.StatusOK {
		return ""
	}
	return body
}

func requireRuntimeE2EBody(t *testing.T, target, want string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var lastBody string
	var lastStatus int
	var lastErr error
	for time.Now().Before(deadline) {
		lastBody, lastStatus, lastErr = runtimeE2EGetResult(target)
		if lastErr == nil && lastStatus == http.StatusOK && lastBody == want {
			return
		}
		time.Sleep(25 * time.Millisecond)
	}
	t.Fatalf("runtime response did not become %q: last status=%d body=%q err=%v", want, lastStatus, lastBody, lastErr)
}

func requireRuntimeE2EPostContains(t *testing.T, target, requestBody, want string, timeout time.Duration) {
	t.Helper()
	deadline := time.Now().Add(timeout)
	var lastBody string
	var lastStatus int
	var lastErr error
	for time.Now().Before(deadline) {
		request, err := http.NewRequest(http.MethodPost, target, strings.NewReader(requestBody))
		if err != nil {
			t.Fatal(err)
		}
		request.Header.Set("Content-Type", "application/json")
		response, err := (&http.Client{Timeout: time.Second}).Do(request)
		if err == nil {
			lastStatus = response.StatusCode
			body, readErr := io.ReadAll(response.Body)
			_ = response.Body.Close()
			lastBody = string(body)
			lastErr = readErr
		} else {
			lastErr = err
		}
		if lastErr == nil && lastStatus == http.StatusOK && strings.Contains(lastBody, want) {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("runtime RPC response did not contain %q: last status=%d body=%q err=%v", want, lastStatus, lastBody, lastErr)
}

func requireRuntimeE2EPortClosed(t *testing.T, port int) {
	t.Helper()
	require.Eventually(t, func() bool {
		connection, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", port), 50*time.Millisecond)
		if err == nil {
			_ = connection.Close()
			return false
		}
		return true
	}, 10*time.Second, 25*time.Millisecond)
}

func runtimeE2EGetResult(target string) (string, int, error) {
	client := &http.Client{Timeout: 200 * time.Millisecond}
	response, err := client.Get(target)
	if err != nil {
		return "", 0, err
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", response.StatusCode, err
	}
	return string(body), response.StatusCode, nil
}
