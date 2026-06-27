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
	"os/signal"
	"strconv"
	"time"
)

import (
	"github.com/spf13/cobra"
)

import (
	config2 "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/core"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	configPath    string
	apiConfigPath string
)

var (
	rootCmd = &cobra.Command{
		Use:   "dubbogo pixiu admin",
		Short: "Dubbogo pixiu admin is the control panel of pixiu gateway.",
		Long: "dubbgo pixiu admin is used to manage the visual interface of dubbogo pixiu, supporting login, user management, \n" +
			"plugin management, service configuration, API key management, interface authority management \n" +
			"(appKey authorization, interface authority, online and offline). \n" +
			"(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo",
		Version: config2.Version,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			initDefaultValue()
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := config2.LoadAPIConfigFromFile(apiConfigPath)
			if err != nil {
				return fmt.Errorf("load admin config error: %w", err)
			}
			// Start server in a goroutine
			errCh := make(chan error, 1)
			go func() {
				errCh <- Start()
			}()
			// gracefully shutdown
			sigint := make(chan os.Signal, 1)
			signal.Notify(sigint, os.Interrupt)
			select {
			case <-sigint:
				Stop()
				return nil
			case err := <-errCh:
				Stop()
				return err
			}
		},
	}
)

// Start start init etcd client and start admin http server
func Start() error {
	return core.RunServer()
}

func Stop() {
	config2.CloseEtcdClient()
}

// init Init startCmd
func init() {
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", os.Getenv("DUBBOGO_PIXIU_CONFIG"), "Load configuration from `FILE`")
	rootCmd.PersistentFlags().StringVarP(&apiConfigPath, "api-config", "a", os.Getenv("DUBBOGO_PIXIU_API_CONFIG"), "Load api configuration from `FILE`")

}

func getRootCmd() *cobra.Command {
	return rootCmd
}

func initDefaultValue() {
	if configPath == "" {
		configPath = "configs/admin_config.yaml"
	}

	if apiConfigPath == "" {
		apiConfigPath = "configs/api_config.yaml"
	}
}

// main admin run method
func main() {
	app := getRootCmd()

	if err := app.Execute(); err != nil {
		logger.Errorf("command failed: %v", err)
		os.Exit(1)
	}
}
