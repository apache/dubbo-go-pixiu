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
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/global"
	"github.com/apache/dubbo-go-pixiu/admin/utils"
)

// errSentinel is a distinct error used to assert exact propagation.
var errSentinel = errors.New("sentinel startup error")

type shutdownContextServer struct {
	shutdownContextErr chan error
	stopped            chan struct{}
}

func (s *shutdownContextServer) ListenAndServe() error {
	<-s.stopped
	return http.ErrServerClosed
}

func (s *shutdownContextServer) Shutdown(ctx context.Context) error {
	s.shutdownContextErr <- ctx.Err()
	close(s.stopped)
	return nil
}

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

func TestRunHTTPServerUsesLiveContextForGracefulShutdown(t *testing.T) {
	previousLogger := global.LOG
	global.LOG = zap.NewNop()
	t.Cleanup(func() { global.LOG = previousLogger })

	server := &shutdownContextServer{
		shutdownContextErr: make(chan error, 1),
		stopped:            make(chan struct{}),
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- runHTTPServer(ctx, server, ":0") }()

	cancel()

	select {
	case err := <-server.shutdownContextErr:
		assert.NoError(t, err, "Shutdown must receive a live context")
	case <-time.After(2 * time.Second):
		t.Fatal("HTTP server shutdown was not called")
	}

	select {
	case err := <-done:
		assert.NoError(t, err)
	case <-time.After(2 * time.Second):
		t.Fatal("runHTTPServer did not return after shutdown")
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

func TestViper_InvalidTypedValueReturnsError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "invalid.yaml")
	if err := os.WriteFile(path, []byte("system:\n  addr: not-an-integer\n"), 0o600); err != nil {
		t.Fatalf("write invalid config: %v", err)
	}

	_, err := Viper(path)
	assert.Error(t, err)
}
