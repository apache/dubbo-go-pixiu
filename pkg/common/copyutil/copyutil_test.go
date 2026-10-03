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

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestCloneJSONLikeDefensiveCopy(t *testing.T) {
	input := map[string]any{
		"nested": map[string]any{"roles": []any{"reader"}},
		"labels": map[string]string{
			"tenant": "acme",
		},
		"scopes": []string{"read"},
	}

	cloned := CloneStringAnyMap(input)
	cloned["nested"].(map[string]any)["roles"].([]any)[0] = "admin"
	cloned["labels"].(map[string]string)["tenant"] = "globex"
	cloned["scopes"].([]string)[0] = "write"

	assert.Equal(t, "reader", input["nested"].(map[string]any)["roles"].([]any)[0])
	assert.Equal(t, "acme", input["labels"].(map[string]string)["tenant"])
	assert.Equal(t, "read", input["scopes"].([]string)[0])
}
