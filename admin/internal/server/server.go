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

package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

import (
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/handler"
	"github.com/apache/dubbo-go-pixiu/admin/internal/wire"
)

const welcomeMessage = `
	Welcome DUBBOGO-PIXIU-ADMIN
	Default doc address: http://127.0.0.1%s/swagger/index.html
	Default running address: http://127.0.0.1:8080
`

// Server represents the admin server with all its dependencies.
type Server struct {
	app     *wire.App
	cleanup func()
}

// New creates a new Server instance using wire for dependency injection.
func New(configPath string) (*Server, error) {
	app, cleanup, err := wire.InitializeApp(wire.ConfigPath(configPath))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize app: %w", err)
	}

	return &Server{
		app:     app,
		cleanup: cleanup,
	}, nil
}

// Run starts the server and blocks until shutdown.
func (s *Server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup HTTP server
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(s.corsMiddleware())
	router.Use(s.app.Handler.Language()) // i18n language detection

	// Setup routes
	s.app.Handler.SetupRouter(router)

	// Register xDS server for instance queries
	handler.RegisterXDSServer(s.app.XDSServer)

	// Serve static files for the web UI
	router.StaticFS("/static", http.Dir("web/dist"))
	router.NoRoute(func(c *gin.Context) {
		c.File("web/dist/index.html")
	})

	// Use System.Addr for port (backward compatible with existing config)
	port := s.app.Config.System.Addr
	if port == 0 {
		port = 8080
	}
	address := fmt.Sprintf(":%d", port)

	httpServer := &http.Server{
		Addr:    address,
		Handler: router,
	}

	// Channel for errors
	errCh := make(chan error, 2)

	// Start HTTP server
	go func() {
		s.app.Logger.Info("HTTP server starting", zap.String("address", address))
		fmt.Printf(welcomeMessage, address)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Start xDS server
	go func() {
		s.app.Logger.Info("xDS server starting", zap.Uint("port", 18000))
		if err := s.app.XDSServer.Start(ctx); err != nil {
			errCh <- fmt.Errorf("xDS server error: %w", err)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		s.app.Logger.Info("shutdown signal received")
	case err := <-errCh:
		s.app.Logger.Error("server error", zap.Error(err))
		return err
	}

	// Graceful shutdown
	s.app.Logger.Info("shutting down servers...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		s.app.Logger.Error("HTTP server shutdown error", zap.Error(err))
	}

	// Close stores (wire cleanup handles etcd client)
	s.app.MySQL.Close()
	s.app.Etcd.Close()

	// Run wire cleanup
	if s.cleanup != nil {
		s.cleanup()
	}

	s.app.Logger.Info("server stopped")
	return nil
}

// corsMiddleware returns a CORS middleware handler.
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, token, username, Accept-Language")
		c.Header("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
