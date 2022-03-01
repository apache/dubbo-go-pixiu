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

package consistenthashing

import (
	"hash/crc32"
	"sort"
	"strconv"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type HashRing struct {
	nodes       map[uint32]*model.Endpoint
	replicaNum  int
	sortedNodes []uint32
}

func NewHashRing(nodes []*model.Endpoint, replicaNum int) *HashRing {
	h := new(HashRing)
	h.replicaNum = replicaNum
	h.nodes = make(map[uint32]*model.Endpoint)
	h.sortedNodes = []uint32{}
	h.addNodes(nodes)
	return h
}

func (h *HashRing) hashKey(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (h *HashRing) getNode(key string) *model.Endpoint {
	if len(h.nodes) == 0 {
		return nil
	}

	hashKey := h.hashKey(key)
	nodes := h.sortedNodes

	for _, node := range nodes {
		if hashKey < node {
			return h.nodes[node]
		}
	}

	return h.nodes[nodes[0]]
}

func (h *HashRing) addNode(node *model.Endpoint) {
	for i := 0; i < h.replicaNum; i++ {
		key := h.hashKey(strconv.Itoa(i) + node.ID)

		h.nodes[key] = node
		h.sortedNodes = append(h.sortedNodes, key)
	}

	sort.Slice(h.sortedNodes, func(i, j int) bool { return h.sortedNodes[i] < h.sortedNodes[j] })
}

func (h *HashRing) addNodes(masterNodes []*model.Endpoint) {
	for _, node := range masterNodes {
		h.addNode(node)
	}
}

// todo server down how to remove node ?
func (h *HashRing) removeNode(masterNodes string) {
	for i := 0; i < h.replicaNum; i++ {
		key := h.hashKey(strconv.Itoa(i) + masterNodes)
		delete(h.nodes, key)

		if index := h.getIndexForKey(key); index != -1 {
			h.sortedNodes = append(h.sortedNodes[:index], h.sortedNodes[index+1:]...)
		}
	}
}

func (h *HashRing) getIndexForKey(key uint32) int {

	for i, node := range h.sortedNodes {
		if node == key {
			return i
		}
	}

	return -1
}
