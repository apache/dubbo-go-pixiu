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

func TestQueryStringsMapper(t *testing.T) {
	r, _ := http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "queryStrings.id",
			MapTo: "0",
		},
		{
			Name:  "queryStrings.name",
			MapTo: "1",
		},
		{
			Name:  "queryStrings.age",
			MapTo: "jk",
		},
	}
	req := client.NewReq(context.TODO(), r, api)

	var params []interface{}
	qs := queryStringsMapper{}
	err := qs.Map(api.IntegrationRequest.MappingParams[0], req, &params, nil)
	assert.Nil(t, err)
	assert.Equal(t, params[0], "12345")
	err = qs.Map(api.IntegrationRequest.MappingParams[1], req, &params, nil)
	assert.EqualError(t, err, "Query parameter [name] does not exist")
	err = qs.Map(api.IntegrationRequest.MappingParams[2], req, &params, nil)
	assert.EqualError(t, err, "Parameter mapping {queryStrings.age jk} incorrect")

	r, _ = http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	api = mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "queryStrings.id",
			MapTo: "1",
		},
		{
			Name:  "queryStrings.age",
			MapTo: "0",
		},
	}
	req = client.NewReq(context.TODO(), r, api)
	params = []interface{}{}
	err = qs.Map(api.IntegrationRequest.MappingParams[0], req, &params, nil)
	assert.Nil(t, err)
	assert.Equal(t, params[1], "12345")
	assert.Nil(t, params[0])
	err = qs.Map(api.IntegrationRequest.MappingParams[1], req, &params, nil)
	assert.Nil(t, err)
	assert.Equal(t, params[1], "12345")
	assert.Equal(t, params[0], "19")
}

func TestHeaderMapper(t *testing.T) {
	r, _ := http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	r.Header.Set("Auth", "1234567")
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "headers.Auth",
			MapTo: "0",
		},
	}
	hm := headerMapper{}
	target := []interface{}{}
	req := client.NewReq(context.TODO(), r, api)

	err := hm.Map(api.IntegrationRequest.MappingParams[0], req, &target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target[0], "1234567")

	err = hm.Map(config.MappingParam{Name: "headers.Test", MapTo: "0"}, req, &target, nil)
	assert.EqualError(t, err, "Header Test not found")
}

func TestBodyMapper(t *testing.T) {
	r, _ := http.NewRequest("POST", "/mock/test?id=12345&age=19", bytes.NewReader([]byte(`{"sex": "male", "name":{"firstName": "Joe", "lastName": "Biden"}}`)))
	r.Header.Set("Auth", "1234567")
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "requestBody.sex",
			MapTo: "0",
		},
		{
			Name:  "requestBody.name.lastName",
			MapTo: "1",
		},
		{
			Name:  "requestBody.name",
			MapTo: "2",
		},
	}
	bm := bodyMapper{}
	target := []interface{}{}
	req := client.NewReq(context.TODO(), r, api)

	err := bm.Map(api.IntegrationRequest.MappingParams[0], req, &target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target[0], "male")

	err = bm.Map(api.IntegrationRequest.MappingParams[1], req, &target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target[1], "Biden")

	err = bm.Map(api.IntegrationRequest.MappingParams[2], req, &target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target[2], map[string]interface{}(map[string]interface{}{"firstName": "Joe", "lastName": "Biden"}))
}

func TestURIMapper(t *testing.T) {
	r, _ := http.NewRequest("POST", "/mock/12345/joe&age=19", bytes.NewReader([]byte(`{"sex": "male", "name":{"firstName": "Joe", "lastName": "Biden"}}`)))
	r.Header.Set("Auth", "1234567")
	api := mock.GetMockAPI(config.MethodGet, "/mock/:id/:name")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:  "requestBody.sex",
			MapTo: "0",
		},
		{
			Name:  "requestBody.name.lastName",
			MapTo: "1",
		},
		{
			Name:  "uri.name",
			MapTo: "2",
		},
		{
			Name:  "uri.id",
			MapTo: "3",
		},
	}
	um := uriMapper{}
	target := []interface{}{}
	req := client.NewReq(context.TODO(), r, api)
	err := um.Map(api.IntegrationRequest.MappingParams[3], req, &target, nil)
	assert.Nil(t, err)
	err = um.Map(api.IntegrationRequest.MappingParams[2], req, &target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target[2], "joe")
	assert.Equal(t, target[3], "12345")
}

func TestValidateTarget(t *testing.T) {
	target := []interface{}{}
	val, err := validateTarget(&target)
	assert.Nil(t, err)
	assert.NotNil(t, val)
	_, err = validateTarget(target)
	assert.EqualError(t, err, "Target params must be a non-nil pointer")
	target2 := ""
	_, err = validateTarget(&target2)
	assert.EqualError(t, err, "Target params for dubbo backend must be *[]interface{}")
}
