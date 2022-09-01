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

package sentinel

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

// loggerWrapper
type loggerWrapper struct {
}

// GetWrappedLogger
func GetWrappedLogger() loggerWrapper {
	return loggerWrapper{}
}

func (l loggerWrapper) Debug(msg string, keysAndValues ...interface{}) {
	logger.Debugf(msg, keysAndValues...)
}

func (l loggerWrapper) Info(msg string, keysAndValues ...interface{}) {
	logger.Infof(msg, keysAndValues...)
}

func (l loggerWrapper) Warn(msg string, keysAndValues ...interface{}) {
	logger.Warnf(msg, keysAndValues...)
}

func (l loggerWrapper) Error(err error, msg string, keysAndValues ...interface{}) {
	logger.Warnf(msg, err, keysAndValues)
}

// InfoEnabled todo logger should implements this method
func (l loggerWrapper) InfoEnabled() bool {
	return true
}

func (l loggerWrapper) ErrorEnabled() bool {
	return true
}

func (l loggerWrapper) DebugEnabled() bool {
	return true
}

func (l loggerWrapper) WarnEnabled() bool {
	return true
}
