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

package matcher

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config/ratelimit"
)

var _matcher *Matcher

// Init matcher
func Init() {
	_matcher = NewMatcher()
}

//PathMatcher according the url path find APIResource name
type PathMatcher interface {
	load(apis []ratelimit.Resource)

	match(path string) (string, bool)
}

type Matcher struct {
	matchers []PathMatcher
}

func NewMatcher() *Matcher {
	return &Matcher{
		matchers: []PathMatcher{
			&Exact{},
			&Regex{},
		},
	}
}

// Load load api resource for matchers
func Load(apis []ratelimit.Resource) {
	for _, v := range _matcher.matchers {
		v.load(apis)
	}
}

// Match match resource via url path
func Match(path string) (string, bool) {
	for _, m := range _matcher.matchers {
		if res, ok := m.match(path); ok {
			return res, ok
		}
	}
	return "", false
}
