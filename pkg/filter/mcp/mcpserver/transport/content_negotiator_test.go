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

package transport

import (
	"testing"
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewContentNegotiator(t *testing.T) {
	cn := NewContentNegotiator()
	if cn == nil {
		t.Fatal("NewContentNegotiator returned nil")
	}
}

func TestParseAcceptHeader(t *testing.T) {
	cn := NewContentNegotiator()

	tests := []struct {
		name         string
		acceptHeader string
		wantJSON     bool
		wantSSE      bool
	}{
		{
			name:         "empty header defaults to JSON",
			acceptHeader: "",
			wantJSON:     true,
			wantSSE:      false,
		},
		{
			name:         "application/json only",
			acceptHeader: "application/json",
			wantJSON:     true,
			wantSSE:      false,
		},
		{
			name:         "text/event-stream only",
			acceptHeader: "text/event-stream",
			wantJSON:     false,
			wantSSE:      true,
		},
		{
			name:         "both formats",
			acceptHeader: "application/json, text/event-stream",
			wantJSON:     true,
			wantSSE:      true,
		},
		{
			name:         "wildcard accepts all",
			acceptHeader: "*/*",
			wantJSON:     true,
			wantSSE:      true,
		},
		{
			name:         "application wildcard",
			acceptHeader: "application/*",
			wantJSON:     true,
			wantSSE:      false,
		},
		{
			name:         "text wildcard",
			acceptHeader: "text/*",
			wantJSON:     false,
			wantSSE:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJSON, gotSSE := cn.parseAcceptHeader(tt.acceptHeader)
			if gotJSON != tt.wantJSON {
				t.Errorf("parseAcceptHeader() JSON = %v, want %v", gotJSON, tt.wantJSON)
			}
			if gotSSE != tt.wantSSE {
				t.Errorf("parseAcceptHeader() SSE = %v, want %v", gotSSE, tt.wantSSE)
			}
		})
	}
}

func TestSupportsSSE(t *testing.T) {
	cn := NewContentNegotiator()

	tests := []struct {
		name         string
		acceptHeader string
		want         bool
	}{
		{"text/event-stream", "text/event-stream", true},
		{"with charset", "text/event-stream;charset=utf-8", true},
		{"wildcard", "*/*", true},
		{"text wildcard", "text/*", true},
		{"JSON only", "application/json", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cn.SupportsSSE(tt.acceptHeader); got != tt.want {
				t.Errorf("SupportsSSE() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSupportsJSON(t *testing.T) {
	cn := NewContentNegotiator()

	tests := []struct {
		name         string
		acceptHeader string
		want         bool
	}{
		{"application/json", "application/json", true},
		{"with charset", "application/json;charset=utf-8", true},
		{"wildcard", "*/*", true},
		{"application wildcard", "application/*", true},
		{"SSE only", "text/event-stream", false},
		{"empty defaults to JSON", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cn.SupportsJSON(tt.acceptHeader); got != tt.want {
				t.Errorf("SupportsJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNegotiateResponse(t *testing.T) {
	cn := NewContentNegotiator()

	tests := []struct {
		name         string
		acceptHeader string
		method       string
		hasSession   bool
		want         ResponseFormat
	}{
		{
			name:         "no accept header defaults to JSON",
			acceptHeader: "",
			method:       "initialize",
			hasSession:   false,
			want:         ResponseFormatJSON,
		},
		{
			name:         "JSON only returns JSON",
			acceptHeader: "application/json",
			method:       "initialize",
			hasSession:   true,
			want:         ResponseFormatJSON,
		},
		{
			name:         "SSE only with session",
			acceptHeader: "text/event-stream",
			method:       "initialize",
			hasSession:   true,
			want:         ResponseFormatSSE, // Returns SSE if session exists and only SSE accepted
		},
		{
			name:         "tool call with both and session prefers SSE",
			acceptHeader: "application/json, text/event-stream",
			method:       string(mcp.MethodToolsCall),
			hasSession:   true,
			want:         ResponseFormatSSE,
		},
		{
			name:         "tool call without session uses JSON",
			acceptHeader: "application/json, text/event-stream",
			method:       string(mcp.MethodToolsCall),
			hasSession:   false,
			want:         ResponseFormatJSON,
		},
		{
			name:         "any method with both and session prefers SSE",
			acceptHeader: "application/json, text/event-stream",
			method:       string(mcp.MethodToolsList),
			hasSession:   true,
			want:         ResponseFormatSSE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := cn.NegotiateResponse(tt.acceptHeader, tt.hasSession)
			if got != tt.want {
				t.Errorf("NegotiateResponse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPreferredContentType(t *testing.T) {
	cn := NewContentNegotiator()

	tests := []struct {
		name   string
		format ResponseFormat
		want   string
	}{
		{
			name:   "JSON format",
			format: ResponseFormatJSON,
			want:   "application/json",
		},
		{
			name:   "SSE format",
			format: ResponseFormatSSE,
			want:   "text/event-stream",
		},
		{
			name:   "Accepted format",
			format: ResponseFormatAccepted,
			want:   "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cn.GetPreferredContentType(tt.format); got != tt.want {
				t.Errorf("GetPreferredContentType() = %v, want %v", got, tt.want)
			}
		})
	}
}
