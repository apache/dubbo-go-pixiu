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

package xds

import (
	"context"
	"net"
	"testing"
	"time"
)

import (
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoveryv3 "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionv3 "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	controllog "github.com/envoyproxy/go-control-plane/pkg/log"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	controlserver "github.com/envoyproxy/go-control-plane/pkg/server/v3"

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	xdsmodel "github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
)

func TestSnapshotServedThroughExtensionConfigDiscovery(t *testing.T) {
	result, err := NewSnapshotBuilder(fakeResourceLoader{clusters: []config.Cluster{{
		Name:    "backend",
		Type:    "Static",
		Address: "127.0.0.1",
		Port:    20880,
		ID:      1,
	}}}).Build("42")
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)
	snapshotCache := cache.NewSnapshotCache(false, cache.IDHash{}, controllog.NewDefaultLogger())
	require.NoError(t, snapshotCache.SetSnapshot(ctx, "gateway-a", result.Snapshot))

	grpcServer := grpc.NewServer()
	extensionv3.RegisterExtensionConfigDiscoveryServiceServer(
		grpcServer,
		controlserver.NewServer(ctx, snapshotCache, nil),
	)
	listener := bufconn.Listen(1024 * 1024)
	go func() {
		_ = grpcServer.Serve(listener)
	}()
	t.Cleanup(func() {
		grpcServer.Stop()
		_ = listener.Close()
	})

	dialContext := func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
	conn, err := grpc.NewClient("bufnet",
		grpc.WithContextDialer(dialContext),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	requestCtx, requestCancel := context.WithTimeout(ctx, 5*time.Second)
	defer requestCancel()
	response, err := extensionv3.NewExtensionConfigDiscoveryServiceClient(conn).FetchExtensionConfigs(
		requestCtx,
		&discoveryv3.DiscoveryRequest{
			Node:          &corev3.Node{Id: "gateway-a"},
			TypeUrl:       resource.ExtensionConfigType,
			ResourceNames: []string{constant.ClusterType},
		},
	)
	require.NoError(t, err)
	require.Equal(t, "42", response.VersionInfo)
	require.Len(t, response.Resources, 1)

	typedConfig := &corev3.TypedExtensionConfig{}
	require.NoError(t, response.Resources[0].UnmarshalTo(typedConfig))
	require.Equal(t, constant.ClusterType, typedConfig.Name)
	clusters := &xdsmodel.PixiuExtensionClusters{}
	require.NoError(t, typedConfig.TypedConfig.UnmarshalTo(clusters))
	require.Len(t, clusters.Clusters, 1)
	require.Equal(t, "backend", clusters.Clusters[0].Name)
}
