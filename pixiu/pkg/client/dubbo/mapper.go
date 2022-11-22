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
	"io"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
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

type dubboTarget struct {
	Values []interface{} // the slice contains the parameters.
	Types  []string      // the slice contains the parameters' types. It should match the values one by one.
}

// pre-allocate proper memory according to the params' usability.
func newDubboTarget(mps []config.MappingParam) *dubboTarget {
	length := 0

	for i := 0; i < len(mps); i++ {
		isGeneric, v := getGenericMapTo(mps[i].MapTo)
		if isGeneric && v != optionKeyValues {
			continue
		}
		length++
	}

	if length > 0 {
		val := make([]interface{}, length)
		target := &dubboTarget{
			Values: val,
			Types:  make([]string, length),
		}
		return target
	}
	return nil
}

type queryStringsMapper struct{}

// nolint
func (qm queryStringsMapper) Map(mp config.MappingParam, c *client.Request, target interface{}, option client.RequestOption) error {
	t, err := validateTarget(target)
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
	if err != nil && option == nil {
		return errors.Errorf("Parameter mapping %v incorrect", mp)
	}
	qValue := queryValues.Get(key[0])
	if len(qValue) == 0 {
		return errors.Errorf("Query parameter %s does not exist", key)
	}

	return setTargetWithOpt(c, option, t, pos, qValue, mp.MapType)
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
	if err != nil && option == nil {
		return errors.Errorf("Parameter mapping %+v incorrect", mp)
	}
	header := c.IngressRequest.Header.Get(key[0])
	if len(header) == 0 {
		return errors.Errorf("Header %s not found", key[0])
	}

	return setTargetWithOpt(c, option, rv, pos, header, mp.MapType)
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
	if err != nil && option == nil {
		return errors.Errorf("Parameter mapping %v incorrect, parameters for Dubbo backend must be mapped to an int to represent position", mp)
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
	val, err := client.GetMapValue(mapBody, keys)

	if err := setTargetWithOpt(c, option, rv, pos, val, mp.MapType); err != nil {
		return errors.Wrap(err, "set target fail")
	}

	c.IngressRequest.Body = io.NopCloser(bytes.NewBuffer(rawBody))
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
	if err != nil && option == nil {
		return errors.Errorf("Parameter mapping %v incorrect", mp)
	}
	uriValues := router.GetURIParams(&c.API, *c.IngressRequest.URL)

	return setTargetWithOpt(c, option, rv, pos, uriValues.Get(keys[0]), mp.MapType)
}

// validateTarget verify if the incoming target for the Map function
// can be processed as expected.
func validateTarget(target interface{}) (*dubboTarget, error) {
	val, ok := target.(*dubboTarget)
	if !ok {
		return nil, errors.New("Target params for dubbo backend must be *dubbogoTarget")
	}
	return val, nil
}

func setTargetWithOpt(req *client.Request, option client.RequestOption,
	target *dubboTarget, pos int, value interface{}, targetType string) error {
	if option != nil {
		return setGenericTarget(req, option, target, value, targetType)
	}
	value, err := mapTypes(targetType, value)
	if err != nil {
		return err
	}
	setCommonTarget(target, pos, value, targetType)
	return nil
}

func setGenericTarget(req *client.Request, option client.RequestOption,
	target *dubboTarget, value interface{}, targetType string) error {
	var err error
	switch option.(type) {
	case *groupOpt, *versionOpt, *interfaceOpt, *applicationOpt, *methodOpt:
		err = option.Action(req, value)
	case *valuesOpt:
		err = option.Action(target, [2]interface{}{value, targetType})
	case *paramTypesOpt:
		err = option.Action(target, value)
	}
	return err
}

func setCommonTarget(target *dubboTarget, pos int, value interface{}, targetType string) {
	// if the mapTo position is greater than the numbers of usable parameters,
	// extend the values and types slices. It changes the address of the the target.
	if cap(target.Values) <= pos {
		list := make([]interface{}, pos+1-len(target.Values))
		typeList := make([]string, pos+1-len(target.Types))
		target.Values = append(target.Values, list...)
		target.Types = append(target.Types, typeList...)
	}
	target.Values[pos] = value
	target.Types[pos] = targetType
}

func mapTypes(jType string, originVal interface{}) (interface{}, error) {
	targetType, ok := constant.JTypeMapper[jType]
	if !ok {
		return nil, errors.Errorf("Invalid parameter type: %s", jType)
	}
	switch targetType {
	case reflect.TypeOf(""):
		return cast.ToStringE(originVal)
	case reflect.TypeOf(int(0)):
		return cast.ToIntE(originVal)
	case reflect.TypeOf(int8(0)):
		return cast.ToInt8E(originVal)
	case reflect.TypeOf(int16(16)):
		return cast.ToInt16E(originVal)
	case reflect.TypeOf(int32(0)):
		return cast.ToInt32E(originVal)
	case reflect.TypeOf(int64(0)):
		return cast.ToInt64E(originVal)
	case reflect.TypeOf(float32(0)):
		return cast.ToFloat32E(originVal)
	case reflect.TypeOf(float64(0)):
		return cast.ToFloat64E(originVal)
	case reflect.TypeOf(true):
		return cast.ToBoolE(originVal)
	case reflect.TypeOf(time.Time{}):
		return cast.ToTimeE(originVal)
	default:
		return originVal, nil
	}
}

// getGenericMapTo will parse the mapTo field, if the mapTo value is
// opt.xxx, the "opt." prefix will identify the param mapTo generic field,
// supporting generic fields: interface, group, application, method, version,
// values, types
func getGenericMapTo(mapTo string) (isGeneric bool, genericField string) {
	fields := strings.Split(mapTo, ".")
	if len(fields) != 2 || fields[0] != "opt" {
		return false, ""
	}
	if _, ok := DefaultMapOption[fields[1]]; !ok {
		return false, ""
	}
	return true, fields[1]
}
