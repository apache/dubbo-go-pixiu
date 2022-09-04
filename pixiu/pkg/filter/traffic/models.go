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

package traffic

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type splitAction func(header string, expectedValue string, r *http.Request) bool

type ActionWrapper struct {
	name   CanaryHeaders
	action splitAction
}

type Actions struct {
	wrappers []*ActionWrapper
}

var actions Actions

func initActions() {
	// canaryByHeader
	headerAction := &ActionWrapper{
		name: canaryByHeader,
		action: func(header string, value string, r *http.Request) bool {
			if v := r.Header.Get(header); v != "" {
				return v == value
			}
			return false
		},
	}

	// canaryWeight
	weightAction := &ActionWrapper{
		name: canaryWeight,
		action: func(header string, value string, r *http.Request) bool {
			rand.Seed(time.Now().UnixNano())
			num := rand.Intn(100) + 1
			intValue, _ := strconv.Atoi(value)
			return num <= intValue
		},
	}

	actions.wrappers = append(actions.wrappers, headerAction)
	actions.wrappers = append(actions.wrappers, weightAction)
}

func spilt(request *http.Request, rules map[CanaryHeaders]*Cluster) *Cluster {
	for _, wrapper := range actions.wrappers {
		// call action, decide to return cluster or not
		if cluster, exist := rules[wrapper.name]; exist {
			header, value := cluster.Canary, cluster.Value
			if wrapper.action(header, value, request) {
				return cluster
			}
		}
	}
	return nil
}
