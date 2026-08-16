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
	"strconv"
	"time"
)

import (
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
	Clusters  []*clusterpb.Cluster
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
			// A CDS response is authoritative for the watched cluster set. Do not
			// retain endpoint references for clusters omitted from the response.
			edsResources[resource.ClusterType] = make(map[string]refEndpoint)
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
				assignments := make([]*endpointpb.ClusterLoadAssignment, 0, len(any2))
				for _, one := range any2 {
					assignment := &endpointpb.ClusterLoadAssignment{}
					if err := one.UnmarshalTo(assignment); err != nil {
						logger.Warnf("can not decode EDS resource: %v", err)
						return
					}
					assignments = append(assignments, assignment)
				}
				extCluster, err := convertClusterLoadAssignments(assignments, clusterRefEndpoints)
				if err != nil {
					logger.Warnf("can not convert EDS resources: %v", err)
					return
				}

				//make output
				output <- &DeltaResources{
					NewResources: []*ProtoAny{
						{
							typeConfig: &envoyconfigcorev3.TypedExtensionConfig{
								Name: constant.ClusterType,
								TypedConfig: func() *anypb.Any { //make any.Any from extCluster
									a, err := anypb.New(extCluster)
									if err != nil {
										logger.Warnf("can not make anypb.Any %v", err)
										return nil
									}
									return a
								}(),
							},
						},
					},
					RemovedResources: nil,
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
