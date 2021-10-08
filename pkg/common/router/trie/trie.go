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
	"github.com/apache/dubbo-go-pixiu/pkg/common/util/stringutil"
	"github.com/pkg/errors"
	"sync"
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
func NewTrieWithDefault(path string, defVal interface{}) Trie {
	ret := Trie{root: Node{endOfPath: false, matchStr: ""}}
	_, _ = ret.Put(path, defVal)
	return ret
}

// Node
type Node struct {
	matchStr         string
	childInitOnce    sync.Once
	children         map[string]*Node
	PathVariablesSet map[string]*Node
	PathVariableNode *Node
	MatchAllNode     *Node
	endOfPath        bool
	bizInfo          interface{}
}

//IsEmpty put key and values into trie as map.
func (trie *Trie) IsEmpty() bool {
	return trie.root.IsEmpty()
}

//Put put key and values into trie as map.
func (trie *Trie) Put(withOutHost string, bizInfo interface{}) (bool, error) {
	parts := stringutil.Split(withOutHost)
	return trie.root.Put(parts, bizInfo)
}

//Get get values according key.pathVariable not supported.
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

//Match get values according url , pathVariable supported.
func (trie Trie) Match(withOutHost string) (*Node, []string, bool) {
	parts := stringutil.Split(withOutHost)
	node, param, ok := trie.root.Match(parts)
	length := len(param)
	for i := 0; i < length/2; i++ {
		temp := param[length-1-i]
		param[length-1-i] = param[i]
		param[i] = temp
	}
	return node, param, ok
}

//Remove remove key and value from trie.  logic delete can't release memory, rebuild if necessary when lake of memory.
func (trie Trie) Remove(withOutHost string) error {
	n, _, _, e := trie.Get(withOutHost)
	if e != nil {
		return e
	}
	if n != nil {
		n.endOfPath = false
	}
	return errors.Errorf("path not exists.")
}

//Contains return if key exists in trie
func (trie Trie) Contains(withOutHost string) (bool, error) {
	parts := stringutil.Split(withOutHost)
	ret, _, _, e := trie.root.Get(parts)
	if e != nil {
		return true, e
	}
	return !(ret == nil), nil
}

//不对外暴露，不推荐外部使用

//Put node put
func (node *Node) Put(keys []string, bizInfo interface{}) (bool, error) {
	//空节点初始化
	node.childInitOnce.Do(func() {
		if node.children == nil {
			node.children = map[string]*Node{}
		}
	})

	if len(keys) == 0 {
		return true, nil
	}

	key := keys[0]
	// isReal 代表是否是输入url 最末尾那段,对应trie 上的节点是否真实存在。
	isReal := len(keys) == 1
	//屏蔽 通配和非统配的细节put 方法，不涉及递归
	isSuccess := node.put(key, isReal, bizInfo)
	//递归退出条件
	if !isSuccess {
		return false, nil
	}
	childKeys := keys[1:]

	//递归体
	if stringutil.IsPathVariableOrWildcard(key) {
		return node.PathVariableNode.Put(childKeys, bizInfo)
	} else if stringutil.IsMatchAll(key) {
		return isSuccess, nil
	} else {
		return node.children[key].Put(childKeys, bizInfo)
	}

}

//IsEmpty return true if empty
func (node *Node) IsEmpty() bool {
	if node.children == nil && node.matchStr == "" && node.PathVariableNode == nil && node.PathVariablesSet == nil {
		return true
	}
	return false
}

//GetBizInfo get info
func (node *Node) GetBizInfo() interface{} {
	return node.bizInfo
}

//Match node match
func (node *Node) Match(parts []string) (*Node, []string, bool) {
	key := parts[0]
	childKeys := parts[1:]
	// isEnd 代表是否是输入url 最末尾那段,对应trie 上的节点是否真实存在。
	isEnd := len(childKeys) == 0
	if isEnd {
		//退出条件
		if node.children != nil && node.children[key] != nil && node.children[key].endOfPath {
			return node.children[key], []string{}, true
		}
		//不能直接return 需要一次回朔 O（2n）    trie下存在：/aaa/bbb/xxxxx/ccc/ddd  /aaa/bbb/:id/ccc   输入url：/aaa/bbb/xxxxx/ccc
		if node.PathVariableNode != nil {
			if node.PathVariableNode.endOfPath {
				return node.PathVariableNode, []string{key}, true
			}
		}

	} else {
		//递归体
		if node.children != nil && node.children[key] != nil {
			n, param, ok := node.children[key].Match(childKeys)
			if ok {
				return n, param, ok
			}
		}
		//同理需要回朔
		if node.PathVariableNode != nil {
			n, param, ok := node.PathVariableNode.Match(childKeys)
			param = append(param, key)
			if ok {
				return n, param, ok
			}
		}
	}
	if node.children != nil && node.children[key] != nil && node.children[key].MatchAllNode != nil {
		return node.children[key].MatchAllNode, []string{}, true
	}
	if node.MatchAllNode != nil {
		return node.MatchAllNode, []string{}, true
	}
	return nil, nil, false
}

//Get node get
func (node *Node) Get(keys []string) (*Node, []string, bool, error) {
	key := keys[0]
	childKeys := keys[1:]
	isReal := len(childKeys) == 0
	if isReal {
		//退出条件
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
		//递归体
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

func (node *Node) put(key string, isReal bool, bizInfo interface{}) bool {
	//不涉及递归，屏蔽变量 非变量 put 细节
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

func (node *Node) putPathVariable(pathVariable string, isReal bool, bizInfo interface{}) bool {
	//变量put
	if node.PathVariableNode == nil {
		node.PathVariableNode = &Node{endOfPath: false}
	}
	if node.PathVariableNode.endOfPath && isReal {
		//已经有一个同路径变量结尾的url 冲突
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

func (node *Node) putNode(matchStr string, isReal bool, bizInfo interface{}) bool {
	//普通node put
	selfNode := &Node{endOfPath: isReal, matchStr: matchStr}
	old := node.children[matchStr]
	if old != nil {
		if old.endOfPath && isReal {
			//已经有一个同路径的url 冲突
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
	node.children[matchStr] = selfNode
	return true
}

func (node *Node) putMatchAllNode(matchStr string, isReal bool, bizInfo interface{}) bool {
	//普通node put
	selfNode := &Node{endOfPath: isReal, matchStr: matchStr}
	old := node.MatchAllNode
	if old != nil {
		if old.endOfPath && isReal {
			//已经有一个同路径的url 冲突
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
