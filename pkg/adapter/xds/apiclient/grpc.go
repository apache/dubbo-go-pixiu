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
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"

	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	"github.com/pkg/errors"

	"google.golang.org/grpc"

	"google.golang.org/protobuf/types/known/anypb"
)

// agent name to talk with xDS server
const xdsAgentName = "dubbo-go-pixiu"

type (
	GrpcApiClient struct {
		config             model.ApiConfigSource
		grpcMg             GrpcClusterManager
		node               *model.Node
		xDSExtensionClient extensionpb.ExtensionConfigDiscoveryServiceClient
		resourceNames      []ResourceTypeName
		exitCh             chan struct{}
		xdsState           xdsState
	}
	xdsState struct {
		nonce        string
		deltaVersion map[string]string
	}

	GrpcCluster interface {
		GetConnect() *grpc.ClientConn
	}
	GrpcClusterManager interface {
		GetGrpcCluster(name string) (GrpcCluster, error)
	}
)

func CreateGrpcApiClient(config *model.ApiConfigSource, node *model.Node,
	grpcMg GrpcClusterManager,
	exitCh chan struct{},
	typeNames ...ResourceTypeName) *GrpcApiClient {
	v := &GrpcApiClient{
		config:        *config,
		node:          node,
		resourceNames: typeNames,
		grpcMg:        grpcMg,
		exitCh:        exitCh,
	}
	v.init()
	return v
}

// Fetch get config data from discovery service and return Any type.
func (g *GrpcApiClient) Fetch(localVersion string) ([]*ProtoAny, error) {
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
	for _, _resource := range clsRsp.Resources {
		elems, err := g.decodeSource(_resource)
		if err != nil {
			return nil, err
		}
		extensions = append(extensions, elems)
	}
	return extensions, nil
}

func (g *GrpcApiClient) decodeSource(_resource *anypb.Any) (*ProtoAny, error) {
	extension := envoy_config_core_v3.TypedExtensionConfig{}
	err := _resource.UnmarshalTo(&extension)
	if err != nil {
		return nil, errors.Wrapf(err, "typed extension as expected.(%s)", g.resourceNames)
	}
	elems := &ProtoAny{&extension}
	return elems, nil
}

func (g *GrpcApiClient) makeNode() *envoy_config_core_v3.Node {
	return &envoy_config_core_v3.Node{
		Id:            g.node.Id,
		Cluster:       g.node.Cluster,
		UserAgentName: xdsAgentName,
	}
}

func (g *GrpcApiClient) Delta() (chan *DeltaResources, error) {
	outputCh := make(chan *DeltaResources)
	return outputCh, g.runDelta(outputCh)
	//return outputCh, g.runStream(outputCh)
}

func (g *GrpcApiClient) runDelta(output chan<- *DeltaResources) error {
	var delta extensionpb.ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient
	backoff := func() {
		for {
			//back off
			var err error
			delta, err = g.sendInitDeltaRequest()
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

				err = g.subscribeOnGoingChang(delta)
				if err != nil {
					logger.Error("can not recv delta discovery request", err)
					break
				}
			}
			backoff()
		}
	}()

	return nil
}

func (g *GrpcApiClient) handleDeltaResponse(resp *discoverypb.DeltaDiscoveryResponse, output chan<- *DeltaResources) {
	// save the xds state
	g.xdsState.deltaVersion = make(map[string]string, 1)
	g.xdsState.nonce = resp.Nonce

	resources := &DeltaResources{
		NewResources:    make([]*ProtoAny, 0, 1),
		RemovedResource: make([]string, 0, 1),
	}
	logger.Infof("get xDS message nonce, %s", resp.Nonce)
	for _, res := range resp.RemovedResources {
		logger.Infof("remove resource found ", res)
		resources.RemovedResource = append(resources.RemovedResource, res)
	}

	for _, res := range resp.Resources {
		logger.Infof("new resource found %s version=%s", res.Name, res.Version)
		g.xdsState.deltaVersion[res.Name] = res.Version
		elems, err := g.decodeSource(res.Resource)
		if err != nil {
			logger.Infof("can not decode source %s version=%s", res.Name, res.Version, err)
		}
		resources.NewResources = append(resources.NewResources, elems)
	}
	//notify the resource change handler
	output <- resources
}

func (g *GrpcApiClient) subscribeOnGoingChang(delta extensionpb.ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient) error {
	err := delta.Send(&discoverypb.DeltaDiscoveryRequest{
		Node:                    g.makeNode(),
		TypeUrl:                 resource.ExtensionConfigType,
		InitialResourceVersions: g.xdsState.deltaVersion,
		ResponseNonce:           g.xdsState.nonce,
	})
	return err
}

func (g *GrpcApiClient) sendInitDeltaRequest() (extensionpb.ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient, error) {
	delta, err := g.xDSExtensionClient.DeltaExtensionConfigs(context.Background())
	if err != nil {
		return nil, errors.Wrapf(err, "can not start delta stream with xds server ")
	}
	err = delta.Send(&discoverypb.DeltaDiscoveryRequest{
		Node:                     g.makeNode(),
		TypeUrl:                  resource.ExtensionConfigType,
		ResourceNamesSubscribe:   g.resourceNames,
		ResourceNamesUnsubscribe: nil,
		InitialResourceVersions:  g.xdsState.deltaVersion,
		ResponseNonce:            g.xdsState.nonce,
		ErrorDetail:              nil,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "can not send delta discovery request")
	}
	return delta, nil
}

func (g *GrpcApiClient) runStream(output chan<- *DeltaResources) error {
	delta, err := g.xDSExtensionClient.StreamExtensionConfigs(context.Background())
	if err != nil {
		return errors.Wrapf(err, "can not start delta stream with xds server ")
	}
	err = delta.Send(&discoverypb.DiscoveryRequest{
		VersionInfo:   "", //todo load local version
		Node:          g.makeNode(),
		ResourceNames: g.resourceNames,
		TypeUrl:       resource.ExtensionConfigType, //"type.googleapis.com/pixiu.config.listener.v3.Listener", //resource.ListenerType,
	})
	if err != nil {
		return errors.Wrapf(err, "can not send delta discovery request")
	}

	ch := make(chan *discoverypb.DiscoveryResponse)
	//get message
	go func() {
		for {
			resp, err := delta.Recv()
			if err != nil { //todo backoff retry
				logger.Error("can not recv delta discovery request", err)
				break
			}
			ch <- resp
		}
	}()

	go func() {
		defer func() {
			close(output)
		}()
	LOOP:
		for {
			select {
			case <-g.exitCh:
				logger.Infof("stop recv delta of xds (%s)", g.resourceNames)
				if err := delta.CloseSend(); err != nil {
					logger.Errorf("close Send error ", err)
				}
				break LOOP
			case resp := <-ch:
				resources := &DeltaResources{
					NewResources:    make([]*ProtoAny, 0, 1),
					RemovedResource: make([]string, 0, 1),
				}
				logger.Infof("get xDS message nonce, %s", resp.Nonce)

				//for _, res := range resp.RemovedResources {
				//	logger.Infof("remove resource found ", res)
				//	resources.RemovedResource = append(resources.RemovedResource, res)
				//}

				for _, res := range resp.Resources {
					elems, err := g.decodeSource(res)
					if err != nil {
						logger.Infof("can not decode source %s", res, err)
					}
					resources.NewResources = append(resources.NewResources, elems)
				}
				output <- resources
			}
		}
	}()
	return nil
}

func (g *GrpcApiClient) init() {
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
	g.xDSExtensionClient = extensionpb.NewExtensionConfigDiscoveryServiceClient(cluster.GetConnect())
}
