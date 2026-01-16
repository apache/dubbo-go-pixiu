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

package model

import (
	"time"
)

// PixiuInstance represents a connected Pixiu gateway instance.
type PixiuInstance struct {
	NodeID      string            `json:"node_id"`
	Address     string            `json:"address"`
	Cluster     string            `json:"cluster"`
	Version     string            `json:"version"`
	Metadata    map[string]string `json:"metadata"`
	Status      string            `json:"status"` // connected, disconnected
	LastSeen    time.Time         `json:"last_seen"`
	ConnectedAt time.Time         `json:"connected_at"`
}

// Instance status constants
const (
	InstanceStatusConnected    = "connected"
	InstanceStatusDisconnected = "disconnected"
)

// InstanceInfo represents simplified instance info for API response.
type InstanceInfo struct {
	NodeID      string            `json:"node_id"`
	Address     string            `json:"address"`
	Cluster     string            `json:"cluster"`
	Version     string            `json:"version"`
	Metadata    map[string]string `json:"metadata"`
	Status      string            `json:"status"`
	LastSeen    string            `json:"last_seen"`
	ConnectedAt string            `json:"connected_at"`
	Uptime      string            `json:"uptime"`
}
