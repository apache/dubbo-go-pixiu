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

package host

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	http2 "github.com/apache/dubbo-go-pixiu/pkg/common/http"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	// Kind is the kind of plugin.
	Kind = constant.HTTPHOSTFilter
)

func init() {
	extension.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// HostFilter is http filter instance
	HostFilter struct {
		cfg *Config
	}
	// Config describe the config of HostFilter
	Config struct {
		Host string `yaml:"host" json:"host"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter(hcm *http2.HttpConnectionManager, config interface{}, bs *model.Bootstrap) (extension.HttpFilter, error) {
	specConfig := config.(Config)
	return &HostFilter{cfg: &specConfig}, nil
}

func (hf *HostFilter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(hf.Handle)
	return nil
}

func (hf *HostFilter) Handle(c *contexthttp.HttpContext) {
	c.Request.Host = hf.cfg.Host
	c.Next()
}
