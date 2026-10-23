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
	"math"
	"testing"
)

func TestPermutationPositionPreservesUint32Arithmetic(t *testing.T) {
	cases := []struct {
		name               string
		size, offset, skip uint32
		overflow           bool
	}{
		{"small", 7, 6, 3, false},
		{"medium", 10007, 10006, 9999, false},
		{"overflow", 131071, 131070, 131070, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			p := permutation{offset: tc.offset, skip: tc.skip}
			overflowed := false
			for j := uint32(0); j < tc.size; j++ {
				// Preserve the old uint32 multiply/add contract, including wraparound.
				want := (tc.offset + j*tc.skip) % tc.size
				wide := uint64(tc.offset) + uint64(j)*uint64(tc.skip)
				if wide > math.MaxUint32 {
					overflowed = true
					// Ensure this fixture distinguishes wrapped from widened arithmetic.
					if uint32(wide%uint64(tc.size)) != want && p.position(int(j), tc.size) == uint32(wide%uint64(tc.size)) {
						t.Fatal("position silently widened uint32 arithmetic")
					}
				}
				if got := p.position(int(j), tc.size); got != want {
					t.Fatalf("position %d: got %d, want %d", j, got, want)
				}
			}
			if tc.overflow && !overflowed {
				t.Fatal("fixture did not exercise uint32 overflow")
			}

		})
	}
}
