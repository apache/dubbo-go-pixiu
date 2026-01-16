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
	"fmt"
	"os"
	"strconv"
	"time"
)

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/app"
)

const Version = "1.0.0"

var (
	configPath string
)

var (
	rootCmd = &cobra.Command{
		Use:   "dubbogo pixiu admin",
		Short: "Dubbogo pixiu admin is the control panel of pixiu gateway.",
		Long: "dubbgo pixiu admin is used to manage the visual interface of dubbogo pixiu, supporting login, user management, \n" +
			"plugin management, service configuration, API key management, interface authority management \n" +
			"(appKey authorization, interface authority, online and offline). \n" +
			"(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo",
		Version: Version,
		PreRun: func(cmd *cobra.Command, args []string) {
			initDefaultValue()
		},
		Run: func(cmd *cobra.Command, args []string) {
			srv, err := app.New(configPath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to create server: %v\n", err)
				os.Exit(1)
			}

			if err := srv.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "server error: %v\n", err)
				os.Exit(1)
			}
		},
	}
)

// init Init startCmd
func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", os.Getenv("DUBBOGO_PIXIU_CONFIG"), "Load configuration from `FILE`")
}

func getRootCmd() *cobra.Command {
	return rootCmd
}

func initDefaultValue() {
	if configPath == "" {
		configPath = "configs/admin_config.yaml"
	}
}

// main admin run method
func main() {
	app := getRootCmd()

	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Execute()
}
