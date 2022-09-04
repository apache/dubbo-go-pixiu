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
	"os"
)

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
)

var (
	gatewayCmd = &cobra.Command{
		Use:   "gateway",
		Short: "Run dubbo go server in gateway mode",
	}

	startGatewayCmd = &cobra.Command{
		Use:     "start",
		Short:   "Start gateway",
		Version: Version,
		PreRun: func(cmd *cobra.Command, args []string) {
			initDefaultValue()

			err := deploy.initialize()
			if err != nil {
				panic(err)
			}
		},
		Run: func(cmd *cobra.Command, args []string) {

			err := deploy.start()
			if err != nil {
				panic(err)
			}
		},
	}
)

// init Init startCmd
func init() {
	startGatewayCmd.PersistentFlags().StringVarP(&configPath, constant.ConfigPathKey, "c", os.Getenv(constant.EnvDubbogoPixiuConfig), "Load configuration from `FILE`")
	startGatewayCmd.PersistentFlags().StringVarP(&apiConfigPath, constant.ApiConfigPathKey, "a", os.Getenv(constant.EnvDubbogoPixiuApiConfig), "Load api configuration from `FILE`")
	startGatewayCmd.PersistentFlags().StringVarP(&logConfigPath, constant.LogConfigPathKey, "g", os.Getenv(constant.EnvDubbogoPixiuLogConfig), "Load log configuration from `FILE`")
	startGatewayCmd.PersistentFlags().StringVarP(&logLevel, constant.LogLevelKey, "l", os.Getenv(constant.EnvDubbogoPixiuLogLevel), "dubbogo pixiu log level, trace|debug|info|warning|error|critical")
	startGatewayCmd.PersistentFlags().StringVarP(&limitCpus, constant.LimitCpusKey, "m", os.Getenv(constant.EnvDubbogoPixiuLimitCpus), "dubbogo pixiu schedule threads count")
	startGatewayCmd.PersistentFlags().StringVarP(&logFormat, constant.LogFormatKey, "f", os.Getenv(constant.EnvDubbogoPixiuLogFormat), "dubbogo pixiu log format, currently useless")

	gatewayCmd.AddCommand(startGatewayCmd)
}
