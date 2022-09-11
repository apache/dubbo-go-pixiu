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

import (
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
)

// option keys
const (
	optionKeyTypes       = "types"
	optionKeyGroup       = "group"
	optionKeyVersion     = "version"
	optionKeyInterface   = "interface"
	optionKeyApplication = "application"
	optionKeyMethod      = "method"
	optionKeyValues      = "values"
)

// DefaultMapOption default map opt
var DefaultMapOption = client.MapOption{
	optionKeyTypes:       &paramTypesOpt{},
	optionKeyGroup:       &groupOpt{},
	optionKeyVersion:     &versionOpt{},
	optionKeyInterface:   &interfaceOpt{},
	optionKeyApplication: &applicationOpt{},
	optionKeyMethod:      &methodOpt{},
	optionKeyValues:      &valuesOpt{},
}

type groupOpt struct{}

// nolint
func (opt *groupOpt) Action(target, val interface{}) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("Group value is not string")
	}
	r, ok := target.(*client.Request)
	if !ok {
		return errors.New("Target is not *client.Request in value options")
	}
	r.API.IntegrationRequest.DubboBackendConfig.Group = v
	return nil
}

type versionOpt struct{}

// nolint
func (opt *versionOpt) Action(target, val interface{}) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("Version value is not string")
	}
	r, ok := target.(*client.Request)
	if !ok {
		return errors.New("Target is not *client.Request in value options")
	}
	r.API.IntegrationRequest.DubboBackendConfig.Version = v
	return nil
}

type methodOpt struct{}

// nolint
func (opt *methodOpt) Action(target, val interface{}) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("Method value is not string")
	}
	r, ok := target.(*client.Request)
	if !ok {
		return errors.New("Target is not *client.Request in value options")
	}
	r.API.IntegrationRequest.DubboBackendConfig.Method = v
	return nil
}

type applicationOpt struct{}

// nolint
func (opt *applicationOpt) Action(target, val interface{}) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("Application value is not string")
	}
	r, ok := target.(*client.Request)
	if !ok {
		return errors.New("Target is not *client.Request in value options")
	}
	r.API.IntegrationRequest.DubboBackendConfig.ApplicationName = v
	return nil
}

type interfaceOpt struct{}

// nolint
func (opt *interfaceOpt) Action(target, val interface{}) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("Interface value is not string")
	}
	r, ok := target.(*client.Request)
	if !ok {
		return errors.New("Target is not *client.Request in value options")
	}
	r.API.IntegrationRequest.DubboBackendConfig.Interface = v
	return nil
}

type valuesOpt struct{}

// Action of valuesOpt retrieve value from [2]interface{} then assign to target, which the first element is the
// parameter values to dubbo generic call. the second element is string, which is the types
// for the generic call, it could be empty or types sep from ','. If empty, it should retrieve types from
// another generic option - types.
func (opt *valuesOpt) Action(target, val interface{}) error {
	dubboTarget, ok := target.(*dubboTarget)
	if !ok {
		return errors.New("Target is not dubboTarget in value options")
	}
	v, ok := val.([2]interface{})
	if !ok {
		return errors.New("The value must be [2]interface{}")
	}
	var toVals []interface{}
	toTypes := []string{}

	if t, tok := v[1].(string); tok && len(t) != 0 {
		toTypes = strings.Split(t, ",")
	}
	if val, vok := v[0].([]interface{}); vok {
		toVals = val
	} else {
		toVals = []interface{}{v[0]}
	}
	if !(len(toTypes) != 0 && len(toTypes) == len(toVals)) {
		dubboTarget.Types = toTypes
		dubboTarget.Values = toVals
		return nil
	}
	for i := range toVals {
		trimType := strings.TrimSpace(toTypes[i])
		if _, ok = constant.JTypeMapper[trimType]; ok {
			toTypes[i] = trimType
		} else {
			return errors.Errorf("Types invalid %s", trimType)
		}
		var err error
		toVals[i], err = mapTypes(toTypes[i], toVals[i])
		if err != nil {
			return errors.WithStack(err)
		}
	}
	dubboTarget.Types = toTypes
	dubboTarget.Values = toVals
	return nil
}

type paramTypesOpt struct{}

// Action for paramTypesOpt override the other param types mapping/config.
// The val must be string(e.g. "int, object"), and will then assign to the target.(dubboTarget).Types
func (opt *paramTypesOpt) Action(target, val interface{}) error {
	v, ok := val.(string)
	if !ok {
		return errors.New("The val type must be string")
	}
	types := strings.Split(v, ",")
	dubboTarget, ok := target.(*dubboTarget)
	if !ok {
		return errors.New("Target is not dubboTarget in target parameter")
	}
	for i := range types {
		trimType := strings.TrimSpace(types[i])
		if len(trimType) == 0 {
			continue
		}
		if _, ok = constant.JTypeMapper[trimType]; !ok {
			return errors.Errorf("Types invalid %s", trimType)
		}
		types[i] = trimType
	}
	dubboTarget.Types = types
	return nil
}
