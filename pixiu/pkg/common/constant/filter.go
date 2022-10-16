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

package constant

import (
	"time"
)

var (
	Default403Body = []byte("403 for bidden")
	Default404Body = []byte("404 page not found")
	Default405Body = []byte("405 method not allowed")
	Default406Body = []byte("406 api not up")
	Default503Body = []byte("503 service unavailable")
)

const (
	// nolint
	// FileDateFormat
	FileDateFormat = "2006-01-02"
	// nolint
	// MessageDateLayout
	MessageDateLayout = "2006-01-02 15:04:05"
	// nolint
	// LogFileMode
	LogFileMode = 0o600
	// nolint
	// buffer
	LogDataBuffer = 5000
	// console
	Console = "console"

	DefaultReqTimeout = 10 * time.Second
)
