// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"strings"
	"testing"
)

import (
	"istio.io/pkg/version"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/kube"
)

var meshInfo = version.MeshInfo{
	{Component: "Pilot", Info: version.BuildInfo{Version: "1.0.0", GitRevision: "gitSHA123", GolangVersion: "go1.10", BuildStatus: "Clean", GitTag: "Tag"}},
	{Component: "Injector", Info: version.BuildInfo{Version: "1.0.1", GitRevision: "gitSHAabc", GolangVersion: "go1.10.1", BuildStatus: "Modified", GitTag: "OtherTag"}},
	{Component: "Citadel", Info: version.BuildInfo{Version: "1.2", GitRevision: "gitSHA321", GolangVersion: "go1.11.0", BuildStatus: "Clean", GitTag: "Tag"}},
}

func TestVersion(t *testing.T) {
	kubeClientWithRevision = mockExecClientVersionTest

	cases := []testCase{
		{ // case 0 client-side only, normal output
			args: strings.Split("version --remote=false --short=false", " "),
			// ignore the output, all output checks are now in istio/pkg
		},
		{ // case 1 remote, normal output
			args: strings.Split("version --remote=true --short=false --output=", " "),
			// ignore the output, all output checks are now in istio/pkg
		},
		{ // case 2 bogus arg
			args:           strings.Split("version --typo", " "),
			expectedOutput: "Error: unknown flag: --typo\n",
			wantException:  true,
		},
		{ // case 3 bogus output arg
			args:           strings.Split("version --output xyz", " "),
			expectedOutput: "Error: --output must be 'yaml' or 'json'\n",
			wantException:  true,
		},
		{ // case 4 remote, --revision flag
			args: strings.Split("version --remote=true --short=false --revision canary", " "),
			// ignore the output, all output checks are now in istio/pkg
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("case %d %s", i, strings.Join(c.args, " ")), func(t *testing.T) {
			verifyOutput(t, c)
		})
	}
}

func mockExecClientVersionTest(_, _ string, _ string) (kube.ExtendedClient, error) {
	return kube.MockClient{
		IstioVersions: &meshInfo,
	}, nil
}
