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

package listener_test

import (
	"net"
	"testing"
)

import (
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	_ "github.com/apache/dubbo-go-pixiu/pkg/listener/grpc"
	_ "github.com/apache/dubbo-go-pixiu/pkg/listener/http"
	_ "github.com/apache/dubbo-go-pixiu/pkg/listener/http2"
	_ "github.com/apache/dubbo-go-pixiu/pkg/listener/tcp"
	_ "github.com/apache/dubbo-go-pixiu/pkg/listener/triple"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestListenerStartReturnsOccupiedPortError(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = occupied.Close() })
	port := occupied.Addr().(*net.TCPAddr).Port

	protocols := []struct {
		name     string
		protocol model.ProtocolType
	}{
		{name: "HTTP", protocol: model.ProtocolTypeHTTP},
		{name: "HTTP2", protocol: model.ProtocolTypeHTTP2},
		{name: "TCP", protocol: model.ProtocolTypeTCP},
		{name: "gRPC", protocol: model.ProtocolTypeGRPC},
		{name: "Triple", protocol: model.ProtocolTypeTriple},
	}
	for _, protocol := range protocols {
		t.Run(protocol.name, func(t *testing.T) {
			config := &model.Listener{
				Name:        "occupied-" + protocol.name,
				Protocol:    protocol.protocol,
				ProtocolStr: protocol.name,
				Address: model.Address{SocketAddress: model.SocketAddress{
					Address: "127.0.0.1",
					Port:    port,
				}},
			}
			service, err := listener.CreateListenerService(config, &model.Bootstrap{})
			require.NoError(t, err)
			require.Error(t, service.Start(), "Start must not return success before binding the port")
			_ = service.Close()
		})
	}
}
