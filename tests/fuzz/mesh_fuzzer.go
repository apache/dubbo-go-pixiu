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
	"github.com/apache/dubbo-go-pixiu/pkg/config/mesh"
)

func FuzzParseMeshNetworks(data []byte) int {
	_, _ = mesh.ParseMeshNetworks(string(data))
	return 1
}

func FuzzValidateMeshConfig(data []byte) int {
	_, _ = mesh.ApplyMeshConfigDefaults(string(data))
	return 1
}
