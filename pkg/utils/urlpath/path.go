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

package urlpath

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"strings"
)

//Split url split to []string by "/"
func Split(path string) []string {
	return strings.Split(strings.TrimLeft(path, constant.PathSlash), constant.PathSlash)
}

//VariableName extract VariableName      (:id, name = id)
func VariableName(key string) string {
	return strings.TrimPrefix(key, constant.PathParamIdentifier)
}

//IsPathVariable return if is a PathVariable     (:id, true)
func IsPathVariable(key string) bool {
	if key == "" {
		return false
	}
	if strings.HasPrefix(key, constant.PathParamIdentifier) {
		return true
	}
	//return key[0] == '{' && key[len(key)-1] == '}'
	return false
}
