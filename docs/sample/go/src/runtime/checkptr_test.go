// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package runtime_test

import (
	"internal/testenv"
	"os/exec"
	"strings"
	"testing"
)

func TestCheckPtr(t *testing.T) {
	t.Parallel()
	testenv.MustHaveGoRun(t)

	exe, err := buildTestProg(t, "testprog", "-gcflags=all=-d=checkptr=1")
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		cmd  string
		want string
	}{
		{"CheckPtrAlignmentPtr", "fatal error: checkptr: unsafe pointer conversion\n"},
		{"CheckPtrAlignmentNoPtr", ""},
		{"CheckPtrArithmetic", "fatal error: checkptr: unsafe pointer arithmetic\n"},
		{"CheckPtrArithmetic2", "fatal error: checkptr: unsafe pointer arithmetic\n"},
		{"CheckPtrSize", "fatal error: checkptr: unsafe pointer conversion\n"},
		{"CheckPtrSmall", "fatal error: checkptr: unsafe pointer arithmetic\n"},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.cmd, func(t *testing.T) {
			t.Parallel()
			got, err := testenv.CleanCmdEnv(exec.Command(exe, tc.cmd)).CombinedOutput()
			if err != nil {
				t.Log(err)
			}
			if tc.want == "" {
				if len(got) > 0 {
					t.Errorf("output:\n%s\nwant no output", got)
				}
				return
			}
			if !strings.HasPrefix(string(got), tc.want) {
				t.Errorf("output:\n%s\n\nwant output starting with: %s", got, tc.want)
			}
		})
	}
}
