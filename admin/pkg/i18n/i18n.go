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

package i18n

import (
	"strings"
)

// Language constants
const (
	LangEnglish = "en"
	LangChinese = "zh"
)

// T translates a message key to the specified language.
// Falls back to English if the key is not found in the requested language.
// Returns the key itself if not found in any language.
func T(lang, key string) string {
	if msgs, ok := messages[lang]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	// Fallback to English
	if msgs, ok := messages[LangEnglish]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}
	return key
}

// ParseAcceptLanguage parses the Accept-Language header and returns the preferred language.
// Supports formats like "zh-CN,zh;q=0.9,en;q=0.8" or simple "zh" / "en".
func ParseAcceptLanguage(header string) string {
	if header == "" {
		return LangEnglish
	}

	// Simple detection: if contains "zh", use Chinese
	lower := strings.ToLower(header)
	if strings.Contains(lower, "zh") {
		return LangChinese
	}

	return LangEnglish
}
