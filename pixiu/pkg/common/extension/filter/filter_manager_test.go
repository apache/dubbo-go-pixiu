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
	"fmt"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	contexthttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

const (
	DEMO = "dgp.filters.demo"
	// Kind is the kind of plugin.
	Kind = DEMO
)

func init() {
	RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// HeaderFilter is http filter instance
	DemoFilterFactory struct {
		conf *Config
	}
	DemoFilter struct {
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

func (p *Plugin) CreateFilterFactory() (HttpFilterFactory, error) {
	return &DemoFilterFactory{conf: &Config{Foo: "default foo", Bar: "default bar"}}, nil
}

func (f *DemoFilter) Decode(ctx *contexthttp.HttpContext) FilterStatus {
	logger.Info("decode phase: ", f.str)

	runes := []rune(f.str)
	for i := 0; i < len(runes)/2; i += 1 {
		runes[i], runes[len(runes)-1-i] = runes[len(runes)-1-i], runes[i]
	}
	f.str = string(runes)

	return Continue
}

func (f *DemoFilter) Encode(ctx *contexthttp.HttpContext) FilterStatus {
	logger.Info("encode phase: ", f.str)
	return Continue
}

func (f *DemoFilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain FilterChain) error {
	c := f.conf
	str := fmt.Sprintf("%s is drinking in the %s", c.Foo, c.Bar)
	filter := &DemoFilter{str: str}

	chain.AppendDecodeFilters(filter)
	chain.AppendEncodeFilters(filter)
	return nil
}

func (f *DemoFilterFactory) Config() interface{} {
	return f.conf
}

func (f *DemoFilterFactory) Apply() error {
	return nil
}

func TestApply(t *testing.T) {
	fm := NewEmptyFilterManager()

	conf := map[string]interface{}{}
	conf["foo"] = "Cat"
	conf["bar"] = "The Walnut"
	f, err := fm.Apply(DEMO, conf)
	assert.Nil(t, err)

	baseContext := &contexthttp.HttpContext{}
	chain := NewDefaultFilterChain()
	_ = f.PrepareFilterChain(baseContext, chain)
	chain.OnDecode(baseContext)
}

func TestLoad(t *testing.T) {
	fm := NewEmptyFilterManager()
	conf := map[string]interface{}{}
	conf["foo"] = "Cat"
	conf["bar"] = "The Walnut"

	filtersConf := []*model.HTTPFilter{
		{
			Name:   DEMO,
			Config: conf,
		},
	}

	runFilter(t, fm, filtersConf)

	conf["foo"] = "Dog"
	conf["bar"] = "The Toilet"
	filtersConf = []*model.HTTPFilter{
		{
			Name:   DEMO,
			Config: conf,
		},
	}
	runFilter(t, fm, filtersConf)
}

func runFilter(t *testing.T, fm *FilterManager, filtersConf []*model.HTTPFilter) {
	fm.ReLoad(filtersConf)

	filters := fm.GetFactory()
	assert.Equal(t, len(filtersConf), len(filters))

	baseContext := &contexthttp.HttpContext{}
	baseContext.Reset()

	chain := fm.CreateFilterChain(baseContext)
	chain.OnDecode(baseContext)
	chain.OnEncode(baseContext)
}
