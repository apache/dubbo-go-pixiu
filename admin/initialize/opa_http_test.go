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

package initialize

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func useDirectHTTPTransport(t *testing.T) {
	t.Helper()
	prev := http.DefaultTransport
	transport := directHTTPTransport()
	http.DefaultTransport = transport
	t.Cleanup(func() {
		http.DefaultTransport = prev
		transport.CloseIdleConnections()
	})
}

func directHTTPTransport() *http.Transport {
	if base, ok := http.DefaultTransport.(*http.Transport); ok && base != nil {
		transport := base.Clone()
		transport.Proxy = nil
		return transport
	}
	return &http.Transport{Proxy: nil}
}

func startLoopbackHTTPServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	listener, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen on 127.0.0.1:0: %v", err)
	}
	srv := httptest.NewUnstartedServer(handler)
	srv.Listener = listener
	srv.Start()
	t.Cleanup(srv.Close)
	return srv
}
