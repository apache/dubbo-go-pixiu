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
	"strings"
)

import (
	"github.com/pkg/errors"

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
}

// CreateGrpcProxyConnectionManager create gRPC proxy connection manager
func CreateGrpcProxyConnectionManager(config *model.GRPCConnectionManagerConfig) *GrpcProxyConnectionManager {
	filterManager := NewGrpcFilterManager(config.GrpcFilters)
	gcm := &GrpcProxyConnectionManager{
		config:        config,
		filterManager: filterManager,
	}
	gcm.routerCoordinator = router2.CreateRouterCoordinator(&config.RouteConfig)
	return gcm
}

// determineStreamType determines the StreamType based on client and server streaming flags
func determineStreamType(isClientStream, isServerStream bool) grpcCtx.StreamType {
	switch {
	case isClientStream && isServerStream:
		return grpcCtx.BidirectionalStream
	case isClientStream:
		return grpcCtx.ClientStream
	case isServerStream:
		return grpcCtx.ServerStream
	default:
		return grpcCtx.UnaryCall
	}
}

// OnStreamRPC handles a streaming RPC call.
func (gcm *GrpcProxyConnectionManager) OnStreamRPC(stream model.RPCStream, info *model.RPCStreamInfo) error {
	ctx := stream.Context()
	fullMethod := info.FullMethod

	// Create gRPC context
	grpcCtx := &grpcCtx.GrpcContext{
		Context:    ctx,
		Stream:     stream,
		StreamType: determineStreamType(info.IsClientStream, info.IsServerStream),
	}

	// Extract service and method names for context, not for routing.
	serviceName := gcm.extractServiceName(ctx, fullMethod)
	methodName := gcm.extractMethodName(fullMethod)
	grpcCtx.ServiceName = serviceName
	grpcCtx.MethodName = methodName

	// Set metadata
	gcm.extractAndSetMetadata(ctx, grpcCtx)

	// Route to backend service using the full method as path and "POST" as the HTTP method.
	if err := gcm.routeRequest(grpcCtx, fullMethod, "POST"); err != nil {
		return err
	}

	// Process request through filter chain
	gcm.handleGrpcInvocation(grpcCtx)

	return grpcCtx.Error
}

func (gcm *GrpcProxyConnectionManager) Close() error {
	var firstErr error
	filterChain := gcm.filterManager.filters

	for _, f := range filterChain {
		if err := f.Close(); err != nil {
			logger.Warnf("Failed to close gRPC filter: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
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

// routeRequest routes the request to a backend service using path and method.
func (gcm *GrpcProxyConnectionManager) routeRequest(grpcContext *grpcCtx.GrpcContext, path, method string) error {
	ra, err := gcm.routerCoordinator.RouteByPathAndName(path, method)
	if err != nil {
		return errors.Errorf("gRPC route not found for path: %s, method: %s", path, method)
	}

	grpcContext.Route = ra
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

// handleGrpcClose handle gRPC request through filter chain
func (gcm *GrpcProxyConnectionManager) handleGrpcClose() error {
	var firstErr error
	filterChain := gcm.filterManager.filters

	for _, f := range filterChain {
		if err := f.Close(); err != nil {
			logger.Warnf("Failed to close gRPC filter: %v", err)
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

// extractMethodName extracts the method name from the full gRPC method string.
// For example, from "/package.Service/Method", it returns "Method".
func (gcm *GrpcProxyConnectionManager) extractMethodName(fullMethod string) string {
	lastSlash := strings.LastIndex(fullMethod, "/")
	if lastSlash == -1 || lastSlash == len(fullMethod)-1 {
		return fullMethod // Return the original string if format is unexpected
	}
	return fullMethod[lastSlash+1:]
}

// extractServiceName extract service name from context or method name
func (gcm *GrpcProxyConnectionManager) extractServiceName(ctx context.Context, fullMethod string) string {
	// Try to get service name from metadata first
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if service := md.Get("grpc-service"); len(service) > 0 {
			return service[0]
		}
	}

	// Fallback to parsing from the full method string, e.g., "/package.Service/Method" -> "package.Service"
	// Trim leading slash for consistency
	fullMethod = strings.TrimPrefix(fullMethod, "/")

	lastSlash := strings.LastIndex(fullMethod, "/")
	if lastSlash > 0 {
		return fullMethod[:lastSlash]
	}

	// If no slash is found, or it's at the beginning, the format is unexpected.
	return "unknown.service"
}
