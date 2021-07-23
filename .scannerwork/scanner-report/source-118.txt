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

package client

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestParseMapSource(t *testing.T) {
	from, key, err := ParseMapSource("queryStrings.id")
	assert.Nil(t, err)
	assert.Equal(t, from, "queryStrings")
	assert.Equal(t, key[0], "id")

	from, key, err = ParseMapSource("headers.id")
	assert.Nil(t, err)
	assert.Equal(t, from, "headers")
	assert.Equal(t, key[0], "id")

	from, key, err = ParseMapSource("requestBody.user.id")
	assert.Nil(t, err)
	assert.Equal(t, from, "requestBody")
	assert.Equal(t, key[0], "user")
	assert.Equal(t, key[1], "id")

	_, _, err = ParseMapSource("what.user.id")
	assert.EqualError(t, err, "Parameter mapping config incorrect. Please fix it")
	_, _, err = ParseMapSource("requestBody.*userid")
	assert.EqualError(t, err, "Parameter mapping config incorrect. Please fix it")
}

func TestGetMapValue(t *testing.T) {
	testMap := map[string]interface{}{
		"Test": "test",
		"structure": map[string]interface{}{
			"name": "joe",
			"age":  77,
		},
	}
	val, err := GetMapValue(testMap, []string{"Test"})
	assert.Nil(t, err)
	assert.Equal(t, val, "test")
	val, err = GetMapValue(testMap, []string{"test"})
	assert.Nil(t, val)
	assert.EqualError(t, err, "test does not exist in request body")
	val, err = GetMapValue(testMap, []string{"structure"})
	assert.Nil(t, err)
	assert.Equal(t, val, testMap["structure"])
	val, err = GetMapValue(testMap, []string{"structure", "name"})
	assert.Nil(t, err)
	assert.Equal(t, val, "joe")
	val, err = GetMapValue(map[string]interface{}{}, []string{"structure"})
	assert.Nil(t, val)
	assert.EqualError(t, err, "structure does not exist in request body")

	val, err = GetMapValue(map[string]interface{}{"structure": "test"}, []string{"structure", "name"})
	assert.Nil(t, val)
	assert.EqualError(t, err, "structure is not a map structure. It contains test")
}
