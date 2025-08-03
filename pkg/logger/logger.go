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
	// only use in test case, so just load default config
	if control == nil {
		control = new(logController)
		InitLogger(nil)
	}
}

// PaddedCallerEncoder is a custom caller encoder that ensures that all file paths are displayed at the same length
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

// InitLog load from config path
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
	err = yaml.UnmarshalYML(confFileStream, conf)
	if err != nil {
		InitLogger(nil)
		return perrors.New(fmt.Sprintf("[Unmarshal]init pixiuLogger error: %v", err))
	}

	InitLogger(conf)

	return nil
}

func InitLogger(conf *zap.Config) {
	var zapLoggerConfig zap.Config
	if conf == nil {
		zapLoggerConfig = zap.NewDevelopmentConfig()
		zapLoggerEncoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "pixiuLogger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   PaddedCallerEncoder,
			// EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		zapLoggerConfig.EncoderConfig = zapLoggerEncoderConfig
	} else {
		zapLoggerConfig = *conf
		// Set up a custom encoder directly without checking the original value
		zapLoggerConfig.EncoderConfig.EncodeCaller = PaddedCallerEncoder
	}
	zapLogger, _ := zapLoggerConfig.Build(zap.AddCallerSkip(2))
	l := &pixiuLogger{zapLogger.Sugar(), &zapLoggerConfig}

	control.updateLogger(l)
}

// SetLoggerLevel safely changes the log level in a concurrent manner.
func SetLoggerLevel(level string) bool {
	return control.setLoggerLevel(level)
}

func HotReload(conf *zap.Config) error {
	InitLogger(conf)
	return nil
}
