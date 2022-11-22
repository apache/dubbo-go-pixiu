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
	"sync"
	"sync/atomic"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	getty "github.com/apache/dubbo-getty"
	hessian "github.com/apache/dubbo-go-hessian2"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

const (
	// WritePkg_Timeout the timeout of write pkg
	WritePkg_Timeout = 5 * time.Second
)

var (
	errTooManySessions = perrors.New("too many sessions")
)

type rpcSession struct {
	session getty.Session
	reqNum  int32
}

// AddReqNum add request total num
func (s *rpcSession) AddReqNum(num int32) {
	atomic.AddInt32(&s.reqNum, num)
}

// GetReqNum get request total num
func (s *rpcSession) GetReqNum() int32 {
	return atomic.LoadInt32(&s.reqNum)
}

// ServerHandler package handler
type ServerHandler struct {
	ls             *TcpListenerService
	sessionMap     map[getty.Session]*rpcSession
	rwlock         sync.RWMutex
	maxSessionNum  int
	sessionTimeout time.Duration
}

// NewServerPackageHandler create serverHandler
func NewServerPackageHandler(ls *TcpListenerService) *ServerHandler {
	return &ServerHandler{
		// todo listener param
		maxSessionNum:  1000,
		sessionTimeout: 60 * time.Second,
		ls:             ls,
		sessionMap:     make(map[getty.Session]*rpcSession),
	}
}

// OnOpen called when session opens
func (h *ServerHandler) OnOpen(session getty.Session) error {
	var err error
	h.rwlock.RLock()
	if h.maxSessionNum <= len(h.sessionMap) {
		err = errTooManySessions
	}
	h.rwlock.RUnlock()
	if err != nil {
		return perrors.WithStack(err)
	}

	logger.Infof("got session:%s", session.Stat())
	h.rwlock.Lock()
	h.sessionMap[session] = &rpcSession{session: session}
	h.rwlock.Unlock()
	return nil
}

// OnError called when err
func (h *ServerHandler) OnError(session getty.Session, err error) {
	logger.Infof("session{%s} got error{%v}, will be closed.", session.Stat(), err)
	h.rwlock.Lock()
	delete(h.sessionMap, session)
	h.rwlock.Unlock()
}

// OnError called when session close
func (h *ServerHandler) OnClose(session getty.Session) {
	logger.Infof("session{%s} is closing......", session.Stat())
	h.rwlock.Lock()
	delete(h.sessionMap, session)
	h.rwlock.Unlock()
}

// OnMessage called when session receive new pkg
func (h *ServerHandler) OnMessage(session getty.Session, pkg interface{}) {
	h.ls.gShutdownConfig.AddActiveCount(1)
	defer h.ls.gShutdownConfig.AddActiveCount(-1)

	h.rwlock.Lock()
	if _, ok := h.sessionMap[session]; ok {
		h.sessionMap[session].AddReqNum(1)
	}
	h.rwlock.Unlock()

	decodeResult, drOK := pkg.(*remoting.DecodeResult)
	if !drOK {
		logger.Errorf("illegal package{%#v}", pkg)
		return
	}
	if !decodeResult.IsRequest {
		res := decodeResult.Result.(*remoting.Response)
		if res.Event {
			logger.Debugf("get rpc heartbeat response{%#v}", res)
			if res.Error != nil {
				logger.Errorf("rpc heartbeat response{error: %#v}", res.Error)
			}
			res.Handle()
			return
		}
		logger.Errorf("illegal package but not heartbeat. {%#v}", pkg)
		return
	}
	req := decodeResult.Result.(*remoting.Request)

	resp := remoting.NewResponse(req.ID, req.Version)
	resp.Status = hessian.Response_OK
	resp.Event = req.Event
	resp.SerialID = req.SerialID
	resp.Version = "2.0.2"

	// heartbeat
	if req.Event {
		logger.Debugf("get rpc heartbeat request{%#v}", resp)
		reply(session, resp)
		return
	}

	defer func() {
		if e := recover(); e != nil {
			resp.Status = hessian.Response_SERVER_ERROR
			if err, ok := e.(error); ok {
				logger.Errorf("OnMessage panic: %+v", perrors.WithStack(err))
				resp.Error = perrors.WithStack(err)
			} else if err, ok := e.(string); ok {
				logger.Errorf("OnMessage panic: %+v", perrors.New(err))
				resp.Error = perrors.New(err)
			} else {
				logger.Errorf("OnMessage panic: %+v, this is impossible.", e)
				resp.Error = fmt.Errorf("OnMessage panic unknow exception. %+v", e)
			}

			if !req.TwoWay {
				return
			}
			reply(session, resp)
		}
	}()

	if h.ls.gShutdownConfig.RejectRequest {
		err := perrors.Errorf("Pixiu is preparing to close, reject all new requests")
		resp.Result = protocol.RPCResult{
			Err: err,
		}
		reply(session, resp)
		return
	}

	invoc, ok := req.Data.(*invocation.RPCInvocation)
	if !ok {
		panic("create invocation occur some exception for the type is not suitable one.")
	}
	attachments := invoc.Attachments()
	attachments["local-addr"] = session.LocalAddr()
	attachments["remote-addr"] = session.RemoteAddr()

	result, err := h.ls.FilterChain.OnData(invoc)
	if err != nil {
		resp.Error = fmt.Errorf("OnData panic unknow exception. %+v", err)
		if !req.TwoWay {
			return
		}
		reply(session, resp)
	}

	if !req.TwoWay {
		return
	}
	if result != nil {
		ptr, _ := result.(*protocol.RPCResult)
		resp.Result = *ptr
	} else {
		resp.Result = nil
	}
	reply(session, resp)
}

func reply(session getty.Session, resp *remoting.Response) {
	if totalLen, sendLen, err := session.WritePkg(resp, WritePkg_Timeout); err != nil {
		if sendLen != 0 && totalLen != sendLen {
			logger.Warnf("start to close the session at replying because %d of %d bytes data is sent success. err:%+v", sendLen, totalLen, err)
			go session.Close()
		}
		logger.Errorf("WritePkg error: %#v, %#v", perrors.WithStack(err), resp)
	}
}

// OnCron cron
func (h *ServerHandler) OnCron(session getty.Session) {

	var (
		flag   bool
		active time.Time
	)

	h.rwlock.RLock()
	if _, ok := h.sessionMap[session]; ok {
		active = session.GetActive()
		if h.sessionTimeout.Nanoseconds() < time.Since(active).Nanoseconds() {
			flag = true
			logger.Warnf("session{%s} timeout{%s}, reqNum{%d}",
				session.Stat(), time.Since(active).String(), h.sessionMap[session].reqNum)
		}
	}
	h.rwlock.RUnlock()

	if flag {
		h.rwlock.Lock()
		delete(h.sessionMap, session)
		h.rwlock.Unlock()
		session.Close()
	}
}
