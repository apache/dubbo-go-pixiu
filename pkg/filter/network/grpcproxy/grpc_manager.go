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

package grpcproxy

import (
	"context"
)
import (
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	grpcCtx "github.com/apache/dubbo-go-pixiu/pkg/context/grpc"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// GrpcProxyConnectionManager network filter for gRPC proxy, similar to DubboProxyConnectionManager
type GrpcProxyConnectionManager struct {
	filter.EmptyNetworkFilter
	config            *model.GRPCConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
	filterManager     *GrpcFilterManager
	clientConnPool    map[string]*grpc.ClientConn
}

// CreateGrpcProxyConnectionManager create gRPC proxy connection manager
func CreateGrpcProxyConnectionManager(config *model.GRPCConnectionManagerConfig) *GrpcProxyConnectionManager {
	filterManager := NewGrpcFilterManager(config.GrpcFilters)
	gcm := &GrpcProxyConnectionManager{
		config:         config,
		filterManager:  filterManager,
		clientConnPool: make(map[string]*grpc.ClientConn),
	}
	gcm.routerCoordinator = router2.CreateRouterCoordinator(&config.RouteConfig)
	return gcm
}

// determineStreamType determines the StreamType based on client and server streaming flags
func determineStreamType(isClientStream, isServerStream bool) grpcCtx.StreamType {
	if isClientStream && isServerStream {
		return grpcCtx.BidirectionalStream
	} else if isClientStream {
		return grpcCtx.ClientStream
	} else if isServerStream {
		return grpcCtx.ServerStream
	}
	return grpcCtx.UnaryCall
}

// OnUnaryRPC handles a unary RPC call.
func (gcm *GrpcProxyConnectionManager) OnUnaryRPC(ctx context.Context, fullMethod string, req any) (interface{}, error) {
	// Create gRPC context
	grpcCtx := &grpcCtx.GrpcContext{
		Context:     ctx,
		MethodName:  fullMethod,
		Arguments:   []any{req}, // Unary call has only one request parameter
		StreamType:  grpcCtx.UnaryCall,
		IsStreaming: false,
	}

	// Extract service name
	serviceName := gcm.extractServiceName(ctx, fullMethod)
	grpcCtx.ServiceName = serviceName

	// Set metadata
	gcm.extractAndSetMetadata(ctx, grpcCtx)

	// Route to backend service
	if err := gcm.routeRequest(grpcCtx, serviceName, fullMethod); err != nil {
		return nil, err
	}

	// Process request through filter chain
	gcm.handleGrpcInvocation(grpcCtx)

	if grpcCtx.Error != nil {
		return nil, grpcCtx.Error
	}

	return grpcCtx.Result, nil
}

// OnStreamRPC handles a streaming RPC call.
func (gcm *GrpcProxyConnectionManager) OnStreamRPC(stream model.RPCStream, info *model.RPCStreamInfo) error {
	ctx := stream.Context()
	fullMethod := info.FullMethod

	// Create gRPC context
	grpcCtx := &grpcCtx.GrpcContext{
		Context:        ctx,
		MethodName:     fullMethod,
		IsStream:       true,
		IsClientStream: info.IsClientStream,
		IsServerStream: info.IsServerStream,
		Stream:         stream,
		StreamType:     determineStreamType(info.IsClientStream, info.IsServerStream),
		IsStreaming:    true,
	}

	// Extract service name
	serviceName := gcm.extractServiceName(ctx, fullMethod)
	grpcCtx.ServiceName = serviceName

	// Set metadata
	gcm.extractAndSetMetadata(ctx, grpcCtx)

	// Route to backend service
	if err := gcm.routeRequest(grpcCtx, serviceName, fullMethod); err != nil {
		return err
	}

	// Process request through filter chain
	gcm.handleGrpcInvocation(grpcCtx)

	return grpcCtx.Error
}

// extractAndSetMetadata extracts and sets gRPC metadata
func (gcm *GrpcProxyConnectionManager) extractAndSetMetadata(ctx context.Context, grpcContext *grpcCtx.GrpcContext) {
	grpcAttachment := make(map[string]any)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for k, v := range md {
			if len(v) > 0 {
				grpcAttachment[k] = v[0]
			}
		}
	}
	grpcContext.Attachments = grpcAttachment
}

// routeRequest routes the request to backend service
func (gcm *GrpcProxyConnectionManager) routeRequest(grpcContext *grpcCtx.GrpcContext, serviceName, methodName string) error {
	ra, err := gcm.routerCoordinator.RouteByPathAndName(serviceName, methodName)
	if err != nil {
		return errors.Errorf("gRPC route not found: %s/%s", serviceName, methodName)
	}

	grpcContext.Route = ra
	logger.Debugf("[dubbo-go-pixiu] gRPC choose endpoint from cluster: %v", ra.Cluster)
	return nil
}

// handleGrpcInvocation handle gRPC request through filter chain
func (gcm *GrpcProxyConnectionManager) handleGrpcInvocation(ctx *grpcCtx.GrpcContext) {
	filterChain := gcm.filterManager.filters

	// recover any err when filterChain run
	defer func() {
		if err := recover(); err != nil {
			logger.Warnf("[dubbo-go-pixiu] gRPC filter chain panic: %+v", err)
			ctx.SetError(errors.Errorf("gRPC filter chain panic: %v", err))
		}
	}()

	for _, f := range filterChain {
		status := f.Handle(ctx)
		switch status {
		case filter.Continue:
			continue
		case filter.Stop:
			return
		}
	}
}

// extractServiceName extract service name from context or method name
func (gcm *GrpcProxyConnectionManager) extractServiceName(ctx context.Context, methodName string) string {
	// Try to get service name from metadata
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if service := md.Get("grpc-service"); len(service) > 0 {
			return service[0]
		}
	}

	// Parse service name from method name, format is usually /package.Service/Method
	if len(methodName) > 0 && methodName[0] == '/' {
		methodName = methodName[1:]
	}

	lastSlash := -1
	for i := len(methodName) - 1; i >= 0; i-- {
		if methodName[i] == '/' {
			lastSlash = i
			break
		}
	}

	if lastSlash >= 0 {
		return methodName[:lastSlash]
	}

	return "unknown.service"
}
