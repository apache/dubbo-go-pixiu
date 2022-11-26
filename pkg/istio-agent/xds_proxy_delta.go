// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package istioagent

import (
	"context"
	"time"
)

import (
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	google_rpc "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	any "google.golang.org/protobuf/types/known/anypb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/features"
	istiogrpc "github.com/apache/dubbo-go-pixiu/pilot/pkg/grpc"
	v3 "github.com/apache/dubbo-go-pixiu/pilot/pkg/xds/v3"
	"github.com/apache/dubbo-go-pixiu/pkg/istio-agent/metrics"
	"github.com/apache/dubbo-go-pixiu/pkg/wasm"
)

// sendDeltaRequest is a small wrapper around sending to con.requestsChan. This ensures that we do not
// block forever on
func (con *ProxyConnection) sendDeltaRequest(req *discovery.DeltaDiscoveryRequest) {
	select {
	case con.deltaRequestsChan <- req:
	case <-con.stopChan:
	}
}

// requests from envoy
// for aditya:
// downstream -> envoy (anything "behind" xds proxy)
// upstream -> istiod (in front of xds proxy)?
func (p *XdsProxy) DeltaAggregatedResources(downstream discovery.AggregatedDiscoveryService_DeltaAggregatedResourcesServer) error {
	proxyLog.Debugf("accepted delta xds connection from envoy, forwarding to upstream")

	con := &ProxyConnection{
		upstreamError:      make(chan error, 2), // can be produced by recv and send
		downstreamError:    make(chan error, 2), // can be produced by recv and send
		deltaRequestsChan:  make(chan *discovery.DeltaDiscoveryRequest, 10),
		deltaResponsesChan: make(chan *discovery.DeltaDiscoveryResponse, 10),
		stopChan:           make(chan struct{}),
		downstreamDeltas:   downstream,
	}
	p.RegisterStream(con)
	defer p.UnregisterStream(con)

	// Handle downstream xds
	initialRequestsSent := false
	go func() {
		// Send initial request
		p.connectedMutex.RLock()
		initialRequest := p.initialDeltaRequest
		p.connectedMutex.RUnlock()

		for {
			// From Envoy
			req, err := downstream.Recv()
			if err != nil {
				select {
				case con.downstreamError <- err:
				case <-con.stopChan:
				}
				return
			}
			// forward to istiod
			con.sendDeltaRequest(req)
			if !initialRequestsSent && req.TypeUrl == v3.ListenerType {
				// fire off an initial NDS request
				if _, f := p.handlers[v3.NameTableType]; f {
					con.sendDeltaRequest(&discovery.DeltaDiscoveryRequest{
						TypeUrl: v3.NameTableType,
					})
				}
				// fire off an initial PCDS request
				if _, f := p.handlers[v3.ProxyConfigType]; f {
					con.sendDeltaRequest(&discovery.DeltaDiscoveryRequest{
						TypeUrl: v3.ProxyConfigType,
					})
				}
				// Fire of a configured initial request, if there is one
				if initialRequest != nil {
					con.sendDeltaRequest(initialRequest)
				}
				initialRequestsSent = true
			}
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	upstreamConn, err := p.buildUpstreamConn(ctx)
	if err != nil {
		proxyLog.Errorf("failed to connect to upstream %s: %v", p.istiodAddress, err)
		metrics.IstiodConnectionFailures.Increment()
		return err
	}
	defer upstreamConn.Close()

	xds := discovery.NewAggregatedDiscoveryServiceClient(upstreamConn)
	ctx = metadata.AppendToOutgoingContext(context.Background(), "ClusterID", p.clusterID)
	for k, v := range p.xdsHeaders {
		ctx = metadata.AppendToOutgoingContext(ctx, k, v)
	}
	// We must propagate upstream termination to Envoy. This ensures that we resume the full XDS sequence on new connection
	return p.HandleDeltaUpstream(ctx, con, xds)
}

func (p *XdsProxy) HandleDeltaUpstream(ctx context.Context, con *ProxyConnection, xds discovery.AggregatedDiscoveryServiceClient) error {
	deltaUpstream, err := xds.DeltaAggregatedResources(ctx,
		grpc.MaxCallRecvMsgSize(defaultClientMaxReceiveMessageSize))
	if err != nil {
		// Envoy logs errors again, so no need to log beyond debug level
		proxyLog.Debugf("failed to create delta upstream grpc client: %v", err)
		// Increase metric when xds connection error, for example: forgot to restart ingressgateway or sidecar after changing root CA.
		metrics.IstiodConnectionErrors.Increment()
		return err
	}
	proxyLog.Infof("connected to delta upstream XDS server: %s", p.istiodAddress)
	defer proxyLog.Debugf("disconnected from delta XDS server: %s", p.istiodAddress)

	con.upstreamDeltas = deltaUpstream

	// handle responses from istiod
	go func() {
		for {
			resp, err := deltaUpstream.Recv()
			if err != nil {
				select {
				case con.upstreamError <- err:
				case <-con.stopChan:
				}
				return
			}
			select {
			case con.deltaResponsesChan <- resp:
			case <-con.stopChan:
			}
		}
	}()

	go p.handleUpstreamDeltaRequest(con)
	go p.handleUpstreamDeltaResponse(con)

	// todo wasm load conversion
	for {
		select {
		case err := <-con.upstreamError:
			// error from upstream Istiod.
			if istiogrpc.IsExpectedGRPCError(err) {
				proxyLog.Debugf("upstream terminated with status %v", err)
				metrics.IstiodConnectionCancellations.Increment()
			} else {
				proxyLog.Warnf("upstream terminated with unexpected error %v", err)
				metrics.IstiodConnectionErrors.Increment()
			}
			return err
		case err := <-con.downstreamError:
			// error from downstream Envoy.
			if istiogrpc.IsExpectedGRPCError(err) {
				proxyLog.Debugf("downstream terminated with status %v", err)
				metrics.EnvoyConnectionCancellations.Increment()
			} else {
				proxyLog.Warnf("downstream terminated with unexpected error %v", err)
				metrics.EnvoyConnectionErrors.Increment()
			}
			// On downstream error, we will return. This propagates the error to downstream envoy which will trigger reconnect
			return err
		case <-con.stopChan:
			proxyLog.Debugf("stream stopped")
			return nil
		}
	}
}

func (p *XdsProxy) handleUpstreamDeltaRequest(con *ProxyConnection) {
	defer func() {
		_ = con.upstreamDeltas.CloseSend()
	}()
	for {
		select {
		case req := <-con.deltaRequestsChan:
			proxyLog.Debugf("delta request for type url %s", req.TypeUrl)
			metrics.XdsProxyRequests.Increment()
			if req.TypeUrl == v3.ExtensionConfigurationType {
				p.ecdsLastNonce.Store(req.ResponseNonce)
			}
			if err := sendUpstreamDelta(con.upstreamDeltas, req); err != nil {
				proxyLog.Errorf("upstream send error for type url %s: %v", req.TypeUrl, err)
				con.upstreamError <- err
				return
			}
		case <-con.stopChan:
			return
		}
	}
}

func (p *XdsProxy) handleUpstreamDeltaResponse(con *ProxyConnection) {
	forwardEnvoyCh := make(chan *discovery.DeltaDiscoveryResponse, 1)
	for {
		select {
		case resp := <-con.deltaResponsesChan:
			// TODO: separate upstream response handling from requests sending, which are both time costly
			proxyLog.Debugf("response for type url %s", resp.TypeUrl)
			metrics.XdsProxyResponses.Increment()
			if h, f := p.handlers[resp.TypeUrl]; f {
				if len(resp.Resources) == 0 {
					// Empty response, nothing to do
					// This assumes internal types are always singleton
					break
				}
				err := h(resp.Resources[0].Resource)
				var errorResp *google_rpc.Status
				if err != nil {
					errorResp = &google_rpc.Status{
						Code:    int32(codes.Internal),
						Message: err.Error(),
					}
				}
				// Send ACK/NACK
				con.sendDeltaRequest(&discovery.DeltaDiscoveryRequest{
					TypeUrl:       resp.TypeUrl,
					ResponseNonce: resp.Nonce,
					ErrorDetail:   errorResp,
				})
				continue
			}
			switch resp.TypeUrl {
			case v3.ExtensionConfigurationType:
				if features.WasmRemoteLoadConversion {
					// If Wasm remote load conversion feature is enabled, rewrite and send.
					go p.deltaRewriteAndForward(con, resp, func(resp *discovery.DeltaDiscoveryResponse) {
						// Forward the response using the thread of `handleUpstreamResponse`
						// to prevent concurrent access to forwardToEnvoy
						select {
						case forwardEnvoyCh <- resp:
						case <-con.stopChan:
						}
					})
				} else {
					// Otherwise, forward ECDS resource update directly to Envoy.
					forwardDeltaToEnvoy(con, resp)
				}
			default:
				forwardDeltaToEnvoy(con, resp)
			}
		case resp := <-forwardEnvoyCh:
			forwardDeltaToEnvoy(con, resp)
		case <-con.stopChan:
			return
		}
	}
}

func (p *XdsProxy) deltaRewriteAndForward(con *ProxyConnection, resp *discovery.DeltaDiscoveryResponse, forward func(resp *discovery.DeltaDiscoveryResponse)) {
	resources := make([]*any.Any, 0, len(resp.Resources))
	for i := range resp.Resources {
		resources = append(resources, resp.Resources[i].Resource)
	}
	sendNack := wasm.MaybeConvertWasmExtensionConfig(resources, p.wasmCache)
	if sendNack {
		proxyLog.Debugf("sending NACK for ECDS resources %+v", resp.Resources)
		con.sendDeltaRequest(&discovery.DeltaDiscoveryRequest{
			TypeUrl:       v3.ExtensionConfigurationType,
			ResponseNonce: resp.Nonce,
			ErrorDetail: &google_rpc.Status{
				// TODO(bianpengyuan): make error message more informative.
				Message: "failed to fetch wasm module",
			},
		})
		return
	}
	proxyLog.Debugf("forward ECDS resources %+v", resp.Resources)
	forward(resp)
}

func forwardDeltaToEnvoy(con *ProxyConnection, resp *discovery.DeltaDiscoveryResponse) {
	if err := sendDownstreamDelta(con.downstreamDeltas, resp); err != nil {
		select {
		case con.downstreamError <- err:
			proxyLog.Errorf("downstream send error: %v", err)
		default:
			proxyLog.Debugf("downstream error channel full, but get downstream send error: %v", err)
		}

		return
	}
}

func sendUpstreamDelta(deltaUpstream discovery.AggregatedDiscoveryService_DeltaAggregatedResourcesClient,
	req *discovery.DeltaDiscoveryRequest) error {
	return istiogrpc.Send(deltaUpstream.Context(), func() error { return deltaUpstream.Send(req) })
}

func sendDownstreamDelta(deltaDownstream discovery.AggregatedDiscoveryService_DeltaAggregatedResourcesServer,
	res *discovery.DeltaDiscoveryResponse) error {
	return istiogrpc.Send(deltaDownstream.Context(), func() error { return deltaDownstream.Send(res) })
}

func (p *XdsProxy) PersistDeltaRequest(req *discovery.DeltaDiscoveryRequest) {
	var ch chan *discovery.DeltaDiscoveryRequest
	var stop chan struct{}

	p.connectedMutex.Lock()
	if p.connected != nil {
		ch = p.connected.deltaRequestsChan
		stop = p.connected.stopChan
	}
	p.initialDeltaRequest = req
	p.connectedMutex.Unlock()

	// Immediately send if we are currently connect
	if ch != nil {
		select {
		case ch <- req:
		case <-stop:
		}
	}
}
