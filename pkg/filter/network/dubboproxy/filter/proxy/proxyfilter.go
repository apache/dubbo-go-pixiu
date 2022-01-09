package proxy

import (
	"context"
	"reflect"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"

	tpconst "github.com/dubbogo/triple/pkg/common/constant"
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	dubbo2 "github.com/apache/dubbo-go-pixiu/pkg/context/dubbo"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	Kind = constant.DubboProxyFilter
)

func init() {
	filter.RegisterDubboFilterPlugin(&Plugin{})
}

type (
	Plugin struct {
	}

	Config struct {
		Protocol string `default:"protocol" yaml:"protocol" json:"protocol" mapstructure:"protocol"`
	}

	Filter struct {
		Config *Config
	}
)

func (p Plugin) Kind() string {
	return Kind
}

// CreateFilterFactory return the filter callback
func (p Plugin) CreateFilter(config interface{}) (filter.DubboFilter, error) {
	cfg := config.(*Config)
	return Filter{Config: cfg}, nil
}

// Config Expose the config so that Filter Manger can inject it, so it must be a pointer
func (p Plugin) Config() interface{} {
	return &Config{}
}

func (f Filter) Handle(ctx *dubbo2.RpcContext) filter.FilterStatus {
	switch f.Config.Protocol {
	case dubbo.DUBBO:
		return f.sendDubboRequest(ctx)
	case tripleConstant.TRIPLE:
		return f.sendTripleRequest(ctx)
	default:
		return filter.Stop
	}
}

func (f Filter) sendDubboRequest(ctx *dubbo2.RpcContext) filter.FilterStatus {

	ra := ctx.Route
	clusterName := ra.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		ctx.SetError(errors.Errorf("Requested dubbo rpc invocation endpoint not found"))
		return filter.Stop
	}

	invoc := ctx.RpcInvocation
	url, err := common.NewURL(endpoint.Address.GetAddress(),
		common.WithProtocol(dubbo.DUBBO), common.WithParamsValue(dubboConstant.SerializationKey, dubboConstant.Hessian2Serialization),
		common.WithParamsValue(dubboConstant.GenericFilterKey, "true"),
		common.WithParamsValue(dubboConstant.InterfaceKey, invoc.AttachmentsByKey(dubboConstant.InterfaceKey, "")),
		common.WithParamsValue(dubboConstant.ReferenceFilterKey, "generic,filter"),
		// dubboAttachment must contains group and version info
		common.WithParamsValue(dubboConstant.GroupKey, invoc.AttachmentsByKey(dubboConstant.GroupKey, "")),
		common.WithParamsValue(dubboConstant.VersionKey, invoc.AttachmentsByKey(dubboConstant.VersionKey, "")),
		common.WithPath(invoc.AttachmentsByKey(dubboConstant.InterfaceKey, "")),
	)

	if err != nil {
		ctx.SetError(err)
		return filter.Stop
	}

	dubboProtocol := dubbo.NewDubboProtocol()

	// TODO: will lead to Error when failed to connect server
	invoker := dubboProtocol.Refer(url)

	var resp interface{}
	invoc.SetReply(&resp)

	invCtx := context.Background()
	result := invoker.Invoke(invCtx, invoc)
	result.SetAttachments(invoc.Attachments())
	if result.Error() != nil {
		ctx.SetError(result.Error())
		return filter.Stop
	}

	value := reflect.ValueOf(result.Result())
	result.SetResult(value.Elem().Interface())
	ctx.SetResult(result)
	return filter.Continue
}

func (f Filter) sendTripleRequest(ctx *dubbo2.RpcContext) filter.FilterStatus {

	invoc := ctx.RpcInvocation

	ra := ctx.Route

	clusterName := ra.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		ctx.SetError(errors.Errorf("Requested dubbo rpc invocation endpoint not found"))
		return filter.Stop
	}

	path := invoc.Attachment(dubboConstant.PathKey).(string)
	// create URL from RpcInvocation
	url, err := common.NewURL(endpoint.Address.GetAddress(),
		common.WithProtocol(tpconst.TRIPLE), common.WithParamsValue(dubboConstant.SerializationKey, dubboConstant.Hessian2Serialization),
		common.WithParamsValue(dubboConstant.GenericFilterKey, "true"),
		common.WithParamsValue(dubboConstant.AppVersionKey, "3.0.0"),
		common.WithParamsValue(dubboConstant.InterfaceKey, path),
		common.WithParamsValue(dubboConstant.ReferenceFilterKey, "generic,filter"),
		common.WithPath(path),
	)

	if err != nil {
		ctx.SetError(err)
		return filter.Stop
	}

	invoker, err := dubbo3.NewDubboInvoker(url)

	if err != nil {
		ctx.SetError(err)
		return filter.Stop
	}

	invCtx := context.Background()

	var resp interface{}
	invoc.SetReply(&resp)
	result := invoker.Invoke(invCtx, invoc)

	if result.Error() != nil {
		ctx.SetError(result.Error())
		return filter.Stop
	}

	// when upstream server down, the result and error in result are both nil
	if result.Result() == nil {
		ctx.SetError(errors.New("result from upstream server is nil"))
		return filter.Stop
	}

	result.SetAttachments(invoc.Attachments())
	value := reflect.ValueOf(result.Result())
	result.SetResult(value.Elem().Interface())
	ctx.SetResult(result)
	return filter.Continue
}
