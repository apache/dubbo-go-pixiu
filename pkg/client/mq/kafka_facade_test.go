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
	"errors"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/IBM/sarama"

	"github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/require"
)

// fakeConsumerGroup is a sarama.ConsumerGroup implementation that does not
// require a real Kafka broker. Its Consume blocks until ctx is canceled, then
// returns ctx.Err(), mirroring sarama's behavior on context cancellation.
type fakeConsumerGroup struct {
	mu       sync.Mutex
	consumes int
	closes   int
	started  chan context.Context
	consume  func(context.Context) error
	onClose  func() error
}

func (f *fakeConsumerGroup) Consume(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) error {
	f.mu.Lock()
	f.consumes++
	f.mu.Unlock()
	if f.started != nil {
		f.started <- ctx
	}
	if f.consume != nil {
		return f.consume(ctx)
	}
	<-ctx.Done()
	return ctx.Err()
}

func (f *fakeConsumerGroup) Errors() <-chan error { return nil }

func (f *fakeConsumerGroup) Close() error {
	f.mu.Lock()
	f.closes++
	f.mu.Unlock()
	if f.onClose != nil {
		return f.onClose()
	}
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
	t.Cleanup(func() { stopTestFacade(t, facade) })
	return facade, fake
}

func stopTestFacade(t *testing.T, facade *KafkaConsumerFacade) {
	t.Helper()
	done := make(chan struct{})
	go func() {
		facade.Stop()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("Stop did not finish within 3s")
	}
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

func TestNewKafkaConsumerFacadeProductionConstructor(t *testing.T) {
	fake := &fakeConsumerGroup{}
	patches := gomonkey.ApplyFunc(sarama.NewConsumerGroup, func(brokers []string, group string, config *sarama.Config) (sarama.ConsumerGroup, error) {
		require.Equal(t, []string{"broker:9092"}, brokers)
		require.Equal(t, "test-group", group)
		require.Equal(t, "test-client", config.ClientID)
		return fake, nil
	})
	defer patches.Reset()
	config := KafkaConsumerConfig{Brokers: []string{"broker:9092"}, ClientID: "test-client", ProtocolVersion: "2.1.0"}
	facade, err := NewKafkaConsumerFacade(config, "test-group")
	require.NoError(t, err)
	require.NotNil(t, facade.consumerManager)
	stopTestFacade(t, facade)
	config.ProtocolVersion = "invalid"
	_, err = NewKafkaConsumerFacade(config, "test-group")
	require.Error(t, err)

	factoryErr := errors.New("consumer group creation failed")
	facade, err = newKafkaConsumerFacade(config, "test-group", func(KafkaConsumerConfig, string) (sarama.ConsumerGroup, error) {
		return nil, factoryErr
	})
	require.Nil(t, facade)
	require.ErrorIs(t, err, factoryErr)
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
	stopTestFacade(t, facade)
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

	stopTestFacade(t, facade)

	fake.mu.Lock()
	consumes := fake.consumes
	fake.mu.Unlock()
	if consumes == 0 {
		t.Fatal("expected Consume to have been invoked at least once")
	}
}

func TestKafkaConsumerFacadeCancellationStopsWorkers(t *testing.T) {
	for _, cancelParent := range []bool{false, true} {
		name := "registered subscription cancellation"
		if cancelParent {
			name = "parent context cancellation"
		}
		t.Run(name, func(t *testing.T) {
			facade, fake := newTestFacade(t)
			fake.started = make(chan context.Context, 1)
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			topics := []string{"topic"}
			const group = "group"
			if err := facade.Subscribe(ctx, WithTopics(topics), WithConsumerGroup(group)); err != nil {
				t.Fatalf("Subscribe failed: %v", err)
			}
			var consumeCtx context.Context
			select {
			case consumeCtx = <-fake.started:
			case <-time.After(3 * time.Second):
				t.Fatal("Consume did not start")
			}
			key := GetConsumerManagerKey(topics, group)
			if cancelParent {
				cancel()
			} else {
				// This is the same cancellation path used by a failed health check.
				facade.mu.RLock()
				entry := facade.consumerManager[key]
				facade.mu.RUnlock()
				facade.removeConsumer(key, entry)
			}
			select {
			case <-consumeCtx.Done():
			case <-time.After(3 * time.Second):
				t.Fatal("subscription cancellation did not reach Consume")
			}
			done := make(chan struct{})
			go func() {
				facade.wg.Wait()
				close(done)
			}()
			select {
			case <-done:
			case <-time.After(3 * time.Second):
				t.Fatal("subscription workers did not exit after cancellation")
			}
			facade.mu.RLock()
			remaining := len(facade.consumerManager)
			facade.mu.RUnlock()
			if remaining != 0 {
				t.Fatalf("cancellation left %d stale consumers", remaining)
			}
			stopTestFacade(t, facade)
		})
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
			facade.consumerManager[GetConsumerManagerKey([]string{"topic"}, "group")] = &consumerSubscription{cancel: func() {}}
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

func TestKafkaConsumerFacadeReplacesSubscriptionSafely(t *testing.T) {
	facade, fake := newTestFacade(t)
	oldStarted := make(chan struct{})
	newStarted := make(chan context.Context, 1)
	oldCanceled := make(chan struct{})
	releaseOld := make(chan struct{})
	var calls atomic.Int32
	fake.consume = func(ctx context.Context) error {
		if calls.Add(1) == 1 {
			close(oldStarted)
			<-ctx.Done()
			close(oldCanceled)
			<-releaseOld
		} else {
			newStarted <- ctx
			<-ctx.Done()
		}
		return ctx.Err()
	}
	// Release the old worker even if an assertion fails, before facade cleanup.
	t.Cleanup(func() {
		select {
		case <-releaseOld:
		default:
			close(releaseOld)
		}
	})
	opts := []Option{WithTopic("topic"), WithConsumerGroup("group")}
	require.NoError(t, facade.Subscribe(context.Background(), opts...))
	select {
	case <-oldStarted:
	case <-time.After(3 * time.Second):
		t.Fatal("old subscription did not start")
	}
	key := GetConsumerManagerKey([]string{"topic"}, "group")
	facade.mu.RLock()
	oldEntry := facade.consumerManager[key]
	facade.mu.RUnlock()
	require.NoError(t, facade.Subscribe(context.Background(), opts...))
	select {
	case <-oldCanceled:
	case <-time.After(3 * time.Second):
		t.Fatal("replacement did not cancel the old subscription")
	}
	var newCtx context.Context
	select {
	case newCtx = <-newStarted:
	case <-time.After(3 * time.Second):
		t.Fatal("replacement did not start")
	}
	facade.mu.RLock()
	entry := facade.consumerManager[key]
	facade.mu.RUnlock()
	require.NotNil(t, entry, "old health-check cleanup must not delete the replacement")
	// Delayed consume-loop cleanup must also leave the replacement untouched.
	facade.removeConsumer(key, oldEntry)
	close(releaseOld)
	facade.mu.RLock()
	current := facade.consumerManager[key]
	facade.mu.RUnlock()
	require.Same(t, entry, current)
	require.NoError(t, newCtx.Err())
	stopTestFacade(t, facade)
	require.ErrorIs(t, newCtx.Err(), context.Canceled)
}

func TestKafkaConsumerFacadeStopClosesBeforeWaiting(t *testing.T) {
	facade, fake := newTestFacade(t)
	started := make(chan struct{})
	closed := make(chan struct{})
	fake.consume = func(context.Context) error {
		close(started)
		<-closed // Deliberately ignores context cancellation; only Close releases it.
		return sarama.ErrClosedConsumerGroup
	}
	fake.onClose = func() error {
		close(closed)
		return errors.New("close failed")
	}
	require.NoError(t, facade.Subscribe(context.Background(), WithTopic("topic")))
	select {
	case <-started:
	case <-time.After(3 * time.Second):
		t.Fatal("Consume did not start")
	}
	stopTestFacade(t, facade)
	stopTestFacade(t, facade)
	fake.mu.Lock()
	defer fake.mu.Unlock()
	require.Equal(t, 1, fake.closes, "even a failed Close must only run once")
}

func TestKafkaConsumerFacadeConcurrentSubscribeAndStop(t *testing.T) {
	facade, fake := newTestFacade(t)
	start := make(chan struct{})
	var workers sync.WaitGroup
	for i := 0; i < 32; i++ {
		workers.Add(2)
		go func() {
			defer workers.Done()
			<-start
			// Either a successful registration or a stopped error is valid.
			_ = facade.Subscribe(context.Background(), WithTopic("topic"))
		}()
		go func() {
			defer workers.Done()
			<-start
			facade.Stop()
		}()
	}
	close(start)
	done := make(chan struct{})
	go func() {
		workers.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(3 * time.Second):
		t.Fatal("concurrent Subscribe/Stop did not finish")
	}
	require.Error(t, facade.Subscribe(context.Background(), WithTopic("topic")))
	facade.mu.RLock()
	require.Empty(t, facade.consumerManager)
	facade.mu.RUnlock()
	fake.mu.Lock()
	require.Equal(t, 1, fake.closes)
	fake.mu.Unlock()
}

func TestKafkaConsumerFacadeHealthCheckLifecycle(t *testing.T) {
	for _, blocked := range []bool{false, true} {
		name := "unhealthy consumer cancels both workers"
		if blocked {
			name = "stop cancels in-flight health request"
		}
		t.Run(name, func(t *testing.T) {
			ticks := make(chan time.Time, 1)
			ticker := time.NewTicker(time.Hour)
			ticker.C = ticks
			defer ticker.Stop()
			patches := gomonkey.ApplyFunc(time.NewTicker, func(time.Duration) *time.Ticker { return ticker })
			defer patches.Reset()
			requestStarted := make(chan struct{}, 5)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestStarted <- struct{}{}
				if blocked {
					<-r.Context().Done()
					return
				}
				w.WriteHeader(http.StatusServiceUnavailable)
			}))
			defer server.Close()
			facade, fake := newTestFacade(t)
			fake.started = make(chan context.Context, 1)
			// Cleanup must cancel any request before closing the HTTP server.
			defer stopTestFacade(t, facade)
			require.NoError(t, facade.Subscribe(context.Background(), WithTopic("topic"), WithCheckUrl(server.URL)))
			var consumeCtx context.Context
			select {
			case consumeCtx = <-fake.started:
			case <-time.After(3 * time.Second):
				t.Fatal("Consume did not start")
			}
			ticks <- time.Now()
			select {
			case <-requestStarted:
			case <-time.After(3 * time.Second):
				t.Fatal("health check did not start")
			}
			if blocked {
				stopTestFacade(t, facade)
			}
			select {
			case <-consumeCtx.Done():
			case <-time.After(3 * time.Second):
				t.Fatal("health-check shutdown did not cancel Consume")
			}
			stopTestFacade(t, facade)
			facade.mu.RLock()
			remaining := len(facade.consumerManager)
			facade.mu.RUnlock()
			require.Zero(t, remaining, "health-check shutdown must remove the subscription")
		})
	}
}
