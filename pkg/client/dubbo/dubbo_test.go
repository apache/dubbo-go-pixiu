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
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/mock"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
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
			Name:  "queryStrings.id",
			MapTo: "0",
		},
		{
			Name:  "queryStrings.age",
			MapTo: "1",
		},
	}
	req := client.NewReq(context.TODO(), r, api)
	params, err := dClient.MapParams(req)
	assert.Nil(t, err)
	assert.Equal(t, params.([]interface{})[0], "12345")
	assert.Equal(t, params.([]interface{})[1], "19")

	r, _ = http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	api = mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "queryStrings.id",
			MapTo: "0",
		},
		{
			Name:  "queryStrings.age",
			MapTo: "1",
		},
		{
			Name:  "headers.Auth",
			MapTo: "2",
		},
	}
	r.Header.Set("Auth", "1234567")
	req = client.NewReq(context.TODO(), r, api)
	params, err = dClient.MapParams(req)
	assert.Nil(t, err)
	assert.Equal(t, params.([]interface{})[0], "12345")
	assert.Equal(t, params.([]interface{})[1], "19")
	assert.Equal(t, params.([]interface{})[2], "1234567")

	r, _ = http.NewRequest("POST", "/mock/test?id=12345&age=19", bytes.NewReader([]byte(`{"sex": "male", "name":{"firstName": "Joe", "lastName": "Biden"}}`)))
	api = mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "queryStrings.id",
			MapTo: "0",
		},
		{
			Name:  "queryStrings.age",
			MapTo: "1",
		},
		{
			Name:  "headers.Auth",
			MapTo: "2",
		},
		{
			Name:  "requestBody.sex",
			MapTo: "3",
		},
		{
			Name:  "requestBody.name.firstName",
			MapTo: "4",
		},
	}
	r.Header.Set("Auth", "1234567")
	req = client.NewReq(context.TODO(), r, api)
	params, err = dClient.MapParams(req)
	assert.Nil(t, err)
	assert.Equal(t, params.([]interface{})[0], "12345")
	assert.Equal(t, params.([]interface{})[1], "19")
	assert.Equal(t, params.([]interface{})[2], "1234567")
	assert.Equal(t, params.([]interface{})[3], "male")
	assert.Equal(t, params.([]interface{})[4], "Joe")
}
