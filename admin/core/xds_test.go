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

package core

import (
	"context"
	"errors"
	"net"
	"strings"
	"testing"
)

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	adminxds "github.com/apache/dubbo-go-pixiu/admin/xds"
)

type fakeSnapshotPublisher struct {
	calls    int
	failures int
	onCall   func(int)
}

func (p *fakeSnapshotPublisher) Publish(context.Context) error {
	p.calls++
	if p.onCall != nil {
		p.onCall(p.calls)
	}
	if p.failures > 0 {
		p.failures--
		return errors.New("invalid candidate")
	}
	return nil
}

type fakeConfigWatcher struct {
	channels []clientv3.WatchChan
	calls    int
	onCall   func(int)
}

func (w *fakeConfigWatcher) WatchWithPrefix(string) (clientv3.WatchChan, error) {
	if w.calls >= len(w.channels) {
		return nil, errors.New("unexpected watch attempt")
	}
	ch := w.channels[w.calls]
	w.calls++
	if w.onCall != nil {
		w.onCall(w.calls)
	}
	return ch, nil
}

func TestConsumeConfigWatchPublishesEachEventBatch(t *testing.T) {
	ch := make(chan clientv3.WatchResponse, 2)
	ch <- clientv3.WatchResponse{Events: []*clientv3.Event{{}}}
	ch <- clientv3.WatchResponse{Events: []*clientv3.Event{{}, {}}}
	close(ch)

	publisher := &fakeSnapshotPublisher{failures: 1}
	err := consumeConfigWatch(context.Background(), ch, publisher)
	if err == nil || !strings.Contains(err.Error(), "watch channel closed") {
		t.Fatalf("unexpected terminal watch error: %v", err)
	}
	if publisher.calls != 2 {
		t.Fatalf("publish calls: want 2 event batches, got %d", publisher.calls)
	}
}

func TestConsumeConfigWatchIgnoresProgressResponses(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	ch := make(chan clientv3.WatchResponse, 1)
	ch <- clientv3.WatchResponse{}
	cancel()

	publisher := &fakeSnapshotPublisher{}
	if err := consumeConfigWatch(ctx, ch, publisher); err != nil {
		t.Fatalf("consume canceled watch: %v", err)
	}
	if publisher.calls != 0 {
		t.Fatalf("progress response triggered %d publications", publisher.calls)
	}
}

func TestConsumeConfigWatchReturnsWatchError(t *testing.T) {
	ch := make(chan clientv3.WatchResponse, 1)
	ch <- clientv3.WatchResponse{Canceled: true, CompactRevision: 3}
	close(ch)

	err := consumeConfigWatch(context.Background(), ch, &fakeSnapshotPublisher{})
	if err == nil || !strings.Contains(err.Error(), "required revision has been compacted") {
		t.Fatalf("unexpected compacted watch error: %v", err)
	}
}

func TestWatchConfigWithRetryReconnectsAndResyncs(t *testing.T) {
	first := make(chan clientv3.WatchResponse, 1)
	first <- clientv3.WatchResponse{Canceled: true, CompactRevision: 3}
	close(first)
	second := make(chan clientv3.WatchResponse)

	watcher := &fakeConfigWatcher{channels: []clientv3.WatchChan{first, second}}
	ctx, cancel := context.WithCancel(context.Background())
	publisher := &fakeSnapshotPublisher{onCall: func(call int) {
		if call == 2 {
			cancel()
		}
	}}

	watchConfigWithRetry(ctx, watcher, "/pixiu/config/api", publisher, 0)

	if watcher.calls != 2 {
		t.Fatalf("watch attempts: want 2, got %d", watcher.calls)
	}
	if publisher.calls != 2 {
		t.Fatalf("initial and reconnect publications: want 2, got %d", publisher.calls)
	}
}

func TestWatchConfigWithRetryEstablishesWatchBeforeInitialPublish(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sequence := make([]string, 0, 2)
	watcher := &fakeConfigWatcher{
		channels: []clientv3.WatchChan{make(chan clientv3.WatchResponse)},
		onCall: func(int) {
			sequence = append(sequence, "watch")
		},
	}
	publisher := &fakeSnapshotPublisher{onCall: func(int) {
		sequence = append(sequence, "publish")
		cancel()
	}}

	watchConfigWithRetry(ctx, watcher, "/pixiu/config/api", publisher, 0)

	if got := strings.Join(sequence, ","); got != "watch,publish" {
		t.Fatalf("startup sequence: want watch,publish, got %s", got)
	}
}

func TestStartXDSServerRecordsBindFailure(t *testing.T) {
	occupied, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("reserve xDS port: %v", err)
	}
	defer occupied.Close()
	port := uint(occupied.Addr().(*net.TCPAddr).Port)

	previousBootstrap := adminconfig.Bootstrap
	previousStatus := adminxds.DefaultStatusStore
	t.Cleanup(func() {
		adminconfig.Bootstrap = previousBootstrap
		adminxds.DefaultStatusStore = previousStatus
	})
	adminconfig.Bootstrap = &adminconfig.AdminBootstrap{XDS: adminconfig.XDSConfig{
		ListenPort: port,
		NodeID:     "bind-failure",
	}}
	adminxds.DefaultStatusStore = adminxds.NewStatusStore("")

	err = StartxDsServer()
	if err == nil {
		t.Fatal("expected occupied xDS port to fail")
	}
	status := adminxds.DefaultStatusStore.Snapshot()
	if status.Listening || status.ListenError == "" || status.NodeID != "bind-failure" {
		t.Fatalf("bind failure was not reflected in status: %+v", status)
	}
}
