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
	"slices"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/util/stringutil"
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

func TestTrie_MatchWithPredicateFallsBackInTraversalOrder(t *testing.T) {
	tr := NewTrie()
	if ok, err := tr.Put("GET/items/42", "exact"); !ok || err != nil {
		t.Fatalf("Put exact = (%v, %v)", ok, err)
	}
	if ok, err := tr.Put("GET/items/**", "prefix"); !ok || err != nil {
		t.Fatalf("Put prefix = (%v, %v)", ok, err)
	}

	// The original API still chooses the exact route without a predicate.
	node, _, ok := tr.Match("GET/items/42")
	if !ok || node.GetBizInfo() != "exact" {
		t.Fatalf("Match() = (%#v, %v), want exact", node, ok)
	}

	node, _, ok = tr.MatchWithPredicate("GET/items/42", func(candidate *Node) bool {
		return candidate.GetBizInfo() == "prefix"
	})
	if !ok || node.GetBizInfo() != "prefix" {
		t.Fatalf("MatchWithPredicate() = (%#v, %v), want prefix after exact rejection", node, ok)
	}
}

func TestTrie_MatchWithPredicateFallbackPreservesParameters(t *testing.T) {
	tr := NewTrie()
	if ok, err := tr.Put("GET/a/:name/x", "literal-branch"); !ok || err != nil {
		t.Fatalf("Put literal branch = (%v, %v)", ok, err)
	}
	if ok, err := tr.Put("GET/:first/:second/x", "parameter-branch"); !ok || err != nil {
		t.Fatalf("Put parameter branch = (%v, %v)", ok, err)
	}

	node, params, ok := tr.MatchWithPredicate("GET/a/b/x", matchBizInfo("literal-branch"))
	if !ok || node.GetBizInfo() != "literal-branch" || !slices.Equal(params, []string{"b"}) {
		t.Fatalf("literal match = (%#v, %v, %v), want literal-branch with [b]", node, params, ok)
	}

	node, params, ok = tr.MatchWithPredicate("GET/a/b/x", matchBizInfo("parameter-branch"))
	if !ok || node.GetBizInfo() != "parameter-branch" || !slices.Equal(params, []string{"a", "b"}) {
		t.Fatalf("fallback match = (%#v, %v, %v), want parameter-branch with [a b]", node, params, ok)
	}
}

func TestTrie_MatchWithPredicateMatchAllFallback(t *testing.T) {
	tr := NewTrie()
	if ok, err := tr.Put("GET/items/:id", "parameter"); !ok || err != nil {
		t.Fatalf("Put parameter = (%v, %v)", ok, err)
	}
	if ok, err := tr.Put("GET/items/**", "prefix"); !ok || err != nil {
		t.Fatalf("Put prefix = (%v, %v)", ok, err)
	}

	node, params, ok := tr.MatchWithPredicate("GET/items/42", matchBizInfo("prefix"))
	if !ok || node.GetBizInfo() != "prefix" || len(params) != 0 {
		t.Fatalf("fallback match = (%#v, %v, %v), want prefix without params", node, params, ok)
	}

	node, _, ok = tr.MatchWithPredicate("GET/items/42", func(*Node) bool { return false })
	if ok || node != nil {
		t.Fatalf("rejected match = (%#v, %v), want no match", node, ok)
	}
}

func matchBizInfo(want any) func(*Node) bool {
	return func(candidate *Node) bool {
		return candidate.GetBizInfo() == want
	}
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

func TestTrieSoftDelete(t *testing.T) {
	trie := NewTrie()

	ret, _ := trie.Put("/a/b/c", "route1")
	assert.True(t, ret)

	ret, _ = trie.Put("/a/b/c/d", "route2")
	assert.True(t, ret)

	ret, _ = trie.Put("/a/b/e", "route3")
	assert.True(t, ret)

	node, _, ok := trie.Match("/a/b/c")
	assert.True(t, ok)
	assert.Equal(t, "route1", node.GetBizInfo())

	_, err := trie.Remove("/a/b/c")
	assert.NoError(t, err)

	node, _, ok = trie.Match("/a/b/c")
	assert.False(t, ok)
	assert.Nil(t, node)

	node, _, ok = trie.Match("/a/b/c/d")
	assert.True(t, ok)
	assert.Equal(t, "route2", node.GetBizInfo())

	node, _, ok = trie.Match("/a/b/e")
	assert.True(t, ok)
	assert.Equal(t, "route3", node.GetBizInfo())

	node, _, ok = trie.Match("/a/b/c/")
	assert.False(t, ok)
	assert.Nil(t, node)
}
