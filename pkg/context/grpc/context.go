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
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// GrpcContext gRPC request context, similar to dubbo.RpcContext
type GrpcContext struct {
	Context     context.Context
	MethodName  string
	ServiceName string
	Arguments   []any
	Attachments map[string]any
	Route       *model.RouteAction
	Result      any
	Error       error
	StreamType  StreamType
	Stream      model.RPCStream
}

// StreamType defines the type of gRPC stream
type StreamType int

const (
	UnaryCall StreamType = 0 + iota
	ClientStream
	ServerStream
	BidirectionalStream
)

// SetError set error to context
func (gc *GrpcContext) SetError(err error) {
	gc.Error = err
}

// SetResult set result to context
func (gc *GrpcContext) SetResult(result any) {
	gc.Result = result
}

// GetAttachment get attachment by key
func (gc *GrpcContext) GetAttachment(key string) (any, bool) {
	if gc.Attachments == nil {
		return nil, false
	}
	val, exists := gc.Attachments[key]
	return val, exists
}

// SetAttachment set attachment
func (gc *GrpcContext) SetAttachment(key string, value any) {
	if gc.Attachments == nil {
		gc.Attachments = make(map[string]any)
	}
	gc.Attachments[key] = value
}

// SetRoute set route
func (gc *GrpcContext) SetRoute(route *model.RouteAction) {
	gc.Route = route
}

// GetRoute get route
func (gc *GrpcContext) GetRoute() *model.RouteAction {
	return gc.Route
}

// SetContext set golang context
func (gc *GrpcContext) SetContext(ctx context.Context) {
	gc.Context = ctx
}

// GetContext get golang context
func (gc *GrpcContext) GetContext() context.Context {
	return gc.Context
}

// GenerateHash generate hash for cache key
func (gc *GrpcContext) GenerateHash() string {
	return gc.ServiceName + "." + gc.MethodName
}

// IsUnary check if it's unary call
func (gc *GrpcContext) IsUnary() bool {
	return gc.StreamType == UnaryCall
}

// IsClientStreaming check if it's client streaming
func (gc *GrpcContext) IsClientStreaming() bool {
	return gc.StreamType == ClientStream
}

// IsServerStreaming check if it's server streaming
func (gc *GrpcContext) IsServerStreaming() bool {
	return gc.StreamType == ServerStream
}

// IsBidirectionalStreaming check if it's bidirectional streaming
func (gc *GrpcContext) IsBidirectionalStreaming() bool {
	return gc.StreamType == BidirectionalStream
}

// SetStream sets the RPCStream for streaming operations
func (gc *GrpcContext) SetStream(stream model.RPCStream) {
	gc.Stream = stream
}

// GetStream gets the RPCStream for streaming operations
func (gc *GrpcContext) GetStream() model.RPCStream {
	return gc.Stream
}
