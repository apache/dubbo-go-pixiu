//go:build !windows

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

package cmd

import (
	"bytes"
	"context"
	"errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
	"testing"
	"time"
)

import (
	"github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	_ "github.com/apache/dubbo-go-pixiu/pkg/listener/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const shutdownProbeEnv = "PIXIU_SHUTDOWN_PROBE"

func writeShutdownMarker(path, value string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(value); err != nil {
		panic(err)
	}
}

type shutdownProbePlugin struct{ path string }

func (*shutdownProbePlugin) Kind() string { return "shutdown-probe" }
func (p *shutdownProbePlugin) CreateAdapter(*model.Adapter) (adapter.Adapter, error) {
	return &shutdownProbeAdapter{path: p.path}, nil
}

type shutdownProbeAdapter struct{ path string }

func (*shutdownProbeAdapter) Config() any  { return &struct{}{} }
func (*shutdownProbeAdapter) Apply() error { return nil }
func (*shutdownProbeAdapter) Start()       {}
func (a *shutdownProbeAdapter) Stop()      { writeShutdownMarker(a.path, "adapter\n") }

// Exercise the real signal handler and os.Exit in a separate process: checking
// only stop() directly would miss a hook that is never reached during shutdown.
func TestGatewaySignalShutdown(t *testing.T) {
	if os.Getenv(shutdownProbeEnv) == "child" {
		port, err := strconv.Atoi(os.Getenv("PIXIU_SHUTDOWN_PORT"))
		require.NoError(t, err)
		path := os.Getenv("PIXIU_SHUTDOWN_MARKER")
		patches := gomonkey.ApplyMethod(&config.ConfigManager{}, "Close", func(*config.ConfigManager) {
			writeShutdownMarker(path, "config-start\n")
			if os.Getenv("PIXIU_SHUTDOWN_BLOCK_CLOSE") == "1" {
				select {}
			}
			writeShutdownMarker(path, "config-closed\n")
		})
		defer patches.Reset()
		adapter.RegisterAdapterPlugin(&shutdownProbePlugin{path: path})
		bootstrap := &model.Bootstrap{StaticResources: model.StaticResources{
			ShutdownConfig: &model.ShutdownConfig{Timeout: "2s"},
			Adapters:       []*model.Adapter{{Name: "shutdown-probe", Config: map[string]any{}}},
			Listeners:      []*model.Listener{{Name: "shutdown-probe", Protocol: model.ProtocolTypeHTTP, ProtocolStr: "HTTP", Address: model.Address{SocketAddress: model.SocketAddress{Address: "127.0.0.1", Port: port}}}},
		}}
		config.SetBootstrap(bootstrap)
		d := &DefaultDeployer{bootstrap: bootstrap, configManger: config.NewConfigManger()}
		require.NoError(t, d.start())
		t.Fatal("gateway returned without receiving a signal")
	}
	for _, blocked := range []bool{false, true} {
		name := "cleanup precedes process exit"
		if blocked {
			name = "timeout still bounds blocked cleanup"
		}
		t.Run(name, func(t *testing.T) {
			socket, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			port := socket.Addr().(*net.TCPAddr).Port
			require.NoError(t, socket.Close())
			path := filepath.Join(t.TempDir(), "closed")
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			command := exec.CommandContext(ctx, os.Args[0], "-test.run=^TestGatewaySignalShutdown$")
			command.Env = append(os.Environ(), shutdownProbeEnv+"=child", "PIXIU_SHUTDOWN_PORT="+strconv.Itoa(port), "PIXIU_SHUTDOWN_MARKER="+path)
			if blocked {
				command.Env = append(command.Env, "PIXIU_SHUTDOWN_BLOCK_CLOSE=1")
			}
			var output bytes.Buffer
			command.Stdout = &output
			command.Stderr = &output
			require.NoError(t, command.Start())
			done := make(chan struct{})
			var waitErr error
			go func() { waitErr = command.Wait(); close(done) }()
			t.Cleanup(func() { _ = command.Process.Kill(); <-done })
			// A bound listener proves startup installed the signal handler. Readiness
			// is checked across processes, avoiding unsynchronized global-state reads.
			ticker := time.NewTicker(20 * time.Millisecond)
			defer ticker.Stop()
		ready:
			for {
				select {
				case <-done:
					t.Fatalf("gateway exited before readiness: %v\n%s", waitErr, output.String())
				case <-ctx.Done():
					t.Fatal("gateway did not become ready")
				case <-ticker.C:
					connection, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(port)), 50*time.Millisecond)
					if err == nil {
						_ = connection.Close()
						break ready
					}
				}
			}
			require.NoError(t, command.Process.Signal(syscall.SIGTERM))
			select {
			case <-done:
			case <-ctx.Done():
				t.Fatal("signal shutdown did not finish")
			}
			require.NoError(t, ctx.Err(), output.String())
			if blocked && waitErr != nil {
				var exitErr *exec.ExitError
				require.True(t, errors.As(waitErr, &exitErr), output.String())
				require.Equal(t, 1, exitErr.ExitCode(), output.String())
			} else {
				require.NoError(t, waitErr, output.String())
			}
			marker, err := os.ReadFile(path)
			require.NoError(t, err, output.String())
			expected := "adapter\nconfig-start\n"
			if !blocked {
				expected += "config-closed\n"
			}
			require.Equal(t, expected, string(marker), output.String())
		})
	}
}
