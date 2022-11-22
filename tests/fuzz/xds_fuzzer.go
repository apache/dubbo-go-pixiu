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

package fuzz

import (
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/config/kube/crd"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pilot/pkg/xds"
)

func init() {
	testing.Init()
}

func FuzzXds(data []byte) int {
	t := &testing.T{}

	// Use crd.ParseInputs to verify data
	_, _, err := crd.ParseInputs(string(data))
	if err != nil {
		return 0
	}
	s := xds.NewFakeDiscoveryServer(t, xds.FakeOptions{
		ConfigString: string(data),
	})
	proxy := s.SetupProxy(&model.Proxy{
		ConfigNamespace: "app",
	})
	_ = s.Listeners(proxy)
	_ = s.Routes(proxy)
	return 1
}
