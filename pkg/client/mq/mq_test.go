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
	"net/http/httptest"
	"testing"
	"time"
)

import (
	"github.com/agiledragon/gomonkey/v2"

	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

func TestMQSubscriptionOutlivesRegistrationRequest(t *testing.T) {
	for _, unsubscribe := range []bool{false, true} {
		name := "client cancellation"
		if unsubscribe {
			name = "unsubscribe"
		}
		t.Run(name, func(t *testing.T) {
			facade, fake := newTestFacade(t)
			fake.started = make(chan context.Context, 1)
			patches := gomonkey.ApplyFunc(NewKafkaConsumerFacade, func(KafkaConsumerConfig, string) (*KafkaConsumerFacade, error) {
				return facade, nil
			})
			defer patches.Reset()
			group := t.Name()
			t.Cleanup(func() { consumerFacadeMap.Delete(group) })
			clientCtx, cancelClient := context.WithCancel(context.Background())
			defer cancelClient()
			mqClient := Client{ctx: clientCtx}
			body, err := json.Marshal(MQSubscribeRequest{ConsumerGroup: group, TopicList: []string{"topic"}})
			require.NoError(t, err)
			requestCtx, cancelRequest := context.WithCancel(context.Background())
			defer cancelRequest()
			api := router.API{}
			api.Path = "subscribe/kafka/topic"
			request := &client.Request{
				Context:        requestCtx,
				IngressRequest: httptest.NewRequest("POST", "/subscribe/kafka/topic", bytes.NewReader(body)).WithContext(requestCtx),
				API:            api,
				Timeout:        time.Nanosecond,
			}
			_, err = mqClient.Call(request)
			require.NoError(t, err)
			cancelRequest()
			var consumeCtx context.Context
			select {
			case consumeCtx = <-fake.started:
			case <-time.After(3 * time.Second):
				t.Fatal("registered consumer did not start")
			}
			require.NoError(t, consumeCtx.Err(), "returning from Call must not cancel a subscription")
			_, hasDeadline := consumeCtx.Deadline()
			require.False(t, hasDeadline, "HTTP registration timeout must not bound subscription lifetime")
			if unsubscribe {
				body, err = json.Marshal(MQUnSubscribeRequest{ConsumerGroup: group})
				require.NoError(t, err)
				api.Path = "unsubscribe/kafka/topic"
				_, err = mqClient.Call(&client.Request{
					IngressRequest: httptest.NewRequest("POST", "/unsubscribe/kafka/topic", bytes.NewReader(body)),
					API:            api,
				})
				require.NoError(t, err)
				_, exists := consumerFacadeMap.Load(group)
				require.False(t, exists)
			} else {
				cancelClient()
			}
			select {
			case <-consumeCtx.Done():
			case <-time.After(3 * time.Second):
				t.Fatal("subscription did not stop")
			}
			stopTestFacade(t, facade)
		})
	}
}
