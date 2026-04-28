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
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

import (
	"gopkg.in/yaml.v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	reloadMutex    sync.Mutex
	configPath     string
	reloadServer   *http.Server
	reloadServerMu sync.Mutex
	reloadSecret   string // shared secret for authentication
)

func SetConfigPath(path string) {
	configPath = path
}

func SetReloadSecret(secret string) {
	reloadSecret = secret
}

// ReloadHandler handles HTTP reload requests
type ReloadHandler struct{}

// checkAuth validates the request authentication
func checkAuth(r *http.Request) bool {
	// Allow localhost without authentication
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && (host == "127.0.0.1" || host == "::1" || host == "localhost") {
		return true
	}

	// Check shared secret if configured
	if reloadSecret != "" {
		token := r.Header.Get("X-Reload-Token")
		return token == reloadSecret
	}

	// If no secret configured and not localhost, deny
	return false
}

// ServeHTTP handles the reload HTTP request
func (h *ReloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	if !checkAuth(r) {
		logger.Warnf("Unauthorized reload request from %s", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	logger.Info("Received reload request via HTTP")

	var err error
	// Try to read from body first (handles chunked encoding where ContentLength == -1)
	// If body is empty, fallback to file reload
	if r.Body != nil {
		content, readErr := io.ReadAll(r.Body)

		if readErr != nil {
			logger.Errorf("Failed to read request body: %v", readErr)
			http.Error(w, fmt.Sprintf("Failed to read request body: %v", readErr), http.StatusBadRequest)
			return
		}

		if len(content) > 0 {
			// Body has content, use it
			reloadMutex.Lock()
			err = reloadFromYAML(content)
			reloadMutex.Unlock()
		} else {
			// Body is empty, fallback to file reload
			err = triggerConfigReload()
		}
	} else {
		err = triggerConfigReload()
	}

	if err != nil {
		logger.Errorf("Reload failed: %v", err)
		http.Error(w, fmt.Sprintf("Reload failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeJSONResponse(w, "success", "Configuration reloaded successfully")
	logger.Info("Reload completed successfully")
}

func writeJSONResponse(w http.ResponseWriter, status, message string) {
	response := map[string]string{
		"status":  status,
		"message": message,
		"time":    time.Now().Format(time.RFC3339),
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("Failed to encode response: %v", err)
	}
}

// HealthHandler handles health check requests
type HealthHandler struct{}

// ServeHTTP handles the health check HTTP request
func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !checkAuth(r) {
		logger.Warnf("Unauthorized health check request from %s", r.RemoteAddr)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	writeJSONResponse(w, "healthy", "")
}

// triggerConfigReload reloads configuration from file
func triggerConfigReload() error {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()

	if configPath == "" {
		return fmt.Errorf("config path not set")
	}

	logger.Infof("Reloading configuration from: %s", configPath)

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	return reloadFromYAML(content)
}

// reloadFromYAML performs the actual reload from YAML content
func reloadFromYAML(content []byte) error {
	newConfig := &model.Bootstrap{}
	if err := yaml.Unmarshal(content, newConfig); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	if err := config.Adapter(newConfig); err != nil {
		return fmt.Errorf("failed to adapt config: %w", err)
	}

	oldConfig := config.GetBootstrap()
	if oldConfig == nil {
		return fmt.Errorf("current config is nil")
	}

	logger.Infof("Old config has %d listeners, new config has %d listeners",
		len(oldConfig.StaticResources.Listeners), len(newConfig.StaticResources.Listeners))

	// Execute reloaders serially with CheckUpdate validation
	for _, reloader := range coordinator.reloaders {
		// Check if update is needed before executing reload
		if !reloader.CheckUpdate(oldConfig, newConfig) {
			logger.Debugf("No update needed for %T", reloader)
			continue
		}

		logger.Infof("Triggering reload for %T", reloader)
		if err := reloader.HotReload(oldConfig, newConfig); err != nil {
			logger.Errorf("Hot reload failed for %T: %v", reloader, err)
			return fmt.Errorf("reload failed for %T: %w", reloader, err)
		}
	}

	config.SetBootstrap(newConfig)
	logger.Info("Configuration reloaded successfully and global config updated")
	return nil
}

// StartReloadServer starts the HTTP server for reload endpoint
// port: the port to listen on (0 means use default 8888)
// secret: optional shared secret for authentication (empty means localhost-only)
func StartReloadServer(port int, secret string) error {
	reloadServerMu.Lock()
	defer reloadServerMu.Unlock()

	if reloadServer != nil {
		return fmt.Errorf("reload server already running")
	}

	if port <= 0 {
		port = 8888 // default port
	}

	if secret != "" {
		reloadSecret = secret
	}

	addr := fmt.Sprintf(":%d", port)
	logger.Infof("Starting reload HTTP server on %s", addr)

	// Try to bind the port first to catch errors early
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to bind port %d: %w", port, err)
	}

	mux := http.NewServeMux()
	mux.Handle("/-/reload", &ReloadHandler{})
	mux.Handle("/-/health", &HealthHandler{})

	reloadServer = &http.Server{
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start server in goroutine
	go func() {
		if err := reloadServer.Serve(listener); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Reload HTTP server failed: %v", err)
		}
	}()

	logger.Infof("Reload HTTP server started successfully on port %d", port)
	return nil
}

// StopReloadServer gracefully stops the reload HTTP server
func StopReloadServer(timeout time.Duration) error {
	reloadServerMu.Lock()
	defer reloadServerMu.Unlock()

	if reloadServer == nil {
		return nil
	}

	logger.Info("Stopping reload HTTP server")

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err := reloadServer.Shutdown(ctx)
	reloadServer = nil

	if err != nil {
		logger.Errorf("Error stopping reload server: %v", err)
		return err
	}

	logger.Info("Reload HTTP server stopped successfully")
	return nil
}
