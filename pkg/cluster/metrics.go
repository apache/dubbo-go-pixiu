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
)

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// Endpoint snapshot publication metrics. Instruments are bound lazily to the
// global MeterProvider on first publish; when metrics are disabled the provider
// is a no-op, so recording stays a cheap no-op and never blocks publication.
var (
	snapshotMetricsOnce sync.Once

	snapshotPublishTotal  metric.Int64Counter
	snapshotEndpointCount metric.Int64Gauge
	snapshotHealthyCount  metric.Int64Gauge
)

func initSnapshotMetrics() {
	snapshotMetricsOnce.Do(func() {
		meter := otel.GetMeterProvider().Meter("pixiu")

		var err error
		if snapshotPublishTotal, err = meter.Int64Counter("pixiu_cluster_snapshot_publish_total",
			metric.WithDescription("Total number of cluster endpoint snapshots successfully published.")); err != nil {
			logger.Errorf("[dubbo-go-pixiu] register pixiu_cluster_snapshot_publish_total failed: %v", err)
		}
		if snapshotEndpointCount, err = meter.Int64Gauge("pixiu_cluster_snapshot_endpoint_count",
			metric.WithDescription("Total endpoints in the latest published cluster snapshot.")); err != nil {
			logger.Errorf("[dubbo-go-pixiu] register pixiu_cluster_snapshot_endpoint_count failed: %v", err)
		}
		if snapshotHealthyCount, err = meter.Int64Gauge("pixiu_cluster_snapshot_healthy_endpoint_count",
			metric.WithDescription("Healthy endpoints in the latest published cluster snapshot.")); err != nil {
			logger.Errorf("[dubbo-go-pixiu] register pixiu_cluster_snapshot_healthy_endpoint_count failed: %v", err)
		}
	})
}

// recordSnapshotPublish reports a successful snapshot publication for a cluster.
// Call exactly once per successful CompareAndSwap of the published snapshot. The
// only label is the cluster name, which is bounded by configuration rather than
// request traffic, so cardinality stays low. The health-update callers invoke
// this while holding the cluster's healthMu, so the body must stay
// allocation-light and must not block.
func recordSnapshotPublish(clusterName string, snapshot *EndpointSnapshot) {
	initSnapshotMetrics()

	attrs := metric.WithAttributes(attribute.String("cluster", clusterName))
	ctx := context.Background()

	if snapshotPublishTotal != nil {
		snapshotPublishTotal.Add(ctx, 1, attrs)
	}
	if snapshotEndpointCount != nil {
		snapshotEndpointCount.Record(ctx, int64(snapshot.EndpointCount()), attrs)
	}
	if snapshotHealthyCount != nil {
		snapshotHealthyCount.Record(ctx, int64(snapshot.HealthyEndpointCount()), attrs)
	}
}
