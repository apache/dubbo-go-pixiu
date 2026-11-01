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
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/util/stringutil"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// Trie
type Trie struct {
	root Node
}

// NewTrie
func NewTrie() Trie {
	return Trie{root: Node{endOfPath: false, matchStr: ""}}
}

// NewTrieWithDefault
func NewTrieWithDefault(path string, defVal any) Trie {
	ret := Trie{root: Node{endOfPath: false, matchStr: ""}}
	_, _ = ret.Put(path, defVal)
	return ret
}

// Node
type Node struct {
	matchStr         string           // abc match abc, :a match all words as a variable names a , * match all words  ,** match all words and children.
	children         map[string]*Node // in path /a/b/c  , b is child of a , c is child of b
	PathVariablesSet map[string]*Node // in path /:a/b/c/:d , :a is a path variable node of level1 , :d is path variable node of level4
	PathVariableNode *Node            // in path /:a/b/c/:d , /b/c/:d is a child tree of pathVariable node :a ,and some special logic for match pathVariable it better not store in children.
	MatchAllNode     *Node            // /a/b/**  /** is a match all Node.
	endOfPath        bool             // if true means a real path exists ,  /a/b/c/d only node of d is true, a,b,c is false.
	bizInfo          any              // route info and any other info store here.
}

func (trie *Trie) Clear() bool {
	return trie.root.Clear()
}

// IsEmpty put key and values into trie as map.
func (trie *Trie) IsEmpty() bool {
	return trie.root.IsEmpty()
}

// Put put key and values into trie as map.
func (trie *Trie) Put(withOutHost string, bizInfo any) (bool, error) {
	if bizInfo == nil {
		return false, errors.Errorf("data to put should not be nil.")
	}
	parts := stringutil.Split(withOutHost)
	return trie.root.internalPut(parts, bizInfo)
}

// Put put key and values into trie as map.
func (trie *Trie) PutOrUpdate(withOutHost string, bizInfo any) (bool, error) {
	if bizInfo == nil {
		return false, errors.Errorf("data to put should not be nil.")
	}
	parts := stringutil.Split(withOutHost)
	if _, err := trie.Remove(withOutHost); err != nil {
		logger.Infof("PutOrUpdate function remote (withOutHost{%s}) = error{%v}", withOutHost, err)
	}
	//if n != nil {
	//	//TODO: log n.bizInfo for trouble shooting
	//}
	return trie.root.internalPut(parts, bizInfo)
}

// Get get values according key.pathVariable not supported.
func (trie Trie) Get(withOutHost string) (*Node, []string, bool, error) {
	parts := stringutil.Split(withOutHost)
	node, param, ok, e := trie.root.Get(parts)
	length := len(param)
	for i := 0; i < length/2; i++ {
		temp := param[length-1-i]
		param[length-1-i] = param[i]
		param[i] = temp
	}
	return node, param, ok, e
}

// Match get values according url , pathVariable supported.
// returns:
// *Node this node in path, if not exists return nil
// []string   reversed array of pathVariable value   /valA/valB/valC matchs  /:aa/:bb/:cc  returns array of (valC,valB,valA)
// bool is ok
// error
func (trie Trie) Match(withOutHost string) (*Node, []string, bool) {
	return trie.MatchWithPredicate(withOutHost, nil)
}

// MatchWithPredicate returns the first matching node accepted by predicate.
// A rejected match continues along the same traversal order as Match.
func (trie Trie) MatchWithPredicate(withOutHost string, predicate func(*Node) bool) (*Node, []string, bool) {
	withOutHost = strings.Split(withOutHost, "?")[0]
	parts := stringutil.Split(withOutHost)
	node, param, ok := trie.root.match(parts, predicate)
	length := len(param)
	for i := 0; i < length/2; i++ {
		temp := param[length-1-i]
		param[length-1-i] = param[i]
		param[i] = temp
	}
	return node, param, ok
}

// Remove remove key and value from trie.  logic delete can't release memory, rebuild if necessary when lake of memory.
func (trie Trie) Remove(withOutHost string) (*Node, error) {
	n, _, _, e := trie.Get(withOutHost)
	if e != nil {
		return nil, e
	}
	if n != nil {
		n.endOfPath = false
	}
	return n, nil
}

// Contains return if key exists in trie
func (trie Trie) Contains(withOutHost string) (bool, error) {
	parts := stringutil.Split(withOutHost)
	ret, _, _, e := trie.root.Get(parts)
	if e != nil {
		return true, e
	}
	return ret != nil, nil
}

// Put node put
func (node *Node) internalPut(keys []string, bizInfo any) (bool, error) {
	// empty node initialization
	if node.children == nil {
		node.children = map[string]*Node{}
	}
	if len(keys) == 0 {
		return true, nil
	}

	key := keys[0]
	// isReal is the end of url path, means node is a place of url end,so the path with parentNode has a real url exists.
	isReal := len(keys) == 1
	isSuccess := node.put(key, isReal, bizInfo)

	if !isSuccess {
		return false, nil
	}
	childKeys := keys[1:]

	if stringutil.IsPathVariableOrWildcard(key) {
		return node.PathVariableNode.internalPut(childKeys, bizInfo)
	} else if stringutil.IsMatchAll(key) {
		return isSuccess, nil
	} else {
		return node.children[key].internalPut(childKeys, bizInfo)
	}

}

func (node *Node) Clear() bool {
	*node = Node{}
	return true
}

// IsEmpty return true if empty
func (node *Node) IsEmpty() bool {
	if node.children == nil && node.matchStr == "" && node.PathVariableNode == nil && node.PathVariablesSet == nil {
		return true
	}
	return false
}

// GetBizInfo get info
func (node *Node) GetBizInfo() any {
	return node.bizInfo
}

//Match node match

func (node *Node) Match(parts []string) (*Node, []string, bool) {
	return node.match(parts, nil)
}

func (node *Node) match(parts []string, predicate func(*Node) bool) (*Node, []string, bool) {
	key := parts[0]
	if len(parts) == 1 {
		return node.matchEnd(key, predicate)
	}
	return node.matchRemaining(key, parts[1:], predicate)
}

func (node *Node) matchEnd(key string, predicate func(*Node) bool) (*Node, []string, bool) {
	if candidate, ok := node.matchLiteralEnd(key, predicate); ok {
		return candidate, []string{}, true
	}
	// Trie parameter nodes also match a literal request segment.
	if candidate, ok := node.matchPathVariableEnd(predicate); ok {
		return candidate, []string{key}, true
	}
	return node.matchAll(key, predicate)
}

func (node *Node) matchRemaining(key string, childKeys []string, predicate func(*Node) bool) (*Node, []string, bool) {
	if candidate, params, ok := node.matchLiteralChild(key, childKeys, predicate); ok {
		return candidate, params, true
	}
	if candidate, params, ok := node.matchPathVariableChild(key, childKeys, predicate); ok {
		return candidate, params, true
	}
	return node.matchAll(key, predicate)
}

func (node *Node) matchLiteralEnd(key string, predicate func(*Node) bool) (*Node, bool) {
	if node.children == nil {
		return nil, false
	}
	candidate := node.children[key]
	if candidate == nil || !candidate.endOfPath {
		return nil, false
	}
	if predicate != nil && !predicate(candidate) {
		return nil, false
	}
	return candidate, true
}

func (node *Node) matchPathVariableEnd(predicate func(*Node) bool) (*Node, bool) {
	if node.PathVariableNode == nil || !node.PathVariableNode.endOfPath {
		return nil, false
	}
	if predicate != nil && !predicate(node.PathVariableNode) {
		return nil, false
	}
	return node.PathVariableNode, true
}

func (node *Node) matchLiteralChild(
	key string,
	childKeys []string,
	predicate func(*Node) bool,
) (*Node, []string, bool) {
	if node.children == nil {
		return nil, nil, false
	}
	child := node.children[key]
	if child == nil {
		return nil, nil, false
	}
	return child.match(childKeys, predicate)
}

func (node *Node) matchPathVariableChild(
	key string,
	childKeys []string,
	predicate func(*Node) bool,
) (*Node, []string, bool) {
	if node.PathVariableNode == nil {
		return nil, nil, false
	}
	candidate, params, ok := node.PathVariableNode.match(childKeys, predicate)
	params = append(params, key)
	return candidate, params, ok
}

func (node *Node) matchAll(key string, predicate func(*Node) bool) (*Node, []string, bool) {
	if candidate := node.matchAllChild(key); matchesPredicate(candidate, predicate) {
		return candidate, []string{}, true
	}
	if matchesPredicate(node.MatchAllNode, predicate) {
		return node.MatchAllNode, []string{}, true
	}
	return nil, nil, false
}

func (node *Node) matchAllChild(key string) *Node {
	if node.children == nil {
		return nil
	}
	child := node.children[key]
	if child == nil {
		return nil
	}
	return child.MatchAllNode
}

func matchesPredicate(candidate *Node, predicate func(*Node) bool) bool {
	if candidate == nil {
		return false
	}
	return predicate == nil || predicate(candidate)
}

// Get node get
// returns:
// *Node this node in path, if not exists return nil
// []string key reversed array of pathVariable   /:aa/:bb/:cc  returns array of (cc,bb,aa)
// bool is ok
// error
func (node *Node) Get(keys []string) (*Node, []string, bool, error) {
	key := keys[0]
	childKeys := keys[1:]
	isReal := len(childKeys) == 0
	if isReal {
		// exit condition
		if stringutil.IsPathVariableOrWildcard(key) {
			if node.PathVariableNode == nil || !node.PathVariableNode.endOfPath {
				return nil, nil, false, nil
			}
			return node.PathVariableNode, []string{stringutil.VariableName(key)}, true, nil
		} else if stringutil.IsMatchAll(key) {
			return node.MatchAllNode, nil, true, nil
		} else {
			if node.children == nil {
				return nil, nil, false, nil
			}
			return node.children[key], nil, true, nil
		}
	} else {

		if stringutil.IsPathVariableOrWildcard(key) {
			if node.PathVariableNode == nil {
				return nil, nil, false, nil
			}
			retNode, pathVariableList, ok, e := node.PathVariableNode.Get(childKeys)
			newList := []string{key}
			copy(newList[1:], pathVariableList)
			return retNode, newList, ok, e
		} else if stringutil.IsMatchAll(key) {
			return nil, nil, false, errors.Errorf("router configuration is empty")
		} else {
			if node.children == nil || node.children[key] == nil {
				return nil, nil, false, nil
			}
			return node.children[key].Get(childKeys)
		}
	}

}

func (node *Node) put(key string, isReal bool, bizInfo any) bool {

	if !stringutil.IsPathVariableOrWildcard(key) {
		if stringutil.IsMatchAll(key) {
			return node.putMatchAllNode(key, isReal, bizInfo)
		} else {
			return node.putNode(key, isReal, bizInfo)
		}
	}
	pathVariable := stringutil.VariableName(key)
	return node.putPathVariable(pathVariable, isReal, bizInfo)
}

func (node *Node) putPathVariable(pathVariable string, isReal bool, bizInfo any) bool {
	//path variable put
	if node.PathVariableNode == nil {
		node.PathVariableNode = &Node{endOfPath: false}
	}
	if node.PathVariableNode.endOfPath && isReal {
		//has a node with same path exists. conflicted.
		return false
	}
	if isReal {
		node.PathVariableNode.bizInfo = bizInfo
		node.PathVariableNode.matchStr = pathVariable
	}
	node.PathVariableNode.endOfPath = node.PathVariableNode.endOfPath || isReal
	if node.PathVariablesSet == nil {
		node.PathVariablesSet = map[string]*Node{}
	}
	node.PathVariablesSet[pathVariable] = node.PathVariableNode
	return true
}

func (node *Node) putNode(matchStr string, isReal bool, bizInfo any) bool {

	selfNode := &Node{endOfPath: isReal, matchStr: matchStr}
	old := node.children[matchStr]
	if old != nil {
		if old.endOfPath && isReal {
			// already has one same path url
			return false
		}
		selfNode = old
	} else {
		old = selfNode
	}

	if isReal {
		selfNode.bizInfo = bizInfo
	}
	selfNode.endOfPath = isReal || old.endOfPath
	node.children[matchStr] = selfNode
	return true
}

func (node *Node) putMatchAllNode(matchStr string, isReal bool, bizInfo any) bool {

	selfNode := &Node{endOfPath: isReal, matchStr: matchStr}
	old := node.MatchAllNode
	if old != nil {
		if old.endOfPath && isReal {
			// already has one same path url
			return false
		}
		selfNode = old
	} else {
		old = selfNode
	}

	if isReal {
		selfNode.bizInfo = bizInfo
	}
	selfNode.endOfPath = selfNode.endOfPath || old.endOfPath
	node.MatchAllNode = selfNode
	return true
}
