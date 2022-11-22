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

package authz

import (
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/test/util/assert"
)

func TestNamespaceMatch(t *testing.T) {
	assert.Equal(t, namespaceMatch("test-login", "*"), true)

	assert.Equal(t, namespaceMatch("test-login", "test-*"), true)
	assert.Equal(t, namespaceMatch("test-login", "*-test"), false)

	assert.Equal(t, namespaceMatch("test-login", "login-*"), false)
	assert.Equal(t, namespaceMatch("test-login", "*-login"), true)
}
