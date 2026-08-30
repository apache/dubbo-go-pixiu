/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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

package schema

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuiltinRegistry(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	assert.Len(t, registry.List(), 4)

	for _, kind := range []string{KindResource, KindMethod, KindListener, KindCluster} {
		objectSchema, ok := registry.Lookup(APIVersionV1Alpha1, kind)
		require.True(t, ok, kind)
		assert.Equal(t, kind, objectSchema.Kind)
		assert.NotEmpty(t, objectSchema.RuntimeAdapter)
	}
}

func TestRegistryRegisterField(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	retryPolicy := FieldSchema{
		Type: FieldTypeObject,
		Properties: map[string]*FieldSchema{
			"maxAttempts": {
				Type:     FieldTypeInteger,
				Required: true,
				Minimum:  floatPointer(1),
			},
		},
		Runtime: &RuntimeBinding{Adapter: "example-retry-adapter", Path: "retryPolicy"},
	}
	require.NoError(t, registry.RegisterField(
		APIVersionV1Alpha1,
		KindMethod,
		"spec.extensions.retryPolicy",
		retryPolicy,
	))

	objectSchema, ok := registry.Lookup(APIVersionV1Alpha1, KindMethod)
	require.True(t, ok)
	require.NotNil(t, objectSchema.Fields["extensions"].Properties["retryPolicy"])
	assert.Equal(t, "example-retry-adapter", objectSchema.Fields["extensions"].Properties["retryPolicy"].Runtime.Adapter)

	// Lookup returns a defensive copy; callers cannot mutate the registry.
	delete(objectSchema.Fields["extensions"].Properties, "retryPolicy")
	objectSchema, ok = registry.Lookup(APIVersionV1Alpha1, KindMethod)
	require.True(t, ok)
	assert.Contains(t, objectSchema.Fields["extensions"].Properties, "retryPolicy")

	err = registry.RegisterField(APIVersionV1Alpha1, KindMethod, "extensions.retryPolicy", retryPolicy)
	assert.ErrorIs(t, err, ErrFieldAlreadyRegistered)
}

func TestRegistryRegisterFieldRejectsUnknownParent(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)

	err = registry.RegisterField(APIVersionV1Alpha1, KindMethod, "missing.value", FieldSchema{Type: FieldTypeString})
	require.Error(t, err)
	assert.False(t, errors.Is(err, ErrFieldAlreadyRegistered))
}

func TestRegisteredFieldParticipatesInValidation(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	require.NoError(t, registry.RegisterField(
		APIVersionV1Alpha1,
		KindCluster,
		"extensions.retryPolicy",
		FieldSchema{
			Type: FieldTypeObject,
			Properties: map[string]*FieldSchema{
				"maxAttempts": integerField(true, floatPointer(1), nil),
			},
		},
	))

	object := validClusterObject()
	object.Spec["extensions"] = map[string]any{
		"retryPolicy": map[string]any{"maxAttempts": 0},
	}
	err = registry.Validate(object)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "spec.extensions.retryPolicy.maxAttempts")
}
