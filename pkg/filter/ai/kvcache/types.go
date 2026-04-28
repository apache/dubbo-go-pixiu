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

package kvcache

type CacheStats struct {
	Size      int     `json:"size"`
	HitRate   float64 `json:"hit_rate"`
	HitCount  int64   `json:"hit_count"`
	MissCount int64   `json:"miss_count"`
}

type LoadMetrics struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	RequestRate float64 `json:"request_rate"`
}

type LookupResponse struct {
	EventID    string                 `json:"event_id"`
	LayoutInfo map[string]CacheLayout `json:"layout_info"`
}

type CacheLayout struct {
	Location   string `json:"0"`
	TokenCount int    `json:"1"`
}
