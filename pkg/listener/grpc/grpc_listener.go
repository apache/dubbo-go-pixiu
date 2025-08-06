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

package grpc

import (
	"context"
	"crypto/tls"
	"net"
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/codec/grpc/passthrough"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Register the gRPC listener service factory
func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeGRPC, newGrpcListenerService)
}

// Constants for gRPC listener
const (
	defaultTLSTimeout      = 20 * time.Second
	defaultGracePeriod     = 5 * time.Second
	defaultMinKeepalive    = 30 * time.Second
	defaultShutdownTimeout = 5 * time.Second
)

// GrpcListenerService implements the ListenerService interface for gRPC
type GrpcListenerService struct {
	listener.BaseListenerService
	server          *grpc.Server
	listener        net.Listener
	grpcConfig      *model.GrpcConfig
	gShutdownConfig *listener.ListenerGracefulShutdownConfig
	closeOnce       sync.Once
}

// newGrpcListenerService creates a new gRPC listener service
func newGrpcListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	// Create network filter chain
	fc := filterchain.CreateNetworkFilterChain(lc.FilterChain)

	// Initialize service with base configuration
	ls := &GrpcListenerService{
		BaseListenerService: listener.BaseListenerService{
			Config:      lc,
			FilterChain: fc,
		},
		gShutdownConfig: &listener.ListenerGracefulShutdownConfig{},
	}

	// Parse gRPC specific configuration
	grpcConfig := model.MapInGrpcStruct(lc.Config)
	ls.grpcConfig = grpcConfig

	// Build server options with a proxy handler for unknown services
	opts := buildGrpcServerOptions(grpcConfig, ls)

	// Create and configure gRPC server
	server := grpc.NewServer(opts...)
	registerProxyServices(server)
	ls.server = server

	ls.logConfiguration()

	return ls, nil
}

// Start initializes and starts the gRPC server
func (ls *GrpcListenerService) Start() error {
	address := ls.Config.Address.SocketAddress.GetAddress()

	// Create network listener
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return errors.Wrapf(err, "failed to listen on %s", address)
	}
	ls.listener = listener

	//ls.logConfiguration()

	// Start server in a goroutine
	go ls.serveGrpc(listener)

	// The listener is ready as soon as net.Listen succeeds and the goroutine is running.
	// We don't need a fixed sleep here as it's unreliable.
	logger.Infof("gRPC listener successfully started at %s", address)

	return nil
}

// serveGrpc runs the gRPC server on the provided listener
func (ls *GrpcListenerService) serveGrpc(listener net.Listener) {
	logger.Info("gRPC server starting to serve...")
	if err := ls.server.Serve(listener); err != nil {
		logger.Errorf("gRPC server serve error: %v", err)
	} else {
		logger.Info("gRPC server stopped gracefully")
	}
}

// proxyStreamHandler handles all unknown gRPC streams and forwards them through the filter chain.
// This is the core of the gRPC proxy functionality.
func (ls *GrpcListenerService) proxyStreamHandler(srv any, ss grpc.ServerStream) error {
	start := time.Now()

	// The full method name is available in the stream's context.
	fullMethod, ok := grpc.MethodFromServerStream(ss)
	if !ok {
		return errors.New("could not determine method from stream")
	}

	// Check if server is shutting down
	if ls.gShutdownConfig.RejectRequest {
		logger.Warnf("Rejecting gRPC stream request %s during shutdown", fullMethod)
		return errors.New("server is shutting down")
	}

	// Track active request count
	ls.gShutdownConfig.ActiveCount++
	defer func() {
		ls.gShutdownConfig.ActiveCount--
	}()

	// Since we don't have StreamInfo here, we must rely on the filter chain to get it if needed.
	// For a pure proxy, we just need to forward the stream.
	stream := &RPCStreamImpl{ServerStream: ss}

	// The filter chain needs RPCStreamInfo. Let's create a basic one.
	// Since we cannot determine the exact stream type (unary, client-stream, server-stream)
	// at this transparent proxy layer without parsing descriptors, we assume it's a
	// bidirectional stream by default. The downstream filters can re-infer this if needed.
	streamInfo := &model.RPCStreamInfo{
		FullMethod:     fullMethod,
		IsClientStream: true, // Assume client streaming
		IsServerStream: true, // Assume server streaming
	}

	// Process stream through filter chain
	err := ls.FilterChain.OnStreamRPC(stream, streamInfo)

	// Log request completion
	duration := time.Since(start)
	if err != nil {
		logger.Errorf("gRPC stream request %s failed: %v (took %v)", fullMethod, err, duration)
	} else {
		logger.Debugf("gRPC stream for %s completed in %v", fullMethod, duration)
	}

	return err
}

// logConfiguration logs the current gRPC server configuration
func (ls *GrpcListenerService) logConfiguration() {
	if grpcConfig := ls.grpcConfig; grpcConfig != nil {
		logger.Infof("gRPC server config: MaxRecvSize=%dMB, MaxSendSize=%dMB, TLS=%t",
			grpcConfig.MaxReceiveMessageSize/(1024*1024),
			grpcConfig.MaxSendMessageSize/(1024*1024),
			grpcConfig.EnableTLS)
	}
}

// cleanup closes the filter chain and other resources. It's designed to be called once.
func (ls *GrpcListenerService) cleanup() {
	ls.closeOnce.Do(func() {
		logger.Info("Cleaning up gRPC listener resources...")
		if ls.FilterChain != nil {
			if err := ls.FilterChain.Close(); err != nil {
				logger.Warnf("Error closing filter chain: %v", err)
			}
		}
	})
}

// Close stops the gRPC server immediately. It is a hard stop.
func (ls *GrpcListenerService) Close() error {
	logger.Info("Forcefully closing gRPC listener...")
	if ls.server != nil {
		ls.server.Stop()
	}
	ls.cleanup()
	return nil
}

// ShutDown gracefully shuts down the gRPC listener
func (ls *GrpcListenerService) ShutDown(wg any) error {
	logger.Info("gRPC listener shutdown initiated")
	waitGroup := wg.(*sync.WaitGroup)
	defer waitGroup.Done()

	// Get shutdown timeout from configuration
	timeout := config.GetBootstrap().GetShutdownConfig().GetTimeout()
	if timeout <= 0 {
		logger.Info("No shutdown timeout configured, stopping immediately")
		if ls.server != nil {
			ls.server.Stop()
		}
		return nil
	}

	// Start graceful shutdown
	ls.gShutdownConfig.RejectRequest = true
	deadline := time.Now().Add(timeout)
	logger.Infof("Graceful shutdown initiated with timeout: %v", timeout)

	// Wait for active requests to complete
	ls.waitForActiveRequests(deadline)

	// Gracefully stop the server
	ls.gracefulStopServer()

	// Clean up resources
	ls.cleanup()

	logger.Info("gRPC listener shutdown completed")
	return nil
}

// waitForActiveRequests waits for active requests to complete until deadline
func (ls *GrpcListenerService) waitForActiveRequests(deadline time.Time) {
	for time.Now().Before(deadline) && ls.gShutdownConfig.ActiveCount > 0 {
		time.Sleep(100 * time.Millisecond)
		logger.Infof("waiting for active gRPC invocation count = %d", ls.gShutdownConfig.ActiveCount)
	}

	if ls.gShutdownConfig.ActiveCount > 0 {
		logger.Warnf("Shutdown timeout reached, forcing stop with %d active requests", ls.gShutdownConfig.ActiveCount)
	} else {
		logger.Info("All active requests completed, proceeding with graceful shutdown")
	}
}

// gracefulStopServer attempts to gracefully stop the gRPC server with timeout
func (ls *GrpcListenerService) gracefulStopServer() {
	if ls.server == nil {
		return
	}

	// Use goroutine for graceful shutdown to avoid blocking
	done := make(chan struct{})
	go func() {
		ls.server.GracefulStop()
		close(done)
	}()

	// Wait for graceful shutdown or timeout
	select {
	case <-done:
		logger.Info("gRPC server gracefully stopped")
	case <-time.After(defaultShutdownTimeout):
		logger.Warn("Graceful stop timeout, forcing stop")
		ls.server.Stop()
	}
}

// Refresh updates the filter chain configuration
func (ls *GrpcListenerService) Refresh(c model.Listener) error {
	fc := filterchain.CreateNetworkFilterChain(c.FilterChain)
	ls.FilterChain = fc
	return nil
}

// buildGrpcServerOptions creates gRPC server options from config
func buildGrpcServerOptions(config *model.GrpcConfig, ls *GrpcListenerService) []grpc.ServerOption {
	var opts []grpc.ServerOption

	// Force the server to use a passthrough codec for all requests.
	// This is essential for a transparent proxy, as the server does not know the message types.
	// It allows the server to receive raw bytes, which can then be forwarded by the proxy filter.
	opts = append(opts, grpc.ForceServerCodec(passthrough.Codec{}))

	// Set the handler for unknown services to the proxy handler
	opts = append(opts, grpc.UnknownServiceHandler(ls.proxyStreamHandler))

	// Set message size limits
	opts = append(opts, grpc.MaxRecvMsgSize(config.MaxReceiveMessageSize))
	opts = append(opts, grpc.MaxSendMsgSize(config.MaxSendMessageSize))

	// Configure keepalive parameters
	configureKeepalive(config, &opts)

	// Configure TLS if enabled
	configureTLS(config, &opts)

	return opts
}

// configureKeepalive sets up keepalive parameters for the gRPC server
func configureKeepalive(config *model.GrpcConfig, opts *[]grpc.ServerOption) {
	idleTimeout, err := time.ParseDuration(config.IdleTimeout)
	if config.IdleTimeout != "" && err != nil {
		logger.Warnf("Invalid gRPC idle_timeout format: %s, keepalive disabled.", config.IdleTimeout)
		return
	}
	if config.IdleTimeout == "" {
		// Use a default or skip if not provided
		return
	}

	maxConnectionAge, err := time.ParseDuration(config.MaxConnectionAge)
	if config.MaxConnectionAge != "" && err != nil {
		logger.Warnf("Invalid gRPC max_connection_age format: %s, using infinite.", config.MaxConnectionAge)
		maxConnectionAge = 0 // Or some other default
	}

	// Server parameters
	kasp := keepalive.ServerParameters{
		Time:                  idleTimeout,
		Timeout:               defaultTLSTimeout,
		MaxConnectionAge:      maxConnectionAge,
		MaxConnectionAgeGrace: defaultGracePeriod,
	}
	*opts = append(*opts, grpc.KeepaliveParams(kasp))

	// Enforcement policy
	kaep := keepalive.EnforcementPolicy{
		MinTime:             defaultMinKeepalive,
		PermitWithoutStream: true,
	}
	*opts = append(*opts, grpc.KeepaliveEnforcementPolicy(kaep))
}

// configureTLS sets up TLS credentials for the gRPC server if enabled
func configureTLS(config *model.GrpcConfig, opts *[]grpc.ServerOption) {
	if !config.EnableTLS || config.TLS == nil || config.TLS.CertFile == "" || config.TLS.KeyFile == "" {
		return
	}

	creds, err := loadTLSCredentials(config.TLS)
	if err != nil {
		logger.Warnf("Failed to load TLS credentials: %v, starting without TLS", err)
		return
	}

	*opts = append(*opts, grpc.Creds(creds))
	logger.Infof("gRPC server TLS enabled with cert: %s", config.TLS.CertFile)
}

// loadTLSCredentials loads TLS credentials from certificate and key files
func loadTLSCredentials(tlsConfig *model.TLSConfig) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(tlsConfig.CertFile, tlsConfig.KeyFile)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	return credentials.NewTLS(config), nil
}

// registerProxyServices registers services with the gRPC server
func registerProxyServices(server *grpc.Server) {
	// Enable gRPC reflection for debugging and service discovery
	reflection.Register(server)
	logger.Info("gRPC reflection service registered")

	// Actual request handling is done in interceptors
	logger.Info("gRPC proxy service configured with interceptors")
}

// RPCStreamImpl implements the model.RPCStream interface
type RPCStreamImpl struct {
	grpc.ServerStream
}

// Context implements model.RPCStream interface
func (s *RPCStreamImpl) Context() context.Context {
	return s.ServerStream.Context()
}

// SendMsg implements model.RPCStream interface
func (s *RPCStreamImpl) SendMsg(m any) error {
	return s.ServerStream.SendMsg(m)
}

// RecvMsg implements model.RPCStream interface
func (s *RPCStreamImpl) RecvMsg(m any) error {
	return s.ServerStream.RecvMsg(m)
}
