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

// inspired by dubbogo/remoting/getty
package tcp

import (
	"fmt"
	"net"
	"sync"
	"time"
)

import (
	getty "github.com/apache/dubbo-getty"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const listenerReadyTimeout = 5 * time.Second

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeTCP, newTcpListenerService)
}

type (
	// ListenerService the facade of a listener
	TcpListenerService struct {
		*listener.BaseListenerService
		server          getty.Server
		gShutdownConfig *listener.ListenerGracefulShutdownConfig
	}
)

func newTcpListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	serverOpts := []getty.ServerOption{getty.WithLocalAddress(lc.Address.SocketAddress.GetAddress())}
	// todo taskPoolMode
	server := getty.NewTCPServer(serverOpts...)

	fc, err := filterchain.BuildNetworkFilterChain(lc.FilterChain)
	if err != nil {
		return nil, err
	}
	return &TcpListenerService{
		BaseListenerService: listener.NewBaseListenerService(lc, fc),
		server:              server,
		gShutdownConfig:     &listener.ListenerGracefulShutdownConfig{},
	}, nil
}

// Start start tcp server
func (ls *TcpListenerService) Start() error {
	startErr := make(chan error, 1)
	ready := make(chan struct{})
	var readyOnce sync.Once
	go func() {
		defer func() {
			if recovered := recover(); recovered != nil {
				startErr <- fmt.Errorf("start TCP listener: %v", recovered)
			}
		}()
		ls.server.RunEventLoop(func(session getty.Session) error {
			readyOnce.Do(func() { close(ready) })
			return ls.newSession(session)
		})
	}()

	address := ls.Config.Address.SocketAddress.GetAddress()
	deadline := time.NewTimer(listenerReadyTimeout)
	defer deadline.Stop()
	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case err := <-startErr:
			return err
		case <-ready:
			return nil
		case <-ticker.C:
			conn, err := net.DialTimeout("tcp", address, 50*time.Millisecond)
			if err == nil {
				_ = conn.Close()
			}
		case <-deadline.C:
			return fmt.Errorf("TCP listener %s did not become ready within %s", address, listenerReadyTimeout)
		}
	}
}

func (ls *TcpListenerService) Close() error {
	ls.server.Close()
	return ls.CloseFilterChain()
}

func (ls *TcpListenerService) ShutDown(wg any) error {
	timeout := config.GetBootstrap().GetShutdownConfig().GetTimeout()
	if timeout <= 0 {
		return nil
	}
	// stop accept request
	ls.gShutdownConfig.SetRejectRequests(true)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) && ls.gShutdownConfig.GetActiveCount() > 0 {
		// sleep 100 ms and check it again
		time.Sleep(100 * time.Millisecond)
		logger.Infof("waiting for active invocation count = %d", ls.gShutdownConfig.GetActiveCount())
	}
	wg.(*sync.WaitGroup).Done()
	ls.server.Close()
	return nil
}

func (ls *TcpListenerService) newSession(session getty.Session) (err error) {

	tcpConn, ok := session.Conn().(*net.TCPConn)
	if !ok {
		panic(fmt.Sprintf("newSession: %s, session.conn{%#v} is not tcp connection", session.Stat(), session.Conn()))
	}
	// default value from dubbo getty config
	// todo: make parameter for tcp listener config
	if err = tcpConn.SetNoDelay(true); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlive(true); err != nil {
		return err
	}
	if err = tcpConn.SetKeepAlivePeriod(10 * time.Second); err != nil {
		return err
	}
	if err = tcpConn.SetReadBuffer(262144); err != nil {
		return err
	}
	if err = tcpConn.SetWriteBuffer(524288); err != nil {
		return err
	}

	session.SetMaxMsgLen(128 * 1024) // max message package length is 128k
	session.SetReadTimeout(time.Second)
	session.SetWriteTimeout(time.Second)
	session.SetWaitTime(time.Second)

	session.SetPkgHandler(NewPackageHandler(ls))
	session.SetEventListener(NewServerPackageHandler(ls))
	return nil
}
