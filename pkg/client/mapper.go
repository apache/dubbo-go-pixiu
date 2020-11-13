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
	"regexp"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
)

// ParamMapper defines the interface about how to map the params in the inbound request.
type ParamMapper interface {
	Map(config.MappingParam, Request, interface{}) error
}

// ParseMapSource parses the source parameter config in the mappingParams
// the source parameter in config could be queryStrings.*, headers.*, requestBody.*
func ParseMapSource(source string) (from string, params []string, err error) {
	reg := regexp.MustCompile(`^([queryStrings|headers|requestBody][\w|\d]+)\.([\w|\d|\.]+)$`)
	if !reg.MatchString(source) {
		return "", nil, errors.New("Parameter mapping config incorrect. Please fix it")
	}
	ps := reg.FindStringSubmatch(source)
	return ps[1], strings.Split(ps[2], "."), nil
}
