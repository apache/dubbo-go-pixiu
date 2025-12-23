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

package router

import (
	"net/url"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

func TestGetURIParams(t *testing.T) {
	api := API{
		URLPattern: "/mock/:id/:name",
		Method:     getMockMethod(constant.Get),
	}
	u, _ := url.Parse("https://test.com/mock/12345/Joe")
	values := GetURIParams(&api, *u)
	assert.Equal(t, values.Get("id"), "12345")
	assert.Equal(t, values.Get("name"), "Joe")

	u, _ = url.Parse("https://test.com/Mock/12345/Joe")
	values = GetURIParams(&api, *u)
	assert.Equal(t, values.Get("id"), "12345")
	assert.Equal(t, values.Get("name"), "Joe")

	u, _ = url.Parse("https://test.com/Mock/12345/Joe&jack")
	values = GetURIParams(&api, *u)
	assert.Equal(t, values.Get("id"), "12345")
	assert.Equal(t, values.Get("name"), "Joe&jack")

	u, _ = url.Parse("https://test.com/mock/12345")
	values = GetURIParams(&api, *u)
	assert.Nil(t, values)

	api.URLPattern = "/mock/test"
	u, _ = url.Parse("https://test.com/mock/12345/Joe?status=up")
	values = GetURIParams(&api, *u)
	assert.Nil(t, values)

	api.URLPattern = "/mock/test"
	u, _ = url.Parse("https://test.com/mock/12345/Joe&Jack")
	values = GetURIParams(&api, *u)
	assert.Nil(t, values)
}

func TestIsWildCardBackendPath(t *testing.T) {
	mockAPI := &API{
		URLPattern: "/mock/:id/:name",
		Method:     getMockMethod(constant.Get),
	}
	mockAPI.Path = "/mock/:id"
	assert.True(t, IsWildCardBackendPath(mockAPI))

	mockAPI.Path = "/mock/test"
	assert.False(t, IsWildCardBackendPath(mockAPI))

	mockAPI.Path = ""
	assert.False(t, IsWildCardBackendPath(mockAPI))
}
