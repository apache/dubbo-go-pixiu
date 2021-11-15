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
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/client/mq"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
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

	// Filter is http filter instance
	Filter struct {
		cfg *mq.Config
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &mq.Config{}}, nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(ctx *http.HttpContext) {
	mqClient := mq.NewSingletonMQClient(*f.cfg)
	req := client.NewReq(ctx.Request.Context(), ctx.Request, *ctx.GetAPI())
	resp, err := mqClient.Call(req)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] event client call err:%v!", err)
		ctx.Err = err
		ctx.Next()
		return
	}
	logger.Debugf("[dubbo-go-pixiu] event client call resp:%v", resp)
	ctx.SourceResp = resp
	ctx.Next()
}

func (f *Filter) Apply() error {
	mq.NewSingletonMQClient(*f.cfg)
	return nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}
