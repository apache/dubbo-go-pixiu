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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMapTypesStringToInt(t *testing.T) {
	val, err := MapTypes("int", "19")
	require.NoError(t, err)
	assert.Equal(t, 19, val)
}

func TestMapTypesEmptyTypePassThrough(t *testing.T) {
	input := map[string]any{"name": "pixiu"}

	val, err := MapTypes("", input)
	require.NoError(t, err)
	assert.Equal(t, input, val)
}

func TestMapTypesWrapperFQNToPrimitive(t *testing.T) {
	val, err := MapTypes("java.lang.Integer", "19")
	require.NoError(t, err)
	assert.Equal(t, 19, val)
}

func TestMapTypesStringToBool(t *testing.T) {
	val, err := MapTypes("boolean", "true")
	require.NoError(t, err)
	assert.Equal(t, true, val)
}

func TestCoerceDirectInvokeValueWrapperFQN(t *testing.T) {
	val, err := CoerceDirectInvokeValue("java.lang.Integer", "42")
	require.NoError(t, err)
	assert.Equal(t, 42, val)
}

func TestCoerceDirectInvokeValueWrapperFQNArray(t *testing.T) {
	val, err := CoerceDirectInvokeValue("java.lang.Integer[]", []any{"1", "2"})
	require.NoError(t, err)
	assert.Equal(t, []any{1, 2}, val)
}

func TestNormalizeReferenceProtocol(t *testing.T) {
	assert.Equal(t, "tri", NormalizeReferenceProtocol("triple"))
	assert.Equal(t, "tri", NormalizeReferenceProtocol("  TRI  "))
}

func TestDirectURLProtocol(t *testing.T) {
	val, err := DirectURLProtocol("tri://127.0.0.1:50051")
	require.NoError(t, err)
	assert.Equal(t, "tri", val)
}

func TestInferJavaClassNames(t *testing.T) {
	val := InferJavaClassNames([]any{"name", int64(1), true})
	assert.Equal(t, []string{"java.lang.String", "java.lang.Long", "java.lang.Boolean"}, val)
}
