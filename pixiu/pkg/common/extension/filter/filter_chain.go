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

package filter

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
)

// FilterChain
type FilterChain interface {
	AppendDecodeFilters(f ...HttpDecodeFilter)
	AppendEncodeFilters(f ...HttpEncodeFilter)

	OnDecode(ctx *http.HttpContext)
	OnEncode(ctx *http.HttpContext)
}

type defaultFilterChain struct {
	decodeFilters      []HttpDecodeFilter
	decodeFiltersIndex int

	encodeFilters      []HttpEncodeFilter
	encodeFiltersIndex int
}

func NewDefaultFilterChain() FilterChain {
	return &defaultFilterChain{
		decodeFilters:      []HttpDecodeFilter{},
		decodeFiltersIndex: 0,
		encodeFilters:      []HttpEncodeFilter{},
		encodeFiltersIndex: 0,
	}
}

func (c *defaultFilterChain) AppendDecodeFilters(f ...HttpDecodeFilter) {
	c.decodeFilters = append(c.decodeFilters, f...)
}

// AppendEncodeFilters append encode filters in reverse order
func (c *defaultFilterChain) AppendEncodeFilters(f ...HttpEncodeFilter) {
	for i := len(f) - 1; i >= 0; i-- {
		c.encodeFilters = append([]HttpEncodeFilter{f[i]}, c.encodeFilters...)
	}
}

func (c *defaultFilterChain) OnDecode(ctx *http.HttpContext) {
	for ; c.decodeFiltersIndex < len(c.decodeFilters); c.decodeFiltersIndex++ {
		filterStatus := c.decodeFilters[c.decodeFiltersIndex].Decode(ctx)

		switch filterStatus {
		case Continue:
			continue
		case Stop:
			return
		}
	}
}

func (c *defaultFilterChain) OnEncode(ctx *http.HttpContext) {
	for ; c.encodeFiltersIndex < len(c.encodeFilters); c.encodeFiltersIndex++ {
		filterStatus := c.encodeFilters[c.encodeFiltersIndex].Encode(ctx)

		switch filterStatus {
		case Continue:
			continue
		case Stop:
			return
		}
	}
}
