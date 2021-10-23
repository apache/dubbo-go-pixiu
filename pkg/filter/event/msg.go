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

// Msq Request Action Type Enum

type MQType string

type MQAction int

const (
	MQActionPublish = 1 + iota
	MQActionSubscribe
	MQActionUnSubscribe
	MQActionConsumeAck
)

var MQActionIntToStr = map[MQAction]string{
	MQActionPublish:     "publish",
	MQActionSubscribe:   "subscribe",
	MQActionUnSubscribe: "unsubscribe",
	MQActionConsumeAck:  "consumeack",
}

var MQActionStrToInt = map[string]MQAction{
	"publish":     MQActionPublish,
	"subscribe":   MQActionSubscribe,
	"unsubscribe": MQActionUnSubscribe,
	"consumeack":  MQActionConsumeAck,
}

// MQRequest url format http://domain/publish/broker/topic
type MQRequest struct {
	Partition  int      `json:"partition"`   // for kafka
	ConsumeUrl string   `json:"consume_url"` // not empty when subscribe msg, eg: http://10.0.0.1:11451/consume
	CheckUrl   string   `json:"check_url"`   // not empty when subscribe msg, eg: http://10.0.0.1:11451/health
	Msg        []string `json:"msg"`         // not empty when publish msg, msg body
}

type MQMsgPush struct {
	Msg []string `json:"msg"`
}
