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

package hotreload

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// LoggerReloader implements the HotReloader interface for reloading logger configurations.
type LoggerReloader struct{}

// CheckUpdate compares the old and new logger configurations to determine if a reload is needed.
func (r *LoggerReloader) CheckUpdate(oldConfig, newConfig *model.Bootstrap) bool {
	oc := oldConfig.Log
	nc := newConfig.Log

	if oc == nil && nc != nil {
		return true
	}

	if oc != nil && nc == nil {
		return false
	}

	// Check if any logger configuration fields have changed.
	if oc.Level != nc.Level ||
		oc.Development != nc.Development ||
		oc.DisableCaller != nc.DisableCaller ||
		oc.DisableStacktrace != nc.DisableStacktrace ||
		oc.Encoding != nc.Encoding {
		return true
	}

	// Check sampling configuration.
	if !r.checkSampling(oc.Sampling, nc.Sampling) {
		return true
	}

	// Check encoder configuration.
	if !r.checkEncoderConfig(oc.EncoderConfig, nc.EncoderConfig) {
		return true
	}

	// Check output paths.
	if !equal(oc.OutputPaths, nc.OutputPaths) {
		return true
	}

	return false
}

// HotReload applies the new logger configuration.
func (r *LoggerReloader) HotReload(oldConfig, newConfig *model.Bootstrap) error {
	if err := logger.HotReload(newConfig.Log.Build()); err != nil {
		logger.Errorf("Failed to reload logger configuration: %v", err)
		return err
	}
	return nil
}

// checkSampling compares the old and new sampling configurations.
func (r *LoggerReloader) checkSampling(oldSampling, newSampling model.SamplingConfig) bool {
	return oldSampling.Initial == newSampling.Initial && oldSampling.Thereafter == newSampling.Thereafter
}

// checkEncoderConfig compares the old and new encoder configurations.
func (r *LoggerReloader) checkEncoderConfig(oldEncoderConfig, newEncoderConfig model.EncoderConfig) bool {
	return oldEncoderConfig.MessageKey == newEncoderConfig.MessageKey &&
		oldEncoderConfig.LevelKey == newEncoderConfig.LevelKey &&
		oldEncoderConfig.TimeKey == newEncoderConfig.TimeKey &&
		oldEncoderConfig.NameKey == newEncoderConfig.NameKey &&
		oldEncoderConfig.CallerKey == newEncoderConfig.CallerKey &&
		oldEncoderConfig.FunctionKey == newEncoderConfig.FunctionKey &&
		oldEncoderConfig.StacktraceKey == newEncoderConfig.StacktraceKey &&
		oldEncoderConfig.SkipLineEnding == newEncoderConfig.SkipLineEnding &&
		oldEncoderConfig.LineEnding == newEncoderConfig.LineEnding &&
		oldEncoderConfig.EncodeLevel == newEncoderConfig.EncodeLevel &&
		oldEncoderConfig.EncodeTime == newEncoderConfig.EncodeTime &&
		oldEncoderConfig.EncodeDuration == newEncoderConfig.EncodeDuration &&
		oldEncoderConfig.EncodeCaller == newEncoderConfig.EncodeCaller &&
		oldEncoderConfig.EncodeName == newEncoderConfig.EncodeName &&
		oldEncoderConfig.ConsoleSeparator == newEncoderConfig.ConsoleSeparator
}
