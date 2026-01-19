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

package proxy

import (
	"strings"
)

// Triple-specific header prefix
const (
	TripleHeaderPrefix = "tri-"
)

// ExtractTripleMetadata extracts Triple-specific metadata from attachments.
// This is used to extract tri-* headers for routing decisions based on
// service version, group, etc. Note: This filter is only used by gRPC listener,
// which handles gRPC protocol. Triple protocol should use triple listener.
// This function is kept for potential cross-protocol scenarios where Triple
// headers might be present in gRPC requests for routing purposes.
func ExtractTripleMetadata(attachments map[string]any) map[string]string {
	meta := make(map[string]string)
	for k, v := range attachments {
		key := strings.ToLower(k)
		if strings.HasPrefix(key, TripleHeaderPrefix) {
			if str, ok := v.(string); ok {
				meta[key] = str
			}
		}
	}
	return meta
}

// IsTripleHeader checks if a header key is Triple-specific
func IsTripleHeader(key string) bool {
	return strings.HasPrefix(strings.ToLower(key), TripleHeaderPrefix)
}
