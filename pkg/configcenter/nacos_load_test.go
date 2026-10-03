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

package configcenter

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"sync"
	"testing"
)

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap/zapcore"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// isNacosRunning checks whether the Nacos server is running.
// It returns true if Nacos is running, otherwise false.
func isNacosRunning(t *testing.T) bool {
	t.Helper()
	_, err := getNacosConfigClient(getBootstrap())
	return err == nil
}

// TestNewNacosConfig tests the creation of a new Nacos configuration.
// If Nacos is not running, the test is skipped.
func TestNewNacosConfig(t *testing.T) {
	if !isNacosRunning(t) {
		t.Skip("Nacos is not running, skipping the test.")
		return
	}

	Convey("Test NewNacosConfig", t, func() {
		cfg := getBootstrap()

		// Test successful creation of NacosConfig.
		_, err := NewNacosConfig(cfg)
		So(err, ShouldBeNil)

		// Test creation failure when Nacos server configurations are missing.
		cfg.Nacos.ServerConfigs = nil
		_, err = NewNacosConfig(cfg)
		So(err, ShouldNotBeNil)
	})
}

// TestNacosConfig_onChange tests the onChange method of NacosConfig.
func TestNacosConfig_onChange(t *testing.T) {
	Convey("TestNacosConfig_onChange", t, func() {
		cfg := getBootstrap()
		c, err := NewNacosConfig(cfg)
		So(err, ShouldBeNil)

		client, ok := c.(*NacosConfig)
		So(ok, ShouldBeTrue)

		// Verify the current working directory.
		wd, err := os.Getwd()
		So(err, ShouldBeNil)

		paths := strings.Split(wd, "/")
		So(paths[len(paths)-1], ShouldEqual, "configcenter")

		// Open the configuration file for testing.
		file, err := os.Open(fmt.Sprintf("/%s/configs/conf.yaml", path.Join(paths[:len(paths)-2]...)))
		So(err, ShouldBeNil)
		defer func() { So(file.Close(), ShouldBeNil) }()

		conf, err := io.ReadAll(file)
		So(err, ShouldBeNil)

		Convey("Test onChange with valid input", func() {
			So(client.remoteConfig, ShouldBeNil)
			client.onChange(Namespace, Group, DataId, string(conf))
			So(client.remoteConfig, ShouldNotBeNil)
		})

		Convey("Test onChange with empty input", func() {
			// Suppress logs during this test.
			logger.SetLoggerLevel(zapcore.FatalLevel)

			client.remoteConfig = nil
			client.onChange(Namespace, Group, DataId, "")
			So(client.remoteConfig, ShouldBeNil)

			// Restore the logger level.
			logger.SetLoggerLevel(zapcore.InfoLevel)
		})
	})
}

// closeRecorderConfig embeds the SDK config-client interface and records
// CloseClient calls so NacosConfig.Close can be asserted without a live Nacos
// server.
type closeRecorderConfig struct {
	config_client.IConfigClient

	mu         sync.Mutex
	closeCount int
}

func (c *closeRecorderConfig) CloseClient() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closeCount++
}

func (c *closeRecorderConfig) closeCalls() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.closeCount
}

// TestNacosConfig_Close verifies the v2 close path: Close must delegate to the
// SDK config client's CloseClient so the gRPC connection is released on
// shutdown. See AlexStocks' [P1] review on PR #982.
func TestNacosConfig_Close(t *testing.T) {
	t.Run("closes the underlying client", func(t *testing.T) {
		recorder := &closeRecorderConfig{}
		cfg := &NacosConfig{client: recorder}

		assert.Equal(t, 0, recorder.closeCalls(), "client must not be closed before Close")

		cfg.Close()

		assert.Equal(t, 1, recorder.closeCalls(),
			"Close must delegate to the underlying config client's CloseClient")
	})

	t.Run("no-op when client is nil", func(t *testing.T) {
		cfg := &NacosConfig{client: nil}
		assert.NotPanics(t, func() { cfg.Close() })
	})
}
