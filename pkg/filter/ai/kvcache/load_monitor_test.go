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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestNewLoadMonitor(t *testing.T) {
	lm := NewLoadMonitor()

	assert.NotNil(t, lm)
	assert.NotZero(t, lm.window)
	assert.NotZero(t, lm.last)
}

func TestLoadMonitor_RecordRequestNil(t *testing.T) {
	var lm *LoadMonitor
	// Should not panic
	lm.RecordRequest()
}

func TestLoadMonitor_RecordRequest(t *testing.T) {
	lm := NewLoadMonitor()

	initialCount := lm.count
	lm.RecordRequest()

	assert.Equal(t, initialCount+1, lm.count)
}

func TestLoadMonitor_GetMetrics(t *testing.T) {
	lm := NewLoadMonitor()

	// Record some requests
	lm.RecordRequest()
	lm.RecordRequest()
	lm.RecordRequest()

	metrics := lm.Snapshot()

	assert.NotNil(t, metrics)
	assert.GreaterOrEqual(t, metrics.RequestRate, 0.0)
	assert.GreaterOrEqual(t, metrics.MemoryUsage, 0.0)
	assert.Equal(t, 0.0, metrics.CPUUsage) // CPUUsage is hardcoded to 0
}

func TestLoadMonitor_MultipleInstances(t *testing.T) {
	// Verify that each call to NewLoadMonitor creates a new instance
	lm1 := NewLoadMonitor()
	lm2 := NewLoadMonitor()

	assert.NotSame(t, lm1, lm2)

	// Record requests on lm1
	lm1.RecordRequest()
	lm1.RecordRequest()

	// lm2 should have independent count
	assert.Equal(t, int64(2), lm1.count)
	assert.Equal(t, int64(0), lm2.count)
}
