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
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

import (
	pxCtx "github.com/apache/dubbo-go-pixiu/pkg/context"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func init() {
	RegisterFilterFactory(DEMO, newDemoFilter)
}

const DEMO = "dgp.filters.demo"

type DemoFilter struct {
	conf *Config
}

type Config struct {
	Foo string `json:"foo,omitempty" yaml:"foo,omitempty"`
	Bar string `json:"bar,omitempty" yaml:"bar,omitempty"`
}

func newDemoFilter() filter.Factory {
	return &DemoFilter{conf: &Config{Foo: "default foo", Bar: "default bar"}}
}

func (d *DemoFilter) Config() interface{} {
	return d.conf
}

func (d *DemoFilter) Apply() (filter.Filter, error) {
	//init every sth before return the func
	c := d.conf
	str := fmt.Sprintf("%s is drinking in the %s", c.Foo, c.Bar)
	//return the filter func
	return func(ctx context.Context) {
		logger.Info(str)
	}, nil
}

func TestApply(t *testing.T) {
	fm := NewFilterManager()

	conf := map[string]interface{}{}
	conf["foo"] = "Cat"
	conf["bar"] = "The Walnut"
	apply, err := fm.Apply(DEMO, conf)
	assert.Nil(t, err)

	baseContext := pxCtx.NewBaseContext()
	apply(baseContext)
}

func TestLoad(t *testing.T) {
	fm := NewFilterManager()
	conf := map[string]interface{}{}
	conf["foo"] = "Cat"
	conf["bar"] = "The Walnut"

	filtersConf := []*config.Filter{
		{
			Name:   DEMO,
			Config: conf,
		},
	}
	fm.Load(filtersConf)

	filters := fm.GetFilters()
	assert.Equal(t, len(filtersConf), len(filters))

	baseContext := pxCtx.NewBaseContext()
	for i := range filters {
		filters[i](baseContext)
	}
}
