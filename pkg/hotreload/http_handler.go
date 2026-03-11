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
	"encoding/json"
	"fmt"
	"io"
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
	reloadMutex sync.Mutex
	configPath  string
)

func SetConfigPath(path string) {
	configPath = path
}

// ReloadHandler handles HTTP reload requests
type ReloadHandler struct{}

// ServeHTTP handles the reload HTTP request
func (h *ReloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed, use POST", http.StatusMethodNotAllowed)
		return
	}

	logger.Info("Received reload request via HTTP")

	var err error
	if r.ContentLength > 0 {
		err = triggerConfigReloadFromBody(r)
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

// triggerConfigReloadFromBody reloads configuration from HTTP request body
func triggerConfigReloadFromBody(r *http.Request) error {
	reloadMutex.Lock()
	defer reloadMutex.Unlock()

	logger.Info("Reloading configuration from request body")

	content, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	defer func() {
		if err := r.Body.Close(); err != nil {
			logger.Errorf("Failed to close request body: %v", err)
		}
	}()

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

	wg := &sync.WaitGroup{}
	var reloadErrors []error
	errorMutex := &sync.Mutex{}

	for _, reloader := range coordinator.reloaders {
		logger.Infof("Triggering reload for %T", reloader)
		wg.Add(1)
		go func(r HotReloader) {
			defer wg.Done()
			if err := r.HotReload(oldConfig, newConfig); err != nil {
				logger.Errorf("Hot reload failed for %T: %v", r, err)
				errorMutex.Lock()
				reloadErrors = append(reloadErrors, err)
				errorMutex.Unlock()
			}
		}(reloader)
	}

	wg.Wait()

	if len(reloadErrors) > 0 {
		return fmt.Errorf("reload completed with %d errors", len(reloadErrors))
	}

	config.SetBootstrap(newConfig)
	logger.Info("Configuration reloaded successfully and global config updated")
	return nil
}

// StartReloadServer starts the HTTP server for reload endpoint
func StartReloadServer(port int) error {
	mux := http.NewServeMux()

	mux.Handle("/-/reload", &ReloadHandler{})
	mux.Handle("/-/health", &HealthHandler{})

	addr := fmt.Sprintf(":%d", port)
	logger.Infof("Starting reload HTTP server on %s", addr)

	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf("Reload HTTP server failed: %v", err)
		}
	}()

	return nil
}
