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

package filter

import (
	"io/ioutil"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

func init() {
	extension.SetFilterFunc(constant.RemoteCallFilter, RemoteCall())
}

// RemoteCall http 2 dubbo
func RemoteCall() context.FilterFunc {
	return func(c context.Context) {
		doRemoteCall(c.(*http.HttpContext))
	}
}

func doRemoteCall(c *http.HttpContext) {
	api := c.GetApi()
	cl, e := client.SingletonPool().GetClient(api.IType)
	if e != nil {
		c.WriteFail()
		c.AbortWithError("", e)
	}
	if bytes, err := ioutil.ReadAll(c.Request.Body); err != nil {
		logger.Errorf("[dubboproxy go] read body err:%v!", err)
		c.WriteFail()
		c.Abort()
	} else {
		if resp, err := cl.Call(client.NewRequest(bytes, api)); err != nil {
			logger.Errorf("[dubboproxy go] client do err:%v!", err)
			c.WriteFail()
			c.Abort()
		} else {
			c.WriteResponse(resp)
			c.Next()
		}
	}
}
