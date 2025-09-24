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

package util

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// ParseResult holds the result of URL parsing
type ParseResult struct {
	Host         string
	Port         int
	UsedFallback bool
	FallbackInfo string
}

// URLParseError represents a general URL parsing error
type URLParseError struct {
	URL    string
	Reason string
}

func (e *URLParseError) Error() string {
	return fmt.Sprintf("failed to parse URL '%s': %s", e.URL, e.Reason)
}

// ParseHostPortFromURL extracts host and port from a raw URL or host:port string.
// It supports forms like:
// - http://host:port/path?query
// - https://host:port
// - host:port
// Returns ParseResult with host, port, and fallback information, or error for invalid formats.
// Fallback ports: HTTP(80), HTTPS(443), others(8080)
func ParseHostPortFromURL(raw string) (*ParseResult, error) {
	if raw == "" {
		return nil, &URLParseError{URL: raw, Reason: "empty URL"}
	}

	addr := strings.TrimSpace(raw)

	if strings.Contains(addr, constant.ProtocolSlash) {
		u, err := url.Parse(addr)
		if err != nil {
			return nil, &URLParseError{URL: raw, Reason: fmt.Sprintf("invalid URL format: %v", err)}
		}
		if u.Host == "" {
			return nil, &URLParseError{URL: raw, Reason: "missing host in URL"}
		}
		return parseHostPortWithFallback(u.Host, u.Scheme, raw)
	}

	return parseHostPortWithFallback(addr, "", raw)
}

func parseHostPortWithFallback(hostport, scheme, originalURL string) (*ParseResult, error) {
	host, portStr, err := net.SplitHostPort(hostport)
	if err != nil {
		// No port specified, use fallback based on scheme
		host = strings.TrimSpace(hostport)
		if host == "" {
			return nil, &URLParseError{URL: originalURL, Reason: "missing host"}
		}

		fallbackPort := getFallbackPort(scheme)
		fallbackInfo := fmt.Sprintf("used fallback port %d", fallbackPort)
		if scheme != "" {
			fallbackInfo = fmt.Sprintf("used fallback port %d for protocol %s", fallbackPort, scheme)
		}

		return &ParseResult{
			Host:         host,
			Port:         fallbackPort,
			UsedFallback: true,
			FallbackInfo: fallbackInfo,
		}, nil
	}

	port, err := strconv.Atoi(portStr)
	if err != nil || port <= 0 {
		return nil, &URLParseError{URL: originalURL, Reason: fmt.Sprintf("invalid port '%s'", portStr)}
	}

	host = strings.TrimSpace(host)
	if host == "" {
		return nil, &URLParseError{URL: originalURL, Reason: "missing host"}
	}

	return &ParseResult{
		Host:         host,
		Port:         port,
		UsedFallback: false,
		FallbackInfo: "",
	}, nil
}

func getFallbackPort(scheme string) int {
	switch strings.ToLower(scheme) {
	case "http":
		return 80
	case "https":
		return 443
	default:
		return 8080 // Default for MCP servers
	}
}

// goTmplArgRe is the regex pattern for Go template arguments
var goTmplArgRe = regexp.MustCompile(`\{\{\.args\.(?P<name>[a-zA-Z0-9_\-]+)}}`)

// ExtractPathFromURL extracts the path component from a URL string.
// Supports MCP template parameter conversion: {{.args.name}} -> {name}
func ExtractPathFromURL(raw string) string {
	if raw == "" {
		return constant.PathSlash
	}
	s := strings.TrimSpace(raw)

	// Prefer url.Parse to extract the path
	if i := strings.Index(s, constant.ProtocolSlash); i >= 0 {
		if u, err := url.Parse(s); err == nil {
			path := u.Path
			if path == "" {
				path = constant.PathSlash
			}
			return ReplaceGoTemplateArgsInPath(path)
		}
		// Fallback: remove the scheme and process
		s = s[i+3:]
	}

	// Handle host[:port]/path form without a scheme
	slash := strings.IndexByte(s, '/')
	if slash >= 0 {
		// If the colon appears before the first slash, treat the portion after the slash as the path
		colon := strings.IndexByte(s, ':')
		if colon >= 0 && colon < slash {
			path := s[slash:]
			if path == "" {
				return constant.PathSlash
			}
			return ReplaceGoTemplateArgsInPath(path)
		}
		// Otherwise, it is a path or relative path
		if s[0] != '/' {
			return ReplaceGoTemplateArgsInPath(constant.PathSlash + s[slash+1:])
		}
		return ReplaceGoTemplateArgsInPath(s[slash:])
	}

	// No slash found, return root path
	return constant.PathSlash
}

// ReplaceGoTemplateArgsInPath converts Go template args {{.args.name}} to standard format {name}
func ReplaceGoTemplateArgsInPath(path string) string {
	if path == "" {
		return constant.PathSlash
	}
	return goTmplArgRe.ReplaceAllString(path, `{$1}`)
}

// ValidateNacosAddresses validates comma-separated Nacos addresses
func ValidateNacosAddresses(addresses string) error {
	if strings.TrimSpace(addresses) == "" {
		return fmt.Errorf("nacos addresses cannot be empty")
	}

	var validCount int
	var errors []string

	for _, part := range strings.Split(addresses, ",") {
		addr := strings.TrimSpace(part)
		if addr == "" {
			continue
		}

		if _, err := ParseHostPortFromURL(addr); err != nil {
			errors = append(errors, fmt.Sprintf("invalid address '%s': %v", addr, err))
		} else {
			validCount++
		}
	}

	if validCount == 0 {
		return fmt.Errorf("no valid nacos addresses found: %s", strings.Join(errors, "; "))
	}

	if len(errors) > 0 {
		logger.Warnf("[dubbo-go-pixiu] some nacos addresses are invalid but continuing with valid ones: %s", strings.Join(errors, "; "))
	}

	return nil
}
