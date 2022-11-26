// Copyright Istio Authors.
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

package install

import (
	"fmt"
)

import (
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
)

import (
	"github.com/apache/dubbo-go-pixiu/istioctl/pkg/clioptions"
	"github.com/apache/dubbo-go-pixiu/istioctl/pkg/util/formatting"
	"github.com/apache/dubbo-go-pixiu/istioctl/pkg/verifier"
	"github.com/apache/dubbo-go-pixiu/operator/cmd/mesh"
	"github.com/apache/dubbo-go-pixiu/pkg/config/constants"
)

// NewVerifyCommand creates a new command for verifying Istio Installation Status
func NewVerifyCommand() *cobra.Command {
	var (
		kubeConfigFlags = &genericclioptions.ConfigFlags{
			Context:    strPtr(""),
			Namespace:  strPtr(""),
			KubeConfig: strPtr(""),
		}

		filenames      = []string{}
		istioNamespace string
		opts           clioptions.ControlPlaneOptions
		manifestsPath  string
	)
	verifyInstallCmd := &cobra.Command{
		Use:   "verify-install [-f <deployment or istio operator file>] [--revision <revision>]",
		Short: "Verifies Istio Installation Status",
		Long: `
verify-install verifies Istio installation status against the installation file
you specified when you installed Istio. It loops through all the installation
resources defined in your installation file and reports whether all of them are
in ready status. It will report failure when any of them are not ready.

If you do not specify an installation it will check for an IstioOperator resource
and will verify if pods and services defined in it are present.

Note: For verifying whether your cluster is ready for Istio installation, see
istioctl experimental precheck.
`,
		Example: `  # Verify that Istio is installed correctly via Istio Operator
  istioctl verify-install

  # Verify the deployment matches a custom Istio deployment configuration
  istioctl verify-install -f $HOME/istio.yaml

  # Verify the deployment matches the Istio Operator deployment definition
  istioctl verify-install --revision <canary>

  # Verify the installation of specific revision
  istioctl verify-install -r 1-9-0`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(filenames) > 0 && opts.Revision != "" {
				cmd.Println(cmd.UsageString())
				return fmt.Errorf("supply either a file or revision, but not both")
			}
			return nil
		},
		RunE: func(c *cobra.Command, args []string) error {
			installationVerifier, err := verifier.NewStatusVerifier(istioNamespace, manifestsPath,
				*kubeConfigFlags.KubeConfig, *kubeConfigFlags.Context, filenames, opts)
			if err != nil {
				return err
			}
			if formatting.IstioctlColorDefault(c.OutOrStdout()) {
				installationVerifier.Colorize()
			}
			return installationVerifier.Verify()
		},
	}

	flags := verifyInstallCmd.PersistentFlags()
	flags.StringVarP(&istioNamespace, "istioNamespace", "i", constants.IstioSystemNamespace,
		"Istio system namespace")
	kubeConfigFlags.AddFlags(flags)
	flags.StringSliceVarP(&filenames, "filename", "f", filenames, "Istio YAML installation file.")
	verifyInstallCmd.PersistentFlags().StringVarP(&manifestsPath, "manifests", "d", "", mesh.ManifestsFlagHelpStr)
	opts.AttachControlPlaneFlags(verifyInstallCmd)
	return verifyInstallCmd
}

func strPtr(val string) *string {
	return &val
}
