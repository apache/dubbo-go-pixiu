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

package shutdown

import (
	"sync"
)

// hooks belongs to one process lifetime. All hooks must be registered before
// signal handling starts. Running them more than once waits for the first run.
type hooks struct {
	mu        sync.Mutex
	once      sync.Once
	callbacks []func()
}

func (h *hooks) register(fn func()) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.callbacks = append(h.callbacks, fn)
}

func (h *hooks) run() {
	h.once.Do(func() {
		h.mu.Lock()
		callbacks := h.callbacks
		h.callbacks = nil
		h.mu.Unlock()
		for i := len(callbacks) - 1; i >= 0; i-- {
			callbacks[i]()
		}
	})
}

var processHooks hooks

// RegisterHook registers cleanup before the gateway starts accepting signals.
// Hooks run in reverse registration order so dependents close before owners.
func RegisterHook(fn func()) {
	processHooks.register(fn)
}

// RunHooks stops background resource producers before listeners drain. The
// signal handler's overall timeout remains active and can interrupt cleanup.
func RunHooks() {
	processHooks.run()
}
