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

package deployment

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

import (
	"github.com/hashicorp/go-multierror"
	"golang.org/x/sync/errgroup"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/test"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo/common/ports"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/echo/deployment"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/components/namespace"
	"github.com/apache/dubbo-go-pixiu/pkg/test/framework/resource"
)

// Config for new echo deployment.
type Config struct {
	// Echos is the target Echos for the newly created echo apps. If nil, a new Echos
	// instance will be created.
	Echos *Echos

	// NamespaceCount indicates the number of echo namespaces to be generated.
	// Ignored if Namespaces is non-empty. Defaults to 1.
	NamespaceCount int

	// Namespaces is the user-provided list of echo namespaces. If empty, NamespaceCount
	// namespaces will be generated.
	Namespaces []namespace.Instance

	// NoExternalNamespace if true, no external namespace will be generated and no external echo
	// instance will be deployed. Ignored if ExternalNamespace is non-nil.
	NoExternalNamespace bool

	// ExternalNamespace the namespace to use for the external deployment. If nil, a namespace
	// will be generated unless NoExternalNamespace is specified.
	ExternalNamespace namespace.Instance
}

func (c *Config) fillDefaults(ctx resource.Context) error {
	// Create the namespaces concurrently.
	g, _ := errgroup.WithContext(context.TODO())

	if c.Echos == nil {
		c.Echos = &Echos{}
	}

	if len(c.Namespaces) > 0 {
		c.NamespaceCount = len(c.Namespaces)
	} else if c.NamespaceCount <= 0 {
		c.NamespaceCount = 1
	}

	// Create the echo namespaces.
	if len(c.Namespaces) == 0 {
		c.Namespaces = make([]namespace.Instance, c.NamespaceCount)
		if c.NamespaceCount == 1 {
			// If only using a single namespace, preserve the "echo" prefix.
			g.Go(func() (err error) {
				c.Namespaces[0], err = namespace.New(ctx, namespace.Config{
					Prefix: "echo",
					Inject: true,
				})
				return
			})
		} else {
			for i := 0; i < c.NamespaceCount; i++ {
				i := i
				g.Go(func() (err error) {
					c.Namespaces[i], err = namespace.New(ctx, namespace.Config{
						Prefix: fmt.Sprintf("echo%d", i+1),
						Inject: true,
					})
					return
				})
			}
		}
	}

	// Create the external namespace, if necessary.
	if c.ExternalNamespace == nil && !c.NoExternalNamespace {
		g.Go(func() (err error) {
			c.ExternalNamespace, err = namespace.New(ctx, namespace.Config{
				Prefix: "external",
				Inject: false,
			})
			return
		})
	}

	// Wait for the namespaces to be created.
	return g.Wait()
}

// SingleNamespaceView is a simplified view of Echos for tests that only require a single namespace.
type SingleNamespaceView struct {
	// Include the echos at the top-level, so there is no need for accessing sub-structures.
	EchoNamespace

	// External (out-of-mesh) deployments
	External External

	// All echo instances
	All echo.Services
}

// TwoNamespaceView is a simplified view of Echos for tests that require 2 namespaces.
type TwoNamespaceView struct {
	// Ns1 contains the echo deployments in the first namespace
	Ns1 EchoNamespace

	// Ns2 contains the echo deployments in the second namespace
	Ns2 EchoNamespace

	// Ns1AndNs2 contains just the echo services in Ns1 and Ns2 (excludes External).
	Ns1AndNs2 echo.Services

	// External (out-of-mesh) deployments
	External External

	// All echo instances
	All echo.Services
}

// Echos is a common set of echo deployments to support integration testing.
type Echos struct {
	// NS is the list of echo namespaces.
	NS []EchoNamespace

	// External (out-of-mesh) deployments
	External External

	// All echo instances.
	All echo.Services
}

// New echo deployment with the given configuration.
func New(ctx resource.Context, cfg Config) (*Echos, error) {
	if err := cfg.fillDefaults(ctx); err != nil {
		return nil, err
	}

	apps := cfg.Echos
	apps.NS = make([]EchoNamespace, len(cfg.Namespaces))
	for i, ns := range cfg.Namespaces {
		apps.NS[i].Namespace = ns
	}
	apps.External.Namespace = cfg.ExternalNamespace

	builder := deployment.New(ctx).WithClusters(ctx.Clusters()...)
	for _, n := range apps.NS {
		builder = n.build(ctx, builder)
	}

	if !cfg.NoExternalNamespace {
		builder = apps.External.build(builder)
	}

	echos, err := builder.Build()
	if err != nil {
		return nil, err
	}

	apps.All = echos.Services()

	g := multierror.Group{}
	for i := 0; i < len(apps.NS); i++ {
		i := i
		g.Go(func() error {
			return apps.NS[i].loadValues(ctx, echos, apps)
		})
	}

	if !cfg.NoExternalNamespace {
		g.Go(func() error {
			return apps.External.loadValues(echos)
		})
	}

	if err := g.Wait().ErrorOrNil(); err != nil {
		return nil, err
	}

	return apps, nil
}

// NewOrFail calls New and fails if an error is returned.
func NewOrFail(t test.Failer, ctx resource.Context, cfg Config) *Echos {
	t.Helper()
	out, err := New(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// SingleNamespaceView converts this Echos into a SingleNamespaceView.
func (d Echos) SingleNamespaceView() SingleNamespaceView {
	return SingleNamespaceView{
		EchoNamespace: d.NS[0],
		External:      d.External,
		All:           d.NS[0].All.Append(d.External.All.Services()),
	}
}

// TwoNamespaceView converts this Echos into a TwoNamespaceView.
func (d Echos) TwoNamespaceView() TwoNamespaceView {
	ns1AndNs2 := d.NS[0].All.Append(d.NS[1].All)
	return TwoNamespaceView{
		Ns1:       d.NS[0],
		Ns2:       d.NS[1],
		Ns1AndNs2: ns1AndNs2,
		External:  d.External,
		All:       ns1AndNs2.Append(d.External.All.Services()),
	}
}

func (d Echos) namespaces(excludes ...namespace.Instance) []string {
	var out []string
	for _, n := range d.NS {
		include := true
		for _, e := range excludes {
			if n.Namespace.Name() == e.Name() {
				include = false
				break
			}
		}
		if include {
			out = append(out, n.Namespace.Name())
		}
	}

	sort.Strings(out)
	return out
}

func serviceEntryPorts() []echo.Port {
	var res []echo.Port
	for _, p := range ports.All().GetServicePorts() {
		if strings.HasPrefix(p.Name, "auto") {
			// The protocol needs to be set in common.EchoPorts to configure the echo deployment
			// But for service entry, we want to ensure we set it to "" which will use sniffing
			p.Protocol = ""
		}
		res = append(res, p)
	}
	return res
}

// SetupSingleNamespace calls Setup and returns a SingleNamespaceView.
func SetupSingleNamespace(view *SingleNamespaceView) resource.SetupFn {
	return func(ctx resource.Context) error {
		// Perform a setup with 1 namespace.
		var apps Echos
		if err := Setup(&apps, Config{NamespaceCount: 1})(ctx); err != nil {
			return err
		}

		// Store the view.
		*view = apps.SingleNamespaceView()
		return nil
	}
}

// SetupTwoNamespaces calls Setup and returns a TwoNamespaceView.
func SetupTwoNamespaces(view *TwoNamespaceView) resource.SetupFn {
	return func(ctx resource.Context) error {
		// Perform a setup with 2 namespaces.
		var apps Echos
		if err := Setup(&apps, Config{NamespaceCount: 2})(ctx); err != nil {
			return err
		}

		// Store the view.
		*view = apps.TwoNamespaceView()
		return nil
	}
}

// Setup function for writing to a global deployment variable.
func Setup(apps *Echos, cfg Config) resource.SetupFn {
	return func(ctx resource.Context) error {
		// Set the target for the deployments.
		cfg.Echos = apps

		_, err := New(ctx, cfg)
		if err != nil {
			return err
		}

		return nil
	}
}

// TODO(nmittler): should ctx.Settings().Skip(echo.Delta) do all of this?
func skipDeltaXDS(ctx resource.Context) bool {
	return ctx.Settings().Skip(echo.Delta) || !ctx.Settings().Revisions.AtLeast("1.12")
}
