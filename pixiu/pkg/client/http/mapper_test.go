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
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/mock"
)

func TestQueryMapper(t *testing.T) {
	qs := queryStringsMapper{}
	r, _ := http.NewRequest("GET", "/mock/test?id=12345&age=19&name=joe&nickName=trump", bytes.NewReader([]byte("")))
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "queryStrings.id",
			MapTo: "headers.Id",
		},
		{
			Name:  "queryStrings.name",
			MapTo: "queryStrings.name",
		},
		{
			Name:  "queryStrings.age",
			MapTo: "requestBody.age",
		},
		{
			Name:  "queryStrings.nickName",
			MapTo: "requestBody.nickName",
		},
	}
	req := client.NewReq(context.TODO(), r, api)

	target := newRequestParams()
	err := qs.Map(api.IntegrationRequest.MappingParams[0], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Header.Get("Id"), "12345")

	err = qs.Map(api.IntegrationRequest.MappingParams[1], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Query.Get("name"), "joe")

	err = qs.Map(api.IntegrationRequest.MappingParams[2], req, target, nil)
	assert.Nil(t, err)
	err = qs.Map(api.IntegrationRequest.MappingParams[3], req, target, nil)
	assert.Nil(t, err)
	rawBody, _ := io.ReadAll(target.Body)
	assert.Equal(t, string(rawBody), "{\"age\":\"19\",\"nickName\":\"trump\"}")

	err = qs.Map(config.MappingParam{Name: "queryStrings.doesNotExistField", MapTo: "queryStrings.whatever"}, req, target, nil)
	assert.EqualError(t, err, "doesNotExistField in query parameters not found")
}

func TestHeaderMapper(t *testing.T) {
	hm := headerMapper{}
	r, _ := http.NewRequest("GET", "/mock/test?id=12345&age=19&name=joe&nickName=trump", bytes.NewReader([]byte("")))
	r.Header.Set("Auth", "xxxx12345xxx")
	r.Header.Set("Token", "ttttt12345ttt")
	r.Header.Set("Origin-Passcode", "whoseyourdaddy")
	r.Header.Set("Pokemon-Name", "Pika")
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "headers.Auth",
			MapTo: "headers.Auth",
		},
		{
			Name:  "headers.Token",
			MapTo: "headers.Token",
		},
		{
			Name:  "headers.Origin-Passcode",
			MapTo: "queryStrings.originPasscode",
		},
		{
			Name:  "headers.Pokemon-Name",
			MapTo: "requestBody.pokeMonName",
		},
	}
	req := client.NewReq(context.TODO(), r, api)

	target := newRequestParams()
	err := hm.Map(api.IntegrationRequest.MappingParams[0], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Header.Get("Auth"), "xxxx12345xxx")
	err = hm.Map(api.IntegrationRequest.MappingParams[1], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Header.Get("Token"), "ttttt12345ttt")

	err = hm.Map(api.IntegrationRequest.MappingParams[2], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Query.Get("originPasscode"), "whoseyourdaddy")

	err = hm.Map(api.IntegrationRequest.MappingParams[3], req, target, nil)
	assert.Nil(t, err)
	rawBody, err := io.ReadAll(target.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody), "{\"pokeMonName\":\"Pika\"}")

	err = hm.Map(config.MappingParam{Name: "headers.doesNotExistField", MapTo: "headers.whatever"}, req, target, nil)
	assert.EqualError(t, err, "Header doesNotExistField not found")
}

func TestBodyMapper(t *testing.T) {
	bm := bodyMapper{}
	r, _ := http.NewRequest("POST", "/mock/test", bytes.NewReader([]byte("{\"id\":\"12345\",\"age\":\"19\",\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\",\"nickName\":\"trump\"}}")))
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "requestBody.id",
			MapTo: "headers.Id",
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
	req := client.NewReq(context.TODO(), r, api)

	target := newRequestParams()
	err := bm.Map(api.IntegrationRequest.MappingParams[0], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Header.Get("Id"), "12345")

	target = newRequestParams()
	err = bm.Map(api.IntegrationRequest.MappingParams[1], req, target, nil)
	assert.Nil(t, err)

	err = bm.Map(api.IntegrationRequest.MappingParams[2], req, target, nil)
	assert.Nil(t, err)

	err = bm.Map(api.IntegrationRequest.MappingParams[3], req, target, nil)
	assert.Nil(t, err)
	rawBody, err := io.ReadAll(target.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody), "{\"age\":\"19\",\"nickName\":\"trump\",\"testStruct\":{\"name\":\"mock\",\"nickName\":\"trump\",\"test\":\"happy\"}}")
}

func TestURIMap(t *testing.T) {
	um := uriMapper{}
	r, _ := http.NewRequest("POST", "/mock/test/12345", bytes.NewReader([]byte("{\"age\":\"19\",\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\",\"nickName\":\"trump\"}}")))
	api := mock.GetMockAPI(config.MethodGet, "/mock/test/:id")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "uri.id",
			MapTo: "headers.id",
		},
	}
	req := client.NewReq(context.TODO(), r, api)

	target := newRequestParams()
	err := um.Map(api.IntegrationRequest.MappingParams[0], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Header.Get("Id"), "12345")

	target = newRequestParams()
	err = um.Map(config.MappingParam{
		Name:  "uri.id",
		MapTo: "uri.id",
	}, req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.URIParams.Get("id"), "12345")
}

func TestSetTarget(t *testing.T) {
	emptyRequestParams := newRequestParams()
	err := setTarget(emptyRequestParams, constant.Headers, "Auth", "1234565")
	assert.Nil(t, err)
	assert.Equal(t, emptyRequestParams.Header.Get("Auth"), "1234565")
	err = setTarget(emptyRequestParams, constant.Headers, "Auth", 1234565)
	assert.EqualError(t, err, "headers only accepts string")

	err = setTarget(emptyRequestParams, constant.QueryStrings, "id", "123")
	assert.Nil(t, err)
	assert.Equal(t, emptyRequestParams.Query.Get("id"), "123")
	err = setTarget(emptyRequestParams, constant.QueryStrings, "id", 123)
	assert.EqualError(t, err, "queryStrings only accepts string")

	err = setTarget(emptyRequestParams, constant.RequestBody, "testStruct", map[string]interface{}{"test": "happy", "name": "mock"})
	assert.Nil(t, err)
	rawBody, err := io.ReadAll(emptyRequestParams.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody), "{\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\"}}")

	err = setTarget(emptyRequestParams, constant.RequestBody, "testStruct.secondLayer", map[string]interface{}{"test": "happy", "name": "mock"})
	assert.Nil(t, err)
	rawBody, err = io.ReadAll(emptyRequestParams.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody), "{\"testStruct\":{\"secondLayer\":{\"name\":\"mock\",\"test\":\"happy\"}}}")

	err = setTarget(emptyRequestParams, constant.RequestBody, "testStruct.secondLayer.thirdLayer", map[string]interface{}{"test": "happy", "name": "mock"})
	assert.Nil(t, err)
	rawBody, err = io.ReadAll(emptyRequestParams.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody), "{\"testStruct\":{\"secondLayer\":{\"thirdLayer\":{\"name\":\"mock\",\"test\":\"happy\"}}}}")

	err = setTarget(emptyRequestParams, constant.RequestURI, "id", "12345")
	assert.Nil(t, err)
	assert.Equal(t, emptyRequestParams.URIParams.Get("id"), "12345")

	nonEmptyRequestParams := newRequestParams()
	nonEmptyRequestParams.Body = io.NopCloser(bytes.NewReader([]byte("{\"testStruct\":\"abcde\"}")))
	err = setTarget(nonEmptyRequestParams, constant.RequestBody, "testStruct", map[string]interface{}{"test": "happy", "name": "mock"})
	assert.Nil(t, err)
	rawBody, err = io.ReadAll(nonEmptyRequestParams.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody), "{\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\"}}")

	nonEmptyRequestParams = newRequestParams()
	nonEmptyRequestParams.Body = io.NopCloser(bytes.NewReader([]byte("{\"otherStructure\":\"abcde\"}")))
	err = setTarget(nonEmptyRequestParams, constant.RequestBody, "testStruct", map[string]interface{}{"test": "happy", "name": "mock"})
	assert.Nil(t, err)
	rawBody, err = io.ReadAll(nonEmptyRequestParams.Body)
	assert.Nil(t, err)
	assert.Equal(t, string(rawBody), "{\"otherStructure\":\"abcde\",\"testStruct\":{\"name\":\"mock\",\"test\":\"happy\"}}")
}

func TestValidateTarget(t *testing.T) {
	requestP := newRequestParams()
	p, e := validateTarget(requestP)
	assert.NotNil(t, p)
	assert.Nil(t, e)

	requestP = nil
	p, e = validateTarget(requestP)
	assert.Nil(t, p)
	assert.NotNil(t, e)

	requestP2 := []int{}
	p, e = validateTarget(requestP2)
	assert.Nil(t, p)
	assert.NotNil(t, e)

	requestP3 := struct{}{}
	p, e = validateTarget(requestP3)
	assert.Nil(t, p)
	assert.NotNil(t, e)
}
