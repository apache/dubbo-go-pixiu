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
	"github.com/dubbogo/dubbo-go-proxy/pkg/utils/urlpath"
)

// Trie
type Trie struct {
	root Node
}

// NewTrie
func NewTrie() Trie {
	return Trie{root: Node{endOfPath: false, matchStr: "*"}}
}

// Node
type Node struct {
	matchStr         string           //冗余信息 通配节点冗余变量名，普通节点冗余节点名
	children         map[string]*Node //子树
	PathVariablesSet map[string]*Node //变量名集合 找不到set先用map todo
	PathVariableNode *Node            //通配变量节点后的子树
	endOfPath        bool             //是否是路径末尾
	bizInfo          interface{}      //随便塞啥
}

//Trie对外暴露 推荐外部使用

//Put put key and values into trie as map.
func (trie *Trie) Put(withOutHost string, bizInfo interface{}) bool {
	parts := urlpath.Split(withOutHost)
	return trie.root.Put(parts, bizInfo)
}

//Get get values according key.pathVariable not supported.
func (trie Trie) Get(withOutHost string) (*Node, []string, bool) {
	parts := urlpath.Split(withOutHost)
	node, param, ok := trie.root.Get(parts)
	length := len(param)
	for i := 0; i < length/2; i++ {
		temp := param[length-1-i]
		param[length-1-i] = param[i]
		param[i] = temp
	}
	return node, param, ok
}

//Match get values according url , pathVariable supported.
func (trie Trie) Match(withOutHost string) (*Node, []string, bool) {
	parts := urlpath.Split(withOutHost)
	node, param, ok := trie.root.Match(parts)
	length := len(param)
	for i := 0; i < length/2; i++ {
		temp := param[length-1-i]
		param[length-1-i] = param[i]
		param[i] = temp
	}
	return node, param, ok
}

//Remove remove key and value from trie. 不释放内存，释放内存需要使用方rebuild 整个字典树
func (trie Trie) Remove(withOutHost string) {
	n, _, _ := trie.Get(withOutHost)
	if n != nil {
		n.endOfPath = false
	}
}

//Contains return if key exists in trie
func (trie Trie) Contains(withOutHost string) bool {
	parts := urlpath.Split(withOutHost)
	ret, _, _ := trie.root.Get(parts)
	return !(ret == nil)
}

//不对外暴露，不推荐外部使用

//Put node put
func (node *Node) Put(keys []string, bizInfo interface{}) bool {
	//空节点初始化
	if node.children == nil {
		node.children = map[string]*Node{}
	}
	if len(keys) == 0 {
		return true
	}

	key := keys[0]
	// isReal 代表是否是输入url 最末尾那段,对应trie 上的节点是否真实存在。
	isReal := len(keys) == 1
	//屏蔽 通配和非统配的细节put 方法，不涉及递归
	isSuccess := node.put(key, isReal, bizInfo)
	//递归退出条件
	if !isSuccess {
		return false
	}
	childKeys := keys[1:]

	//递归体
	if urlpath.IsPathVariable(key) {
		return node.PathVariableNode.Put(childKeys, bizInfo)
	} else {
		return node.children[key].Put(childKeys, bizInfo)
	}

}

//GetBizInfo get info
func (node *Node) GetBizInfo() interface{} {
	return node.bizInfo
}

//Match node match
func (node *Node) Match(parts []string) (*Node, []string, bool) {
	key := parts[0]
	childKeys := parts[1:]
	// isReal 代表是否是输入url 最末尾那段,对应trie 上的节点是否真实存在。
	isReal := len(childKeys) == 0
	if isReal {
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
		return nil, nil, false
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
		return nil, nil, false
	}
}

//Get node get
func (node *Node) Get(keys []string) (*Node, []string, bool) {
	key := keys[0]
	childKeys := keys[1:]
	isReal := len(childKeys) == 0
	if isReal {
		//退出条件
		if urlpath.IsPathVariable(key) {
			if node.PathVariableNode == nil || !node.PathVariableNode.endOfPath {
				return nil, nil, false
			}
			return node.PathVariableNode, []string{urlpath.VariableName(key)}, true
		} else {
			if node.children == nil {
				return nil, nil, false
			}
			return node.children[key], nil, true
		}
	} else {
		//递归体
		if urlpath.IsPathVariable(key) {
			if node.PathVariableNode == nil {
				return nil, nil, false
			}
			retNode, pathVariableList, ok := node.PathVariableNode.Get(childKeys)
			newList := []string{key}
			copy(newList[1:], pathVariableList)
			return retNode, newList, ok
		} else {
			if node.children == nil || node.children[key] == nil {
				return nil, nil, false
			}
			return node.children[key].Get(childKeys)
		}
	}

}

func (node *Node) put(key string, isReal bool, bizInfo interface{}) bool {
	//不涉及递归，屏蔽变量 非变量 put 细节
	if !urlpath.IsPathVariable(key) {
		return node.putNode(key, isReal, bizInfo)
	}
	pathVariable := urlpath.VariableName(key)
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
