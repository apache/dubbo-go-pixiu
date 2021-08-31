package grpcproxy

import (
	"context"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc/status"
)

func Invoke(ctx context.Context, stub grpcdynamic.Stub, mthDesc *desc.MethodDescriptor, grpcReq proto.Message) (proto.Message, error) {
	var resp proto.Message
	var err error
	// TODO(Kenway): use a better way to use bi-direction stream
	// Bi-direction Stream
	if mthDesc.IsServerStreaming() && mthDesc.IsClientStreaming() {
		resp, err = invokeBiDirectionStream(ctx, stub, mthDesc, grpcReq)
	} else if mthDesc.IsClientStreaming() {
		resp, err = invokeClientStream(ctx, stub, mthDesc, grpcReq)
	} else if mthDesc.IsServerStreaming() {
		resp, err = invokeServerStream(ctx, stub, mthDesc, grpcReq)
	} else {
		resp, err = invokeUnary(ctx, stub, mthDesc, grpcReq)
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

func invokeUnary(ctx context.Context, stub grpcdynamic.Stub, mthDesc *desc.MethodDescriptor, grpcReq proto.Message) (proto.Message, error) {
	return stub.InvokeRpc(ctx, mthDesc, grpcReq)
}
