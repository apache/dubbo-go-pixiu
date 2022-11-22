// Copyright Istio Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package forwarder

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
)

import (
	"github.com/gorilla/websocket"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/test/echo/common"
)

var _ protocol = &websocketProtocol{}

type websocketProtocol struct {
	*Config
}

func newWebsocketProtocol(r *Config) (protocol, error) {
	return &websocketProtocol{
		Config: r,
	}, nil
}

func (c *websocketProtocol) makeRequest(ctx context.Context, req *request) (string, error) {
	wsReq := make(http.Header)

	var outBuffer bytes.Buffer
	outBuffer.WriteString(fmt.Sprintf("[%d] Url=%s\n", req.RequestID, req.URL))
	writeHeaders(req.RequestID, req.Header, outBuffer, wsReq.Add)

	// Set the special header to trigger the upgrade to WebSocket.
	common.SetWebSocketHeader(wsReq)

	if req.Message != "" {
		outBuffer.WriteString(fmt.Sprintf("[%d] Echo=%s\n", req.RequestID, req.Message))
	}

	dialContext := func(network, addr string) (net.Conn, error) {
		return newDialer().Dial(network, addr)
	}
	if len(c.UDS) > 0 {
		dialContext = func(network, addr string) (net.Conn, error) {
			return newDialer().Dial("unix", c.UDS)
		}
	}

	dialer := &websocket.Dialer{
		TLSClientConfig:  c.tlsConfig,
		NetDial:          dialContext,
		HandshakeTimeout: c.timeout,
	}

	conn, _, err := dialer.Dial(req.URL, wsReq)
	if err != nil {
		// timeout or bad handshake
		return outBuffer.String(), err
	}
	defer func() {
		_ = conn.Close()
	}()

	// Apply per-request timeout to calculate deadline for reads/writes.
	ctx, cancel := context.WithTimeout(ctx, req.Timeout)
	defer cancel()

	// Apply the deadline to the connection.
	deadline, _ := ctx.Deadline()
	if err := conn.SetWriteDeadline(deadline); err != nil {
		return outBuffer.String(), err
	}
	if err := conn.SetReadDeadline(deadline); err != nil {
		return outBuffer.String(), err
	}

	err = conn.WriteMessage(websocket.TextMessage, []byte(req.Message))
	if err != nil {
		return outBuffer.String(), err
	}

	_, resp, err := conn.ReadMessage()
	if err != nil {
		return outBuffer.String(), err
	}

	for _, line := range strings.Split(string(resp), "\n") {
		if line != "" {
			outBuffer.WriteString(fmt.Sprintf("[%d body] %s\n", req.RequestID, line))
		}
	}

	return outBuffer.String(), nil
}

func (c *websocketProtocol) Close() error {
	return nil
}
