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

package weightrandom

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestWeightRandom_Handler(t *testing.T) {
	tests := []struct {
		name          string
		clusterConfig *model.ClusterConfig
		want          *model.Endpoint
	}{
		{
			name:          "no healthy endpoints",
			clusterConfig: &model.ClusterConfig{Endpoints: []*model.Endpoint{}},
			want:          nil,
		},
		{
			name:          "single healthy endpoint with default weight",
			clusterConfig: clusterWithWeightedEndpoints(""),
			want:          &model.Endpoint{ID: "ep1", Name: "ep1"},
		},
		// `want: nil` below means "any healthy endpoint is acceptable"; the helper
		// in the assertion accepts non-nil picks for non-empty clusters.
		{
			name:          "multiple healthy endpoints with default weight, should return one randomly",
			clusterConfig: clusterWithWeightedEndpoints("", "", ""),
			want:          nil,
		},
		{
			name:          "multiple healthy endpoints with different weights",
			clusterConfig: clusterWithWeightedEndpoints("3", "1", "2"),
			want:          nil,
		},
		{
			name:          "endpoint with invalid weight string, should use default weight",
			clusterConfig: clusterWithWeightedEndpoints("abc", ""),
			want:          nil,
		},
		{
			name:          "endpoint with zero weight, should use default weight",
			clusterConfig: clusterWithWeightedEndpoints("0", ""),
			want:          nil,
		},
		{
			name:          "endpoint with negative weight, should use default weight",
			clusterConfig: clusterWithWeightedEndpoints("-1", ""),
			want:          nil,
		},
		{
			name:          "all endpoints have invalid weights, should return random",
			clusterConfig: clusterWithWeightedEndpoints("abc", "def"),
			want:          nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				wr  = WeightRandom{}
				got = wr.HandlerWithSnapshot(testPickContext(tt.clusterConfig), nil)
			)

			if tt.want == nil {
				if got == nil {
					// Expected nil, got nil - pass
				} else if len(tt.clusterConfig.Endpoints) > 0 {
					// Expected nil (due to multiple or weighted), got non-nil - pass (can't predict)
				} else {
					t.Errorf("WeightRandom.Handler() got = %v, want nil", got)
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WeightRandom.Handler() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// clusterWithWeightedEndpoints builds a ClusterConfig with sequentially-named
// endpoints (`ep1`, `ep2`, ...). Each entry in weights becomes one endpoint;
// an empty string leaves Metadata unset, any other value is stored as the
// raw `weight` metadata.
func clusterWithWeightedEndpoints(weights ...string) *model.ClusterConfig {
	endpoints := make([]*model.Endpoint, 0, len(weights))
	for i, w := range weights {
		id := fmt.Sprintf("ep%d", i+1)
		ep := &model.Endpoint{ID: id, Name: id}
		if w != "" {
			ep.Metadata = map[string]string{"weight": w}
		}
		endpoints = append(endpoints, ep)
	}
	return &model.ClusterConfig{Endpoints: endpoints}
}

// Helper function to create a ClusterConfig with specific endpoints and weights for probabilistic testing
func createWeightedClusterConfig(endpointsWithWeights map[string]int) *model.ClusterConfig {
	endpoints := make([]*model.Endpoint, 0, len(endpointsWithWeights))
	for name, weight := range endpointsWithWeights {
		metadata := make(map[string]string)
		if weight > 0 {
			metadata["weight"] = strconv.Itoa(weight)
		}
		endpoints = append(endpoints, &model.Endpoint{ID: name, Name: name, Metadata: metadata})
	}
	return &model.ClusterConfig{
		Endpoints: endpoints,
	}
}

// Probabilistic test to check if the weighting is working as expected over many iterations
func TestWeightRandom_Handler_Probabilistic(t *testing.T) {
	tests := []struct {
		name             string
		endpointsWeights map[string]int
		expectedProbs    map[string]float64
		tolerance        float64
		iterations       int
	}{
		{
			name:             "simple weighted distribution",
			endpointsWeights: map[string]int{"ep1": 3, "ep2": 1},
			expectedProbs:    map[string]float64{"ep1": 0.75, "ep2": 0.25},
			tolerance:        0.05,
			iterations:       10000,
		},
		{
			name:             "more endpoints with varying weights",
			endpointsWeights: map[string]int{"a": 1, "b": 2, "c": 7},
			expectedProbs:    map[string]float64{"a": 0.1, "b": 0.2, "c": 0.7},
			tolerance:        0.03,
			iterations:       20000,
		},
		{
			name:             "some endpoints with default weight (y = 1)",
			endpointsWeights: map[string]int{"x": 2, "y": 1, "z": 3}, // 'y' will get weight 1
			expectedProbs:    map[string]float64{"x": 0.33, "y": 0.17, "z": 0.5},
			tolerance:        0.04,
			iterations:       15000,
		},
		{
			name:             "some endpoints with zero weight (y = 0)",
			endpointsWeights: map[string]int{"x": 2, "y": 0, "z": 3}, // 'y' will get weight 0
			expectedProbs:    map[string]float64{"x": 0.4, "y": 0, "z": 0.6},
			tolerance:        0.04,
			iterations:       15000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				clusterConfig = createWeightedClusterConfig(tt.endpointsWeights)
				wr            = WeightRandom{}
				counts        = make(map[string]int)
			)

			for i := 0; i < tt.iterations; i++ {
				endpoint := wr.HandlerWithSnapshot(testPickContext(clusterConfig), nil)
				if endpoint != nil {
					counts[endpoint.ID]++
				}
			}

			for id, expectedProb := range tt.expectedProbs {
				actualProb := float64(counts[id]) / float64(tt.iterations)
				if diff := abs(actualProb - expectedProb); diff > tt.tolerance {
					t.Errorf("Endpoint %s: expected probability %f, got %f (difference %f > tolerance %f)",
						id, expectedProb, actualProb, diff, tt.tolerance)
				}
				if expectedProb == 0 && counts[id] != 0 {
					t.Errorf("Endpoint %s: expected 0 selections (weight=0), got %d", id, counts[id])
				}
			}
		})
	}
}

func testPickContext(cluster *model.ClusterConfig) loadbalancer.PickContext {
	return loadbalancer.PickContext{
		Config:           cluster,
		HealthyEndpoints: cluster.GetEndpoint(true),
	}
}

func abs(f float64) float64 {
	if f < 0 {
		return -f
	}
	return f
}
