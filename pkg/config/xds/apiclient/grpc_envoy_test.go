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
	"testing"
)

import (
	clusterpb "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	corepb "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointpb "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	"github.com/stretchr/testify/require"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestAggGrpcApiClient_GetClusterResourceReference(t *testing.T) {
	client := &AggGrpcApiClient{dubboServiceFilter: map[string]struct{}{}}
	references := make(map[resource.Type]map[string]refEndpoint)
	first := testEnvoyEDSCluster(t, "outbound|20880||orders.default.svc.cluster.local", "orders-eds", clusterpb.Cluster_ROUND_ROBIN)
	second := testEnvoyEDSCluster(t, "outbound|20881||orders.default.svc.cluster.local", "orders-eds", clusterpb.Cluster_RANDOM)

	client.getClusterResourceReference(first, references)
	client.getClusterResourceReference(second, references)

	ref := references[resource.ClusterType]["orders-eds"]
	require.True(t, ref.IsPending)
	require.Equal(t, []*clusterpb.Cluster{first, second}, ref.Clusters)
	require.Equal(t, "orders", client.readServiceNameOfCluster(first))
	require.Empty(t, client.readServiceNameOfCluster(&clusterpb.Cluster{}))

	filtered := &AggGrpcApiClient{dubboServiceFilter: map[string]struct{}{"payments": {}}}
	filteredReferences := make(map[resource.Type]map[string]refEndpoint)
	filtered.getClusterResourceReference(first, filteredReferences)
	require.Empty(t, filteredReferences)
}

func TestConvertClusterLoadAssignments(t *testing.T) {
	first := testEnvoyEDSCluster(t, "cluster-a", "orders-eds", clusterpb.Cluster_ROUND_ROBIN)
	second := testEnvoyEDSCluster(t, "cluster-b", "orders-eds", clusterpb.Cluster_RING_HASH)
	references := map[string]refEndpoint{
		"orders-eds": {IsPending: true, Clusters: []*clusterpb.Cluster{first, second}},
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
