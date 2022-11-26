//go:build integ
// +build integ

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

package mtlsk8sca

import (
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/istio"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/label"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/resource"
	"github.com/apache/dubbo-go-pixiu/tests/integration/security/util"
)

var (
	inst istio.Instance
	apps = &util.EchoDeployments{}
)

func TestMain(m *testing.M) {
	// nolint: staticcheck
	framework.
		NewSuite(m).
		RequireSingleCluster().
		RequireMultiPrimary().
		Label(label.CustomSetup).
		// https://github.com/istio/istio/issues/22161. 1.22 drops support for legacy-unknown signer
		RequireMaxVersion(21).
		Setup(istio.Setup(&inst, setupConfig)).
		Setup(func(ctx resource.Context) error {
			// TODO: due to issue https://github.com/istio/istio/issues/25286,
			// currently VM does not work in this test
			return util.SetupApps(ctx, inst, apps, false)
		}).
		Run()
}

func setupConfig(_ resource.Context, cfg *istio.Config) {
	if cfg == nil {
		return
	}
	cfg.Values["global.pilotCertProvider"] = "kubernetes"
	cfg.DeployEastWestGW = false
}
