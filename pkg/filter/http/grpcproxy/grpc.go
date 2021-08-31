package grpcproxy

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	stdHttp "net/http"
	"strings"
	"sync"

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPGrpcProxyFilter
)

var (
	fs fileSource
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is grpc filter plugin.
	Plugin struct {
	}

	// Filter is grpc filter instance
	Filter struct {
		cfg *Config
		// hold grpc.ClientConns
		pool sync.Pool
	}

	// Config describe the config of AccessFilter
	Config struct {
		Proto string  `yaml:"proto_descriptor" json:"proto_descriptor"`
		rules []*Rule `yaml:"rules" json:"rules"`
	}

	Rule struct {
		Selector string `yaml:"selector" json:"selector"`
		Match    Match  `yaml:"match" json:"match"`
	}

	Match struct {
		method string `yaml:"method" json:"method"`
	}
)

func (ap *Plugin) Kind() string {
	return Kind
}

func (ap *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &Config{}}, nil
}

func (af *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(af.Handle)
	return nil
}

// Handle use the default http to grpc transcoding strategy https://cloud.google.com/endpoints/docs/grpc/transcoding
func (af *Filter) Handle(c *http.HttpContext) {
	paths := strings.Split(c.Request.URL.Path, "/")
	if len(paths) < 2 {
		writeResp(c, stdHttp.StatusBadRequest, "request path invalid")
		return
	}
	svc := paths[0]
	mth := paths[1]

	dscp, err := fs.FindSymbol(svc + "." + mth)
	if err != nil {
		writeResp(c, stdHttp.StatusMethodNotAllowed, "method not found")
		return
	}

	svcDesc, ok := dscp.(*desc.ServiceDescriptor)
	if !ok {
		writeResp(c, stdHttp.StatusMethodNotAllowed, fmt.Sprintf("service not expose, %s", svc))
		return
	}

	mthDesc := svcDesc.FindMethodByName(mth)

	// TODO(Kenway): Can extension registry being cached ?
	var extReg dynamic.ExtensionRegistry
	registered := make(map[string]bool)
	err = RegisterExtension(&extReg, mthDesc.GetInputType(), registered)
	if err != nil {
		writeResp(c, stdHttp.StatusInternalServerError, "register extension failed")
		return
	}

	err = RegisterExtension(&extReg, mthDesc.GetOutputType(), registered)
	if err != nil {
		writeResp(c, stdHttp.StatusInternalServerError, "register extension failed")
		return
	}

	msgFac := dynamic.NewMessageFactoryWithExtensionRegistry(&extReg)
	grpcReq := msgFac.NewMessage(mthDesc.GetInputType())

	err = jsonToProtoMsg(c.Request.Body, grpcReq)
	if err != nil {
		writeResp(c, stdHttp.StatusBadRequest, "convert to proto msg failed")
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var clientConn *grpc.ClientConn
	clientConn = af.pool.Get().(*grpc.ClientConn)
	if clientConn == nil {
		// TODO(Kenway): Support Credential and TLS
		clientConn, err = grpc.DialContext(ctx, c.GetAPI().IntegrationRequest.Host, grpc.WithInsecure())
		if err != nil || clientConn == nil {
			writeResp(c, stdHttp.StatusServiceUnavailable, "connect to grpc server failed")
			return
		}
		af.pool.Put(clientConn)
	}

	stub := grpcdynamic.NewStubWithMessageFactory(clientConn, msgFac)

	resp, err := Invoke(ctx, stub, mthDesc, grpcReq)
	if st, ok := status.FromError(err); !ok || isServerError(st) {
		writeResp(c, stdHttp.StatusInternalServerError, err.Error())
		return
	}

	res, err := protoMsgToJson(resp)
	if err != nil {
		writeResp(c, stdHttp.StatusInternalServerError, "serialize proto msg to json failed")
		return
	}
	writeResp(c, stdHttp.StatusOK, res)
}

func RegisterExtension(extReg *dynamic.ExtensionRegistry, msgDesc *desc.MessageDescriptor, registered map[string]bool) error {
	msgType := msgDesc.GetFullyQualifiedName()
	if _, ok := registered[msgType]; ok {
		return nil
	}

	if len(msgDesc.GetExtensionRanges()) > 0 {
		fds, err := fs.AllExtensionsForType(msgType)
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
			err := RegisterExtension(extReg, fd.GetMessageType(), registered)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func jsonToProtoMsg(reader io.Reader, msg proto.Message) error {
	return jsonpb.Unmarshal(reader, msg)
}

func protoMsgToJson(msg proto.Message) (string, error) {
	m := jsonpb.Marshaler{}
	return m.MarshalToString(msg)
}

func writeResp(c *http.HttpContext, code int, msg string) {
	resp, _ := json.Marshal(http.ErrResponse{Message: msg})
	c.WriteJSONWithStatus(code, resp)
}

func isServerError(st *status.Status) bool {
	return st.Code() == codes.DeadlineExceeded || st.Code() == codes.ResourceExhausted || st.Code() == codes.Internal ||
		st.Code() == codes.Unavailable
}

func (af *Filter) Config() interface{} {
	return af.cfg
}

func (af *Filter) Apply() error {
	gc := config.GetBootstrap().StaticResources.GRPCConfig
	err := af.initFromFileDescriptor([]string{gc.ProtoPath}, gc.FileNames...)
	if err != nil {
		logger.Errorf("failed to init proto files, %s", err.Error())
	}

	return nil
}
