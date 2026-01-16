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

package config

import (
	"os"
	"path/filepath"
	"strings"
	"time"
)

import (
	etcdv3 "github.com/dubbogo/gost/database/kv/etcd/v3"

	perrors "github.com/pkg/errors"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"gopkg.in/natefinch/lumberjack.v2"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
)

// Config holds all configuration for the admin server.
type Config struct {
	Server ServerConfig `yaml:"server" json:"server" mapstructure:"server"`
	Etcd   EtcdConfig   `yaml:"etcd" json:"etcd" mapstructure:"etcd"`
	MySQL  MySQLConfig  `yaml:"mysql" json:"mysql" mapstructure:"mysql"`
	Zap    ZapConfig    `yaml:"zap" json:"zap" mapstructure:"zap"`
	System SystemConfig `yaml:"system" json:"system" mapstructure:"system"`
	JWT    JWTConfig    `yaml:"jwt" json:"jwt" mapstructure:"jwt"`
}

// JWTConfig holds JWT authentication configuration.
type JWTConfig struct {
	SignKey     string `yaml:"sign-key" json:"sign-key" mapstructure:"sign-key"`
	ExpireHours int    `yaml:"expire-hours" json:"expire-hours" mapstructure:"expire-hours"`
	Issuer      string `yaml:"issuer" json:"issuer" mapstructure:"issuer"`
}

// GetSignKey returns the JWT sign key, with fallback to environment variable.
func (c *JWTConfig) GetSignKey() string {
	if c.SignKey != "" {
		// Check if it's an environment variable reference
		if len(c.SignKey) > 2 && c.SignKey[0] == '$' && c.SignKey[1] == '{' {
			envVar := c.SignKey[2 : len(c.SignKey)-1]
			if val := os.Getenv(envVar); val != "" {
				return val
			}
		}
		return c.SignKey
	}
	// Fallback to environment variable
	if val := os.Getenv("JWT_SIGN_KEY"); val != "" {
		return val
	}
	// Default fallback (should be changed in production)
	return "pixiu-admin-default-key-change-me"
}

// GetExpireHours returns the token expire hours with default.
func (c *JWTConfig) GetExpireHours() int {
	if c.ExpireHours > 0 {
		return c.ExpireHours
	}
	return 24 // Default 24 hours
}

// GetIssuer returns the JWT issuer with default.
func (c *JWTConfig) GetIssuer() string {
	if c.Issuer != "" {
		return c.Issuer
	}
	return "pixiu-admin"
}

// ServerConfig holds HTTP server configuration.
type ServerConfig struct {
	Address string `yaml:"address" json:"address" mapstructure:"address"`
	Port    int    `yaml:"port" json:"port" mapstructure:"port"`
}

// EtcdConfig holds etcd client configuration.
type EtcdConfig struct {
	Address string `yaml:"address" json:"address" mapstructure:"address"`
	Path    string `yaml:"path" json:"path" mapstructure:"path"`
}

// MySQLConfig holds MySQL database configuration.
type MySQLConfig struct {
	Host     string `yaml:"host" json:"host" mapstructure:"host"`
	Port     int    `yaml:"port" json:"port" mapstructure:"port"`
	Username string `yaml:"username" json:"username" mapstructure:"username"`
	Password string `yaml:"password" json:"password" mapstructure:"password"`
	Database string `yaml:"dbname" json:"dbname" mapstructure:"dbname"`
}

// ZapConfig holds zap logger configuration (compatible with existing config).
type ZapConfig struct {
	Level        string `yaml:"level" json:"level" mapstructure:"level"`
	Format       string `yaml:"format" json:"format" mapstructure:"format"`
	Prefix       string `yaml:"prefix" json:"prefix" mapstructure:"prefix"`
	Director     string `yaml:"director" json:"director" mapstructure:"director"`
	ShowLine     bool   `yaml:"show-line" json:"show-line" mapstructure:"show-line"`
	EncodeLevel  string `yaml:"encode-level" json:"encode-level" mapstructure:"encode-level"`
	LogInConsole bool   `yaml:"log-in-console" json:"log-in-console" mapstructure:"log-in-console"`
}

// SystemConfig holds system configuration.
type SystemConfig struct {
	Env    string `yaml:"env" json:"env" mapstructure:"env"`
	Addr   int    `yaml:"addr" json:"addr" mapstructure:"addr"`
	DBType string `yaml:"db-type" json:"db-type" mapstructure:"db-type"`
}

// LogConfig is an alias for ZapConfig for backward compatibility.
type LogConfig = ZapConfig

// BaseInfo represents basic gateway information.
type BaseInfo struct {
	Name           string `json:"name" yaml:"name"`
	Description    string `json:"description" yaml:"description"`
	PluginFilePath string `json:"pluginFilePath" yaml:"pluginFilePath"`
}

// Load reads configuration from the given file path.
func Load(path string) (*Config, error) {
	if path == "" {
		return nil, perrors.New("config file path is empty")
	}

	cfg := &Config{}
	if err := yaml.UnmarshalYMLConfig(path, cfg); err != nil {
		return nil, perrors.Wrapf(err, "failed to load config from %s", path)
	}

	return cfg, nil
}

// NewEtcdClient creates a new etcd client from configuration.
func NewEtcdClient(cfg EtcdConfig) (*etcdv3.Client, error) {
	client, err := etcdv3.NewConfigClientWithErr(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(20*time.Second),
		etcdv3.WithEndpoints(strings.Split(cfg.Address, ",")...),
	)
	if err != nil {
		return nil, perrors.Wrap(err, "failed to create etcd client")
	}
	return client, nil
}

// NewLogger creates a new zap logger from configuration.
func NewLogger(cfg ZapConfig) (*zap.Logger, error) {
	level := zapcore.InfoLevel
	switch strings.ToLower(cfg.Level) {
	case "debug":
		level = zapcore.DebugLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Handle color level encoder
	if strings.Contains(strings.ToLower(cfg.EncodeLevel), "color") {
		encoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	var encoder zapcore.Encoder
	if cfg.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	var cores []zapcore.Core

	// Console output
	if cfg.LogInConsole {
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	// File output
	if cfg.Director != "" {
		_ = os.MkdirAll(cfg.Director, 0755)
		logPath := filepath.Join(cfg.Director, "admin.log")
		writer := &lumberjack.Logger{
			Filename:  logPath,
			MaxSize:   100, // MB
			MaxAge:    30,  // days
			LocalTime: true,
			Compress:  true,
		}
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(writer), level))
	}

	if len(cores) == 0 {
		// Default to console if nothing configured
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), level))
	}

	core := zapcore.NewTee(cores...)

	opts := []zap.Option{}
	if cfg.ShowLine {
		opts = append(opts, zap.AddCaller())
	}

	return zap.New(core, opts...), nil
}
