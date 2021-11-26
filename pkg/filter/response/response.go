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

package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	// Kind is the kind of plugin.
	Kind = constant.HTTPResponseFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// HeaderFilter is http filter instance
	Filter struct {
		cfg *Config
	}
	// Config describe the config of Filter
	Config struct {
		Strategy string `json:"strategy,omitempty" yaml:"strategy,omitempty"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilterFactory, error) {
	return &Filter{cfg: &Config{}}, nil
}

func (f *Filter) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(c *contexthttp.HttpContext) {
	f.doResponse(c)
}

func (f *Filter) Config() interface{} {
	return f.cfg
}

func (f *Filter) Apply() error {
	if f.cfg.Strategy == "" {
		strategy := defaultNewParams()
		f.cfg.Strategy = strategy
	}
	return nil
}

func (f *Filter) doResponse(ctx *contexthttp.HttpContext) {
	// error do first
	if ctx.Err != nil {
		bt, _ := json.Marshal(contexthttp.ErrResponse{Message: ctx.Err.Error()})
		ctx.SourceResp = bt
		ctx.TargetResp = &client.Response{Data: bt}
		ctx.WriteJSONWithStatus(http.StatusServiceUnavailable, bt)
		ctx.Abort()
		return
	}

	// http type
	r, ok := ctx.SourceResp.(*http.Response)
	if ok {
		ctx.TargetResp = &client.Response{Data: r}
		byts, err := ioutil.ReadAll(r.Body)
		if err != nil {
			ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
			ctx.WriteWithStatus(r.StatusCode, []byte(r.Status))
			ctx.Abort()
			return
		}
		for k, v := range r.Header {
			ctx.AddHeader(k, v[0])
		}
		ctx.WriteWithStatus(r.StatusCode, byts)
		ctx.Abort()
		return
	}

	ctx.TargetResp = f.newResponse(ctx.SourceResp)
	ctx.WriteResponse(*ctx.TargetResp)
	ctx.Abort()
}

func (f *Filter) newResponse(data interface{}) *client.Response {
	hump := false
	if f.cfg.Strategy == constant.ResponseStrategyHump {
		hump = true
	}
	r, err := dealResp(data, hump)
	if err != nil {
		return &client.Response{Data: data}
	}

	return &client.Response{Data: r}
}

func defaultNewParams() string {
	strategy := os.Getenv(constant.EnvResponseStrategy)
	if len(strategy) != 0 {
		strategy = constant.ResponseStrategyNormal
	}
	return strategy
}

func dealResp(in interface{}, HumpToLine bool) (interface{}, error) {
	if in == nil {
		return in, nil
	}
	switch reflect.TypeOf(in).Kind() {
	case reflect.Map:
		if _, ok := in.(map[interface{}]interface{}); ok {
			m := mapIItoMapSI(in)
			if HumpToLine {
				m = humpToLine(m)
			}
			return m, nil
		} else if inm, ok := in.(map[string]interface{}); ok {
			if HumpToLine {
				m := humpToLine(in)
				return m, nil
			}
			return inm, nil
		}
	case reflect.Slice:
		value := reflect.ValueOf(in)
		newTemps := make([]interface{}, 0, value.Len())
		for i := 0; i < value.Len(); i++ {
			if value.Index(i).CanInterface() {
				newTemp, e := dealResp(value.Index(i).Interface(), HumpToLine)
				if e != nil {
					return nil, e
				}
				newTemps = append(newTemps, newTemp)
			} else {
				return nil, errors.New(fmt.Sprintf("unexpect err,value:%+v", value))
			}
		}
		return newTemps, nil
	default:
		return in, nil
	}
	return in, nil
}

func mapIItoMapSI(in interface{}) interface{} {
	var inMap map[interface{}]interface{}
	if v, ok := in.(map[interface{}]interface{}); !ok {
		return in
	} else {
		inMap = v
	}
	outMap := make(map[string]interface{}, len(inMap))

	for k, v := range inMap {
		if v == nil {
			continue
		}
		s := fmt.Sprint(k)
		if s == "class" {
			// ignore the "class" field
			continue
		}
		vt := reflect.TypeOf(v)
		switch vt.Kind() {
		case reflect.Map:
			if _, ok := v.(map[interface{}]interface{}); ok {
				v = mapIItoMapSI(v)
			}
		case reflect.Slice:
			vl := reflect.ValueOf(v)
			os := make([]interface{}, 0, vl.Len())
			for i := 0; i < vl.Len(); i++ {
				if vl.Index(i).CanInterface() {
					osv := mapIItoMapSI(vl.Index(i).Interface())
					os = append(os, osv)
				}
			}
			v = os
		}
		outMap[s] = v

	}
	return outMap
}

// traverse all the keys in the map and change the hump to underline
func humpToLine(in interface{}) interface{} {
	var m map[string]interface{}
	if v, ok := in.(map[string]interface{}); ok {
		m = v
	} else {
		return in
	}

	out := make(map[string]interface{}, len(m))
	for k1, v1 := range m {
		x := humpToUnderline(k1)

		if v1 == nil {
			out[x] = v1
		} else if reflect.TypeOf(v1).Kind() == reflect.Struct {
			out[x] = humpToLine(struct2Map(v1))
		} else if reflect.TypeOf(v1).Kind() == reflect.Slice {
			value := reflect.ValueOf(v1)
			newTemps := make([]interface{}, 0, value.Len())
			for i := 0; i < value.Len(); i++ {
				newTemp := humpToLine(value.Index(i).Interface())
				newTemps = append(newTemps, newTemp)
			}
			out[x] = newTemps
		} else if reflect.TypeOf(v1).Kind() == reflect.Map {
			out[x] = humpToLine(v1)
		} else {
			out[x] = v1
		}
	}
	return out
}

func humpToUnderline(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}

func struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	data := make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
