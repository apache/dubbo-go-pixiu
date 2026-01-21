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

package kvcache

import (
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

type FilterFactory struct {
	cfg   *Config
	store *RedisStore
}

func (f *FilterFactory) Config() any {
	return f.cfg
}

func (f *FilterFactory) Apply() error {
	if f.cfg == nil {
		f.cfg = &Config{}
	}
	if f.cfg.DefaultMode == "" {
		f.cfg.DefaultMode = ModeReuse
	}
	if f.cfg.DefaultUseCache == nil {
		value := true
		f.cfg.DefaultUseCache = &value
	}
	if f.cfg.DropClientHeaders == nil {
		value := true
		f.cfg.DropClientHeaders = &value
	}
	if f.cfg.Enabled && strings.EqualFold(f.cfg.Backend, "redis") {
		f.store = NewRedisStore(f.cfg.Redis, f.cfg.DefaultTTL)
		if f.store == nil {
			logger.Warnf("[dubbo-go-pixiu] kvcache filter enabled but redis config invalid")
		}
	}
	return nil
}

func (f *FilterFactory) PrepareFilterChain(_ *contexthttp.HttpContext, chain filter.FilterChain) error {
	cfgCopy := *f.cfg
	filterInstance := &Filter{
		cfg:   &cfgCopy,
		store: f.store,
	}
	chain.AppendDecodeFilters(filterInstance)
	chain.AppendEncodeFilters(filterInstance)
	return nil
}
