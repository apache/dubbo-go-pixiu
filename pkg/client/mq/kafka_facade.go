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
	"net/http"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/event"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

import (
	"github.com/Shopify/sarama"
	perrors "github.com/pkg/errors"
)

func NewKafkaConsumerFacade(config event.KafkaConsumerConfig) (*KafkaConsumerFacade, error) {
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
	consumer, err := sarama.NewConsumer(config.Brokers, c)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumerFacade{consumer: consumer, httpClient: &http.Client{Timeout: 5 * time.Second}, done: make(chan struct{})}, nil
}

type KafkaConsumerFacade struct {
	consumer        sarama.Consumer
	consumerManager map[string]func()
	rwLock          sync.RWMutex
	httpClient      *http.Client
	wg              sync.WaitGroup
	done            chan struct{}
}

func (f *KafkaConsumerFacade) Subscribe(ctx context.Context, opts ...Option) error {
	cOpt := DefaultOptions()
	cOpt.ApplyOpts(opts...)
	partConsumer, err := f.consumer.ConsumePartition(cOpt.Topic, int32(cOpt.Partition), sarama.OffsetOldest)
	if err != nil {
		return err
	}
	c, cancel := context.WithCancel(ctx)
	key := GetConsumerManagerKey(cOpt.Topic, int32(cOpt.Partition))
	f.rwLock.Lock()
	defer f.rwLock.Unlock()
	f.consumerManager[key] = cancel
	f.wg.Add(2)
	go f.ConsumePartitions(c, partConsumer, cOpt.ConsumeUrl)
	go f.checkConsumerIsAlive(c, key, cOpt.CheckUrl)
	return nil
}

// ConsumePartitions consume function
func (f *KafkaConsumerFacade) ConsumePartitions(ctx context.Context, partConsumer sarama.PartitionConsumer, consumeUrl string) {
	defer f.wg.Done()

	for {
		select {
		case <-f.done:
			logger.Info()
		case msg := <-partConsumer.Messages():
			data, err := json.Marshal(event.MQMsgPush{Msg: []string{string(msg.Value)}})
			if err != nil {
				logger.Warn()
				continue
			}

			req, err := http.NewRequest(http.MethodPost, consumeUrl, bytes.NewReader(data))
			if err != nil {
				logger.Warn()
				continue
			}

			for i := 0; i < 5; i++ {
				err := func() error {
					resp, err := f.httpClient.Do(req)
					if err != nil {
						return err
					}
					defer resp.Body.Close()

					if resp.StatusCode == http.StatusOK {
						return nil
					}
					return perrors.New("failed send msg to consumer with status code " + strconv.Itoa(resp.StatusCode))
				}()
				if err != nil {
					logger.Warn(err.Error())
					time.Sleep(10 * time.Millisecond)
				} else {
					break
				}
			}
		}
	}
}

// checkConsumerIsAlive make sure consumer is alive or would be removed from consumer list
func (f *KafkaConsumerFacade) checkConsumerIsAlive(ctx context.Context, key string, checkUrl string) {
	defer f.wg.Done()

	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-f.done:
			logger.Info()
		case <-ticker.C:
			lastCheck := 0
			for i := 0; i < 5; i++ {
				err := func() error {
					req, err := http.NewRequest(http.MethodGet, checkUrl, bytes.NewReader([]byte{}))
					if err != nil {
						logger.Warn()
					}

					resp, err := f.httpClient.Do(req)
					if err != nil {
						logger.Warn()
					}
					defer resp.Body.Close()

					lastCheck = resp.StatusCode
					if resp.StatusCode != http.StatusOK {
						return perrors.New("failed check consumer alive or not with status code " + strconv.Itoa(resp.StatusCode))
					}

					logger.Warn()
					return nil
				}()
				if err != nil {
					logger.Warn()
					time.Sleep(10 * time.Millisecond)
				} else {
					break
				}
			}

			if lastCheck != http.StatusOK {
				f.consumerManager[key]()
				delete(f.consumerManager, key)
			}

		}
	}
}

func (f *KafkaConsumerFacade) UnSubscribe(opts ...Option) error {
	cOpt := DefaultOptions()
	cOpt.ApplyOpts(opts...)
	key := GetConsumerManagerKey(cOpt.Topic, int32(cOpt.Partition))
	if cancel, ok := f.consumerManager[key]; !ok {
		return perrors.New("consumer goroutine not found")
	} else {
		cancel()
		delete(f.consumerManager, key)
	}
	return nil
}

func (f *KafkaConsumerFacade) Stop() {
	close(f.done)
	f.wg.Wait()
}

func NewKafkaProviderFacade(config event.KafkaProducerConfig) (*KafkaProducerFacade, error) {
	c := sarama.NewConfig()
	c.Producer.Return.Successes = true
	c.Producer.Return.Errors = true
	c.Producer.RequiredAcks = sarama.WaitForLocal
	c.Metadata.Full = config.Metadata.Full
	c.Metadata.Retry.Max = config.Metadata.Retry.Max
	c.Metadata.Retry.Backoff = config.Metadata.Retry.Backoff
	c.Producer.MaxMessageBytes = config.Producer.MaxMessageBytes
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
		pMsgs = append(pMsgs, &sarama.ProducerMessage{Topic: pOpt.Topic, Value: sarama.StringEncoder(msg)})
	}

	return k.producer.SendMessages(pMsgs)
}
