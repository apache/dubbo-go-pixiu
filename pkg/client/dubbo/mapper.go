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
	"strconv"
)

import (
	"github.com/pkg/errors"
)

import (
	"encoding/json"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"io/ioutil"
	"reflect"
)

type queryStringsMapper struct{}

func (qm queryStringsMapper) Map(mp config.MappingParam, c client.Request, target interface{}) error {
	rv, err := validateTarget(target)
	if err != nil {
		return err
	}
	c.IngressRequest.ParseForm()
	_, key, err := client.ParseMapSource(mp.Name)
	if err != nil {
		return err
	}
	pos, err := strconv.Atoi(mp.MapTo)
	if err != nil {
		return errors.Errorf("Parameter mapping %v incorrect", mp)
	}
	formValue := c.IngressRequest.Form.Get(key[0])
	if len(formValue) == 0 {
		return errors.Errorf("Query parameter %s does not exist", key)
	}

	setTarget(rv, pos, formValue)

	return nil
}

type headerMapper struct{}

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

func (bm bodyMapper) Map(mp config.MappingParam, c client.Request, target interface{}) error {
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

	body, err := c.IngressRequest.GetBody()
	if err != nil {
		return err
	}
	rawBody, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	mapBody := map[string]interface{}{}
	json.Unmarshal(rawBody, &mapBody)
	val, err := getMapValue(mapBody, keys)

	setTarget(rv, pos, val)
	return nil
}

func getMapValue(sourceMap map[string]interface{}, keys []string) (interface{}, error) {
	for i, key := range keys {
		_, ok := sourceMap[key]
		if !ok {
			return nil, errors.Errorf("%s does not exist in request body", key)
		}
		rvalue := reflect.ValueOf(sourceMap[key])
		if rvalue.Type().Kind() != reflect.Map {
			return rvalue.Interface(), nil
		}
		return getMapValue(sourceMap[key].(map[string]interface{}), keys[i+1:])
	}
	return nil, nil
}
