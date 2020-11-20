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

package client

// Client represents the interface of http/dubbo clients
type Client interface {
	Init() error
	Close() error
	Call(req *Request) (res interface{}, err error)

	// MapParams mapping param, uri, query, body ...
	MapParams(req *Request) (reqData interface{}, err error)
}

// MapOption option map, key : name, value : option
type MapOption map[string]IOption

var DefaultMapOption = MapOption{
	"types":       &dubboParamTypesOpt{},
	"group":       &dubboGroupOpt{},
	"version":     &dubboVersionOpt{},
	"interface":   &dubboInterfaceOpt{},
	"application": &dubboApplicationOpt{},
	"method":      &dubboMethodOpt{},
}

// IOption option interface.
type IOption interface {
	// Usable if option can use
	Usable() bool
	// SetUsable set usable
	SetUsable(b bool)
	// Action do with val for special
	Action(req *Request, val interface{})
}

// CommonOption
type CommonOption struct {
	usable bool
	IOption
}

// Usable get usable
func (opt *CommonOption) Usable() bool {
	return opt.usable
}

// SetUsable set usable
func (opt *CommonOption) SetUsable(b bool) {
	opt.usable = b
}

type dubboParamTypesOpt struct {
	CommonOption
}

// nolint
func (opt *dubboParamTypesOpt) Action(req *Request, val interface{}) {
	v, ok := val.([]interface{})
	if !ok {
		return
	}

	var pt []string
	for i := range v {
		ptv, ok := v[i].(string)
		if ok {
			pt = append(pt, ptv)
		}
	}

	req.API.IntegrationRequest.DubboBackendConfig.ParamTypes = pt
}

type dubboGroupOpt struct {
	CommonOption
}

// nolint
func (opt *dubboGroupOpt) Action(req *Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.Group = v
}

type dubboVersionOpt struct {
	CommonOption
}

// nolint
func (opt *dubboVersionOpt) Action(req *Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.Version = v
}

type dubboMethodOpt struct {
	CommonOption
}

// nolint
func (opt *dubboMethodOpt) Action(req *Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.Method = v
}

type dubboApplicationOpt struct {
	CommonOption
}

// nolint
func (opt *dubboApplicationOpt) Action(req *Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.ApplicationName = v
}

type dubboInterfaceOpt struct {
	CommonOption
}

// nolint
func (opt *dubboInterfaceOpt) Action(req *Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.Interface = v
}
