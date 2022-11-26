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
	"bytes"
	"context"
	"net/http"
	"regexp"
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

func TestReg(t *testing.T) {
	reg := regexp.MustCompile(`^(query|uri)`)
	b := reg.MatchString("Stringquery")
	assert.False(t, b)
	b = reg.MatchString("queryString")
	assert.True(t, b)
	b = reg.MatchString("queryuri")
	assert.True(t, b)

	reg = regexp.MustCompile(`^(query|uri)(String|Int|Long|Double|Time)`)
	b = reg.MatchString("Stringquery")
	assert.False(t, b)
	b = reg.MatchString("queryString")
	assert.True(t, b)
	b = reg.MatchString("queryuri")
	assert.False(t, b)
	b = reg.MatchString("queryuriInt")
	assert.False(t, b)

	reg = regexp.MustCompile(`^(query|uri)(String|Int|Long|Double|Time)\.([\w|\d]+)`)
	b = reg.MatchString("Stringquery")
	assert.False(t, b)
	b = reg.MatchString("queryString")
	assert.False(t, b)
	b = reg.MatchString("queryString.name")
	assert.True(t, b)
	b = reg.MatchString("queryuriInt.age")
	assert.False(t, b)

	reg = regexp.MustCompile(`^([query|uri][\w|\d]+)\.([\w|\d]+)$`)
	b = reg.MatchString("queryStrings.name")
	assert.True(t, b)
	b = reg.MatchString("uriInt.")
	assert.False(t, b)
	b = reg.MatchString("queryStrings")
	assert.False(t, b)

	param := reg.FindStringSubmatch("queryString.name")
	for _, p := range param {
		t.Log(p)
	}
}

func TestClose(t *testing.T) {
	client := SingletonDubboClient()
	client.GenericServicePool["key1"] = nil
	client.GenericServicePool["key2"] = nil
	client.GenericServicePool["key3"] = nil
	client.GenericServicePool["key4"] = nil
	assert.Equal(t, 4, len(client.GenericServicePool))
	client.Close()
	assert.Equal(t, 0, len(client.GenericServicePool))
}

func TestMappingParams(t *testing.T) {
	dClient := NewDubboClient()
	r, _ := http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:    "queryStrings.id",
			MapTo:   "0",
			MapType: "string",
		},
		{
			Name:    "queryStrings.age",
			MapTo:   "1",
			MapType: "int",
		},
	}
	req := client.NewReq(context.TODO(), r, api)
	params, err := dClient.MapParams(req)
	assert.Nil(t, err)
	assert.Equal(t, params.(*dubboTarget).Values[0], "12345")
	assert.Equal(t, params.(*dubboTarget).Values[1], int(19))

	r, _ = http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	api = mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:    "queryStrings.id",
			MapTo:   "0",
			MapType: "string",
		},
		{
			Name:    "queryStrings.age",
			MapTo:   "1",
			MapType: "int",
		},
		{
			Name:    "headers.Auth",
			MapTo:   "2",
			MapType: "string",
		},
	}
	r.Header.Set("Auth", "1234567")
	req = client.NewReq(context.TODO(), r, api)
	params, err = dClient.MapParams(req)
	assert.Nil(t, err)
	assert.Equal(t, params.(*dubboTarget).Values[0], "12345")
	assert.Equal(t, params.(*dubboTarget).Values[1], int(19))
	assert.Equal(t, params.(*dubboTarget).Values[2], "1234567")

	r, _ = http.NewRequest("POST", "/mock/test?id=12345&age=19", bytes.NewReader([]byte(`{"sex": "male", "name":{"firstName": "Joe", "lastName": "Biden"}}`)))
	api = mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:    "queryStrings.id",
			MapTo:   "0",
			MapType: "string",
		},
		{
			Name:    "queryStrings.age",
			MapTo:   "1",
			MapType: "int",
		},
		{
			Name:    "headers.Auth",
			MapTo:   "2",
			MapType: "string",
		},
		{
			Name:    "requestBody.sex",
			MapTo:   "3",
			MapType: "string",
		},
		{
			Name:    "requestBody.name.firstName",
			MapTo:   "4",
			MapType: "java.lang.String",
		},
	}
	r.Header.Set("Auth", "1234567")
	req = client.NewReq(context.TODO(), r, api)
	params, err = dClient.MapParams(req)
	assert.Nil(t, err)
	assert.Equal(t, params.(*dubboTarget).Values[0], "12345")
	assert.Equal(t, params.(*dubboTarget).Values[1], int(19))
	assert.Equal(t, params.(*dubboTarget).Values[2], "1234567")
	assert.Equal(t, params.(*dubboTarget).Values[3], "male")
	assert.Equal(t, params.(*dubboTarget).Values[4], "Joe")

	r, _ = http.NewRequest("POST", "/mock/test?id=12345&age=19", bytes.NewReader([]byte(`{"sex": "male", "name":{"firstName": "Joe", "lastName": "Biden"}}`)))
	api = mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "queryStrings.id",
			MapTo: "opt.method",
		},
		{
			Name:  "queryStrings.age",
			MapTo: "opt.application",
		},
		{
			Name:  "headers.Auth",
			MapTo: "opt.group",
		},
		{
			Name:  "requestBody.sex",
			MapTo: "opt.values",
		},
		{
			Name:  "requestBody.name.firstName",
			MapTo: "opt.interface",
		},
	}
	r.Header.Set("Auth", "1234567")
	req = client.NewReq(context.TODO(), r, api)
	_, err = dClient.MapParams(req)
	assert.Nil(t, err)
	assert.Equal(t, req.API.Method.IntegrationRequest.DubboBackendConfig.ApplicationName, "19")
	assert.Equal(t, req.API.Method.IntegrationRequest.DubboBackendConfig.Group, "1234567")
	assert.Equal(t, req.API.Method.IntegrationRequest.DubboBackendConfig.Interface, "Joe")
	assert.Equal(t, req.API.Method.IntegrationRequest.DubboBackendConfig.Method, "12345")
}

func TestBuildOption(t *testing.T) {
	mp := config.MappingParam{
		Name:    "queryStrings.id",
		MapTo:   "0",
		MapType: "string",
	}
	option := buildOption(mp)
	assert.Nil(t, option)

	mp = config.MappingParam{
		Name:    "queryStrings.id",
		MapTo:   "opt.whatsoever",
		MapType: "",
	}
	option = buildOption(mp)
	assert.Nil(t, option)

	mp = config.MappingParam{
		Name:    "queryStrings.id",
		MapTo:   "opt.interface",
		MapType: "",
	}
	option = buildOption(mp)
	assert.NotNil(t, option)
	_, ok := option.(*interfaceOpt)
	assert.True(t, ok)
}

func TestApply(t *testing.T) {
	dClient := NewDubboClient()
	err := dClient.Apply()
	assert.Nil(t, err)
}
