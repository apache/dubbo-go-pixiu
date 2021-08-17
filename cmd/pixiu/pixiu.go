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
	_ "net/http/pprof"
	"runtime"
	"strconv"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

import (
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/metadata/service/inmemory"
	"github.com/spf13/cobra"
)

var (
	// Version server version
	Version = "0.3.0"

	flagToLogLevel = map[string]string{
		"trace":    "TRACE",
		"debug":    "DEBUG",
		"info":     "INFO",
		"warning":  "WARN",
		"error":    "ERROR",
		"critical": "FATAL",
	}

	configPath    string
	apiConfigPath string
	logConfigPath string
	logLevel      string

	// CURRENTLY USELESS
	logFormat string

	limitCpus string

	// Currently default set to false, wait for up coming support
	initFromRemote = false
)

// main server run method
func main() {
	app := getRootCmd()

	// ignore error so we don't exit non-zero and break gfmrun README example tests
	_ = app.Execute()
}

func getRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "dubbogo server",
		Short: "Dubbogo server is a lightweight gateway.",
		Long: "dubbo-go-server is a gateway that mainly focuses on providing gateway solution to your Dubbo and RESTful \n" +
			"services. It supports HTTP-to-Dubbo and HTTP-to-HTTP proxy and more protocols will be supported in the near \n" +
			"future. \n" +
			"(c) " + strconv.Itoa(time.Now().Year()) + " Dubbogo",
		Version: Version,
	}

	rootCmd.AddCommand(gatewayCmd)
	rootCmd.AddCommand(sideCarCmd)

	return rootCmd
}

// initDefaultValue If not set both in args and env, set default values
func initDefaultValue() {
	if configPath == "" {
		configPath = constant.DefaultConfigPath
	}

	if apiConfigPath == "" {
		apiConfigPath = constant.DefaultApiConfigPath
	}

	if logConfigPath != "" {
		logConfigPath = constant.DefaultLogConfigPath
	}

	if logLevel == "" {
		logLevel = constant.DefaultLogLevel
	}

	if limitCpus == "" {
		limitCpus = constant.DefaultLimitCpus
	}

	if logFormat != "" {
		logFormat = constant.DefaultLogFormat
	}
}

// initLog
func initLog() error {
	err := logger.InitLog(logConfigPath)
	if err != nil {
		// cause `logger.InitLog` already handle init failed, so just use logger to log
		return err
	}

	if level, ok := flagToLogLevel[logLevel]; ok {
		logger.SetLoggerLevel(level)
	} else {
		logger.SetLoggerLevel(flagToLogLevel[constant.DefaultLogLevel])
		return fmt.Errorf("logLevel is invalid, set log level to default: %s", constant.DefaultLogLevel)
	}
	return nil
}

// initApiConfig return value of the bool is for the judgment of whether is a api meta data error, a kind of silly (?)
func initApiConfig() (*model.Bootstrap, bool, error) {
	bootstrap := config.Load(configPath)
	return bootstrap, false, nil
}

func initLimitCpus() error {
	limitCpuNumber, err := strconv.ParseInt(limitCpus, 10, 64)
	if err != nil {
		return err
	}
	if limitCpuNumber <= 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	} else {
		runtime.GOMAXPROCS(int(limitCpuNumber))
	}
	return nil
}
