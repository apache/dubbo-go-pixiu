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

package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/mock"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
)

func TestMapParams(t *testing.T) {
	hClient := NewHTTPClient()
	r, _ := http.NewRequest("POST", "/mock/test?team=theBoys", bytes.NewReader([]byte("{\"id\":\"12345\",\"age\":\"19\",\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\",\"nickName\":\"trump\"}}")))
	r.Header.Set("Auth", "12345")
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	req := client.NewReq(context.TODO(), r, api)

	val, err := hClient.MapParams(req)
	assert.Nil(t, err)
	p, _ := val.(*requestParams)
	assert.Equal(t, p.Query.Encode(), "team=theBoys")
	assert.Equal(t, p.Header.Get("Auth"), "12345")
	rawBody, err := ioutil.ReadAll(p.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody), "{\"id\":\"12345\",\"age\":\"19\",\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\",\"nickName\":\"trump\"}}")

	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "queryStrings.team",
			MapTo: "queryStrings.team",
		},
		{
			Name:  "requestBody.id",
			MapTo: "headers.Id",
		},
		{
			Name:  "headers.Auth",
			MapTo: "queryStrings.auth",
		},
		{
			Name:  "requestBody.age",
			MapTo: "requestBody.age",
		},
		{
			Name:  "requestBody.testStruct",
			MapTo: "requestBody.testStruct",
		},
		{
			Name:  "requestBody.testStruct.nickName",
			MapTo: "requestBody.nickName",
		},
	}
	api.IntegrationRequest.HTTPBackendConfig.Scheme = "https"
	api.IntegrationRequest.HTTPBackendConfig.Host = "localhost"
	r, _ = http.NewRequest("POST", "/mock/test?team=theBoys", bytes.NewReader([]byte("{\"id\":\"12345\",\"age\":\"19\",\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\",\"nickName\":\"trump\"}}")))
	r.Header.Set("Auth", "12345")
	req = client.NewReq(context.TODO(), r, api)
	val, err = hClient.MapParams(req)
	assert.Nil(t, err)
	p, _ = val.(*requestParams)
	assert.Equal(t, p.Header.Get("Id"), "12345")
	assert.Equal(t, p.Query.Get("auth"), "12345")
	assert.Equal(t, p.Query.Get("team"), "theBoys")
	rawBody, err = ioutil.ReadAll(p.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody), "{\"age\":\"19\",\"nickName\":\"trump\",\"testStruct\":{\"name\":\"mock\",\"nickName\":\"trump\",\"test\":\"happy\"}}")

	hClient.Call(req)
}
