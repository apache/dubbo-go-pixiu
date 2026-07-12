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

package core

import (
	"context"
	"errors"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	envoyServer "github.com/envoyproxy/go-control-plane/pkg/server/v3"

	"github.com/apache/dubbo-go-pixiu/admin/global"
	"github.com/apache/dubbo-go-pixiu/admin/utils"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// errSentinel is a distinct error used to assert exact propagation.
var errSentinel = errors.New("sentinel startup error")

// TestRunCoordinated_FirstErrorWins verifies that the first worker to fail is
// propagated and that its peers observe the cancellation and exit.
func TestRunCoordinated_FirstErrorWins(t *testing.T) {
	peerDone := make(chan struct{})

	workers := []func(context.Context) error{
		// Failing worker: returns the sentinel error immediately.
		func(ctx context.Context) error {
			return errSentinel
		},
		// Blocking worker: should be canceled by the coordinator and exit.
		func(ctx context.Context) error {
			<-ctx.Done()
			close(peerDone)
			return nil
		},
	}

	done := make(chan error, 1)
	go func() { done <- runCoordinated(context.Background(), workers...) }()

	select {
	case err := <-done:
		if !assert.ErrorIs(t, err, errSentinel) {
			t.Errorf("expected sentinel error, got %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("runCoordinated did not return in time")
	}

	select {
	case <-peerDone:
		// peer observed cancellation
	case <-time.After(2 * time.Second):
		t.Fatal("peer worker was not canceled")
	}
}

// TestRunCoordinated_NoErrorReturnsNil verifies the happy path: workers that
// return nil produce a nil result.
func TestRunCoordinated_NoErrorReturnsNil(t *testing.T) {
	workers := []func(context.Context) error{
		func(context.Context) error { return nil },
		func(context.Context) error { return nil },
	}

	done := make(chan error, 1)
	go func() { done <- runCoordinated(context.Background(), workers...) }()

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("runCoordinated did not return in time")
	}
}

// TestRunCoordinated_OnlyFirstErrorPropagated verifies that when both workers
// fail, only the first reported error is returned.
func TestRunCoordinated_OnlyFirstErrorPropagated(t *testing.T) {
	errFirst := errors.New("first")
	errSecond := errors.New("second")
	firstReported := make(chan struct{})

	workers := []func(context.Context) error{
		func(context.Context) error {
			close(firstReported)
			return errFirst
		},
		func(ctx context.Context) error {
			// Wait until the first error has been reported, then also fail.
			<-firstReported
			return errSecond
		},
	}

	err := runCoordinated(context.Background(), workers...)
	assert.ErrorIs(t, err, errFirst)
	assert.NotErrorIs(t, err, errSecond)
}

// TestRunCoordinated_ParentCancelPropagates verifies that canceling the
// parent context stops the workers and runCoordinated returns nil.
func TestRunCoordinated_ParentCancelPropagates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	peerDone := make(chan struct{})

	workers := []func(context.Context) error{
		func(ctx context.Context) error {
			<-ctx.Done()
			close(peerDone)
			return nil
		},
		func(ctx context.Context) error {
			<-ctx.Done()
			return nil
		},
	}

	done := make(chan error, 1)
	go func() { done <- runCoordinated(ctx, workers...) }()

	cancel()

	select {
	case <-peerDone:
	case <-time.After(2 * time.Second):
		t.Fatal("workers were not canceled by parent cancel")
	}

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("runCoordinated did not return after parent cancel")
	}
}

// writeTempConfig writes a minimal admin config with the given system addr and
// returns its path.
func writeTempConfig(t *testing.T, addr int) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, utils.ConfigFile)
	content := "system:\n  addr: " + strconv.Itoa(addr) + "\n"
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write temp config: %v", err)
	}
	return path
}

// TestViper_NonDefaultPathHonored is the P1 test: an explicitly passed config
// path must be the one loaded, not the default config.yaml.
func TestViper_NonDefaultPathHonored(t *testing.T) {
	// Ensure no env override interferes.
	t.Setenv(utils.ConfigEnv, "")

	path := writeTempConfig(t, 9999)
	vp, err := Viper(path)
	assert.NoError(t, err)
	assert.NotNil(t, vp)
	assert.Equal(t, 9999, global.CONFIG.System.Addr, "config from the explicit path must be loaded")

	// Sanity: the loaded file is the one we wrote, not the default.
	assert.Equal(t, path, vp.ConfigFileUsed())
}

// TestViper_MissingFileReturnsError verifies that a non-existent explicit path
// is surfaced as an error rather than silently falling back to a default.
func TestViper_MissingFileReturnsError(t *testing.T) {
	t.Setenv(utils.ConfigEnv, "")

	_, err := Viper(filepath.Join(t.TempDir(), "does-not-exist.yaml"))
	assert.Error(t, err)
}

// TestRunXDSServer_PortBindFailure is the P0 test: when the xDS port is already
// bound, runXDSServer must return the listen error instead of hanging.
func TestRunXDSServer_PortBindFailure(t *testing.T) {
	// Occupy a free port on all interfaces so the xDS server — which listens
	// on :<port> (all interfaces) — cannot bind it. Binding 127.0.0.1 only
	// would NOT conflict with an all-interfaces bind on macOS.
	occupier, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("occupy port: %v", err)
	}
	defer occupier.Close()
	occupiedPort := uint(occupier.Addr().(*net.TCPAddr).Port)

	// Build a real xDS server backed by an empty snapshot cache, mirroring
	// StartxDsServer. The listen call fails before Serve ever runs.
	snap := cache.NewSnapshotCache(false, cache.IDHash{}, logger.GetLogger())
	srv := envoyServer.NewServer(context.Background(), snap, nil)

	done := make(chan error, 1)
	go func() { done <- runXDSServer(context.Background(), srv, occupiedPort) }()

	select {
	case err := <-done:
		assert.Error(t, err, "port already in use must surface as an error")
	case <-time.After(2 * time.Second):
		t.Fatal("runXDSServer hung instead of returning the bind error")
	}
}
