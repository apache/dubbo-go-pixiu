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
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
)

// defaultMappingParams default http to dubbo config
var defaultMappingParams = []config.MappingParam{
	{
		Name:  "requestBody.values",
		MapTo: "opt.values",
	}, {
		Name:  "requestBody.types",
		MapTo: "opt.types",
	}, {
		Name:  "uri.application",
		MapTo: "opt.application",
	}, {
		Name:  "uri.interface",
		MapTo: "opt.interface",
	}, {
		Name:  "queryStrings.method",
		MapTo: "opt.method",
	}, {
		Name:  "queryStrings.group",
		MapTo: "opt.group",
	}, {
		Name:  "queryStrings.version",
		MapTo: "opt.version",
	},
}

// getMappingParams If configuration items are set in the api_config configuration file, then the configuration in api_config will be used first,
// if not, the default will be used, and a total rule for obtaining the original data will be generated at the end
func getMappingParams(mps []config.MappingParam) []config.MappingParam {
	newMappingParam := make([]config.MappingParam, len(mps))
	copy(newMappingParam, mps)
	for _, defaultMappingParam := range defaultMappingParams {
		equal := false
		for _, customMappingParam := range mps {
			if customMappingParam.Name == defaultMappingParam.Name {
				equal = true
			}
		}
		if equal {
			continue
		}
		newMappingParam = append(newMappingParam, defaultMappingParam)
	}
	return newMappingParam
}
