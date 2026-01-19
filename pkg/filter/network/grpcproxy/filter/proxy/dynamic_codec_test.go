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
	"strings"
	"testing"
)

import (
	"google.golang.org/protobuf/proto"

	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"

	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
)

// buildTestMethodDescriptor creates a test method descriptor for unit testing
func buildTestMethodDescriptor(t *testing.T) protoreflect.MethodDescriptor {
	// Define a simple proto file descriptor
	fileDescProto := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("test.proto"),
		Package: proto.String("test"),
		MessageType: []*descriptorpb.DescriptorProto{
			{
				Name: proto.String("TestRequest"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   proto.String("name"),
						Number: proto.Int32(1),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					},
					{
						Name:   proto.String("id"),
						Number: proto.Int32(2),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_INT64.Enum(),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					},
				},
			},
			{
				Name: proto.String("TestResponse"),
				Field: []*descriptorpb.FieldDescriptorProto{
					{
						Name:   proto.String("message"),
						Number: proto.Int32(1),
						Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
						Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
					},
				},
			},
		},
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("TestService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:       proto.String("TestMethod"),
						InputType:  proto.String(".test.TestRequest"),
						OutputType: proto.String(".test.TestResponse"),
					},
				},
			},
		},
	}

	// Build file descriptor
	fileDesc, err := protodesc.NewFile(fileDescProto, nil)
	if err != nil {
		t.Fatalf("Failed to create file descriptor: %v", err)
	}

	// Get service and method
	svcDesc := fileDesc.Services().ByName("TestService")
	if svcDesc == nil {
		t.Fatal("TestService not found")
	}

	methodDesc := svcDesc.Methods().ByName("TestMethod")
	if methodDesc == nil {
		t.Fatal("TestMethod not found")
	}

	return methodDesc
}

func TestDynamicCodec_Name(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	codec := NewDynamicCodec(methodDesc)

	if got := codec.Name(); got != "dynamic_proto" {
		t.Errorf("Name() = %q, want %q", got, "dynamic_proto")
	}
}

func TestDynamicCodec_MarshalBytes(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	codec := NewDynamicCodec(methodDesc)

	// Test marshaling raw bytes (passthrough)
	data := []byte{0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
	result, err := codec.Marshal(data)
	if err != nil {
		t.Fatalf("Marshal([]byte) error = %v", err)
	}
	if string(result) != string(data) {
		t.Errorf("Marshal([]byte) = %v, want %v", result, data)
	}
}

func TestDynamicCodec_MarshalDynamicMessage(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	codec := NewDynamicCodec(methodDesc)

	// Create a dynamic message
	inputType := methodDesc.Input()
	dynMsg := dynamicpb.NewMessage(inputType)

	// Set field values
	nameField := inputType.Fields().ByName("name")
	dynMsg.Set(nameField, protoreflect.ValueOfString("test"))

	// Marshal as DynamicMessage wrapper
	dm := &DynamicMessage{
		Message:    dynMsg,
		Descriptor: inputType,
	}

	result, err := codec.Marshal(dm)
	if err != nil {
		t.Fatalf("Marshal(*DynamicMessage) error = %v", err)
	}

	// Verify by unmarshaling
	newMsg := dynamicpb.NewMessage(inputType)
	if err := proto.Unmarshal(result, newMsg); err != nil {
		t.Fatalf("Unmarshal verification error = %v", err)
	}

	if got := newMsg.Get(nameField).String(); got != "test" {
		t.Errorf("Unmarshaled name = %q, want %q", got, "test")
	}
}

func TestDynamicCodec_MarshalUnsupportedType(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	codec := NewDynamicCodec(methodDesc)

	// Try to marshal unsupported type
	_, err := codec.Marshal(123)
	if err == nil {
		t.Error("Marshal(int) should return error")
	}
}

func TestDynamicCodec_UnmarshalToBytes(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	codec := NewDynamicCodec(methodDesc)

	data := []byte{0x0a, 0x05, 0x68, 0x65, 0x6c, 0x6c, 0x6f}
	var result []byte
	err := codec.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal(to []byte) error = %v", err)
	}
	if string(result) != string(data) {
		t.Errorf("Unmarshal result = %v, want %v", result, data)
	}
}

func TestDynamicCodec_UnmarshalToDynamicMessage(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	codec := NewDynamicCodec(methodDesc)

	// Create test data
	inputType := methodDesc.Input()
	srcMsg := dynamicpb.NewMessage(inputType)
	nameField := inputType.Fields().ByName("name")
	idField := inputType.Fields().ByName("id")
	srcMsg.Set(nameField, protoreflect.ValueOfString("hello"))
	srcMsg.Set(idField, protoreflect.ValueOfInt64(42))

	data, err := proto.Marshal(srcMsg)
	if err != nil {
		t.Fatalf("Failed to marshal source message: %v", err)
	}

	// Unmarshal to DynamicMessage
	var dm *DynamicMessage
	err = codec.Unmarshal(data, &dm)
	if err != nil {
		t.Fatalf("Unmarshal(to *DynamicMessage) error = %v", err)
	}

	if dm == nil {
		t.Fatal("Unmarshal result is nil")
	}

	// Verify fields
	if got := dm.Message.Get(nameField).String(); got != "hello" {
		t.Errorf("name = %q, want %q", got, "hello")
	}
	if got := dm.Message.Get(idField).Int(); got != 42 {
		t.Errorf("id = %d, want %d", got, 42)
	}
}

func TestDynamicCodec_UnmarshalUnsupportedType(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	codec := NewDynamicCodec(methodDesc)

	data := []byte{0x0a, 0x05}
	var result int
	err := codec.Unmarshal(data, &result)
	if err == nil {
		t.Error("Unmarshal(to int) should return error")
	}
}

func TestDynamicMessage_GetField(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	inputType := methodDesc.Input()
	dynMsg := dynamicpb.NewMessage(inputType)

	nameField := inputType.Fields().ByName("name")
	dynMsg.Set(nameField, protoreflect.ValueOfString("test_value"))

	dm := &DynamicMessage{
		Message:    dynMsg,
		Descriptor: inputType,
	}

	// Test GetField with existing field
	val, ok := dm.GetField("name")
	if !ok {
		t.Error("GetField(name) returned false")
	}
	if val.String() != "test_value" {
		t.Errorf("GetField(name) = %q, want %q", val.String(), "test_value")
	}

	// Test GetField with non-existing field
	_, ok = dm.GetField("nonexistent")
	if ok {
		t.Error("GetField(nonexistent) should return false")
	}
}

func TestDynamicMessage_GetFieldString(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	inputType := methodDesc.Input()
	dynMsg := dynamicpb.NewMessage(inputType)

	nameField := inputType.Fields().ByName("name")
	dynMsg.Set(nameField, protoreflect.ValueOfString("hello_world"))

	dm := &DynamicMessage{
		Message:    dynMsg,
		Descriptor: inputType,
	}

	// Test existing field
	val, ok := dm.GetFieldString("name")
	if !ok {
		t.Error("GetFieldString(name) returned false")
	}
	if val != "hello_world" {
		t.Errorf("GetFieldString(name) = %q, want %q", val, "hello_world")
	}

	// Test non-existing field
	_, ok = dm.GetFieldString("nonexistent")
	if ok {
		t.Error("GetFieldString(nonexistent) should return false")
	}
}

func TestDynamicMessage_GetFieldInt(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	inputType := methodDesc.Input()
	dynMsg := dynamicpb.NewMessage(inputType)

	idField := inputType.Fields().ByName("id")
	dynMsg.Set(idField, protoreflect.ValueOfInt64(12345))

	dm := &DynamicMessage{
		Message:    dynMsg,
		Descriptor: inputType,
	}

	// Test existing field
	val, ok := dm.GetFieldInt("id")
	if !ok {
		t.Error("GetFieldInt(id) returned false")
	}
	if val != 12345 {
		t.Errorf("GetFieldInt(id) = %d, want %d", val, 12345)
	}

	// Test non-existing field
	_, ok = dm.GetFieldInt("nonexistent")
	if ok {
		t.Error("GetFieldInt(nonexistent) should return false")
	}
}

func TestDynamicMessage_SetField(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	inputType := methodDesc.Input()
	dynMsg := dynamicpb.NewMessage(inputType)

	dm := &DynamicMessage{
		Message:    dynMsg,
		Descriptor: inputType,
	}

	// Set existing field
	ok := dm.SetField("name", protoreflect.ValueOfString("new_value"))
	if !ok {
		t.Error("SetField(name) returned false")
	}

	// Verify
	val, _ := dm.GetFieldString("name")
	if val != "new_value" {
		t.Errorf("After SetField, name = %q, want %q", val, "new_value")
	}

	// Set non-existing field
	ok = dm.SetField("nonexistent", protoreflect.ValueOfString("value"))
	if ok {
		t.Error("SetField(nonexistent) should return false")
	}
}

func TestDynamicMessage_ToBytes(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	inputType := methodDesc.Input()
	dynMsg := dynamicpb.NewMessage(inputType)

	nameField := inputType.Fields().ByName("name")
	dynMsg.Set(nameField, protoreflect.ValueOfString("serialize_test"))

	dm := &DynamicMessage{
		Message:    dynMsg,
		Descriptor: inputType,
	}

	data, err := dm.ToBytes()
	if err != nil {
		t.Fatalf("ToBytes() error = %v", err)
	}

	// Verify by unmarshaling
	newMsg := dynamicpb.NewMessage(inputType)
	if err := proto.Unmarshal(data, newMsg); err != nil {
		t.Fatalf("Unmarshal verification error = %v", err)
	}

	if got := newMsg.Get(nameField).String(); got != "serialize_test" {
		t.Errorf("After ToBytes, name = %q, want %q", got, "serialize_test")
	}
}

func TestDynamicMessage_ToJSON(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	inputType := methodDesc.Input()
	dynMsg := dynamicpb.NewMessage(inputType)

	nameField := inputType.Fields().ByName("name")
	dynMsg.Set(nameField, protoreflect.ValueOfString("json_test"))

	dm := &DynamicMessage{
		Message:    dynMsg,
		Descriptor: inputType,
	}

	data, err := dm.ToJSON()
	if err != nil {
		t.Fatalf("ToJSON() error = %v", err)
	}

	// Verify it returns valid JSON with the expected field
	if len(data) == 0 {
		t.Error("ToJSON() returned empty data")
	}

	// Should contain the field name in JSON format
	jsonStr := string(data)
	if !strings.Contains(jsonStr, "name") || !strings.Contains(jsonStr, "json_test") {
		t.Errorf("ToJSON() = %s, expected to contain 'name' and 'json_test'", jsonStr)
	}
}

func TestResponseCodec_Name(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	codec := NewResponseCodec(methodDesc)

	if got := codec.Name(); got != "dynamic_proto" {
		t.Errorf("Name() = %q, want %q", got, "dynamic_proto")
	}
}

func TestResponseCodec_MarshalAndUnmarshal(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	codec := NewResponseCodec(methodDesc)

	// Create output message
	outputType := methodDesc.Output()
	dynMsg := dynamicpb.NewMessage(outputType)
	msgField := outputType.Fields().ByName("message")
	dynMsg.Set(msgField, protoreflect.ValueOfString("response_value"))

	dm := &DynamicMessage{
		Message:    dynMsg,
		Descriptor: outputType,
	}

	// Marshal
	data, err := codec.Marshal(dm)
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}

	// Unmarshal
	var result *DynamicMessage
	err = codec.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal error = %v", err)
	}

	// Verify
	val, ok := result.GetFieldString("message")
	if !ok {
		t.Error("GetFieldString(message) returned false")
	}
	if val != "response_value" {
		t.Errorf("message = %q, want %q", val, "response_value")
	}
}

func TestDecodeRequest(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	inputType := methodDesc.Input()

	// Create test message
	srcMsg := dynamicpb.NewMessage(inputType)
	nameField := inputType.Fields().ByName("name")
	srcMsg.Set(nameField, protoreflect.ValueOfString("decode_test"))

	data, err := proto.Marshal(srcMsg)
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}

	// Decode
	result, err := DecodeRequest(methodDesc, data)
	if err != nil {
		t.Fatalf("DecodeRequest error = %v", err)
	}

	// Verify
	val, ok := result.GetFieldString("name")
	if !ok {
		t.Error("GetFieldString(name) returned false")
	}
	if val != "decode_test" {
		t.Errorf("name = %q, want %q", val, "decode_test")
	}
}

func TestDecodeRequest_InvalidData(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)

	// Invalid protobuf data
	invalidData := []byte{0xff, 0xff, 0xff}
	_, err := DecodeRequest(methodDesc, invalidData)
	if err == nil {
		t.Error("DecodeRequest with invalid data should return error")
	}
}

func TestDecodeResponse(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	outputType := methodDesc.Output()

	// Create test message
	srcMsg := dynamicpb.NewMessage(outputType)
	msgField := outputType.Fields().ByName("message")
	srcMsg.Set(msgField, protoreflect.ValueOfString("response_decode_test"))

	data, err := proto.Marshal(srcMsg)
	if err != nil {
		t.Fatalf("Marshal error = %v", err)
	}

	// Decode
	result, err := DecodeResponse(methodDesc, data)
	if err != nil {
		t.Fatalf("DecodeResponse error = %v", err)
	}

	// Verify
	val, ok := result.GetFieldString("message")
	if !ok {
		t.Error("GetFieldString(message) returned false")
	}
	if val != "response_decode_test" {
		t.Errorf("message = %q, want %q", val, "response_decode_test")
	}
}

func TestEncodeMessage(t *testing.T) {
	methodDesc := buildTestMethodDescriptor(t)
	inputType := methodDesc.Input()

	dynMsg := dynamicpb.NewMessage(inputType)
	nameField := inputType.Fields().ByName("name")
	dynMsg.Set(nameField, protoreflect.ValueOfString("encode_test"))

	dm := &DynamicMessage{
		Message:    dynMsg,
		Descriptor: inputType,
	}

	data, err := EncodeMessage(dm)
	if err != nil {
		t.Fatalf("EncodeMessage error = %v", err)
	}

	// Verify by decoding
	result, err := DecodeRequest(methodDesc, data)
	if err != nil {
		t.Fatalf("DecodeRequest verification error = %v", err)
	}

	val, _ := result.GetFieldString("name")
	if val != "encode_test" {
		t.Errorf("After EncodeMessage, name = %q, want %q", val, "encode_test")
	}
}
