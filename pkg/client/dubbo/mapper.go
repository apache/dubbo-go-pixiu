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
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

var mappers = map[string]client.ParamMapper{
	constant.QueryStrings: queryStringsMapper{},
	constant.Headers:      headerMapper{},
	constant.RequestBody:  bodyMapper{},
	constant.RequestURI:   uriMapper{},
}

type queryStringsMapper struct{}

// nolint
func (qm queryStringsMapper) Map(mp config.MappingParam, c *client.Request, target interface{}, option client.RequestOption) error {
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

	return setTargetWithOpt(c, option, rv, pos, qValue, c.API.IntegrationRequest.ParamTypes[pos])
}

type headerMapper struct{}

// nolint
func (hm headerMapper) Map(mp config.MappingParam, c *client.Request, target interface{}, option client.RequestOption) error {
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

	return setTargetWithOpt(c, option, rv, pos, header, c.API.IntegrationRequest.ParamTypes[pos])
}

type bodyMapper struct{}

// nolint
func (bm bodyMapper) Map(mp config.MappingParam, c *client.Request, target interface{}, option client.RequestOption) error {
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

	if err := setTargetWithOpt(c, option, rv, pos, val, c.API.IntegrationRequest.ParamTypes[pos]); err != nil {
		return errors.Wrap(err, "set target fail")
	}

	c.IngressRequest.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
	return nil
}

type uriMapper struct{}

// nolint
func (um uriMapper) Map(mp config.MappingParam, c *client.Request, target interface{}, option client.RequestOption) error {
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
	uriValues := router.GetURIParams(&c.API, *c.IngressRequest.URL)

	return setTargetWithOpt(c, option, rv, pos, uriValues.Get(keys[0]), c.API.IntegrationRequest.ParamTypes[pos])
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

func setTargetWithOpt(req *client.Request, option client.RequestOption, rv reflect.Value, pos int, value interface{}, targetType string) error {
	value, err := mapTypes(targetType, value)
	if err != nil {
		return err
	}
	newPos := pos

	if option != nil {
		option.Action(req, value)

		if option.VirtualPos() != 0 {
			newPos = option.VirtualPos()
		}

		if option.Usable() {
			setTarget(rv, newPos, value)
		}

		return nil
	}

	setTarget(rv, newPos, value)

	return nil
}

func setTarget(rv reflect.Value, pos int, value interface{}) {
	if rv.Kind() != reflect.Ptr && rv.Type().Name() != "" && rv.CanAddr() {
		rv = rv.Addr()
	} else {
		rv = rv.Elem()
	}

	tempValue := rv.Interface().([]interface{})

	// for dubbo values split, like RequestOption
	// When config mapTo -1, values is single object, set to 0 position, values is array, will auto split len(values).
	//
	if pos == -1 {
		v, ok := value.([]interface{})
		if ok {
			npos := len(v) - 1
			if len(tempValue) <= npos {
				list := make([]interface{}, npos+1-len(tempValue))
				tempValue = append(tempValue, list...)
			}
			for i := range v {
				s := v[i]
				tempValue[i] = s
			}
			rv.Set(reflect.ValueOf(tempValue))
			return
		}

		pos = 0
	}

	if len(tempValue) <= pos {
		list := make([]interface{}, pos+1-len(tempValue))
		tempValue = append(tempValue, list...)
	}
	tempValue[pos] = value
	rv.Set(reflect.ValueOf(tempValue))
}

func mapTypes(jType string, originVal interface{}) (interface{}, error) {
	targetType, ok := constant.JTypeMapper[jType]
	if !ok {
		return nil, errors.Errorf("Invalid parameter type: %s", jType)
	}
	switch targetType {
	case reflect.TypeOf(""):
		return cast.ToStringE(originVal)
	case reflect.TypeOf(int32(0)):
		return cast.ToInt32E(originVal)
	case reflect.TypeOf(int64(0)):
		return cast.ToInt64E(originVal)
	case reflect.TypeOf(float64(0)):
		return cast.ToFloat64E(originVal)
	case reflect.TypeOf(true):
		return cast.ToBoolE(originVal)
	case reflect.TypeOf(time.Time{}):
		return cast.ToBoolE(originVal)
	default:
		return originVal, nil
	}
}
