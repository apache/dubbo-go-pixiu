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

package benchmark

import (
	"os"
	"strings"
	"testing"
)

func TestTripleBenchmarkReadmesAgree(t *testing.T) {
	english, err := os.ReadFile("README.md")
	if err != nil {
		t.Fatalf("read English benchmark README: %v", err)
	}
	chinese, err := os.ReadFile("README_CN.md")
	if err != nil {
		t.Fatalf("read Chinese benchmark README: %v", err)
	}

	for _, test := range []struct {
		name           string
		englishSection string
		englishMethod  string
		chineseMethod  string
	}{
		{"direct GetUser", "Triple Direct", "GetUser", "GetUser（直连）"},
		{"via Pixiu GetUser", "Triple via Pixiu", "GetUser", "GetUser（经 Pixiu）"},
		{"via Pixiu GetUsers", "Triple via Pixiu", "GetUsers", "GetUsers（经 Pixiu）"},
		{"via Pixiu SayHello", "Triple via Pixiu", "SayHello", "SayHello（经 Pixiu）"},
	} {
		t.Run(test.name, func(t *testing.T) {
			englishValues := readmeTableValues(t, string(english), "### "+test.englishSection, test.englishMethod)
			chineseValues := readmeTableValues(t, string(chinese), "## Triple", test.chineseMethod)
			if strings.Join(englishValues, "|") != strings.Join(chineseValues, "|") {
				t.Fatalf("Triple %s values differ: English=%v Chinese=%v", test.englishMethod, englishValues, chineseValues)
			}
		})
	}
}

func readmeTableValues(t *testing.T, readme, section, method string) []string {
	t.Helper()
	sectionStart := strings.Index(readme, section)
	if sectionStart < 0 {
		t.Fatalf("section %q not found", section)
	}
	sectionBody := readme[sectionStart:]
	if nextSection := strings.Index(sectionBody[1:], "\n## "); nextSection >= 0 {
		sectionBody = sectionBody[:nextSection+1]
	}

	for _, line := range strings.Split(sectionBody, "\n") {
		columns := strings.Split(line, "|")
		if len(columns) < 4 || strings.TrimSpace(columns[1]) != method {
			continue
		}
		values := make([]string, 0, len(columns)-3)
		for _, column := range columns[2 : len(columns)-1] {
			values = append(values, strings.ReplaceAll(strings.TrimSpace(column), "**", ""))
		}
		return values
	}
	t.Fatalf("method %q not found in section %q", method, section)
	return nil
}
