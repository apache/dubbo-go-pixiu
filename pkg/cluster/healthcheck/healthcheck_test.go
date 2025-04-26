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

package healthcheck

import (
	"net"
	"testing"
	"time"
)

func TestNormalizeAddress(t *testing.T) {
	tests := []struct {
		name        string
		address     string
		port        string
		wantAddress string
		wantErr     bool
	}{
		{
			name:        "port is empty, address has port",
			address:     "localhost:8080",
			port:        "",
			wantAddress: "localhost:8080",
			wantErr:     false,
		},
		{
			name:        "port is empty, address has no port",
			address:     "localhost",
			port:        "",
			wantAddress: "",
			wantErr:     true,
		},
		{
			name:        "port is not empty, address has no port",
			address:     "localhost",
			port:        "80",
			wantAddress: "localhost:80",
			wantErr:     false,
		},
		{
			name:        "port is not empty, address has same port",
			address:     "localhost:80",
			port:        "80",
			wantAddress: "localhost:80",
			wantErr:     false,
		},
		{
			name:        "port is not empty, address has different port",
			address:     "localhost:8080",
			port:        "80",
			wantAddress: "localhost:80",
			wantErr:     false,
		},
		{
			name:        "invalid address format for empty port",
			address:     "[::1]", // IPv6 without port
			port:        "",
			wantAddress: "",
			wantErr:     true,
		},
		{
			name:        "valid IPv6 address with port",
			address:     "[::1]:8080",
			port:        "",
			wantAddress: "[::1]:8080",
			wantErr:     false,
		},
		{
			name:        "port is not empty, valid IPv6 address without port",
			address:     "[::1]",
			port:        "80",
			wantAddress: "[::1]:80",
			wantErr:     false,
		},
		{
			name:        "port is not empty, valid IPv6 address with different port",
			address:     "[::1]:8080",
			port:        "80",
			wantAddress: "[::1]:80",
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAddress, gotErr := normalizeAddress(tt.address, tt.port)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("normalizeAddress(%q, %q) error = %v, wantErr %v", tt.address, tt.port, gotErr, tt.wantErr)
				return
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("normalizeAddress(%q, %q) gotAddress = %q, want %q", tt.address, tt.port, gotAddress, tt.wantAddress)
			}
		})
	}
}

func TestCheckTcpConn(t *testing.T) {
	// We need a way to simulate a successful and a failed TCP connection.
	// We can achieve this by setting up a temporary listener for the success case
	// and using an invalid address for the failure case.

	// Success case: Set up a temporary listener
	listener, err := net.Listen("tcp", "localhost:0") // Listen on a random available port
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	defer listener.Close()
	addr := listener.Addr().String()
	host, portStr, _ := net.SplitHostPort(addr)

	t.Run("successful connection", func(t *testing.T) {
		go func() {
			conn, _ := listener.Accept() // Accept the incoming connection
			if conn != nil {
				conn.Close()
			}
		}()
		success := CheckTcpConn(host, portStr, 100*time.Millisecond)
		if !success {
			t.Errorf("CheckTcpConn(%q, %q, ...) should return true for a successful connection", host, portStr)
		}
	})

	// Failure case 1: Invalid address format
	t.Run("failed connection due to invalid address format", func(t *testing.T) {
		success := CheckTcpConn("invalid address", "80", 100*time.Millisecond)
		if success {
			t.Errorf("CheckTcpConn(%q, %q, ...) should return false for an invalid address format", "invalid address", "80")
		}
	})

	// Failure case 2: Connection timeout
	t.Run("failed connection due to timeout", func(t *testing.T) {
		// Use a non-routable local address to ensure a timeout
		success := CheckTcpConn("127.0.0.1", "80", 100*time.Millisecond)
		if success {
			t.Errorf("CheckTcpConn(%q, %q, ...) should return false due to timeout", "127.0.0.1", "80")
		}
	})

	// Test with empty port (should fail due to normalizeAddress)
	t.Run("failed with empty port and no port in address", func(t *testing.T) {
		success := CheckTcpConn("localhost", "", 100*time.Millisecond)
		if success {
			t.Errorf("CheckTcpConn(%q, %q, ...) should return false when port is empty and address has no port", "localhost", "")
		}
	})

	// Test with empty port and address has port (should succeed)
	t.Run("successful with empty port and port in address", func(t *testing.T) {
		go func() {
			conn, _ := listener.Accept()
			if conn != nil {
				conn.Close()
			}
		}()
		success := CheckTcpConn(addr, "", 100*time.Millisecond)
		if !success {
			t.Errorf("CheckTcpConn(%q, %q, ...) should return true when port is empty and address has port", addr, "")
		}
	})
}
