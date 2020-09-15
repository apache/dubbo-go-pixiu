package main

import (
	"runtime"
)

import (
	"github.com/urfave/cli"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/proxy"
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

	cmdStart = cli.Command{
		Name:  "start",
		Usage: "start dubbogo proxy",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:   "config, c",
				Usage:  "Load configuration from `FILE`",
				EnvVar: "DUBBOGO_PROXY_CONFIG",
				Value:  "configs/conf.yaml",
			},
			cli.StringFlag{
				Name:   "log-config, lc",
				Usage:  "Load log configuration from `FILE`",
				EnvVar: "LOG_FILE",
				Value:  "configs/log.yml",
			},
			cli.StringFlag{
				Name:   "api-config, ac",
				Usage:  "Load api description configuration from `FILE`",
				EnvVar: "API_FILE",
				Value:  "configs/api_config.yml",
			},
			cli.StringFlag{
				Name:   "log-level, l",
				Usage:  "dubbogo proxy log level, trace|debug|info|warning|error|critical",
				EnvVar: "LOG_LEVEL",
			},
			cli.StringFlag{
				Name:  "log-format, lf",
				Usage: "dubbogo proxy log format, currently useless",
			},
			cli.StringFlag{
				Name:  "limit-cpus, limc",
				Usage: "dubbogo proxy schedule threads count",
			},
		},
		Action: func(c *cli.Context) error {
			configPath := c.String("config")
			flagLogLevel := c.String("log-level")
			logConfPath := c.String("log-config")
			apiConfPath := c.String("api-config")

			bootstrap := config.Load(configPath)
			if logLevel, ok := flagToLogLevel[flagLogLevel]; ok {
				logger.SetLoggerLevel(logLevel)
			}

			logger.InitLog(logConfPath)

			limitCpus := c.Int("limit-cpus")
			if limitCpus <= 0 {
				runtime.GOMAXPROCS(runtime.NumCPU())
			} else {
				runtime.GOMAXPROCS(limitCpus)
			}

			proxy.Start(bootstrap)
			return nil
		},
	}

	cmdStop = cli.Command{
		Name:  "stop",
		Usage: "stop dubbogo proxy",
		Action: func(c *cli.Context) error {
			return nil
		},
	}

	cmdReload = cli.Command{
		Name:  "reload",
		Usage: "reconfiguration",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
)
