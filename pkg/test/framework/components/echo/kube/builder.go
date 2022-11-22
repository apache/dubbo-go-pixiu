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
	"github.com/hashicorp/go-multierror"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/cluster"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/resource"
)

func init() {
	echo.RegisterFactory(cluster.Kubernetes, build)
}

func build(ctx resource.Context, configs []echo.Config) (echo.Instances, error) {
	instances := make([]echo.Instance, len(configs))

	g := multierror.Group{}
	for i, cfg := range configs {
		i, cfg := i, cfg
		g.Go(func() (err error) {
			instances[i], err = newInstance(ctx, cfg)
			return
		})
	}

	err := g.Wait().ErrorOrNil()
	return instances, err
}
