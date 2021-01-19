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

package tire

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestName(t *testing.T) {
	tire := NewTire()
	ret := tire.Put("/path1/:pathvarible1/path2/:pathvarible2", nil)
	assert.True(t, ret)
	ret = tire.Put("/path2/:pathvarible1/path2/:pathvarible2", nil)
	assert.True(t, ret)
	ret = tire.Put("/path2/3/path2/:pathvarible2", nil)
	assert.True(t, ret)

	ret = tire.Put("/path2/3/path2/:pathvarible2", nil)
	assert.False(t, ret)

	ret = tire.Put("/path2/3/path2/:pathvarible2/3", nil)
	assert.True(t, ret)
	ret = tire.Put("/path2/3/path2/:432423/3", nil)
	assert.False(t, ret)
	ret = tire.Put("/path2/3/path2/:432423/3/a/b/c/d/:fdsa", nil)
	assert.True(t, ret)

	ret = tire.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsa", nil)
	assert.True(t, ret)

	ret = tire.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsafdsafsdafsda", nil)
	assert.False(t, ret)

	ret = tire.Put("/path1/:pathvarible1/path2/:pathvarible2/:fdsa", nil)
	assert.True(t, ret)

	ret = tire.Put("/path1/:432/path2/:34", nil)
	assert.False(t, ret)
}
