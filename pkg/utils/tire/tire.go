package tire

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/utils/urlPath"
)

type Info struct {
	BizInfo interface{}
}

type Tire struct {
	root Node
}

func newRealNode(key string) Node {
	return Node{endOfPath: true, matchStr: key}
}

func newMidNode(key string) Node {
	return Node{endOfPath: false, matchStr: key}
}

// https://hsot:port/path1/{pathvarible1}/path2/{pathvarible2}

type Node struct {
	matchStr         string
	children         map[string]*Node
	PathVariablesSet map[string]*Node
	PathVariableNode *Node
	endOfPath        bool
	bizInfo          Info
}

func (tire Tire) Put(withOutHost string) bool {
	parts := urlPath.Split(withOutHost)
	return tire.root.Put(parts)
}

func (tire Tire) Get(withOutHost string) (*Node, []string) {
	parts := urlPath.Split(withOutHost)
	return tire.root.Get(parts)
}

func (tire Tire) Contains(withOutHost string) bool {
	parts := urlPath.Split(withOutHost)
	ret, _ := tire.root.Get(parts)
	return !(ret == nil)
}

func (node Node) Put(keys []string) bool {
	if node.children == nil {
		node.children = map[string]*Node{}
	}
	key := keys[0]
	childKeys := keys[1:]
	isReal := len(childKeys) == 0
	isSuccess := node.put(key, isReal)
	if !isSuccess {
		return false
	}
	return node.children[key].Put(childKeys)
}

func (node Node) Get(keys []string) (*Node, []string) {
	key := keys[0]
	childKeys := keys[1:]
	isReal := len(childKeys) == 0
	if isReal {
		if isPathVariable(key) {
			if node.PathVariableNode.endOfPath == true {
				return node.PathVariableNode, []string{key[1 : len(key)-1]}
			} else {
				return nil, nil
			}
			return node.PathVariableNode.Get(childKeys)
		} else {
			return node.children[key], nil
		}
	} else {
		if isPathVariable(key) {
			if node.PathVariableNode == nil {
				return nil, nil
			}
			retNode, pathVariableList := node.PathVariableNode.Get(childKeys)
			return retNode, append(pathVariableList, key[1:len(key)-1])
		} else {
			if node.children == nil || node.children[key] == nil {
				return nil, nil
			}
			return node.children[key].Get(childKeys)
		}
	}

}

func (node Node) put(key string, isReal bool) bool {
	if isPathVariable(key) {
		pathVariable := key[1 : len(key)-1]
		return node.putPathVariable(pathVariable, isReal)
	} else {
		return node.putNode(key, isReal)
	}
}

func (node Node) putPathVariable(pathVariable string, isReal bool) bool {
	if node.PathVariableNode == nil {
		node.PathVariableNode = &Node{endOfPath: false}
	}
	if node.PathVariableNode.endOfPath && isReal {
		//已经有一个同路径变量结尾的url 冲突
		return false
	}
	node.PathVariableNode.endOfPath = node.PathVariableNode.endOfPath || isReal
	if node.PathVariablesSet == nil {
		node.PathVariablesSet = map[string]*Node{}
	}
	node.PathVariablesSet[pathVariable] = node.PathVariableNode
	return true
}

func (node Node) putNode(matchStr string, isReal bool) bool {
	selfNode := Node{endOfPath: isReal, matchStr: matchStr}
	old := node.children[matchStr]
	if old.endOfPath && isReal {
		//已经有一个同路径的url 冲突
		return false
	}
	selfNode.endOfPath = selfNode.endOfPath || old.endOfPath
	node.children[matchStr] = &selfNode
	return true
}

func isPathVariable(key string) bool {
	return key[0] == '{' && key[len(key)] == '}'
}
