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
	"slices"
	"testing"
)

func TestCompileHeaders(t *testing.T) {
	precompiled := HeaderMatcher{Name: "X-Precompiled", Regex: true}
	if err := precompiled.SetValueRegex("^ready$"); err != nil {
		t.Fatal(err)
	}
	headers := []HeaderMatcher{
		{Name: "X-Visible", Values: []string{"blue", "green"}},
		{Name: "X-Cached", Regex: true, Values: []string{"^t[0-9]+$"}},
		precompiled,
		{Name: "X-Empty-Regex", Regex: true},
		{Name: "X-Invalid", Regex: true, Values: []string{"["}},
		{Name: "X-No-Values"},
	}

	compiled := compileHeaders(headers)
	wantNames := []string{"X-Visible", "X-Cached", "X-Precompiled", "X-Empty-Regex", "X-No-Values"}
	names := make([]string, 0, len(compiled))
	for _, header := range compiled {
		names = append(names, header.Name)
	}
	if !slices.Equal(names, wantNames) {
		t.Fatalf("compiled header names = %v, want %v", names, wantNames)
	}
	if !slices.Equal(compiled[0].Values, []string{"blue", "green"}) {
		t.Fatalf("plain values = %v, want [blue green]", compiled[0].Values)
	}
	if compiled[1].Regex == nil || !compiled[1].Regex.MatchString("t17") {
		t.Fatal("cached regex did not match")
	}
	if compiled[2].Regex == nil || !compiled[2].Regex.MatchString("ready") {
		t.Fatal("precompiled regex did not match")
	}
	if compiled[3].Regex != nil || len(compiled[3].Values) != 0 {
		t.Fatalf("empty regex = (%v, %v), want no condition", compiled[3].Regex, compiled[3].Values)
	}
	if compiled[4].Regex != nil || len(compiled[4].Values) != 0 {
		t.Fatalf("empty values = (%v, %v), want no condition", compiled[4].Regex, compiled[4].Values)
	}
}
