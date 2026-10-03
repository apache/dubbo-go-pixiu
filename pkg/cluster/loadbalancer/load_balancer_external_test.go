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

package loadbalancer_test

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/roundrobin"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// embeddingPlugin is an out-of-tree balancer that embeds an in-tree balancer.
// Method promotion gives it SnapshotOptIn, so structural interface satisfaction
// alone would wrongly grant it the snapshot fast paths. The runtime trust check
// keys on the concrete type's package, which is this external test package, so
// the plugin must still be treated as untrusted.
type embeddingPlugin struct {
	roundrobin.RoundRobin
	seenHealthyEndpoints []*model.Endpoint
}

func (p *embeddingPlugin) HandlerWithSnapshot(c loadbalancer.PickContext, _ model.LbPolicy) *model.Endpoint {
	p.seenHealthyEndpoints = c.HealthyEndpoints
	if len(c.HealthyEndpoints) == 0 {
		return nil
	}
	c.HealthyEndpoints[0].Metadata["weight"] = "mutated"
	return c.HealthyEndpoints[0]
}

// TestEmbeddingInTreeBalancerCannotOptIntoFastPaths is the trust-boundary
// regression for the embedding bypass: an external type that promotes
// SnapshotOptIn from an embedded in-tree balancer must not receive the
// healthy-only or zero-copy fast paths.
func TestEmbeddingInTreeBalancerCannotOptIntoFastPaths(t *testing.T) {
	plugin := &embeddingPlugin{}

	// HealthyOnly fast path must be denied: external embedder still gets the
	// full snapshot.
	assert.True(t, loadbalancer.NeedsAllEndpoints(plugin),
		"embedding an in-tree balancer must not promote the healthy-only fast path")

	healthy := &model.Endpoint{ID: "healthy", Metadata: map[string]string{"weight": "1"}}
	cluster := &model.ClusterConfig{
		Name:      "embedding-trust-boundary",
		Endpoints: []*model.Endpoint{healthy},
	}

	got := loadbalancer.PickEndpoint(plugin, loadbalancer.PickContext{
		AllEndpoints:     []*model.Endpoint{healthy},
		Config:           cluster,
		HealthyEndpoints: []*model.Endpoint{healthy},
	}, nil)

	// ZeroCopy fast path must be denied: the plugin mutates what it receives,
	// and that mutation must not escape to the snapshot-owned endpoint.
	if assert.NotNil(t, got) {
		assert.NotSame(t, healthy, got, "embedder must receive a defensive copy")
	}
	if assert.Len(t, plugin.seenHealthyEndpoints, 1) {
		assert.NotSame(t, healthy, plugin.seenHealthyEndpoints[0])
	}
	assert.Equal(t, map[string]string{"weight": "1"}, healthy.Metadata,
		"embedder mutation must not escape to the snapshot-owned endpoint")
}
