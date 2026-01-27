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

import (
	"runtime"
	"sync"
	"time"
)

var (
	monitorOnce   sync.Once
	globalMonitor *LoadMonitor
)

type LoadMonitor struct {
	window time.Duration
	last   time.Time
	count  int64
	rate   float64
	mutex  sync.Mutex
}

func NewLoadMonitor() *LoadMonitor {
	monitorOnce.Do(func() {
		globalMonitor = &LoadMonitor{
			window: time.Second,
			last:   time.Now(),
		}
	})
	return globalMonitor
}

func (lm *LoadMonitor) RecordRequest() {
	if lm == nil {
		return
	}
	lm.mutex.Lock()
	lm.count++
	lm.mutex.Unlock()
}

func (lm *LoadMonitor) Snapshot() LoadMetrics {
	if lm == nil {
		return LoadMetrics{}
	}
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	now := time.Now()
	elapsed := now.Sub(lm.last)
	if elapsed >= lm.window && elapsed > 0 {
		lm.rate = float64(lm.count) / elapsed.Seconds()
		lm.count = 0
		lm.last = now
	}
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	memUsage := 0.0
	if ms.Sys > 0 {
		memUsage = float64(ms.Alloc) / float64(ms.Sys)
	}
	return LoadMetrics{
		CPUUsage:    0,
		MemoryUsage: memUsage,
		RequestRate: lm.rate,
	}
}
