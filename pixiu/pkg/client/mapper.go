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

import (
	"reflect"
	"regexp"
	"strings"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
)

// ParamMapper defines the interface about how to map the params in the inbound request.
type ParamMapper interface {
	// Map implements how the request parameters map to the target parameters described by config.MappingParam
	Map(config.MappingParam, *Request, interface{}, RequestOption) error
}

// ParseMapSource parses the source parameter config in the mappingParams
// the source parameter in config could be queryStrings.*, headers.*, requestBody.*
func ParseMapSource(source string) (from string, params []string, err error) {
	reg := regexp.MustCompile(`^([uri|queryStrings|headers|requestBody][\w|\d]+)\.([\w|\d|\.|\-]+)$`)
	if !reg.MatchString(source) {
		return "", nil, errors.New("Parameter mapping config incorrect. Please fix it")
	}
	ps := reg.FindStringSubmatch(source)
	return ps[1], strings.Split(ps[2], "."), nil
}

// GetMapValue return the value from map base on the path
func GetMapValue(sourceMap map[string]interface{}, keys []string) (interface{}, error) {
	if len(keys) > 0 && keys[0] == constant.DefaultBodyAll {
		return sourceMap, nil
	}
	_, ok := sourceMap[keys[0]]
	if !ok {
		return nil, errors.Errorf("%s does not exist in request body", keys[0])
	}
	rvalue := reflect.ValueOf(sourceMap[keys[0]])
	if ok && len(keys) == 1 {
		return rvalue.Interface(), nil
	}
	if rvalue.Type().Kind() != reflect.Map && len(keys) == 1 {
		return rvalue.Interface(), nil
	}
	deeperStruct, ok := sourceMap[keys[0]].(map[string]interface{})
	if !ok {
		return nil, errors.Errorf("%s is not a map structure. It contains %v", keys[0], sourceMap[keys[0]])
	}
	return GetMapValue(deeperStruct, keys[1:])
}
