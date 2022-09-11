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

package main

import (
	"fmt"
	"os"

	// import all known client auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/apache/dubbo-go-pixiu/istioctl/cmd"
	"istio.io/pkg/log"
)

func main() {
	if err := cmd.ConfigAndEnvProcessing(); err != nil {
		fmt.Fprintf(os.Stderr, "Could not initialize: %v\n", err)
		exitCode := cmd.GetExitCode(err)
		os.Exit(exitCode)
	}

	rootCmd := cmd.GetRootCmd(os.Args[1:])

	log.EnableKlogWithCobra()

	if err := rootCmd.Execute(); err != nil {
		exitCode := cmd.GetExitCode(err)
		os.Exit(exitCode)
	}
}
