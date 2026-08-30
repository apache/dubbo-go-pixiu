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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuiltinRegistryContainsOnlyAdminRouteBinding(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	require.Len(t, registry.List(), 1)

	objectSchema, ok := registry.Lookup(KindAdminRouteBinding)
	require.True(t, ok)
	assert.Equal(t, KindAdminRouteBinding, objectSchema.Kind)
	assert.Contains(t, objectSchema.Fields, "entry")
	assert.Contains(t, objectSchema.Fields, "target")
	assert.Contains(t, objectSchema.Fields, "params")
	assert.Contains(t, objectSchema.Fields, "publish")
}

func TestRegistryAddsTypedExtensionWithoutExposingMutableSchema(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	require.NoError(t, registry.RegisterField(
		KindAdminRouteBinding,
		"spec.extensions.retryPolicy",
		FieldSchema{
			Type: FieldTypeObject,
			Properties: map[string]*FieldSchema{
				"maxAttempts": {
					Type:     FieldTypeInteger,
					Required: true,
					Minimum:  floatPointer(1),
				},
			},
		},
	))

	objectSchema, ok := registry.Lookup(KindAdminRouteBinding)
	require.True(t, ok)
	delete(objectSchema.Fields["extensions"].Properties, "retryPolicy")
	objectSchema, ok = registry.Lookup(KindAdminRouteBinding)
	require.True(t, ok)
	assert.Contains(t, objectSchema.Fields["extensions"].Properties, "retryPolicy")

	object := validRouteBindingObject()
	object.Spec["extensions"] = map[string]any{
		"retryPolicy": map[string]any{"maxAttempts": 0},
	}
	err = registry.Validate(object)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "spec.extensions.retryPolicy.maxAttempts")
}

func TestRegistryRejectsUnregisteredExtension(t *testing.T) {
	registry, err := NewBuiltinRegistry()
	require.NoError(t, err)
	object := validRouteBindingObject()
	object.Spec["extensions"] = map[string]any{
		"retryPolicy": map[string]any{"maxAttempts": 3},
	}

	err = registry.Validate(object)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "spec.extensions.retryPolicy")
}
