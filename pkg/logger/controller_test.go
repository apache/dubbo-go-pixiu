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
	"testing"
)

import (
	"go.uber.org/zap/zapcore"
)

func TestParseLevelAndSet(t *testing.T) {
	cfg, _ := newDevConfigToFile(t)
	InitLogger(cfg)

	tests := []struct {
		in       string
		ok       bool
		zapLevel zapcore.Level
	}{
		{"debug", true, zapcore.DebugLevel},
		{"trace", true, zapcore.DebugLevel},
		{"INFO", true, zapcore.InfoLevel},
		{"Warn", true, zapcore.WarnLevel},
		{"Warning", true, zapcore.WarnLevel},
		{"error", true, zapcore.ErrorLevel},
		{"panic", true, zapcore.PanicLevel},
		{"dpanic", true, zapcore.DPanicLevel},
		{"fatal", true, zapcore.FatalLevel},
		{"critical", true, zapcore.FatalLevel},
		{"unknown", false, zapcore.InfoLevel}, // parseLevel default fallback to info
	}

	for _, tt := range tests {
		lvl := ParseLogLevel(tt.in)
		if lvl != tt.zapLevel {
			t.Fatalf("SetLoggerLevel(%q) want %v", tt.in, tt.ok)
		}
	}

	// make sure able to write
	GetLogger().Info("still alive")
	_ = GetLogger().Sync()
}

func TestInitLoggerNil(t *testing.T) {
	InitLogger(nil)
	if GetLogger() == nil || GetLogger().SugaredLogger == nil {
		t.Fatalf("GetLogger returned nil")
	}
	// retrigger InitLogger(nil)
	InitLogger(nil)
	GetLogger().Debug("dev init ok")
	_ = GetLogger().Sync()
}

// HotReload(nil) fallback to dev
func TestHotReloadNil(t *testing.T) {
	if err := HotReload(nil); err != nil {
		t.Fatalf("HotReload(nil) unexpected error: %v", err)
	}
	GetLogger().Warn("after hot reload nil")
	_ = GetLogger().Sync()
}

func TestLoggerBasicUsage(t *testing.T) {
	cfg, _ := newDevConfigToFile(t)
	InitLogger(cfg)

	log := GetLogger()
	log = &PixiuLogger{SugaredLogger: log.With("k", "v"), config: log.config}
	log.Infow("with fields", "a", 1)

	_ = log.Sync()
}
