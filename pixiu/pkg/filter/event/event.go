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

package event

import (
	"fmt"
	sdkhttp "net/http"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client/mq"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPEventFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}

	// FilterFactory is http filter instance
	FilterFactory struct {
		cfg *mq.Config
	}
	// Filter is http filter instance
	Filter struct {
		cfg *mq.Config
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &mq.Config{}}, nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{cfg: factory.cfg}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus {
	f.cfg.Timeout = ctx.Timeout
	mqClient := mq.NewSingletonMQClient(*f.cfg)
	req := client.NewReq(ctx.Request.Context(), ctx.Request, *ctx.GetAPI())
	req.Timeout = ctx.Timeout
	resp, err := mqClient.Call(req)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] event client call err:%v!", err)
		ctx.SendLocalReply(sdkhttp.StatusInternalServerError, []byte(fmt.Sprintf("event client call err:%v", err)))
		return filter.Stop
	}
	logger.Debugf("[dubbo-go-pixiu] event client call resp:%v", resp)
	ctx.SourceResp = resp
	return filter.Continue
}

func (factory *FilterFactory) Apply() error {
	mq.NewSingletonMQClient(*factory.cfg)
	return nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}
