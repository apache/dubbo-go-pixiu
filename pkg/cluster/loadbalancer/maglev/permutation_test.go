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
)

import (
	"github.com/stretchr/testify/assert"
)

const ConsistencyToleration = 3

func WithLookUpTable(tableSize, nodeCount int, function func(*LookUpTable)) {
	table := createTableWithNodes(tableSize, nodeCount)
	table.Populate()
	if function != nil {
		function(table)
	}
}

func createTableWithNodes(tableSize, nodeCount int) *LookUpTable {
	nodes := make([]string, 0, nodeCount)
	for i := 1; i <= nodeCount; i++ {
		nodes = append(nodes, createEndpoint(i))
	}
	table, _ := NewLookUpTable(tableSize, nodes)
	return table
}

func createEndpoint(id int) string {
	return fmt.Sprintf("192.168.1.%d:%d", id, 1000+id)
}

func TestCreateLookUpTableWithoutTableSize(t *testing.T) {
	testCases := []struct {
		nodeCount  int
		expectSize int
	}{
		{
			nodeCount:  0,
			expectSize: 307,
		},
		{
			nodeCount:  1,
			expectSize: 307,
		},
		{
			nodeCount:  3,
			expectSize: 307,
		},
		{
			nodeCount:  5,
			expectSize: 503,
		},
		{
			nodeCount:  15,
			expectSize: 1511,
		},
		{
			nodeCount:  100001,
			expectSize: 0,
		},
	}
	for i, tc := range testCases {
		table := createTableWithNodes(0, tc.nodeCount)
		if i == 5 {
			assert.Nil(t, table, "Got non nil table")
		} else {
			assert.Equal(t, tc.expectSize, table.size, "Wrong default table size")
		}
	}
}

func TestLookUpTable_Populate(t *testing.T) {
	testCases := []struct {
		nodeCount int
		size      int
	}{
		{
			nodeCount: 0,
			size:      11,
		},
		{
			nodeCount: 3,
			size:      313,
		},
		{
			nodeCount: 5,
			size:      503,
		},
		{
			nodeCount: 10,
			size:      1013,
		},
		{
			nodeCount: 50,
			size:      5021,
		},
		{
			nodeCount: 250,
			size:      25447,
		},
	}

	for _, tc := range testCases {
		WithLookUpTable(tc.size, tc.nodeCount, func(table *LookUpTable) {
			checkUnPopulateSlot(t, table)
			checkMissingEndpoint(t, table)
			checkConsistency(t, table)
		})
	}
}

func TestLookUpTable_Get(t *testing.T) {
	tableSize, nodeCount := 1033, 10
	key := "/this/is/a/test"

	table := createTableWithNodes(tableSize, nodeCount)

	_, err := table.Get(key)
	assert.NotNil(t, err, "Got endpoint before populating")

	table.Populate()

	ep, err := table.Get(key)
	assert.Nil(t, err, "Fail to get endpoint")
	assert.NotEqual(t, "", ep, "Wrong endpoint")
}

func TestLookUpTable_Add(t *testing.T) {
	testCases := []struct {
		nodeCount int
		size      int
		add       []int
	}{
		{
			nodeCount: 0,
			size:      11,
			add:       []int{1, 2, 3, 4, 5},
		},
		{
			nodeCount: 3,
			size:      313,
			add:       []int{4, 5, 6},
		},
		{
			nodeCount: 50,
			size:      5021,
			add:       []int{51, 52, 53, 54},
		},
	}

	for _, tc := range testCases {
		WithLookUpTable(tc.size, tc.nodeCount, func(table *LookUpTable) {
			for _, add := range tc.add {
				ep := createEndpoint(add)
				table.Add(ep)
				checkUnPopulateSlot(t, table)
				checkMissingEndpoint(t, table)
				checkConsistency(t, table)
			}
		})
	}
}

func TestLookUpTable_Remove(t *testing.T) {
	testCases := []struct {
		nodeCount int
		size      int
		delete    []int
		wanted    []bool
	}{
		{
			nodeCount: 0,
			size:      11,
			delete:    []int{2},
			wanted:    []bool{false},
		},
		{
			nodeCount: 1,
			size:      101,
			delete:    []int{2, 1, 2},
			wanted:    []bool{false, true, false},
		},
		{
			nodeCount: 3,
			size:      307,
			delete:    []int{1, 2, 3, 4},
			wanted:    []bool{true, true, true, false},
		},
		{
			nodeCount: 10,
			size:      1013,
			delete:    []int{4, 5, 6, 11},
			wanted:    []bool{true, true, true, false},
		},
		{
			nodeCount: 50,
			size:      5021,
			delete:    []int{54, 15, 1, 59, 5, 1},
			wanted:    []bool{false, true, true, false, true, false},
		},
	}

	for id, tc := range testCases {
		WithLookUpTable(tc.size, tc.nodeCount, func(table *LookUpTable) {
			for i, del := range tc.delete {
				ep := createEndpoint(del)
				ret := table.Remove(ep)
				assert.Equal(t, tc.wanted[i], ret,
					"Wrong state removing: %s in case %d", del, id)
				checkUnPopulateSlot(t, table)
				checkMissingEndpoint(t, table)
				checkEndpointNotIn(t, table, ep)
				checkConsistency(t, table)
			}
		})
	}
}

func checkUnPopulateSlot(t *testing.T, table *LookUpTable) {
	if table.endpointNum == 0 {
		return
	}

	for i, s := range table.slots {
		assert.NotEqual(t, "", s,
			"Unpopulating slot in table(node=%d, factor=%d) at slot: %d", table.endpointNum, table.size, i)
	}
}

func checkMissingEndpoint(t *testing.T, table *LookUpTable) {
	for _, p := range table.permutations {
		assert.Equal(t, true, p.hit > 0,
			"Missing endpoint: %s in table(node=%d, factor=%d)", table.slots[p.index], table.endpointNum, table.size)
	}
}

func checkEndpointNotIn(t *testing.T, table *LookUpTable, endpoint string) {
	for i, slot := range table.slots {
		assert.NotEqual(t, endpoint, slot,
			"Unexpect endpoint %s in slot %d", endpoint, i)
	}
}

func checkConsistency(t *testing.T, table *LookUpTable) {
	if table.endpointNum == 0 {
		return
	}

	avg := table.size / table.endpointNum
	dist := make(map[string]int)
	for _, slot := range table.slots {
		dist[slot]++
	}
	for k, v := range dist {
		assert.NotEqual(t, true, v > avg+ConsistencyToleration || v < avg-ConsistencyToleration,
			"%s with distributions %d not in %d +/- %d", k, v, avg, ConsistencyToleration)
	}
}
