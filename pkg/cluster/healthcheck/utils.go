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
	"net"
	"strings"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func CheckTcpConn(address string, port string, timeout time.Duration) bool {
	normalizedAddress, err := normalizeAddress(address, port)
	if err != nil {
		logger.Infof("[health check] address format for address \"%s\" failed, %s", address, err.Error())
		return false
	}

	conn, err := net.DialTimeout("tcp", normalizedAddress, timeout)
	if err != nil {
		logger.Infof("[health check] health check for address \"%s\" failed, %s", normalizedAddress, err.Error())
		return false
	}
	defer conn.Close()
	return true
}

// normalizeAddress normalizes the address by ensuring it has the correct port.
// If the port field is empty, it will check if the address already has a port.
//   - If the address has a port, it will return the address as is.
//   - If the address does not have a port, it will return err.
//
// If the port field is not empty, it will check if the address's port matches the provided port.
//   - If it matches, it will return the address as is.
//   - If it does not match, it will return the address with the new port.
func normalizeAddress(address string, port string) (string, error) {
	if port == "" {
		_, _, err := net.SplitHostPort(address)
		if err != nil {
			return "", err
		}
		return address, nil
	}

	host, existingPort, err := net.SplitHostPort(address)
	if err != nil {
		if strings.Contains(err.Error(), "missing port in address") {
			return net.JoinHostPort(strings.Trim(address, "[]"), port), nil
		}
		return "", err
	}

	if existingPort != port {
		return net.JoinHostPort(host, port), nil
	}

	return address, nil
}
