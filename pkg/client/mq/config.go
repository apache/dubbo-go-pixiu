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
	"time"
)

type (
	Config struct {
		ClientID            string              `yaml:"client_id" json:"client_id"`
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
		Brokers         []string      `yaml:"brokers" json:"brokers"`
		ProtocolVersion string        `yaml:"protocol_version" json:"protocol_version"`
		Metadata        Metadata      `yaml:"metadata" json:"metadata"`
		Producer        Producer      `yaml:"producer" json:"producer"`
		Timeout         time.Duration `yaml:"timeout" json:"timeout"`
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
