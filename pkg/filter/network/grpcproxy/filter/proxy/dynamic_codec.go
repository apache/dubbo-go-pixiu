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

package proxy

import (
	"fmt"
)

import (
	"google.golang.org/protobuf/encoding/protojson"

	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/protobuf/types/dynamicpb"
)

// DynamicCodec is a gRPC codec that uses protobuf reflection to encode/decode messages.
// This enables the gateway to inspect and manipulate message content.
type DynamicCodec struct {
	methodDesc protoreflect.MethodDescriptor
}

// NewDynamicCodec creates a new DynamicCodec for the given method descriptor
func NewDynamicCodec(methodDesc protoreflect.MethodDescriptor) *DynamicCodec {
	return &DynamicCodec{
		methodDesc: methodDesc,
	}
}

// Marshal encodes a message to bytes
func (c *DynamicCodec) Marshal(v any) ([]byte, error) {
	switch msg := v.(type) {
	case *dynamicpb.Message:
		return proto.Marshal(msg)
	case proto.Message:
		return proto.Marshal(msg)
	case []byte:
		// Already bytes, pass through
		return msg, nil
	case *DynamicMessage:
		return proto.Marshal(msg.Message)
	default:
		return nil, fmt.Errorf("dynamic codec: cannot marshal type %T", v)
	}
}

// Unmarshal decodes bytes into a message
func (c *DynamicCodec) Unmarshal(data []byte, v any) error {
	switch msg := v.(type) {
	case **DynamicMessage:
		// Create dynamic message for input type
		inputType := c.methodDesc.Input()
		dynMsg := dynamicpb.NewMessage(inputType)
		if err := proto.Unmarshal(data, dynMsg); err != nil {
			return fmt.Errorf("dynamic codec: unmarshal error: %w", err)
		}
		*msg = &DynamicMessage{
			Message:    dynMsg,
			Descriptor: inputType,
		}
		return nil
	case *[]byte:
		// Passthrough mode
		*msg = data
		return nil
	case proto.Message:
		return proto.Unmarshal(data, msg)
	default:
		return fmt.Errorf("dynamic codec: cannot unmarshal into type %T", v)
	}
}

// Name returns the codec name
func (c *DynamicCodec) Name() string {
	return "dynamic_proto"
}

// DynamicMessage wraps a dynamicpb.Message with its descriptor
type DynamicMessage struct {
	Message    *dynamicpb.Message
	Descriptor protoreflect.MessageDescriptor
}

// GetField retrieves a field value by name
func (dm *DynamicMessage) GetField(name string) (protoreflect.Value, bool) {
	fd := dm.Descriptor.Fields().ByName(protoreflect.Name(name))
	if fd == nil {
		return protoreflect.Value{}, false
	}
	return dm.Message.Get(fd), true
}

// GetFieldString retrieves a string field value by name
func (dm *DynamicMessage) GetFieldString(name string) (string, bool) {
	val, ok := dm.GetField(name)
	if !ok {
		return "", false
	}
	return val.String(), true
}

// GetFieldInt retrieves an int64 field value by name
func (dm *DynamicMessage) GetFieldInt(name string) (int64, bool) {
	val, ok := dm.GetField(name)
	if !ok {
		return 0, false
	}
	return val.Int(), true
}

// SetField sets a field value by name
func (dm *DynamicMessage) SetField(name string, value protoreflect.Value) bool {
	fd := dm.Descriptor.Fields().ByName(protoreflect.Name(name))
	if fd == nil {
		return false
	}
	dm.Message.Set(fd, value)
	return true
}

// ToBytes serializes the message to bytes
func (dm *DynamicMessage) ToBytes() ([]byte, error) {
	return proto.Marshal(dm.Message)
}

// ToJSON converts the message to JSON (for logging/debugging)
func (dm *DynamicMessage) ToJSON() ([]byte, error) {
	return protojson.Marshal(dm.Message)
}

// ResponseCodec is a codec for handling response messages
type ResponseCodec struct {
	methodDesc protoreflect.MethodDescriptor
}

// NewResponseCodec creates a codec for response messages
func NewResponseCodec(methodDesc protoreflect.MethodDescriptor) *ResponseCodec {
	return &ResponseCodec{
		methodDesc: methodDesc,
	}
}

// Marshal encodes a response message
func (c *ResponseCodec) Marshal(v any) ([]byte, error) {
	switch msg := v.(type) {
	case *dynamicpb.Message:
		return proto.Marshal(msg)
	case proto.Message:
		return proto.Marshal(msg)
	case []byte:
		return msg, nil
	case *DynamicMessage:
		return proto.Marshal(msg.Message)
	default:
		return nil, fmt.Errorf("response codec: cannot marshal type %T", v)
	}
}

// Unmarshal decodes response bytes
func (c *ResponseCodec) Unmarshal(data []byte, v any) error {
	switch msg := v.(type) {
	case **DynamicMessage:
		// Create dynamic message for output type
		outputType := c.methodDesc.Output()
		dynMsg := dynamicpb.NewMessage(outputType)
		if err := proto.Unmarshal(data, dynMsg); err != nil {
			return fmt.Errorf("response codec: unmarshal error: %w", err)
		}
		*msg = &DynamicMessage{
			Message:    dynMsg,
			Descriptor: outputType,
		}
		return nil
	case *[]byte:
		*msg = data
		return nil
	case proto.Message:
		return proto.Unmarshal(data, msg)
	default:
		return fmt.Errorf("response codec: cannot unmarshal into type %T", v)
	}
}

// Name returns the codec name
func (c *ResponseCodec) Name() string {
	return "dynamic_proto"
}

// DecodeRequest decodes a request using the method descriptor
func DecodeRequest(methodDesc protoreflect.MethodDescriptor, data []byte) (*DynamicMessage, error) {
	inputType := methodDesc.Input()
	msg := dynamicpb.NewMessage(inputType)
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("failed to decode request: %w", err)
	}
	return &DynamicMessage{
		Message:    msg,
		Descriptor: inputType,
	}, nil
}

// DecodeResponse decodes a response using the method descriptor
func DecodeResponse(methodDesc protoreflect.MethodDescriptor, data []byte) (*DynamicMessage, error) {
	outputType := methodDesc.Output()
	msg := dynamicpb.NewMessage(outputType)
	if err := proto.Unmarshal(data, msg); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &DynamicMessage{
		Message:    msg,
		Descriptor: outputType,
	}, nil
}

// EncodeMessage encodes a dynamic message to bytes
func EncodeMessage(msg *DynamicMessage) ([]byte, error) {
	return proto.Marshal(msg.Message)
}
