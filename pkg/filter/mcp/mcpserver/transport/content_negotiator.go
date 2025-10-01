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
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/mark3labs/mcp-go/mcp"
)

// ResponseFormat represents the response format type
type ResponseFormat int

const (
	ResponseFormatJSON     ResponseFormat = iota // Traditional JSON response
	ResponseFormatSSE                            // Server-Sent Events response
	ResponseFormatAccepted                       // 202 Accepted (no immediate response)
)

// MCP methods that prefer streaming responses
var streamingMethods = []string{
	string(mcp.MethodResourcesRead),
	"resources/subscribe",
	"tools/stream",
	"notifications/subscribe",
}

// ContentNegotiator handles HTTP content negotiation for MCP responses
type ContentNegotiator struct{}

// NewContentNegotiator creates a new content negotiator
func NewContentNegotiator() *ContentNegotiator {
	return &ContentNegotiator{}
}

// NegotiateResponse determines the appropriate response format based on request and Accept header
func (cn *ContentNegotiator) NegotiateResponse(acceptHeader string, jsonrpcReq mcp.JSONRPCRequest, hasSession bool) ResponseFormat {
	supportsJSON, supportsSSE := cn.parseAcceptHeader(acceptHeader)

	// If no Accept header specified, default to JSON for backward compatibility
	if acceptHeader == "" || (!supportsJSON && !supportsSSE) {
		return ResponseFormatJSON
	}

	// Note: For true JSON-RPC notifications, we would check if ID is nil,
	// but mcp.RequestId type doesn't support direct nil comparison.
	// In MCP context, method type usually determines response handling.

	switch {
	case supportsJSON && !supportsSSE:
		// Only JSON is supported
		return ResponseFormatJSON
	case !supportsJSON && supportsSSE:
		// Only SSE is supported
		if hasSession {
			return ResponseFormatSSE
		}
		// Fall back to JSON if no session available
		return ResponseFormatJSON
	case supportsJSON && supportsSSE:
		// Both are supported, make intelligent choice based on request type
		if cn.shouldPreferSSE(jsonrpcReq, hasSession) {
			return ResponseFormatSSE
		}
		return ResponseFormatJSON
	default:
		// Should never reach here due to early return at line 57
		return ResponseFormatJSON
	}
}

// SupportsSSE checks if the Accept header includes text/event-stream
func (cn *ContentNegotiator) SupportsSSE(acceptHeader string) bool {
	_, supportsSSE := cn.parseAcceptHeader(acceptHeader)
	return supportsSSE
}

// SupportsJSON checks if the Accept header includes application/json
func (cn *ContentNegotiator) SupportsJSON(acceptHeader string) bool {
	supportsJSON, _ := cn.parseAcceptHeader(acceptHeader)
	return supportsJSON
}

// parseAcceptHeader parses the Accept header and returns support for JSON and SSE
func (cn *ContentNegotiator) parseAcceptHeader(acceptHeader string) (supportsJSON, supportsSSE bool) {
	if acceptHeader == "" {
		// Default to JSON support for backward compatibility
		return true, false
	}

	acceptHeader = strings.ToLower(acceptHeader)

	// Check for JSON support
	if strings.Contains(acceptHeader, constant.HeaderValueApplicationJson) ||
		strings.Contains(acceptHeader, constant.MediaTypeApplicationWild) ||
		strings.Contains(acceptHeader, constant.MediaTypeWildcard) {
		supportsJSON = true
	}

	// Check for SSE support
	if strings.Contains(acceptHeader, constant.HeaderValueTextEventStream) ||
		strings.Contains(acceptHeader, constant.MediaTypeTextWild) ||
		strings.Contains(acceptHeader, constant.MediaTypeWildcard) {
		supportsSSE = true
	}

	return supportsJSON, supportsSSE
}

// shouldPreferSSE determines if SSE should be preferred over JSON for a given request
func (cn *ContentNegotiator) shouldPreferSSE(jsonrpcReq mcp.JSONRPCRequest, hasSession bool) bool {
	// Only prefer SSE if we have an active session
	if !hasSession {
		return false
	}

	// Prefer SSE for tool calls as they may involve long-running operations
	if jsonrpcReq.Method == string(mcp.MethodToolsCall) {
		return true
	}

	// Prefer SSE for methods that might generate server-to-client notifications
	for _, method := range streamingMethods {
		if jsonrpcReq.Method == method {
			return true
		}
	}

	// Default to JSON for other methods
	return false
}

// GetPreferredContentType returns the Content-Type header value for the given format
func (cn *ContentNegotiator) GetPreferredContentType(format ResponseFormat) string {
	switch format {
	case ResponseFormatSSE:
		return constant.HeaderValueTextEventStream
	case ResponseFormatJSON, ResponseFormatAccepted:
		return constant.HeaderValueApplicationJson
	default:
		return constant.HeaderValueApplicationJson
	}
}

// ValidateAcceptHeaderForMethod validates that the Accept header is appropriate for the HTTP method
func (cn *ContentNegotiator) ValidateAcceptHeaderForMethod(method, acceptHeader string) error {
	switch strings.ToUpper(method) {
	case constant.Get:
		// GET requests for SSE must accept text/event-stream
		if !cn.SupportsSSE(acceptHeader) {
			return &AcceptHeaderError{
				Method:       method,
				AcceptHeader: acceptHeader,
				Required:     constant.HeaderValueTextEventStream,
			}
		}
	case constant.Post:
		// POST requests should accept both application/json and text/event-stream
		supportsJSON, supportsSSE := cn.parseAcceptHeader(acceptHeader)
		if !supportsJSON && !supportsSSE {
			return &AcceptHeaderError{
				Method:       method,
				AcceptHeader: acceptHeader,
				Required:     constant.HeaderValueApplicationJson + " or " + constant.HeaderValueTextEventStream,
			}
		}
	}
	return nil
}

// AcceptHeaderError represents an Accept header validation error
type AcceptHeaderError struct {
	Method       string
	AcceptHeader string
	Required     string
}

func (e *AcceptHeaderError) Error() string {
	if e.AcceptHeader == "" {
		return "missing Accept header for " + e.Method + " request, required: " + e.Required
	}
	return "invalid Accept header '" + e.AcceptHeader + "' for " + e.Method + " request, required: " + e.Required
}
