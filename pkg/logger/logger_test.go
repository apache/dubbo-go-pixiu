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
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newDevConfigToFile(t *testing.T) (*zap.Config, string) {
	t.Helper()
	dir := t.TempDir()
	out := filepath.Join(dir, "zap.log")

	cfg := zap.NewDevelopmentConfig()
	// output to a temp file
	cfg.OutputPaths = []string{out}
	cfg.ErrorOutputPaths = []string{out}
	cfg.EncoderConfig.StacktraceKey = "stacktrace"
	cfg.EncoderConfig.CallerKey = "caller"

	return &cfg, out
}

func readAll(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s failed: %v", path, err)
	}
	return string(b)
}

func TestDisableStacktraceTrue_NoStackEvenOnError(t *testing.T) {
	cfg, out := newDevConfigToFile(t)
	cfg.DisableStacktrace = true // YAML equals to：disableStacktrace: true
	InitLogger(cfg)

	log := GetLogger()
	log.Error("boom")
	_ = log.Sync()

	got := readAll(t, out)
	line := ""
	for _, l := range strings.Split(got, "\n") {
		if strings.Contains(l, "boom") {
			line = l
			break
		}
	}
	if line == "" {
		t.Fatalf("error line not found")
	}
	if strings.Contains(strings.ToLower(line), "stacktrace") {
		t.Fatalf("disableStacktrace=true: must NOT output stacktrace, got:\n%s", line)
	}
}

func TestSetLoggerLevel_DoesNotRebuildAndTakesEffect(t *testing.T) {
	cfg, out := newDevConfigToFile(t)
	InitLogger(cfg)

	before := GetLogger().SugaredLogger

	// dynamic set to error
	ok := SetLoggerLevel("error")
	if !ok {
		t.Fatalf("SetLoggerLevel returned false")
	}

	after := GetLogger().SugaredLogger
	if before != after {
		t.Fatalf("SetLoggerLevel should NOT rebuild logger; pointer changed: %p -> %p", before, after)
	}

	// write new：info should not appear，error should appear
	log := GetLogger()
	log.Info("info should be filtered")
	log.Error("error should appear")
	_ = log.Sync()

	got := readAll(t, out)
	if strings.Contains(got, "info should be filtered") {
		t.Fatalf("info should NOT appear after level set to error:\n%s", got)
	}
	if !strings.Contains(got, "error should appear") {
		t.Fatalf("error should appear but not found:\n%s", got)
	}
}

func TestHotReload_RebuildsAndSwitchesSink(t *testing.T) {
	// cfg1 -> out1
	cfg1, out1 := newDevConfigToFile(t)
	InitLogger(cfg1)
	l1 := GetLogger().SugaredLogger

	GetLogger().Info("hello-1")
	_ = GetLogger().Sync()

	// cfg2 -> out2（new sink）
	cfg2, out2 := newDevConfigToFile(t)
	// to split it，set lever to info
	HotReload(cfg2)
	l2 := GetLogger().SugaredLogger

	if l1 == l2 {
		t.Fatalf("HotReload should rebuild logger; got same pointer %p", l1)
	}

	GetLogger().Info("hello-2")
	_ = GetLogger().Sync()

	got1 := readAll(t, out1)
	got2 := readAll(t, out2)

	if !strings.Contains(got1, "hello-1") {
		t.Fatalf("out1 should contain hello-1 but not found:\n%s", got1)
	}
	if strings.Contains(got1, "hello-2") {
		t.Fatalf("out1 should NOT contain hello-2 after reload:\n%s", got1)
	}
	if !strings.Contains(got2, "hello-2") {
		t.Fatalf("out2 should contain hello-2 but not found:\n%s", got2)
	}
}

func TestPaddedCallerEncoder_FixedWidthAtLeast30(t *testing.T) {
	caller := zapcore.EntryCaller{
		Defined: true,
		File:    "a/b/c.go",
		Line:    7,
	}
	collector := &stringCollector{}
	PaddedCallerEncoder(caller, collector)

	if len(collector.items) == 0 {
		t.Fatalf("collector got no items")
	}
	got := collector.items[0]
	if len(got) < 30 {
		t.Fatalf("caller not padded to >=30, got len=%d val=%q", len(got), got)
	}
}

// -------------------------- helpers --------------------------

type stringCollector struct {
	items []string
}

// mock zapcore.PrimitiveArrayEncoder
func (s *stringCollector) AppendString(v string) { s.items = append(s.items, v) }

func (s *stringCollector) AppendBool(bool)                      {}
func (s *stringCollector) AppendByteString([]byte)              {}
func (s *stringCollector) AppendComplex128(complex128)          {}
func (s *stringCollector) AppendComplex64(complex64)            {}
func (s *stringCollector) AppendDuration(time.Duration)         {}
func (s *stringCollector) AppendFloat64(float64)                {}
func (s *stringCollector) AppendFloat32(float32)                {}
func (s *stringCollector) AppendInt(int)                        {}
func (s *stringCollector) AppendInt64(int64)                    {}
func (s *stringCollector) AppendInt32(int32)                    {}
func (s *stringCollector) AppendInt16(int16)                    {}
func (s *stringCollector) AppendInt8(int8)                      {}
func (s *stringCollector) AppendTime(time.Time)                 {}
func (s *stringCollector) AppendUint(uint)                      {}
func (s *stringCollector) AppendUint64(uint64)                  {}
func (s *stringCollector) AppendUint32(uint32)                  {}
func (s *stringCollector) AppendUint16(uint16)                  {}
func (s *stringCollector) AppendUint8(uint8)                    {}
func (s *stringCollector) AppendUintptr(uintptr)                {}
func (s *stringCollector) AppendReflected(any)                  {}
func (s *stringCollector) AppendArray(zapcore.ArrayMarshaler)   {}
func (s *stringCollector) AppendObject(zapcore.ObjectMarshaler) {}
func (s *stringCollector) AppendBinary([]byte)                  {}
func (s *stringCollector) AppendComplex(complex128)             {}
func (s *stringCollector) AppendDurationRef(time.Duration)      {}
func (s *stringCollector) AppendTimeLayout(time.Time, string)   {}
func (s *stringCollector) AppendIP(ip any)                      {}
func (s *stringCollector) AppendIPNet(net any)                  {}
func (s *stringCollector) AppendMAC(mac any)                    {}
func (s *stringCollector) AppendHex(any)                        {}
func (s *stringCollector) AppendFloat(any)                      {}
func (s *stringCollector) Cap() int                             { return 0 }
func (s *stringCollector) Len() int                             { return len(s.items) }
func (s *stringCollector) Truncate(int)                         {}
