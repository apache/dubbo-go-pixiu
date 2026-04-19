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

package rand

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Force Rand to pick index 1 so the test checks which slice length it uses.
func TestRand_TwoHealthyEndpoints_CanSelectSecondEndpoint(t *testing.T) {
	originalRandIntn := randIntn
	randIntn = func(n int) int {
		assert.Equal(t, 2, n)
		return 1
	}
	defer func() {
		randIntn = originalRandIntn
	}()

	cluster := &model.ClusterConfig{
		Endpoints: []*model.Endpoint{
			{ID: "ep-1"},
			{ID: "ep-2"},
		},
	}

	var got *model.Endpoint
	assert.NotPanics(t, func() {
		got = Rand{}.Handler(cluster, nil)
	})
	if assert.NotNil(t, got) {
		assert.Equal(t, "ep-2", got.ID)
	}
}

// With only one healthy endpoint left, Rand must index into the healthy slice, not the full list.
func TestRand_PartiallyUnhealthyEndpoints_DoesNotIndexPastHealthySlice(t *testing.T) {
	originalRandIntn := randIntn
	randIntn = func(n int) int {
		assert.Equal(t, 1, n)
		return 0
	}
	defer func() {
		randIntn = originalRandIntn
	}()

	cluster := &model.ClusterConfig{
		Endpoints: []*model.Endpoint{
			{ID: "ep-1", UnHealthy: true},
			{ID: "ep-2"},
			{ID: "ep-3", UnHealthy: true},
		},
	}

	var got *model.Endpoint
	assert.NotPanics(t, func() {
		got = Rand{}.Handler(cluster, nil)
	})
	if assert.NotNil(t, got) {
		assert.Equal(t, "ep-2", got.ID)
	}
}

// Rand should return nil before calling randIntn when no healthy endpoints remain.
func TestRand_AllEndpointsUnhealthy_ReturnsNilWithoutPanic(t *testing.T) {
	originalRandIntn := randIntn
	randIntn = func(n int) int {
		t.Fatalf("randIntn should not be called, got n=%d", n)
		return 0
	}
	defer func() {
		randIntn = originalRandIntn
	}()

	cluster := &model.ClusterConfig{
		Endpoints: []*model.Endpoint{
			{ID: "ep-1", UnHealthy: true},
			{ID: "ep-2", UnHealthy: true},
		},
	}

	var got *model.Endpoint
	assert.NotPanics(t, func() {
		got = Rand{}.Handler(cluster, nil)
	})
	assert.Nil(t, got)
}
