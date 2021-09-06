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
	"github.com/golang/protobuf/proto" //nolint
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	perrors "github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Invoke(ctx context.Context, stub grpcdynamic.Stub, mthDesc *desc.MethodDescriptor, grpcReq proto.Message, opts ...grpc.CallOption) (proto.Message, error) {
	var resp proto.Message
	var err error
	// Bi-direction Stream
	if mthDesc.IsServerStreaming() && mthDesc.IsClientStreaming() {
		err = perrors.New("currently not support bi-direction stream")
		// resp, err = invokeBiDirectionStream(ctx, stub, mthDesc, grpcReq)
	} else if mthDesc.IsClientStreaming() {
		err = perrors.New("currently not support client side stream")
		// resp, err = invokeClientStream(ctx, stub, mthDesc, grpcReq)
	} else if mthDesc.IsServerStreaming() {
		err = perrors.New("currently not support server side stream")
		// resp, err = invokeServerStream(ctx, stub, mthDesc, grpcReq)
	} else {
		resp, err = invokeUnary(ctx, stub, mthDesc, grpcReq, opts...)
	}

	return resp, err
}

func invokeBiDirectionStream(ctx context.Context, stub grpcdynamic.Stub, mthDesc *desc.MethodDescriptor, grpcReq proto.Message) (proto.Message, error) {
	stream, err := stub.InvokeRpcBidiStream(ctx, mthDesc)
	if _, ok := status.FromError(err); !ok {
		return nil, err
	}
	defer func(stream *grpcdynamic.BidiStream) {
		_ = stream.CloseSend()
	}(stream)

	err = stream.SendMsg(grpcReq)
	if _, ok := status.FromError(err); !ok {
		return nil, err
	}

	return stream.RecvMsg()
}

func invokeClientStream(ctx context.Context, stub grpcdynamic.Stub, mthDesc *desc.MethodDescriptor, grpcReq proto.Message) (proto.Message, error) {
	stream, err := stub.InvokeRpcClientStream(ctx, mthDesc)
	if _, ok := status.FromError(err); !ok {
		return nil, err
	}

	err = stream.SendMsg(grpcReq)
	if _, ok := status.FromError(err); !ok {
		return nil, err
	}

	return stream.CloseAndReceive()
}

func invokeServerStream(ctx context.Context, stub grpcdynamic.Stub, mthDesc *desc.MethodDescriptor, grpcReq proto.Message) (proto.Message, error) {
	stream, err := stub.InvokeRpcServerStream(ctx, mthDesc, grpcReq)
	if _, ok := status.FromError(err); !ok {
		return nil, err
	}

	return stream.RecvMsg()
}

func invokeUnary(ctx context.Context, stub grpcdynamic.Stub, mthDesc *desc.MethodDescriptor, grpcReq proto.Message, opts ...grpc.CallOption) (proto.Message, error) {
	return stub.InvokeRpc(ctx, mthDesc, grpcReq, opts...)
}
