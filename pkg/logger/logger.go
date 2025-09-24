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

package logger

import (
	"fmt"
	"os"
	"path"
	"strings"
)

import (
	perrors "github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
)

var control *logController

type pixiuLogger struct {
	*zap.SugaredLogger
	config *zap.Config
}

func init() {
	// only used in test/bootstrap; keep a sane default
	if control == nil {
		control = new(logController)
		InitLogger(nil)
	}
}

// PaddedCallerEncoder aligns caller path to a fixed width for prettier console output.
func PaddedCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {

	callerPath := caller.TrimmedPath()

	// Set a fixed length, and if the path is too short, add a space after it
	const fixedLength = 30
	if len(callerPath) < fixedLength {
		padding := strings.Repeat(" ", fixedLength-len(callerPath))
		callerPath = callerPath + padding
	}

	enc.AppendString(callerPath)
}

// helper: build with unified pixiu options.
// - always AddCaller + AddCallerSkip(2)
// - AddStacktrace(Error+) only when stacktrace is not disabled in cfg
func buildWithPixiuOptions(cfg *zap.Config) (*zap.Logger, error) {
	opts := []zap.Option{
		zap.AddCaller(),
		zap.AddCallerSkip(2),
	}
	if !cfg.DisableStacktrace {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	return cfg.Build(opts...)
}

// InitLog loads from YAML file; falls back to development defaults when file is absent/invalid.
func InitLog(logConfFile string) error {
	if logConfFile == "" {
		InitLogger(nil)
		return perrors.New("log configure file name is nil")
	}
	if path.Ext(logConfFile) != ".yml" {
		InitLogger(nil)
		return perrors.New(fmt.Sprintf("log configure file name %s suffix must be .yml", logConfFile))
	}

	confFileStream, err := os.ReadFile(logConfFile)
	if err != nil {
		InitLogger(nil)
		return perrors.New(fmt.Sprintf("os.ReadFile file:%s, error:%v", logConfFile, err))
	}

	conf := &zap.Config{}
	if err := yaml.UnmarshalYML(confFileStream, conf); err != nil {
		InitLogger(nil)
		return perrors.New(fmt.Sprintf("[Unmarshal]init pixiuLogger error: %v", err))
	}

	InitLogger(conf)

	return nil
}

// InitLogger initializes logger. Default is development-style (console, debug),
// but we force stacktrace to Error+ only, and enable caller with our custom encoder.
// If a config is supplied, we respect it and only normalize caller encoder and stacktrace threshold.
func InitLogger(conf *zap.Config) {
	var cfg zap.Config

	if conf == nil {
		// Default: development style
		cfg = zap.NewDevelopmentConfig()

		// Normalize/override keys & encoders for consistent console output
		cfg.EncoderConfig.TimeKey = "time"
		cfg.EncoderConfig.LevelKey = "level"
		cfg.EncoderConfig.NameKey = "pixiuLogger"
		cfg.EncoderConfig.CallerKey = "caller"
		cfg.EncoderConfig.MessageKey = "message"
		cfg.EncoderConfig.StacktraceKey = "stacktrace"

		cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		cfg.EncoderConfig.EncodeCaller = PaddedCallerEncoder

		// Keep console encoding & debug level as dev style implies.
		// cfg.Encoding = "console" // dev config already does this
	} else {
		cfg = *conf
		// Unify caller encoder regardless of YAML to keep alignment style
		cfg.EncoderConfig.EncodeCaller = PaddedCallerEncoder
	}

	z, err := buildWithPixiuOptions(&cfg)
	if err != nil {
		z = zap.NewNop()
	}
	l := &pixiuLogger{z.Sugar(), &cfg}
	control.updateLogger(l)
}

// SetLoggerLevel changes the level at runtime without rebuilding logger.
func SetLoggerLevel(level string) bool {
	return control.setLoggerLevel(level)
}

// HotReload rebuilds from a new zap.Config (e.g., re-read YAML).
func HotReload(conf *zap.Config) error {
	InitLogger(conf)
	return nil
}

// GetLogger exposes the current sugared logger.
func GetLogger() *pixiuLogger {
	return control.logger
}
