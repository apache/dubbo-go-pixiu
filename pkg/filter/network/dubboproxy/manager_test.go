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

package dubboproxy

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"

	"github.com/dubbogo/grpc-go/metadata"

	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestDubboProxyConnectionManager_OnData_InvalidType(t *testing.T) {
	// Create a minimal manager with required dependencies
	dcm := &DubboProxyConnectionManager{
		config:            &model.DubboProxyConnectionManagerConfig{},
		routerCoordinator: nil, // not needed for this test
	}

	// Test with nil data
	result, err := dcm.OnData(nil)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid invocation type")
	assert.Contains(t, err.Error(), "expected *invocation.RPCInvocation")

	// Test with wrong type - string
	result, err = dcm.OnData("invalid")
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid invocation type")

	// Test with wrong type - int
	result, err = dcm.OnData(123)
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid invocation type")

	// Test with wrong pointer type
	result, err = dcm.OnData(&struct{}{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invalid invocation type")
}

func TestDubboProxyConnectionManager_OnData_ValidTypePassesTypeCheck(t *testing.T) {
	// This test verifies that OnData accepts valid RPCInvocation type
	// without panicking on type assertion (the fix in this PR)
	// Note: Full route functionality requires mocking routerCoordinator which is complex

	// Create a manager with proper initialization
	cfg := &model.DubboProxyConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{},
	}
	dcm := CreateDubboProxyConnectionManager(cfg)

	// Create valid RPCInvocation
	invoc := invocation.NewRPCInvocationWithOptions(
		invocation.WithMethodName("testMethod"),
		invocation.WithArguments([]interface{}{}),
	)

	// OnData should not panic on type assertion (the key fix in this PR)
	// It will return an error about route not found, which is acceptable
	result, err := dcm.OnData(invoc)

	// The key assertion: we get a proper error, not a panic
	// Type assertion passed, so we should not see "invalid invocation type" error
	if err != nil {
		assert.NotContains(t, err.Error(), "invalid invocation type",
			"Type assertion should have passed for valid RPCInvocation")
	}
	// Result may be nil due to route error, which is expected
	_ = result // We don't care about result for this test
}

func TestDubboProxyConnectionManager_OnTripleData_EmptyMetadataValue(t *testing.T) {
	dcm := &DubboProxyConnectionManager{
		config:            &model.DubboProxyConnectionManagerConfig{},
		routerCoordinator: nil,
	}

	// Create context with metadata where the values slice is empty
	// This tests the "if len(values) == 0" check in manager.go:114
	md := metadata.MD{
		"test-key": []string{}, // Empty slice triggers the error
	}
	ctx := metadata.NewIncomingContext(context.Background(), md)

	result, err := dcm.OnTripleData(ctx, "testMethod", []interface{}{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "empty metadata value for key: test-key")
}

func TestDubboProxyConnectionManager_OnTripleData_MissingInterfaceKey(t *testing.T) {
	dcm := &DubboProxyConnectionManager{
		config:            &model.DubboProxyConnectionManagerConfig{},
		routerCoordinator: nil,
	}

	// Create context with metadata but no interface key
	md := metadata.Pairs("some-key", "some-value")
	ctx := metadata.NewIncomingContext(context.Background(), md)

	result, err := dcm.OnTripleData(ctx, "testMethod", []interface{}{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "missing or invalid interface key")
}

func TestDubboProxyConnectionManager_OnTripleData_InvalidInterfaceKeyType(t *testing.T) {
	dcm := &DubboProxyConnectionManager{
		config:            &model.DubboProxyConnectionManagerConfig{},
		routerCoordinator: nil,
	}

	// Create context with metadata where interface key exists but is not a string
	// Note: metadata.Pairs always stores string values, so we need to manually construct
	md := metadata.MD{
		"interface": []string{}, // Empty slice - this will trigger empty metadata error
	}
	ctx := metadata.NewIncomingContext(context.Background(), md)

	result, err := dcm.OnTripleData(ctx, "testMethod", []interface{}{})
	assert.Error(t, err)
	assert.Nil(t, result)
	// Should get empty metadata value error
	assert.Contains(t, err.Error(), "empty metadata value")
}

func TestDubboProxyConnectionManager_OnTripleData_NoMetadata(t *testing.T) {
	dcm := &DubboProxyConnectionManager{
		config:            &model.DubboProxyConnectionManagerConfig{},
		routerCoordinator: nil,
	}

	// Create context without metadata
	ctx := context.Background()

	result, err := dcm.OnTripleData(ctx, "testMethod", []interface{}{})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "missing or invalid interface key")
}

func TestCreateDubboProxyConnectionManager(t *testing.T) {
	cfg := &model.DubboProxyConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{},
		TimeoutStr:  "5s",
	}

	dcm := CreateDubboProxyConnectionManager(cfg)
	assert.NotNil(t, dcm)
	assert.NotNil(t, dcm.config)
	assert.NotNil(t, dcm.routerCoordinator)
	assert.NotNil(t, dcm.codec)
	assert.NotNil(t, dcm.filterManager)
}

func TestDubboProxyConnectionManager_OnTripleData_ValidMetadataTypeCheckPasses(t *testing.T) {
	// This test verifies that OnTripleData accepts valid metadata without type assertion panic
	cfg := &model.DubboProxyConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{},
	}
	dcm := CreateDubboProxyConnectionManager(cfg)

	// Create context with valid metadata including interface key
	md := metadata.Pairs(
		"interface", "com.example.TestService",
		"version", "1.0.0",
		"group", "testGroup",
	)
	ctx := metadata.NewIncomingContext(context.Background(), md)

	// OnTripleData should not panic on type assertions (the fix in this PR)
	// It will return an error about route not found, which is acceptable
	result, err := dcm.OnTripleData(ctx, "testMethod", []interface{}{})

	// The key assertion: we get a proper error, not a panic
	// Type assertions passed, so we should not see "missing or invalid interface key" error
	if err != nil {
		assert.NotContains(t, err.Error(), "missing or invalid interface key",
			"Type assertion for interface key should have passed")
		assert.NotContains(t, err.Error(), "empty metadata value",
			"Metadata value check should have passed")
	}
	// Result may be nil due to route error, which is expected
	_ = result
}

func TestDubboProxyConnectionManager_OnEncode(t *testing.T) {
	// Test OnEncode with invalid type
	cfg := &model.DubboProxyConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{},
	}
	dcm := CreateDubboProxyConnectionManager(cfg)

	// Test with invalid type
	_, err := dcm.OnEncode("invalid")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid rpc response")
}
