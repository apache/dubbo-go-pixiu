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

package filter

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestEmptyNetworkFilter_OnDecode(t *testing.T) {
	enf := &EmptyNetworkFilter{}
	data := []byte("test data")

	result, n, err := enf.OnDecode(data)

	if result != nil {
		t.Errorf("OnDecode should return nil, got %v", result)
	}
	if n != 0 {
		t.Errorf("OnDecode should return 0, got %d", n)
	}
	if err == nil {
		t.Error("OnDecode should return an error")
	}
	if !strings.Contains(err.Error(), "OnDecode not implemented") {
		t.Errorf("OnDecode error should contain 'OnDecode not implemented', got %s", err.Error())
	}
}

func TestEmptyNetworkFilter_OnEncode(t *testing.T) {
	enf := &EmptyNetworkFilter{}

	result, err := enf.OnEncode("test")

	if result != nil {
		t.Errorf("OnEncode should return nil, got %v", result)
	}
	if err == nil {
		t.Error("OnEncode should return an error")
	}
	if !strings.Contains(err.Error(), "OnEncode not implemented") {
		t.Errorf("OnEncode error should contain 'OnEncode not implemented', got %s", err.Error())
	}
}

func TestEmptyNetworkFilter_OnData(t *testing.T) {
	enf := &EmptyNetworkFilter{}

	result, err := enf.OnData("test data")

	if result != nil {
		t.Errorf("OnData should return nil, got %v", result)
	}
	if err == nil {
		t.Error("OnData should return an error")
	}
	if !strings.Contains(err.Error(), "OnData not implemented") {
		t.Errorf("OnData error should contain 'OnData not implemented', got %s", err.Error())
	}
}

func TestEmptyNetworkFilter_OnTripleData(t *testing.T) {
	enf := &EmptyNetworkFilter{}
	ctx := context.Background()

	result, err := enf.OnTripleData(ctx, "testMethod", []any{"arg1", "arg2"})

	if result != nil {
		t.Errorf("OnTripleData should return nil, got %v", result)
	}
	if err == nil {
		t.Error("OnTripleData should return an error")
	}
	if !strings.Contains(err.Error(), "OnTripleData not implemented") {
		t.Errorf("OnTripleData error should contain 'OnTripleData not implemented', got %s", err.Error())
	}
}

func TestEmptyNetworkFilter_OnUnaryRPC(t *testing.T) {
	enf := &EmptyNetworkFilter{}
	ctx := context.Background()

	result, err := enf.OnUnaryRPC(ctx, "/test/service/method", "request")

	if result != nil {
		t.Errorf("OnUnaryRPC should return nil, got %v", result)
	}
	if err == nil {
		t.Error("OnUnaryRPC should return an error")
	}
	if !strings.Contains(err.Error(), "OnUnaryRPC not implemented") {
		t.Errorf("OnUnaryRPC error should contain 'OnUnaryRPC not implemented', got %s", err.Error())
	}
}

func TestEmptyNetworkFilter_OnStreamRPC(t *testing.T) {
	enf := &EmptyNetworkFilter{}
	stream := &mockRPCStream{}
	info := &model.RPCStreamInfo{}

	err := enf.OnStreamRPC(stream, info)

	if err == nil {
		t.Error("OnStreamRPC should return an error")
	}
	if !strings.Contains(err.Error(), "OnStreamRPC not implemented") {
		t.Errorf("OnStreamRPC error should contain 'OnStreamRPC not implemented', got %s", err.Error())
	}
}

func TestEmptyNetworkFilter_ServeHTTP(t *testing.T) {
	enf := &EmptyNetworkFilter{}
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	enf.ServeHTTP(w, req)

	if w.Code != http.StatusNotImplemented {
		t.Errorf("ServeHTTP should return status %d, got %d", http.StatusNotImplemented, w.Code)
	}
	body := w.Body.String()
	if !strings.Contains(body, "ServeHTTP not implemented") {
		t.Errorf("ServeHTTP response should contain 'ServeHTTP not implemented', got %s", body)
	}
}

func TestEmptyNetworkFilter_Close(t *testing.T) {
	enf := &EmptyNetworkFilter{}

	err := enf.Close()

	if err != nil {
		t.Errorf("Close should return nil, got %v", err)
	}
}

// mockRPCStream is a mock implementation of model.RPCStream for testing
type mockRPCStream struct{}

func (m *mockRPCStream) Context() context.Context {
	return context.Background()
}

func (m *mockRPCStream) SendMsg(msg interface{}) error {
	return nil
}

func (m *mockRPCStream) RecvMsg(msg interface{}) error {
	return nil
}
