// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"fmt"
	"strings"
)

import (
	rbacpb "github.com/envoyproxy/go-control-plane/envoy/config/rbac/v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/security/authz/matcher"
	sm "github.com/apache/dubbo-go-pixiu/pilot/pkg/security/model"
	"github.com/apache/dubbo-go-pixiu/pkg/spiffe"
)

type generator interface {
	permission(key, value string, forTCP bool) (*rbacpb.Permission, error)
	principal(key, value string, forTCP bool) (*rbacpb.Principal, error)
}

type destIPGenerator struct{}

func (destIPGenerator) permission(_, value string, _ bool) (*rbacpb.Permission, error) {
	cidrRange, err := matcher.CidrRange(value)
	if err != nil {
		return nil, err
	}
	return permissionDestinationIP(cidrRange), nil
}

func (destIPGenerator) principal(_, _ string, _ bool) (*rbacpb.Principal, error) {
	return nil, fmt.Errorf("unimplemented")
}

type destPortGenerator struct{}

func (destPortGenerator) permission(_, value string, _ bool) (*rbacpb.Permission, error) {
	portValue, err := convertToPort(value)
	if err != nil {
		return nil, err
	}
	return permissionDestinationPort(portValue), nil
}

func (destPortGenerator) principal(_, _ string, _ bool) (*rbacpb.Principal, error) {
	return nil, fmt.Errorf("unimplemented")
}

type connSNIGenerator struct{}

func (connSNIGenerator) permission(_, value string, _ bool) (*rbacpb.Permission, error) {
	m := matcher.StringMatcher(value)
	return permissionRequestedServerName(m), nil
}

func (connSNIGenerator) principal(_, _ string, _ bool) (*rbacpb.Principal, error) {
	return nil, fmt.Errorf("unimplemented")
}

type envoyFilterGenerator struct{}

func (envoyFilterGenerator) permission(key, value string, _ bool) (*rbacpb.Permission, error) {
	// Split key of format "experimental.envoy.filters.a.b[c]" to "envoy.filters.a.b" and "c".
	parts := strings.SplitN(strings.TrimSuffix(strings.TrimPrefix(key, "experimental."), "]"), "[", 2)

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid key: %v", key)
	}

	// If value is of format [v], create a list matcher.
	// Else, if value is of format v, create a string matcher.
	if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
		m := matcher.MetadataListMatcher(parts[0], parts[1:], matcher.StringMatcher(strings.Trim(value, "[]")))
		return permissionMetadata(m), nil
	}
	m := matcher.MetadataStringMatcher(parts[0], parts[1], matcher.StringMatcher(value))
	return permissionMetadata(m), nil
}

func (envoyFilterGenerator) principal(_, _ string, _ bool) (*rbacpb.Principal, error) {
	return nil, fmt.Errorf("unimplemented")
}

type srcIPGenerator struct{}

func (srcIPGenerator) permission(_, _ string, _ bool) (*rbacpb.Permission, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (srcIPGenerator) principal(_, value string, _ bool) (*rbacpb.Principal, error) {
	cidr, err := matcher.CidrRange(value)
	if err != nil {
		return nil, err
	}
	return principalDirectRemoteIP(cidr), nil
}

type remoteIPGenerator struct{}

func (remoteIPGenerator) permission(_, _ string, _ bool) (*rbacpb.Permission, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (remoteIPGenerator) principal(_, value string, _ bool) (*rbacpb.Principal, error) {
	cidr, err := matcher.CidrRange(value)
	if err != nil {
		return nil, err
	}
	return principalRemoteIP(cidr), nil
}

type srcNamespaceGenerator struct{}

func (srcNamespaceGenerator) permission(_, _ string, _ bool) (*rbacpb.Permission, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (srcNamespaceGenerator) principal(_, value string, _ bool) (*rbacpb.Principal, error) {
	v := strings.Replace(value, "*", ".*", -1)
	m := matcher.StringMatcherRegex(fmt.Sprintf(".*/ns/%s/.*", v))
	return principalAuthenticated(m), nil
}

type srcPrincipalGenerator struct{}

func (srcPrincipalGenerator) permission(_, _ string, _ bool) (*rbacpb.Permission, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (srcPrincipalGenerator) principal(key, value string, _ bool) (*rbacpb.Principal, error) {
	m := matcher.StringMatcherWithPrefix(value, spiffe.URIPrefix)
	return principalAuthenticated(m), nil
}

type requestPrincipalGenerator struct{}

func (requestPrincipalGenerator) permission(_, _ string, _ bool) (*rbacpb.Permission, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (requestPrincipalGenerator) principal(key, value string, forTCP bool) (*rbacpb.Principal, error) {
	if forTCP {
		return nil, fmt.Errorf("%q is HTTP only", key)
	}

	m := matcher.MetadataStringMatcher(sm.AuthnFilterName, key, matcher.StringMatcher(value))
	return principalMetadata(m), nil
}

type requestAudiencesGenerator struct{}

func (requestAudiencesGenerator) permission(key, value string, forTCP bool) (*rbacpb.Permission, error) {
	return requestPrincipalGenerator{}.permission(key, value, forTCP)
}

func (requestAudiencesGenerator) principal(key, value string, forTCP bool) (*rbacpb.Principal, error) {
	return requestPrincipalGenerator{}.principal(key, value, forTCP)
}

type requestPresenterGenerator struct{}

func (requestPresenterGenerator) permission(key, value string, forTCP bool) (*rbacpb.Permission, error) {
	return requestPrincipalGenerator{}.permission(key, value, forTCP)
}

func (requestPresenterGenerator) principal(key, value string, forTCP bool) (*rbacpb.Principal, error) {
	return requestPrincipalGenerator{}.principal(key, value, forTCP)
}

type requestHeaderGenerator struct{}

func (requestHeaderGenerator) permission(_, _ string, _ bool) (*rbacpb.Permission, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (requestHeaderGenerator) principal(key, value string, forTCP bool) (*rbacpb.Principal, error) {
	if forTCP {
		return nil, fmt.Errorf("%q is HTTP only", key)
	}

	header, err := extractNameInBrackets(strings.TrimPrefix(key, attrRequestHeader))
	if err != nil {
		return nil, err
	}
	m := matcher.HeaderMatcher(header, value)
	return principalHeader(m), nil
}

type requestClaimGenerator struct{}

func (requestClaimGenerator) permission(_, _ string, _ bool) (*rbacpb.Permission, error) {
	return nil, fmt.Errorf("unimplemented")
}

func (requestClaimGenerator) principal(key, value string, forTCP bool) (*rbacpb.Principal, error) {
	if forTCP {
		return nil, fmt.Errorf("%q is HTTP only", key)
	}

	claims, err := extractNameInNestedBrackets(strings.TrimPrefix(key, attrRequestClaims))
	if err != nil {
		return nil, err
	}
	// Generate a metadata list matcher for the given path keys and value.
	// On proxy side, the value should be of list type.
	m := MetadataMatcherForJWTClaims(claims, matcher.StringMatcher(value))
	return principalMetadata(m), nil
}

type hostGenerator struct{}

func (hg hostGenerator) permission(key, value string, forTCP bool) (*rbacpb.Permission, error) {
	if forTCP {
		return nil, fmt.Errorf("%q is HTTP only", key)
	}

	return permissionHeader(matcher.HostMatcher(hostHeader, value)), nil
}

func (hostGenerator) principal(key, value string, forTCP bool) (*rbacpb.Principal, error) {
	return nil, fmt.Errorf("unimplemented")
}

type pathGenerator struct{}

func (g pathGenerator) permission(key, value string, forTCP bool) (*rbacpb.Permission, error) {
	if forTCP {
		return nil, fmt.Errorf("%q is HTTP only", key)
	}

	m := matcher.PathMatcher(value)
	return permissionPath(m), nil
}

func (pathGenerator) principal(key, value string, forTCP bool) (*rbacpb.Principal, error) {
	return nil, fmt.Errorf("unimplemented")
}

type methodGenerator struct{}

func (methodGenerator) permission(key, value string, forTCP bool) (*rbacpb.Permission, error) {
	if forTCP {
		return nil, fmt.Errorf("%q is HTTP only", key)
	}

	m := matcher.HeaderMatcher(methodHeader, value)
	return permissionHeader(m), nil
}

func (methodGenerator) principal(key, value string, forTCP bool) (*rbacpb.Principal, error) {
	return nil, fmt.Errorf("unimplemented")
}
