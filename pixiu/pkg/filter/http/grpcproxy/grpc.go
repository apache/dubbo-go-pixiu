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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	ct "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
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
		// hold grpc.ClientConns, key format: cluster name + "." + endpoint
		pools map[string]*sync.Pool

		extReg     *dynamic.ExtensionRegistry
		registered map[string]bool
	}
	Filter struct {
		cfg *Config
		// grpc descriptor source factory
		descriptor *Descriptor
		// hold grpc.ClientConns, key format: cluster name + "." + endpoint
		pools map[string]*sync.Pool

		extReg     *dynamic.ExtensionRegistry
		registered map[string]bool
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
	return &FilterFactory{cfg: &Config{DescriptorSourceStrategy: AUTO}, descriptor: &Descriptor{}}, nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{cfg: factory.cfg, descriptor: factory.descriptor, pools: factory.pools, extReg: factory.extReg, registered: factory.registered}
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

	e := server.GetClusterManager().PickEndpoint(re.Cluster)
	if e == nil {
		logger.Errorf("%s err {cluster not exists}", loggerHeader)
		c.SendLocalReply(stdHttp.StatusServiceUnavailable, []byte("cluster not exists"))
		return filter.Stop
	}
	// timeout for Dial and Invoke
	ctx, cancel := context.WithTimeout(c.Ctx, c.Timeout)
	defer cancel()
	ep := e.Address.GetAddress()

	p, ok := f.pools[strings.Join([]string{re.Cluster, ep}, ".")]
	if !ok {
		p = &sync.Pool{}
	}

	clientConn, ok = p.Get().(*grpc.ClientConn)
	if !ok || clientConn == nil {
		// TODO(Kenway): Support Credential and TLS
		clientConn, err = grpc.DialContext(ctx, ep, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil || clientConn == nil {
			logger.Errorf("%s err {failed to connect to grpc service provider}", loggerHeader)
			c.SendLocalReply(stdHttp.StatusServiceUnavailable, []byte((fmt.Sprintf("%s", err))))
			return filter.Stop
		}
	}

	// get DescriptorSource, contain file and reflection
	source, err := f.descriptor.getDescriptorSource(context.WithValue(ctx, ct.ContextKey(GrpcClientConnKey), clientConn), f.cfg)
	if err != nil {
		logger.Errorf("%s err %s : %s ", loggerHeader, "get desc source fail", err)
		c.SendLocalReply(stdHttp.StatusInternalServerError, []byte("service not config proto file or the server not support reflection API"))
		return filter.Stop
	}
	//put DescriptorSource concurrent, del if no need
	ctx = context.WithValue(ctx, ct.ContextKey(DescriptorSourceKey), source)

	dscp, err := source.FindSymbol(svc)
	if err != nil {
		logger.Errorf("%s err {%s}", loggerHeader, "request path invalid")
		c.SendLocalReply(stdHttp.StatusBadRequest, []byte("method not allow"))
		return filter.Stop
	}

	svcDesc, ok := dscp.(*desc.ServiceDescriptor)
	if !ok {
		logger.Errorf("%s err {service not expose, %s}", loggerHeader, svc)
		c.SendLocalReply(stdHttp.StatusBadRequest, []byte(fmt.Sprintf("service not expose, %s", svc)))
		return filter.Stop
	}

	mthDesc := svcDesc.FindMethodByName(mth)

	err = f.registerExtension(source, mthDesc)
	if err != nil {
		logger.Errorf("%s err {%s}", loggerHeader, "register extension failed")
		c.SendLocalReply(stdHttp.StatusInternalServerError, []byte(fmt.Sprintf("%s", err)))
		return filter.Stop
	}

	msgFac := dynamic.NewMessageFactoryWithExtensionRegistry(f.extReg)
	grpcReq := msgFac.NewMessage(mthDesc.GetInputType())

	err = jsonToProtoMsg(c.Request.Body, grpcReq)
	if err != nil && !errors.Is(err, io.EOF) {
		logger.Errorf("%s err {failed to convert json to proto msg, %s}", loggerHeader, err.Error())
		c.SendLocalReply(stdHttp.StatusInternalServerError, []byte(fmt.Sprintf("%s", err)))
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
	if st, ok := status.FromError(err); !ok || isServerError(st) {
		logger.Errorf("%s err {failed to invoke grpc service provider, %s}", loggerHeader, err.Error())
		c.SendLocalReply(stdHttp.StatusServiceUnavailable, []byte(fmt.Sprintf("%s", err)))
		return filter.Stop
	}

	res, err := protoMsgToJson(resp)
	if err != nil {
		logger.Errorf("%s err {failed to convert proto msg to json, %s}", loggerHeader, err.Error())
		c.SendLocalReply(stdHttp.StatusInternalServerError, []byte(fmt.Sprintf("%s", err)))
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
	p.Put(clientConn)
	return filter.Continue
}

func (f *Filter) registerExtension(source DescriptorSource, mthDesc *desc.MethodDescriptor) error {
	err := RegisterExtension(source, f.extReg, mthDesc.GetInputType(), f.registered)
	if err != nil {
		return perrors.New("register extension failed")
	}

	err = RegisterExtension(source, f.extReg, mthDesc.GetOutputType(), f.registered)
	if err != nil {
		return perrors.New("register extension failed")
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

func (factory *FilterFactory) Config() interface{} {
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

func configCheck(cfg *Config) error {
	if len(cfg.DescriptorSourceStrategy.Val()) == 0 {
		return perrors.Errorf("grpc descriptor source config `descriptor_source_strategy` is `%s`, maybe set it `%s`", cfg.DescriptorSourceStrategy.String(), AUTO)
	}
	return nil
}
