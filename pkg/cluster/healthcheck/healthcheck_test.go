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
	"log"
	"net"
	"strconv"
	"testing"
	"time"
)

func TestTcpConn(t *testing.T) {

	go func() {
		listener, err := net.Listen("tcp", ":12347")
		if err != nil {
			log.Fatalf("listen failed: %v", err)
		}

		for {
			client, err := listener.Accept()
			if err != nil {
				log.Printf("accept new client failed: %v", err)
				continue
			}

			client.Close()
		}
	}()

	time.Sleep(100 * time.Millisecond) // Give the server some time to start

	type TestCase struct {
		name     string
		addr     string
		port     int
		timeout  time.Duration
		expected bool
	}

	testCases := []TestCase{
		{
			name:     "12347",
			addr:     "127.0.0.1",
			port:     12347,
			timeout:  100 * time.Millisecond,
			expected: true,
		},
		{
			name:     "12348",
			addr:     "127.0.0.1",
			port:     12348,
			timeout:  100 * time.Millisecond,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := CheckTcpConn(tc.addr, strconv.Itoa(tc.port), tc.timeout)
			if actual != tc.expected {
				t.Fail()
			}
		})
	}
}
