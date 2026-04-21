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

package dubbo

import (
	"testing"
	"time"
)

var _ DubboClient = (*Client)(nil)

func TestDubboOutboundRequestFields(t *testing.T) {
	req := &DubboOutboundRequest{
		Service:       "com.example.UserService",
		Method:        "GetUser",
		Group:         "gray",
		Version:       "1.0.0",
		Address:       "127.0.0.1:20880",
		Protocol:      "dubbo",
		Serialization: "hessian2",
		Arguments:     []any{"user-1", 2},
		ParamTypes:    []string{"java.lang.String", "int"},
		Attachments:   map[string]any{"traceparent": "00-test"},
		Timeout:       3 * time.Second,
	}

	if req.Service != "com.example.UserService" {
		t.Fatalf("unexpected service: %q", req.Service)
	}
	if req.Method != "GetUser" {
		t.Fatalf("unexpected method: %q", req.Method)
	}
	if req.Group != "gray" {
		t.Fatalf("unexpected group: %q", req.Group)
	}
	if req.Version != "1.0.0" {
		t.Fatalf("unexpected version: %q", req.Version)
	}
	if req.Address != "127.0.0.1:20880" {
		t.Fatalf("unexpected address: %q", req.Address)
	}
	if req.Protocol != "dubbo" {
		t.Fatalf("unexpected protocol: %q", req.Protocol)
	}
	if req.Serialization != "hessian2" {
		t.Fatalf("unexpected serialization: %q", req.Serialization)
	}
	if len(req.Arguments) != 2 {
		t.Fatalf("unexpected arguments: %#v", req.Arguments)
	}
	if len(req.ParamTypes) != 2 {
		t.Fatalf("unexpected param types: %#v", req.ParamTypes)
	}
	if req.Attachments["traceparent"] != "00-test" {
		t.Fatalf("unexpected attachments: %#v", req.Attachments)
	}
	if req.Timeout != 3*time.Second {
		t.Fatalf("unexpected timeout: %s", req.Timeout)
	}
}

func TestDubboOutboundRequestKeepsExplicitZeroParamTypes(t *testing.T) {
	req := &DubboOutboundRequest{
		ParamTypes: []string{},
	}

	if req.ParamTypes == nil {
		t.Fatal("expected explicit empty param types slice to be preserved")
	}
	if len(req.ParamTypes) != 0 {
		t.Fatalf("expected empty param types slice, got %#v", req.ParamTypes)
	}
}

func TestDubboOutboundRequestKeepsNilParamTypes(t *testing.T) {
	req := &DubboOutboundRequest{
		ParamTypes: nil,
	}

	if req.ParamTypes != nil {
		t.Fatalf("expected nil param types slice, got %#v", req.ParamTypes)
	}
}
