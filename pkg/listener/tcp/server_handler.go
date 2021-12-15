// inspired by dubbogo/remoting/getty
package tcp

import (
	getty "github.com/apache/dubbo-getty"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	perrors "github.com/pkg/errors"
	"sync"
	"sync/atomic"
	"time"
)

var (
	errTooManySessions      = perrors.New("too many sessions")
	errHeartbeatReadTimeout = perrors.New("heartbeat read timeout")
)

type rpcSession struct {
	session getty.Session
	reqNum  int32
}

func (s *rpcSession) AddReqNum(num int32) {
	atomic.AddInt32(&s.reqNum, num)
}

func (s *rpcSession) GetReqNum() int32 {
	return atomic.LoadInt32(&s.reqNum)
}

type ServerHandler struct {
	ls             *TcpListenerService
	sessionMap     map[getty.Session]*rpcSession
	rwlock         sync.RWMutex
	maxSessionNum  int
	sessionTimeout time.Duration
	timeoutTimes   int
}

func NewServerPackageHandler(ls *TcpListenerService) *ServerHandler {
	return &ServerHandler{
		// todo listener param
		maxSessionNum:  1000,
		sessionTimeout: 60 * time.Second,
		ls:             ls,
		sessionMap:     make(map[getty.Session]*rpcSession),
	}
}

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

func (h *ServerHandler) OnError(session getty.Session, err error) {
	logger.Infof("session{%s} got error{%v}, will be closed.", session.Stat(), err)
	h.rwlock.Lock()
	delete(h.sessionMap, session)
	h.rwlock.Unlock()
}

func (h *ServerHandler) OnClose(session getty.Session) {
	logger.Infof("session{%s} is closing......", session.Stat())
	h.rwlock.Lock()
	delete(h.sessionMap, session)
	h.rwlock.Unlock()
}

func (h *ServerHandler) OnMessage(session getty.Session, pkg interface{}) {
	h.rwlock.Lock()
	if _, ok := h.sessionMap[session]; ok {
		h.sessionMap[session].reqNum++
	}
	h.rwlock.Unlock()

	// don't following the rule of getty data handle chain, do nothing.
}

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

	//heartbeatCallBack := func(err error) {
	//	if err != nil {
	//		logger.Warnf("failed to send heartbeat, error{%v}", err)
	//		if h.timeoutTimes >= 3 {
	//			h.rwlock.Lock()
	//			delete(h.sessionMap, session)
	//			h.rwlock.Unlock()
	//			session.Close()
	//			return
	//		}
	//		h.timeoutTimes++
	//		return
	//	}
	//	h.timeoutTimes = 0
	//}

	//if err := heartbeat(session, h.server.conf.heartbeatTimeout, heartbeatCallBack); err != nil {
	//	logger.Warnf("failed to send heartbeat, error{%v}", err)
	//}
}
