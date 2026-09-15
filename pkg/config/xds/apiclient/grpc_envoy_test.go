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
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

import (
	clusterpb "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corepb "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointpb "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	grpcstatus "google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

import (
	xdsmodel "github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
)

func TestAggGrpcApiClient_GetClusterResourceReference(t *testing.T) {
	client := &AggGrpcApiClient{dubboServiceFilter: map[string]struct{}{}}
	references := make(map[string]refEndpoint)
	first := testEnvoyEDSCluster(t, "outbound|20880||orders.default.svc.cluster.local", "orders-eds", clusterpb.Cluster_ROUND_ROBIN)
	second := testEnvoyEDSCluster(t, "outbound|20881||orders.default.svc.cluster.local", "orders-eds", clusterpb.Cluster_RANDOM)

	client.getClusterResourceReference(first, references)
	client.getClusterResourceReference(second, references)

	ref := references["orders-eds"]
	require.Equal(t, []*clusterpb.Cluster{first, second}, ref.Clusters)
	require.Equal(t, "orders", client.readServiceNameOfCluster(first))
	require.Empty(t, client.readServiceNameOfCluster(&clusterpb.Cluster{}))

	filtered := &AggGrpcApiClient{dubboServiceFilter: map[string]struct{}{"payments": {}}}
	filteredReferences := make(map[string]refEndpoint)
	filtered.getClusterResourceReference(first, filteredReferences)
	require.Empty(t, filteredReferences)
}

func TestAggGrpcApiClient_SkipsUnsupportedCDSClusters(t *testing.T) {
	client := &AggGrpcApiClient{dubboServiceFilter: map[string]struct{}{}}
	staticResource, err := anypb.New(&clusterpb.Cluster{
		Name:                 "local-management",
		ClusterDiscoveryType: &clusterpb.Cluster_Type{Type: clusterpb.Cluster_STATIC},
		LbPolicy:             clusterpb.Cluster_ROUND_ROBIN,
	})
	require.NoError(t, err)
	edsResource, err := anypb.New(testEnvoyEDSCluster(t, "orders", "orders-eds", clusterpb.Cluster_ROUND_ROBIN))
	require.NoError(t, err)

	references, err := client.decodeCDSReferences([]*anypb.Any{staticResource, edsResource})
	require.NoError(t, err)
	require.NotContains(t, references, "local-management")
	require.Contains(t, references, "orders-eds")
}

func TestConvertClusterLoadAssignments(t *testing.T) {
	first := testEnvoyEDSCluster(t, "cluster-a", "orders-eds", clusterpb.Cluster_ROUND_ROBIN)
	second := testEnvoyEDSCluster(t, "cluster-b", "orders-eds", clusterpb.Cluster_RING_HASH)
	references := map[string]refEndpoint{
		"orders-eds": {Clusters: []*clusterpb.Cluster{first, second}},
	}
	endpointMeta, err := structpb.NewStruct(map[string]any{"canary": true})
	require.NoError(t, err)
	assignment := &endpointpb.ClusterLoadAssignment{
		ClusterName: "orders-eds",
		Endpoints: []*endpointpb.LocalityLbEndpoints{
			{
				Locality: &corepb.Locality{Region: "cn", Zone: "hangzhou-a"},
				LbEndpoints: []*endpointpb.LbEndpoint{
					testEnvoyEndpoint("10.0.0.1", 20880, corepb.HealthStatus_HEALTHY, 3, endpointMeta),
					testEnvoyEndpoint("10.0.0.3", 20882, corepb.HealthStatus_DEGRADED, 2, nil),
				},
			},
			{
				Locality: &corepb.Locality{Region: "cn", Zone: "hangzhou-b", SubZone: "rack-2"},
				LbEndpoints: []*endpointpb.LbEndpoint{
					testEnvoyEndpoint("10.0.0.2", 20881, corepb.HealthStatus_UNHEALTHY, 1, nil),
				},
			},
		},
	}

	converted, err := convertClusterLoadAssignments([]*endpointpb.ClusterLoadAssignment{assignment}, references)
	require.NoError(t, err)
	require.Len(t, converted.Clusters, 2)
	require.Equal(t, "cluster-a", converted.Clusters[0].Name)
	require.Equal(t, "RoundRobin", converted.Clusters[0].LbStr)
	require.Equal(t, "cluster-b", converted.Clusters[1].Name)
	require.Equal(t, "RingHashing", converted.Clusters[1].LbStr)
	require.Len(t, converted.Clusters[0].Endpoints, 3)
	require.NotSame(t, converted.Clusters[0].Endpoints[0], converted.Clusters[1].Endpoints[0])

	healthy := converted.Clusters[0].Endpoints[0]
	require.Equal(t, "10.0.0.1:20880", healthy.Id)
	require.Equal(t, "3", healthy.Metadata[endpointWeightMetadataKey])
	require.Equal(t, "cn", healthy.Metadata[endpointRegionMetadataKey])
	require.Equal(t, "hangzhou-a", healthy.Metadata[endpointZoneMetadataKey])
	require.Equal(t, "true", healthy.Metadata["envoy.lb/canary"])

	unhealthy := converted.Clusters[0].Endpoints[2]
	require.Equal(t, "true", unhealthy.Metadata[endpointHealthMetadataKey])
	require.Equal(t, "rack-2", unhealthy.Metadata[endpointSubZoneMetadataKey])
}

func TestConvertClusterLoadAssignmentsRejectsUnsupportedData(t *testing.T) {
	references := map[string]refEndpoint{
		"orders-eds": {Clusters: []*clusterpb.Cluster{
			testEnvoyEDSCluster(t, "cluster-a", "orders-eds", clusterpb.Cluster_LEAST_REQUEST),
		}},
	}
	assignment := &endpointpb.ClusterLoadAssignment{
		ClusterName: "orders-eds",
		Endpoints: []*endpointpb.LocalityLbEndpoints{{
			LbEndpoints: []*endpointpb.LbEndpoint{{
				HostIdentifier: &endpointpb.LbEndpoint_EndpointName{EndpointName: "named-only"},
			}},
		}},
	}

	_, err := convertClusterLoadAssignments([]*endpointpb.ClusterLoadAssignment{assignment}, references)
	require.ErrorContains(t, err, "socket-address")

	assignment.Endpoints[0].LbEndpoints[0] = testEnvoyEndpoint("10.0.0.1", 20880, corepb.HealthStatus_HEALTHY, 1, nil)
	_, err = convertClusterLoadAssignments([]*endpointpb.ClusterLoadAssignment{assignment}, references)
	require.ErrorContains(t, err, "LEAST_REQUEST")

	references["orders-eds"].Clusters[0] = testEnvoyEDSCluster(t, "cluster-a", "orders-eds", clusterpb.Cluster_ROUND_ROBIN)
	assignment.Endpoints[0].LbEndpoints[0] = testEnvoyEndpoint("10.0.0.1", 65536, corepb.HealthStatus_HEALTHY, 1, nil)
	_, err = convertClusterLoadAssignments([]*endpointpb.ClusterLoadAssignment{assignment}, references)
	require.ErrorContains(t, err, "invalid socket address")
}

func TestAggGrpcApiClient_ConsumeADSStream(t *testing.T) {
	cluster := testEnvoyEDSCluster(t, "cluster-a", "orders-eds", clusterpb.Cluster_ROUND_ROBIN)
	clusterResource, err := anypb.New(cluster)
	require.NoError(t, err)
	assignment := &endpointpb.ClusterLoadAssignment{
		ClusterName: "orders-eds",
		Endpoints: []*endpointpb.LocalityLbEndpoints{{
			LbEndpoints: []*endpointpb.LbEndpoint{
				testEnvoyEndpoint("10.0.0.1", 20880, corepb.HealthStatus_HEALTHY, 1, nil),
			},
		}},
	}
	endpointResource, err := anypb.New(assignment)
	require.NoError(t, err)

	stream := &recordingADSStream{
		ctx: context.Background(),
		responses: []*discoverypb.DiscoveryResponse{
			{TypeUrl: resource.ClusterType, VersionInfo: "cds-v1", Nonce: "cds-nonce", Resources: []*anypb.Any{clusterResource}},
			{TypeUrl: resource.EndpointType, VersionInfo: "eds-v2", Nonce: "eds-nonce", Resources: []*anypb.Any{endpointResource}},
		},
	}
	client := &AggGrpcApiClient{dubboServiceFilter: map[string]struct{}{}}
	states := map[string]*adsResourceState{
		resource.ClusterType:  {},
		resource.EndpointType: {},
	}
	references := make(map[string]refEndpoint)
	output := make(chan *DeltaResources, 1)
	done := make(chan error, 1)
	go func() {
		done <- client.consumeADSStream(context.Background(), stream, states, references, output)
	}()

	update := <-output
	require.Len(t, stream.sent, 2)
	require.Equal(t, resource.ClusterType, stream.sent[0].TypeUrl)
	require.Equal(t, "cds-v1", stream.sent[0].VersionInfo)
	require.Equal(t, "cds-nonce", stream.sent[0].ResponseNonce)
	require.Equal(t, resource.EndpointType, stream.sent[1].TypeUrl)
	require.Equal(t, []string{"orders-eds"}, stream.sent[1].ResourceNames)
	require.Empty(t, stream.sent[1].ResponseNonce)

	converted := &xdsmodel.PixiuExtensionClusters{}
	require.NoError(t, update.NewResources[0].To(converted))
	require.Len(t, converted.Clusters, 1)
	require.Equal(t, "cluster-a", converted.Clusters[0].Name)
	update.Complete(nil)

	require.ErrorIs(t, <-done, io.EOF)
	require.Len(t, stream.sent, 3)
	require.Equal(t, "eds-v2", stream.sent[2].VersionInfo)
	require.Equal(t, "eds-nonce", stream.sent[2].ResponseNonce)
	require.Equal(t, "cds-v1", states[resource.ClusterType].versionInfo)
	require.Equal(t, "eds-v2", states[resource.EndpointType].versionInfo)
}

func TestAggGrpcApiClient_MakeADSNACK(t *testing.T) {
	client := &AggGrpcApiClient{}
	state := &adsResourceState{versionInfo: "last-good", resourceNames: []string{"orders-eds"}}
	req := client.makeADSRequest(resource.EndpointType, state, "rejected-nonce", stderr.New("bad endpoint"))

	require.Equal(t, "last-good", req.VersionInfo)
	require.Equal(t, "rejected-nonce", req.ResponseNonce)
	require.Equal(t, []string{"orders-eds"}, req.ResourceNames)
	require.Equal(t, int32(codes.InvalidArgument), req.ErrorDetail.Code)
}

func TestAggGrpcApiClient_NACKsInvalidCDSBeforeRecordingVersion(t *testing.T) {
	invalid := testEnvoyEDSCluster(t, "cluster-a", "orders-eds", clusterpb.Cluster_LEAST_REQUEST)
	resourceAny, err := anypb.New(invalid)
	require.NoError(t, err)
	stream := &recordingADSStream{
		ctx: context.Background(),
		responses: []*discoverypb.DiscoveryResponse{{
			TypeUrl:     resource.ClusterType,
			VersionInfo: "invalid-cds-v2",
			Nonce:       "invalid-cds-nonce",
			Resources:   []*anypb.Any{resourceAny},
		}},
	}
	client := &AggGrpcApiClient{dubboServiceFilter: map[string]struct{}{}}
	states := map[string]*adsResourceState{
		resource.ClusterType:  {versionInfo: "last-good-cds-v1"},
		resource.EndpointType: {},
	}
	output := make(chan *DeltaResources, 1)

	err = client.consumeADSStream(context.Background(), stream, states, make(map[string]refEndpoint), output)
	require.ErrorIs(t, err, io.EOF)
	require.Len(t, stream.sent, 1)
	require.Equal(t, resource.ClusterType, stream.sent[0].TypeUrl)
	require.Equal(t, "last-good-cds-v1", stream.sent[0].VersionInfo)
	require.Equal(t, "invalid-cds-nonce", stream.sent[0].ResponseNonce)
	require.Equal(t, int32(codes.InvalidArgument), stream.sent[0].ErrorDetail.Code)
	require.ErrorContains(t, grpcstatus.ErrorProto(stream.sent[0].ErrorDetail), "LEAST_REQUEST")
	require.Equal(t, "last-good-cds-v1", states[resource.ClusterType].versionInfo)
	select {
	case <-output:
		t.Fatal("invalid CDS must not be published to the cluster manager")
	default:
	}
}

type reconnectADSServer struct {
	discoverypb.UnimplementedAggregatedDiscoveryServiceServer

	mu               sync.Mutex
	streamCount      int
	clusterResource  *anypb.Any
	endpointResource *anypb.Any
	reconnected      chan []*discoverypb.DiscoveryRequest
}

func (s *reconnectADSServer) StreamAggregatedResources(stream discoverypb.AggregatedDiscoveryService_StreamAggregatedResourcesServer) error {
	s.mu.Lock()
	s.streamCount++
	streamNumber := s.streamCount
	s.mu.Unlock()

	if streamNumber == 1 {
		initialCDS, err := stream.Recv()
		if err != nil {
			return err
		}
		if initialCDS.TypeUrl != resource.ClusterType {
			return stderr.New("first request is not CDS")
		}
		if err := stream.Send(&discoverypb.DiscoveryResponse{
			TypeUrl:     resource.ClusterType,
			VersionInfo: "cds-v1",
			Nonce:       "cds-nonce-1",
			Resources:   []*anypb.Any{s.clusterResource},
		}); err != nil {
			return err
		}
		cdsACK, err := stream.Recv()
		if err != nil {
			return err
		}
		edsSubscription, err := stream.Recv()
		if err != nil {
			return err
		}
		if cdsACK.VersionInfo != "cds-v1" || cdsACK.ResponseNonce != "cds-nonce-1" || edsSubscription.TypeUrl != resource.EndpointType {
			return stderr.New("client did not ACK CDS and subscribe to EDS")
		}
		if err := stream.Send(&discoverypb.DiscoveryResponse{
			TypeUrl:     resource.EndpointType,
			VersionInfo: "eds-v1",
			Nonce:       "eds-nonce-1",
			Resources:   []*anypb.Any{s.endpointResource},
		}); err != nil {
			return err
		}
		edsACK, err := stream.Recv()
		if err != nil {
			return err
		}
		if edsACK.VersionInfo != "eds-v1" || edsACK.ResponseNonce != "eds-nonce-1" {
			return stderr.New("client did not ACK EDS")
		}
		return nil
	}

	requests := make([]*discoverypb.DiscoveryRequest, 0, 2)
	for len(requests) < 2 {
		request, err := stream.Recv()
		if err != nil {
			return err
		}
		requests = append(requests, request)
	}
	select {
	case s.reconnected <- requests:
	case <-stream.Context().Done():
		return stream.Context().Err()
	}
	<-stream.Context().Done()
	return nil
}

func TestAggGrpcApiClient_ReconnectsWithLastAcceptedVersions(t *testing.T) {
	clusterResource, err := anypb.New(testEnvoyEDSCluster(t, "cluster-a", "orders-eds", clusterpb.Cluster_ROUND_ROBIN))
	require.NoError(t, err)
	endpointResource, err := anypb.New(&endpointpb.ClusterLoadAssignment{
		ClusterName: "orders-eds",
		Endpoints: []*endpointpb.LocalityLbEndpoints{{
			LbEndpoints: []*endpointpb.LbEndpoint{
				testEnvoyEndpoint("10.0.0.1", 20880, corepb.HealthStatus_HEALTHY, 1, nil),
			},
		}},
	})
	require.NoError(t, err)

	managementServer := &reconnectADSServer{
		clusterResource:  clusterResource,
		endpointResource: endpointResource,
		reconnected:      make(chan []*discoverypb.DiscoveryRequest, 1),
	}
	listener := bufconn.Listen(1024 * 1024)
	grpcServer := grpc.NewServer()
	discoverypb.RegisterAggregatedDiscoveryServiceServer(grpcServer, managementServer)
	go func() { _ = grpcServer.Serve(listener) }()
	t.Cleanup(func() {
		grpcServer.Stop()
		_ = listener.Close()
	})

	// grpc.NewClient defaults to the DNS resolver, which cannot resolve the
	// synthetic "bufnet" target; the passthrough scheme keeps the endpoint
	// verbatim so the custom dialer receives it.
	conn, err := grpc.NewClient(
		"passthrough:///bufnet",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) { return listener.Dial() }),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = conn.Close() })

	exitCh := make(chan struct{})
	client := &AggGrpcApiClient{
		GrpcExtensionApiClient: GrpcExtensionApiClient{
			exitCh:  exitCh,
			typeUrl: resource.ClusterType,
		},
		xDSAggClient:       discoverypb.NewAggregatedDiscoveryServiceClient(conn),
		dubboServiceFilter: map[string]struct{}{},
	}
	updates := make(chan *DeltaResources)
	require.NoError(t, client.pipeline(updates))

	select {
	case update := <-updates:
		update.Complete(nil)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for initial CDS/EDS application")
	}

	select {
	case requests := <-managementServer.reconnected:
		require.Equal(t, resource.ClusterType, requests[0].TypeUrl)
		require.Equal(t, "cds-v1", requests[0].VersionInfo)
		require.Equal(t, resource.EndpointType, requests[1].TypeUrl)
		require.Equal(t, "eds-v1", requests[1].VersionInfo)
		require.Equal(t, []string{"orders-eds"}, requests[1].ResourceNames)
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for ADS reconnect")
	}

	close(exitCh)
	select {
	case _, ok := <-updates:
		require.False(t, ok)
	case <-time.After(5 * time.Second):
		t.Fatal("ADS pipeline did not stop")
	}
}

func testEnvoyEDSCluster(t *testing.T, name string, serviceName string, policy clusterpb.Cluster_LbPolicy) *clusterpb.Cluster {
	t.Helper()
	istio, err := structpb.NewStruct(map[string]any{
		"services": []any{map[string]any{"name": "orders"}},
	})
	require.NoError(t, err)
	return &clusterpb.Cluster{
		Name:                 name,
		ClusterDiscoveryType: &clusterpb.Cluster_Type{Type: clusterpb.Cluster_EDS},
		EdsClusterConfig:     &clusterpb.Cluster_EdsClusterConfig{ServiceName: serviceName},
		LbPolicy:             policy,
		Metadata: &corepb.Metadata{FilterMetadata: map[string]*structpb.Struct{
			"istio": istio,
		}},
	}
}

func testEnvoyEndpoint(address string, port uint32, health corepb.HealthStatus, weight uint32, metadata *structpb.Struct) *endpointpb.LbEndpoint {
	var endpointMetadata *corepb.Metadata
	if metadata != nil {
		endpointMetadata = &corepb.Metadata{FilterMetadata: map[string]*structpb.Struct{"envoy.lb": metadata}}
	}
	return &endpointpb.LbEndpoint{
		HostIdentifier: &endpointpb.LbEndpoint_Endpoint{Endpoint: &endpointpb.Endpoint{
			Address: &corepb.Address{Address: &corepb.Address_SocketAddress{SocketAddress: &corepb.SocketAddress{
				Address:       address,
				PortSpecifier: &corepb.SocketAddress_PortValue{PortValue: port},
			}}},
		}},
		HealthStatus:        health,
		LoadBalancingWeight: wrapperspb.UInt32(weight),
		Metadata:            endpointMetadata,
	}
}

type recordingADSStream struct {
	ctx       context.Context
	sent      []*discoverypb.DiscoveryRequest
	responses []*discoverypb.DiscoveryResponse
	recvIndex int
}

func (s *recordingADSStream) Send(req *discoverypb.DiscoveryRequest) error {
	s.sent = append(s.sent, req)
	return nil
}

func (s *recordingADSStream) Recv() (*discoverypb.DiscoveryResponse, error) {
	if s.recvIndex >= len(s.responses) {
		return nil, io.EOF
	}
	resp := s.responses[s.recvIndex]
	s.recvIndex++
	return resp, nil
}

func (s *recordingADSStream) Header() (metadata.MD, error) { return nil, nil }
func (s *recordingADSStream) Trailer() metadata.MD         { return nil }
func (s *recordingADSStream) CloseSend() error             { return nil }
func (s *recordingADSStream) Context() context.Context     { return s.ctx }
func (s *recordingADSStream) SendMsg(any) error            { return nil }
func (s *recordingADSStream) RecvMsg(any) error            { return io.EOF }

var _ discoverypb.AggregatedDiscoveryService_StreamAggregatedResourcesClient = (*recordingADSStream)(nil)
