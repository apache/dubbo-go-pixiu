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

package main

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/require"

	"go.etcd.io/etcd/server/v3/embed"
)

import (
	config2 "github.com/apache/dubbo-go-pixiu/admin/config"
)

func TestAdminRejectsUnsafeSigningKeyBeforeLoadingConfiguration(t *testing.T) {
	command := getRootCmd()
	originalConfigPath, originalAPIConfigPath := configPath, apiConfigPath
	originalBootstrap := config2.Bootstrap
	t.Cleanup(func() {
		configPath, apiConfigPath = originalConfigPath, originalAPIConfigPath
		config2.Bootstrap = originalBootstrap
		command.SetArgs(nil)
	})
	configurationPath := filepath.Join(t.TempDir(), "admin.yaml")
	require.NoError(t, os.WriteFile(configurationPath, []byte("name: signing-key-regression\n"), 0o600))
	_, err := config2.LoadAPIConfigFromFile(configurationPath)
	require.NoError(t, err)
	config2.Bootstrap = originalBootstrap

	for _, key := range []string{"dubbo-go-pixiu", "short-key"} {
		t.Run(key, func(t *testing.T) {
			t.Setenv("DUBBOGO_PIXIU_JWT_SIGN_KEY", key)
			command.SetArgs([]string{"--config", configurationPath})
			err := command.Execute()
			require.ErrorContains(t, err, "DUBBOGO_PIXIU_JWT_SIGN_KEY")
			require.True(t, config2.Bootstrap == originalBootstrap, "configuration loaded before key validation")

			err = Start("")
			require.ErrorContains(t, err, "DUBBOGO_PIXIU_JWT_SIGN_KEY")
		})
	}
}

func TestAdminSignalShutdown(t *testing.T) {
	if signalName := os.Getenv("PIXIU_TEST_ADMIN_SIGNAL"); signalName != "" {
		testAdminSignalShutdown(t, signalName)
		return
	}
	for _, signalName := range []string{"SIGINT", "SIGTERM"} {
		t.Run(signalName, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			executable, err := os.Executable()
			require.NoError(t, err)
			cmd := exec.CommandContext(ctx, executable, "-test.run=^TestAdminSignalShutdown$")
			cmd.Env = append(os.Environ(), "PIXIU_TEST_ADMIN_SIGNAL="+signalName)
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, "%s", output)
		})
	}
}

// Use a subprocess so real OS signals and global server state stay isolated.
func testAdminSignalShutdown(t *testing.T, signalName string) {
	t.Helper()
	etcdConfig := embed.NewConfig()
	etcdConfig.Dir = t.TempDir()
	etcdConfig.LogLevel = "error"
	etcdConfig.LCUrls = []url.URL{{Scheme: "http", Host: "127.0.0.1:0"}}
	etcdConfig.ACUrls = etcdConfig.LCUrls
	etcdConfig.LPUrls = []url.URL{{Scheme: "http", Host: "127.0.0.1:0"}}
	etcdConfig.APUrls = etcdConfig.LPUrls
	etcdConfig.InitialCluster = etcdConfig.InitialClusterFromName(etcdConfig.Name)
	etcd, err := embed.StartEtcd(etcdConfig)
	require.NoError(t, err)
	defer etcd.Close()
	select {
	case <-etcd.Server.ReadyNotify():
	case <-time.After(10 * time.Second):
		t.Fatal("embedded etcd did not start")
	}

	httpListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	xdsListener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	httpPort := httpListener.Addr().(*net.TCPAddr).Port
	xdsPort := xdsListener.Addr().(*net.TCPAddr).Port
	require.NoError(t, httpListener.Close())
	require.NoError(t, xdsListener.Close())

	dir := t.TempDir()
	adminPath := filepath.Join(dir, "admin.yaml")
	apiPath := filepath.Join(dir, "api.yaml")
	require.NoError(t, os.WriteFile(adminPath, []byte(fmt.Sprintf("system:\n  addr: %d\nzap:\n  director: %s\n", httpPort, dir)), 0o600))
	require.NoError(t, os.WriteFile(apiPath, []byte(fmt.Sprintf("etcd:\n  address: %s\nxds:\n  listen_port: %d\n", etcd.Clients[0].Addr(), xdsPort)), 0o600))
	t.Setenv("DUBBOGO_PIXIU_JWT_SIGN_KEY", "admin-signal-test-key-at-least-32-bytes-long")

	addresses := []string{fmt.Sprintf("127.0.0.1:%d", httpPort), fmt.Sprintf("127.0.0.1:%d", xdsPort)}
	signalErr := make(chan error, 1)
	go func() {
		deadline := time.Now().Add(10 * time.Second)
		for _, address := range addresses {
			for {
				conn, err := net.DialTimeout("tcp", address, 100*time.Millisecond)
				if err == nil {
					_ = conn.Close()
					break
				}
				if time.Now().After(deadline) {
					signalErr <- fmt.Errorf("server %s did not start: %w", address, err)
					return
				}
				time.Sleep(10 * time.Millisecond)
			}
		}
		sig := syscall.SIGINT
		if signalName == "SIGTERM" {
			sig = syscall.SIGTERM
		}
		signalErr <- syscall.Kill(os.Getpid(), sig)
	}()
	command := getRootCmd()
	command.SetArgs([]string{"--config", adminPath, "--api-config", apiPath})
	require.NoError(t, command.Execute())
	require.NoError(t, <-signalErr)
	require.Nil(t, config2.Client, "etcd client must close after servers stop")
	for _, address := range addresses {
		conn, err := net.DialTimeout("tcp", address, time.Second)
		if conn != nil {
			_ = conn.Close()
		}
		require.Error(t, err, "command returned while %s was still listening", address)
	}
}
