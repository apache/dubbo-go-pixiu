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
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/mock"
)

func TestQueryStringsMapper(t *testing.T) {
	r, _ := http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:    "queryStrings.id",
			MapTo:   "0",
			MapType: "string",
		},
		{
			Name:    "queryStrings.name",
			MapTo:   "1",
			MapType: "string",
		},
		{
			Name:    "queryStrings.age",
			MapTo:   "jk",
			MapType: "int",
		},
	}

	req := client.NewReq(context.TODO(), r, api)

	params := newDubboTarget(api.IntegrationRequest.MappingParams)
	qs := queryStringsMapper{}
	// Giving valid mapping params
	err := qs.Map(api.IntegrationRequest.MappingParams[0], req, params, nil)
	// it should not return error
	assert.Nil(t, err)
	// it should update the target value in target position from corresponding query value in request.
	assert.Equal(t, params.Values[0], "12345")
	assert.Equal(t, params.Types[0], "string")
	// Giving valid mapping params and same target
	err = qs.Map(api.IntegrationRequest.MappingParams[1], req, params, nil)
	// it should return error when request does not contain the source parameter
	assert.EqualError(t, err, "Query parameter [name] does not exist")
	// Giving invalid mapping params that is not a number and same target
	err = qs.Map(api.IntegrationRequest.MappingParams[2], req, params, nil)
	// it should return error that points out the mapping param
	assert.EqualError(t, err, "Parameter mapping {queryStrings.age jk int} incorrect")

	r, _ = http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	api = mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:    "queryStrings.id",
			MapTo:   "1",
			MapType: "string",
		},
		{
			Name:    "queryStrings.age",
			MapTo:   "0",
			MapType: "int",
		},
	}

	req = client.NewReq(context.TODO(), r, api)
	params = newDubboTarget(api.IntegrationRequest.MappingParams)
	err = qs.Map(api.IntegrationRequest.MappingParams[0], req, params, nil)
	assert.Nil(t, err)
	assert.Equal(t, params.Values[1], "12345")
	assert.Equal(t, params.Types[1], "string")
	assert.Nil(t, params.Values[0])
	err = qs.Map(api.IntegrationRequest.MappingParams[1], req, params, nil)
	assert.Nil(t, err)
	assert.Equal(t, params.Types[0], "int")
	assert.Equal(t, params.Values[0], 19)
}

func TestHeaderMapper(t *testing.T) {
	r, _ := http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	r.Header.Set("Auth", "1234567")
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:    "headers.Auth",
			MapTo:   "0",
			MapType: "string",
		},
	}
	hm := headerMapper{}
	target := newDubboTarget(api.IntegrationRequest.MappingParams)
	req := client.NewReq(context.TODO(), r, api)

	err := hm.Map(api.IntegrationRequest.MappingParams[0], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Values[0], "1234567")
	assert.Equal(t, target.Types[0], "string")

	err = hm.Map(config.MappingParam{Name: "headers.Test", MapTo: "0"}, req, target, nil)
	assert.EqualError(t, err, "Header Test not found")
}

func TestBodyMapper(t *testing.T) {
	r, _ := http.NewRequest("POST", "/mock/test?id=12345&age=19", bytes.NewReader([]byte(`{"sex": "male", "name":{"firstName": "Joe", "lastName": "Biden"}}`)))
	r.Header.Set("Auth", "1234567")
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:    "requestBody.sex",
			MapTo:   "0",
			MapType: "string",
		},
		{
			Name:    "requestBody.name.lastName",
			MapTo:   "1",
			MapType: "string",
		},
		{
			Name:    "requestBody.name",
			MapTo:   "2",
			MapType: "object",
		},
	}
	bm := bodyMapper{}
	target := newDubboTarget(api.IntegrationRequest.MappingParams)
	req := client.NewReq(context.TODO(), r, api)

	err := bm.Map(api.IntegrationRequest.MappingParams[0], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Values[0], "male")
	assert.Equal(t, target.Types[0], "string")

	err = bm.Map(api.IntegrationRequest.MappingParams[1], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Values[1], "Biden")
	assert.Equal(t, target.Types[1], "string")

	err = bm.Map(api.IntegrationRequest.MappingParams[2], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Types[2], "object")
	assert.Equal(t, target.Values[2], map[string]interface{}(map[string]interface{}{
		"firstName": "Joe", "lastName": "Biden",
	}))
}

func TestURIMapper(t *testing.T) {
	r, _ := http.NewRequest("POST", "/mock/12345/joe?age=19", bytes.NewReader([]byte(
		`{"sex": "male", "name":{"firstName": "Joe", "lastName": "Biden"}}`)))
	r.Header.Set("Auth", "1234567")
	api := mock.GetMockAPI(config.MethodGet, "/mock/:id/:name")
	api.IntegrationRequest.MappingParams = []config.MappingParam{
		{
			Name:    "requestBody.sex",
			MapTo:   "0",
			MapType: "string",
		},
		{
			Name:    "requestBody.name.lastName",
			MapTo:   "1",
			MapType: "string",
		},
		{
			Name:    "uri.name",
			MapTo:   "2",
			MapType: "object",
		},
		{
			Name:    "uri.id",
			MapTo:   "3",
			MapType: "string",
		},
	}

	um := uriMapper{}
	target := newDubboTarget(api.IntegrationRequest.MappingParams)
	req := client.NewReq(context.TODO(), r, api)
	err := um.Map(api.IntegrationRequest.MappingParams[3], req, target, nil)
	assert.Nil(t, err)
	err = um.Map(api.IntegrationRequest.MappingParams[2], req, target, nil)
	assert.Nil(t, err)
	assert.Equal(t, target.Values[2], "joe")
	assert.Equal(t, target.Types[2], "object")
	assert.Equal(t, target.Values[3], "12345")
	assert.Equal(t, target.Types[3], "string")
}

func TestValidateTarget(t *testing.T) {
	target := newDubboTarget([]config.MappingParam{
		{
			Name:    "requestBody.sex",
			MapTo:   "0",
			MapType: "string",
		},
	})
	val, err := validateTarget(target)
	assert.Nil(t, err)
	assert.NotNil(t, val)
	_, err = validateTarget(*target)
	assert.EqualError(t, err, "Target params for dubbo backend must be *dubbogoTarget")
	target2 := ""
	_, err = validateTarget(target2)
	assert.EqualError(t, err, "Target params for dubbo backend must be *dubbogoTarget")
}

func TestMapType(t *testing.T) {
	_, err := mapTypes("strings", 123)
	assert.EqualError(t, err, "Invalid parameter type: strings")

	val, err := mapTypes("string", 123)
	assert.Nil(t, err)
	assert.Equal(t, val, "123")
	_, err = mapTypes("string", []int{123, 222})
	assert.EqualError(t, err, "unable to cast []int{123, 222} of type []int to string")

	val, err = mapTypes("int", "123")
	assert.Nil(t, err)
	assert.Equal(t, val, 123)
	val, err = mapTypes("int", 123.6)
	assert.Nil(t, err)
	assert.Equal(t, val, 123)
	_, err = mapTypes("int", "123a")
	assert.EqualError(t, err, "unable to cast \"123a\" of type string to int64")

	val, err = mapTypes("object", map[string]string{"abc": "123"})
	assert.Nil(t, err)
	assert.Equal(t, val, map[string]string{"abc": "123"})
	val, err = mapTypes("object", struct{ Abc string }{Abc: "123"})
	assert.Nil(t, err)
	assert.Equal(t, val, struct{ Abc string }{Abc: "123"})
	val, err = mapTypes("object", 123.6)
	assert.Nil(t, err)
	assert.Equal(t, val, 123.6)
}

func TestNewDubboTarget(t *testing.T) {
	mps := []config.MappingParam{
		{
			Name:  "string1",
			MapTo: "0",
		},
		{
			Name:  "string2",
			MapTo: "opt.values",
		},
	}
	target := newDubboTarget(mps)
	assert.NotNil(t, target)
	assert.Equal(t, len(target.Values), 2)

	mps = []config.MappingParam{
		{
			Name:  "string1",
			MapTo: "opt.interface",
		},
	}
	target = newDubboTarget(mps)
	assert.Nil(t, target)
}

func TestSetCommonTarget(t *testing.T) {
	vals := make([]interface{}, 10)
	types := make([]string, 10)
	target := &dubboTarget{
		Values: vals,
		Types:  types,
	}
	setCommonTarget(target, 1, 123, "int")
	assert.Equal(t, target.Values[1], 123)
	assert.Equal(t, target.Types[1], "int")
	assert.Nil(t, target.Values[0])
	assert.Equal(t, target.Types[0], "")
	setCommonTarget(target, 10, "123", "string")
	assert.Equal(t, target.Values[10], "123")
	assert.Equal(t, target.Types[10], "string")
}

func TestSetGenericTarget(t *testing.T) {
	api := mock.GetMockAPI(config.MethodGet, "/mock/test")
	r, _ := http.NewRequest("GET", "/mock/test?id=12345&age=19", bytes.NewReader([]byte("")))
	req := client.NewReq(context.TODO(), r, api)

	target := &dubboTarget{
		Values: make([]interface{}, 3),
		Types:  make([]string, 3),
	}

	opt := DefaultMapOption[optionKeyValues]
	err := setGenericTarget(req, opt, target, []interface{}{1, "abc", struct{ Name string }{"joe"}}, "int, string, object")
	assert.Nil(t, err)
	assert.Equal(t, target.Values[0], 1)
	assert.Equal(t, target.Values[1], "abc")
	assert.Equal(t, target.Values[2], struct{ Name string }{"joe"})
	assert.Equal(t, target.Types[0], "int")
	assert.Equal(t, target.Types[1], "string")
	assert.Equal(t, target.Types[2], "object")

	opt = DefaultMapOption[optionKeyTypes]
	err = setGenericTarget(req, opt, target, "int, object, object", "")
	assert.Nil(t, err)
	assert.Equal(t, target.Types[0], "int")
	assert.Equal(t, target.Types[1], "object")
	assert.Equal(t, target.Types[2], "object")

	opt = DefaultMapOption[optionKeyInterface]
	err = setGenericTarget(req, opt, target, "testingInterface", "")
	assert.Nil(t, err)
	assert.Equal(t, req.API.IntegrationRequest.Interface, "testingInterface")

	opt = DefaultMapOption[optionKeyApplication]
	err = setGenericTarget(req, opt, target, "testingApplication", "")
	assert.Nil(t, err)
	assert.Equal(t, req.API.IntegrationRequest.ApplicationName, "testingApplication")
}

func TestGetGenericMapTo(t *testing.T) {
	isGeneric, gMapTo := getGenericMapTo("1")
	assert.False(t, isGeneric)
	assert.Equal(t, gMapTo, "")

	isGeneric, gMapTo = getGenericMapTo("opt.interface")
	assert.True(t, isGeneric)
	assert.Equal(t, gMapTo, "interface")

	isGeneric, gMapTo = getGenericMapTo("opt.whatever")
	assert.False(t, isGeneric)
	assert.Equal(t, gMapTo, "")
}
