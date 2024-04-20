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
	stdHttp "net/http"
	"regexp"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/router/trie"
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
		Prefix string `yaml:"prefix" json:"prefix" mapstructure:"prefix"`
		Path   string `yaml:"path" json:"path" mapstructure:"path"`
		// Regex   string          `yaml:"regex" json:"regex" mapstructure:"regex"` TODO: next version
		Methods []string        `yaml:"methods" json:"methods" mapstructure:"methods"`
		Headers []HeaderMatcher `yaml:"headers,omitempty" json:"headers,omitempty" mapstructure:"headers"`
		// pathRE  *regexp.Regexp
	}

	// RouteAction match route should do
	RouteAction struct {
		Cluster                     string `yaml:"cluster" json:"cluster" mapstructure:"cluster"`
		ClusterNotFoundResponseCode int    `yaml:"cluster_not_found_response_code" json:"cluster_not_found_response_code" mapstructure:"cluster_not_found_response_code"`
	}

	// RouteConfiguration
	RouteConfiguration struct {
		RouteTrie trie.Trie `yaml:"-" json:"-" mapstructure:"-"`
		Routes    []*Router `yaml:"routes" json:"routes" mapstructure:"routes"`
		Dynamic   bool      `yaml:"dynamic" json:"dynamic" mapstructure:"dynamic"`
	}

	// HeaderMatcher include Name header key, Values header value, Regex regex value
	HeaderMatcher struct {
		Name    string   `yaml:"name" json:"name" mapstructure:"name"`
		Values  []string `yaml:"values" json:"values" mapstructure:"values"`
		Regex   bool     `yaml:"regex" json:"regex" mapstructure:"regex"`
		valueRE *regexp.Regexp
	}
)

func NewRouterMatchPrefix(name string) RouterMatch {
	return RouterMatch{Prefix: "/" + name + "/"}
}

func (rc *RouteConfiguration) RouteByPathAndMethod(path, method string) (*RouteAction, error) {
	if rc.RouteTrie.IsEmpty() {
		return nil, errors.Errorf("router configuration is empty")
	}

	node, _, _ := rc.RouteTrie.Match(stringutil.GetTrieKey(method, path))
	if node == nil {
		return nil, errors.Errorf("route failed for %s, no rules matched.", stringutil.GetTrieKey(method, path))
	}
	if node.GetBizInfo() == nil {
		return nil, errors.Errorf("action is nil. please check your configuration.")
	}
	ret := (node.GetBizInfo()).(RouteAction)

	return &ret, nil
}

func (rc *RouteConfiguration) Route(req *stdHttp.Request) (*RouteAction, error) {
	return rc.RouteByPathAndMethod(req.URL.Path, req.Method)
}

// MatchHeader used when there's only headers to match
func (rm *RouterMatch) MatchHeader(req *stdHttp.Request) bool {
	if len(rm.Methods) > 0 {
		for _, method := range rm.Methods {
			if method == req.Method {
				goto HEADER
			}
		}
		return false
	}
HEADER:
	for _, header := range rm.Headers {
		if val := req.Header.Get(header.Name); len(val) > 0 {
			if header.MatchValues(val) {
				return true
			}
		}
	}
	return false
}

// MatchValues match values in header, including regex type
func (hm *HeaderMatcher) MatchValues(dst string) bool {
	if hm.Regex && hm.valueRE != nil {
		return hm.valueRE.MatchString(dst)
	}

	for _, src := range hm.Values {
		if src == dst {
			return true
		}
	}
	return false
}

// SetValueRegex compile the regex, disable regex if it failed
func (hm *HeaderMatcher) SetValueRegex(regex string) error {
	r, err := regexp.Compile(regex)
	if err == nil {
		hm.valueRE = r
		return nil
	}
	hm.Regex = false
	return err
}
