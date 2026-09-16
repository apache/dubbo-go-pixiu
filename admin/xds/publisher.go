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

package xds

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

import (
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
)

type snapshotBuild interface {
	Build(version string) (*SnapshotBuildResult, error)
}

type snapshotCacheWriter interface {
	SetSnapshot(ctx context.Context, node string, snapshot cache.ResourceSnapshot) error
}

// SnapshotPublisher serializes candidate builds and cache publication. Its
// version advances only when SetSnapshot accepts the candidate, so a failed
// build or publication cannot replace the last-good snapshot metadata.
type SnapshotPublisher struct {
	mu          sync.Mutex
	nodeID      string
	lastVersion uint64
	builder     snapshotBuild
	cache       snapshotCacheWriter
	status      *StatusStore
}

func NewSnapshotPublisher(nodeID string, builder *SnapshotBuilder, snapshotCache cache.SnapshotCache, status *StatusStore) *SnapshotPublisher {
	// Seed versions from wall-clock nanoseconds so an Admin restart does not
	// normally reuse a version a still-reconnecting client has already ACKed.
	// The mutex-protected counter then guarantees monotonic versions for this
	// publisher even when multiple reloads happen within one clock tick.
	return newSnapshotPublisher(nodeID, builder, snapshotCache, status, uint64(time.Now().UnixNano()))
}

func newSnapshotPublisher(nodeID string, builder snapshotBuild, cache snapshotCacheWriter, status *StatusStore, initialVersion uint64) *SnapshotPublisher {
	return &SnapshotPublisher{
		nodeID:      nodeID,
		lastVersion: initialVersion,
		builder:     builder,
		cache:       cache,
		status:      status,
	}
}

// Publish builds and installs the next candidate snapshot. Concurrent calls
// are serialized to preserve source order and monotonically increasing
// versions.
func (p *SnapshotPublisher) Publish(ctx context.Context) error {
	if p == nil || p.builder == nil || p.cache == nil || p.status == nil {
		return fmt.Errorf("xDS snapshot publisher is not fully configured")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	nextVersion := p.lastVersion + 1
	version := strconv.FormatUint(nextVersion, 10)
	result, err := p.builder.Build(version)
	if err != nil {
		return p.recordError(fmt.Errorf("build xDS snapshot version %s: %w", version, err))
	}
	if result == nil || result.Snapshot == nil {
		return p.recordError(fmt.Errorf("build xDS snapshot version %s: builder returned a nil snapshot", version))
	}
	if result.Version != version {
		return p.recordError(fmt.Errorf("build xDS snapshot version %s: builder returned version %q", version, result.Version))
	}
	if err := p.cache.SetSnapshot(ctx, p.nodeID, result.Snapshot); err != nil {
		return p.recordError(fmt.Errorf("publish xDS snapshot version %s for node %q: %w", version, p.nodeID, err))
	}

	p.lastVersion = nextVersion
	p.status.RecordSuccess(result.Version, result.ListenerCount, result.ClusterCount)
	return nil
}

func (p *SnapshotPublisher) recordError(err error) error {
	p.status.RecordError(err)
	return err
}
