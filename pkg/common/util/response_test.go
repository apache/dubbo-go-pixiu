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

package util

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestNewDubboMapResponse(t *testing.T) {
	slice := append(make([]string, 0), "read")
	resp := map[string]interface{}{
		"age":   18,
		"iD":    0o001,
		"name":  "tc",
		"time":  nil,
		"hobby": slice,
		"father": map[string]interface{}{
			"name": "bob",
		},
	}
	result, err := dealResp(resp, false)
	if err != nil {
		t.Error("failed to deal resp")
	}
	r := result.(map[string]interface{})
	assert.Equal(t, 18, r["age"])
	assert.Equal(t, 1, r["iD"])
	assert.Equal(t, "tc", r["name"])
	assert.Equal(t, "read", r["hobby"].([]string)[0])
	assert.Equal(t, "bob", r["father"].(map[string]interface{})["name"])
}

func TestNewDubboMapHumpResponse(t *testing.T) {
	slice := append(make([]string, 0), "readBook")
	resp := map[string]interface{}{
		"age":         18,
		"iD":          0o001,
		"name":        "tc",
		"time":        nil,
		"simpleHobby": slice,
		"father": map[string]interface{}{
			"name": "bob",
		},
	}
	result, err := dealResp(resp, true)
	if err != nil {
		t.Error("failed to deal resp")
	}
	r := result.(map[string]interface{})
	assert.Equal(t, 18, r["age"])
	assert.Equal(t, 1, r["i_d"])
	assert.Equal(t, "tc", r["name"])
	assert.Equal(t, "readBook", r["simple_hobby"].([]interface{})[0])
	assert.Equal(t, "bob", r["father"].(map[string]interface{})["name"])
}

func TestNewDubboSliceResponse(t *testing.T) {
	resp := append(make([]map[string]interface{}, 0), map[string]interface{}{
		"age":  18,
		"iD":   0o001,
		"name": "tc",
		"time": nil,
		"father": map[string]interface{}{
			"name": "bob",
		},
	})
	result, err := dealResp(resp, false)
	if err != nil {
		t.Error("failed to deal resp")
	}
	r := result.([]interface{})[0].(map[string]interface{})
	assert.Equal(t, 18, r["age"])
	assert.Equal(t, 1, r["iD"])
	assert.Equal(t, "tc", r["name"])
	assert.Equal(t, "bob", r["father"].(map[string]interface{})["name"])
}
