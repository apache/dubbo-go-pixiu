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
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/router"
)

var mappers = map[string]client.ParamMapper{
	constant.QueryStrings: queryStringsMapper{},
	constant.Headers:      headerMapper{},
	constant.RequestBody:  bodyMapper{},
	constant.RequestURI:   uriMapper{},
}

type requestParams struct {
	Header    http.Header
	Query     url.Values
	Body      io.ReadCloser
	URIParams url.Values
}

func newRequestParams() *requestParams {
	return &requestParams{
		Header:    http.Header{},
		Query:     url.Values{},
		Body:      io.NopCloser(bytes.NewReader([]byte(""))),
		URIParams: url.Values{},
	}
}

type queryStringsMapper struct{}

// nolint
func (qs queryStringsMapper) Map(mp config.MappingParam, c *client.Request, rawTarget interface{}, option client.RequestOption) error {
	target, fromKey, to, toKey, err := mapPrepare(mp, rawTarget)
	if err != nil {
		return err
	}
	queryValues, err := url.ParseQuery(c.IngressRequest.URL.RawQuery)
	if err != nil {
		return errors.Wrap(err, "Error happened when parsing the query paramters")
	}
	rawValue := queryValues.Get(fromKey[0])
	if len(rawValue) == 0 {
		return errors.Errorf("%s in query parameters not found", fromKey[0])
	}
	setTarget(target, to, toKey[0], rawValue)
	return nil
}

type headerMapper struct{}

// nolint
func (hm headerMapper) Map(mp config.MappingParam, c *client.Request, rawTarget interface{}, option client.RequestOption) error {
	target, fromKey, to, toKey, err := mapPrepare(mp, rawTarget)
	if err != nil {
		return err
	}

	rawHeader := c.IngressRequest.Header.Get(fromKey[0])
	if len(rawHeader) == 0 {
		return errors.Errorf("Header %s not found", fromKey[0])
	}
	setTarget(target, to, toKey[0], rawHeader)
	return nil
}

type bodyMapper struct{}

// nolint
func (bm bodyMapper) Map(mp config.MappingParam, c *client.Request, rawTarget interface{}, option client.RequestOption) error {
	// TO-DO: add support for content-type other than application/json
	target, fromKey, to, toKey, err := mapPrepare(mp, rawTarget)
	if err != nil {
		return err
	}

	rawBody, err := io.ReadAll(c.IngressRequest.Body)
	defer func() {
		c.IngressRequest.Body = io.NopCloser(bytes.NewReader(rawBody))
	}()
	if err != nil {
		return err
	}
	mapBody := map[string]interface{}{}
	json.Unmarshal(rawBody, &mapBody)
	val, err := client.GetMapValue(mapBody, fromKey)
	if err != nil {
		return errors.Wrapf(err, "Error when get body value from key %s", fromKey)
	}
	setTarget(target, to, strings.Join(toKey, constant.Dot), val)
	return nil
}

type uriMapper struct{}

// nolint
func (um uriMapper) Map(mp config.MappingParam, c *client.Request, rawTarget interface{}, option client.RequestOption) error {
	target, fromKey, to, toKey, err := mapPrepare(mp, rawTarget)
	if err != nil {
		return err
	}
	if to == constant.RequestURI {
	}
	uriValues := router.GetURIParams(&c.API, *c.IngressRequest.URL)
	if uriValues == nil {
		return errors.New("No URI parameters found")
	}
	setTarget(target, to, strings.Join(toKey, constant.Dot), uriValues.Get(fromKey[0]))
	return nil
}

func validateTarget(target interface{}) (*requestParams, error) {
	val, ok := target.(*requestParams)
	if !ok || val == nil {
		return nil, errors.New("Target params must be a requestParams pointer")
	}
	return val, nil
}

func setTarget(target *requestParams, to string, key string, val interface{}) error {
	valType := reflect.TypeOf(val)
	if (to == constant.Headers || to == constant.QueryStrings) && valType.Kind() != reflect.String {
		return errors.Errorf("%s only accepts string", to)
	}
	switch to {
	case constant.Headers:
		target.Header.Set(key, val.(string))
	case constant.RequestURI:
		target.URIParams.Set(key, val.(string))
	case constant.QueryStrings:
		target.Query.Set(key, val.(string))
	case constant.RequestBody:
		rawBody, err := io.ReadAll(target.Body)
		defer func() {
			target.Body = io.NopCloser(bytes.NewReader(rawBody))
		}()
		if err != nil {
			return errors.New("Raw body read failed")
		}
		mapBody := map[string]interface{}{}
		if len(rawBody) != 0 {
			err = json.Unmarshal(rawBody, &mapBody)
			if err != nil {
				return errors.New("Raw body parse failed")
			}
		}
		setMapWithPath(mapBody, key, val)
		rawBody, err = json.Marshal(mapBody)
		if err != nil {
			return errors.New("Stringify map to body failed")
		}
	default:
		return errors.Errorf("Mapping target to %s does not support", to)
	}
	return nil
}

// setMapWithPath set the value to the target map. If the origin targetMap has
// {"abc": "cde": {"f":1, "g":2}} and the path is abc, value is "123", then the
// targetMap will be updated to {"abc", "123"}
func setMapWithPath(targetMap map[string]interface{}, path string, val interface{}) map[string]interface{} {
	keys := strings.Split(path, constant.Dot)

	_, ok := targetMap[keys[0]]
	if len(keys) == 1 {
		targetMap[keys[0]] = val
		return targetMap
	}
	if !ok && len(keys) != 1 {
		targetMap[keys[0]] = make(map[string]interface{})
		targetMap[keys[0]] = setMapWithPath(targetMap[keys[0]].(map[string]interface{}), strings.Join(keys[1:], constant.Dot), val)
	}
	return targetMap
}

func mapPrepare(mp config.MappingParam, rawTarget interface{}) (target *requestParams, fromKey []string, to string, toKey []string, err error) {
	// ensure the target is a pointer and type is requestParams
	target, err = validateTarget(rawTarget)
	if err != nil {
		return nil, nil, "", nil, err
	}
	// retrieve the mapping values' origin param name
	_, fromKey, err = client.ParseMapSource(mp.Name)
	if err != nil {
		return nil, nil, "", nil, err
	}
	// retrieve the mapping values' target param name and param types(header/uri/query/request body)
	to, toKey, err = client.ParseMapSource(mp.MapTo)
	if err != nil {
		return nil, nil, "", nil, err
	}
	return target, fromKey, to, toKey, nil
}
