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

package timeout

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPTimeoutFilter
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
		cfg *Config
	}

	Config struct {
		Timeout time.Duration `yaml:"timeout" json:"timeout"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	specConfig := Config{constant.DefaultTimeout}
	return &Filter{cfg: &specConfig}, nil
}

func (f *Filter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(hc *contexthttp.HttpContext) {
	ctx, cancel := context.WithTimeout(hc.Ctx, f.getTimeout(hc.Timeout))
	defer cancel()
	// Channel capacity must be greater than 0.
	// Otherwise, if the parent coroutine quit due to timeout,
	// the child coroutine may never be able to quit.
	finishChan := make(chan struct{}, 1)
	go func() {
		// panic by recovery
		hc.Next()
		finishChan <- struct{}{}
	}()

	select {
	// timeout do.
	case <-ctx.Done():
		logger.Warnf("api:%s request timeout", hc.GetUrl())
		bt, _ := json.Marshal(contexthttp.ErrResponse{Message: http.ErrHandlerTimeout.Error()})
		hc.SourceResp = bt
		hc.TargetResp = &client.Response{Data: bt}
		hc.WriteJSONWithStatus(http.StatusGatewayTimeout, bt)
		hc.Abort()
	case <-finishChan:
		// finish call do something.
	}
}

func (f *Filter) Config() interface{} {
	return nil
}

func (f *Filter) Apply() error {
	if f.cfg.Timeout <= 0 {
		f.cfg.Timeout = constant.DefaultTimeout
	}
	return nil
}

func (f *Filter) getTimeout(t time.Duration) time.Duration {
	if t <= 0 {
		return f.cfg.Timeout
	}

	return t
}
