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

package trie

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTire_Put(t *testing.T) {
	trie := NewTrie()
	ret := trie.Put("/path1/:pathvarible1/path2/:pathvarible2", nil)
	assert.True(t, ret)
	ret = trie.Put("/path2/:pathvarible1/path2/:pathvarible2", nil)
	assert.True(t, ret)
	ret = trie.Put("/path2/3/path2/:pathvarible2", nil)
	assert.True(t, ret)

	ret = trie.Put("/path2/3/path2/:pathvarible2", nil)
	assert.False(t, ret)

	ret = trie.Put("/path2/3/path2/:pathvarible2/3", nil)
	assert.True(t, ret)
	ret = trie.Put("/path2/3/path2/:432423/3", nil)
	assert.False(t, ret)
	ret = trie.Put("/path2/3/path2/:432423/3/a/b/c/d/:fdsa", nil)
	assert.True(t, ret)

	ret = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsa", nil)
	assert.True(t, ret)

	ret = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsafdsafsdafsda", nil)
	assert.False(t, ret)

	ret = trie.Put("/path1/:pathvarible1/path2/:pathvarible2/:fdsa", nil)
	assert.True(t, ret)

	ret = trie.Put("/path1/:432/path2/:34", nil)
	assert.False(t, ret)
}

func TestTire_MatchAndGet(t *testing.T) {
	trie := NewTrie()
	ret := trie.Put("/path1/:pathvarible1/path2/:pathvarible2", "test1")
	assert.True(t, ret)
	ret = trie.Put("/path2/:pathvarible1/path2/:pathvarible2", nil)
	assert.True(t, ret)
	ret = trie.Put("/path2/3/path2/:pathvarible2", nil)
	assert.True(t, ret)

	ret = trie.Put("/path2/3/path2/:pathvarible2", nil)
	assert.False(t, ret)

	ret = trie.Put("/path2/3/path2/:pathvarible2/3", nil)
	assert.True(t, ret)
	ret = trie.Put("/path2/3/path2/:432423/3", nil)
	assert.False(t, ret)
	ret = trie.Put("/path2/3/path2/:432423/3/a/b/c/d/:fdsa", nil)
	assert.True(t, ret)

	ret = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsa", nil)
	assert.True(t, ret)

	ret = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsafdsafsdafsda", nil)
	assert.False(t, ret)

	ret = trie.Put("/path1/:pathvarible1/path2/:pathvarible2/:fdsa", nil)
	assert.True(t, ret)

	ret = trie.Put("/path1/:432/path2/:34", nil)
	assert.False(t, ret)

	node, param, ok := trie.Match("/path1/v1/path2/v2")
	assert.True(t, ok)
	assert.True(t, len(*param) == 2)
	assert.True(t, node.GetBizInfo() == "test1")

	_, _, ok = trie.Get("/path1/v1/path2/v2")
	assert.False(t, ok)

	node, _, ok = trie.Get("/path1/:p/path2/:p2")
	assert.True(t, ok)
	assert.True(t, node.GetBizInfo() == "test1")
}
