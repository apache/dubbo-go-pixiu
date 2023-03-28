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

package maglev

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	Factors = []int{1, 2, 10, 100, 500}
	Nodes   = []int{0, 5, 50, 250, 500}
	Cases   []TestCase
)

type TestCase struct {
	name  string
	table *LookUpTable
}

func init() {
	Cases = make([]TestCase, 0, len(Factors)*len(Nodes))
	for _, factor := range Factors {
		for _, node := range Nodes {
			Cases = append(Cases, TestCase{
				name:  fmt.Sprintf("%d nodes with factor=%d", node, factor),
				table: createTableWithNodes(factor, node),
			})
		}
	}
}

func createTableWithNodes(factor, nodeCount int) *LookUpTable {
	nodes := make([]string, 0, nodeCount)
	for i := 1; i <= nodeCount; i++ {
		nodes = append(nodes, fmt.Sprintf("192.168.1.%d:%d", i, 1000+i))
	}
	return NewLookUpTable(factor, nodes)
}

func TestLookUpTable_Populate(t *testing.T) {
	for _, tc := range Cases {
		tc.table.Populate()
		for i, s := range tc.table.slot {
			assert.NotEqual(t, "", s,
				fmt.Sprintf("Found unpopulating slot in case: '%s' at slot: %d", tc.name, i))
		}
		for _, p := range tc.table.permutations {
			assert.Equal(t, true, p.hit > 0,
				fmt.Sprintf("Found missing endpoint: %s in case: '%s'", tc.table.slot[p.index], tc.name))
		}
	}
}
