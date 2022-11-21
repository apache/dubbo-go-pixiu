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

package kube

import (
	"fmt"
	"strings"
)

import (
	kubeCore "k8s.io/api/core/v1"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config/constants"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo"
)

type podSelector struct {
	Label string
	Value string
}

func (s podSelector) String() string {
	return s.Label + "=" + s.Value
}

func (s podSelector) MatchesPod(pod *kubeCore.Pod) bool {
	return pod.ObjectMeta.Labels[s.Label] == s.Value
}

func newPodSelector(cfg echo.Config) podSelector {
	label := "app"
	if cfg.DeployAsVM {
		label = constants.TestVMLabel
	}
	return podSelector{
		Label: label,
		Value: cfg.Service,
	}
}

func serviceAccount(cfg echo.Config) string {
	if cfg.ServiceAccount {
		return cfg.Service
	}
	if cfg.DeployAsVM {
		return "default"
	}
	return ""
}

// workloadHasSidecar returns true if the input endpoint is deployed with sidecar injected based on the config.
func workloadHasSidecar(cfg echo.Config, podName string) bool {
	// Match workload first.
	for _, w := range cfg.Subsets {
		if strings.HasPrefix(podName, fmt.Sprintf("%v-%v", cfg.Service, w.Version)) {
			return w.Annotations.GetBool(echo.SidecarInject) &&
				!strings.HasPrefix(w.Annotations.Get(echo.SidecarInjectTemplates), "grpc-")
		}
	}
	return true
}
