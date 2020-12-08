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
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	contexthttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter"
)

// nolint
func Init() {
	extension.SetFilterFunc(constant.ResponseFilter, responseFilterFunc())
}

func responseFilterFunc() selfcontext.FilterFunc {
	return New(defaultNewParams()).Do()
}

func defaultNewParams() string {
	strategy := os.Getenv(constant.EnvResponseStrategy)
	if len(strategy) != 0 {
		strategy = constant.ResponseStrategyNormal
	}
	return strategy
}

type responseFilter struct {
	strategy string
}

// New create timeout filter.
func New(strategy string) filter.Filter {
	return &responseFilter{
		strategy: strategy,
	}
}

// Do execute responseFilter filter logic.
func (f responseFilter) Do() selfcontext.FilterFunc {
	return func(c selfcontext.Context) {
		f.doResponse(c.(*contexthttp.HttpContext))
	}
}

func (f responseFilter) doResponse(ctx *contexthttp.HttpContext) {
	// error do first
	if ctx.Err != nil {
		bt, _ := json.Marshal(filter.ErrResponse{Message: ctx.Err.Error()})
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

func (f responseFilter) newResponse(data interface{}) *client.Response {
	hump := false
	if f.strategy == constant.ResponseStrategyHump {
		hump = true
	}
	r, err := dealResp(data, hump)
	if err != nil {
		return &client.Response{Data: data}
	}

	return &client.Response{Data: r}
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
	var inMap = make(map[interface{}]interface{})
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

	var out = make(map[string]interface{}, len(m))
	for k1, v1 := range m {
		x := humpToUnderline(k1)

		if v1 == nil {
			out[x] = v1
		} else if reflect.TypeOf(v1).Kind() == reflect.Struct {
			out[x] = humpToLine(struct2Map(v1))
		} else if reflect.TypeOf(v1).Kind() == reflect.Slice {
			value := reflect.ValueOf(v1)
			var newTemps = make([]interface{}, 0, value.Len())
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

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}
