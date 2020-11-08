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
	"regexp"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
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
