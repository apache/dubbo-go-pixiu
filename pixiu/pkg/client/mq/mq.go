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
	"encoding/json"
	"io"
	"strings"
	"sync"
)

import (
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

var (
	mqClient          *Client
	once              sync.Once
	consumerFacadeMap sync.Map
)

func NewSingletonMQClient(config Config) *Client {
	if mqClient == nil {
		once.Do(func() {
			var err error
			mqClient, err = NewMQClient(config)
			if err != nil {
				logger.Errorf("create mq client failed, %s", err.Error())
			}
		})
	}
	return mqClient
}

func NewMQClient(config Config) (*Client, error) {
	var c *Client
	ctx := context.Background()
	switch config.MqType {
	case constant.MQTypeKafka:
		config.KafkaProducerConfig.Timeout = config.Timeout
		pf, err := NewKafkaProviderFacade(config.KafkaProducerConfig)
		if err != nil {
			return nil, err
		}
		c = &Client{
			ctx:                 ctx,
			producerFacade:      pf,
			kafkaConsumerConfig: config.KafkaConsumerConfig,
		}
	case constant.MQTypeRocketMQ:
		return nil, perrors.New("rocketmq not support")
	}

	return c, nil
}

type Client struct {
	ctx                 context.Context
	producerFacade      ProducerFacade
	kafkaConsumerConfig KafkaConsumerConfig
}

func (c Client) Apply() error {
	panic("implement me")
}

func (c Client) Close() error {
	return nil
}

func (c Client) Call(req *client.Request) (res interface{}, err error) {
	body, err := io.ReadAll(req.IngressRequest.Body)
	if err != nil {
		return nil, err
	}

	paths := strings.Split(req.API.Path, "/")
	if len(paths) < 3 {
		return nil, perrors.New("failed to send message, broker or Topic not found")
	}

	switch MQActionStrToInt[paths[0]] {
	case MQActionPublish:
		var pReq MQProduceRequest
		err = json.Unmarshal(body, &pReq)
		if err != nil {
			return nil, err
		}
		err = c.producerFacade.Send(pReq.Msg, WithTopic(pReq.Topic))
		if err != nil {
			return nil, err
		}
	case MQActionSubscribe:
		var cReq MQSubscribeRequest
		err = json.Unmarshal(body, &cReq)
		if err != nil {
			return nil, err
		}
		if _, ok := consumerFacadeMap.Load(cReq.ConsumerGroup); !ok {
			facade, err := NewKafkaConsumerFacade(c.kafkaConsumerConfig, cReq.ConsumerGroup)
			if err != nil {
				return nil, err
			}
			consumerFacadeMap.Store(cReq.ConsumerGroup, facade)
			if f, ok := consumerFacadeMap.Load(cReq.ConsumerGroup); ok {
				cf := f.(ConsumerFacade)
				ctx, cancel := context.WithTimeout(c.ctx, req.Timeout)
				defer cancel()
				err = cf.Subscribe(ctx, WithTopics(cReq.TopicList), WithConsumeUrl(cReq.ConsumeUrl), WithCheckUrl(cReq.CheckUrl), WithConsumerGroup(cReq.ConsumerGroup))
				if err != nil {
					facade.Stop()
					consumerFacadeMap.Delete(cReq.ConsumerGroup)
					return nil, err
				}
			}
		}
	case MQActionUnSubscribe:
		var cReq MQUnSubscribeRequest
		err = json.Unmarshal(body, &cReq)
		if err != nil {
			return nil, err
		}
		if facade, ok := consumerFacadeMap.Load(cReq.ConsumerGroup); ok {
			facade.(ConsumerFacade).Stop()
			consumerFacadeMap.Delete(cReq.ConsumerGroup)
			return nil, err
		}
	default:
		return nil, perrors.New("failed to get mq action")
	}

	return nil, nil
}

func (c Client) MapParams(req *client.Request) (reqData interface{}, err error) {
	return nil, perrors.New("map params does not support on mq mqClient")
}
