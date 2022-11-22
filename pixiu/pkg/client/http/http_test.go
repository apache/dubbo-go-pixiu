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
	"io"
	"net/http"
	"testing"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/mock"
)

func TestMapParams(t *testing.T) {
	hClient := NewHTTPClient()
	r, _ := http.NewRequest("POST", "/mock/test?team=theBoys", bytes.NewReader([]byte(
		"{\"id\":\"12345\",\"age\":\"19\",\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\",\"nickName\":\"trump\"}}")))
	r.Header.Set("Auth", "12345")
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	req := client.NewReq(context.TODO(), r, api)

	val, err := hClient.MapParams(req)
	assert.Nil(t, err)
	p, _ := val.(*requestParams)
	assert.Equal(t, p.Query.Encode(), "team=theBoys")
	assert.Equal(t, p.Header.Get("Auth"), "12345")
	rawBody, err := io.ReadAll(p.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody),
		"{\"id\":\"12345\",\"age\":\"19\",\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\",\"nickName\":\"trump\"}}")

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
	api.IntegrationRequest.HTTPBackendConfig.Schema = "https"
	api.IntegrationRequest.HTTPBackendConfig.Host = "localhost"
	r, _ = http.NewRequest("POST", "/mock/test?team=theBoys", bytes.NewReader([]byte(
		"{\"id\":\"12345\",\"age\":\"19\",\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\",\"nickName\":\"trump\"}}")))
	r.Header.Set("Auth", "12345")
	req = client.NewReq(context.TODO(), r, api)
	val, err = hClient.MapParams(req)
	assert.Nil(t, err)
	p, _ = val.(*requestParams)
	assert.Equal(t, p.Header.Get("Id"), "12345")
	assert.Equal(t, p.Query.Get("auth"), "12345")
	assert.Equal(t, p.Query.Get("team"), "theBoys")
	rawBody, err = io.ReadAll(p.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody),
		"{\"age\":\"19\",\"nickName\":\"trump\",\"testStruct\":{\"name\":\"mock\",\"nickName\":\"trump\",\"test\":\"happy\"}}")
}

func TestParseURL(t *testing.T) {
	hClient := NewHTTPClient()
	requestParams := newRequestParams()
	requestParams.URIParams.Set("id", "12345")
	r, _ := http.NewRequest("POST", "/mock/test/12345", bytes.NewReader([]byte("")))
	api := mock.GetMockAPI(config.MethodGet, "/mock/test/:id")
	api.IntegrationRequest.RequestType = "http"
	api.IntegrationRequest.HTTPBackendConfig.Schema = "http"
	api.IntegrationRequest.HTTPBackendConfig.Host = "abc.com"
	api.IntegrationRequest.HTTPBackendConfig.Path = "/:id"
	req := client.NewReq(context.TODO(), r, api)
	parsedURL, err := hClient.parseURL(req, *requestParams)
	assert.Equal(t, parsedURL, "http://abc.com/12345")
	assert.Nil(t, err)

	requestParams = newRequestParams()
	requestParams.URIParams.Set("id", "12345")
	requestParams.Query.Set("name", "Joe")
	parsedURL, err = hClient.parseURL(req, *requestParams)
	assert.Equal(t, parsedURL, "http://abc.com/12345?name=Joe")
	assert.Nil(t, err)

	requestParams = newRequestParams()
	requestParams.URIParams.Set("id", "12345")
	req.API.HTTPBackendConfig.Path = ""
	parsedURL, err = hClient.parseURL(req, *requestParams)
	assert.Equal(t, parsedURL, "http://abc.com")
	assert.Nil(t, err)

	requestParams = newRequestParams()
	requestParams.URIParams.Set("id", "12345")
	requestParams.Query.Set("name", "Joe")
	parsedURL, err = hClient.parseURL(req, *requestParams)
	assert.Equal(t, parsedURL, "http://abc.com?name=Joe")
	assert.Nil(t, err)
}
