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

package grpc

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	stdHttp "net/http"
)

import (
	"github.com/dubbogo/grpc-go/codes"
	"github.com/dubbogo/grpc-go/status"
	"golang.org/x/net/http2"
	"google.golang.org/protobuf/proto"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	router2 "github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/router"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
)

// GrpcConnectionManager network filter for grpc
type GrpcConnectionManager struct {
	filter.EmptyNetworkFilter
	config            *model.GRPCConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
}

// CreateGrpcConnectionManager create grpc connection manager
func CreateGrpcConnectionManager(hcmc *model.GRPCConnectionManagerConfig) *GrpcConnectionManager {
	hcm := &GrpcConnectionManager{config: hcmc}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(&hcmc.RouteConfig)
	return hcm
}

// ServeHTTP handle request and response
func (gcm *GrpcConnectionManager) ServeHTTP(w stdHttp.ResponseWriter, r *stdHttp.Request) {

	ra, err := gcm.routerCoordinator.RouteByPathAndName(r.RequestURI, r.Method)
	if err != nil {
		logger.Infof("GrpcConnectionManager can't find route %v", err)
		gcm.writeStatus(w, status.New(codes.NotFound, fmt.Sprintf("proxy can't find route error = %v", err)))
		return
	}

	logger.Debugf("[dubbo-go-pixiu] client choose endpoint from cluster :%v", ra.Cluster)

	clusterName := ra.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		logger.Infof("GrpcConnectionManager can't find endpoint in cluster")
		gcm.writeStatus(w, status.New(codes.Unknown, "can't find endpoint in cluster"))
		return
	}
	ctx := context.Background()
	// timeout
	ctx, cancel := context.WithTimeout(ctx, gcm.config.Timeout)
	defer cancel()
	newReq := r.Clone(ctx)
	newReq.URL.Scheme = "http"
	newReq.URL.Host = endpoint.Address.GetAddress()

	// todo: need cache?
	forwarder := gcm.newHttpForwarder()
	res, err := forwarder.Forward(newReq)

	if err != nil {
		logger.Infof("GrpcConnectionManager forward request error %v", err)
		gcm.writeStatus(w, status.New(codes.Unknown, fmt.Sprintf("forward error not = %v", err)))
		return
	}

	if err := gcm.response(w, res); err != nil {
		logger.Infof("GrpcConnectionManager response  error %v", err)
	}
}

func (gcm *GrpcConnectionManager) writeStatus(w stdHttp.ResponseWriter, status *status.Status) {
	w.Header().Set("Grpc-Status", fmt.Sprintf("%d", status.Code()))
	w.Header().Set("Grpc-Message", status.Message())
	w.Header().Set("Content-Type", "application/grpc")

	if p := status.Proto(); p != nil && len(p.Details) > 0 {
		stBytes, err := proto.Marshal(p)
		if err != nil {
			logger.Warnf("GrpcConnectionManager writeStatus status proto marshal error: %s", err.Error())
		} else {
			w.Header().Set("Grpc-Status-Details-Bin", base64.RawStdEncoding.EncodeToString(stBytes))
		}
	}

	w.WriteHeader(stdHttp.StatusOK)
}

func (gcm *GrpcConnectionManager) response(w stdHttp.ResponseWriter, res *stdHttp.Response) error {
	defer res.Body.Close()

	copyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	w.Write(bytes)

	for k, vv := range res.Trailer {
		k = stdHttp.TrailerPrefix + k
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	return nil
}

func (gcm *GrpcConnectionManager) newHttpForwarder() *HttpForwarder {
	transport := &http2.Transport{
		DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
		AllowHTTP: true,
	}
	return &HttpForwarder{transport: transport}
}

func copyHeader(dst, src stdHttp.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
