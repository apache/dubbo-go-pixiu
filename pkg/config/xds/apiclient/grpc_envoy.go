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
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"strconv"
	"time"
)

import (
	clusterpb "github.com/envoyproxy/go-control-plane/envoy/config/cluster/v3"
	envoyconfigcorev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	endpointpb "github.com/envoyproxy/go-control-plane/envoy/config/endpoint/v3"
	discoverypb "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	"github.com/pkg/errors"

	"google.golang.org/grpc/codes"
	grpcstatus "google.golang.org/grpc/status"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	xdsmodel "github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
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
	case constant.ListenerType:
		v.typeUrl = resource.ListenerType
	case constant.ClusterType:
		v.typeUrl = resource.ClusterType
	case constant.EndpointType:
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
	adsResourceState struct {
		versionInfo   string
		resourceNames []string
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
	Clusters  []*clusterpb.Cluster
}

func (g *AggGrpcApiClient) pipeline(output chan *DeltaResources) error {
	if g.typeUrl != resource.ClusterType {
		close(output)
		return errors.Errorf("standard ADS currently supports CDS/EDS, got %s", g.typeUrl)
	}

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
	states := map[string]*adsResourceState{
		resource.ClusterType:  {resourceNames: append([]string(nil), g.resourceNames...)},
		resource.EndpointType: {},
	}
	references := make(map[string]refEndpoint)
	stream, err := g.connectADS(ctx, states, references)
	if err != nil {
		cancel()
		close(output)
		if errors.Is(err, context.Canceled) {
			return nil
		}
		return err
	}

	go func() {
		defer cancel()
		defer close(output)
		for {
			if err := g.consumeADSStream(ctx, stream, states, references, output); err != nil && !errors.Is(err, context.Canceled) {
				logger.Errorf("ADS stream closed: %v", err)
			}
			if ctx.Err() != nil {
				return
			}
			stream, err = g.connectADS(ctx, states, references)
			if err != nil {
				return
			}
		}
	}()
	return nil
}

func (g *AggGrpcApiClient) connectADS(ctx context.Context, states map[string]*adsResourceState, references map[string]refEndpoint) (discoverypb.AggregatedDiscoveryService_StreamAggregatedResourcesClient, error) {
	for {
		stream, err := g.xDSAggClient.StreamAggregatedResources(ctx)
		if err == nil {
			if err = stream.Send(g.makeADSRequest(resource.ClusterType, states[resource.ClusterType], "", nil)); err == nil {
				if len(references) == 0 {
					return stream, nil
				}
				states[resource.EndpointType].resourceNames = referenceNames(references)
				if err = stream.Send(g.makeADSRequest(resource.EndpointType, states[resource.EndpointType], "", nil)); err == nil {
					return stream, nil
				}
			}
			_ = stream.CloseSend()
		}
		logger.Errorf("can not start ADS stream; retrying in 1s: %v", err)
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (g *AggGrpcApiClient) consumeADSStream(ctx context.Context, stream discoverypb.AggregatedDiscoveryService_StreamAggregatedResourcesClient, states map[string]*adsResourceState, references map[string]refEndpoint, output chan<- *DeltaResources) error {
	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}
		if resp == nil {
			return errors.New("ADS response is nil")
		}
		switch resp.TypeUrl {
		case resource.ClusterType:
			newReferences, applyErr := g.decodeCDSReferences(resp.Resources)
			if applyErr == nil && len(newReferences) == 0 {
				applyErr = publishADSClusters(ctx, output, &xdsmodel.PixiuExtensionClusters{})
			}
			if applyErr == nil {
				clear(references)
				for name, ref := range newReferences {
					references[name] = ref
				}
				states[resource.ClusterType].versionInfo = resp.VersionInfo
			}
			if err := stream.Send(g.makeADSRequest(resource.ClusterType, states[resource.ClusterType], resp.Nonce, applyErr)); err != nil {
				return err
			}
			if applyErr != nil || len(references) == 0 {
				continue
			}
			states[resource.EndpointType].resourceNames = referenceNames(references)
			if err := stream.Send(g.makeADSRequest(resource.EndpointType, states[resource.EndpointType], "", nil)); err != nil {
				return err
			}

		case resource.EndpointType:
			clusters, applyErr := decodeAndConvertEDS(resp.Resources, references)
			if applyErr == nil {
				applyErr = publishADSClusters(ctx, output, clusters)
			}
			if applyErr == nil {
				states[resource.EndpointType].versionInfo = resp.VersionInfo
			}
			if err := stream.Send(g.makeADSRequest(resource.EndpointType, states[resource.EndpointType], resp.Nonce, applyErr)); err != nil {
				return err
			}

		default:
			return errors.Errorf("unexpected ADS response type %q", resp.TypeUrl)
		}
	}
}

func (g *AggGrpcApiClient) decodeCDSReferences(resources []*anypb.Any) (map[string]refEndpoint, error) {
	nested := map[resource.Type]map[string]refEndpoint{resource.ClusterType: {}}
	for index, raw := range resources {
		cluster := &clusterpb.Cluster{}
		if raw == nil {
			return nil, errors.Errorf("CDS resource %d is nil", index)
		}
		if err := raw.UnmarshalTo(cluster); err != nil {
			return nil, errors.Wrapf(err, "can not decode CDS resource %d", index)
		}
		g.getClusterResourceReference(cluster, nested)
	}
	return nested[resource.ClusterType], nil
}

func decodeAndConvertEDS(resources []*anypb.Any, references map[string]refEndpoint) (*xdsmodel.PixiuExtensionClusters, error) {
	assignments := make([]*endpointpb.ClusterLoadAssignment, 0, len(resources))
	for _, raw := range resources {
		assignment := &endpointpb.ClusterLoadAssignment{}
		if raw == nil {
			return nil, errors.New("EDS resource is nil")
		}
		if err := raw.UnmarshalTo(assignment); err != nil {
			return nil, errors.Wrap(err, "can not decode EDS resource")
		}
		assignments = append(assignments, assignment)
	}
	return convertClusterLoadAssignments(assignments, references)
}

func publishADSClusters(ctx context.Context, output chan<- *DeltaResources, clusters *xdsmodel.PixiuExtensionClusters) error {
	payload, err := anypb.New(clusters)
	if err != nil {
		return errors.Wrap(err, "encode converted ADS clusters")
	}
	update := newDeltaResources()
	update.NewResources = append(update.NewResources, NewProtoAny(&envoyconfigcorev3.TypedExtensionConfig{
		Name:        constant.ClusterType,
		TypedConfig: payload,
	}))
	select {
	case output <- update:
	case <-ctx.Done():
		return ctx.Err()
	}
	select {
	case err = <-update.applyResult:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (g *AggGrpcApiClient) makeADSRequest(typeURL string, state *adsResourceState, nonce string, responseErr error) *discoverypb.DiscoveryRequest {
	req := &discoverypb.DiscoveryRequest{
		Node:          g.makeNode(),
		VersionInfo:   state.versionInfo,
		ResourceNames: append([]string(nil), state.resourceNames...),
		TypeUrl:       typeURL,
		ResponseNonce: nonce,
	}
	if responseErr != nil {
		req.ErrorDetail = grpcstatus.New(codes.InvalidArgument, responseErr.Error()).Proto()
	}
	return req
}

func referenceNames(references map[string]refEndpoint) []string {
	names := make([]string, 0, len(references))
	for name := range references {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
}

// readServiceNameOfCluster get service name of k8s
func (g *AggGrpcApiClient) readServiceNameOfCluster(c *clusterpb.Cluster) string {
	if c == nil || c.Metadata == nil {
		return ""
	}
	istio := c.Metadata.FilterMetadata["istio"]
	if istio == nil {
		return ""
	}
	services := istio.Fields["services"].GetListValue()
	if services == nil || len(services.Values) == 0 {
		return ""
	}
	service := services.Values[0].GetStructValue()
	if service == nil {
		return ""
	}
	return service.Fields["name"].GetStringValue()
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
	if c == nil {
		return
	}
	logger.Infof("cluster name ==>%s", c.Name)
	serviceName := g.readServiceNameOfCluster(c)
	_, serviceSelected := g.dubboServiceFilter[serviceName]
	_, clusterSelected := g.dubboServiceFilter[c.Name]
	if len(g.dubboServiceFilter) > 0 && !serviceSelected && !clusterSelected {
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

			ref := edsResources[resource.ClusterType][name]
			ref.IsPending = true
			ref.Clusters = append(ref.Clusters, c)
			edsResources[resource.ClusterType][name] = ref
		} else {
			logger.Infof("cluster type %s not supported", typ.Type.String())
		}
	}
}

const (
	endpointHealthMetadataKey  = "pixiu.io/unhealthy"
	endpointWeightMetadataKey  = "envoy.lb/weight"
	endpointRegionMetadataKey  = "envoy.locality/region"
	endpointZoneMetadataKey    = "envoy.locality/zone"
	endpointSubZoneMetadataKey = "envoy.locality/sub_zone"
)

func convertClusterLoadAssignments(assignments []*endpointpb.ClusterLoadAssignment, references map[string]refEndpoint) (*xdsmodel.PixiuExtensionClusters, error) {
	result := &xdsmodel.PixiuExtensionClusters{Clusters: make([]*xdsmodel.Cluster, 0, len(assignments))}
	for _, assignment := range assignments {
		if assignment == nil || assignment.ClusterName == "" {
			return nil, errors.New("EDS assignment must have a cluster name")
		}
		ref, ok := references[assignment.ClusterName]
		if !ok || len(ref.Clusters) == 0 {
			continue
		}

		endpoints, err := convertLoadBalancingEndpoints(assignment)
		if err != nil {
			return nil, errors.Wrapf(err, "convert EDS assignment %q", assignment.ClusterName)
		}
		for _, cluster := range ref.Clusters {
			lb, err := convertLoadBalancingPolicy(cluster.GetLbPolicy())
			if err != nil {
				return nil, errors.Wrapf(err, "convert cluster %q", cluster.GetName())
			}
			result.Clusters = append(result.Clusters, &xdsmodel.Cluster{
				Name:    cluster.GetName(),
				TypeStr: "Static",
				LbStr:   lb,
				EdsClusterConfig: &xdsmodel.EdsClusterConfig{
					ServiceName: assignment.ClusterName,
				},
				Endpoints: cloneXDSEndpoints(endpoints),
			})
		}
	}
	return result, nil
}

func convertLoadBalancingEndpoints(assignment *endpointpb.ClusterLoadAssignment) ([]*xdsmodel.Endpoint, error) {
	var result []*xdsmodel.Endpoint
	for _, localityEndpoints := range assignment.Endpoints {
		if localityEndpoints == nil {
			continue
		}
		for _, lbEndpoint := range localityEndpoints.LbEndpoints {
			if lbEndpoint == nil {
				continue
			}
			endpoint := lbEndpoint.GetEndpoint()
			if endpoint == nil || endpoint.Address == nil || endpoint.Address.GetSocketAddress() == nil {
				return nil, errors.New("only socket-address EDS endpoints are supported")
			}
			address := endpoint.Address.GetSocketAddress()
			if address.Address == "" || address.GetPortValue() == 0 {
				return nil, errors.New("EDS endpoint socket address and port must be set")
			}

			metadata := endpointMetadata(lbEndpoint, localityEndpoints.Locality)
			id := fmt.Sprintf("%s:%d", address.Address, address.GetPortValue())
			result = append(result, &xdsmodel.Endpoint{
				Id:   id,
				Name: id,
				Address: &xdsmodel.SocketAddress{
					Address: address.Address,
					Port:    int64(address.GetPortValue()),
				},
				Metadata: metadata,
			})
		}
	}
	return result, nil
}

func endpointMetadata(endpoint *endpointpb.LbEndpoint, locality *envoyconfigcorev3.Locality) map[string]string {
	metadata := make(map[string]string)
	if endpoint.LoadBalancingWeight != nil {
		metadata[endpointWeightMetadataKey] = strconv.FormatUint(uint64(endpoint.LoadBalancingWeight.Value), 10)
	}
	switch endpoint.HealthStatus {
	case envoyconfigcorev3.HealthStatus_UNHEALTHY, envoyconfigcorev3.HealthStatus_DRAINING, envoyconfigcorev3.HealthStatus_TIMEOUT:
		metadata[endpointHealthMetadataKey] = "true"
	}
	if locality != nil {
		if locality.Region != "" {
			metadata[endpointRegionMetadataKey] = locality.Region
		}
		if locality.Zone != "" {
			metadata[endpointZoneMetadataKey] = locality.Zone
		}
		if locality.SubZone != "" {
			metadata[endpointSubZoneMetadataKey] = locality.SubZone
		}
	}
	if endpoint.Metadata != nil {
		for filter, values := range endpoint.Metadata.FilterMetadata {
			for key, value := range values.Fields {
				encoded, err := json.Marshal(value.AsInterface())
				if err == nil {
					metadata[filter+"/"+key] = string(encoded)
				}
			}
		}
	}
	return metadata
}

func convertLoadBalancingPolicy(policy clusterpb.Cluster_LbPolicy) (string, error) {
	switch policy {
	case clusterpb.Cluster_ROUND_ROBIN:
		return "RoundRobin", nil
	case clusterpb.Cluster_RANDOM:
		return "Rand", nil
	case clusterpb.Cluster_RING_HASH:
		return "RingHashing", nil
	case clusterpb.Cluster_MAGLEV:
		return "MaglevHashing", nil
	default:
		return "", errors.Errorf("unsupported Envoy load balancing policy %s", policy.String())
	}
}

func cloneXDSEndpoints(endpoints []*xdsmodel.Endpoint) []*xdsmodel.Endpoint {
	result := make([]*xdsmodel.Endpoint, 0, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint == nil {
			result = append(result, nil)
			continue
		}
		result = append(result, proto.Clone(endpoint).(*xdsmodel.Endpoint))
	}
	return result
}
