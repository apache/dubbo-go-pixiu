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
		{"unknown", false, "info"}, // parseLevel 默认回退到 info（并返回 false）
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
	log = &pixiuLogger{SugaredLogger: log.SugaredLogger.With("k", "v"), config: log.config}
	log.Infow("with fields", "a", 1)

	_ = log.Sync()
}
