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

package mcpserver

import (
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func markToolsListChangedAndFlush(sseHandler *transport.SSEHandler, session *transport.MCPSession) (bool, bool, error) {
	if session == nil {
		return false, false, nil
	}
	if session.MarkToolsListChangedPending() == 0 {
		return false, false, nil
	}
	if err := flushToolsListChangedNotification(sseHandler, session); err != nil {
		return true, false, err
	}
	return true, true, nil
}

func flushToolsListChangedNotification(sseHandler *transport.SSEHandler, session *transport.MCPSession) error {
	if session == nil {
		return nil
	}
	if sseHandler == nil {
		return fmt.Errorf("SSE handler not configured")
	}
	for {
		version, pending := session.PendingToolsListChangedVersion()
		if !pending || !session.HasPipeWriter() {
			return nil
		}
		notification := map[string]any{
			"jsonrpc": "2.0",
			"method":  toolsListChangedMethod,
		}
		if err := sseHandler.SendSSEMessage(session, notification); err != nil {
			return err
		}
		session.MarkToolsListChangedNotified(version)
		logger.Debugf("[dubbo-go-pixiu] mcp server sent tools/list_changed")
	}
}
