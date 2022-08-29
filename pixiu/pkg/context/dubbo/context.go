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

package dubbo

import (
	"context"
)

import (
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

// RpcContext the rpc invocation context
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
