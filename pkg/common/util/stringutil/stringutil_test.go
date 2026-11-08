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

package stringutil_test

import (
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/trie"
	"github.com/apache/dubbo-go-pixiu/pkg/common/util/stringutil"
)

func TestGetTrieKeyURLNormalization(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"empty", "", "GET"},
		{"root", "/", "GET"},
		{"relative path", "api/users", "GET/api/users"},
		{"trailing slash", "/api/users/", "GET/api/users"},
		{"absolute path", "http://example.com/api/users", "GET/api/users"},
		{"relative embedded scheme", "/proxy/http://example.com", "GET/proxy/http://example.com"},
		{"relative embedded scheme path", "/proxy/http://example.com/a", "GET/proxy/http://example.com/a"},
		{"relative embedded scheme without leading slash", "proxy/http://example.com", "GET/proxy/http://example.com"},
		{"relative embedded scheme query", "/proxy/http://example.com?next=/callback", "GET/proxy/http://example.com"},
		{"absolute URL with embedded scheme path", "https://gateway.example/proxy/http://example.com/a", "GET/proxy/http://example.com/a"},
		{"uppercase scheme", "HTTP://example.com/api/users", "GET/api/users"},
		{"scheme punctuation", "git+ssh://example.com/api/users", "GET/api/users"},
		{"invalid leading scheme", "1http://example.com", "GET/1http://example.com"},
		{"relative escaped path with scheme", "/proxy/http://example.com/a%2Fb", "GET/proxy/http://example.com/a%2Fb"},
		{"absolute root", "http://example.com", "GET"},
		{"https root", "https://example.com", "GET"},
		{"ipv6 root", "http://[::1]:8080", "GET"},
		{"absolute root query", "https://example.com?name=value", "GET"},
		{"absolute query slash", "https://example.com?next=/callback", "GET"},
		{"root slash query", "https://example.com/?next=/callback", "GET"},
		{"path query", "/api/users?name=value", "GET/api/users"},
		{"query URL", "/api/users?next=https://example.com/callback", "GET/api/users"},
		{"path slash query", "/api/users/?name=value", "GET/api/users"},
		{"escaped slash", "https://example.com/api/a%2Fb", "GET/api/a%2Fb"},
		{"path parameter", "/api/:id", "GET/api/:id"},
		{"wildcard", "/api/**", "GET/api/**"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := stringutil.GetTrieKey("GET", tt.path); got != tt.want {
				t.Fatalf("GetTrieKey(GET, %q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestGetTrieKeyAbsoluteRootMatchesRootRoute(t *testing.T) {
	routes := trie.NewTrie()
	if ok, err := routes.Put(stringutil.GetTrieKey("GET", "/"), "root"); !ok || err != nil {
		t.Fatalf("Put root = (%v, %v)", ok, err)
	}
	if ok, err := routes.Put(stringutil.GetTrieKey("GET", "/callback"), "callback"); !ok || err != nil {
		t.Fatalf("Put callback = (%v, %v)", ok, err)
	}
	for _, path := range []string{"https://example.com", "https://example.com?next=/callback"} {
		t.Run(path, func(t *testing.T) {
			node, _, ok := routes.Match(stringutil.GetTrieKey("GET", path))
			if !ok || node == nil || node.GetBizInfo() != "root" {
				t.Fatalf("absolute root %q matched (%v, %v), want root", path, node, ok)
			}
		})
	}
}

func TestGetTrieKeyEmbeddedSchemeDoesNotMatchRootRoute(t *testing.T) {
	routes := trie.NewTrie()
	if ok, err := routes.Put("GET", "root"); !ok || err != nil {
		t.Fatalf("Put root = (%v, %v)", ok, err)
	}
	for _, path := range []string{"/proxy/http://example.com", "proxy/http://example.com", "/proxy/http://example.com/a"} {
		t.Run(path, func(t *testing.T) {
			if node, _, ok := routes.Match(stringutil.GetTrieKey("GET", path)); ok {
				t.Fatalf("relative path %q unexpectedly matched %v", path, node.GetBizInfo())
			}
		})
	}
}
