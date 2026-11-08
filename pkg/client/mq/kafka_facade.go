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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/IBM/sarama"

	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

type kafkaErrors struct {
	count int
	err   string
}

func (ke kafkaErrors) Error() string {
	return fmt.Sprintf("Failed to deliver %d messages due to %s", ke.count, ke.err)
}

// consumerGroupFactory builds a sarama.ConsumerGroup from broker config.
// It is abstracted so tests can inject a fake group without a real Kafka broker.
type consumerGroupFactory func(config KafkaConsumerConfig, consumerGroup string) (sarama.ConsumerGroup, error)

// defaultConsumerGroupFactory builds a real sarama consumer group.
func defaultConsumerGroupFactory(config KafkaConsumerConfig, consumerGroup string) (sarama.ConsumerGroup, error) {
	c := sarama.NewConfig()
	c.ClientID = config.ClientID
	c.Metadata.Full = config.Metadata.Full
	c.Metadata.Retry.Max = config.Metadata.Retry.Max
	c.Metadata.Retry.Backoff = config.Metadata.Retry.Backoff
	if config.ProtocolVersion != "" {
		version, err := sarama.ParseKafkaVersion(config.ProtocolVersion)
		if err != nil {
			return nil, err
		}
		c.Version = version
	}
	return sarama.NewConsumerGroup(config.Brokers, consumerGroup, c)
}

// NewKafkaConsumerFacade creates a KafkaConsumerFacade backed by a real sarama
// consumer group.
func NewKafkaConsumerFacade(config KafkaConsumerConfig, consumerGroup string) (*KafkaConsumerFacade, error) {
	return newKafkaConsumerFacade(config, consumerGroup, defaultConsumerGroupFactory)
}

// newKafkaConsumerFacade is the testable constructor: it accepts a consumer
// group factory so the consumerManager initialization and shutdown behavior can
// be exercised without a live Kafka broker.
func newKafkaConsumerFacade(config KafkaConsumerConfig, consumerGroup string, factory consumerGroupFactory) (*KafkaConsumerFacade, error) {
	client, err := factory(config, consumerGroup)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumerFacade{
		consumerGroup:   client,
		consumerManager: make(map[string]*consumerSubscription),
		httpClient:      &http.Client{Timeout: 5 * time.Second},
		done:            make(chan struct{}),
	}, nil
}

type consumerSubscription struct {
	cancel context.CancelFunc
}

type KafkaConsumerFacade struct {
	consumerGroup   sarama.ConsumerGroup
	consumerManager map[string]*consumerSubscription
	mu              sync.RWMutex // protects subscriptions, stopped, and wg.Add
	stopped         bool
	httpClient      *http.Client
	wg              sync.WaitGroup
	done            chan struct{}
	stopOnce        sync.Once
}

// Subscribe starts a subscription whose lifetime follows ctx. Repeating a
// topic/group key cancels and replaces the previous subscription. A stopped
// facade cannot be reused.
func (f *KafkaConsumerFacade) Subscribe(ctx context.Context, opts ...Option) error {
	cOpt := DefaultOptions()
	cOpt.ApplyOpts(opts...)
	key := GetConsumerManagerKey(cOpt.TopicList, cOpt.ConsumerGroup)
	f.mu.Lock()
	if f.stopped {
		f.mu.Unlock()
		return perrors.New("kafka consumer facade is stopped")
	}
	// Replace the old subscription atomically. Each worker retains its entry's
	// identity, so its eventual cleanup cannot remove a newer subscription.
	c, cancel := context.WithCancel(ctx)
	entry := &consumerSubscription{cancel: cancel}
	if previous := f.consumerManager[key]; previous != nil {
		previous.cancel()
	}
	f.consumerManager[key] = entry
	// Stop takes the same lock before waiting, so Add cannot race with Wait.
	f.wg.Add(2)
	f.mu.Unlock()
	go f.consumeLoop(c, cOpt.TopicList, &consumerGroupHandler{cOpt.ConsumeUrl, f.httpClient}, key, entry)
	go f.checkConsumerIsAlive(c, key, cOpt.CheckUrl, entry)
	return nil
}

// consumeLoop repeatedly joins the consumer group until either the subscription
// context is canceled (e.g. consumer deemed unhealthy, or parent shutdown) or
// the facade is stopping via f.done. It always decrements the WaitGroup on exit
// and removes its consumerManager entry so Stop()'s wg.Wait() can return and no
// stale cancel funcs are left behind.
func (f *KafkaConsumerFacade) consumeLoop(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler, key string, entry *consumerSubscription) {
	defer f.wg.Done()
	defer f.removeConsumer(key, entry)

	for {
		// Consume blocks until the session ends (rebalance, ctx cancel, or
		// Close). On a clean ctx cancellation we stop; otherwise we rejoin.
		if err := f.consumerGroup.Consume(ctx, topics, handler); err != nil {
			logger.Warnf("failed to consume the msg from kafka, %s", err.Error())
		}

		select {
		case <-f.done:
			logger.Info("shutdown the consume loop")
			return
		case <-ctx.Done():
			logger.Infof("shutdown the consume loop due to %s", ctx.Err().Error())
			return
		default:
			// session ended (e.g. rebalance) with an active context; rejoin
		}
	}
}

type consumerGroupHandler struct {
	consumerUrl string
	httpClient  *http.Client
}

func (c *consumerGroupHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		session.MarkMessage(msg, "")
		data, err := json.Marshal(MQMsgPush{Msg: []string{string(msg.Value)}})
		if err != nil {
			logger.Warn()
			continue
		}

		req, err := http.NewRequest(http.MethodPost, c.consumerUrl, bytes.NewReader(data))
		if err != nil {
			logger.Warn()
			continue
		}
		err = func() error {
			resp, err := c.httpClient.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusOK {
				return nil
			}
			return perrors.New("failed send msg to consumer")
		}()
		if err != nil {
			logger.Warn(err.Error())
		}
	}
	return nil
}

// removeConsumer cancels and deletes only the subscription owned by the caller.
// A worker from a replaced subscription must not cancel or delete its successor.
func (f *KafkaConsumerFacade) removeConsumer(key string, entry *consumerSubscription) {
	f.mu.Lock()
	if current, ok := f.consumerManager[key]; ok && current == entry {
		entry.cancel()
		delete(f.consumerManager, key)
	}
	f.mu.Unlock()
}

// checkConsumerIsAlive periodically checks the consumer liveness endpoint. When
// the consumer is deemed unhealthy it cancels the consume loop (via the shared
// subscription context) and removes the manager entry. It also removes the
// entry on either shutdown signal so no stale cancel funcs are left behind.
func (f *KafkaConsumerFacade) checkConsumerIsAlive(ctx context.Context, key string, checkUrl string, entry *consumerSubscription) {
	defer f.wg.Done()
	defer f.removeConsumer(key, entry)

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-f.done:
			return
		case <-ctx.Done():
			return
		case <-ticker.C:
			lastCheck := 0
			for i := 0; i < 5; i++ {
				err := func() error {
					req, err := http.NewRequestWithContext(ctx, http.MethodGet, checkUrl, nil)
					if err != nil {
						return perrors.Wrap(err, "create consumer health request")
					}
					resp, err := f.httpClient.Do(req)
					if err != nil {
						return perrors.Wrap(err, "check consumer health")
					}
					defer resp.Body.Close()
					lastCheck = resp.StatusCode
					if resp.StatusCode != http.StatusOK {
						return perrors.New("failed check consumer alive or not with status code " + strconv.Itoa(resp.StatusCode))
					}
					return nil
				}()
				if err != nil {
					select {
					case <-ctx.Done():
						return
					case <-time.After(10 * time.Millisecond):
					}
				} else {
					break
				}
			}

			if lastCheck != http.StatusOK {
				// Consumer is unhealthy: cancel the shared subscription context
				// to stop the consume loop and drop the manager entry.
				return
			}
		}
	}
}

func (f *KafkaConsumerFacade) UnSubscribe(opts ...Option) error {
	return nil
}

func (f *KafkaConsumerFacade) Stop() {
	f.stopOnce.Do(func() {
		// Cancel any still-registered consumers before signaling shutdown so the
		// consume loops can exit their Consume() calls promptly.
		f.mu.Lock()
		f.stopped = true
		for key, entry := range f.consumerManager {
			entry.cancel()
			delete(f.consumerManager, key)
		}
		close(f.done)
		f.mu.Unlock()

		// Close before waiting: sarama also uses its closed channel to release
		// Consume during session setup/rejoin. Waiting first could prevent the
		// very Close call needed to unblock it. Sarama's network timeouts bound
		// its outstanding broker requests. Close and Wait both run only once.
		if err := f.consumerGroup.Close(); err != nil {
			logger.Warnf("failed to close kafka consumer group: %s", err.Error())
		}
		f.wg.Wait()
	})
}

func NewKafkaProviderFacade(config KafkaProducerConfig) (*KafkaProducerFacade, error) {
	c := sarama.NewConfig()
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true
	c.Producer.RequiredAcks = sarama.WaitForLocal
	c.Metadata.Full = config.Metadata.Full
	c.Metadata.Retry.Max = config.Metadata.Retry.Max
	c.Metadata.Retry.Backoff = config.Metadata.Retry.Backoff
	c.Producer.MaxMessageBytes = config.Producer.MaxMessageBytes
	c.Producer.Timeout = config.Timeout
	if config.ProtocolVersion != "" {
		version, err := sarama.ParseKafkaVersion(config.ProtocolVersion)
		if err != nil {
			return nil, err
		}
		c.Version = version
	}
	producer, err := sarama.NewSyncProducer(config.Brokers, c)
	if err != nil {
		return nil, err
	}
	return &KafkaProducerFacade{producer: producer}, nil
}

type KafkaProducerFacade struct {
	producer sarama.SyncProducer
}

func (k *KafkaProducerFacade) Send(msgs []string, opts ...Option) error {
	pOpt := DefaultOptions()
	pOpt.ApplyOpts(opts...)

	pMsgs := make([]*sarama.ProducerMessage, 0)
	for _, msg := range msgs {
		pMsgs = append(pMsgs, &sarama.ProducerMessage{Topic: pOpt.TopicList[0], Value: sarama.StringEncoder(msg)})
	}
	err := k.producer.SendMessages(pMsgs)
	if err != nil {
		if value, ok := err.(sarama.ProducerErrors); ok {
			if len(value) > 0 {
				return kafkaErrors{len(value), value[0].Err.Error()}
			}
		}
		return err
	}
	return nil
}
