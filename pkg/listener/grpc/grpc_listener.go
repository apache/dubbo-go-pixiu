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

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

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
	defaultStartupWait     = 100 * time.Millisecond
	defaultShutdownTimeout = 5 * time.Second
)

// GrpcListenerService implements the ListenerService interface for gRPC
type GrpcListenerService struct {
	listener.BaseListenerService
	server          *grpc.Server
	listener        net.Listener
	gShutdownConfig *listener.ListenerGracefulShutdownConfig
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

	// Build server options with interceptors
	opts := buildGrpcServerOptions(grpcConfig)
	opts = append(opts, grpc.UnaryInterceptor(ls.unaryInterceptor))
	opts = append(opts, grpc.StreamInterceptor(ls.streamInterceptor))

	// Create and configure gRPC server
	server := grpc.NewServer(opts...)
	registerProxyServices(server)
	ls.server = server

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

	logger.Infof("gRPC listener starting at %s", address)
	ls.logConfiguration()

	// Start server in a goroutine
	go ls.serveGrpc(listener)

	// Wait briefly to ensure server starts
	time.Sleep(defaultStartupWait)
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

// logConfiguration logs the current gRPC server configuration
func (ls *GrpcListenerService) logConfiguration() {
	if grpcConfig, ok := ls.Config.Config.(model.GrpcConfig); ok {
		logger.Infof("gRPC server config: MaxRecvSize=%dMB, MaxSendSize=%dMB, TLS=%t",
			grpcConfig.MaxReceiveMessageSize/(1024*1024),
			grpcConfig.MaxSendMessageSize/(1024*1024),
			grpcConfig.EnableTLS)
	}
}

// Close stops the gRPC server and closes the listener
func (ls *GrpcListenerService) Close() error {
	// Stop gRPC server
	if ls.server != nil {
		ls.server.Stop()
	}

	// Close network listener
	if ls.listener != nil {
		return ls.listener.Close()
	}

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
func buildGrpcServerOptions(config *model.GrpcConfig) []grpc.ServerOption {
	var opts []grpc.ServerOption

	// Set message size limits
	opts = append(opts, grpc.MaxRecvMsgSize(config.MaxReceiveMessageSize))
	opts = append(opts, grpc.MaxSendMsgSize(config.MaxSendMessageSize))

	// Configure keepalive parameters
	configureKeepalive(config, &opts)

	// Configure compression if enabled
	if config.EnableCompression {
		// gRPC supports gzip compression by default
		// Additional compression algorithms can be added as needed
	}

	// Configure TLS if enabled
	configureTLS(config, &opts)

	return opts
}

// configureKeepalive sets up keepalive parameters for the gRPC server
func configureKeepalive(config *model.GrpcConfig, opts *[]grpc.ServerOption) {
	idleTimeout, err1 := time.ParseDuration(config.IdleTimeout)
	maxConnectionAge, err2 := time.ParseDuration(config.MaxConnectionAge)

	if err1 == nil && err2 == nil {
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
func (s *RPCStreamImpl) SendMsg(m interface{}) error {
	return s.ServerStream.SendMsg(m)
}

// RecvMsg implements model.RPCStream interface
func (s *RPCStreamImpl) RecvMsg(m interface{}) error {
	return s.ServerStream.RecvMsg(m)
}

// unaryInterceptor handles unary RPC calls
func (ls *GrpcListenerService) unaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	start := time.Now()
	logger.Debugf("gRPC unary request: %s", info.FullMethod)

	// Check if server is shutting down
	if ls.gShutdownConfig.RejectRequest {
		logger.Warnf("Rejecting gRPC unary request %s during shutdown", info.FullMethod)
		return nil, errors.New("server is shutting down")
	}

	// Track active request count
	ls.gShutdownConfig.ActiveCount++
	defer func() {
		ls.gShutdownConfig.ActiveCount--
	}()

	// Process request through filter chain
	result, err := ls.FilterChain.OnUnaryRPC(ctx, info.FullMethod, req)

	// Log request completion
	duration := time.Since(start)
	if err != nil {
		logger.Errorf("gRPC unary request %s failed: %v (took %v)", info.FullMethod, err, duration)
	} else {
		logger.Debugf("gRPC unary request %s completed (took %v)", info.FullMethod, duration)
	}

	return result, err
}

// streamInterceptor handles streaming RPC calls
func (ls *GrpcListenerService) streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	start := time.Now()
	logger.Debugf("gRPC stream request: %s (ClientStream: %v, ServerStream: %v)",
		info.FullMethod, info.IsClientStream, info.IsServerStream)

	// Check if server is shutting down
	if ls.gShutdownConfig.RejectRequest {
		logger.Warnf("Rejecting gRPC stream request %s during shutdown", info.FullMethod)
		return errors.New("server is shutting down")
	}

	// Track active request count
	ls.gShutdownConfig.ActiveCount++
	defer func() {
		ls.gShutdownConfig.ActiveCount--
	}()

	// Create stream wrapper and info
	stream := &RPCStreamImpl{ServerStream: ss}
	streamInfo := &model.RPCStreamInfo{
		FullMethod:     info.FullMethod,
		IsClientStream: info.IsClientStream,
		IsServerStream: info.IsServerStream,
	}

	// Process stream through filter chain
	err := ls.FilterChain.OnStreamRPC(stream, streamInfo)

	// Log request completion
	duration := time.Since(start)
	if err != nil {
		logger.Errorf("gRPC stream request %s failed: %v (took %v)", info.FullMethod, err, duration)
	} else {
		logger.Debugf("gRPC stream request %s completed (took %v)", info.FullMethod, duration)
	}

	return err
}
