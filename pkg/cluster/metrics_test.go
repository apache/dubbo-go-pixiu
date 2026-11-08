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

package cluster

import (
	"context"
	"sync"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/metric/metricdata"
)

// installSnapshotMetricsReader binds the snapshot instruments to a fresh
// ManualReader so a test can inspect the recorded measurements, and resets the
// global provider and instrument state on cleanup so later callers rebind.
//
// It mutates package-global instruments and the process-global MeterProvider,
// so tests that use it must not call t.Parallel(). Additionally, any test in
// this package that constructs a cluster (even without calling this helper)
// will trigger recordSnapshotPublish → initSnapshotMetrics, so tests that
// construct clusters must also not call t.Parallel() to avoid racing on the
// sync.Once and global instrument variables.
func installSnapshotMetricsReader(t *testing.T) *sdkmetric.ManualReader {
	t.Helper()

	reader := sdkmetric.NewManualReader()
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(reader))

	prevProvider := otel.GetMeterProvider()

	otel.SetMeterProvider(provider)
	snapshotMetricsOnce = sync.Once{}
	snapshotPublishTotal, snapshotEndpointCount, snapshotHealthyCount = nil, nil, nil

	t.Cleanup(func() {
		otel.SetMeterProvider(prevProvider)
		snapshotMetricsOnce = sync.Once{}
		snapshotPublishTotal, snapshotEndpointCount, snapshotHealthyCount = nil, nil, nil
	})

	return reader
}

func collectSnapshotMetrics(t *testing.T, reader *sdkmetric.ManualReader) map[string]metricdata.Metrics {
	t.Helper()

	var rm metricdata.ResourceMetrics
	require.NoError(t, reader.Collect(context.Background(), &rm))

	out := make(map[string]metricdata.Metrics)
	for _, sm := range rm.ScopeMetrics {
		for _, m := range sm.Metrics {
			out[m.Name] = m
		}
	}
	return out
}

func sumForCluster(t *testing.T, m metricdata.Metrics, cluster string) int64 {
	t.Helper()
	data, ok := m.Data.(metricdata.Sum[int64])
	require.True(t, ok, "metric %s is not an int64 Sum", m.Name)
	for _, dp := range data.DataPoints {
		if v, ok := dp.Attributes.Value(attribute.Key("cluster")); ok && v.AsString() == cluster {
			return dp.Value
		}
	}
	t.Fatalf("no data point for cluster %q in metric %s", cluster, m.Name)
	return 0
}

func gaugeForCluster(t *testing.T, m metricdata.Metrics, cluster string) int64 {
	t.Helper()
	data, ok := m.Data.(metricdata.Gauge[int64])
	require.True(t, ok, "metric %s is not an int64 Gauge", m.Name)
	for _, dp := range data.DataPoints {
		if v, ok := dp.Attributes.Value(attribute.Key("cluster")); ok && v.AsString() == cluster {
			return dp.Value
		}
	}
	t.Fatalf("no data point for cluster %q in metric %s", cluster, m.Name)
	return 0
}

func TestSnapshotMetricsRecordPublishCountAndSizes(t *testing.T) {
	reader := installSnapshotMetricsReader(t)

	healthy := testEndpoint("ep-1", "127.0.0.1", 18080)
	unhealthy := testEndpoint("ep-2", "127.0.0.2", 18081)
	runtimeCluster := NewCluster(testCluster("snapshot-metrics", healthy, unhealthy))

	// One publish from NewCluster's initial RefreshEndpointsFrom.
	require.True(t, runtimeCluster.UpdateEndpointHealth(unhealthy.ID, unhealthy.Address.GetAddress(), false))

	metrics := collectSnapshotMetrics(t, reader)

	publish, ok := metrics["pixiu_cluster_snapshot_publish_total"]
	require.True(t, ok, "publish total metric missing")
	assert.Equal(t, int64(2), sumForCluster(t, publish, "snapshot-metrics"))

	count, ok := metrics["pixiu_cluster_snapshot_endpoint_count"]
	require.True(t, ok, "endpoint count metric missing")
	assert.Equal(t, int64(2), gaugeForCluster(t, count, "snapshot-metrics"))

	healthyCount, ok := metrics["pixiu_cluster_snapshot_healthy_endpoint_count"]
	require.True(t, ok, "healthy endpoint count metric missing")
	assert.Equal(t, int64(1), gaugeForCluster(t, healthyCount, "snapshot-metrics"))
}

func TestSnapshotMetricsDoNotCountNoOpHealthUpdate(t *testing.T) {
	reader := installSnapshotMetricsReader(t)

	endpoint := testEndpoint("ep-1", "127.0.0.1", 18082)
	runtimeCluster := NewCluster(testCluster("snapshot-metrics-noop", endpoint))

	// Endpoint already healthy; setting it healthy again is a no-op that does
	// not swap the snapshot and so must not increment the publish counter.
	require.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), true))

	metrics := collectSnapshotMetrics(t, reader)
	publish, ok := metrics["pixiu_cluster_snapshot_publish_total"]
	require.True(t, ok, "publish total metric missing")
	assert.Equal(t, int64(1), sumForCluster(t, publish, "snapshot-metrics-noop"))
}

func TestSnapshotMetricsRecordAddressHealthPublish(t *testing.T) {
	reader := installSnapshotMetricsReader(t)

	// Two endpoints share one address, so an address-keyed health flip marks
	// both unhealthy in a single publish.
	first := testEndpoint("ep-1", "127.0.0.1", 18083)
	second := testEndpoint("ep-2", "127.0.0.1", 18083)
	runtimeCluster := NewCluster(testCluster("snapshot-metrics-address", first, second))

	// One publish from NewCluster's initial RefreshEndpointsFrom.
	require.True(t, runtimeCluster.UpdateEndpointAddressHealth(first.Address.GetAddress(), false))

	metrics := collectSnapshotMetrics(t, reader)

	publish, ok := metrics["pixiu_cluster_snapshot_publish_total"]
	require.True(t, ok, "publish total metric missing")
	assert.Equal(t, int64(2), sumForCluster(t, publish, "snapshot-metrics-address"))

	count, ok := metrics["pixiu_cluster_snapshot_endpoint_count"]
	require.True(t, ok, "endpoint count metric missing")
	assert.Equal(t, int64(2), gaugeForCluster(t, count, "snapshot-metrics-address"))

	healthyCount, ok := metrics["pixiu_cluster_snapshot_healthy_endpoint_count"]
	require.True(t, ok, "healthy endpoint count metric missing")
	assert.Equal(t, int64(0), gaugeForCluster(t, healthyCount, "snapshot-metrics-address"))
}
