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

package roundrobin

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// RoundRobin should return nil instead of indexing an empty healthy slice.
func TestRoundRobin_AllEndpointsUnhealthy_ReturnsNilWithoutPanic(t *testing.T) {
	cluster := &model.ClusterConfig{
		Endpoints: []*model.Endpoint{
			{ID: "ep-1", UnHealthy: true},
			{ID: "ep-2", UnHealthy: true},
		},
	}

	var got *model.Endpoint
	assert.NotPanics(t, func() {
		got = RoundRobin{}.Handler(cluster, nil)
	})
	assert.Nil(t, got)
}

func TestRoundRobin_RepeatedPicksFollowHealthyOrder(t *testing.T) {
	cluster := &model.ClusterConfig{
		Endpoints: []*model.Endpoint{
			{ID: "ep-1"},
			{ID: "ep-2"},
			{ID: "ep-3"},
		},
	}

	rr := RoundRobin{}
	got := []string{
		rr.Handler(cluster, nil).ID,
		rr.Handler(cluster, nil).ID,
		rr.Handler(cluster, nil).ID,
		rr.Handler(cluster, nil).ID,
	}

	assert.Equal(t, []string{"ep-1", "ep-2", "ep-3", "ep-1"}, got)
}

func TestRoundRobin_MixedHealthyEndpointsHonorNonZeroCursor(t *testing.T) {
	cluster := &model.ClusterConfig{
		Endpoints: []*model.Endpoint{
			{ID: "ep-1", UnHealthy: true},
			{ID: "ep-2"},
			{ID: "ep-3", UnHealthy: true},
			{ID: "ep-4"},
		},
		PrePickEndpointIndex: 3,
	}

	rr := RoundRobin{}
	first := rr.Handler(cluster, nil)
	second := rr.Handler(cluster, nil)

	if assert.NotNil(t, first) {
		assert.Equal(t, "ep-4", first.ID)
	}
	if assert.NotNil(t, second) {
		assert.Equal(t, "ep-2", second.ID)
	}
}
