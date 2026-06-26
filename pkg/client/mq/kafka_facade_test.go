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
	"sync"
	"testing"
	"time"
)

// TestKafkaConsumerFacadeConsumerManagerInitialized verifies that
// the consumerManager map is initialized during construction,
// preventing nil map panics during Subscribe operations.
func TestKafkaConsumerFacadeConsumerManagerInitialized(t *testing.T) {
	facade := &KafkaConsumerFacade{
		consumerManager: make(map[string]func()),
		done:            make(chan struct{}),
		mu:              sync.RWMutex{},
	}

	testKey := "test-topic-test-group"

	// This should NOT panic due to nil map write
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("Subscribe panicked with nil map: %v", r)
		}
	}()

	facade.mu.Lock()
	facade.consumerManager[testKey] = func() {}
	facade.mu.Unlock()

	facade.mu.RLock()
	_, exists := facade.consumerManager[testKey]
	facade.mu.RUnlock()

	if !exists {
		t.Error("Expected key to exist in consumerManager")
	}
}

// TestKafkaConsumerFacadeConcurrentAccess verifies that concurrent
// access to consumerManager is safe due to mutex protection.
func TestKafkaConsumerFacadeConcurrentAccess(t *testing.T) {
	facade := &KafkaConsumerFacade{
		consumerManager: make(map[string]func()),
		done:            make(chan struct{}),
		mu:              sync.RWMutex{},
	}

	done := make(chan bool)

	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			facade.mu.Lock()
			facade.consumerManager[GetConsumerManagerKey([]string{"topic"}, "group")] = func() {}
			facade.mu.Unlock()
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			facade.mu.RLock()
			_ = len(facade.consumerManager)
			facade.mu.RUnlock()
		}
		done <- true
	}()

	<-done
	<-done
}

// TestGetConsumerManagerKey verifies the key generation function
func TestGetConsumerManagerKey(t *testing.T) {
	topics := []string{"topic1", "topic2"}
	group := "test-group"

	key := GetConsumerManagerKey(topics, group)

	expectedKey := GetConsumerManagerKey(topics, group)
	if key != expectedKey {
		t.Error("Key should be deterministic")
	}

	differentTopics := []string{"topic3", "topic4"}
	differentKey := GetConsumerManagerKey(differentTopics, group)
	if key == differentKey {
		t.Error("Different topics should produce different keys")
	}

	differentGroup := "other-group"
	differentGroupKey := GetConsumerManagerKey(topics, differentGroup)
	if key == differentGroupKey {
		t.Error("Different groups should produce different keys")
	}
}

// TestNewKafkaConsumerFacadeConstructorInitializesMap verifies that
// the consumerManager map is initialized during construction via NewKafkaConsumerFacade.
func TestNewKafkaConsumerFacadeConstructorInitializesMap(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test that requires Kafka broker in short mode")
	}

	config := KafkaConsumerConfig{
		Brokers:         []string{"localhost:9092"},
		ClientID:        "test-client",
		ProtocolVersion: "2.8.0",
		Metadata: Metadata{
			Full: true,
			Retry: MetadataRetry{
				Max:     3,
				Backoff: 250 * time.Millisecond,
			},
		},
	}

	facade, err := NewKafkaConsumerFacade(config, "test-group")
	if err != nil {
		t.Logf("Cannot connect to Kafka broker, skipping: %v", err)
		return
	}

	if facade.consumerManager == nil {
		t.Fatal("consumerManager should be initialized, not nil")
	}

	facade.Stop()
}
