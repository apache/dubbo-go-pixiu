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

package uds

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestUdsListener(t *testing.T) {
	l, err := NewListener("./etc/istio/proxy/test")
	if err != nil {
		t.Fatalf("unexpected error %v", err)
	}
	defer l.Close()
	conn, err := connect("./etc/istio/proxy/test")
	if err != nil {
		t.Fatalf("failed to connect %v", err)
	}
	conn.Close()
}

func connect(socket string) (*grpc.ClientConn, error) {
	var opts []grpc.DialOption

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", socket)
	}))

	conn, err := grpc.Dial(socket, opts...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
