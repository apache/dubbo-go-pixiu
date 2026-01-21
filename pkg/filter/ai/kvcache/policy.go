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

import "strings"

func normalizeMode(mode string) string {
	switch strings.ToLower(strings.TrimSpace(mode)) {
	case ModeReuse:
		return ModeReuse
	case ModeRefresh:
		return ModeRefresh
	case ModeFlush:
		return ModeFlush
	default:
		return ""
	}
}

func defaultPolicy(cfg *Config) Policy {
	useCache := true
	if cfg.DefaultUseCache != nil {
		useCache = *cfg.DefaultUseCache
	}
	mode := normalizeMode(cfg.DefaultMode)
	if mode == "" {
		mode = ModeReuse
	}
	return Policy{
		UseCache: useCache,
		Mode:     mode,
		TTL:      cfg.DefaultTTL,
	}
}

func resolveMode(policy Policy) string {
	if !policy.UseCache {
		return ModeRefresh
	}
	mode := normalizeMode(policy.Mode)
	if mode == "" {
		return ModeReuse
	}
	return mode
}

func resolveCacheResult(value string) string {
	switch strings.ToUpper(strings.TrimSpace(value)) {
	case CacheHit:
		return CacheHit
	case CacheMiss:
		return CacheMiss
	case CacheError:
		return CacheError
	default:
		return ""
	}
}

func modelAllowed(model string, targets []string) bool {
	if len(targets) == 0 {
		return true
	}
	if model == "" {
		return false
	}
	for _, target := range targets {
		if model == target {
			return true
		}
	}
	return false
}
