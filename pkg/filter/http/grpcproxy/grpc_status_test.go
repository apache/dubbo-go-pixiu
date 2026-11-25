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

package grpcproxy

import (
	"net/http"
	"testing"
)

import "google.golang.org/grpc/codes"

func TestGRPCStatusToHTTPStatus(t *testing.T) {
	tests := []struct {
		code codes.Code
		want int
	}{
		{code: codes.OK, want: http.StatusOK},
		{code: codes.Canceled, want: 499},
		{code: codes.Unknown, want: http.StatusInternalServerError},
		{code: codes.InvalidArgument, want: http.StatusBadRequest},
		{code: codes.DeadlineExceeded, want: http.StatusGatewayTimeout},
		{code: codes.NotFound, want: http.StatusNotFound},
		{code: codes.AlreadyExists, want: http.StatusConflict},
		{code: codes.PermissionDenied, want: http.StatusForbidden},
		{code: codes.ResourceExhausted, want: http.StatusTooManyRequests},
		{code: codes.FailedPrecondition, want: http.StatusBadRequest},
		{code: codes.Aborted, want: http.StatusConflict},
		{code: codes.OutOfRange, want: http.StatusBadRequest},
		{code: codes.Unimplemented, want: http.StatusNotImplemented},
		{code: codes.Internal, want: http.StatusInternalServerError},
		{code: codes.Unavailable, want: http.StatusServiceUnavailable},
		{code: codes.DataLoss, want: http.StatusInternalServerError},
		{code: codes.Unauthenticated, want: http.StatusUnauthorized},
		{code: codes.Code(99), want: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.code.String(), func(t *testing.T) {
			if got := grpcStatusToHTTPStatus(tt.code); got != tt.want {
				t.Fatalf("grpcStatusToHTTPStatus(%s) = %d, want %d", tt.code, got, tt.want)
			}
		})
	}
}
