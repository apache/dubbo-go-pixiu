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

package hotreload

import (
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// HotReloader defines the interface for the HotReload module.
type HotReloader interface {
	// CheckUpdate checks if the configuration has changed and needs to be reloaded.
	CheckUpdate(oldConfig, newConfig *model.Bootstrap) bool

	// HotReload performs the hot reload of the configuration.
	HotReload(oldConfig, newConfig *model.Bootstrap) error
}

// Coordinator listens for configuration file changes and notifies registered reloaders
// to perform hot reload when the configuration changes.
type Coordinator struct {
	reloaders []HotReloader         // List of registered reloaders
	boot      *model.Bootstrap      // Current configuration
	manager   *config.ConfigManager // Configuration manager
}

var coordinator = Coordinator{reloaders: []HotReloader{&LoggerReloader{}, &RouteReloader{}}}

// StartHotReload initializes the hot reload process.
// It should be called when the project starts, e.g., in cmd/gateway.go.
func StartHotReload(manager *config.ConfigManager, boot *model.Bootstrap) {
	if manager == nil || boot == nil {
		logger.Warn("ConfigManager or Bootstrap is nil, hot reload will not start")
		return
	}

	coordinator.manager = manager
	coordinator.boot = boot

	go coordinator.HotReload()
}

// HotReload periodically checks for configuration updates and triggers hot reload if changes are detected.
func (c *Coordinator) HotReload() {
	for {
		time.Sleep(constant.CheckConfigInterval)

		boot := c.manager.ViewRemoteConfig()
		if boot == nil {
			continue
		}

		c.hotReload(boot)
	}
}

// hotReload checks for configuration changes and triggers hot reload for registered reloaders.
func (c *Coordinator) hotReload(newBoot *model.Bootstrap) {
	changed := false
	wg := &sync.WaitGroup{}

	for _, reloader := range c.reloaders {
		if reloader.CheckUpdate(c.boot, newBoot) {
			changed = true
			wg.Add(1)
			go func(r HotReloader) {
				defer wg.Done()
				if err := r.HotReload(c.boot, newBoot); err != nil {
					logger.Errorf("Hot reload failed: %v", err)
				}
			}(reloader)
		}
	}

	wg.Wait()
	if changed {
		c.boot = newBoot
	}
}

// equal checks if two string slices are equal.
func equal(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
