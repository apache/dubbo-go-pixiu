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

package event

import (
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPEventFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}

	// Filter is http filter instance
	Filter struct {
		cfg *Config
	}

	Config struct {
		ClientId            string              `yaml:"client_id" json:"client_id"`
		Endpoints           string              `yaml:"endpoints" json:"endpoints"`
		MqType              MQType              `yaml:"type" json:"type"`
		Retry               int                 `yaml:"retry" json:"retry" default:"5"`
		Timeout             time.Duration       `yaml:"timeout" json:"timeout" default:"2s"`
		KafkaConsumerConfig KafkaConsumerConfig `yaml:"kafka_consumer_config" json:"kafka_consumer_config"`
		KafkaProducerConfig KafkaProducerConfig `yaml:"kafka_producer_config" json:"kafka_producer_config"`
	}

	KafkaConsumerConfig struct {
		Brokers         []string `yaml:"brokers" json:"brokers"`
		ProtocolVersion string   `yaml:"protocol_version" json:"protocol_version"`
		ClientID        string   `yaml:"client_id" json:"client_id"`
		Metadata        Metadata `yaml:"metadata" json:"metadata"`
	}

	KafkaProducerConfig struct {
		Brokers         []string `yaml:"brokers" json:"brokers"`
		ProtocolVersion string   `yaml:"protocol_version" json:"protocol_version"`
		Metadata        Metadata `yaml:"metadata" json:"metadata"`
		Producer        Producer `yaml:"producer" json:"producer"`
	}

	Metadata struct {
		Full  bool          `yaml:"full" json:"full"`
		Retry MetadataRetry `yaml:"retry" json:"retry"`
	}

	MetadataRetry struct {
		Max     int           `yaml:"max" json:"max"`
		Backoff time.Duration `yaml:"backoff" json:"backoff"`
	}

	Producer struct {
		MaxMessageBytes int `yaml:"max_message_bytes" json:"max_message_bytes"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &Config{}}, nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(ctx *http.HttpContext) {
	// TODO handle request
}

func (f *Filter) Apply() error {
	// TODO init mq client here
	return nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}
