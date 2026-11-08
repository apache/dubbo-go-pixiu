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

// Package snapshotopt carries the opt-in token for snapshot load-balancer fast
// paths. It lives under internal/ on purpose: only in-tree balancers (packages
// rooted at pkg/cluster/loadbalancer) can import it and therefore name Token,
// so only trusted balancers can declare the SnapshotOptIn method that hands
// these flags to the runtime. External plugins cannot construct a Token, so
// they cannot opt into zero-copy or healthy-only access and always receive
// defensively copied snapshot endpoints.
package snapshotopt

// Token tells the snapshot pick path which fast paths a trusted balancer opts
// into. The zero value (both false) is the safe default: full snapshot,
// defensively copied.
type Token struct {
	// ZeroCopy is true when the balancer never mutates or retains snapshot
	// endpoints, so the runtime may hand it the snapshot-owned slices directly
	// instead of defensive copies.
	ZeroCopy bool
	// HealthyOnly is true when the balancer never reads PickContext.AllEndpoints,
	// so the runtime can skip populating the full (including unhealthy) set.
	HealthyOnly bool
}
