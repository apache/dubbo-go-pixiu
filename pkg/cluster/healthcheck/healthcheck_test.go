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

func TestCreateHealthCheckParsesInitialDelaySeconds(t *testing.T) {
	hc := CreateHealthCheck(&model.ClusterConfig{}, model.HealthCheckConfig{
		TimeoutConfig:       "1s",
		IntervalConfig:      "30s",
		InitialDelaySeconds: "10",
		HealthyThreshold:    1,
		UnhealthyThreshold:  1,
	})

	if hc.initialDelay != 10*time.Second {
		t.Fatalf("initialDelay = %s, want 10s", hc.initialDelay)
	}
}

func TestCreateHealthCheckDefaultsInitialDelayWhenUnset(t *testing.T) {
	hc := CreateHealthCheck(&model.ClusterConfig{}, model.HealthCheckConfig{
		TimeoutConfig:      "1s",
		IntervalConfig:     "30s",
		HealthyThreshold:   1,
		UnhealthyThreshold: 1,
	})

	if hc.initialDelay != DefaultFirstInterval {
		t.Fatalf("initialDelay = %s, want %s", hc.initialDelay, DefaultFirstInterval)
	}
}

// newTestChecker builds an EndpointChecker driving a stub HealthChecker
// with the supplied thresholds and an event-capture callback. No
// goroutines are started — the test drives HandleSuccess / HandleFailure
// directly, so behavior is deterministic and isolated from the timer
// loop.
func newTestChecker(healthyThreshold, unhealthyThreshold uint32) (*EndpointChecker, *[]EndpointHealthEvent) {
	events := make([]EndpointHealthEvent, 0)
	eventsPtr := &events
	hc := &HealthChecker{
		healthyThreshold:   healthyThreshold,
		unhealthyThreshold: unhealthyThreshold,
		onEndpointHealth: func(ev EndpointHealthEvent) {
			*eventsPtr = append(*eventsPtr, ev)
		},
	}
	c := &EndpointChecker{
		HealthChecker: hc,
		endpointID:    "test-ep",
		endpointAddr:  "127.0.0.1:18080",
	}
	return c, eventsPtr
}

// TestHandleSuccessHonorsHealthyThreshold proves the configured
// healthyThreshold actually takes effect. Prior to v1.2 the threshold
// was compared against an uninitialized uint32 field that was always
// zero, so the first success flipped state regardless of config.
func TestHandleSuccessHonorsHealthyThreshold(t *testing.T) {
	c, events := newTestChecker(3, 5)

	// First two successes must NOT emit a healthy event.
	c.HandleSuccess()
	c.HandleSuccess()
	if len(*events) != 0 {
		t.Fatalf("expected no health events before threshold; got %d", len(*events))
	}

	// Third success crosses the threshold (>= 3).
	c.HandleSuccess()
	if len(*events) != 1 {
		t.Fatalf("expected exactly one healthy event at threshold; got %d", len(*events))
	}
	if !(*events)[0].Healthy {
		t.Fatalf("expected first event to be healthy=true; got %+v", (*events)[0])
	}
}

// TestHandleFailureHonorsUnhealthyThresholdForBothFailureModes proves
// the configured unhealthyThreshold takes effect AND that both
// HandleFailure(false) and HandleFailure(true) feed the same counter.
// Prior to v1.2 the false branch flipped state immediately (no counter)
// and the true branch compared against an uninitialized threshold.
func TestHandleFailureHonorsUnhealthyThresholdForBothFailureModes(t *testing.T) {
	c, events := newTestChecker(3, 4)

	// Mix two non-timeout failures and one timeout — still below threshold.
	c.HandleFailure(false)
	c.HandleFailure(false)
	c.HandleFailure(true)
	if len(*events) != 0 {
		t.Fatalf("expected no unhealthy event below threshold; got %d", len(*events))
	}

	// Fourth failure (any mode) crosses the threshold (>= 4).
	c.HandleFailure(false)
	if len(*events) != 1 {
		t.Fatalf("expected exactly one unhealthy event at threshold; got %d", len(*events))
	}
	if (*events)[0].Healthy {
		t.Fatalf("expected first event to be healthy=false; got %+v", (*events)[0])
	}
}

// TestHandleSuccessAndFailureResetEachOtherCounters keeps the existing
// reset behavior locked: a success run clears unHealthCount, a failure
// run clears healthCount, so transitions require fresh consecutive
// signals after a flap.
func TestHandleSuccessAndFailureResetEachOtherCounters(t *testing.T) {
	c, events := newTestChecker(2, 2)

	c.HandleSuccess()      // healthy=1
	c.HandleFailure(false) // unhealthy=1, healthy reset to 0
	c.HandleSuccess()      // healthy=1 (not 2, because failure reset it)
	if len(*events) != 0 {
		t.Fatalf("expected no events yet (each counter at 1, threshold 2); got %d", len(*events))
	}

	c.HandleSuccess() // healthy=2 -> flip
	if len(*events) != 1 || !(*events)[0].Healthy {
		t.Fatalf("expected one healthy event after two consecutive successes; got %+v", *events)
	}
}

// TestHandleTimeoutRoutesThroughHandleFailure ensures the legacy
// HandleTimeout entry point still works for external Checker
// implementations and feeds the unified failure counter.
func TestHandleTimeoutRoutesThroughHandleFailure(t *testing.T) {
	c, events := newTestChecker(3, 2)

	c.HandleTimeout()
	if len(*events) != 0 {
		t.Fatalf("expected no event below threshold; got %d", len(*events))
	}
	c.HandleTimeout()
	if len(*events) != 1 || (*events)[0].Healthy {
		t.Fatalf("expected unhealthy event after two timeouts; got %+v", *events)
	}
}
