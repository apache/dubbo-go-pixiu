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

package mcpserver

import (
	"strconv"
	"strings"
)

const (
	DefaultRegistryName = "default"
	DefaultServerID     = "default"
)

// ServerSource identifies one dynamic registry publication source. The
// registry adapter is process-global today, so source identity must include the
// registry namespace and server ID instead of relying on a server-local key.
type ServerSource struct {
	Registry string
	ServerID string
}

func NewServerSource(registry, serverID string) ServerSource {
	return ServerSource{Registry: registry, ServerID: serverID}.Normalize()
}

func (s ServerSource) Normalize() ServerSource {
	registry := strings.TrimSpace(s.Registry)
	if registry == "" {
		registry = DefaultRegistryName
	}
	serverID := strings.TrimSpace(s.ServerID)
	if serverID == "" {
		serverID = DefaultServerID
	}
	return ServerSource{Registry: registry, ServerID: serverID}
}

func (s ServerSource) Key() string {
	source := s.Normalize()
	return sanitizeSourcePart(source.Registry) + "/" + sanitizeSourcePart(source.ServerID)
}

func (s ServerSource) String() string {
	return s.Key()
}

func sanitizeSourcePart(value string) string {
	if value == "" {
		return "_"
	}
	var b strings.Builder
	for _, r := range value {
		switch {
		case r >= 'a' && r <= 'z':
			b.WriteRune(r)
		case r >= 'A' && r <= 'Z':
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == '-' || r == '_' || r == '.':
			b.WriteRune(r)
		default:
			b.WriteString("%")
			b.WriteString(strconv.FormatInt(int64(r), 16))
		}
	}
	return b.String()
}
