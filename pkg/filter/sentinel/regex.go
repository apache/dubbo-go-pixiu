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

package sentinel

import (
	"regexp"
	"sync"
)

type Regex struct {
	apiNames map[string]string

	mu sync.RWMutex
}

func (p *Regex) load(apis []*Resource) {
	m := make(map[string]string, len(apis))

	for _, api := range apis {
		apiName := api.Name
		for _, item := range api.Items {
			if item.MatchStrategy == REGEX {
				m[item.Pattern] = apiName
			}
		}
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.apiNames = m
}

func (p *Regex) match(path string) (string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for k, v := range p.apiNames {
		matched, _ := regexp.MatchString(k, path)
		if matched {
			return v, true
		}
	}
	return "", false
}
