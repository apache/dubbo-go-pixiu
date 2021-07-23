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

package http

import (
	"regexp"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// HttpHeaderMatch
func HttpHeaderMatch(c *HttpContext, hm model.HeaderMatcher) bool {
	if hm.Name == "" {
		return true
	}

	if hm.Value == "" {
		if c.GetHeader(hm.Name) == "" {
			return true
		}
	} else {
		if hm.Regex {
			// TODO
			return true
		} else {
			if c.GetHeader(hm.Name) == hm.Value {
				return true
			}
		}
	}

	return false
}

// HttpRouteMatch
func HttpRouteMatch(c *HttpContext, rm model.RouterMatch) bool {
	if rm.Prefix != "" {
		if !strings.HasPrefix(c.GetUrl(), rm.Path) {
			return false
		}
	}

	if rm.Path != "" {
		if c.GetUrl() != rm.Path {
			return false
		}
	}

	if rm.Regex != "" {
		if !regexp.MustCompile(rm.Regex).MatchString(c.GetUrl()) {
			return false
		}
	}

	return true
}

// HttpRouteActionMatch
func HttpRouteActionMatch(c *HttpContext, ra model.RouteAction) bool {
	conf := config.GetBootstrap()

	if ra.Cluster == "" || !conf.ExistCluster(ra.Cluster) {
		return false
	}

	return true
}
