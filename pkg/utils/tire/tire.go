package tire

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/utils/urlPath"
)

type Tire struct {
	root Node
}

func NewTire() Tire {
	return Tire{root: Node{endOfPath: false, matchStr: "*"}}
}

// https://hsot:port/path1/{pathvarible1}/path2/{pathvarible2}

type Node struct {
	matchStr         string           //冗余信息暂时无用，rebuild 需要
	children         map[string]*Node //子树
	PathVariablesSet map[string]*Node //变量名集合 找不到set先用map todo
	PathVariableNode *Node            //通配变量节点后的子树
	endOfPath        bool             //是否是路径末尾
	bizInfo          interface{}      //随便塞啥
}

func (tire *Tire) Put(withOutHost string, bizInfo interface{}) bool {
	parts := urlPath.Split(withOutHost)

	return tire.root.Put(parts, bizInfo)
}
func (tire Tire) Get(withOutHost string) (*Node, []string, bool) {
	parts := urlPath.Split(withOutHost)
	return tire.root.Get(parts)
}

func (tire Tire) Match(withOutHost string) (*Node, *[]string, bool) {
	parts := urlPath.Split(withOutHost)
	return tire.root.Match(parts)
}

//不释放内存，释放内存需要使用方rebuild 整个字典树
func (tire Tire) Remove(withOutHost string) {
	n, _, _ := tire.Get(withOutHost)
	if n != nil {
		n.endOfPath = false
	}
}

func (tire Tire) Contains(withOutHost string) bool {
	parts := urlPath.Split(withOutHost)
	ret, _, _ := tire.root.Get(parts)
	return !(ret == nil)
}

func (node *Node) Put(keys []string, bizInfo interface{}) bool {
	if node.children == nil {
		node.children = map[string]*Node{}
	}
	if len(keys) == 0 {
		return true
	}
	key := keys[0]
	isReal := len(keys) == 1
	isSuccess := node.put(key, isReal, bizInfo)
	if !isSuccess {
		return false
	}
	childKeys := keys[1:]
	if urlPath.IsPathVariable(key) {
		return node.PathVariableNode.Put(childKeys, bizInfo)
	} else {
		return node.children[key].Put(childKeys, bizInfo)
	}

}

func (node *Node) GetBizInfo() interface{} {
	return node.bizInfo
}

func (node *Node) Match(parts []string) (*Node, *[]string, bool) {
	key := parts[0]
	childKeys := parts[1:]
	isReal := len(childKeys) == 0
	if isReal {
		if node.children != nil && node.children[key] != nil && node.children[key].endOfPath {
			return node.children[key], &[]string{}, true
		}
		//不能直接return 需要一次回朔 O（2n）    tire下存在：/aaa/bbb/xxxxx/ccc/ddd  /aaa/bbb/:id/ccc   输入url：/aaa/bbb/xxxxx/ccc
		if node.PathVariableNode != nil {
			if node.PathVariableNode.endOfPath {
				return node.PathVariableNode, &[]string{key}, true
			}
		}
		return nil, nil, false
	} else {
		if node.children != nil && node.children[key] != nil {
			node, param, ok := node.children[key].Match(childKeys)
			if ok {
				return node, param, ok
			}
		}
		//同理需要回朔
		if node.PathVariableNode != nil {
			node, param, ok := node.PathVariableNode.Match(childKeys)
			if ok {
				newParams := append(*param, key)
				param = &newParams
				return node, param, ok
			}
		}
		return nil, nil, false
	}
}

func (node *Node) Get(keys []string) (*Node, []string, bool) {
	key := keys[0]
	childKeys := keys[1:]
	isReal := len(childKeys) == 0
	if isReal {
		if urlPath.IsPathVariable(key) {
			if node.PathVariableNode == nil {
				return nil, nil, false
			}
			if node.PathVariableNode.endOfPath {
				return node.PathVariableNode, []string{urlPath.VariableName(key)}, true
			} else {
				return nil, nil, false
			}
		} else {
			if node.children == nil {
				return nil, nil, false
			}
			return node.children[key], nil, true
		}
	} else {
		if urlPath.IsPathVariable(key) {
			if node.PathVariableNode == nil {
				return nil, nil, false
			}
			retNode, pathVariableList, ok := node.PathVariableNode.Get(childKeys)
			return retNode, append(pathVariableList, key[1:len(key)-1]), ok
		} else {
			if node.children == nil || node.children[key] == nil {
				return nil, nil, false
			}
			return node.children[key].Get(childKeys)
		}
	}

}

func (node *Node) put(key string, isReal bool, bizInfo interface{}) bool {
	if urlPath.IsPathVariable(key) {
		pathVariable := urlPath.VariableName(key)
		return node.putPathVariable(pathVariable, isReal, bizInfo)
	} else {
		return node.putNode(key, isReal, bizInfo)
	}
}

func (node *Node) putPathVariable(pathVariable string, isReal bool, bizInfo interface{}) bool {
	if node.PathVariableNode == nil {
		node.PathVariableNode = &Node{endOfPath: false}
	}
	if node.PathVariableNode.endOfPath && isReal {
		//已经有一个同路径变量结尾的url 冲突
		return false
	}
	if isReal {
		node.PathVariableNode.bizInfo = bizInfo
	}
	node.PathVariableNode.endOfPath = node.PathVariableNode.endOfPath || isReal
	if node.PathVariablesSet == nil {
		node.PathVariablesSet = map[string]*Node{}
	}
	node.PathVariablesSet[pathVariable] = node.PathVariableNode
	return true
}

func (node *Node) putNode(matchStr string, isReal bool, bizInfo interface{}) bool {
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
