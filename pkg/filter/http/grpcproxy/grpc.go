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

package grpcproxy

import (
	"context"
	"errors"
	"fmt"
	"io"
	stdHttp "net/http"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/golang/protobuf/jsonpb" //nolint
	"github.com/golang/protobuf/proto"  //nolint

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"

	perrors "github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	ct "github.com/apache/dubbo-go-pixiu/pkg/context"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPGrpcProxyFilter

	loggerHeader = "[grpc-proxy]"

	// DescriptorSourceKey current ds
	DescriptorSourceKey = "DescriptorSource"

	// GrpcClientConnKey the grpc-client-conn by the coroutine local
	GrpcClientConnKey = "GrpcClientConn"
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

const (
	NONE   = "none"
	AUTO   = "auto"
	LOCAL  = "local"
	REMOTE = "remote"
)

type (
	DescriptorSourceStrategy string

	// Plugin is grpc filter plugin.
	Plugin struct {
	}

	// FilterFactory is grpc filter instance
	FilterFactory struct {
		cfg *Config
		// grpc descriptor source factory
		descriptor *Descriptor
		// hold grpc.ClientConns, key format: cluster name + NUL + endpoint
		connections           *grpcConnectionManager
		removeEndpointHandler func()

		extReg      *dynamic.ExtensionRegistry
		registered  map[string]bool
		extensionMu *sync.RWMutex
	}
	Filter struct {
		cfg *Config
		// grpc descriptor source factory
		descriptor *Descriptor
		// hold grpc.ClientConns, key format: cluster name + NUL + endpoint
		connections *grpcConnectionManager

		extReg      *dynamic.ExtensionRegistry
		registered  map[string]bool
		extensionMu *sync.RWMutex
	}

	// Config describe the config of AccessFilter
	Config struct {
		DescriptorSourceStrategy DescriptorSourceStrategy `yaml:"descriptor_source_strategy" json:"descriptor_source_strategy" default:"auto"`
		Path                     string                   `yaml:"path" json:"path"`
		Rules                    []*Rule                  `yaml:"rules" json:"rules"`     //nolint
		Timeout                  time.Duration            `yaml:"timeout" json:"timeout"` //nolint
	}

	Rule struct {
		Selector string `yaml:"selector" json:"selector"`
		Match    Match  `yaml:"match" json:"match"`
	}

	Match struct {
		Method string `yaml:"method" json:"method"` //nolint
	}
)

func (c DescriptorSourceStrategy) String() string {
	return string(c)
}

func (c DescriptorSourceStrategy) Val() DescriptorSourceStrategy {
	switch c {
	case NONE:
		return NONE
	case AUTO:
		return AUTO
	case LOCAL:
		return LOCAL
	case REMOTE:
		return REMOTE
	}
	return ""
}

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	descriptor := &Descriptor{}
	connections := newGRPCConnectionManager()
	connections.onRemove = descriptor.removeConnection
	var removeEndpointHandler func()
	if clusterManager := server.GetClusterManager(); clusterManager != nil {
		removeEndpointHandler = clusterManager.AddEndpointStateHandler(func(clusterName, endpoint string, present bool, version uint64) {
			connections.UpdateEndpointState(clusterName, endpoint, present, version)
		})
	}
	return &FilterFactory{
		cfg:                   &Config{DescriptorSourceStrategy: AUTO},
		descriptor:            descriptor,
		connections:           connections,
		removeEndpointHandler: removeEndpointHandler,
		extReg:                &dynamic.ExtensionRegistry{},
		registered:            make(map[string]bool),
		extensionMu:           &sync.RWMutex{},
	}, nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	// Deep copy config to avoid pointer sharing (factory.cfg may change at runtime)
	f := &Filter{
		cfg:         factory.cfg.DeepCopy(),
		descriptor:  factory.descriptor,
		connections: factory.connections,
		extReg:      factory.extReg,
		registered:  factory.registered,
		extensionMu: factory.extensionMu,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

// getServiceAndMethod first return value is package.service, second one is method name
func getServiceAndMethod(path string) (string, string) {
	pos := strings.LastIndex(path, "/")
	if pos < 0 {
		return "", ""
	}

	mth := path[pos+1:]
	prefix := strings.TrimSuffix(path, "/"+mth)

	pos = strings.LastIndex(prefix, "/")
	if pos < 0 {
		return "", ""
	}

	svc := prefix[pos+1:]
	return svc, mth
}

// Decode use the default http to grpc transcoding strategy https://cloud.google.com/endpoints/docs/grpc/transcoding
func (f *Filter) Decode(c *http.HttpContext) filter.FilterStatus {
	svc, mth := getServiceAndMethod(c.GetUrl())

	var clientConn *grpc.ClientConn
	var err error

	re := c.GetRouteEntry()
	logger.Debugf("%s client choose endpoint from cluster :%v", loggerHeader, re.Cluster)

	e := server.GetClusterManager().PickEndpoint(re.Cluster, c)
	if e == nil {
		logger.Errorf("%s err {cluster not exists}", loggerHeader)
		errResp := http.ServiceUnavailable.WithError(errors.New("cluster not exists"))
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}
	// timeout for Dial and Invoke
	ctx, cancel := context.WithTimeout(c.Ctx, c.Timeout)
	defer cancel()
	ep := e.Address.GetAddress()

	connectionKey := grpcConnectionKey(re.Cluster, ep)
	clientConn, err = f.connections.Get(ctx, connectionKey, ep)
	if err != nil || clientConn == nil {
		logger.Errorf("%s err {failed to connect to grpc service provider}: %v", loggerHeader, err)
		errResp := http.ServiceUnavailable.WithError(fmt.Errorf("endpoint not found: %w", err))
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// get DescriptorSource, contain file and reflection
	source, err := f.descriptor.getDescriptorSource(context.WithValue(ctx, ct.ContextKey(GrpcClientConnKey), clientConn), f.cfg)
	if err != nil {
		logger.Errorf("%s err %s : %s ", loggerHeader, "get desc source fail", err)
		errResp := http.ConfigurationError.WithError(fmt.Errorf("service not config proto file or the server not support reflection API"))
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}
	//put DescriptorSource concurrent, del if no need
	ctx = context.WithValue(ctx, ct.ContextKey(DescriptorSourceKey), source)

	mthDesc, err := f.descriptor.getMethodDescriptor(source, clientConn, svc, mth)
	if err != nil {
		if _, ok := err.(*serviceNotExposedError); ok {
			logger.Errorf("%s err {service not expose, %s}", loggerHeader, svc)
			errResp := http.BadRequest.WithError(err)
			c.SendLocalReply(errResp.Status, errResp.ToJSON())
			return filter.Stop
		}
		logger.Errorf("%s err {request path invalid, service: %s, method: %s, cause: %v}", loggerHeader, svc, mth, err)
		errResp := http.MethodNotAllowed.New()
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	err = f.registerExtension(source, mthDesc)
	if err != nil {
		logger.Errorf("%s err {%s}", loggerHeader, "register extension failed")
		errResp := http.ConfigurationError.WithError(fmt.Errorf("register extension failed: %w", err))
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	msgFac := dynamic.NewMessageFactoryWithExtensionRegistry(f.extReg)
	grpcReq := msgFac.NewMessage(mthDesc.GetInputType())

	err = jsonToProtoMsg(c.Request.Body, grpcReq)
	if err != nil && !errors.Is(err, io.EOF) {
		logger.Errorf("%s err {failed to convert json to proto msg, %s}", loggerHeader, err.Error())
		errResp := http.BadGateway.WithError(fmt.Errorf("protocol conversion error: %w", err))
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	stub := grpcdynamic.NewStubWithMessageFactory(clientConn, msgFac)

	// metadata in grpc has the same feature in http
	md := mapHeaderToMetadata(c.AllHeaders())
	ctx = metadata.NewOutgoingContext(ctx, md)

	md = metadata.MD{}
	t := metadata.MD{}

	resp, err := Invoke(ctx, stub, mthDesc, grpcReq, grpc.Header(&md), grpc.Trailer(&t))
	// judge err is server side error or not
	if st, ok := status.FromError(err); ok {
		// Handle client-side gRPC errors (e.g., InvalidArgument)
		if st.Code() != codes.OK && !isServerError(st) {
			logger.Errorf("%s err {gRPC client error, code: %s, msg: %s}", loggerHeader, st.Code(), st.Message())
			errResp := http.BadGateway.WithError(fmt.Errorf("gRPC client error: %w", err))
			c.SendLocalReply(errResp.Status, errResp.ToJSON())
			f.connections.Invalidate(connectionKey, clientConn)
			return filter.Stop
		}
		// Handle server-side gRPC errors
		if isServerError(st) {
			if isServerTimeout(st) {
				logger.Errorf("%s err {failed to invoke grpc service provider because timeout, err:%s}", loggerHeader, err.Error())
				errResp := http.GatewayTimeout.WithError(fmt.Errorf("upstream timeout: %w", err))
				c.SendLocalReply(errResp.Status, errResp.ToJSON())
				f.connections.Invalidate(connectionKey, clientConn)
				return filter.Stop
			}
			logger.Errorf("%s err {failed to invoke grpc service provider, %s}", loggerHeader, err.Error())
			errResp := http.ServiceUnavailable.WithError(fmt.Errorf("gRPC invoke error: %w", err))
			c.SendLocalReply(errResp.Status, errResp.ToJSON())
			f.connections.Invalidate(connectionKey, clientConn)
			return filter.Stop
		}
	} else if err != nil {
		// Handle non-gRPC errors
		logger.Errorf("%s err {failed to invoke grpc service provider, %s}", loggerHeader, err.Error())
		errResp := http.ServiceUnavailable.WithError(fmt.Errorf("gRPC invoke error: %w", err))
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		f.connections.Invalidate(connectionKey, clientConn)
		return filter.Stop
	}

	res, err := protoMsgToJson(resp)
	if err != nil {
		logger.Errorf("%s err {failed to convert proto msg to json, %s}", loggerHeader, err.Error())
		errResp := http.BadGateway.WithError(fmt.Errorf("protocol conversion error: %w", err))
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	h := mapMetadataToHeader(md)
	th := mapMetadataToHeader(t)

	// let response filter handle resp
	c.SourceResp = &stdHttp.Response{
		StatusCode: stdHttp.StatusOK,
		Header:     h,
		Body:       io.NopCloser(strings.NewReader(res)),
		Trailer:    th,
		Request:    c.Request,
	}
	return filter.Continue
}

func grpcConnectionKey(cluster, endpoint string) string {
	return cluster + "\x00" + endpoint
}

func (f *Filter) registerExtension(source DescriptorSource, mthDesc *desc.MethodDescriptor) error {
	inputDesc := mthDesc.GetInputType()
	outputDesc := mthDesc.GetOutputType()
	if f.extensionMu == nil || f.extReg == nil || f.registered == nil {
		if err := RegisterExtension(source, f.extReg, inputDesc, f.registered); err != nil {
			return perrors.New("register extension failed")
		}
		if err := RegisterExtension(source, f.extReg, outputDesc, f.registered); err != nil {
			return perrors.New("register extension failed")
		}
		return nil
	}

	inputName := inputDesc.GetFullyQualifiedName()
	outputName := outputDesc.GetFullyQualifiedName()
	f.extensionMu.RLock()
	registered := f.registered[inputName] && f.registered[outputName]
	f.extensionMu.RUnlock()
	if registered {
		return nil
	}

	f.extensionMu.Lock()
	defer f.extensionMu.Unlock()
	if !f.registered[inputName] {
		if err := RegisterExtension(source, f.extReg, inputDesc, f.registered); err != nil {
			return perrors.New("register extension failed")
		}
		f.registered[inputName] = true
	}
	if !f.registered[outputName] {
		if err := RegisterExtension(source, f.extReg, outputDesc, f.registered); err != nil {
			return perrors.New("register extension failed")
		}
		f.registered[outputName] = true
	}
	return nil
}

func RegisterExtension(source DescriptorSource, extReg *dynamic.ExtensionRegistry, msgDesc *desc.MessageDescriptor, registered map[string]bool) error {
	msgType := msgDesc.GetFullyQualifiedName()
	if _, ok := registered[msgType]; ok {
		return nil
	}

	if len(msgDesc.GetExtensionRanges()) > 0 {
		fds, err := source.AllExtensionsForType(msgType)
		if err != nil {
			return fmt.Errorf("failed to find msg type {%s} in file source", msgType)
		}

		err = extReg.AddExtension(fds...)
		if err != nil {
			return fmt.Errorf("failed to register extensions of msgType {%s}, err is {%s}", msgType, err.Error())
		}
	}

	for _, fd := range msgDesc.GetFields() {
		if fd.GetMessageType() != nil {
			err := RegisterExtension(source, extReg, fd.GetMessageType(), registered)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func mapHeaderToMetadata(header stdHttp.Header) metadata.MD {
	md := metadata.MD{}
	for key, val := range header {
		md.Append(key, val...)
	}
	return md
}

func mapMetadataToHeader(md metadata.MD) stdHttp.Header {
	h := stdHttp.Header{}
	for key, val := range md {
		for _, v := range val {
			h.Add(key, v)
		}
	}
	return h
}

func jsonToProtoMsg(reader io.Reader, msg proto.Message) error {
	body, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return jsonpb.UnmarshalString(string(body), msg)
}

func protoMsgToJson(msg proto.Message) (string, error) {
	m := jsonpb.Marshaler{}
	return m.MarshalToString(msg)
}

func isServerError(st *status.Status) bool {
	return st.Code() == codes.DeadlineExceeded || st.Code() == codes.ResourceExhausted || st.Code() == codes.Internal ||
		st.Code() == codes.Unavailable
}

func isServerTimeout(st *status.Status) bool {
	return st.Code() == codes.DeadlineExceeded || st.Code() == codes.Canceled
}

func (factory *FilterFactory) Config() any {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {

	err := configCheck(factory.cfg)
	if err != nil {
		return err
	}

	factory.descriptor.initDescriptorSource(factory.cfg)

	return nil
}

// Close releases all backend connections owned by this filter factory.
func (factory *FilterFactory) Close() error {
	var firstErr error
	if factory.connections != nil {
		if err := factory.connections.Close(); err != nil {
			firstErr = err
		}
	}
	if factory.descriptor != nil {
		factory.descriptor.Close()
	}
	if factory.removeEndpointHandler != nil {
		factory.removeEndpointHandler()
	}
	return firstErr
}

func configCheck(cfg *Config) error {
	if len(cfg.DescriptorSourceStrategy.Val()) == 0 {
		return perrors.Errorf("grpc descriptor source config `descriptor_source_strategy` is `%s`, maybe set it `%s`", cfg.DescriptorSourceStrategy.String(), AUTO)
	}
	return nil
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (config *Config) DeepCopy() *Config {
	if config == nil {
		return nil
	}

	cp := *config

	if config.Rules != nil {
		cp.Rules = make([]*Rule, len(config.Rules))
		for i, r := range config.Rules {
			if r == nil {
				cp.Rules[i] = nil
				continue
			}
			nr := &Rule{
				Selector: r.Selector,
				Match: Match{
					Method: r.Match.Method,
				},
			}
			cp.Rules[i] = nr
		}
	} else {
		cp.Rules = nil
	}

	return &cp
}
