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
		if x == m {
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

// compiledHeaderSlicePool is a pool for temporary []CompiledHeader slices during snapshot building
var compiledHeaderSlicePool = sync.Pool{
	New: func() any {
		s := make([]CompiledHeader, 0, 4) // start with small capacity, grow as needed
		return &s
	},
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

	for _, r := range routes {
		// A) header-only：with Headers, without Path / Prefix
		if r.Match.Path == "" && r.Match.Prefix == "" && len(r.Match.Headers) > 0 {
			hr := HeaderRoute{
				Methods: r.Match.Methods,
				Action:  r.Route,
			}

			// use temporary slice from pool to build compiled headers
			chPtr := compiledHeaderSlicePool.Get().(*[]CompiledHeader)
			ch := (*chPtr)[:0] // reset

			for _, h := range r.Match.Headers {
				c := CompiledHeader{Name: h.Name}
				if h.Regex {
					// 1) the model already has compiled regex (if any) → use it directly
					if h.valueRE != nil {
						c.Regex = h.valueRE
					} else if len(h.Values) > 0 && h.Values[0] != "" {
						// 2) else use global cache/compile (cross-snapshot reuse)
						if re := getRegexpWithCache(h.Values[0]); re != nil {
							c.Regex = re
						} else {
							// invalid regex → skip this header matcher
							logger.Errorf("Header regex compiled fail for %v", h.Values[0])
							continue
						}
					}
				} else {
					// not regex → copy values directly (if any)
					if len(h.Values) > 0 {
						// direct assignment is ok here (string slice)
						c.Values = append(c.Values, h.Values...)
					}
				}
				ch = append(ch, c)
			}

			// move the temporary slice content to snapshot (ownership transferred)
			hr.Headers = make([]CompiledHeader, len(ch))
			copy(hr.Headers, ch)

			// reset and put back the temporary slice to pool
			*chPtr = (*chPtr)[:0]
			compiledHeaderSlicePool.Put(chPtr)

			s.HeaderOnly = append(s.HeaderOnly, hr)
			continue
		}

		// B) Trie
		methods := r.Match.Methods
		if len(methods) == 0 {
			methods = constMethods
		}
		for _, m := range methods {
			t := getTrie(m)
			t.Put(stringutil.GetTrieKeyWithPrefix(m, r.Match.Path, r.Match.Prefix, r.Match.Prefix != ""), r.Route)
		}
	}
	return s
}
