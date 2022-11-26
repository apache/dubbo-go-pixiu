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

package clusters

import (
	adminapi "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/util/protomarshal"
)

// Wrapper is a wrapper around the Envoy Clusters
// It has extra helper functions for handling any/struct/marshal protobuf pain
type Wrapper struct {
	*adminapi.Clusters
}

// MarshalJSON is a custom marshaller to handle protobuf pain
func (w *Wrapper) MarshalJSON() ([]byte, error) {
	return protomarshal.Marshal(w)
}

// UnmarshalJSON is a custom unmarshaller to handle protobuf pain
func (w *Wrapper) UnmarshalJSON(b []byte) error {
	cd := &adminapi.Clusters{}
	err := protomarshal.UnmarshalAllowUnknown(b, cd)
	*w = Wrapper{cd}
	return err
}
