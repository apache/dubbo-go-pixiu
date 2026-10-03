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

package fingerprint

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestJSONStableIgnoresMapIterationOrder(t *testing.T) {
	a := map[string]any{"b": 2, "a": map[string]any{"d": 4, "c": 3}}
	b := map[string]any{"a": map[string]any{"c": 3, "d": 4}, "b": 2}

	fpA, err := JSONStable(a)
	assert.NoError(t, err)
	fpB, err := JSONStable(b)
	assert.NoError(t, err)
	assert.Equal(t, fpA, fpB)
}

func TestStringsOrderedPreservesSliceSemantics(t *testing.T) {
	assert.NotEqual(t, StringsOrdered([]string{"a", "b"}), StringsOrdered([]string{"b", "a"}))
	assert.Equal(t, StringsSorted([]string{"a", "b"}), StringsSorted([]string{"b", "a"}))
}
