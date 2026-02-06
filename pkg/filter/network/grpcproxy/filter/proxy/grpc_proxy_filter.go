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
)

import (
	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/backoff"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/metadata"

	"google.golang.org/protobuf/reflect/protoreflect"
)

import (
	ptcodec "github.com/apache/dubbo-go-pixiu/pkg/common/codec/grpc/passthrough"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	grpcCtx "github.com/apache/dubbo-go-pixiu/pkg/context/grpc"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

// Constants for gRPC proxy filter
const (
	Kind                       = constant.GRPCProxyFilter
	defaultKeepAliveTime       = 300 * time.Second
	defaultKeepAliveTimeout    = 5 * time.Second
	defaultConnectTimeout      = 5 * time.Second
	defaultMaxMsgSize          = 4 * 1024 * 1024 // 4MB
	defaultHealthCheckInterval = 30 * time.Second
	defaultDescriptorCacheTTL  = 5 * time.Minute
)

// ReflectionMode defines the mode of gRPC reflection
const (
	// ReflectionModePassthrough performs transparent binary proxying (default, highest performance)
	ReflectionModePassthrough = "passthrough"
	// ReflectionModeReflection uses gRPC reflection to decode/encode messages (content-aware)
	ReflectionModeReflection = "reflection"
	// ReflectionModeHybrid tries reflection first, falls back to passthrough on failure
	ReflectionModeHybrid = "hybrid"
)

func init() {
	filter.RegisterGrpcFilterPlugin(&Plugin{})
}

type (
	// Plugin gRPC proxy plugin implementation
	Plugin struct{}

	// Config defines the configuration options for the gRPC proxy filter
	Config struct {
		EnableTLS            bool          `yaml:"enable_tls" json:"enable_tls" mapstructure:"enable_tls"`
		TLSCertFile          string        `yaml:"tls_cert_file" json:"tls_cert_file" mapstructure:"tls_cert_file"`
		TLSKeyFile           string        `yaml:"tls_key_file" json:"tls_key_file" mapstructure:"tls_key_file"`
		MaxConcurrentStreams uint32        `yaml:"max_concurrent_streams" json:"max_concurrent_streams" mapstructure:"max_concurrent_streams"`
		KeepAliveTimeStr     string        `yaml:"keepalive_time" json:"keepalive_time" mapstructure:"keepalive_time"`
		KeepAliveTimeoutStr  string        `yaml:"keepalive_timeout" json:"keepalive_timeout" mapstructure:"keepalive_timeout"`
		ConnectTimeoutStr    string        `yaml:"connect_timeout" json:"connect_timeout" mapstructure:"connect_timeout"`
		KeepAliveTime        time.Duration `yaml:"-" json:"-"`
		KeepAliveTimeout     time.Duration `yaml:"-" json:"-"`
		ConnectTimeout       time.Duration `yaml:"-" json:"-"`

		// ReflectionMode: "passthrough" (default), "reflection", or "hybrid"
		ReflectionMode        string        `yaml:"reflection_mode" json:"reflection_mode" mapstructure:"reflection_mode"`
		DescriptorCacheTTLStr string        `yaml:"descriptor_cache_ttl" json:"descriptor_cache_ttl" mapstructure:"descriptor_cache_ttl"`
		DescriptorCacheTTL    time.Duration `yaml:"-" json:"-"`
		ExtractTripleMetadata bool          `yaml:"extract_triple_metadata" json:"extract_triple_metadata" mapstructure:"extract_triple_metadata"`
	}

	// Filter implements the gRPC proxy filter
	Filter struct {
		Config            *Config
		clientConnPool    sync.Map     // address -> *grpc.ClientConn
		mu                sync.RWMutex // protects concurrent operations
		reflectionManager *ReflectionManager
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

	// Parse time durations from strings, with defaults
	cfg.KeepAliveTime = parseDurationWithDefault(cfg.KeepAliveTimeStr, defaultKeepAliveTime)
	cfg.KeepAliveTimeout = parseDurationWithDefault(cfg.KeepAliveTimeoutStr, defaultKeepAliveTimeout)
	cfg.ConnectTimeout = parseDurationWithDefault(cfg.ConnectTimeoutStr, defaultConnectTimeout)
	cfg.DescriptorCacheTTL = parseDurationWithDefault(cfg.DescriptorCacheTTLStr, defaultDescriptorCacheTTL)

	// Set default reflection mode if not specified
	if cfg.ReflectionMode == "" {
		cfg.ReflectionMode = ReflectionModePassthrough
	}

	f := &Filter{Config: cfg}

	// Initialize reflection manager if reflection mode is enabled
	if cfg.ReflectionMode == ReflectionModeReflection || cfg.ReflectionMode == ReflectionModeHybrid {
		f.reflectionManager = NewReflectionManager(cfg.DescriptorCacheTTL)
		logger.Infof("gRPC proxy filter initialized with reflection mode: %s, cache TTL: %s",
			cfg.ReflectionMode, cfg.DescriptorCacheTTL)
	} else {
		logger.Infof("gRPC proxy filter initialized with passthrough mode")
	}

	return f, nil
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

	logger.Debugf("Forwarding gRPC request %s to cluster %s, endpoint %s",
		ctx.ServiceName+"/"+ctx.MethodName, ctx.Route.Cluster, address)

	return f.handleStream(ctx, address)
}

// handleStream handles all types of gRPC calls by creating a full-duplex stream pipe.
func (f *Filter) handleStream(ctx *grpcCtx.GrpcContext, address string) filter.FilterStatus {
	conn, err := f.getOrCreateConnection(address)
	if err != nil {
		ctx.SetError(errors.Errorf("gRPC proxy failed to get connection: %v", err))
		return filter.Stop
	}

	md := make(metadata.MD)
	for k, v := range ctx.Attachments {
		if str, ok := v.(string); ok {
			md.Set(k, str)
		}
	}

	if f.Config.ExtractTripleMetadata {
		tripleMeta := ExtractTripleMetadata(ctx.Attachments)
		if len(tripleMeta) > 0 {
			ctx.SetAttachment("_triple_metadata", tripleMeta)
			logger.Debugf("Extracted Triple metadata: %v", tripleMeta)
		}
	}

	outCtx := metadata.NewOutgoingContext(ctx.Context, md)
	fullMethod := ctx.ServiceName + "/" + ctx.MethodName

	var codecOpt grpc.CallOption
	var useReflection bool

	if f.reflectionManager != nil && f.Config.ReflectionMode != ReflectionModePassthrough {
		methodDesc, err := f.reflectionManager.GetMethodDescriptor(
			ctx.Context, conn, address, ctx.ServiceName, ctx.MethodName)
		if err != nil {
			if f.Config.ReflectionMode == ReflectionModeHybrid {
				logger.Warnf("Reflection failed for %s, falling back to passthrough: %v", fullMethod, err)
				codecOpt = grpc.ForceCodec(ptcodec.Codec{})
				useReflection = false
			} else {
				ctx.SetError(errors.Wrapf(err, "reflection failed for %s", fullMethod))
				return filter.Stop
			}
		} else {
			ctx.SetAttachment("_method_descriptor", methodDesc)
			ctx.SetAttachment("_is_client_streaming", methodDesc.IsStreamingClient())
			logger.Debugf("Reflection succeeded for %s, input: %s, output: %s, clientStream: %v",
				fullMethod, methodDesc.Input().FullName(), methodDesc.Output().FullName(), methodDesc.IsStreamingClient())
			codecOpt = grpc.ForceCodec(NewDynamicCodec(methodDesc))
			useReflection = true
		}
	} else {
		codecOpt = grpc.ForceCodec(ptcodec.Codec{})
		useReflection = false
	}

	ctx.SetAttachment("_use_reflection", useReflection)

	clientStream, err := conn.NewStream(outCtx, &grpc.StreamDesc{
		StreamName:    ctx.MethodName,
		ServerStreams: true,
		ClientStreams: true,
	}, fullMethod, codecOpt)
	if err != nil {
		ctx.SetError(errors.Errorf("failed to create client stream: %v", err))
		return filter.Stop
	}

	if ctx.Stream == nil {
		ctx.SetError(errors.New("no stream available in context"))
		return filter.Stop
	}

	var wg sync.WaitGroup
	wg.Add(2)
	errChan := make(chan error, 2)
	doneChan := make(chan struct{})

	go f.forwardClientToServer(ctx, clientStream, &wg, errChan, doneChan)
	go f.forwardServerToClient(ctx, clientStream, &wg, errChan, doneChan)

	go func() {
		select {
		case <-ctx.Context.Done():
			close(doneChan)
		case err := <-errChan:
			ctx.SetError(err)
			close(doneChan)
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil && ctx.Error == nil {
			ctx.SetError(err)
		}
	}

	if ctx.Error != nil {
		logger.Debugf("gRPC stream for %s completed with error: %v", fullMethod, ctx.Error)
		return filter.Stop
	}

	// The listener already logs the successful completion with duration.
	// logger.Debugf("gRPC stream for %s completed successfully", fullMethod)
	return filter.Continue
}

// forwardClientToServer forwards messages from the incoming client stream to the backend server stream.
func (f *Filter) forwardClientToServer(ctx *grpcCtx.GrpcContext, clientStream grpc.ClientStream, wg *sync.WaitGroup, errChan chan<- error, doneChan <-chan struct{}) {
	defer wg.Done()

	useReflection := false
	if val, ok := ctx.GetAttachment("_use_reflection"); ok {
		useReflection, _ = val.(bool)
	}

	var methodDesc protoreflect.MethodDescriptor
	if useReflection {
		if md, ok := ctx.GetAttachment("_method_descriptor"); ok {
			methodDesc, _ = md.(protoreflect.MethodDescriptor)
		}
	}

	// Check if this is a client-streaming call (client-stream or bidirectional)
	isClientStreaming := false
	if val, ok := ctx.GetAttachment("_is_client_streaming"); ok {
		isClientStreaming, _ = val.(bool)
	}

	if len(ctx.Arguments) > 0 {
		for _, arg := range ctx.Arguments {
			if err := clientStream.SendMsg(arg); err != nil {
				errChan <- errors.Wrap(err, "failed to send initial message")
				return
			}
		}
	}

	// For unary/server-stream: read one message then CloseSend to prevent deadlock
	if !isClientStreaming {
		var msg []byte
		if err := ctx.Stream.RecvMsg(&msg); err != nil {
			if err == io.EOF {
				if err := clientStream.CloseSend(); err != nil {
					logger.Errorf("Error closing send stream to backend: %v", err)
				}
				return
			}
			errChan <- errors.Wrap(err, "error receiving from client")
			return
		}

		if methodDesc != nil {
			if decoded, err := DecodeRequest(methodDesc, msg); err == nil {
				logger.Debugf("Decoded request message: %v", decoded.Message)
			}
		}

		if err := clientStream.SendMsg(msg); err != nil {
			errChan <- errors.Wrap(err, "error forwarding to backend")
			return
		}

		if err := clientStream.CloseSend(); err != nil {
			logger.Errorf("Error closing send stream to backend: %v", err)
		}
		return
	}

	// Client-streaming: continuously forward messages until EOF
	for {
		select {
		case <-doneChan:
			return
		default:
			var msg []byte
			if err := ctx.Stream.RecvMsg(&msg); err != nil {
				if err == io.EOF {
					if err := clientStream.CloseSend(); err != nil {
						logger.Errorf("Error closing send stream to backend: %v", err)
					}
					return
				}
				errChan <- errors.Wrap(err, "error receiving from client")
				return
			}

			if methodDesc != nil {
				if decoded, err := DecodeRequest(methodDesc, msg); err == nil {
					logger.Debugf("Decoded request message: %v", decoded.Message)
				}
			}

			if err := clientStream.SendMsg(msg); err != nil {
				errChan <- errors.Wrap(err, "error forwarding to backend")
				return
			}
		}
	}
}

// forwardServerToClient forwards messages from the backend server stream to the incoming client stream.
func (f *Filter) forwardServerToClient(ctx *grpcCtx.GrpcContext, clientStream grpc.ClientStream, wg *sync.WaitGroup, errChan chan<- error, doneChan <-chan struct{}) {
	defer wg.Done()

	var methodDesc protoreflect.MethodDescriptor
	if val, ok := ctx.GetAttachment("_use_reflection"); ok && val.(bool) {
		if md, ok := ctx.GetAttachment("_method_descriptor"); ok {
			methodDesc, _ = md.(protoreflect.MethodDescriptor)
		}
	}

	if header, err := clientStream.Header(); err == nil {
		if s, ok := ctx.Stream.(grpc.ServerStream); ok {
			s.SetHeader(header)
		}
	}

	for {
		select {
		case <-doneChan:
			return
		default:
			var resp []byte
			err := clientStream.RecvMsg(&resp)
			if err != nil {
				if s, ok := ctx.Stream.(grpc.ServerStream); ok {
					s.SetTrailer(clientStream.Trailer())
				}
				if err != io.EOF {
					errChan <- err
				}
				return
			}

			if methodDesc != nil {
				if decoded, err := DecodeResponse(methodDesc, resp); err == nil {
					logger.Debugf("Decoded response message: %v", decoded.Message)
				}
			}

			if err := ctx.Stream.SendMsg(resp); err != nil {
				errChan <- errors.Wrap(err, "failed to forward response to client")
				return
			}
		}
	}
}

// parseDurationWithDefault parses a string duration and returns a default if empty or invalid.
func parseDurationWithDefault(durationStr string, defaultVal time.Duration) time.Duration {
	if durationStr == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Warnf("Invalid duration format: '%s', using default %s", durationStr, defaultVal)
		return defaultVal
	}
	return d
}

// getOrCreateConnection retrieves an existing connection or creates a new one for a given address.
func (f *Filter) getOrCreateConnection(address string) (*grpc.ClientConn, error) {
	if address == "" {
		return nil, errors.New("cannot create connection to empty address")
	}

	// Optimistic check without a lock. `sync.Map` is safe for concurrent reads.
	if conn, ok := f.clientConnPool.Load(address); ok {
		if grpcConn, ok := conn.(*grpc.ClientConn); ok {
			state := grpcConn.GetState()
			if state != connectivity.Shutdown && state != connectivity.TransientFailure {
				logger.Debugf("Reusing existing connection to %s (state: %s)", address, state.String())
				return grpcConn, nil
			}
			// If the connection is stale, it will be handled by the write-lock path.
			logger.Warnf("Found stale connection to %s in state %s, will create new one", address, state.String())
		}
	}

	// If no valid connection is found, acquire a write lock to create one.
	f.mu.Lock()
	defer f.mu.Unlock()

	// Double-check if another goroutine created the connection while we were waiting for the lock.
	if conn, ok := f.clientConnPool.Load(address); ok {
		if grpcConn, ok := conn.(*grpc.ClientConn); ok {
			state := grpcConn.GetState()
			if state != connectivity.Shutdown && state != connectivity.TransientFailure {
				logger.Debugf("Another goroutine created connection to %s, reusing it", address)
				return grpcConn, nil
			}
			// The existing connection is bad, remove it before creating a new one.
			f.clientConnPool.Delete(address)
		}
	}

	// Create a new connection.
	logger.Infof("Creating new backend connection to %s", address)
	conn, err := f.createConnection(address)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to %s", address)
	}

	// Store the new connection in the pool.
	f.clientConnPool.Store(address, conn)

	// Start a goroutine to monitor the connection's health.
	go f.monitorConnection(address, conn)

	return conn, nil
}

// monitorConnection periodically checks connection health and removes bad connections
func (f *Filter) monitorConnection(cacheKey string, conn *grpc.ClientConn) {
	ticker := time.NewTicker(defaultHealthCheckInterval)
	defer ticker.Stop()

	for range ticker.C {
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

	// Add keepalive options
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:                f.Config.KeepAliveTime,
		Timeout:             f.Config.KeepAliveTimeout,
		PermitWithoutStream: false,
	}))

	// Configure connection timeout
	opts = append(opts, grpc.WithConnectParams(grpc.ConnectParams{
		Backoff:           backoff.DefaultConfig,
		MinConnectTimeout: f.Config.ConnectTimeout,
	}))

	// Set max concurrent streams if configured
	if f.Config.MaxConcurrentStreams > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(defaultMaxMsgSize),
			grpc.MaxCallSendMsgSize(defaultMaxMsgSize),
		))
	}

	// Disable proxy-level retries. Let the client handle retries.
	opts = append(opts, grpc.WithDisableRetry())

	// WithBlock makes Dial block until the connection is established.
	opts = append(opts, grpc.WithBlock())

	// Dial the backend
	// Use a background context as the connection lifecycle is managed by the filter, not a single request
	dialCtx, cancel := context.WithTimeout(context.Background(), f.Config.ConnectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(dialCtx, address, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to dial backend %s", address)
	}

	return conn, nil
}

// createTLSCredentials loads TLS certs for mTLS
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

	// Close reflection manager if initialized
	if f.reflectionManager != nil {
		f.reflectionManager.Close()
		logger.Info("Reflection manager closed")
	}

	var wg sync.WaitGroup
	var closeErrors []error
	var errorMu sync.Mutex

	f.clientConnPool.Range(func(key, value any) bool {
		if conn, ok := value.(*grpc.ClientConn); ok {
			wg.Add(1)
			go func() {
				defer wg.Done()
				logger.Debugf("Closing gRPC connection to %s", conn.Target())
				if err := conn.Close(); err != nil {
					logger.Warnf("Error closing gRPC connection to %s: %v", conn.Target(), err)
					errorMu.Lock()
					closeErrors = append(closeErrors, err)
					errorMu.Unlock()
				}
			}()
		}
		return true
	})

	wg.Wait()
	logger.Info("All gRPC proxy connections have been closed")

	if len(closeErrors) > 0 {
		return closeErrors[0]
	}

	return nil
}
