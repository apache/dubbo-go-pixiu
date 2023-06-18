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

package v1alpha1

import (
	"fmt"
	"time"
)

import (
	"go.uber.org/atomic"
)

/**
Copy From pilot/pkg/xds/discovery.go
*/
type debounceOptions struct {
	// debounceAfter is the delay added to events to wait
	// after a registry/config event for debouncing.
	// This will delay the push by at least this interval, plus
	// the time getting subsequent events. If no change is
	// detected the push will happen, otherwise we'll keep
	// delaying until things settle.
	debounceAfter time.Duration

	// debounceMax is the maximum time to wait for events
	// while debouncing. Defaults to 10 seconds. If events keep
	// showing up with no break for this time, we'll trigger a push.
	debounceMax time.Duration

	// enableDebounce indicates whether EDS pushes should be debounced.
	enableDebounce bool
}

type DebounceHelper struct {
}

func (h *DebounceHelper) Debounce(ch chan *pushRequest, stopCh <-chan struct{}, opts debounceOptions, pushFn func(req *pushRequest), updateSent *atomic.Int64) {
	defer func() {
		err := recover()
		if err != nil {
			log.Infof("Debounce panic caused by: {%+v}", err)
		}
	}()

	var timeChan <-chan time.Time
	var startDebounce time.Time
	var lastConfigUpdateTime time.Time

	pushCounter := 0
	debouncedEvents := 0

	// Keeps track of the push requests. If updates are debounce they will be merged.
	var req *pushRequest

	free := true
	freeCh := make(chan struct{}, 1)

	push := func(req *pushRequest, debouncedEvents int) {
		pushFn(req)
		updateSent.Add(int64(debouncedEvents))
		freeCh <- struct{}{}
	}

	pushWorker := func() {
		eventDelay := time.Since(startDebounce)
		quietTime := time.Since(lastConfigUpdateTime)
		// it has been too long or quiet enough
		if eventDelay >= opts.debounceMax || quietTime >= opts.debounceAfter {
			if req != nil {
				pushCounter++
				if req.ConfigsUpdated == nil {
					log.Infof("Push debounce stable[%d] %d : %v since last change, %v since last push",
						pushCounter, debouncedEvents,
						quietTime, eventDelay)
				} else {
					log.Infof("Push debounce stable[%d] %d for config %s: %v since last change, %v since last push",
						pushCounter, debouncedEvents, PushRequestConfigsUpdated(req),
						quietTime, eventDelay)
				}
				free = false
				go push(req, debouncedEvents)
				req = nil
				debouncedEvents = 0
			}
		} else {
			timeChan = time.After(opts.debounceAfter - quietTime)
		}
	}

	for {
		select {
		case <-freeCh:
			free = true
			pushWorker()
		case r := <-ch:
			if !opts.enableDebounce {
				// trigger push now, just for EDS
				go func(req *pushRequest) {
					pushFn(req)
					updateSent.Inc()
				}(r)
				continue
			}

			lastConfigUpdateTime = time.Now()
			if debouncedEvents == 0 {
				timeChan = time.After(opts.debounceAfter)
				startDebounce = lastConfigUpdateTime
			}
			debouncedEvents++

			req = req.Merge(r)
		case <-timeChan:
			if free {
				pushWorker()
			}
		case <-stopCh:
			return
		}
	}
}

func PushRequestConfigsUpdated(req *pushRequest) string {
	configs := ""
	for key := range req.ConfigsUpdated {
		configs += key.String()
		break
	}
	if len(req.ConfigsUpdated) > 1 {
		more := fmt.Sprintf(" and %d more configs", len(req.ConfigsUpdated)-1)
		configs += more
	}
	return configs
}
