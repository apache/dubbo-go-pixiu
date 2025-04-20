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

package healthcheck

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"net"
	"strconv"
	"strings"
	"time"
)

func TcpConn(tarAddr string, port int, timeout time.Duration) bool {
	if _, _, err := net.SplitHostPort(tarAddr); err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			tarAddr = tarAddr + ":" + strconv.Itoa(port)
		} else {
			logger.Infof("[health check] invalid address format: %s", tarAddr)
			return false
		}
	}
	conn, err := net.DialTimeout("tcp", tarAddr, timeout)
	if err != nil {
		logger.Infof("[health check] http checker for host %s error: %v", tarAddr, err)
		return false
	}
	conn.Close()
	return true
}
