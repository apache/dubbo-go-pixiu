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
	"github.com/apache/dubbo-go-pixiu/pkg/filter/ratelimit/config"
)

var matcher *Matcher

func Init() {
	matcher = NewMatcher()
}

//PathMatcher according the url path find APIResource name
type PathMatcher interface {
	load(apis []config.APIResource)

	match(path string) (string, bool)
}

type Matcher struct {
	matchers []PathMatcher
}

func NewMatcher() *Matcher {
	return &Matcher{
		matchers: []PathMatcher{
			&Accurate{},
			&Regex{},
		},
	}
}

func Load(apis []config.APIResource) {
	for _, v := range matcher.matchers {
		v.load(apis)
	}
}

func Match(path string) (string, bool) {
	for _, m := range matcher.matchers {
		if res, ok := m.match(path); ok {
			return res, ok
		}
	}
	return "", false
}
