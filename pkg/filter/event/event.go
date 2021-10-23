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
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

import (
	"github.com/Shopify/sarama"
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
		ClientId     string        `yaml:"client_id" json:"client_id"`
		Endpoints    string        `yaml:"endpoints" json:"endpoints"`
		MqType       MQType        `yaml:"type" json:"type"`
		KafkaVersion string        `yaml:"kafka_version" json:"kafka_version"`
		Retry        int           `yaml:"retry" json:"retry" default:"5"`
		Timeout      time.Duration `yaml:"timeout" json:"timeout" default:"2s"`
		Offset       string        `yaml:"offset" json:"offset" default:"oldest"` // Offset newest or oldest
	}
)

func (c *Config) ToSaramaConfig() *sarama.Config {
	config := sarama.NewConfig()

	version, err := sarama.ParseKafkaVersion(c.KafkaVersion)
	if err != nil {
		version = sarama.V2_0_0_0
		logger.Warnf("kafka version is invalid, use sarama.V2_0_0_0 instead, err: %s", err.Error())
	}
	config.Version = version

	offset := sarama.OffsetNewest
	switch c.Offset {
	case "newest":
		offset = sarama.OffsetNewest
	case "oldest":
		offset = sarama.OffsetOldest
	default:
		logger.Warn("offset is invalid, use oldest instead")
	}
	config.Consumer.Offsets.Initial = offset

	config.ClientID = "pixiu-kafka"
	if c.ClientId != "" {
		config.ClientID = c.ClientId
	}
	logger.Infof("kafka client id is %s", config.ClientID)

	config.Producer.Retry.Max = c.Retry
	config.Producer.Timeout = c.Timeout

	return config
}

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
