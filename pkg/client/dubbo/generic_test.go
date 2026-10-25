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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	cst "github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

func TestNormalizeGenericMode(t *testing.T) {
	tests := []struct {
		desc    string
		mode    string
		want    string
		wantErr bool
	}{
		{desc: "empty selects the map mode", mode: "", want: cst.GenericModeMap},
		{desc: "map mode", mode: "true", want: cst.GenericModeMap},
		{desc: "gson mode", mode: "gson", want: cst.GenericModeGson},
		{desc: "protobuf json mode", mode: "protobuf-json", want: cst.GenericModeProtobufJSON},
		{desc: "bean mode", mode: "bean", want: cst.GenericModeBean},
		{desc: "value is trimmed and lowercased", mode: "  Protobuf-JSON ", want: cst.GenericModeProtobufJSON},
		{desc: "legacy protobuf is rejected", mode: "protobuf", wantErr: true},
		{desc: "unknown mode is rejected", mode: "unknown", wantErr: true},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := NormalizeGenericMode(test.mode)
			if test.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "not supported")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestIsJSONTextGenericMode(t *testing.T) {
	tests := []struct {
		desc string
		mode string
		want bool
	}{
		{desc: "map mode", mode: cst.GenericModeMap, want: false},
		{desc: "gson mode", mode: cst.GenericModeGson, want: true},
		{desc: "protobuf json mode", mode: cst.GenericModeProtobufJSON, want: true},
		{desc: "bean mode", mode: cst.GenericModeBean, want: false},
		{desc: "empty mode", mode: "", want: false},
		{desc: "value is trimmed and lowercased", mode: " GSON ", want: true},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			assert.Equal(t, test.want, IsJSONTextGenericMode(test.mode))
		})
	}
}

func TestEncodeJSONTextArgument(t *testing.T) {
	tests := []struct {
		desc  string
		value any
		want  string
	}{
		{desc: "object becomes JSON text", value: map[string]any{"name": "test"}, want: `{"name":"test"}`},
		{desc: "text passes through", value: `{"name":"test"}`, want: `{"name":"test"}`},
		{desc: "scalar becomes JSON text", value: 42, want: "42"},
		{desc: "nil becomes JSON null", value: nil, want: "null"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got, err := EncodeJSONTextArgument(test.value)
			require.NoError(t, err)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestUnwrapJSONTextResult(t *testing.T) {
	t.Run("JSON text becomes raw bytes", func(t *testing.T) {
		assert.Equal(t, []byte(`{"name":"test"}`), UnwrapJSONTextResult(`{"name":"test"}`))
	})

	t.Run("plain text stays a string", func(t *testing.T) {
		assert.Equal(t, "hello", UnwrapJSONTextResult("hello"))
	})

	t.Run("non string value stays untouched", func(t *testing.T) {
		assert.Equal(t, 42, UnwrapJSONTextResult(42))
		assert.Nil(t, UnwrapJSONTextResult(nil))
	})
}
