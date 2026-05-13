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

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type normalizeAddressCase struct {
	name        string
	address     string
	port        string
	wantAddress string
	wantErr     bool
}

func normalizeOK(name, address, port, wantAddress string) normalizeAddressCase {
	return normalizeAddressCase{
		name:        name,
		address:     address,
		port:        port,
		wantAddress: wantAddress,
	}
}

func normalizeErr(name, address, port string) normalizeAddressCase {
	return normalizeAddressCase{
		name:    name,
		address: address,
		port:    port,
		wantErr: true,
	}
}

func TestEndpointCheckerHealthHandlersEmitEventsWithoutMutatingEndpoint(t *testing.T) {
	endpoint := &model.Endpoint{
		ID:        "ep-1",
		UnHealthy: true,
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18080,
		},
	}

	events := make([]EndpointHealthEvent, 0, 2)
	hc := &HealthChecker{
		onEndpointHealth: func(event EndpointHealthEvent) {
			events = append(events, event)
		},
	}
	checker := newChecker(endpoint, hc)

	checker.handleHealth()
	checker.handleUnHealth()

	if !endpoint.UnHealthy {
		t.Fatalf("endpoint.UnHealthy was mutated by health handlers")
	}
	if len(events) != 2 {
		t.Fatalf("events length = %d, want 2", len(events))
	}
	if !events[0].Healthy {
		t.Fatalf("first event healthy = false, want true")
	}
	if events[0].EndpointID != endpoint.ID {
		t.Fatalf("first event endpoint ID = %q, want %q", events[0].EndpointID, endpoint.ID)
	}
	if events[0].EndpointAddress != endpoint.Address.GetAddress() {
		t.Fatalf("first event endpoint address = %q, want %q", events[0].EndpointAddress, endpoint.Address.GetAddress())
	}
	if events[1].Healthy {
		t.Fatalf("second event healthy = true, want false")
	}
}

func TestEndpointCheckerHealthHandlersMutateEndpointWithoutListener(t *testing.T) {
	endpoint := &model.Endpoint{
		ID: "ep-1",
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18081,
		},
	}
	checker := newChecker(endpoint, &HealthChecker{})

	checker.handleUnHealth()
	if !endpoint.UnHealthy {
		t.Fatalf("endpoint.UnHealthy = false after unhealth event, want true")
	}

	checker.handleHealth()
	if endpoint.UnHealthy {
		t.Fatalf("endpoint.UnHealthy = true after health event, want false")
	}
}

func TestEndpointCheckerHealthHandlersMutateSameAddressEndpointsWithoutListener(t *testing.T) {
	first := &model.Endpoint{
		ID: "ep-1",
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18082,
		},
	}
	second := &model.Endpoint{
		ID: "ep-2",
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18082,
		},
	}
	otherAddress := &model.Endpoint{
		ID: "ep-3",
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18083,
		},
	}
	checker := newChecker(first, &HealthChecker{
		cluster: &model.ClusterConfig{
			Endpoints: []*model.Endpoint{first, second, otherAddress},
		},
	})

	checker.handleUnHealth()

	if !first.UnHealthy {
		t.Fatalf("first endpoint UnHealthy = false after shared-address unhealth event, want true")
	}
	if !second.UnHealthy {
		t.Fatalf("second endpoint UnHealthy = false after shared-address unhealth event, want true")
	}
	if otherAddress.UnHealthy {
		t.Fatalf("other-address endpoint UnHealthy = true after shared-address unhealth event, want false")
	}

	checker.handleHealth()

	if first.UnHealthy {
		t.Fatalf("first endpoint UnHealthy = true after shared-address health event, want false")
	}
	if second.UnHealthy {
		t.Fatalf("second endpoint UnHealthy = true after shared-address health event, want false")
	}
}

func TestHealthCheckerStopOneKeepsSharedAddressChecker(t *testing.T) {
	first := &model.Endpoint{
		ID: "ep-1",
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18084,
		},
	}
	second := &model.Endpoint{
		ID: "ep-2",
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18084,
		},
	}
	addr := first.Address.GetAddress()
	checker := &EndpointChecker{
		stop: make(chan struct{}),
	}
	hc := &HealthChecker{
		cluster: &model.ClusterConfig{
			Endpoints: []*model.Endpoint{first, second},
		},
		checkers: map[string]*EndpointChecker{
			addr: checker,
		},
	}

	hc.stopCheck(first)
	if _, ok := hc.checkers[addr]; !ok {
		t.Fatalf("shared-address checker was stopped while another endpoint still used the address")
	}

	hc.cluster.Endpoints = []*model.Endpoint{first}
	hc.stopCheck(first)
	if _, ok := hc.checkers[addr]; ok {
		t.Fatalf("checker was kept after the last endpoint using the address was removed")
	}
	select {
	case <-checker.stop:
	default:
		t.Fatalf("checker stop channel was not closed")
	}
}

func TestNormalizeAddress(t *testing.T) {
	tests := []normalizeAddressCase{
		normalizeOK("port is empty, address has port", "localhost:8080", "", "localhost:8080"),
		normalizeErr("port is empty, address has no port", "localhost", ""),
		normalizeOK("port is not empty, address has no port", "localhost", "80", "localhost:80"),
		normalizeOK("port is not empty, address has same port", "localhost:80", "80", "localhost:80"),
		normalizeOK("port is not empty, address has different port", "localhost:8080", "80", "localhost:80"),
		normalizeErr("invalid IPv6 address format for empty port", "[::1]", ""),
		normalizeOK("valid IPv6 address with port", "[::1]:8080", "", "[::1]:8080"),
		normalizeOK("port is not empty, valid IPv6 address without port", "[::1]", "80", "[::1]:80"),
		normalizeOK("port is not empty, valid IPv6 address with different port", "[::1]:8080", "80", "[::1]:80"),
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

func acceptOneTCPConnection(listener net.Listener) {
	go func() {
		conn, _ := listener.Accept()
		if conn != nil {
			conn.Close()
		}
	}()
}

func assertCheckTcpConn(t *testing.T, name, host, port string, want bool) {
	t.Helper()
	t.Run(name, func(t *testing.T) {
		success := CheckTcpConn(host, port, 100*time.Millisecond)
		if success != want {
			t.Errorf("CheckTcpConn(%q, %q, ...) = %t, want %t", host, port, success, want)
		}
	})
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

	acceptOneTCPConnection(listener)
	assertCheckTcpConn(t, "successful connection", host, portStr, true)

	// Failure case 1: Invalid address format
	assertCheckTcpConn(t, "failed connection due to invalid address format", "127.0.0.1:80:90", "80", false)

	// Failure case 2: Connection timeout
	assertCheckTcpConn(t, "failed connection due to timeout", "127.0.0.1", "80", false)

	// Test with empty port (should fail due to normalizeAddress)
	assertCheckTcpConn(t, "failed with empty port and no port in address", "localhost", "", false)

	// Test with empty port and address has port (should succeed)
	acceptOneTCPConnection(listener)
	assertCheckTcpConn(t, "successful with empty port and port in address", addr, "", true)
}
