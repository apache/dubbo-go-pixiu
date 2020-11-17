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
	"encoding/json"
	"io/ioutil"
	"net/url"
	"reflect"
	"strconv"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
)

var mappers = map[string]client.ParamMapper{
	constant.QueryStrings: queryStringsMapper{},
	constant.Headers:      headerMapper{},
	constant.RequestBody:  bodyMapper{},
	constant.RequestURI:   uriMapper{},
}

type queryStringsMapper struct{}

// nolint
func (qm queryStringsMapper) Map(mp config.MappingParam, c client.Request, target interface{}) error {
	rv, err := validateTarget(target)
	if err != nil {
		return err
	}
	queryValues, err := url.ParseQuery(c.IngressRequest.URL.RawQuery)
	if err != nil {
		return errors.Wrap(err, "Error happened when parsing the query paramters")
	}
	_, key, err := client.ParseMapSource(mp.Name)
	if err != nil {
		return err
	}
	pos, err := strconv.Atoi(mp.MapTo)
	if err != nil {
		return errors.Errorf("Parameter mapping %v incorrect", mp)
	}
	qValue := queryValues.Get(key[0])
	if len(qValue) == 0 {
		return errors.Errorf("Query parameter %s does not exist", key)
	}

	setTarget(rv, pos, qValue)

	return nil
}

type headerMapper struct{}

// nolint
func (hm headerMapper) Map(mp config.MappingParam, c client.Request, target interface{}) error {
	rv, err := validateTarget(target)
	if err != nil {
		return err
	}
	_, key, err := client.ParseMapSource(mp.Name)
	pos, err := strconv.Atoi(mp.MapTo)
	if err != nil {
		return errors.Errorf("Parameter mapping %v incorrect", mp)
	}
	header := c.IngressRequest.Header.Get(key[0])
	if len(header) == 0 {
		return errors.Errorf("Header %s not found", key[0])
	}
	setTarget(rv, pos, header)
	return nil
}

type bodyMapper struct{}

// nolint
func (bm bodyMapper) Map(mp config.MappingParam, c client.Request, target interface{}) error {
	// TO-DO: add support for content-type other than application/json
	rv, err := validateTarget(target)
	if err != nil {
		return err
	}
	_, keys, err := client.ParseMapSource(mp.Name)
	if err != nil {
		return err
	}
	pos, err := strconv.Atoi(mp.MapTo)
	if err != nil {
		return errors.Errorf("Parameter mapping %v incorrect", mp)
	}

	rawBody, err := ioutil.ReadAll(c.IngressRequest.Body)
	defer func() {
		c.IngressRequest.Body = ioutil.NopCloser(bytes.NewReader(rawBody))
	}()
	if err != nil {
		return err
	}
	mapBody := map[string]interface{}{}
	json.Unmarshal(rawBody, &mapBody)
	val, err := client.GetMapValue(mapBody, keys)

	setTarget(rv, pos, val)
	c.IngressRequest.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
	return nil
}

type uriMapper struct{}

// nolint
func (um uriMapper) Map(mp config.MappingParam, c client.Request, target interface{}) error {
	rv, err := validateTarget(target)
	if err != nil {
		return err
	}
	_, keys, err := client.ParseMapSource(mp.Name)
	if err != nil {
		return err
	}
	pos, err := strconv.Atoi(mp.MapTo)
	if err != nil {
		return errors.Errorf("Parameter mapping %v incorrect", mp)
	}
	uriValues := c.API.GetURIParams(*c.IngressRequest.URL)
	setTarget(rv, pos, uriValues.Get(keys[0]))
	return nil
}

// validateTarget verify if the incoming target for the Map function
// can be processed as expected.
func validateTarget(target interface{}) (reflect.Value, error) {
	rv := reflect.ValueOf(target)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return rv, errors.New("Target params must be a non-nil pointer")
	}
	if _, ok := target.(*[]interface{}); !ok {
		return rv, errors.New("Target params for dubbo backend must be *[]interface{}")
	}
	return rv, nil
}

func setTarget(rv reflect.Value, pos int, value interface{}) {
	if rv.Kind() != reflect.Ptr && rv.Type().Name() != "" && rv.CanAddr() {
		rv = rv.Addr()
	} else {
		rv = rv.Elem()
	}

	tempValue := rv.Interface().([]interface{})
	if len(tempValue) <= pos {
		list := make([]interface{}, pos+1-len(tempValue))
		tempValue = append(tempValue, list...)
	}
	tempValue[pos] = value
	rv.Set(reflect.ValueOf(tempValue))
}
