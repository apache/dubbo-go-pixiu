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

// TripleLogger does not add user-defined log fields, so there will be no pointer
// reference issues after updates, ensuring that all references are under control.
type TripleLogger struct {
}

func GetTripleLogger() *TripleLogger {
	return &TripleLogger{}
}

func (l *TripleLogger) Info(args ...any) {
	control.info(args...)
}

func (l *TripleLogger) Warn(args ...any) {
	control.warn(args...)
}

func (l *TripleLogger) Error(args ...any) {
	control.error(args...)
}

func (l *TripleLogger) Debug(args ...any) {
	control.debug(args...)
}

func (l *TripleLogger) Infof(fmt string, args ...any) {
	control.infof(fmt, args...)
}

func (l *TripleLogger) Warnf(fmt string, args ...any) {
	control.warnf(fmt, args...)
}

func (l *TripleLogger) Errorf(fmt string, args ...any) {
	control.errorf(fmt, args...)
}

func (l *TripleLogger) Debugf(fmt string, args ...any) {
	control.debugf(fmt, args...)
}
