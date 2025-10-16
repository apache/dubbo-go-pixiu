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

package http

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse represents the unified error response structure
type ErrorResponse struct {
	Status  int    `json:"status"`  // HTTP status code
	Message string `json:"message"` // Standard error message
	Err     error  `json:"-"`       // Underlying error, not marshaled directly
}

// ErrorBuilder builds error responses
type ErrorBuilder struct {
	status  int
	message string
}

// Predefined error response builders
var (
	// 404 - Not Found
	RouteNotFound   = newErrorBuilder(http.StatusNotFound, "Route not found")
	ServiceNotFound = newErrorBuilder(http.StatusNotFound, "Service not found")
	APINotFound     = newErrorBuilder(http.StatusNotFound, "API not found")

	// 400 - Bad Request
	BadRequest = newErrorBuilder(http.StatusBadRequest, "Bad request")

	// 401 - Unauthorized
	Unauthorized = newErrorBuilder(http.StatusUnauthorized, "Unauthorized")

	// 403 - Forbidden
	Forbidden = newErrorBuilder(http.StatusForbidden, "Forbidden")

	// 405 - Method Not Allowed
	MethodNotAllowed = newErrorBuilder(http.StatusMethodNotAllowed, "Method not allowed")

	// 406 - Not Acceptable
	NotAcceptable = newErrorBuilder(http.StatusNotAcceptable, "Not acceptable")

	// 429 - Rate Limited
	RateLimited = newErrorBuilder(http.StatusTooManyRequests, "Rate limited")

	// 500 - Internal Server Error
	InternalError      = newErrorBuilder(http.StatusInternalServerError, "Internal server error")
	ConfigurationError = newErrorBuilder(http.StatusInternalServerError, "Configuration error")

	// 502 - Bad Gateway
	BadGateway = newErrorBuilder(http.StatusBadGateway, "Bad gateway")

	// 503 - Service Unavailable
	ServiceUnavailable = newErrorBuilder(http.StatusServiceUnavailable, "Service unavailable")

	// 504 - Gateway Timeout
	GatewayTimeout = newErrorBuilder(http.StatusGatewayTimeout, "Gateway timeout")
)

func newErrorBuilder(status int, message string) *ErrorBuilder {
	return &ErrorBuilder{
		status:  status,
		message: message,
	}
}

// New creates a standard error response without details
func (eb *ErrorBuilder) New() *ErrorResponse {
	return &ErrorResponse{
		Status:  eb.status,
		Message: eb.message,
	}
}

func (eb *ErrorBuilder) WithError(err error) *ErrorResponse {
	return &ErrorResponse{
		Status:  eb.status,
		Message: eb.message,
		Err:     err,
	}
}

func (eb *ErrorBuilder) GetStatus() int {
	return eb.status
}

func (e *ErrorResponse) ToJSON() []byte {
	type alias struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}
	payload := alias{Status: e.Status, Message: e.Message}
	if e.Err != nil {
		payload.Error = e.Err.Error()
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return []byte(`{"status":500,"message":"Internal server error"}`)
	}
	return data
}

// Error implements the error interface
func (e *ErrorResponse) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %s", e.Status, e.Message, e.Err.Error())
	}
	return fmt.Sprintf("[%d] %s", e.Status, e.Message)
}
