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
	"sync"
	"testing"
)

import (
	"github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/require"

	corepb "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"go.uber.org/mock/gomock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/metadata"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/emptypb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls/mocks"
)

func TestGRPCClusterManager_GetGrpcCluster(t *testing.T) {
	cluster := &model.ClusterConfig{
		Name:    "cluster-1",
		TypeStr: "GRPC",
		Endpoints: []*model.Endpoint{
			{
				Address: model.SocketAddress{
					Address: "localhost",
					Port:    18000,
				},
			},
		},
	}
	type fields struct {
		clusters  *sync.Map
		clusterMg controls.ClusterManager
	}
	type args struct {
		name string
	}
	ctrl := gomock.NewController(t)
	clusterMg := mocks.NewMockClusterManager(ctrl)
	clusterMg.EXPECT().
		HasCluster("cluster-1").AnyTimes().
		Return(true)
	clusterMg.EXPECT().CloneXdsControlStore().DoAndReturn(func() (controls.ClusterStore, error) {
		store := mocks.NewMockClusterStore(ctrl)
		store.EXPECT().Config().Return([]*model.ClusterConfig{cluster})
		return store, nil
	})
	clusterMg.EXPECT().HasCluster("cluster-2").Return(false)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{"test-simple", fields{
			clusters:  &sync.Map{},
			clusterMg: clusterMg,
		}, args{name: "cluster-1"}, nil},
		{"test-not-exist", fields{
			clusters:  &sync.Map{},
			clusterMg: clusterMg,
		}, args{name: "cluster-2"}, ErrClusterNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GRPCClusterManager{
				clusters:  tt.fields.clusters,
				clusterMg: tt.fields.clusterMg,
			}
			got, err := g.GetGrpcCluster(tt.args.name)
			assert := require.New(t)
			if tt.wantErr != nil {
				assert.ErrorIs(err, tt.wantErr)
				return
			}
			assert.NotNil(got)
			//run two times.
			got, err = g.GetGrpcCluster(tt.args.name)
			assert.NoError(err)
			assert.NotNil(got)

		})
	}
}

func TestGRPCCluster_GetConnect(t *testing.T) {
	cluster := &model.ClusterConfig{
		Name:    "cluster-1",
		TypeStr: "GRPC",
		Endpoints: []*model.Endpoint{
			{
				Address: model.SocketAddress{
					Address: "localhost",
					Port:    18000,
				},
			},
		},
	}
	g := GRPCCluster{
		name:   "name",
		config: cluster,
		once:   sync.Once{},
		conn:   nil,
	}

	gconn := &grpc.ClientConn{}
	var state = connectivity.Ready
	patches := gomonkey.ApplyFunc(grpc.DialContext, func(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
		return gconn, nil
	})
	defer patches.Reset()
	patches.ApplyMethod(&grpc.ClientConn{}, "Close", func(_ *grpc.ClientConn) error {
		return nil
	})
	patches.ApplyMethod(&grpc.ClientConn{}, "GetState", func(_ *grpc.ClientConn) connectivity.State {
		return state
	})

	assert := require.New(t)
	conn, err := g.GetConnection()
	assert.NoError(err)
	assert.NotNil(conn)
	assert.NotNil(g.conn)
	assert.True(g.IsAlive())
	err = g.Close()
	assert.NoError(err)
	state = connectivity.Shutdown
	assert.False(g.IsAlive())
}

func TestGrpcExtensionApiClient_HandleDeltaResponse(t *testing.T) {
	typedPayload, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	outerResource, err := anypb.New(&corepb.TypedExtensionConfig{
		Name:        "resource-a",
		TypedConfig: typedPayload,
	})
	require.NoError(t, err)

	client := &GrpcExtensionApiClient{}
	resources, err := client.handleDeltaResponse(&discoverypb.DeltaDiscoveryResponse{
		Nonce:            "nonce-1",
		RemovedResources: []string{"resource-old"},
		Resources: []*discoverypb.Resource{{
			Name:     "resource-a",
			Version:  "v2",
			Resource: outerResource,
		}},
	})
	require.NoError(t, err)
	require.Len(t, resources.NewResources, 1)
	require.Equal(t, []string{"resource-old"}, resources.RemovedResources)

	state := &xdsState{deltaVersion: map[string]string{"resource-old": "v1"}}
	acceptDeltaResources(state, resources)
	require.Equal(t, map[string]string{"resource-a": "v2"}, state.deltaVersion)

	_, err = client.handleDeltaResponse(&discoverypb.DeltaDiscoveryResponse{
		Resources: []*discoverypb.Resource{{
			Name:     "broken",
			Version:  "v3",
			Resource: &anypb.Any{TypeUrl: "invalid", Value: []byte{0xff}},
		}},
	})
	require.Error(t, err)
}

func TestGrpcExtensionApiClient_RespondToDelta(t *testing.T) {
	client := &GrpcExtensionApiClient{node: &model.Node{Id: "node-1"}}
	stream := &recordingDeltaStream{ctx: context.Background()}

	require.NoError(t, client.respondToDelta(stream, "ack-nonce", nil))
	require.NoError(t, client.respondToDelta(stream, "nack-nonce", stderr.New("apply failed")))
	require.Len(t, stream.sent, 2)

	ack := stream.sent[0]
	require.Equal(t, "ack-nonce", ack.ResponseNonce)
	require.Nil(t, ack.ErrorDetail)
	require.Empty(t, ack.InitialResourceVersions)

	nack := stream.sent[1]
	require.Equal(t, "nack-nonce", nack.ResponseNonce)
	require.Equal(t, int32(codes.InvalidArgument), nack.ErrorDetail.Code)
	require.Contains(t, nack.ErrorDetail.Message, "apply failed")
}

func TestGrpcExtensionApiClient_ConsumeWaitsForApplyBeforeAck(t *testing.T) {
	typedPayload, err := anypb.New(&emptypb.Empty{})
	require.NoError(t, err)
	outerResource, err := anypb.New(&corepb.TypedExtensionConfig{
		Name:        "resource-a",
		TypedConfig: typedPayload,
	})
	require.NoError(t, err)

	stream := &recordingDeltaStream{
		ctx: context.Background(),
		responses: []*discoverypb.DeltaDiscoveryResponse{{
			Nonce: "nonce-1",
			Resources: []*discoverypb.Resource{{
				Name:     "resource-a",
				Version:  "v2",
				Resource: outerResource,
			}},
		}},
	}
	client := &GrpcExtensionApiClient{node: &model.Node{Id: "node-1"}}
	state := &xdsState{deltaVersion: map[string]string{"resource-a": "v1"}}
	output := make(chan *DeltaResources, 1)
	done := make(chan error, 1)
	go func() {
		done <- client.consumeDeltaStream(context.Background(), stream, state, output)
	}()

	resources := <-output
	require.Empty(t, stream.sent, "response must not be ACKed before it is applied")
	require.Equal(t, "v1", state.deltaVersion["resource-a"])

	resources.Complete(nil)
	require.ErrorIs(t, <-done, io.EOF)
	require.Len(t, stream.sent, 1)
	require.Equal(t, "nonce-1", stream.sent[0].ResponseNonce)
	require.Nil(t, stream.sent[0].ErrorDetail)
	require.Equal(t, "v2", state.deltaVersion["resource-a"])
}

func TestDeltaResources_CompleteOnce(t *testing.T) {
	resources := newDeltaResources()
	resources.Complete(nil)
	resources.Complete(stderr.New("late result"))
	require.NoError(t, <-resources.applyResult)
}

func TestProtoAny_ToReturnsDecodeError(t *testing.T) {
	resource := NewProtoAny(&corepb.TypedExtensionConfig{
		Name: "broken",
		TypedConfig: &anypb.Any{
			TypeUrl: "type.googleapis.com/google.protobuf.Empty",
			Value:   []byte{0xff},
		},
	})

	require.NotPanics(t, func() {
		require.Error(t, resource.To(&emptypb.Empty{}))
	})
}

type recordingDeltaStream struct {
	ctx       context.Context
	sent      []*discoverypb.DeltaDiscoveryRequest
	responses []*discoverypb.DeltaDiscoveryResponse
	recvIndex int
}

func (s *recordingDeltaStream) Send(req *discoverypb.DeltaDiscoveryRequest) error {
	s.sent = append(s.sent, req)
	return nil
}

func (s *recordingDeltaStream) Recv() (*discoverypb.DeltaDiscoveryResponse, error) {
	if s.recvIndex >= len(s.responses) {
		return nil, io.EOF
	}
	resp := s.responses[s.recvIndex]
	s.recvIndex++
	return resp, nil
}

func (s *recordingDeltaStream) Header() (metadata.MD, error) { return nil, nil }
func (s *recordingDeltaStream) Trailer() metadata.MD         { return nil }
func (s *recordingDeltaStream) CloseSend() error             { return nil }
func (s *recordingDeltaStream) Context() context.Context     { return s.ctx }
func (s *recordingDeltaStream) SendMsg(any) error            { return nil }
func (s *recordingDeltaStream) RecvMsg(any) error            { return io.EOF }

var _ extensionpb.ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient = (*recordingDeltaStream)(nil)
