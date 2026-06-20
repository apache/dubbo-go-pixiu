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
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

// ResponseFormat represents the response format type
type ResponseFormat int

const (
	ResponseFormatJSON     ResponseFormat = iota // Traditional JSON response
	ResponseFormatSSE                            // Server-Sent Events response
	ResponseFormatAccepted                       // 202 Accepted (no immediate response)
)

// ContentNegotiator handles HTTP content negotiation for MCP responses
type ContentNegotiator struct{}

// NewContentNegotiator creates a new content negotiator
func NewContentNegotiator() *ContentNegotiator {
	return &ContentNegotiator{}
}

// NegotiateResponse determines the appropriate response format based on Accept header and session state
func (cn *ContentNegotiator) NegotiateResponse(acceptHeader string, hasSession bool) ResponseFormat {
	_, supportsSSE := cn.parseAcceptHeader(acceptHeader)

	// Per MCP spec: when a session exists and client accepts SSE, use SSE for consistent streaming
	// Otherwise, default to JSON (for backward compatibility and when no session available)
	if supportsSSE && hasSession {
		return ResponseFormatSSE
	}

	return ResponseFormatJSON
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
