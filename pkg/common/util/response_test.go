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
	hessian "github.com/apache/dubbo-go-hessian2"

	"github.com/stretchr/testify/assert"
)

func TestNewDubboMapResponse(t *testing.T) {
	slice := append(make([]string, 0), "read")
	resp := map[string]any{
		"age":   18,
		"iD":    0o001,
		"name":  "tc",
		"time":  nil,
		"hobby": slice,
		"father": map[string]any{
			"name": "bob",
		},
	}
	result, err := dealResp(resp, false)
	if err != nil {
		t.Error("failed to deal resp")
	}
	r := result.(map[string]any)
	assert.Equal(t, 18, r["age"])
	assert.Equal(t, 1, r["iD"])
	assert.Equal(t, "tc", r["name"])
	assert.Equal(t, "read", r["hobby"].([]string)[0])
	assert.Equal(t, "bob", r["father"].(map[string]any)["name"])
}

func TestNewDubboMapHumpResponse(t *testing.T) {
	slice := append(make([]string, 0), "readBook")
	resp := map[string]any{
		"age":         18,
		"iD":          0o001,
		"name":        "tc",
		"time":        nil,
		"simpleHobby": slice,
		"father": map[string]any{
			"name": "bob",
		},
	}
	result, err := dealResp(resp, true)
	if err != nil {
		t.Error("failed to deal resp")
	}
	r := result.(map[string]any)
	assert.Equal(t, 18, r["age"])
	assert.Equal(t, 1, r["i_d"])
	assert.Equal(t, "tc", r["name"])
	assert.Equal(t, "readBook", r["simple_hobby"].([]any)[0])
	assert.Equal(t, "bob", r["father"].(map[string]any)["name"])
}

func TestNewDubboSliceResponse(t *testing.T) {
	resp := append(make([]map[string]any, 0), map[string]any{
		"age":  18,
		"iD":   0o001,
		"name": "tc",
		"time": nil,
		"father": map[string]any{
			"name": "bob",
		},
	})
	result, err := dealResp(resp, false)
	if err != nil {
		t.Error("failed to deal resp")
	}
	r := result.([]any)[0].(map[string]any)
	assert.Equal(t, 18, r["age"])
	assert.Equal(t, 1, r["iD"])
	assert.Equal(t, "tc", r["name"])
	assert.Equal(t, "bob", r["father"].(map[string]any)["name"])
}

func TestNewDubboResponsePreservesNullMapValues(t *testing.T) {
	var nullString *string
	for _, tt := range []struct {
		name     string
		input    any
		want     string
		humpWant string
	}{
		{
			name:     "string map",
			input:    map[string]any{"nullValue": nil, "zero": 0, "false": false, "empty": ""},
			want:     `{"nullValue":null,"zero":0,"false":false,"empty":""}`,
			humpWant: `{"null_value":null,"zero":0,"false":false,"empty":""}`,
		},
		{
			name:     "generic map",
			input:    map[any]any{"nullValue": nil, "zero": 0, "false": false, "empty": "", "class": "java.util.Map"},
			want:     `{"nullValue":null,"zero":0,"false":false,"empty":""}`,
			humpWant: `{"null_value":null,"zero":0,"false":false,"empty":""}`,
		},
		{
			name:     "nested generic map",
			input:    map[any]any{"child": map[any]any{"nullValue": nil}},
			want:     `{"child":{"nullValue":null}}`,
			humpWant: `{"child":{"null_value":null}}`,
		},
		{
			name:     "generic map in slice",
			input:    []any{map[any]any{"nullValue": nil}},
			want:     `[{"nullValue":null}]`,
			humpWant: `[{"null_value":null}]`,
		},
		{
			name:     "nested string map",
			input:    map[any]any{"child": map[string]any{"nullValue": nil}},
			want:     `{"child":{"nullValue":null}}`,
			humpWant: `{"child":{"null_value":null}}`,
		},
		{
			name:     "typed nil value",
			input:    map[any]any{"nullValue": nullString},
			want:     `{"nullValue":null}`,
			humpWant: `{"null_value":null}`,
		},
		{
			name:     "non-string key and class metadata",
			input:    map[any]any{1: nil, "class": nil},
			want:     `{"1":null}`,
			humpWant: `{"1":null}`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Run("without hump conversion", func(t *testing.T) {
				assert.JSONEq(t, tt.want, string(NewDubboResponse(tt.input, false).Data))
			})
			t.Run("with hump conversion", func(t *testing.T) {
				assert.JSONEq(t, tt.humpWant, string(NewDubboResponse(tt.input, true).Data))
			})
		})
	}
}

func TestNewDubboResponsePreservesHessianNullMapValues(t *testing.T) {
	encoder := hessian.NewEncoder()
	if err := encoder.Encode(map[string]any{"nullable": nil, "zero": int32(0)}); err != nil {
		t.Fatal(err)
	}
	decoded, err := hessian.NewDecoder(encoder.Buffer()).Decode()
	if err != nil {
		t.Fatal(err)
	}
	assert.JSONEq(t, `{"nullable":null,"zero":0}`, string(NewDubboResponse(decoded, false).Data))
}
