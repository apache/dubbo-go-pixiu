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
	"os"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/xds"
	xdspb "github.com/dubbo-go-pixiu/pixiu-api/pkg/xds/model"
	clusterpb "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoyconfigcorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointpb "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	listenerpb "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type GrpcApiClientOption func(*AggGrpcApiClient)

func WithIstioService(serviceNames ...string) GrpcApiClientOption {
	return func(g *AggGrpcApiClient) {
		for _, name := range serviceNames {
			g.dubboServiceFilter[name] = struct{}{}
		}
	}
}

// CreateEnvoyGrpcApiClient create Grpc type ApiClient
func CreateEnvoyGrpcApiClient(
	config *model.ApiConfigSource,
	node *model.Node,
	exitCh chan struct{},
	typeName ResourceTypeName,
	opts ...GrpcApiClientOption,
) *AggGrpcApiClient {
	v := &AggGrpcApiClient{}
	v.config = *config
	v.node = node
	v.grpcMg = grpcMg
	v.exitCh = exitCh
	switch typeName {
	case xds.ListenerType:
		v.typeUrl = resource.ListenerType
	case xds.ClusterType:
		v.typeUrl = resource.ClusterType
	case xds.EndpointType:
		v.typeUrl = resource.EndpointType
	default:
		logger.Warnf("typeName should be dubbo-go.pixiu/v1/discovery:cluster or dubbo-go.pixiu/v1/discovery:listener")
		v.typeUrl = resource.ClusterType
	}

	v.dubboServiceFilter = make(map[string]struct{}, 10)
	for _, fn := range opts {
		fn(v)
	}
	v.init()
	return v
}

type (
	AggGrpcApiClient struct {
		GrpcExtensionApiClient
		xDSAggClient       discoverypb.AggregatedDiscoveryServiceClient
		dubboServiceFilter map[string]struct{} // the service name to filter
	}
	discoveryResponseHandler func(any2 []*anypb.Any)
)

func (g *AggGrpcApiClient) init() {
	if len(g.config.ClusterName) == 0 {
		panic("should config one cluster at least")
	}
	//todo implement multiple grpc api services
	if len(g.config.ClusterName) > 1 {
		logger.Warn("defined multiple cluster for xDS api services but only one support.")
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
	g.xDSAggClient = discoverypb.NewAggregatedDiscoveryServiceClient(conn)

}

func (g *AggGrpcApiClient) Fetch(_ string) ([]*ProtoAny, error) {
	//un-support fetch
	return nil, nil
}

func (g *AggGrpcApiClient) Delta() (chan *DeltaResources, error) {
	outputCh := make(chan *DeltaResources)
	return outputCh, g.pipeline(outputCh)
}

type refEndpoint struct {
	IsPending bool
	RawProto  proto.Message
}

func (g *AggGrpcApiClient) pipeline(output chan *DeltaResources) error {
	// all endpoint refer cluster or listener. map[type]map[resource name]endpoint info
	edsResources := make(map[resource.Type]map[string]refEndpoint)
	req := g.makeDiscoveryRequest(g.resourceNames, g.typeUrl)
	var handler discoveryResponseHandler
	switch g.typeUrl {
	case resource.ListenerType:
		handler = func(any2 []*anypb.Any) {
			for _, res := range any2 {
				l := listenerpb.Listener{}
				if err := res.UnmarshalTo(&l); err != nil {
					logger.Warnf("can not decode source %s , %v", res.TypeUrl, res)
					continue
				}
			}
		}
	case resource.ClusterType:
		handler = func(any2 []*anypb.Any) {
			// only one goroutine handle response, no need to lock for local var.
			for _, res := range any2 {
				logger.Infof("new resource found %s", res.TypeUrl)
				c := clusterpb.Cluster{}
				if err := res.UnmarshalTo(&c); err != nil {
					logger.Warnf("can not decode source %s , %v", res.TypeUrl, res)
					continue
				}
				//needn't lock edsResources
				g.getClusterResourceReference(&c, edsResources)
			}

			pendingResourceNames := make([]string, 0)
			clusterRefEndpoints := edsResources[resource.ClusterType]
			// list all pending pendingResourceNames
			for name, b := range clusterRefEndpoints {
				if b.IsPending {
					pendingResourceNames = append(pendingResourceNames, name)
				}
			}

			// do not block, watch new resource at another goroutine
			err := g.runEndpointReferences(pendingResourceNames, func(any2 []*anypb.Any) {
				// run on another goroutine
				extCluster := xdspb.PixiuExtensionClusters{
					Clusters: []*xdspb.Cluster{
						{
							Name:             "",
							TypeStr:          xds.ClusterType,
							Type:             0,
							EdsClusterConfig: nil,
							LbStr:            "",
							Lb:               0,
							HealthChecks:     nil,
							Endpoints:        make([]*xdspb.Endpoint, 0, len(any2)),
						},
					},
				}

				for _, one := range any2 {
					l := endpointpb.ClusterLoadAssignment{}
					if err := one.UnmarshalTo(&l); err != nil {
						logger.Warnf("unmarshal error", err)
						continue
					}
					c := clusterRefEndpoints[l.ClusterName].RawProto.(*clusterpb.Cluster)
					extCluster.Clusters[0].Name = g.readServiceNameOfCluster(c)

					for _, ep := range l.Endpoints {
						address := ep.LbEndpoints[0].GetEndpoint().GetAddress().GetSocketAddress()
						extCluster.Clusters[0].Endpoints = append(extCluster.Clusters[0].Endpoints, &xdspb.Endpoint{
							Id:   "",
							Name: "",
							Address: &xdspb.SocketAddress{
								Address: address.Address,
								Port:    int64(address.GetPortValue()),
							},
							Metadata: nil,
						})
					}
				}

				//make output
				output <- &DeltaResources{
					NewResources: []*ProtoAny{
						{
							typeConfig: &envoyconfigcorev3.TypedExtensionConfig{
								Name: "cluster", //todo cluster name
								TypedConfig: func() *anypb.Any { //make any.Any from extCluster
									a, err := anypb.New(&extCluster)
									if err != nil {
										logger.Warnf("can not make anypb.Any %v", err)
										return nil
									}
									return a
								}(),
							},
						},
					},
					RemovedResource: nil,
				}
			})
			if err != nil { //todo retry
				logger.Errorf("can not run reference request %v", err)
			}
		}

	default:
		return errors.Errorf("nedd listenerType of clusterType but get %s", g.typeUrl)
	}

	if err := g.runDelta(req, handler); err != nil {
		return errors.WithMessagef(err, "start run %s failed", req.TypeUrl)
	}

	return nil
}

// readServiceNameOfCluster get service name of k8s
func (g *AggGrpcApiClient) readServiceNameOfCluster(c *clusterpb.Cluster) string {
	if c.Metadata == nil {
		return ""
	}
	return c.
		Metadata.
		FilterMetadata["istio"].
		Fields["services"].
		GetListValue().
		GetValues()[0].
		GetStructValue().
		Fields["name"].
		GetStringValue()
}

// request EDS for the allResourceNames
func (g *AggGrpcApiClient) runEndpointReferences(allResourceNames []string,
	output discoveryResponseHandler) (err error) {

	//todo reload all request
	req := g.makeDiscoveryRequest(allResourceNames, resource.EndpointType)
	if err := g.runDelta(req, output); err != nil {
		return errors.WithMessagef(err, "start run %s failed", req.TypeUrl)
	}
	return nil
}

// runDelta start 2 goroutine to and watch change
func (g *AggGrpcApiClient) runDelta(req *discoverypb.DiscoveryRequest, output discoveryResponseHandler) error {
	var delta discoverypb.AggregatedDiscoveryService_StreamAggregatedResourcesClient
	var cancel context.CancelFunc
	var xState xdsState
	//read resource list
	backoff := func() {
		xState = xdsState{}
		for {
			//back off
			var err error
			var ctx context.Context // context to sync exitCh
			ctx, cancel = context.WithCancel(context.TODO())
			delta, err = g.sendInitAggRequest(ctx, req, &xState)
			if err != nil {
				logger.Error("can not receive delta discovery request, will back off 1 sec later", err)
				select {
				case <-time.After(1 * time.Second):
				case <-g.exitCh:
					logger.Infof("get close single.")
					return
				}
				continue //backoff
			}
			return //success
		}
	}

	backoff()
	if delta == nil { // delta instance not created because exitCh
		return nil
	}
	go func() {
		//waiting exitCh close
		for range g.exitCh {
		}
		cancel()
	}()
	//get message
	go func() {
		for { // delta response backoff.
			for { //loop consume receive data form xds server(sendInitDeltaRequest)
				resp, err := delta.Recv()
				if err != nil {
					logger.Error("can not receive delta discovery request", err)
					break
				}
				g.handleDeltaResponse(resp, &xState, output)
			}
			backoff()
		}
	}()

	return nil
}

func (g *AggGrpcApiClient) handleDeltaResponse(resp *discoverypb.DiscoveryResponse, xdsState *xdsState, handler discoveryResponseHandler) {
	// save the xds state
	xdsState.deltaVersion = make(map[string]string, 1)
	xdsState.nonce = resp.Nonce
	xdsState.versionInfo = resp.VersionInfo
	handler(resp.Resources)
	//notify the resource change handler
	//output <- resources
}

func (g *AggGrpcApiClient) sendInitAggRequest(ctx context.Context, req *discoverypb.DiscoveryRequest, xState *xdsState) (stream discoverypb.AggregatedDiscoveryService_StreamAggregatedResourcesClient, err error) {
	req.VersionInfo = xState.versionInfo
	req.ResponseNonce = xState.nonce
	stream, err = g.xDSAggClient.StreamAggregatedResources(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetch dynamic resource from remote error. %s", g.resourceNames)
	}

	err = stream.Send(req)
	if err != nil {
		return nil, errors.Wrapf(err, "fetch dynamic resource from remote error. %s", g.resourceNames)
	}
	return
}

func (g *AggGrpcApiClient) makeDiscoveryRequest(resources []string,
	typeUrl string,
) *discoverypb.DiscoveryRequest {
	return &discoverypb.DiscoveryRequest{
		//VersionInfo:   xdsState.versionInfo,
		Node:          g.makeNode(),
		ResourceNames: resources, //[]string{"outbound|20000||dubbo-go-app.default.svc.cluster.local"},
		TypeUrl:       typeUrl,   //"type.googleapis.com/envoy.config.listener.v3.Listener",
		//ResponseNonce: xdsState.nonce,
		ErrorDetail: nil,
	}
}

func (g *AggGrpcApiClient) makeNode() *envoyconfigcorev3.Node {
	podId := os.Getenv("POD_IP")
	if len(podId) == 0 {
		logger.Warnf("expect POD_ID env")
		podId = "0.0.0.0"
	}
	podName := os.Getenv("POD_NAME")
	if len(podName) == 0 {
		logger.Warnf("expect POD_NAME env")
		podName = "pixiu-gateway"
	}
	nsName := os.Getenv("POD_NAMESPACE")
	if len(nsName) == 0 {
		logger.Warnf("expect POD_NAMESPACE env")
		nsName = "default"
	}

	return &envoyconfigcorev3.Node{
		Id:                   "sidecar~" + podId + "~" + podName + "." + nsName + ".svc.cluster.local",
		UserAgentName:        "pixiu",
		Cluster:              "testCluster",
		UserAgentVersionType: &envoyconfigcorev3.Node_UserAgentVersion{UserAgentVersion: "1.45.0"},
		ClientFeatures:       []string{"envoy.lb.does_not_support_overprovisioning"},
		Metadata: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"CLUSTER_ID": {
					Kind: &structpb.Value_StringValue{StringValue: "Kubernetes"},
				},
				"LABELS": {
					Kind: &structpb.Value_StructValue{StructValue: &structpb.Struct{
						Fields: map[string]*structpb.Value{},
					}},
				},
			},
		},
	}
}

// getClusterResourceReference get resources of cluster
func (g *AggGrpcApiClient) getClusterResourceReference(c *clusterpb.Cluster, edsResources map[resource.Type]map[string]refEndpoint) {
	logger.Infof("cluster name ==>%s", c.Name)
	if _, exist := g.dubboServiceFilter[g.readServiceNameOfCluster(c)]; !exist {
		logger.Infof("cluster name ==>%v", c)
		return
	}

	switch typ := c.ClusterDiscoveryType.(type) {
	case *clusterpb.Cluster_Type:
		if typ.Type == clusterpb.Cluster_EDS {
			name := c.Name
			if c.EdsClusterConfig != nil && c.EdsClusterConfig.ServiceName != "" {
				name = c.EdsClusterConfig.ServiceName
			}

			if _, ok := edsResources[resource.ClusterType]; !ok {
				edsResources[resource.ClusterType] = make(map[string]refEndpoint)
			}

			edsResources[resource.ClusterType][name] = refEndpoint{
				IsPending: true,
				RawProto:  c,
			}
		} else {
			logger.Infof("cluster type %s not supported", typ.Type.String())
		}
	}
}
