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

package dubbo

import "github.com/dubbogo/dubbo-go-proxy/pkg/client"

// option keys
const (
	optionKeyTypes       = "types"
	optionKeyGroup       = "group"
	optionKeyVersion     = "version"
	optionKeyInterface   = "interface"
	optionKeyApplication = "application"
	optionKeyMethod      = "method"
)

// DefaultMapOption default map opt
var DefaultMapOption = client.MapOption{
	optionKeyTypes:       &paramTypesOpt{},
	optionKeyGroup:       &groupOpt{},
	optionKeyVersion:     &versionOpt{},
	optionKeyInterface:   &interfaceOpt{},
	optionKeyApplication: &applicationOpt{},
	optionKeyMethod:      &methodOpt{},
}

type paramTypesOpt struct {
	client.CommonOption
}

// nolint
func (opt *paramTypesOpt) Action(req *client.Request, val interface{}) {
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

type groupOpt struct {
	client.CommonOption
}

// nolint
func (opt *groupOpt) Action(req *client.Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.Group = v
}

type versionOpt struct {
	client.CommonOption
}

// nolint
func (opt *versionOpt) Action(req *client.Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.Version = v
}

type methodOpt struct {
	client.CommonOption
}

// nolint
func (opt *methodOpt) Action(req *client.Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.Method = v
}

type applicationOpt struct {
	client.CommonOption
}

// nolint
func (opt *applicationOpt) Action(req *client.Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.ApplicationName = v
}

type interfaceOpt struct {
	client.CommonOption
}

// nolint
func (opt *interfaceOpt) Action(req *client.Request, val interface{}) {
	v, ok := val.(string)
	if !ok {
		return
	}

	req.API.IntegrationRequest.DubboBackendConfig.Interface = v
}
