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
	"errors"
	"strconv"
	"strings"
	"sync"
	"testing"
)

import (
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
)

type recordingSnapshotCache struct {
	failures  int
	nodes     []string
	versions  []string
	snapshots []*cache.Snapshot
}

func (c *recordingSnapshotCache) SetSnapshot(_ context.Context, node string, snapshot cache.ResourceSnapshot) error {
	c.nodes = append(c.nodes, node)
	c.versions = append(c.versions, snapshot.GetVersion(resource.ExtensionConfigType))
	if typedSnapshot, ok := snapshot.(*cache.Snapshot); ok {
		c.snapshots = append(c.snapshots, typedSnapshot)
	}
	if c.failures > 0 {
		c.failures--
		return errors.New("cache unavailable")
	}
	return nil
}

type snapshotBuildFunc func(version string) (*SnapshotBuildResult, error)

func (f snapshotBuildFunc) Build(version string) (*SnapshotBuildResult, error) {
	return f(version)
}

type mutableResourceLoader struct {
	listeners []config.Listener
	clusters  []config.Cluster
}

func (l *mutableResourceLoader) LoadListeners() ([]config.Listener, error) {
	return l.listeners, nil
}

func (l *mutableResourceLoader) LoadClusters() ([]config.Cluster, error) {
	return l.clusters, nil
}

func TestSnapshotPublisherAdvancesVersionOnlyAfterSuccess(t *testing.T) {
	status := NewStatusStore("node-a")
	cacheWriter := &recordingSnapshotCache{failures: 1}
	publisher := newSnapshotPublisher(
		"node-a",
		NewSnapshotBuilder(fakeResourceLoader{}),
		cacheWriter,
		status,
		0,
	)

	if err := publisher.Publish(context.Background()); err == nil {
		t.Fatal("expected first cache publication to fail")
	}
	failed := status.Snapshot()
	if failed.SnapshotVersion != "" || !strings.Contains(failed.LastError, "cache unavailable") {
		t.Fatalf("failed publication replaced last-good state: %+v", failed)
	}

	if err := publisher.Publish(context.Background()); err != nil {
		t.Fatalf("retry publish: %v", err)
	}
	if err := publisher.Publish(context.Background()); err != nil {
		t.Fatalf("second publish: %v", err)
	}

	wantVersions := []string{"1", "1", "2"}
	if strings.Join(cacheWriter.versions, ",") != strings.Join(wantVersions, ",") {
		t.Fatalf("attempted versions: want %v, got %v", wantVersions, cacheWriter.versions)
	}
	if got := status.Snapshot(); got.SnapshotVersion != "2" || got.LastError != "" {
		t.Fatalf("unexpected last-good state: %+v", got)
	}
}

func TestSnapshotPublisherBuildFailureDoesNotTouchCache(t *testing.T) {
	status := NewStatusStore("node-a")
	cacheWriter := &recordingSnapshotCache{}
	publisher := newSnapshotPublisher("node-a", snapshotBuildFunc(func(string) (*SnapshotBuildResult, error) {
		return nil, errors.New("invalid resource")
	}), cacheWriter, status, 0)

	if err := publisher.Publish(context.Background()); err == nil {
		t.Fatal("expected build failure")
	}
	if len(cacheWriter.versions) != 0 {
		t.Fatalf("cache was called for an invalid candidate: %v", cacheWriter.versions)
	}
	if got := status.Snapshot(); !strings.Contains(got.LastError, "invalid resource") {
		t.Fatalf("build failure was not recorded: %+v", got)
	}
}

func TestSnapshotPublisherPublishesDeletionAsEmptyResource(t *testing.T) {
	loader := &mutableResourceLoader{listeners: []config.Listener{{Name: "to-delete"}}}
	status := NewStatusStore("node-a")
	cacheWriter := &recordingSnapshotCache{}
	publisher := newSnapshotPublisher(
		"node-a",
		NewSnapshotBuilder(loader),
		cacheWriter,
		status,
		0,
	)

	if err := publisher.Publish(context.Background()); err != nil {
		t.Fatalf("publish populated snapshot: %v", err)
	}
	loader.listeners = nil
	if err := publisher.Publish(context.Background()); err != nil {
		t.Fatalf("publish empty snapshot: %v", err)
	}

	if len(cacheWriter.snapshots) != 2 {
		t.Fatalf("snapshot count: want 2, got %d", len(cacheWriter.snapshots))
	}
	latest := &SnapshotBuildResult{Snapshot: cacheWriter.snapshots[1]}
	if got := len(unpackListeners(t, latest).Listeners); got != 0 {
		t.Fatalf("deleted listener remained in latest snapshot: %d resources", got)
	}
	if got := status.Snapshot(); got.SnapshotVersion != "2" || got.ListenerCount != 0 {
		t.Fatalf("status did not record empty snapshot: %+v", got)
	}
}

func TestSnapshotPublisherSerializesConcurrentVersions(t *testing.T) {
	status := NewStatusStore("node-a")
	cacheWriter := &recordingSnapshotCache{}
	publisher := newSnapshotPublisher(
		"node-a",
		NewSnapshotBuilder(fakeResourceLoader{}),
		cacheWriter,
		status,
		0,
	)

	const publishes = 25
	var wg sync.WaitGroup
	for i := 0; i < publishes; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := publisher.Publish(context.Background()); err != nil {
				t.Errorf("Publish: %v", err)
			}
		}()
	}
	wg.Wait()

	if len(cacheWriter.versions) != publishes {
		t.Fatalf("publish count: want %d, got %d", publishes, len(cacheWriter.versions))
	}
	for i, version := range cacheWriter.versions {
		if want := strconv.Itoa(i + 1); version != want {
			t.Fatalf("version at index %d: want %s, got %s", i, want, version)
		}
	}
	if got := status.Snapshot().SnapshotVersion; got != strconv.Itoa(publishes) {
		t.Fatalf("last status version: want %d, got %s", publishes, got)
	}
}

func TestSnapshotPublisherRejectsIncompleteConfiguration(t *testing.T) {
	if err := newSnapshotPublisher("node-a", nil, nil, nil, 0).Publish(context.Background()); err == nil {
		t.Fatal("expected incomplete publisher error")
	}
}
