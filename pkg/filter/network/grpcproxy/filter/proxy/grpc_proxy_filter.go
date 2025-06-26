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

package proxy

import (
	"context"
	"crypto/tls"
	"io"
	"sync"
	"time"

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	grpcCtx "github.com/apache/dubbo-go-pixiu/pkg/context/grpc"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

// Constants for gRPC proxy filter
const (
	Kind                    = constant.GRPCProxyFilter
	defaultKeepAliveTime    = 10 * time.Second
	defaultKeepAliveTimeout = 5 * time.Second
	defaultConnectTimeout   = 5 * time.Second
	defaultMaxRetryCount    = 3
)

func init() {
	filter.RegisterGrpcFilterPlugin(&Plugin{})
}

type (
	// Plugin gRPC proxy plugin implementation
	Plugin struct{}

	// Config defines the configuration options for the gRPC proxy filter
	Config struct {
		EnableTLS            bool   `yaml:"enable_tls" json:"enable_tls" mapstructure:"enable_tls"`
		TLSCertFile          string `yaml:"tls_cert_file" json:"tls_cert_file" mapstructure:"tls_cert_file"`
		TLSKeyFile           string `yaml:"tls_key_file" json:"tls_key_file" mapstructure:"tls_key_file"`
		MaxConcurrentStreams uint32 `yaml:"max_concurrent_streams" json:"max_concurrent_streams" mapstructure:"max_concurrent_streams"`
		KeepAliveTime        string `yaml:"keepalive_time" json:"keepalive_time" mapstructure:"keepalive_time"`
		KeepAliveTimeout     string `yaml:"keepalive_timeout" json:"keepalive_timeout" mapstructure:"keepalive_timeout"`
		ConnectTimeout       string `yaml:"connect_timeout" json:"connect_timeout" mapstructure:"connect_timeout"`
	}

	// Filter implements the gRPC proxy filter
	Filter struct {
		Config         *Config
		clientConnPool sync.Map     // address -> *grpc.ClientConn
		mu             sync.RWMutex // protects concurrent operations
	}
)

// Kind return plugin kind
func (p Plugin) Kind() string {
	return Kind
}

// CreateFilter create gRPC proxy filter
func (p Plugin) CreateFilter(config any) (filter.GrpcFilter, error) {
	cfg, ok := config.(*Config)
	if !ok {
		return nil, errors.New("gRPC proxy filter config type error")
	}
	return &Filter{Config: cfg}, nil
}

// Config Expose the config so that Filter Manger can inject it, so it must be a pointer
func (p Plugin) Config() any {
	return &Config{}
}

// Handle processes gRPC invocation by routing to the appropriate backend
func (f *Filter) Handle(ctx *grpcCtx.GrpcContext) filter.FilterStatus {
	// Validate context
	if ctx == nil {
		logger.Error("gRPC proxy received nil context")
		return filter.Stop
	}

	// Get route information
	if ctx.Route == nil {
		ctx.SetError(errors.New("gRPC proxy missing route information"))
		return filter.Stop
	}

	clusterName := ctx.Route.Cluster
	if clusterName == "" {
		ctx.SetError(errors.New("gRPC proxy missing cluster name"))
		return filter.Stop
	}

	// Get cluster manager
	clusterManager := server.GetClusterManager()
	if clusterManager == nil {
		ctx.SetError(errors.New("gRPC proxy cluster manager not initialized"))
		return filter.Stop
	}

	// Select endpoint from cluster
	endpoint := clusterManager.PickEndpoint(clusterName, ctx)
	if endpoint == nil {
		ctx.SetError(errors.Errorf("gRPC proxy can't find endpoint in cluster: %s", clusterName))
		return filter.Stop
	}

	// Get target address
	address := endpoint.Address.GetAddress()
	if address == "" {
		ctx.SetError(errors.New("gRPC proxy got empty endpoint address"))
		return filter.Stop
	}

	logger.Debugf("gRPC proxy forwarding %s.%s to endpoint: %s",
		ctx.ServiceName, ctx.MethodName, address)

	// Process request based on stream type
	switch ctx.StreamType {
	case grpcCtx.UnaryCall:
		return f.handleUnaryCall(ctx, address)
	case grpcCtx.ClientStream:
		return f.handleClientStream(ctx, address)
	case grpcCtx.ServerStream:
		return f.handleServerStream(ctx, address)
	case grpcCtx.BidirectionalStream:
		return f.handleBidirectionalStream(ctx, address)
	default:
		logger.Warnf("Unknown stream type %d, handling as unary call", ctx.StreamType)
		return f.handleUnaryCall(ctx, address)
	}
}

// handleUnaryCall processes a unary gRPC call by forwarding it to the backend
func (f *Filter) handleUnaryCall(ctx *grpcCtx.GrpcContext, address string) filter.FilterStatus {
	start := time.Now()

	// Validate context and arguments
	if ctx.Context == nil {
		ctx.SetError(errors.New("missing context for unary call"))
		return filter.Stop
	}

	// Create full method path
	fullMethod := "/" + ctx.ServiceName + "/" + ctx.MethodName

	// Get or create connection with timeout
	callCtx, cancel := context.WithTimeout(ctx.Context, defaultConnectTimeout)
	defer cancel()

	conn, err := f.getOrCreateConnection(ctx, address)
	if err != nil {
		ctx.SetError(errors.Wrapf(err, "failed to get connection to %s", address))
		return filter.Stop
	}

	// Prepare metadata
	md := metadata.MD{}
	if ctx.Attachments != nil {
		for k, v := range ctx.Attachments {
			switch val := v.(type) {
			case string:
				md.Set(k, val)
			case []string:
				md.Set(k, val...)
			default:
				// Skip non-string metadata
				continue
			}
		}
	}

	// Create outgoing context with metadata and timeout
	outCtx := metadata.NewOutgoingContext(callCtx, md)

	logger.Debugf("gRPC proxy unary call to %s.%s at %s", ctx.ServiceName, ctx.MethodName, address)

	// Get request from arguments
	var req interface{}
	if len(ctx.Arguments) > 0 {
		req = ctx.Arguments[0]
	} else {
		logger.Warn("No arguments provided for unary call")
	}

	// Create response holder
	var resp interface{}

	// Execute the gRPC call with timeout
	err = conn.Invoke(outCtx, fullMethod, req, &resp)
	if err != nil {
		ctx.SetError(errors.Wrapf(err, "gRPC call to %s failed", fullMethod))
		logger.Errorf("gRPC unary call to %s failed: %v", fullMethod, err)
		return filter.Stop
	}

	// Set result
	ctx.SetResult(resp)

	// Log success with timing
	duration := time.Since(start)
	logger.Debugf("gRPC unary call to %s completed successfully in %v", fullMethod, duration)

	return filter.Continue
}

// handleClientStream handle client streaming call
func (f *Filter) handleClientStream(ctx *grpcCtx.GrpcContext, address string) filter.FilterStatus {
	// Get or create connection
	conn, err := f.getOrCreateConnection(ctx, address)
	if err != nil {
		ctx.SetError(errors.Errorf("gRPC proxy failed to get connection: %v", err))
		return filter.Stop
	}

	// Set metadata
	md := make(metadata.MD)
	for k, v := range ctx.Attachments {
		if str, ok := v.(string); ok {
			md.Set(k, str)
		}
	}
	outCtx := metadata.NewOutgoingContext(ctx.Context, md)

	// Create full method path
	fullMethod := "/" + ctx.ServiceName + "/" + ctx.MethodName
	logger.Debugf("[dubbo-go-pixiu] gRPC proxy client stream to %s", fullMethod)

	// Create client stream
	clientStream, err := conn.NewStream(outCtx, &grpc.StreamDesc{
		StreamName:    ctx.MethodName,
		ServerStreams: false,
		ClientStreams: true,
	}, fullMethod)

	if err != nil {
		ctx.SetError(errors.Errorf("failed to create client stream: %v", err))
		return filter.Stop
	}

	// Check if we have a stream to work with
	if ctx.Stream == nil {
		ctx.SetError(errors.New("no stream available in context"))
		return filter.Stop
	}

	// Process client streaming - forward messages from client to backend
	var wg sync.WaitGroup
	wg.Add(1)

	// Channel to signal completion or errors
	errChan := make(chan error, 2)

	// Forward messages from client to backend
	go func() {
		defer wg.Done()
		defer clientStream.CloseSend()

		for {
			var msg interface{}
			err := ctx.Stream.RecvMsg(&msg)
			if err == io.EOF {
				logger.Debug("Client stream closed")
				break
			}
			if err != nil {
				errChan <- errors.Wrap(err, "error receiving from client")
				return
			}

			if err := clientStream.SendMsg(msg); err != nil {
				errChan <- errors.Wrap(err, "error forwarding to backend")
				return
			}
		}
	}()

	// Wait for client streaming to complete
	wg.Wait()

	// Check for errors during streaming
	select {
	case err := <-errChan:
		ctx.SetError(err)
		return filter.Stop
	default:
		// No errors
	}

	// Receive final response from backend
	var resp interface{}
	err = clientStream.RecvMsg(&resp)
	if err != nil && err != io.EOF {
		ctx.SetError(errors.Errorf("failed to receive response: %v", err))
		return filter.Stop
	}

	// Set result and forward to client
	ctx.SetResult(resp)
	if resp != nil {
		if err := ctx.Stream.SendMsg(resp); err != nil && err != io.EOF {
			logger.Errorf("Failed to send response back to client: %v", err)
		}
	}

	logger.Debugf("gRPC client stream to %s completed successfully", fullMethod)
	return filter.Continue
}

// handleServerStream handle server streaming call
func (f *Filter) handleServerStream(ctx *grpcCtx.GrpcContext, address string) filter.FilterStatus {
	// Get or create connection
	conn, err := f.getOrCreateConnection(ctx, address)
	if err != nil {
		ctx.SetError(errors.Errorf("gRPC proxy failed to get connection: %v", err))
		return filter.Stop
	}

	// Set metadata
	md := make(metadata.MD)
	for k, v := range ctx.Attachments {
		if str, ok := v.(string); ok {
			md.Set(k, str)
		}
	}
	outCtx := metadata.NewOutgoingContext(ctx.Context, md)

	// Create full method path
	fullMethod := "/" + ctx.ServiceName + "/" + ctx.MethodName
	logger.Debugf("[dubbo-go-pixiu] gRPC proxy server stream to %s", fullMethod)

	// Create client stream
	clientStream, err := conn.NewStream(outCtx, &grpc.StreamDesc{
		StreamName:    ctx.MethodName,
		ServerStreams: true,
		ClientStreams: false,
	}, fullMethod)

	if err != nil {
		ctx.SetError(errors.Errorf("failed to create client stream: %v", err))
		return filter.Stop
	}

	// Check if we have a stream to work with
	if ctx.Stream == nil {
		ctx.SetError(errors.New("no stream available in context"))
		return filter.Stop
	}

	// Send initial request to backend
	if len(ctx.Arguments) > 0 {
		// Send first argument as request body
		err = clientStream.SendMsg(ctx.Arguments[0])
		if err != nil {
			ctx.SetError(errors.Errorf("failed to send request: %v", err))
			return filter.Stop
		}
	} else {
		// If no arguments provided, try to receive initial message from client
		var initialMsg interface{}
		if err := ctx.Stream.RecvMsg(&initialMsg); err == nil {
			if err := clientStream.SendMsg(initialMsg); err != nil {
				ctx.SetError(errors.Errorf("failed to forward initial message: %v", err))
				return filter.Stop
			}
		}
	}

	// Close send side
	clientStream.CloseSend()

	// Forward server stream responses back to client
	var responses []interface{}
	for {
		var resp interface{}
		err := clientStream.RecvMsg(&resp)
		if err == io.EOF {
			break
		}
		if err != nil {
			ctx.SetError(errors.Errorf("failed to receive stream response: %v", err))
			return filter.Stop
		}

		// Store response for result
		responses = append(responses, resp)

		// Forward to client
		if err := ctx.Stream.SendMsg(resp); err != nil {
			logger.Errorf("Failed to forward response to client: %v", err)
			// Continue trying with other responses
		}
	}

	logger.Debugf("gRPC server stream to %s completed, received %d responses",
		fullMethod, len(responses))

	// Store all responses as result
	ctx.SetResult(responses)
	return filter.Continue
}

// handleBidirectionalStream handle bidirectional streaming call
func (f *Filter) handleBidirectionalStream(ctx *grpcCtx.GrpcContext, address string) filter.FilterStatus {
	// Get or create connection
	conn, err := f.getOrCreateConnection(ctx, address)
	if err != nil {
		ctx.SetError(errors.Errorf("gRPC proxy failed to get connection: %v", err))
		return filter.Stop
	}

	// Set metadata
	md := make(metadata.MD)
	for k, v := range ctx.Attachments {
		if str, ok := v.(string); ok {
			md.Set(k, str)
		}
	}
	outCtx := metadata.NewOutgoingContext(ctx.Context, md)

	// Create full method path
	fullMethod := "/" + ctx.ServiceName + "/" + ctx.MethodName
	logger.Debugf("[dubbo-go-pixiu] gRPC proxy bidirectional stream to %s", fullMethod)

	// Create client stream
	clientStream, err := conn.NewStream(outCtx, &grpc.StreamDesc{
		StreamName:    ctx.MethodName,
		ServerStreams: true,
		ClientStreams: true,
	}, fullMethod)

	if err != nil {
		ctx.SetError(errors.Errorf("failed to create client stream: %v", err))
		return filter.Stop
	}

	// Check if we have a stream to work with
	if ctx.Stream == nil {
		ctx.SetError(errors.New("no stream available in context"))
		return filter.Stop
	}

	// Process bidirectional streaming
	var wg sync.WaitGroup
	wg.Add(2) // One for client->server, one for server->client

	// Channel to signal completion or errors
	errChan := make(chan error, 2)
	doneChan := make(chan struct{})

	// Forward messages from client to backend
	go func() {
		defer wg.Done()

		// Send initial arguments if available
		if len(ctx.Arguments) > 0 {
			for _, arg := range ctx.Arguments {
				if err := clientStream.SendMsg(arg); err != nil {
					errChan <- errors.Wrap(err, "failed to send initial message")
					return
				}
			}
		}

		// Continue forwarding messages from client
		for {
			select {
			case <-doneChan:
				return
			default:
				var msg interface{}
				err := ctx.Stream.RecvMsg(&msg)
				if err == io.EOF {
					clientStream.CloseSend()
					return
				}
				if err != nil {
					errChan <- errors.Wrap(err, "error receiving from client")
					return
				}

				if err := clientStream.SendMsg(msg); err != nil {
					errChan <- errors.Wrap(err, "error forwarding to backend")
					return
				}
			}
		}
	}()

	// Forward responses from backend to client
	go func() {
		defer wg.Done()

		var responses []interface{}
		for {
			select {
			case <-doneChan:
				return
			default:
				var resp interface{}
				err := clientStream.RecvMsg(&resp)
				if err == io.EOF {
					return
				}
				if err != nil {
					errChan <- errors.Wrap(err, "error receiving from backend")
					return
				}

				responses = append(responses, resp)

				if err := ctx.Stream.SendMsg(resp); err != nil {
					logger.Errorf("Failed to forward response to client: %v", err)
					// Continue with other responses
				}
			}
		}
	}()

	// Wait for either error or completion
	go func() {
		select {
		case <-ctx.Context.Done():
			close(doneChan)
		case err := <-errChan:
			ctx.SetError(err)
			close(doneChan)
		}
	}()

	// Wait for both goroutines to complete
	wg.Wait()

	// Check for errors
	select {
	case err := <-errChan:
		ctx.SetError(err)
		return filter.Stop
	default:
		// No errors
	}

	logger.Debugf("gRPC bidirectional stream to %s completed successfully", fullMethod)
	return filter.Continue
}

// getOrCreateConnection retrieves an existing connection or creates a new one
func (f *Filter) getOrCreateConnection(ctx *grpcCtx.GrpcContext, address string) (*grpc.ClientConn, error) {
	if address == "" {
		return nil, errors.New("cannot create connection to empty address")
	}

	// Use address as cache key for connection reuse
	cacheKey := address

	// First try to get existing connection with read lock
	f.mu.RLock()
	if conn, ok := f.clientConnPool.Load(cacheKey); ok {
		f.mu.RUnlock()

		// Check if connection is still valid
		if grpcConn, ok := conn.(*grpc.ClientConn); ok {
			state := grpcConn.GetState()
			if state != connectivity.Shutdown && state != connectivity.TransientFailure {
				logger.Debugf("Reusing existing connection to %s (state: %s)", address, state.String())
				return grpcConn, nil
			}
			logger.Warnf("Found stale connection to %s in state %s, will create new one",
				address, state.String())
		}
	} else {
		f.mu.RUnlock()
	}

	// Need write lock to create new connection
	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check if another goroutine created the connection while we were waiting
	if conn, ok := f.clientConnPool.Load(cacheKey); ok {
		if grpcConn, ok := conn.(*grpc.ClientConn); ok {
			state := grpcConn.GetState()
			if state != connectivity.Shutdown && state != connectivity.TransientFailure {
				logger.Debugf("Another goroutine created connection to %s, reusing it", address)
				return grpcConn, nil
			}
			// Connection exists but is in bad state, remove it
			logger.Debugf("Removing stale connection to %s in state %s", address, state.String())
			f.clientConnPool.Delete(cacheKey)
		}
	}

	// Get cluster manager to find the best endpoint
	clusterManager := server.GetClusterManager()
	if clusterManager == nil || ctx.Route == nil || ctx.Route.Cluster == "" {
		// Cannot use cluster manager, use provided address directly
		logger.Debug("Using direct connection without cluster manager")
		return f.createAndStoreConnection(address, cacheKey)
	}

	// Try to get better endpoint from cluster
	clusterName := ctx.Route.Cluster
	endpoint := clusterManager.PickEndpoint(clusterName, ctx)
	if endpoint == nil {
		logger.Warnf("No endpoint found in cluster %s, using provided address", clusterName)
		return f.createAndStoreConnection(address, cacheKey)
	}

	// Get endpoint address
	endpointAddress := endpoint.Address.GetAddress()
	if endpointAddress == "" {
		logger.Warn("Empty endpoint address from cluster, using provided address")
		return f.createAndStoreConnection(address, cacheKey)
	}

	// Use endpoint address as cache key
	endpointCacheKey := endpointAddress

	// Create connection to endpoint
	logger.Infof("Creating new backend connection to %s for %s.%s",
		endpointAddress, ctx.ServiceName, ctx.MethodName)

	conn, err := f.createConnection(endpointAddress)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to %s", endpointAddress)
	}

	// Store in connection pool
	f.clientConnPool.Store(endpointCacheKey, conn)

	// Start connection health checker
	go f.monitorConnection(endpointCacheKey, conn)

	return conn, nil
}

// createAndStoreConnection creates a new connection and stores it in the pool
func (f *Filter) createAndStoreConnection(address, cacheKey string) (*grpc.ClientConn, error) {
	conn, err := f.createConnection(address)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create connection to %s", address)
	}

	// Store in connection pool
	f.clientConnPool.Store(cacheKey, conn)

	// Start connection health checker
	go f.monitorConnection(cacheKey, conn)

	return conn, nil
}

// monitorConnection periodically checks connection health and removes bad connections
func (f *Filter) monitorConnection(cacheKey string, conn *grpc.ClientConn) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			state := conn.GetState()
			if state == connectivity.Shutdown || state == connectivity.TransientFailure {
				logger.Warnf("Connection to %s is in bad state (%s), removing from pool",
					cacheKey, state.String())

				f.mu.Lock()
				if currentConn, ok := f.clientConnPool.Load(cacheKey); ok {
					if currentConn == conn {
						f.clientConnPool.Delete(cacheKey)
					}
				}
				f.mu.Unlock()

				return
			}
		}
	}
}

// createConnection creates a new gRPC connection with optimized settings
func (f *Filter) createConnection(address string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	// Configure TLS
	if f.Config.EnableTLS {
		if f.Config.TLSCertFile != "" && f.Config.TLSKeyFile != "" {
			creds, err := f.createTLSCredentials()
			if err != nil {
				logger.Warnf("Failed to load TLS credentials: %v, falling back to insecure", err)
				opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
			} else {
				opts = append(opts, grpc.WithTransportCredentials(creds))
			}
		} else {
			logger.Warn("TLS enabled but certificate files not provided, falling back to insecure")
			opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Configure keepalive parameters
	keepAliveTime := defaultKeepAliveTime
	keepAliveTimeout := defaultKeepAliveTimeout

	if f.Config.KeepAliveTime != "" {
		if duration, err := time.ParseDuration(f.Config.KeepAliveTime); err == nil {
			keepAliveTime = duration
		}
	}

	if f.Config.KeepAliveTimeout != "" {
		if duration, err := time.ParseDuration(f.Config.KeepAliveTimeout); err == nil {
			keepAliveTimeout = duration
		}
	}

	// Add keepalive options
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                keepAliveTime,
		Timeout:             keepAliveTimeout,
		PermitWithoutStream: true, // Allow pings even without active streams
	}))

	// Configure connection timeout
	connectTimeout := defaultConnectTimeout
	if f.Config.ConnectTimeout != "" {
		if duration, err := time.ParseDuration(f.Config.ConnectTimeout); err == nil {
			connectTimeout = duration
		}
	}

	opts = append(opts, grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: connectTimeout,
	}))

	// Set max concurrent streams if configured
	if f.Config.MaxConcurrentStreams > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(1024*1024*4), // 4MB
			grpc.MaxCallSendMsgSize(1024*1024*4), // 4MB
		))
	}

	// Add retry policy
	opts = append(opts, grpc.WithDisableRetry())

	// Establish connection with timeout
	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout)
	defer cancel()

	logger.Debugf("Establishing new gRPC connection to %s", address)
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to backend %s", address)
	}

	return conn, nil
}

// createTLSCredentials creates TLS credentials from certificate files
func (f *Filter) createTLSCredentials() (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(f.Config.TLSCertFile, f.Config.TLSKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load TLS key pair")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	return credentials.NewTLS(tlsConfig), nil
}

// Close gracefully closes all connections and cleans up resources
func (f *Filter) Close() error {
	logger.Info("Closing gRPC proxy filter and all connections")

	// Use a wait group to track all connection closures
	var wg sync.WaitGroup
	var closeErrors []error
	var errorMu sync.Mutex

	// Close all connections with timeout
	f.clientConnPool.Range(func(key, value any) bool {
		if conn, ok := value.(*grpc.ClientConn); ok {
			wg.Add(1)
			go func(address string, conn *grpc.ClientConn) {
				defer wg.Done()

				// Create context with timeout for graceful close
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				logger.Debugf("Closing gRPC connection to %s", address)
				if err := conn.Close(); err != nil {
					errorMu.Lock()
					closeErrors = append(closeErrors, errors.Wrapf(err, "error closing connection to %s", address))
					errorMu.Unlock()
				}

				// Wait for connection to actually close or timeout
				<-ctx.Done()
			}(key.(string), conn)
		}
		return true
	})

	// Wait for all connections to close
	wg.Wait()

	// Clear connection pool
	f.clientConnPool = sync.Map{}

	// Report any errors
	if len(closeErrors) > 0 {
		logger.Warnf("Encountered %d errors while closing gRPC connections", len(closeErrors))
		return closeErrors[0]
	}

	logger.Info("Successfully closed all gRPC proxy connections")
	return nil
}
