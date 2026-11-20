/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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

package initialize

import (
	"testing"
)

func TestClusterDraftAndLegacyRoutesRegistered(t *testing.T) {
	want := map[string]bool{
		"GET /config/api/cluster/schema":       false,
		"GET /config/api/cluster/draft/list":   false,
		"GET /config/api/cluster/draft/detail": false,
		"POST /config/api/cluster/draft":       false,
		"PUT /config/api/cluster/draft":        false,
		"POST /config/api/cluster/validate":    false,
		"POST /config/api/cluster/preview":     false,
		"GET /config/api/cluster/status":       false,
		"GET /config/api/cluster/diff":         false,
		"PUT /config/api/cluster/publish":      false,
		"GET /config/api/cluster/history":      false,
		"PUT /config/api/cluster/rollback":     false,
		"DELETE /config/api/cluster/draft":     false,
		"PUT /config/api/cluster":              false,
		"POST /config/api/cluster":             false,
		"DELETE /config/api/cluster":           false,
	}
	for _, route := range Routers().Routes() {
		key := route.Method + " " + route.Path
		if _, exists := want[key]; exists {
			want[key] = true
		}
	}
	for key, registered := range want {
		if !registered {
			t.Errorf("route %s is not registered", key)
		}
	}
}
