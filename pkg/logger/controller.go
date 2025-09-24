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
	"strings"
	"sync"
)

import (
	"go.uber.org/zap/zapcore"
)

// logController governs the logging output or configuration changes throughout the entire project.
type logController struct {
	mu     sync.RWMutex
	logger *pixiuLogger
}

// setLoggerLevel changes the level at runtime without rebuilding the logger.
func (c *logController) setLoggerLevel(level string) bool {
	lvl, ok := c.parseLevel(level)
	if !ok {
		return false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.logger == nil || c.logger.config == nil {
		return false
	}
	c.logger.config.Level.SetLevel(lvl)
	return true
}

// updateLogger swaps the underlying logger atomically.
func (c *logController) updateLogger(l *pixiuLogger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = l
}

func (c *logController) debug(args ...any) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.logger.Debug(args...)
}

func (c *logController) info(args ...any) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.logger.Info(args...)
}

func (c *logController) warn(args ...any) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.logger.Warn(args...)
}

func (c *logController) error(args ...any) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.logger.Error(args...)
}

func (c *logController) debugf(fmt string, args ...any) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.logger.Debugf(fmt, args...)
}

func (c *logController) infof(fmt string, args ...any) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.logger.Infof(fmt, args...)
}

func (c *logController) warnf(fmt string, args ...any) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.logger.Warnf(fmt, args...)
}

func (c *logController) errorf(fmt string, args ...any) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.logger.Errorf(fmt, args...)
}

// parseLevel parses textual level to zapcore.Level.
func (c *logController) parseLevel(level string) (zapcore.Level, bool) {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel, true
	case "info":
		return zapcore.InfoLevel, true
	case "warn", "warning":
		return zapcore.WarnLevel, true
	case "error":
		return zapcore.ErrorLevel, true
	case "dpanic":
		return zapcore.DPanicLevel, true
	case "panic":
		return zapcore.PanicLevel, true
	case "fatal":
		return zapcore.FatalLevel, true
	default:
		return zapcore.InfoLevel, false
	}
}
