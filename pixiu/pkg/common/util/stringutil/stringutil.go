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

package stringutil

import (
	"net"
	"strings"
	"time"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
)

// Split url split to []string by "/"
func Split(path string) []string {
	return strings.Split(strings.TrimLeft(path, constant.PathSlash), constant.PathSlash)
}

// VariableName extract VariableName      (:id, name = id)
func VariableName(key string) string {
	return strings.TrimPrefix(key, constant.PathParamIdentifier)
}

// IsPathVariableOrWildcard return if is a PathVariable     (:id, true)
func IsPathVariableOrWildcard(key string) bool {
	if key == "" {
		return false
	}
	if strings.HasPrefix(key, constant.PathParamIdentifier) {
		return true
	}

	if IsWildcard(key) {
		return true
	}
	//return key[0] == '{' && key[len(key)-1] == '}'
	return false
}

// IsWildcard return if is *
func IsWildcard(key string) bool {
	return key == "*"
}

func IsMatchAll(key string) bool {
	return key == "**"
}

func GetTrieKey(method string, path string) string {
	ret := ""
	//"http://localhost:8882/api/v1/test-dubbo/user?name=tc"
	if strings.Contains(path, constant.ProtocolSlash) {
		path = path[strings.Index(path, constant.ProtocolSlash)+len(constant.ProtocolSlash):]
		path = path[strings.Index(path, constant.PathSlash)+1:]
	}
	if strings.HasPrefix(path, constant.PathSlash) {
		ret = method + path
	} else {
		ret = method + constant.PathSlash + path
	}
	if strings.HasSuffix(ret, constant.PathSlash) {
		ret = ret[0 : len(ret)-1]
	}
	ret = strings.Split(ret, "?")[0]
	return ret
}

func GetIPAndPort(address string) ([]*net.TCPAddr, error) {
	if len(address) <= 0 {
		return nil, errors.Errorf("invalid address, %s", address)
	}
	tcpAddr := make([]*net.TCPAddr, 0)
	for _, addrs := range strings.Split(address, ",") {
		addr, err := net.ResolveTCPAddr("", addrs)
		if err != nil {
			return nil, err
		}
		tcpAddr = append(tcpAddr, addr)
	}

	return tcpAddr, nil
}

func ResolveTimeStr2Time(currentV string, defaultV time.Duration) time.Duration {
	if currentV == "" {
		return defaultV
	} else {
		if duration, err := time.ParseDuration(currentV); err != nil {
			return defaultV
		} else {
			return duration
		}
	}
}
