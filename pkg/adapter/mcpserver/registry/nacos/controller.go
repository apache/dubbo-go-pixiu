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

package nacos

import (
	"context"
	"fmt"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// Retry configuration constants
	MaxRetryAttempts = 3
	RetryDelayMs     = 500
)

// McpController is the MCP server's configuration synchronizer in Nacos.
// It is responsible for discovering, watching, transforming, and applying configurations.
type McpController struct {
	client   *NacosRegistryClient
	onChange func(serverId string, cfg *McpServerConfig)
	watched  map[string]bool
	mu       sync.RWMutex

	// Track the last known server count to detect suspicious empty lists
	lastKnownServerCount int
}

// NewMcpController creates a new MCP controller
func NewMcpController(client *NacosRegistryClient, onChange func(serverId string, cfg *McpServerConfig)) *McpController {
	return &McpController{
		client:               client,
		onChange:             onChange,
		watched:              make(map[string]bool),
		lastKnownServerCount: 0,
	}
}

// Run starts the controller, periodically discovering and watching MCP services
func (c *McpController) Run(ctx context.Context, interval time.Duration) error {
	logger.Infof("Starting MCP controller with interval: %v", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run once immediately
	if err := c.reconcile(); err != nil {
		logger.Errorf("Initial reconcile failed: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Infof("MCP controller stopped")
			return nil
		case <-ticker.C:
			if err := c.reconcile(); err != nil {
				logger.Errorf("Reconcile failed: %v", err)
			}
		}
	}
}

// reconcile coordinates logic: discover services, compute diffs, bind watchers
func (c *McpController) reconcile() error {
	logger.Debugf("[dubbo-go-pixiu] nacos registry starting to list MCP servers")

	// Retrieve all MCP services with retry mechanism
	var servers []BasicMcpServerInfo
	var err error

	for attempt := 0; attempt < MaxRetryAttempts; attempt++ {
		servers, err = c.client.ListMcpServer()
		if err != nil {
			if attempt < MaxRetryAttempts-1 {
				logger.Warnf("Failed to list MCP servers (attempt %d/%d): %v, retrying...",
					attempt+1, MaxRetryAttempts, err)
				time.Sleep(time.Duration(RetryDelayMs*(attempt+1)) * time.Millisecond)
				continue
			}
			return fmt.Errorf("failed to list MCP servers after %d attempts: %w", MaxRetryAttempts, err)
		}
		break
	}

	// Empty list protection: if we previously had servers but now suddenly have none, skip cleanup
	if len(servers) == 0 && c.lastKnownServerCount > 0 {
		logger.Warnf("Detected empty server list, but previously had %d servers. Skipping cleanup to avoid false positives.",
			c.lastKnownServerCount)
		return nil
	}

	// Update known server count
	if len(servers) > 0 {
		c.lastKnownServerCount = len(servers)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	// Determine services to watch
	currentWatched := make(map[string]bool)
	for _, server := range servers {
		currentWatched[server.Id] = true

		// If not watching yet, start watching
		if !c.watched[server.Id] {
			logger.Infof("Starting to watch MCP server: %s (%s)", server.Name, server.Id)

			err := c.client.ListenToMcpServer(server.Id, c.wrapListener(server.Id))
			if err != nil {
				logger.Errorf("Failed to listen to server %s: %v", server.Id, err)
				continue
			}

			c.watched[server.Id] = true
		}
	}

	// Cancel watchers for servers that no longer exist
	for serverId := range c.watched {
		if !currentWatched[serverId] {
			logger.Infof("Stopping watch for MCP server: %s", serverId)

			err := c.client.CancelListenToServer(serverId)
			if err != nil {
				logger.Errorf("Failed to cancel listen for server %s: %v", serverId, err)
			}

			delete(c.watched, serverId)
		}
	}

	return nil
}

// wrapListener wraps a listener callback
func (c *McpController) wrapListener(serverId string) McpServerListener {
	return func(cfg *McpServerConfig) {
		logger.Infof("Received config update for server: %s", serverId)

		if c.onChange != nil {
			c.onChange(serverId, cfg)
		}
	}
}

// Close closes the controller
func (c *McpController) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Cancel all watchers
	for serverId := range c.watched {
		err := c.client.CancelListenToServer(serverId)
		if err != nil {
			logger.Errorf("Failed to cancel listen for server %s: %v", serverId, err)
		}
	}

	// Clear the watch list
	c.watched = make(map[string]bool)

	return nil
}
