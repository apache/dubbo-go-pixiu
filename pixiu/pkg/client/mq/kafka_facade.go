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
	"github.com/Shopify/sarama"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

type kafkaErrors struct {
	count int
	err   string
}

func (ke kafkaErrors) Error() string {
	return fmt.Sprintf("Failed to deliver %d messages due to %s", ke.count, ke.err)
}

func NewKafkaConsumerFacade(config KafkaConsumerConfig, consumerGroup string) (*KafkaConsumerFacade, error) {
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
	client, err := sarama.NewConsumerGroup(config.Brokers, consumerGroup, c)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumerFacade{consumerGroup: client, httpClient: &http.Client{Timeout: 5 * time.Second}, done: make(chan struct{})}, nil
}

type KafkaConsumerFacade struct {
	consumerGroup   sarama.ConsumerGroup
	consumerManager map[string]func()
	httpClient      *http.Client
	wg              sync.WaitGroup
	done            chan struct{}
}

func (f *KafkaConsumerFacade) Subscribe(ctx context.Context, opts ...Option) error {
	cOpt := DefaultOptions()
	cOpt.ApplyOpts(opts...)
	c, cancel := context.WithCancel(ctx)
	key := GetConsumerManagerKey(cOpt.TopicList, cOpt.ConsumerGroup)
	f.consumerManager[key] = cancel
	f.wg.Add(2)
	go f.consumeLoop(ctx, cOpt.TopicList, &consumerGroupHandler{cOpt.ConsumeUrl, f.httpClient})
	go f.checkConsumerIsAlive(c, key, cOpt.CheckUrl)
	return nil
}

func (f *KafkaConsumerFacade) consumeLoop(ctx context.Context, topics []string, handler sarama.ConsumerGroupHandler) {
	for {
		if _, ok := <-f.done; ok {
			logger.Info("shutdown the consume loop")
			break
		}
		if err := f.consumerGroup.Consume(ctx, topics, handler); err != nil {
			logger.Warn("failed to consume the msg from kafka, %s", err.Error())
		}
		if ctx.Err() != nil {
			// log consume stop
			logger.Error("shutdown the consume loop due to %s", ctx.Err().Error())
			break
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

// checkConsumerIsAlive make sure consumer is alive or would be removed from consumer list
func (f *KafkaConsumerFacade) checkConsumerIsAlive(ctx context.Context, key string, checkUrl string) {
	defer f.wg.Done()

	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-f.done:
			ticker.Stop()
		case <-ticker.C:
			lastCheck := 0
			for i := 0; i < 5; i++ {
				err := func() error {
					req, err := http.NewRequest(http.MethodGet, checkUrl, bytes.NewReader([]byte{}))
					if err == nil {
						resp, err := f.httpClient.Do(req)
						if err == nil {
							defer resp.Body.Close()
							lastCheck = resp.StatusCode
							if resp.StatusCode != http.StatusOK {
								return perrors.New("failed check consumer alive or not with status code " + strconv.Itoa(resp.StatusCode))
							}
							return nil
						}
					}
					return perrors.New("failed to check consumer alive due to: " + err.Error())
				}()
				if err != nil {
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
	return nil
}

func (f *KafkaConsumerFacade) Stop() {
	close(f.done)
	f.wg.Wait()
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
