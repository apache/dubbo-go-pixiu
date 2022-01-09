package dubbo

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type RpcContext struct {
	RpcInvocation *invocation.RPCInvocation
	RpcResult     protocol.Result
	Route         *model.RouteAction
	Ctx           context.Context
}

// SetError set error in RpcResult
func (c *RpcContext) SetError(err error) {
	if c.RpcResult == nil {
		c.RpcResult = &protocol.RPCResult{}
	}
	c.RpcResult.SetError(err)
}

// SetResult set successful result in RpcResult
func (c *RpcContext) SetResult(result protocol.Result) {
	c.RpcResult = result
}

// SetInvocation set invocation
func (c *RpcContext) SetInvocation(invocation *invocation.RPCInvocation) {
	c.RpcInvocation = invocation
}

// SetRoute set route
func (c *RpcContext) SetRoute(route *model.RouteAction) {
	c.Route = route
}
