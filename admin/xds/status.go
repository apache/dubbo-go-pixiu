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
	"sync"
	"time"
)

// SnapshotStatus is the last known Admin xDS snapshot publication state.
// LastUpdatedAt describes the last successful publication. A later error does
// not change the last-good version, counts, or timestamp.
type SnapshotStatus struct {
	NodeID          string    `json:"node_id"`
	SnapshotVersion string    `json:"snapshot_version"`
	ListenerCount   int       `json:"listener_count"`
	ClusterCount    int       `json:"cluster_count"`
	LastUpdatedAt   time.Time `json:"last_updated_at"`
	LastError       string    `json:"last_error"`
}

// StatusStore provides a race-safe snapshot of xDS publication state for the
// server, tests, and the future read-only diagnostics endpoint.
type StatusStore struct {
	mu     sync.RWMutex
	status SnapshotStatus
}

// DefaultStatusStore is shared by the Admin xDS server and diagnostics API.
var DefaultStatusStore = NewStatusStore("")

func NewStatusStore(nodeID string) *StatusStore {
	return &StatusStore{status: SnapshotStatus{NodeID: nodeID}}
}

// Reset clears publication state when the xDS server starts with a new
// effective configuration.
func (s *StatusStore) Reset(nodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status = SnapshotStatus{NodeID: nodeID}
}

// RecordSuccess records an accepted and published snapshot and clears any
// error from a previous failed candidate.
func (s *StatusStore) RecordSuccess(version string, listenerCount, clusterCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.SnapshotVersion = version
	s.status.ListenerCount = listenerCount
	s.status.ClusterCount = clusterCount
	s.status.LastUpdatedAt = time.Now().UTC()
	s.status.LastError = ""
}

// RecordError records a failed candidate without replacing the last-good
// snapshot metadata.
func (s *StatusStore) RecordError(err error) {
	if err == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.status.LastError = err.Error()
}

// Snapshot returns a value copy that callers may read without holding the
// store lock.
func (s *StatusStore) Snapshot() SnapshotStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.status
}
