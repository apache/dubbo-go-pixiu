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
	"github.com/spf13/cobra"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	pxruntime "github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/runtime"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/pluginregistry"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
)

var (
	// Version pixiu version
	Version = "0.5.1"

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

// deploy server deployment
var deploy = &DefaultDeployer{}

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

	rootCmd.AddCommand(gatewayCmd)
	rootCmd.AddCommand(sideCarCmd)

	return rootCmd
}

type DefaultDeployer struct {
	bootstrap    *model.Bootstrap
	configManger *config.ConfigManager
}

func (d *DefaultDeployer) initialize() error {

	err := initLog()
	if err != nil {
		logger.Warnf("[startGatewayCmd] failed to init logger, %s", err.Error())
	}

	// load Bootstrap config
	d.bootstrap = d.configManger.LoadBootConfig(configPath)
	if err != nil {
		panic(fmt.Errorf("[startGatewayCmd] failed to get api meta config, %s", err.Error()))
	}

	err = initLimitCpus()
	if err != nil {
		logger.Errorf("[startCmd] failed to get limit cpu number, %s", err.Error())
	}

	return err
}

func (d *DefaultDeployer) start() error {
	server.Start(d.bootstrap)
	return nil
}

func (d *DefaultDeployer) stop() error {
	//TODO implement me
	panic("implement me")
}

// initDefaultValue If not set both in args and env, set default values
func initDefaultValue() {
	if configPath == "" {
		configPath = constant.DefaultConfigPath
	}

	if apiConfigPath == "" {
		apiConfigPath = constant.DefaultApiConfigPath
	}

	if logConfigPath == "" {
		logConfigPath = constant.DefaultLogConfigPath
	}

	if logLevel == "" {
		logLevel = constant.DefaultLogLevel
	}

	if limitCpus == "" {
		limitCpus = constant.DefaultLimitCpus
	}

	if logFormat == "" {
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
func initApiConfig() (*model.Bootstrap, error) {
	bootstrap := config.Load(configPath)
	return bootstrap, nil
}

func initLimitCpus() error {
	limitCpuNumberFromEnv, err := strconv.ParseInt(limitCpus, 10, 64)
	if err != nil {
		return err
	}
	limitCpuNumber := int(limitCpuNumberFromEnv)
	if limitCpuNumber <= 0 {
		limitCpuNumber = pxruntime.GetCPUNum()
	}
	runtime.GOMAXPROCS(limitCpuNumber)
	logger.Infof("GOMAXPROCS set to %v", limitCpuNumber)
	return nil
}

func init() {
	deploy.configManger = config.NewConfigManger()
}
