/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	_ "net/http/pprof"
	"strconv"
	"time"
)

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cmd"
	_ "github.com/apache/dubbo-go-pixiu/pkg/pluginregistry"
)

const (
	// Version pixiu version
	Version = "1.0.0"
)

// main pixiu run method
func main() {
	app := getRootCmd()

	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Execute()
}

func getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dubbogo pixiu",
		Short: "Dubbogo pixiu is a lightweight gateway.",
		Long: "dubbo-go-pixiu is a gateway that mainly focuses on providing gateway solution to your Dubbo and RESTful \n" +
			"services. It supports HTTP-to-Dubbo and HTTP-to-HTTP proxy and more protocols will be supported in the near \n" +
			"future. \n" +
			"(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo",
		Version: Version,
	}

	rootCmd.AddCommand(cmd.GatewayCmd)
	rootCmd.AddCommand(cmd.SideCarCmd)

	return rootCmd
}
