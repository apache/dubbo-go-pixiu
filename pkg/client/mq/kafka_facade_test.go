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

package mq

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/IBM/sarama"
)

// fakeConsumerGroup is a sarama.ConsumerGroup implementation that does not
// require a real Kafka broker. Its Consume blocks until ctx is canceled, then
// returns ctx.Err(), mirroring sarama's behavior on context cancellation.
type fakeConsumerGroup struct {
	mu       sync.Mutex
	consumes int
	closes   int
}

func (f *fakeConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	f.mu.Lock()
	f.consumes++
	f.mu.Unlock()
	<-ctx.Done()
	return ctx.Err()
}

func (f *fakeConsumerGroup) Errors() <-chan error { return nil }

func (f *fakeConsumerGroup) Close() error {
	f.mu.Lock()
	f.closes++
	f.mu.Unlock()
	return nil
}

func (f *fakeConsumerGroup) Pause(map[string][]int32)  {}
func (f *fakeConsumerGroup) Resume(map[string][]int32) {}
func (f *fakeConsumerGroup) PauseAll()                 {}
func (f *fakeConsumerGroup) ResumeAll()                {}

// newTestFacade builds a KafkaConsumerFacade backed by a fakeConsumerGroup via
// the injectable factory, so tests never touch a real Kafka broker.
func newTestFacade(t *testing.T) (*KafkaConsumerFacade, *fakeConsumerGroup) {
	t.Helper()
	fake := &fakeConsumerGroup{}
	facade, err := newKafkaConsumerFacade(KafkaConsumerConfig{}, "test-group", func(_ KafkaConsumerConfig, _ string) (sarama.ConsumerGroup, error) {
		return fake, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if facade.consumerManager == nil {
		t.Fatal("consumerManager should be initialized, not nil")
	}
	return facade, fake
}

// TestNewKafkaConsumerFacadeConstructorInitializesMap verifies that the
// constructor initializes consumerManager, preventing nil map panics during
// Subscribe. It uses the injectable factory so it runs without a Kafka broker.
func TestNewKafkaConsumerFacadeConstructorInitializesMap(t *testing.T) {
	facade, _ := newTestFacade(t)

	facade.mu.RLock()
	_, exists := facade.consumerManager["unused"]
	facade.mu.RUnlock()
	if exists {
		t.Fatal("consumerManager should start empty")
	}
}

// TestKafkaConsumerFacadeConsumerManagerInitialized verifies that writing to
// consumerManager through Subscribe does not panic on a nil map.
func TestKafkaConsumerFacadeConsumerManagerInitialized(t *testing.T) {
	facade, _ := newTestFacade(t)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Subscribe panicked with nil map: %v", r)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := facade.Subscribe(ctx, WithTopics([]string{"topic"}), WithConsumerGroup("group"), WithCheckUrl("http://127.0.0.1:0/health")); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}
	// Stop promptly so the test does not leak goroutines.
	facade.Stop()
}

// TestKafkaConsumerFacadeStopReturns is the regression test for the P0 bug
// where Stop() blocked forever because consumeLoop never called wg.Done(). The
// fake Consume returns when ctx is canceled, so Stop() must complete promptly.
func TestKafkaConsumerFacadeStopReturns(t *testing.T) {
	facade, fake := newTestFacade(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := facade.Subscribe(ctx, WithTopics([]string{"topic"}), WithConsumerGroup("group"), WithCheckUrl("http://127.0.0.1:0/health")); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	done := make(chan struct{})
	go func() {
		facade.Stop()
		close(done)
	}()

	select {
	case <-done:
		// Stop() returned: the consume loop and health check both decremented
		// the WaitGroup on every exit path, as required.
	case <-time.After(3 * time.Second):
		t.Fatal("Stop() did not return within 3s; consume loop likely never called wg.Done()")
	}

	fake.mu.Lock()
	closes := fake.closes
	fake.mu.Unlock()
	if closes != 1 {
		t.Fatalf("expected consumer group to be closed once, got %d", closes)
	}
}

// TestKafkaConsumerFacadeStopCancelsConsumeLoop verifies the P0 fix where the
// cancel func stored in consumerManager actually reaches the consume loop:
// after Stop() the fake's Consume must have returned due to ctx cancellation.
func TestKafkaConsumerFacadeStopCancelsConsumeLoop(t *testing.T) {
	facade, fake := newTestFacade(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := facade.Subscribe(ctx, WithTopics([]string{"topic"}), WithConsumerGroup("group"), WithCheckUrl("http://127.0.0.1:0/health")); err != nil {
		t.Fatalf("Subscribe failed: %v", err)
	}

	facade.Stop()

	fake.mu.Lock()
	consumes := fake.consumes
	fake.mu.Unlock()
	if consumes == 0 {
		t.Fatal("expected Consume to have been invoked at least once")
	}
}

// TestKafkaConsumerFacadeConcurrentAccess verifies that concurrent access to
// consumerManager is safe under the mutex. It uses a bounded wait so a stray
// panic cannot hang the test suite indefinitely.
func TestKafkaConsumerFacadeConcurrentAccess(t *testing.T) {
	facade, _ := newTestFacade(t)

	done := make(chan struct{})

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			facade.mu.Lock()
			facade.consumerManager[GetConsumerManagerKey([]string{"topic"}, "group")] = func() {}
			facade.mu.Unlock()
		}
		done <- struct{}{}
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			facade.mu.RLock()
			_ = len(facade.consumerManager)
			facade.mu.RUnlock()
		}
		done <- struct{}{}
	}()

	for i := 0; i < 2; i++ {
		select {
		case <-done:
		case <-time.After(3 * time.Second):
			t.Fatal("concurrent access test timed out")
		}
	}
}

// TestGetConsumerManagerKey verifies the key generation function is
// deterministic and distinguishes topics and groups.
func TestGetConsumerManagerKey(t *testing.T) {
	topics := []string{"topic1", "topic2"}
	group := "test-group"

	key := GetConsumerManagerKey(topics, group)

	if key != GetConsumerManagerKey(topics, group) {
		t.Error("Key should be deterministic")
	}

	if key == GetConsumerManagerKey([]string{"topic3", "topic4"}, group) {
		t.Error("Different topics should produce different keys")
	}

	if key == GetConsumerManagerKey(topics, "other-group") {
		t.Error("Different groups should produce different keys")
	}
}
