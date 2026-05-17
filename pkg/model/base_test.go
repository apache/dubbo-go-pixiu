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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestSocketAddressEqual(t *testing.T) {
	tests := []struct {
		name              string
		a                 SocketAddress
		b                 SocketAddress
		want              bool
		matchesGetAddress bool
	}{
		{
			name:              "both zero",
			a:                 SocketAddress{},
			b:                 SocketAddress{},
			want:              true,
			matchesGetAddress: true,
		},
		{
			name:              "ip and port equal",
			a:                 SocketAddress{Address: "127.0.0.1", Port: 8080},
			b:                 SocketAddress{Address: "127.0.0.1", Port: 8080},
			want:              true,
			matchesGetAddress: true,
		},
		{
			name:              "same ip different port",
			a:                 SocketAddress{Address: "127.0.0.1", Port: 8080},
			b:                 SocketAddress{Address: "127.0.0.1", Port: 8081},
			want:              false,
			matchesGetAddress: true,
		},
		{
			name:              "same single domain",
			a:                 SocketAddress{Domains: []string{"svc.example.com"}},
			b:                 SocketAddress{Domains: []string{"svc.example.com"}},
			want:              true,
			matchesGetAddress: true,
		},
		{
			name:              "different domains",
			a:                 SocketAddress{Domains: []string{"svc.a.example.com"}},
			b:                 SocketAddress{Domains: []string{"svc.b.example.com"}},
			want:              false,
			matchesGetAddress: true,
		},
		{
			name:              "one has domain other has ip port",
			a:                 SocketAddress{Domains: []string{"svc.example.com"}},
			b:                 SocketAddress{Address: "127.0.0.1", Port: 8080},
			want:              false,
			matchesGetAddress: true,
		},
		{
			name:              "both have single empty domain",
			a:                 SocketAddress{Domains: []string{""}},
			b:                 SocketAddress{Domains: []string{""}},
			want:              true,
			matchesGetAddress: true,
		},
		{
			name:              "ipv4 vs ipv6 formatting",
			a:                 SocketAddress{Address: "127.0.0.1", Port: 8080},
			b:                 SocketAddress{Address: "::1", Port: 8080},
			want:              false,
			matchesGetAddress: true,
		},
		// Counterexamples: Equal != (GetAddress == GetAddress). The
		// matchesGetAddress=false flag tells the assertion below to
		// not enforce string-projection agreement and instead encode
		// the divergence into the contract.
		{
			name: "domain literal collides with ip+port string but Equal is false",
			// "1.2.3.4:5" as a domain entry literally formats to the
			// same string as the IP+port form. They describe different
			// upstreams (domain mode vs IP+port mode) so Equal returns
			// false even though GetAddress strings match.
			a:                 SocketAddress{Domains: []string{"1.2.3.4:5"}},
			b:                 SocketAddress{Address: "1.2.3.4", Port: 5},
			want:              false,
			matchesGetAddress: false,
		},
		{
			name: "domain entry shaped like ip alone collides with zero-port ip form",
			// "1.2.3.4" as a domain entry formats to "1.2.3.4"; IP-only
			// SocketAddress{Address: "1.2.3.4", Port: 0} formats to
			// "1.2.3.4:0". Both string forms differ — picked another
			// pair where the strings actually collide.
			a:                 SocketAddress{Domains: []string{":0"}},
			b:                 SocketAddress{}, // Address="", Port=0 -> ":0"
			want:              false,
			matchesGetAddress: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.Equal(tc.b)
			assert.Equal(t, tc.want, got)

			assert.Equal(t, got, tc.b.Equal(tc.a))

			if tc.matchesGetAddress {
				assert.Equal(t, tc.a.GetAddress() == tc.b.GetAddress(), got,
					"Equal must agree with GetAddress() comparison for case %q", tc.name)
			} else {
				// Counterexample cases: GetAddress strings match yet
				// Equal is false. Encoded as part of the contract — see
				// godoc on Equal / GetAddress.
				assert.NotEqual(t, tc.a.GetAddress() == tc.b.GetAddress(), got,
					"counterexample %q: Equal must diverge from GetAddress here", tc.name)
			}
		})
	}
}
