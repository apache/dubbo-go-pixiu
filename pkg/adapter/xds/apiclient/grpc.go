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
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"dubbo.apache.org/dubbo-go/v3/common/logger"
	envoy_config_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

// agent name to talk with xDS server
const xdsAgentName = "dubbo-go-pixiu"

type (
	GrpcApiClient struct {
		config             model.ApiConfigSource
		grpcMg             GrpcClusterManager
		node               *model.Node
		xDSExtensionClient extensionpb.ExtensionConfigDiscoveryServiceClient
		lastExtension      *envoy_config_core_v3.TypedExtensionConfig
		resourceNames      []ResourceTypeName
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
	typeNames ...ResourceTypeName) *GrpcApiClient {
	v := &GrpcApiClient{
		config:        *config,
		node:          node,
		resourceNames: typeNames,
		grpcMg:        grpcMg,
	}
	v.init()
	return v
}

// Fetch get config data from discovery service and return Any type.
func (g *GrpcApiClient) Fetch(localVersion string) ([]*ProtoAny, error) {
	clsRsp, err := g.xDSExtensionClient.FetchExtensionConfigs(context.Background(), &discoverypb.DiscoveryRequest{
		VersionInfo: localVersion,
		Node: &envoy_config_core_v3.Node{
			Id:            g.node.Id,
			Cluster:       g.node.Cluster,
			UserAgentName: xdsAgentName,
		},
		ResourceNames: g.resourceNames,
		TypeUrl:       resource.ExtensionConfigType, //"type.googleapis.com/pixiu.config.listener.v3.Listener", //resource.ListenerType,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "fetch dynamic resource from remote error. %s", g.resourceNames)
	}
	extensions := make([]*ProtoAny, 0, len(clsRsp.Resources))
	for _, _resource := range clsRsp.Resources {
		extension := envoy_config_core_v3.TypedExtensionConfig{}
		err := _resource.UnmarshalTo(&extension)
		if err != nil {
			return nil, errors.Wrapf(err, "typed extension as expected.(%s)", g.resourceNames)
		}
		extensions = append(extensions, &ProtoAny{&extension})
	}
	return extensions, nil
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
		logger.Fatalf("get cluster for init error. error=%v", err)
	}
	g.xDSExtensionClient = extensionpb.NewExtensionConfigDiscoveryServiceClient(cluster.GetConnect())
}
