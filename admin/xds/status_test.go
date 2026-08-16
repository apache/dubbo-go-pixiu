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
	"errors"
	"sync"
	"testing"
)

func TestStatusStoreRecordsLastGoodSnapshot(t *testing.T) {
	store := NewStatusStore("node-a")
	store.RecordSuccess("7", 2, 3)

	got := store.Snapshot()
	if got.NodeID != "node-a" || got.SnapshotVersion != "7" {
		t.Fatalf("unexpected identity or version: %+v", got)
	}
	if got.ListenerCount != 2 || got.ClusterCount != 3 {
		t.Fatalf("unexpected resource counts: %+v", got)
	}
	if got.LastUpdatedAt.IsZero() {
		t.Fatal("successful publication did not record its update time")
	}
	if got.LastError != "" {
		t.Fatalf("successful publication retained an error: %q", got.LastError)
	}
}

func TestStatusStoreErrorPreservesLastGoodSnapshot(t *testing.T) {
	store := NewStatusStore("node-a")
	store.RecordSuccess("7", 2, 3)
	want := store.Snapshot()

	store.RecordError(errors.New("candidate is inconsistent"))
	got := store.Snapshot()

	if got.NodeID != want.NodeID || got.SnapshotVersion != want.SnapshotVersion ||
		got.ListenerCount != want.ListenerCount || got.ClusterCount != want.ClusterCount ||
		!got.LastUpdatedAt.Equal(want.LastUpdatedAt) {
		t.Fatalf("failed candidate replaced last-good state: want %+v, got %+v", want, got)
	}
	if got.LastError != "candidate is inconsistent" {
		t.Fatalf("unexpected last error: %q", got.LastError)
	}

	store.RecordSuccess("8", 4, 5)
	if got := store.Snapshot(); got.LastError != "" {
		t.Fatalf("successful publication did not clear error: %q", got.LastError)
	}
}

func TestStatusStoreConcurrentAccess(t *testing.T) {
	store := NewStatusStore("node-a")
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			store.RecordSuccess("version", 1, 1)
		}()
		go func() {
			defer wg.Done()
			_ = store.Snapshot()
		}()
	}
	wg.Wait()

	if got := store.Snapshot(); got.NodeID != "node-a" {
		t.Fatalf("concurrent access corrupted status: %+v", got)
	}
}

func TestStatusStoreReset(t *testing.T) {
	store := NewStatusStore("node-a")
	store.RecordSuccess("7", 2, 3)
	store.RecordError(errors.New("old error"))

	store.Reset("node-b")
	if got := store.Snapshot(); got != (SnapshotStatus{NodeID: "node-b"}) {
		t.Fatalf("reset retained publication state: %+v", got)
	}
}
