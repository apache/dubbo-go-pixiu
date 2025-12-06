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
	"os"
)

import (
	"controllers/internal/controller/config"

	"controllers/internal/manager"

	"github.com/go-logr/zapr"

	"github.com/spf13/cobra"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	ctrl "sigs.k8s.io/controller-runtime"
)

func GetRootCmd() *cobra.Command {
	root := newPixiuGatewayController()
	return root
}

func newPixiuGatewayController() *cobra.Command {
	var configPath string
	cmd := &cobra.Command{
		Use:  "pixiu-gateway-controller [command]",
		Long: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := config.ControllerConfig
			if configPath != "" {
				c, err := config.GetConfigFromFile(configPath)
				if err != nil {
					// If config file doesn't exist, log warning and continue with default config
					// This allows the controller to start even if config file is missing
					cmd.Printf("Warning: failed to load config from %s: %v, using default config\n", configPath, err)
				} else {
					cfg = c
					config.SetControllerConfig(c)
				}
			}
			// Command line flags are already bound to config.ControllerConfig,
			// so cfg (which is config.ControllerConfig) already has the updated values

			logLevel, err := zapcore.ParseLevel(cfg.LogLevel)
			if err != nil {
				return err
			}

			// controllers log
			core := zapcore.NewCore(
				zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
				zapcore.AddSync(zapcore.Lock(os.Stderr)),
				logLevel,
			)
			logger := zapr.NewLogger(zap.New(core, zap.AddCaller()))

			logger.Info("controller initialized", "configuration", cfg)
			ctrl.SetLogger(logger.WithName("controller-runtime"))

			ctx := ctrl.LoggerInto(cmd.Context(), logger)
			return manager.Run(ctx, logger)
		},
	}

	cmd.Flags().StringVarP(&configPath, "config-path", "c", "", "configuration file path for pixiu-gateway-controller")
	cmd.Flags().StringVar(&config.ControllerConfig.MetricsAddr, "metrics-bind-address", "0", "The address the metrics endpoint binds to. "+
		"Use :8443 for HTTPS or :8080 for HTTP, or leave as 0 to disable the metrics service.")
	cmd.Flags().StringVar(&config.ControllerConfig.ProbeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	cmd.Flags().StringVar(&config.ControllerConfig.LogLevel, "log-level", config.DefaultLogLevel, "The log level for pixiu-gateway-controller")
	cmd.Flags().StringVar(&config.ControllerConfig.ControllerName, "controller-name", config.DefaultControllerName, "The name of the controller")

	return cmd
}
