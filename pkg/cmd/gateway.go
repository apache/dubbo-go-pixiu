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

package cmd

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
)

import (
	"github.com/spf13/cobra"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	pxruntime "github.com/apache/dubbo-go-pixiu/pkg/common/runtime"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

var (
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

var (
	GatewayCmd = &cobra.Command{
		Use:   "gateway",
		Short: "Run dubbo go pixiu in gateway mode",
	}

	deploy = &DefaultDeployer{
		configManger: config.NewConfigManger(),
	}

	startGatewayCmd = &cobra.Command{
		Use:   "start",
		Short: "Start gateway",
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

	GatewayCmd.AddCommand(startGatewayCmd)
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
