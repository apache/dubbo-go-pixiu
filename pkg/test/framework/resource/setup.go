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

package resource

// SetupFn is a function used for performing setup actions.
type SetupFn func(ctx Context) error

// ShouldSkipFn is a function used for performing skip actions; if it returns true a job is skipped
// Note: function may be called multiple times during the setup process.
type ShouldSkipFn func(ctx Context) bool
