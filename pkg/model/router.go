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
	http2 "net/http"
	"regexp"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/util/stringutil"
)

// Router struct

type (
	Router struct {
		ID    string      `yaml:"id" json:"id" mapstructure:"id"`
		Match RouterMatch `yaml:"match" json:"match" mapstructure:"match"`
		Route RouteAction `yaml:"route" json:"route" mapstructure:"route"`
	}

	// RouterMatch
	RouterMatch struct {
		Prefix  string `yaml:"prefix" json:"prefix" mapstructure:"prefix"`
		Path    string `yaml:"path" json:"path" mapstructure:"path"`
		Regex   string `yaml:"regex" json:"regex" mapstructure:"regex"`
		PathRE  *regexp.Regexp
		Methods []string        `yaml:"methods" json:"methods" mapstructure:"methods"`
		Headers []HeaderMatcher `yaml:"headers" json:"headers" mapstructure:"headers"`
	}

	// RouteAction match route should do
	RouteAction struct {
		Cluster                     string            `yaml:"cluster" json:"cluster" mapstructure:"cluster"`
		ClusterNotFoundResponseCode int               `yaml:"cluster_not_found_response_code" json:"cluster_not_found_response_code" mapstructure:"cluster_not_found_response_code"`
		PrefixRewrite               string            `yaml:"prefix_rewrite" json:"prefix_rewrite" mapstructure:"prefix_rewrite"`
		HostRewrite                 string            `yaml:"host_rewrite" json:"host_rewrite" mapstructure:"host_rewrite"`
		Timeout                     string            `yaml:"timeout" json:"timeout" mapstructure:"timeout"`
		Priority                    int8              `yaml:"priority" json:"priority" mapstructure:"priority"`
		ResponseHeadersToAdd        HeaderValueOption `yaml:"response_headers_to_add" json:"response_headers_to_add" mapstructure:"response_headers_to_add"`          // ResponseHeadersToAdd add response head
		ResponseHeadersToRemove     []string          `yaml:"response_headers_to_remove" json:"response_headers_to_remove" mapstructure:"response_headers_to_remove"` // ResponseHeadersToRemove remove response head
		RequestHeadersToAdd         HeaderValueOption `yaml:"request_headers_to_add" json:"request_headers_to_add" mapstructure:"request_headers_to_add"`             // RequestHeadersToAdd add request head
	}

	// RouteConfiguration
	RouteConfiguration struct {
		Routes  []Router `yaml:"routes" json:"routes" mapstructure:"routes"`
		Dynamic bool     `yaml:"dynamic" json:"dynamic" mapstructure:"dynamic"`
	}

	// Name header key, Value header value, Regex header value is regex
	HeaderMatcher struct {
		Name    string   `yaml:"name" json:"name" mapstructure:"name"`
		Values  []string `yaml:"values" json:"values" mapstructure:"values"`
		Regex   bool     `yaml:"regex" json:"regex" mapstructure:"regex"`
		ValueRE *regexp.Regexp
	}
)

func (rc *RouteConfiguration) Route(req *http2.Request) (*RouteAction, error) {
	if rc.Routes == nil {
		return nil, errors.Errorf("router configuration is empty")
	}

	for _, r := range rc.Routes {
		if r.MatchRouter(req) {
			return &r.Route, nil
		}
	}

	return nil, errors.Errorf("no matched route")
}

func (r *Router) MatchRouter(req *http2.Request) bool {
	if r.Match.matchPath(req) {
		return true
	}

	if r.Match.matchMethod(req) {
		return true
	}

	if r.Match.matchHeader(req) {
		return true
	}

	return false
}

func (rm *RouterMatch) matchPath(req *http2.Request) bool {
	if rm.Path == "" && rm.Prefix == "" && rm.PathRE == nil {
		return true
	}

	path := req.RequestURI

	if rm.Path != "" && rm.Path == path {
		return true
	}
	if rm.Prefix != "" && strings.HasPrefix(path, rm.Prefix) {
		return true
	}
	if rm.PathRE != nil {
		return rm.PathRE.MatchString(path)
	}

	return false
}

func (rm *RouterMatch) matchMethod(req *http2.Request) bool {
	if len(rm.Methods) == 0 {
		return true
	}

	return stringutil.StrInSlice(req.Method, rm.Methods)
}

func (rm *RouterMatch) matchHeader(req *http2.Request) bool {

	for _, h := range rm.Headers {
		v := req.Header.Get(h.Name)
		if stringutil.StrInSlice(v, h.Values) {
			return true
		}

		if h.ValueRE != nil && h.ValueRE.MatchString(v) {
			return true
		}
	}

	return false
}
