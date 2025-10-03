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

func TestDisableStacktraceTrueNoStackEvenOnError(t *testing.T) {
	cfg, out := newDevConfigToFile(t)
	cfg.DisableStacktrace = true // YAML equivalent: disableStacktrace: true
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

func TestSetLoggerLevelDoesNotRebuildAndTakesEffect(t *testing.T) {
	cfg, out := newDevConfigToFile(t)
	InitLogger(cfg)

	before := GetLogger().SugaredLogger

	// dynamic set to error
	ok := SetLoggerLevel(zapcore.ErrorLevel)
	if !ok {
		t.Fatalf("SetLoggerLevel returned false")
	}

	after := GetLogger().SugaredLogger
	if before != after {
		t.Fatalf("SetLoggerLevel should NOT rebuild logger; pointer changed: %p -> %p", before, after)
	}

	// write new: info should not appear, error should appear
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

func TestHotReloadRebuildsAndSwitchesSink(t *testing.T) {
	// cfg1 -> out1
	cfg1, out1 := newDevConfigToFile(t)
	InitLogger(cfg1)
	l1 := GetLogger().SugaredLogger

	const (
		H1 = "hello-1"
		H2 = "hello-2"
	)

	GetLogger().Info(H1)
	_ = GetLogger().Sync()

	// cfg2 -> out2（new sink）
	cfg2, out2 := newDevConfigToFile(t)
	// to split it, set level to info
	HotReload(cfg2)
	l2 := GetLogger().SugaredLogger

	if l1 == l2 {
		t.Fatalf("HotReload should rebuild logger; got same pointer %p", l1)
	}

	GetLogger().Info(H2)
	_ = GetLogger().Sync()

	got1 := readAll(t, out1)
	got2 := readAll(t, out2)

	if !strings.Contains(got1, H1) {
		t.Fatalf("out1 should contain hello-1 but not found:\n%s", got1)
	}
	if strings.Contains(got1, H2) {
		t.Fatalf("out1 should NOT contain hello-2 after reload:\n%s", got1)
	}
	if !strings.Contains(got2, H2) {
		t.Fatalf("out2 should contain hello-2 but not found:\n%s", got2)
	}
}

func TestPaddedCallerEncoderFixedWidthAtLeast30(t *testing.T) {
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

func (s *stringCollector) AppendBool(bool)                      { /*mock*/ }
func (s *stringCollector) AppendByteString([]byte)              { /*mock*/ }
func (s *stringCollector) AppendComplex128(complex128)          { /*mock*/ }
func (s *stringCollector) AppendComplex64(complex64)            { /*mock*/ }
func (s *stringCollector) AppendDuration(time.Duration)         { /*mock*/ }
func (s *stringCollector) AppendFloat64(float64)                { /*mock*/ }
func (s *stringCollector) AppendFloat32(float32)                { /*mock*/ }
func (s *stringCollector) AppendInt(int)                        { /*mock*/ }
func (s *stringCollector) AppendInt64(int64)                    { /*mock*/ }
func (s *stringCollector) AppendInt32(int32)                    { /*mock*/ }
func (s *stringCollector) AppendInt16(int16)                    { /*mock*/ }
func (s *stringCollector) AppendInt8(int8)                      { /*mock*/ }
func (s *stringCollector) AppendTime(time.Time)                 { /*mock*/ }
func (s *stringCollector) AppendUint(uint)                      { /*mock*/ }
func (s *stringCollector) AppendUint64(uint64)                  { /*mock*/ }
func (s *stringCollector) AppendUint32(uint32)                  { /*mock*/ }
func (s *stringCollector) AppendUint16(uint16)                  { /*mock*/ }
func (s *stringCollector) AppendUint8(uint8)                    { /*mock*/ }
func (s *stringCollector) AppendUintptr(uintptr)                { /*mock*/ }
func (s *stringCollector) AppendReflected(any)                  { /*mock*/ }
func (s *stringCollector) AppendArray(zapcore.ArrayMarshaler)   { /*mock*/ }
func (s *stringCollector) AppendObject(zapcore.ObjectMarshaler) { /*mock*/ }
func (s *stringCollector) AppendBinary([]byte)                  { /*mock*/ }
func (s *stringCollector) AppendComplex(complex128)             { /*mock*/ }
func (s *stringCollector) AppendDurationRef(time.Duration)      { /*mock*/ }
func (s *stringCollector) AppendTimeLayout(time.Time, string)   { /*mock*/ }
func (s *stringCollector) AppendIP(ip any)                      { /*mock*/ }
func (s *stringCollector) AppendIPNet(net any)                  { /*mock*/ }
func (s *stringCollector) AppendMAC(mac any)                    { /*mock*/ }
func (s *stringCollector) AppendHex(any)                        { /*mock*/ }
func (s *stringCollector) AppendFloat(any)                      { /*mock*/ }
func (s *stringCollector) Cap() int                             { return 0 }
func (s *stringCollector) Len() int                             { return len(s.items) }
func (s *stringCollector) Truncate(int)                         { /*mock*/ }
