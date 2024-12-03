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

package model

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Log struct {
	Level             string                 `json:"level" yaml:"level"`
	Development       bool                   `json:"development" yaml:"development"`
	DisableCaller     bool                   `json:"disableCaller" yaml:"disableCaller"`
	DisableStacktrace bool                   `json:"disableStacktrace" yaml:"disableStacktrace"`
	Sampling          SamplingConfig         `json:"sampling" yaml:"sampling"`
	Encoding          string                 `json:"encoding" yaml:"encoding"`
	EncoderConfig     EncoderConfig          `json:"encoderConfig" yaml:"encoderConfig"`
	OutputPaths       []string               `json:"outputPaths" yaml:"outputPaths"`
	ErrorOutputPaths  []string               `json:"errorOutputPaths" yaml:"errorOutputPaths"`
	InitialFields     map[string]interface{} `json:"initialFields" yaml:"initialFields"`
}

func (l *Log) Build() *zap.Config {
	lvl, err := zap.ParseAtomicLevel(l.Level)
	if err != nil {
		logger.Errorf("failed parse %s to zap.AtomicLevel", l.Level)
	}

	return &zap.Config{
		Level:             lvl,
		Development:       l.Development,
		DisableCaller:     l.DisableCaller,
		DisableStacktrace: l.DisableStacktrace,
		Sampling:          l.Sampling.build(),
		Encoding:          l.Encoding,
		EncoderConfig:     l.EncoderConfig.build(),
		OutputPaths:       l.OutputPaths,
		ErrorOutputPaths:  l.ErrorOutputPaths,
		InitialFields:     l.InitialFields,
	}
}

type SamplingConfig struct {
	Initial    int                                           `json:"initial" yaml:"initial"`
	Thereafter int                                           `json:"thereafter" yaml:"thereafter"`
	Hook       func(zapcore.Entry, zapcore.SamplingDecision) `json:"-" yaml:"-"`
}

func (s *SamplingConfig) build() *zap.SamplingConfig {
	return &zap.SamplingConfig{
		Initial:    s.Initial,
		Thereafter: s.Thereafter,
	}
}

type EncoderConfig struct {
	MessageKey       string `json:"messageKey" yaml:"messageKey"`
	LevelKey         string `json:"levelKey" yaml:"levelKey"`
	TimeKey          string `json:"timeKey" yaml:"timeKey"`
	NameKey          string `json:"nameKey" yaml:"nameKey"`
	CallerKey        string `json:"callerKey" yaml:"callerKey"`
	FunctionKey      string `json:"functionKey" yaml:"functionKey"`
	StacktraceKey    string `json:"stacktraceKey" yaml:"stacktraceKey"`
	SkipLineEnding   bool   `json:"skipLineEnding" yaml:"skipLineEnding"`
	LineEnding       string `json:"lineEnding" yaml:"lineEnding"`
	EncodeLevel      string `json:"levelEncoder" yaml:"levelEncoder"`
	EncodeTime       string `json:"timeEncoder" yaml:"timeEncoder"`
	EncodeDuration   string `json:"durationEncoder" yaml:"durationEncoder"`
	EncodeCaller     string `json:"callerEncoder" yaml:"callerEncoder"`
	EncodeName       string `json:"nameEncoder" yaml:"nameEncoder"`
	ConsoleSeparator string `json:"consoleSeparator" yaml:"consoleSeparator"`
}

func (e *EncoderConfig) build() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:       e.MessageKey,
		LevelKey:         e.LevelKey,
		TimeKey:          e.TimeKey,
		NameKey:          e.NameKey,
		CallerKey:        e.CallerKey,
		FunctionKey:      e.FunctionKey,
		StacktraceKey:    e.StacktraceKey,
		SkipLineEnding:   e.SkipLineEnding,
		LineEnding:       e.LineEnding,
		EncodeLevel:      e.unmarshalLevelEncoder(),
		EncodeTime:       e.unmarshalTimeEncoder(),
		EncodeDuration:   e.unmarshalDurationEncoder(),
		EncodeCaller:     e.unmarshalCallerEncoder(),
		EncodeName:       e.unmarshalNameEncoder(),
		ConsoleSeparator: e.ConsoleSeparator,
	}
}

func (e *EncoderConfig) unmarshalLevelEncoder() zapcore.LevelEncoder {
	switch e.EncodeLevel {
	case "capital":
		return zapcore.CapitalLevelEncoder
	case "capitalColor":
		return zapcore.CapitalColorLevelEncoder
	case "color":
		return zapcore.LowercaseColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

func (e *EncoderConfig) unmarshalTimeEncoder() zapcore.TimeEncoder {
	switch e.EncodeTime {
	case "rfc3339nano", "RFC3339Nano":
		return zapcore.RFC3339NanoTimeEncoder
	case "rfc3339", "RFC3339":
		return zapcore.RFC3339TimeEncoder
	case "iso8601", "ISO8601":
		return zapcore.ISO8601TimeEncoder
	case "millis":
		return zapcore.EpochMillisTimeEncoder
	case "nanos":
		return zapcore.EpochNanosTimeEncoder
	default:
		return zapcore.EpochTimeEncoder
	}
}

func (e *EncoderConfig) unmarshalDurationEncoder() zapcore.DurationEncoder {
	switch e.EncodeDuration {
	case "string":
		return zapcore.StringDurationEncoder
	case "nanos":
		return zapcore.NanosDurationEncoder
	case "ms":
		return zapcore.MillisDurationEncoder
	default:
		return zapcore.SecondsDurationEncoder
	}
}

func (e *EncoderConfig) unmarshalCallerEncoder() zapcore.CallerEncoder {
	switch e.EncodeCaller {
	case "full":
		return zapcore.FullCallerEncoder
	default:
		return zapcore.ShortCallerEncoder
	}
}

func (e *EncoderConfig) unmarshalNameEncoder() zapcore.NameEncoder {
	switch e.EncodeName {
	case "full":
		return zapcore.FullNameEncoder
	default:
		return zapcore.FullNameEncoder
	}
}
