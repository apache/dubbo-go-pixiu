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

func WithLookUpTable(factor, nodeCount int, function func(*LookUpTable)) {
	table := createTableWithNodes(factor, nodeCount)
	table.Populate()
	if function != nil {
		function(table)
	}
}

func createTableWithNodes(factor, nodeCount int) *LookUpTable {
	nodes := make([]string, 0, nodeCount)
	for i := 1; i <= nodeCount; i++ {
		nodes = append(nodes, createEndpoint(i))
	}
	return NewLookUpTable(factor, nodes)
}

func createEndpoint(id int) string {
	return fmt.Sprintf("192.168.1.%d:%d", id, 1000+id)
}

func TestLookUpTable_Populate(t *testing.T) {
	factors := []int{1, 2, 10, 100, 500}
	nodes := []int{0, 3, 5, 50, 250, 500}

	for _, factor := range factors {
		for _, node := range nodes {
			WithLookUpTable(factor, node, func(table *LookUpTable) {
				checkUnPopulateSlot(t, table)
				checkMissingEndpoint(t, table)
			})
		}
	}
}

func TestLookUpTable_Get(t *testing.T) {
	factor, nodeCount := 10, 10
	key := "/this/is/a/test"

	table := createTableWithNodes(factor, nodeCount)

	_, err := table.Get(key)
	assert.NotNil(t, err, "Got endpoint before populating")

	table.Populate()

	ep, err := table.Get(key)
	assert.Nil(t, err, "Fail to get endpoint")
	assert.NotEqual(t, "", ep, "Wrong endpoint")
}

func TestLookUpTable_Add(t *testing.T) {
	factor := 10
	testCases := []struct {
		nodeCount int
		add       []int
		wanted    []int // expected table size
	}{
		{
			nodeCount: 0,
			add:       []int{1, 2, 3, 4, 5},
			wanted:    []int{10, 20, 30, 40, 50},
		},
		{
			nodeCount: 3,
			add:       []int{4, 5, 6},
			wanted:    []int{40, 50, 60},
		},
		{
			nodeCount: 50,
			add:       []int{51, 52, 53, 54},
			wanted:    []int{510, 520, 530, 540},
		},
	}

	for id, tc := range testCases {
		WithLookUpTable(factor, tc.nodeCount, func(table *LookUpTable) {
			addEndpoints := make([]string, len(tc.add))
			for i, add := range tc.add {
				ep := createEndpoint(add)
				table.Add(ep)
				addEndpoints[i] = ep
				assert.Equal(t, tc.wanted[i], table.size,
					"Wrong size %d of grown table in case %d", table.size, id)
				checkUnPopulateSlot(t, table)
				checkMissingEndpoint(t, table)
			}
		})
	}
}

func TestLookUpTable_Remove(t *testing.T) {
	factor := 10
	testCases := []struct {
		nodeCount       int
		delete          []int
		wanted          []bool
		checkUnPopulate bool
		checkMissing    bool
	}{
		{
			nodeCount: 0,
			delete:    []int{2},
			wanted:    []bool{false},
		},
		{
			nodeCount: 1,
			delete:    []int{2, 1, 2},
			wanted:    []bool{false, true, false},
		},
		{
			nodeCount: 3,
			delete:    []int{1, 2, 3, 4},
			wanted:    []bool{true, true, true, false},
		},
		{
			nodeCount:       10,
			delete:          []int{4, 5, 6, 11},
			wanted:          []bool{true, true, true, false},
			checkMissing:    true,
			checkUnPopulate: true,
		},
		{
			nodeCount:       50,
			delete:          []int{54, 15, 1, 59, 5, 1},
			wanted:          []bool{false, true, true, false, true, false},
			checkMissing:    true,
			checkUnPopulate: true,
		},
	}

	for id, tc := range testCases {
		WithLookUpTable(factor, tc.nodeCount, func(table *LookUpTable) {
			for i, del := range tc.delete {
				ep := createEndpoint(del)
				ret := table.Remove(ep)
				assert.Equal(t, tc.wanted[i], ret,
					"Wrong state removing: %s in case %d", del, id)
				if tc.checkUnPopulate {
					checkUnPopulateSlot(t, table)
				}
				if tc.checkMissing {
					checkMissingEndpoint(t, table)
				}
				checkEndpointNotIn(t, table, ep)
			}
		})
	}
}

func TestLookUpTable_AddAndRemove(t *testing.T) {
	type change struct {
		change          int
		addOrRmv        bool // true is add
		wanted          int
		checkUnPopulate bool
		checkMissing    bool
		checkNotIn      bool
	}

	factor := 10
	testCases := []struct {
		name      string
		nodeCount int
		changes   []change
	}{
		{
			name:      "add from 0 then remove to 0",
			nodeCount: 0,
			changes: []change{
				{1, true, 10, true, true, false},
				{2, true, 20, true, true, false},
				{3, true, 30, true, true, false},
				{3, true, 30, true, true, false},
				{3, false, 30, false, false, true},
				{2, false, 30, false, false, true},
				{1, false, 30, false, false, true},
			},
		},
		{
			name:      "add and remove",
			nodeCount: 3,
			changes: []change{
				{4, true, 40, true, true, false},
				{1, false, 40, true, true, true},
				{5, true, 40, true, true, false},
				{6, true, 50, true, true, false},
				{3, false, 50, true, true, true},
				{2, false, 50, true, true, true},
				{1, false, 50, true, true, true},
			},
		},
		{
			name:      "remove to 0 then add",
			nodeCount: 3,
			changes: []change{
				{3, false, 30, true, true, true},
				{2, false, 30, true, true, true},
				{1, false, 30, false, false, true},
				{4, true, 30, true, true, false},
				{5, true, 30, true, true, false},
				{6, true, 30, true, true, false},
				{7, true, 40, true, true, false},
			},
		},
	}
	for _, tc := range testCases {
		WithLookUpTable(factor, tc.nodeCount, func(table *LookUpTable) {
			for _, ch := range tc.changes {
				ep := createEndpoint(ch.change)
				if ch.addOrRmv {
					table.Add(ep)
				} else {
					table.Remove(ep)
				}
				assert.Equal(t, ch.wanted, table.size,
					"Wrong table size %d in case: %s", table.size, tc.name)
				if ch.checkUnPopulate {
					checkUnPopulateSlot(t, table)
				}
				if ch.checkMissing {
					checkMissingEndpoint(t, table)
				}
				if ch.checkNotIn {
					checkEndpointNotIn(t, table, ep)
				}
			}
		})
	}
}

func checkUnPopulateSlot(t *testing.T, table *LookUpTable) {
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
