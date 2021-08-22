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

package recovery

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// Kind is the kind of plugin.
	Kind = constant.RecoveryFilter
)

func init() {
	extension.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// RecoveryFilter is http filter instance
	RecoveryFilter struct {
		cfg *Config
	}
	// Config describe the config of RecoveryFilter
	Config struct {
		Host string `yaml:"host" json:"host"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (extension.HttpFilter, error) {
	return &RecoveryFilter{}, nil
}

func (rf *RecoveryFilter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(rf.Handle)
	return nil
}

func (rf *RecoveryFilter) Handle(c *http.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			logger.Warnf("[dubbopixiu go] error:%+v", err)
			c.WriteErr(err)
		}
	}()
	c.Next()
}

func (f *RecoveryFilter) Config() interface{} {
	return nil
}

func (f *RecoveryFilter) Apply() error {
	return nil
}

// GetMock return mocked filter
func GetMock() extension.HttpFilter {
	filter := &RecoveryFilter{}
	_ = filter.Apply()
	return filter
}
