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
	runtimemetrics "runtime/metrics"
	"sync"
	"time"
)

type LoadMonitor struct {
	window time.Duration
	last   time.Time
	count  int64
	rate   float64
	mutex  sync.Mutex

	lastCPUSampleAt time.Time
	lastCPUSeconds  float64
	hasCPUSample    bool
}

func NewLoadMonitor() *LoadMonitor {
	return &LoadMonitor{
		window: time.Second,
		last:   time.Now(),
	}
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

	cpuUsage := lm.sampleCPUUsage()
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
		CPUUsage:    cpuUsage,
		MemoryUsage: memUsage,
		RequestRate: lm.rate,
	}
}

func (lm *LoadMonitor) sampleCPUUsage() float64 {
	totalCPUSeconds, gomaxprocs, ok := readRuntimeCPUStats()
	if !ok {
		return 0
	}

	now := time.Now()
	if !lm.hasCPUSample {
		lm.lastCPUSeconds = totalCPUSeconds
		lm.lastCPUSampleAt = now
		lm.hasCPUSample = true
		return 0
	}

	wall := now.Sub(lm.lastCPUSampleAt).Seconds()
	if wall <= 0 || gomaxprocs <= 0 {
		return 0
	}

	cpuDelta := totalCPUSeconds - lm.lastCPUSeconds
	lm.lastCPUSeconds = totalCPUSeconds
	lm.lastCPUSampleAt = now
	if cpuDelta <= 0 {
		return 0
	}

	usage := cpuDelta / (wall * gomaxprocs)
	if usage < 0 {
		return 0
	}
	if usage > 1 {
		return 1
	}
	return usage
}

func readRuntimeCPUStats() (totalCPUSeconds float64, gomaxprocs float64, ok bool) {
	samples := []runtimemetrics.Sample{
		{Name: "/cpu/classes/total:cpu-seconds"},
		{Name: "/sched/gomaxprocs:threads"},
	}
	runtimemetrics.Read(samples)

	totalValue := samples[0].Value
	gomaxValue := samples[1].Value
	if totalValue.Kind() != runtimemetrics.KindFloat64 || gomaxValue.Kind() != runtimemetrics.KindUint64 {
		return 0, 0, false
	}
	return totalValue.Float64(), float64(gomaxValue.Uint64()), true
}
