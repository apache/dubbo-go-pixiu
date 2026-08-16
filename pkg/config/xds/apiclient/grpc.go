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

package apiclient

import (
	"context"
	stderr "errors"
	"sync"
	"time"
)

import (
	envoyconfigcorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	grpcstatus "google.golang.org/grpc/status"

	"google.golang.org/protobuf/types/known/anypb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls"
)

// agent name to talk with xDS server
const xdsAgentName = "dubbo-go-pixiu"

var (
	grpcMg             *GRPCClusterManager
	ErrClusterNotFound = stderr.New("can not find cluster")
)

type (
	GrpcExtensionApiClient struct {
		config             model.ApiConfigSource
		grpcMg             *GRPCClusterManager
		node               *model.Node
		xDSExtensionClient extensionpb.ExtensionConfigDiscoveryServiceClient
		resourceNames      []ResourceTypeName
		typeUrl            string
		exitCh             chan struct{}
	}
	xdsState struct {
		deltaVersion map[string]string
	}
)

func Init(clusterMg controls.ClusterManager) {
	grpcMg = &GRPCClusterManager{
		clusters:  &sync.Map{},
		clusterMg: clusterMg,
	}
}

func Stop() {
	if err := grpcMg.Close(); err != nil { //todo
		logger.Errorf("grpcClusterManager close failed. %v", err)
	}
}

// CreateGrpExtensionApiClient create Grpc type ApiClient
func CreateGrpExtensionApiClient(config *model.ApiConfigSource, node *model.Node,
	exitCh chan struct{},
	typeName ResourceTypeName) *GrpcExtensionApiClient {
	v := &GrpcExtensionApiClient{
		config:  *config,
		node:    node,
		typeUrl: typeName,
		grpcMg:  grpcMg,
		exitCh:  exitCh,
	}
	v.init()
	return v
}

// Fetch get config data from discovery service and return Any type.
func (g *GrpcExtensionApiClient) Fetch(localVersion string) ([]*ProtoAny, error) {
	clsRsp, err := g.xDSExtensionClient.FetchExtensionConfigs(context.Background(), &discoverypb.DiscoveryRequest{
		VersionInfo:   localVersion,
		Node:          g.makeNode(),
		ResourceNames: g.resourceNames,
		TypeUrl:       resource.ExtensionConfigType, //"type.googleapis.com/pixiu.config.listener.v3.Listener", //resource.ListenerType,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "fetch dynamic resource from remote error. %s", g.resourceNames)
	}
	logger.Infof("init the from xds server typeUrl=%s version=%s", clsRsp.TypeUrl, clsRsp.VersionInfo)
	extensions := make([]*ProtoAny, 0, len(clsRsp.Resources))
	for _, clsResource := range clsRsp.Resources {
		elems, err := g.decodeSource(clsResource)
		if err != nil {
			return nil, err
		}
		extensions = append(extensions, elems)
	}
	return extensions, nil
}

func (g *GrpcExtensionApiClient) decodeSource(resource *anypb.Any) (*ProtoAny, error) {
	extension := envoyconfigcorev3.TypedExtensionConfig{}
	err := resource.UnmarshalTo(&extension)
	if err != nil {
		return nil, errors.Wrapf(err, "typed extension as expected.(%s)", g.resourceNames)
	}
	elems := &ProtoAny{typeConfig: &extension}
	return elems, nil
}

func (g *GrpcExtensionApiClient) makeNode() *envoyconfigcorev3.Node {
	return &envoyconfigcorev3.Node{
		Id:            g.node.Id,
		Cluster:       g.node.Cluster,
		UserAgentName: xdsAgentName,
	}
}

func (g *GrpcExtensionApiClient) Delta() (chan *DeltaResources, error) {
	outputCh := make(chan *DeltaResources)
	return outputCh, g.runDelta(outputCh)
}

func (g *GrpcExtensionApiClient) runDelta(output chan<- *DeltaResources) error {
	ctx, cancel := context.WithCancel(context.Background())
	if g.exitCh != nil {
		go func() {
			select {
			case <-g.exitCh:
				cancel()
			case <-ctx.Done():
			}
		}()
	}

	xState := &xdsState{deltaVersion: make(map[string]string)}
	delta, err := g.connectDelta(ctx, xState)
	if err != nil {
		cancel()
		close(output)
		if stderr.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	go func() {
		defer cancel()
		defer close(output)
		for {
			if err := g.consumeDeltaStream(ctx, delta, xState, output); err != nil && !stderr.Is(err, context.Canceled) {
				logger.Errorf("delta extension config stream closed: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			var err error
			delta, err = g.connectDelta(ctx, xState)
			if err != nil {
				return
			}
		}
	}()
	return nil
}

func (g *GrpcExtensionApiClient) connectDelta(ctx context.Context, xState *xdsState) (extensionpb.ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient, error) {
	for {
		delta, err := g.sendInitDeltaRequest(ctx, xState)
		if err == nil {
			return delta, nil
		}
		logger.Errorf("can not start delta extension config stream; retrying in 1s: %v", err)
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (g *GrpcExtensionApiClient) consumeDeltaStream(ctx context.Context, delta extensionpb.ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient, xState *xdsState, output chan<- *DeltaResources) error {
	for {
		resp, err := delta.Recv()
		if err != nil {
			return err
		}

		resources, decodeErr := g.handleDeltaResponse(resp)
		if decodeErr != nil {
			if err := g.respondToDelta(delta, resp.Nonce, decodeErr); err != nil {
				return err
			}
			continue
		}

		select {
		case output <- resources:
		case <-ctx.Done():
			return ctx.Err()
		}

		var applyErr error
		select {
		case applyErr = <-resources.applyResult:
		case <-ctx.Done():
			return ctx.Err()
		}
		if applyErr == nil {
			acceptDeltaResources(xState, resources)
		}
		if err := g.respondToDelta(delta, resp.Nonce, applyErr); err != nil {
			return err
		}
	}
}

func (g *GrpcExtensionApiClient) handleDeltaResponse(resp *discoverypb.DeltaDiscoveryResponse) (*DeltaResources, error) {
	if resp == nil {
		return nil, errors.New("delta discovery response is nil")
	}
	resources := newDeltaResources()
	logger.Infof("get xDS message nonce, %s", resp.Nonce)
	for _, res := range resp.RemovedResources {
		logger.Infof("remove resource found %s", res)
		resources.RemovedResources = append(resources.RemovedResources, res)
	}

	for _, res := range resp.Resources {
		logger.Infof("new resource found %s version=%s", res.Name, res.Version)
		elems, err := g.decodeSource(res.Resource)
		if err != nil {
			return nil, errors.Wrapf(err, "decode resource %s version=%s", res.Name, res.Version)
		}
		resources.NewResources = append(resources.NewResources, elems)
		resources.versions[res.Name] = res.Version
	}
	return resources, nil
}

func acceptDeltaResources(xState *xdsState, resources *DeltaResources) {
	for _, name := range resources.RemovedResources {
		delete(xState.deltaVersion, name)
	}
	for name, version := range resources.versions {
		xState.deltaVersion[name] = version
	}
}

func (g *GrpcExtensionApiClient) respondToDelta(delta extensionpb.ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient, nonce string, applyErr error) error {
	req := &discoverypb.DeltaDiscoveryRequest{
		Node:          g.makeNode(),
		TypeUrl:       resource.ExtensionConfigType,
		ResponseNonce: nonce,
	}
	if applyErr != nil {
		req.ErrorDetail = grpcstatus.New(codes.InvalidArgument, applyErr.Error()).Proto()
	}
	return delta.Send(req)
}

func (g *GrpcExtensionApiClient) sendInitDeltaRequest(ctx context.Context, xState *xdsState) (extensionpb.ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient, error) {
	delta, err := g.xDSExtensionClient.DeltaExtensionConfigs(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "can not start delta stream with xds server ")
	}
	err = delta.Send(&discoverypb.DeltaDiscoveryRequest{
		Node:                     g.makeNode(),
		TypeUrl:                  resource.ExtensionConfigType,
		ResourceNamesSubscribe:   []string{g.typeUrl},
		ResourceNamesUnsubscribe: nil,
		InitialResourceVersions:  xState.deltaVersion,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "can not send delta discovery request")
	}
	return delta, nil
}

func (g *GrpcExtensionApiClient) init() {
	if len(g.config.ClusterName) == 0 {
		panic("should config one cluster at least")
	}
	//todo implement multiple grpc api services
	if len(g.config.ClusterName) > 1 {
		logger.Warnf("defined multiple cluster for xDS api services but only one support.")
	}
	cluster, err := g.grpcMg.GetGrpcCluster(g.config.ClusterName[0])

	if err != nil {
		logger.Errorf("get cluster for init error. error=%v", err)
		panic(err)
	}
	conn, err := cluster.GetConnection()
	if err != nil {
		panic(err)
	}
	g.xDSExtensionClient = extensionpb.NewExtensionConfigDiscoveryServiceClient(conn)
}

type GRPCClusterManager struct {
	clusters  *sync.Map // map[clusterName]*grpcCluster
	clusterMg controls.ClusterManager
}

type GRPCCluster struct {
	name   string //cluster name
	config *model.ClusterConfig
	once   sync.Once
	conn   *grpc.ClientConn
}

// GetGrpcCluster get the cluster or create it first time.
func (g *GRPCClusterManager) GetGrpcCluster(name string) (*GRPCCluster, error) {
	if !g.clusterMg.HasCluster(name) {
		return nil, errors.Wrapf(ErrClusterNotFound, "name = %s", name)
	}

	if load, ok := g.clusters.Load(name); ok {
		grpcCluster := load.(*GRPCCluster) // grpcClusterManager only
		return grpcCluster, nil
	}

	store, err := g.clusterMg.CloneXdsControlStore()
	if err != nil {
		return nil, errors.WithMessagef(err, "clone cluster store failed")
	}

	var clusterCfg *model.ClusterConfig
	for _, cfg := range store.Config() {
		if cfg.Name == name {
			clusterCfg = cfg
			break
		}
	}
	if clusterCfg == nil {
		return nil, errors.Wrapf(ErrClusterNotFound, "name of %s", name)
	}
	newCluster := &GRPCCluster{
		name:   name,
		config: clusterCfg,
	}
	g.clusters.Store(name, newCluster)
	return newCluster, nil
}

func (g *GRPCClusterManager) Close() (err error) {
	//todo enhance the close process when concurrent
	g.clusters.Range(func(_, value any) bool {
		if conn := value.(*grpc.ClientConn); conn != nil {
			if err = conn.Close(); err != nil {
				logger.Errorf("can not close grpc connection.", err)
			}
		}
		return true
	})
	return nil
}

func (g *GRPCCluster) GetConnection() (conn *grpc.ClientConn, err error) {
	g.once.Do(func() {
		creds := insecure.NewCredentials()
		//if *xdsCreds { // todo
		//	log.Println("Using xDS credentials...")
		//	var err error
		//	if creds, err = xdscreds.NewClientCredentials(xdscreds.ClientOptions{FallbackCreds: insecure.NewCredentials()}); err != nil {
		//		log.Fatalf("failed to create client-side xDS credentials: %v", err)
		//	}
		//}
		if len(g.config.Endpoints) == 0 {
			err = errors.Errorf("expect endpoint.")
			return
		}
		endpoint := g.config.Endpoints[0].Address.GetAddress()
		logger.Infof("to connect xds server %s ...", endpoint)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) //todo fix timeout cancel warning
		defer cancel()
		conn, err = grpc.DialContext(ctx, endpoint,
			grpc.WithTransportCredentials(creds),
			grpc.WithBlock(),
		)
		if err != nil {
			err = errors.Errorf("grpc.Dial(%s) failed: %v", endpoint, err)
			return
		}
		logger.Infof("connected xds server (%s)", endpoint)
		g.conn = conn
	})
	return g.conn, nil
}

func (g *GRPCCluster) IsAlive() (alive bool) {
	if g.conn == nil {
		return false
	}
	state := connectivity.Ready
	defer func() {
		if r := recover(); r != nil {
			logger.Warnf("grpc connection state check panicked: %v", r)
			alive = state == connectivity.Ready
		}
	}()
	state = g.conn.GetState()
	alive = state == connectivity.Ready
	return
}

func (g *GRPCCluster) Close() error {
	if g.conn == nil {
		return nil
	}
	if err := g.conn.Close(); err != nil {
		return errors.Wrapf(err, "can not close. %v", g.config)
	}
	g.conn = nil
	return nil
}
