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
 * WITHOmanage.goUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package filter

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func init() {
	extension.RegisterHttpFilter(&Plugin{})
}

const (
	DEMO = "dgp.filters.demo"
	// Kind is the kind of plugin.
	Kind = constant.ResponseFilter
)

func init() {
	extension.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// HeaderFilter is http filter instance
	DemoFilter struct {
		conf *Config
		str string
	}
	// Config describe the config of ResponseFilter
	Config struct {
		Foo string `json:"foo,omitempty" yaml:"foo,omitempty"`
		Bar string `json:"bar,omitempty" yaml:"bar,omitempty"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (extension.HttpFilter, error) {
	return &DemoFilter{conf: &Config{Foo: "default foo", Bar: "default bar"}},nil
}

func (f *DemoFilter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *DemoFilter) Handle(c *contexthttp.HttpContext) {
	logger.Info(f.str)
}

func (f *DemoFilter) Config() interface{} {
	return f.conf
}

func (f *DemoFilter) Apply() error {
	c := f.conf
	f.str = fmt.Sprintf("%s is drinking in the %s", c.Foo, c.Bar)
	//return the filter func
	return nil
}

func TestApply(t *testing.T) {
	fm := NewFilterManager()

	conf := map[string]interface{}{}
	conf["foo"] = "Cat"
	conf["bar"] = "The Walnut"
	f, err := fm.Apply(DEMO, conf)
	assert.Nil(t, err)

	baseContext := &contexthttp.HttpContext{}
	f.Handle(baseContext)
}

func TestLoad(t *testing.T) {
	fm := NewFilterManager()
	conf := map[string]interface{}{}
	conf["foo"] = "Cat"
	conf["bar"] = "The Walnut"

	filtersConf := []*model.Filter{
		{
			Name:   DEMO,
			Config: conf,
		},
	}
	fm.Load(filtersConf)

	filters := fm.GetFilters()
	assert.Equal(t, len(filtersConf), len(filters))

	baseContext := &contexthttp.HttpContext{}

	for i := range filters {
		filters[i].Handle(baseContext)
	}
}
