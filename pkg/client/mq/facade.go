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
	"strings"
)

type ConsumerFacade interface {
	// Subscribe message with specified broker and Topic, then handle msg with handler which send msg to real consumers
	Subscribe(ctx context.Context, option ...Option) error
	UnSubscribe(opts ...Option) error
	Stop()
}

func GetConsumerManagerKey(topic []string, consumerGroup string) string {
	return strings.Join((topic), ".") + "_" + consumerGroup
}

// MQOptions Consumer options
// TODO: Add rocketmq params
type MQOptions struct {
	TopicList     []string
	ConsumeUrl    string
	CheckUrl      string
	ConsumerGroup string
}

func (o *MQOptions) ApplyOpts(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func DefaultOptions() *MQOptions {
	return &MQOptions{}
}

type Option func(o *MQOptions)

func WithTopics(t []string) Option {
	return func(o *MQOptions) {
		o.TopicList = t
	}
}

func WithTopic(t string) Option {
	return func(o *MQOptions) {
		o.TopicList = []string{t}
	}
}

func WithConsumeUrl(ch string) Option {
	return func(o *MQOptions) {
		o.ConsumeUrl = ch
	}
}

func WithCheckUrl(ck string) Option {
	return func(o *MQOptions) {
		o.CheckUrl = ck
	}
}

func WithConsumerGroup(cp string) Option {
	return func(o *MQOptions) {
		o.ConsumerGroup = cp
	}
}

type ProducerFacade interface {
	// Send msg to specified broker and Topic
	Send(msgs []string, opts ...Option) error
}
