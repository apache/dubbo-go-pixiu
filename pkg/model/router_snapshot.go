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
	"regexp"
	"strings"
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/trie"
	"github.com/apache/dubbo-go-pixiu/pkg/common/util/stringutil"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	constMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"}
)

// RouteSnapshot Read-only snapshot for routing
type RouteSnapshot struct {
	// multi-trie for each method, built once and read-only
	MethodTries map[string]*trie.Trie

	// precompiled regex for header-only routes
	HeaderOnly []HeaderRoute
}

type HeaderRoute struct {
	Methods []string
	Headers []CompiledHeader
	Action  RouteAction
}

type RouteCandidate struct {
	Headers []CompiledHeader
	Action  RouteAction
}

type RouteCandidates struct {
	First      RouteCandidate
	Additional []RouteCandidate
	Count      int
}

func (r *RouteCandidates) Add(candidate RouteCandidate) {
	if r.Count == 0 {
		r.First = candidate
		r.Count = 1
		return
	}
	r.Additional = append(r.Additional, candidate)
	r.Count++
}

func (r *RouteCandidates) ForEach(fn func(RouteCandidate) bool) {
	if r.Count == 0 || !fn(r.First) {
		return
	}
	for _, candidate := range r.Additional {
		if !fn(candidate) {
			return
		}
	}
}

type CompiledHeader struct {
	Name   string
	Regex  *regexp.Regexp
	Values []string
}

func MethodAllowed(methods []string, m string) bool {
	if len(methods) == 0 {
		return true
	}
	for _, x := range methods {
		// "*" matches any method (wildcard semantics, important for non-HTTP
		// scenarios like Dubbo/Triple RPC method matching).
		if x == "*" || x == m {
			return true
		}
	}
	return false
}

var regexCache sync.Map // map[string]*regexp.Regexp

func getRegexpWithCache(pat string) *regexp.Regexp {
	if v, ok := regexCache.Load(pat); ok {
		return v.(*regexp.Regexp)
	}
	// Compile fail return nil (caller will ignore this regex)
	re, err := regexp.Compile(pat)
	if err != nil {
		return nil
	}
	if v, ok := regexCache.LoadOrStore(pat, re); ok {
		return v.(*regexp.Regexp)
	}
	return re
}

func ToSnapshot(routes []*Router) *RouteSnapshot {
	s := &RouteSnapshot{
		MethodTries: make(map[string]*trie.Trie, 8),
	}

	// pre-scan header-only routes count
	headerOnlyCount := 0
	for _, r := range routes {
		if r.Match.Path == "" && r.Match.Prefix == "" && len(r.Match.Headers) > 0 {
			headerOnlyCount++
		}
	}

	if headerOnlyCount > 0 {
		s.HeaderOnly = make([]HeaderRoute, 0, headerOnlyCount)
	}

	// part to get or create trie for a method
	getTrie := func(m string) *trie.Trie {
		if t := s.MethodTries[m]; t != nil {
			return t
		}
		nt := trie.NewTrie()
		s.MethodTries[m] = &nt
		return &nt
	}
	routeCandidates := make(map[string]*RouteCandidates, len(routes))
	for _, r := range routes {
		// A) header-only：with Headers, without Path / Prefix
		if r.Match.Path == "" && r.Match.Prefix == "" && len(r.Match.Headers) > 0 {
			hr := HeaderRoute{
				Methods: r.Match.Methods,
				Action:  r.Route,
			}
			hr.Headers = compileHeaders(r.Match.Headers)
			s.HeaderOnly = append(s.HeaderOnly, hr)
			continue
		}

		// B) Trie
		methods := r.Match.Methods
		if len(methods) == 0 {
			methods = constMethods
		}
		headers := compileHeaders(r.Match.Headers)
		for _, m := range methods {
			t := getTrie(m)
			key := stringutil.GetTrieKeyWithPrefix(m, r.Match.Path, r.Match.Prefix, r.Match.Prefix != "")
			identity := trieRouteIdentity(key)
			candidates := routeCandidates[identity]
			if candidates == nil {
				candidates = &RouteCandidates{}
				routeCandidates[identity] = candidates
			}
			// Put every spelling so parameter and single-segment wildcard aliases
			// reach the same candidate list through the existing Trie node.
			_, _ = t.Put(key, candidates)
			candidates.Add(RouteCandidate{
				Headers: headers,
				Action:  r.Route,
			})
		}
	}
	return s
}

// trieRouteIdentity matches the trie's shared node for parameter and single-
// segment wildcard tokens without merging literals or match-all suffixes.
func trieRouteIdentity(key string) string {
	// Scan exact slash-delimited segments first so literal routes keep the
	// original key without allocating a normalized copy.
	firstSlash := strings.IndexByte(key, '/')
	if firstSlash < 0 {
		return key
	}
	for start := firstSlash + 1; start <= len(key); {
		end := strings.IndexByte(key[start:], '/')
		if end < 0 {
			end = len(key)
		} else {
			end += start
		}
		if stringutil.IsPathVariableOrWildcard(key[start:end]) {
			break
		}
		if end == len(key) {
			return key
		}
		start = end + 1
	}
	parts := strings.Split(key, "/")
	for i := 1; i < len(parts); i++ {
		if stringutil.IsPathVariableOrWildcard(parts[i]) {
			parts[i] = "*"
		}
	}
	return strings.Join(parts, "/")
}

func compileHeaders(headers []HeaderMatcher) []CompiledHeader {
	if len(headers) == 0 {
		return nil
	}
	valueCount := 0
	for _, h := range headers {
		if !h.Regex {
			valueCount += len(h.Values)
		}
	}
	compiled := make([]CompiledHeader, 0, len(headers))
	values := make([]string, valueCount)
	valueOffset := 0
	for _, h := range headers {
		c := CompiledHeader{Name: h.Name}
		if h.Regex {
			if h.valueRE != nil {
				c.Regex = h.valueRE
			} else if len(h.Values) > 0 && h.Values[0] != "" {
				if re := getRegexpWithCache(h.Values[0]); re != nil {
					c.Regex = re
				} else {
					logger.Errorf("Header regex compiled fail for %v", h.Values[0])
					continue
				}
			}
		} else if len(h.Values) > 0 {
			end := valueOffset + len(h.Values)
			copy(values[valueOffset:end], h.Values)
			c.Values = values[valueOffset:end]
			valueOffset = end
		}
		compiled = append(compiled, c)
	}
	return compiled
}
