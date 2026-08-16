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
	"strings"
	"testing"
)

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

type fakeSnapshotPublisher struct {
	calls    int
	failures int
}

func (p *fakeSnapshotPublisher) Publish(context.Context) error {
	p.calls++
	if p.failures > 0 {
		p.failures--
		return errors.New("invalid candidate")
	}
	return nil
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
