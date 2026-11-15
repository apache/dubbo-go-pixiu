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

package traffic

import (
	"net/http"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func newTestFactory(t *testing.T, traffics ...*Cluster) *FilterFactory {
	t.Helper()

	hf, err := (&Plugin{}).CreateFilterFactory()
	require.NoError(t, err)
	factory := hf.(*FilterFactory)
	factory.cfg.Traffics = traffics
	require.NoError(t, factory.Apply())
	return factory
}

// decodeCluster runs the filter for a request to /user with a fixed weight
// and returns the cluster the request is routed to.
func decodeCluster(t *testing.T, factory *FilterFactory, weight int, canary string) string {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, "http://localhost/user", nil)
	require.NoError(t, err)
	request.RequestURI = "/user"
	if canary != "" {
		request.Header.Set(string(canaryByHeader), canary)
	}
	ctx := mock.GetMockHTTPContext(request)
	ctx.Route = &model.RouteAction{Cluster: "user"}

	f := &Filter{weight: weight, Rules: factory.rulesMatch(request.RequestURI)}
	f.Decode(ctx)
	return ctx.Route.Cluster
}

func TestTrafficWeightOnlyClusters(t *testing.T) {
	factory := newTestFactory(t,
		&Cluster{Name: "user-v1", Router: "/user", CanaryWeight: 80},
		&Cluster{Name: "user-v2", Router: "/user", CanaryWeight: 20},
	)

	assert.Equal(t, "user-v1", decodeCluster(t, factory, 50, ""))
	assert.Equal(t, "user-v2", decodeCluster(t, factory, 90, ""))
}

func TestTrafficSingleWeightCluster(t *testing.T) {
	factory := newTestFactory(t,
		&Cluster{Name: "user-canary", Router: "/user", CanaryWeight: 20},
	)

	assert.Equal(t, "user-canary", decodeCluster(t, factory, 10, ""))
	// Requests outside the canary weight keep the route's cluster.
	assert.Equal(t, "user", decodeCluster(t, factory, 50, ""))
}

func TestTrafficHeaderAndWeightClusters(t *testing.T) {
	factory := newTestFactory(t,
		&Cluster{Name: "user-v1", Router: "/user", CanaryWeight: 80},
		&Cluster{Name: "user-v2", Router: "/user", CanaryByHeader: "v2", CanaryWeight: 20},
	)

	assert.Equal(t, "user-v1", decodeCluster(t, factory, 50, ""))
	assert.Equal(t, "user-v2", decodeCluster(t, factory, 90, ""))
	assert.Equal(t, "user-v2", decodeCluster(t, factory, 50, "v2"))
	assert.Equal(t, "user-v1", decodeCluster(t, factory, 50, "unknown"))
}
