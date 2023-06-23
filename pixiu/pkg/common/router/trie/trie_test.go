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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/util/stringutil"
)

func TestTrie_Put(t *testing.T) {
	trie := NewTrie()
	ret, _ := trie.Put("/path1/:pathvarible1/path2/:pathvarible2", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path1/:pathvarible1/path2/:pathvarible2/**", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/:pathvarible1/path2/:pathvarible2", "")
	assert.True(t, ret)
	ret, _ = trie.Put("/path2/3/path2/:pathvarible2", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:pathvarible2", "")
	assert.False(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:pathvarible2/3", "")
	assert.True(t, ret)
	ret, _ = trie.Put("/path2/3/path2/:432423/3", "")
	assert.False(t, ret)
	ret, _ = trie.Put("/path2/3/path2/:432423/3/a/b/c/d/:fdsa", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsa", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsafdsafsdafsda", "")
	assert.False(t, ret)

	ret, _ = trie.Put("/path1/:pathvarible1/path2/:pathvarible2/:fdsa", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path1/:432/path2/:34", "")
	assert.False(t, ret)
}

func TestTrie_MatchAndGet(t *testing.T) {
	trie := NewTrie()

	ret, _ := trie.Put("/path1/:pathvarible1/path2/:pathvarible2", "test1")
	assert.True(t, ret)

	_, _ = trie.Put("/a/b", "ab")
	result, _, _ := trie.Match("/a/b")
	assert.Equal(t, result.GetBizInfo(), "ab")

	result, _, _ = trie.Match("/a/b?a=b&c=d")
	assert.Equal(t, result.GetBizInfo(), "ab")

	_, _ = trie.Put("POST/api/v1/**", "ab")
	result, _, _ = trie.Match("POST/api/v1")
	assert.Equal(t, "ab", result.GetBizInfo())

	ret, _ = trie.Put("/path2/:pathvarible1/path2/:pathvarible2", "")
	assert.True(t, ret)
	ret, _ = trie.Put("/path2/3/path2/:pathvarible2", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:pathvarible2", "")
	assert.False(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:pathvarible2/3", "")
	assert.True(t, ret)
	ret, _ = trie.Put("/path2/3/path2/:432423/3", "")
	assert.False(t, ret)
	ret, _ = trie.Put("/path2/3/path2/:432423/3/a/b/c/d/:fdsa", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsa", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsafdsafsdafsda", "")
	assert.False(t, ret)

	ret, _ = trie.Put("/path1/:pathvarible1/path2/:pathvarible2/:fdsa", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path1/:432/path2/:34", "")
	assert.False(t, ret)

	ret, _ = trie.Put("/a/:432/b/:34/**", "test**")
	assert.True(t, ret)

	ret, _ = trie.Put("/a/:432/b/:34/**", "")
	assert.False(t, ret)

	node, param, ok := trie.Match("/a/v1/b/v2/sadf/asdf")
	assert.Equal(t, (param)[0], "v1")
	assert.Equal(t, (param)[1], "v2")
	assert.Equal(t, node.GetBizInfo(), "test**")
	assert.True(t, ok)

	_, param, ok = trie.Match("/a/v1/b/v2/sadf")
	assert.Equal(t, (param)[0], "v1")
	assert.Equal(t, (param)[1], "v2")
	assert.Equal(t, node.GetBizInfo(), "test**")
	assert.True(t, ok)

	node, param, ok = trie.Match("/path1/v1/path2/v2")
	assert.True(t, ok)
	assert.True(t, len(param) == 2)
	assert.Equal(t, (param)[0], "v1")
	assert.Equal(t, (param)[1], "v2")
	assert.True(t, node.GetBizInfo() == "test1")

	_, _, ok, _ = trie.Get("/path1/v1/path2/v2")
	assert.False(t, ok)

	node, _, ok, _ = trie.Get("/path1/:p/path2/:p2")
	assert.True(t, ok)
	assert.True(t, node.GetBizInfo() == "test1")

	node, _, ok = trie.Match("/path1/12/path2/12?a=b")
	assert.True(t, ok)
	assert.True(t, node.GetBizInfo() == "test1")
}

func TestTrie_Clear(t *testing.T) {
	v := "http://baidu.com/aa/bb"
	v = stringutil.GetTrieKey("PUT", v)
	assert.Equal(t, "PUT/aa/bb", v)

	trie := NewTrie()
	ret, _ := trie.Put("/path1/:pathvarible1/path2/:pathvarible2", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path1/:pathvarible1/path2/:pathvarible2/**", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/:pathvarible1/path2/:pathvarible2", "")
	assert.True(t, ret)
	ret, _ = trie.Put("/path2/3/path2/:pathvarible2", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:pathvarible2", "")
	assert.False(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:pathvarible2/3", "")
	assert.True(t, ret)
	ret, _ = trie.Put("/path2/3/path2/:432423/3", "")
	assert.False(t, ret)
	ret, _ = trie.Put("/path2/3/path2/:432423/3/a/b/c/d/:fdsa", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsa", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path2/3/path2/:432423/3/a/b/c/c/:fdsafdsafsdafsda", "")
	assert.False(t, ret)

	ret, _ = trie.Put("/path1/:pathvarible1/path2/:pathvarible2/:fdsa", "")
	assert.True(t, ret)

	ret, _ = trie.Put("/path1/:432/path2/:34", "")

	assert.False(t, ret)
	assert.False(t, trie.IsEmpty())
	trie.Clear()
	assert.True(t, trie.IsEmpty())
}

func TestTrie_ParamMatch(t *testing.T) {
	trie := NewTrie()
	ret, _ := trie.Put("PUT/path1/:pathvarible1/path2/:pathvarible2", "")
	assert.True(t, ret)
	str := "https://www.baidu.com/path1/param1/path2/param2?aaaaa=aaaaa"

	node, _, ok := trie.Match(stringutil.GetTrieKey("PUT", str))
	assert.True(t, ok)
	assert.Equal(t, "", node.GetBizInfo())

	ret, _ = trie.Put("PUT/path1/:pathvarible1/path2", "")
	node, _, ok = trie.Match(stringutil.GetTrieKey("PUT", str))
	assert.True(t, ok)
	assert.Equal(t, "", node.GetBizInfo())
	assert.True(t, ret)
}
