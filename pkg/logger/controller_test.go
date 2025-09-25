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

func TestParseLevelAndSet(t *testing.T) {
	cfg, _ := newDevConfigToFile(t)
	InitLogger(cfg)

	tests := []struct {
		in       string
		ok       bool
		zapLevel string
	}{
		{"debug", true, "debug"},
		{"INFO", true, "info"},
		{"Warn", true, "warn"},
		{"error", true, "error"},
		{"panic", true, "panic"},
		{"fatal", true, "fatal"},
		{"unknown", false, "info"}, // parseLevel default fallback to info
	}

	for _, tt := range tests {
		ok := SetLoggerLevel(tt.in)
		if ok != tt.ok {
			t.Fatalf("SetLoggerLevel(%q) ok=%v, want %v", tt.in, ok, tt.ok)
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
	log = &pixiuLogger{SugaredLogger: log.With("k", "v"), config: log.config}
	log.Infow("with fields", "a", 1)

	_ = log.Sync()
}
