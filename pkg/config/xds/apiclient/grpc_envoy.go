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
	"dubbo.apache.org/dubbo-go/v3/common/logger"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/xds"
	clusterpb "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoyconfigcorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	listenerpb "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	structpb "github.com/golang/protobuf/ptypes/struct"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/types/known/anypb"
	"time"
)

// CreateEnvoyGrpcApiClient create Grpc type ApiClient
func CreateEnvoyGrpcApiClient(config *model.ApiConfigSource, node *model.Node,
	exitCh chan struct{},
	typeNames ...ResourceTypeName) *AggGrpcApiClient {
	v := &AggGrpcApiClient{}
	v.config = *config
	v.node = node
	v.grpcMg = grpcMg
	v.exitCh = exitCh
	switch typeNames[0] {
	case xds.ListenerType:
		v.resourceNames = []string{resource.ListenerType}
	case xds.ClusterType:
		v.resourceNames = []string{resource.ClusterType}
	default:
		logger.Warnf("typeNames should be dubbo-go.pixiu/v1/discovery:cluster or dubbo-go.pixiu/v1/discovery:listener")
		v.resourceNames = []string{resource.ClusterType}
	}
	v.init()
	return v
}

type (
	AggGrpcApiClient struct {
		GrpcApiClient
		xDSAggClient discoverypb.AggregatedDiscoveryServiceClient
		interfaceMap *InterfaceMapHandlerImpl
	}
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

	g.interfaceMap = NewInterfaceMapHandlerImpl("localhost:18080", "") //todo use config url
}

func (g *AggGrpcApiClient) Fetch(localVersion string) ([]*ProtoAny, error) {
	//un-support fetch
	return nil, nil
}

func (g *AggGrpcApiClient) Delta() (chan *DeltaResources, error) {
	outputCh := make(chan *DeltaResources)
	return outputCh, g.runDelta(outputCh)
}

func (g *AggGrpcApiClient) runDelta(output chan<- *DeltaResources) error {
	var delta discoverypb.AggregatedDiscoveryService_StreamAggregatedResourcesClient
	var cancel context.CancelFunc
	backoff := func() {
		for {
			//back off
			var err error
			var ctx context.Context // context to sync exitCh
			ctx, cancel = context.WithCancel(context.TODO())
			delta, err = g.sendInitAggRequest(ctx)
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
			for { //loop consume recv data form xds server(sendInitDeltaRequest)
				resp, err := delta.Recv()
				if err != nil { //todo backoff retry
					logger.Error("can not receive delta discovery request", err)
					break
				}
				g.handleDeltaResponse(resp, output)

				//err = g.subscribeOnGoingChang(delta)
				//if err != nil {
				//	logger.Error("can not recv delta discovery request", err)
				//	break
				//}
			}
			backoff()
		}
	}()

	return nil
}

func (g *AggGrpcApiClient) handleDeltaResponse(resp *discoverypb.DiscoveryResponse, output chan<- *DeltaResources) {
	// save the xds state
	g.xdsState.deltaVersion = make(map[string]string, 1)
	g.xdsState.nonce = resp.Nonce
	g.xdsState.versionInfo = resp.VersionInfo

	resources := &DeltaResources{
		NewResources:    make([]*ProtoAny, 0, 1),
		RemovedResource: make([]string, 0, 1),
	}
	logger.Infof("get xDS message nonce, %s", resp.Nonce)
	//for _, res := range resp.Resources { todo
	//	logger.Infof("remove resource found %x", res)
	//	resources.RemovedResource = append(resources.RemovedResource, res.TypeUrl) //todo
	//}

	dubboService, err := g.interfaceMap.GetServiceUniqueKeyHostAddrMapFromPilot()
	if err != nil {
		logger.Warnf("cannot get dubbo sevice", err)
		return
	}

	logger.Infof("all dubbo service is: %v", dubboService)

	for _, res := range resp.Resources {
		logger.Infof("new resource found %s version=%s", res.TypeUrl)
		//g.xdsState.deltaVersion[res.Name] = res.Version // needn't to do this
		elems, err := g.decodeSource(res)
		if err != nil {
			logger.Warnf("can not decode source %s , %v", res.TypeUrl, res)
		}
		resources.NewResources = append(resources.NewResources, elems)
	}
	//notify the resource change handler
	//output <- resources
}

func (g *AggGrpcApiClient) sendInitAggRequest(ctx context.Context) (stream discoverypb.AggregatedDiscoveryService_StreamAggregatedResourcesClient, err error) {
	stream, err = g.xDSAggClient.StreamAggregatedResources(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "fetch dynamic resource from remote error. %s", g.resourceNames)
	}

	err = stream.Send(&discoverypb.DiscoveryRequest{
		VersionInfo:   "localVersion",
		Node:          g.makeNode(),
		ResourceNames: nil,
		TypeUrl:       resource.FetchRoutes, //g.resourceNames[0], //"type.googleapis.com/envoy.config.listener.v3.Listener",
		ResponseNonce: g.xdsState.nonce,
		ErrorDetail:   nil, //todo use error detail
	})
	if err != nil {
		return nil, errors.Wrapf(err, "fetch dynamic resource from remote error. %s", g.resourceNames)
	}
	return
}

func (g *AggGrpcApiClient) subscribeOnGoingChang(delta discoverypb.AggregatedDiscoveryService_StreamAggregatedResourcesClient) error {
	err := delta.Send(&discoverypb.DiscoveryRequest{
		VersionInfo:   g.xdsState.versionInfo,
		Node:          g.makeNode(),
		ResourceNames: nil,
		TypeUrl:       g.resourceNames[0], //"type.googleapis.com/envoy.config.listener.v3.Listener",
		ResponseNonce: g.xdsState.nonce,
		ErrorDetail:   nil, //todo use error detail
	})
	return err
}

func (g *AggGrpcApiClient) makeNode() *envoyconfigcorev3.Node {
	return &envoyconfigcorev3.Node{
		Id:                   "sidecar~172.1.1.1~sleep-55b5877479-rwcct.default~default.svc.cluster.local", //todo read from sdk
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

func (g *AggGrpcApiClient) decodeSource(oneResource *anypb.Any) (*ProtoAny, error) {
	switch oneResource.GetTypeUrl() {
	case resource.ListenerType:
		l := listenerpb.Listener{}
		if err := oneResource.UnmarshalTo(&l); err != nil {
			return nil, errors.Wrap(err, "unmarshal error")
		}
	case resource.ClusterType:
		l := clusterpb.Cluster{}
		if err := oneResource.UnmarshalTo(&l); err != nil {
			return nil, errors.Wrap(err, "unmarshal error")
		}
		fmt.Println("xxxx->>>", l.Metadata)
	default:
		return nil, errors.Errorf("unkown typeurl of %s", oneResource.TypeUrl)
	}
	elems := &ProtoAny{any: oneResource}
	return elems, nil
}
