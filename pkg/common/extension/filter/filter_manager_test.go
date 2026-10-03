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
	"sync"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
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
	prepareErrorFactory  struct{}
	closeTrackingFactory struct {
		closes         int
		panicOnPrepare bool
	}
	// Config describe the config of ResponseFilter
	Config struct {
		Foo string `json:"foo,omitempty" yaml:"foo,omitempty"`
		Bar string `json:"bar,omitempty" yaml:"bar,omitempty"`
	}
)

func (f *prepareErrorFactory) Config() any  { return &struct{}{} }
func (f *prepareErrorFactory) Apply() error { return nil }
func (f *prepareErrorFactory) PrepareFilterChain(*contexthttp.HttpContext, FilterChain) error {
	return fmt.Errorf("prepare rejected")
}

func (f *closeTrackingFactory) Config() any  { return &struct{}{} }
func (f *closeTrackingFactory) Apply() error { return nil }
func (f *closeTrackingFactory) PrepareFilterChain(*contexthttp.HttpContext, FilterChain) error {
	if f.panicOnPrepare {
		panic("prepare panic")
	}
	return nil
}
func (f *closeTrackingFactory) Close() error {
	f.closes++
	return nil
}

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

func (f *DemoFilterFactory) Config() any {
	return f.conf
}

func (f *DemoFilterFactory) Apply() error {
	return nil
}

func TestApply(t *testing.T) {
	fm := NewEmptyFilterManager()

	conf := map[string]any{}
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
	conf := map[string]any{}
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

func TestReloadSkipsFailedFactoriesAndDoesNotReopenAfterClose(t *testing.T) {
	fm := NewEmptyFilterManager()
	fm.ReLoad([]*model.HTTPFilter{{Name: "missing-filter"}})
	assert.Empty(t, fm.GetFactory())
	assert.NotPanics(t, func() { fm.CreateFilterChain(&contexthttp.HttpContext{}).Release() })

	fm.ReLoad([]*model.HTTPFilter{{Name: DEMO}})
	assert.Len(t, fm.GetFactory(), 1)
	assert.NoError(t, fm.Close())
	fm.ReLoad([]*model.HTTPFilter{{Name: DEMO}})
	assert.Empty(t, fm.GetFactory())
}

func TestReloadAndCloseDoNotReopenManager(t *testing.T) {
	fm := NewEmptyFilterManager()
	filters := []*model.HTTPFilter{{Name: DEMO}}
	var wg sync.WaitGroup
	wg.Go(func() {
		for range 20 {
			fm.ReLoad(filters)
		}
	})
	wg.Go(func() {
		assert.NoError(t, fm.Close())
	})
	wg.Wait()
	assert.Empty(t, fm.GetFactory())
}

func runFilter(t *testing.T, fm *FilterManager, filtersConf []*model.HTTPFilter) {
	assert.NoError(t, fm.ReLoadChecked(filtersConf))

	filters := fm.GetFactory()
	assert.Equal(t, len(filtersConf), len(filters))

	baseContext := &contexthttp.HttpContext{}
	baseContext.Reset()

	chain, err := fm.CreateFilterChainChecked(baseContext)
	assert.NoError(t, err)
	defer chain.Release()
	chain.OnDecode(baseContext)
	chain.OnEncode(baseContext)
}

func TestCreateFilterChainReturnsPrepareError(t *testing.T) {
	factory := HttpFilterFactory(&prepareErrorFactory{})
	fm := NewEmptyFilterManager()
	fm.filtersArray = []*HttpFilterFactory{&factory}

	_, err := fm.CreateFilterChainChecked(&contexthttp.HttpContext{})
	assert.ErrorContains(t, err, "prepare rejected")
}

func TestCreateFilterChainPreservesBestEffortCompatibility(t *testing.T) {
	rejected := HttpFilterFactory(&prepareErrorFactory{})
	accepted := HttpFilterFactory(&DemoFilterFactory{conf: &Config{Foo: "Cat", Bar: "The Walnut"}})
	fm := NewEmptyFilterManager()
	fm.filtersArray = []*HttpFilterFactory{&rejected, &accepted}

	chain := fm.CreateFilterChain(&contexthttp.HttpContext{})
	defer chain.Release()
	legacyChain, ok := chain.(*defaultFilterChain)
	assert.True(t, ok)
	assert.Len(t, legacyChain.decodeFilters, 1)
	assert.Len(t, legacyChain.encodeFilters, 1)
}

func TestCreatedFilterChainsReleaseRetiredFactories(t *testing.T) {
	builders := []struct {
		name  string
		build func(*FilterManager) (LeasedFilterChain, error)
	}{
		{"best effort", func(fm *FilterManager) (LeasedFilterChain, error) {
			return fm.CreateFilterChain(&contexthttp.HttpContext{}), nil
		}},
		{"checked", func(fm *FilterManager) (LeasedFilterChain, error) {
			return fm.CreateFilterChainChecked(&contexthttp.HttpContext{})
		}},
	}
	for _, builder := range builders {
		t.Run(builder.name, func(t *testing.T) {
			tracking := &closeTrackingFactory{}
			factory := HttpFilterFactory(tracking)
			fm := NewEmptyFilterManager()
			fm.filtersArray = []*HttpFilterFactory{&factory}

			chain, err := builder.build(fm)
			assert.NoError(t, err)
			assert.NoError(t, fm.ReLoadChecked(nil))
			assert.Equal(t, 0, tracking.closes)

			chain.Release()
			chain.Release()
			assert.Equal(t, 1, tracking.closes)
		})
	}
}

func TestCreatedFilterChainsReleaseFactoriesOnPreparePanic(t *testing.T) {
	builders := []struct {
		name  string
		build func(*FilterManager)
	}{
		{"best effort", func(fm *FilterManager) {
			fm.CreateFilterChain(&contexthttp.HttpContext{})
		}},
		{"checked", func(fm *FilterManager) {
			_, _ = fm.CreateFilterChainChecked(&contexthttp.HttpContext{})
		}},
	}
	for _, builder := range builders {
		t.Run(builder.name, func(t *testing.T) {
			tracking := &closeTrackingFactory{panicOnPrepare: true}
			factory := HttpFilterFactory(tracking)
			fm := NewEmptyFilterManager()
			fm.filtersArray = []*HttpFilterFactory{&factory}

			assert.PanicsWithValue(t, "prepare panic", func() { builder.build(fm) })
			assert.NoError(t, fm.ReLoadChecked(nil))
			assert.Equal(t, 1, tracking.closes)
		})
	}
}
