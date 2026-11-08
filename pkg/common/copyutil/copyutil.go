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

package copyutil

// CloneJSONLike returns a defensive copy of the JSON-like structures used in
// MCP configs, claims and session plans. Non-container values are returned as-is.
func CloneJSONLike(value any) any {
	switch v := value.(type) {
	case nil:
		return nil
	case map[string]any:
		cp := make(map[string]any, len(v))
		for k, item := range v {
			cp[k] = CloneJSONLike(item)
		}
		return cp
	case []any:
		cp := make([]any, len(v))
		for i := range v {
			cp[i] = CloneJSONLike(v[i])
		}
		return cp
	case []string:
		return CloneStringSlice(v)
	case []int:
		cp := make([]int, len(v))
		copy(cp, v)
		return cp
	case []int64:
		cp := make([]int64, len(v))
		copy(cp, v)
		return cp
	case []float64:
		cp := make([]float64, len(v))
		copy(cp, v)
		return cp
	case []bool:
		cp := make([]bool, len(v))
		copy(cp, v)
		return cp
	case map[string]string:
		return CloneStringMap(v)
	default:
		return v
	}
}

func CloneStringSlice(values []string) []string {
	if values == nil {
		return nil
	}
	cp := make([]string, len(values))
	copy(cp, values)
	return cp
}

func CloneStringMap(values map[string]string) map[string]string {
	if values == nil {
		return nil
	}
	cp := make(map[string]string, len(values))
	for k, v := range values {
		cp[k] = v
	}
	return cp
}

func CloneStringAnyMap(values map[string]any) map[string]any {
	if values == nil {
		return nil
	}
	cp, _ := CloneJSONLike(values).(map[string]any)
	return cp
}
