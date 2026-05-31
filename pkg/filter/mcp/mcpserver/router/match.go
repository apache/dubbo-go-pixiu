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

package router

import (
	"fmt"
	"regexp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// matcher evaluates a PolicyMatch against a SelectionContext. Regex is
// precompiled once at construction so matching is allocation-free at runtime.
type matcher struct {
	match model.PolicyMatch
	re    *regexp.Regexp
}

// newMatcher compiles a PolicyMatch. It returns an error only for an invalid
// regex, so configuration mistakes surface at startup rather than per-request.
func newMatcher(m model.PolicyMatch) (*matcher, error) {
	if err := validatePolicyMatch(m); err != nil {
		return nil, err
	}
	mt := &matcher{match: m}
	if m.Regex != "" {
		re, err := regexp.Compile(m.Regex)
		if err != nil {
			return nil, fmt.Errorf("invalid regex %q: %w", m.Regex, err)
		}
		mt.re = re
	}
	return mt, nil
}

func validatePolicyMatch(m model.PolicyMatch) error {
	operators := 0
	if m.Equals != "" {
		operators++
	}
	if len(m.In) > 0 {
		operators++
	}
	if m.Regex != "" {
		operators++
	}

	if m.MissingClaim != "" {
		if m.Claim != "" || operators > 0 {
			return fmt.Errorf("missing_claim cannot be combined with claim, equals, in, or regex")
		}
		return nil
	}

	if operators > 1 {
		return fmt.Errorf("only one of equals, in, or regex may be configured")
	}
	if operators > 0 && m.Claim == "" {
		return fmt.Errorf("claim is required when equals, in, or regex is configured")
	}
	return nil
}

// alwaysMatches reports whether the match clause is empty (applies to all).
func (m *matcher) alwaysMatches() bool {
	return m.match.Claim == "" && m.match.MissingClaim == ""
}

// matches evaluates the clause against the context's claims.
//
// Semantics:
//   - empty clause            -> always true
//   - missing_claim           -> true when that claim is absent/empty
//   - claim + equals          -> true when claim string-equals value
//   - claim + in              -> true when claim is one of the values
//   - claim + regex           -> true when claim matches the pattern
//   - claim only (no op)      -> true when the claim is present and non-empty
func (m *matcher) matches(sc SelectionContext) bool {
	if m.match.MissingClaim != "" {
		return !hasClaim(sc, m.match.MissingClaim)
	}
	if m.match.Claim == "" {
		return true
	}

	val, ok := claimString(sc, m.match.Claim)
	if !ok {
		return false
	}

	switch {
	case m.match.Equals != "":
		return val == m.match.Equals
	case len(m.match.In) > 0:
		for _, candidate := range m.match.In {
			if val == candidate {
				return true
			}
		}
		return false
	case m.re != nil:
		return m.re.MatchString(val)
	default:
		// Claim specified with no operator: presence check.
		return val != ""
	}
}

// hasClaim reports whether a non-empty claim of the given key exists. It also
// recognizes the well-known synthetic claims sub/tenant promoted onto the
// SelectionContext so policies can reference them uniformly.
func hasClaim(sc SelectionContext, key string) bool {
	v, ok := claimString(sc, key)
	return ok && v != ""
}

// claimString resolves a claim to its string form, consulting both the raw
// claims map and the promoted UserID/Tenant fields.
func claimString(sc SelectionContext, key string) (string, bool) {
	switch key {
	case "sub":
		if sc.UserID != "" {
			return sc.UserID, true
		}
	case "tenant":
		if sc.Tenant != "" {
			return sc.Tenant, true
		}
	}
	if sc.Claims != nil {
		if raw, ok := sc.Claims[key]; ok {
			return fmt.Sprintf("%v", raw), true
		}
	}
	return "", false
}
