package tcp

import (
	"fmt"
	getty "github.com/apache/dubbo-getty"
	"github.com/apache/dubbo-getty/demo/hello"
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"net"
	"time"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeTCP, newTcpListenerService)
}

type (
	// ListenerService the facade of a listener
	TcpListenerService struct {
		listener.BaseListenerService
		server getty.Server
	}
)

func newTcpListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	serverOpts := []getty.ServerOption{getty.WithLocalAddress(lc.Address.SocketAddress.GetAddress())}
	// todo taskPoolMode
	server := getty.NewTCPServer(serverOpts...)

	fc := filterchain.CreateNetworkFilterChain(lc.FilterChain, bs)
	return &TcpListenerService{
		BaseListenerService: listener.BaseListenerService{
			Config:      lc,
			FilterChain: fc,
		},
		server: server,
	}, nil
}

func (ls *TcpListenerService) Start() error {
	go ls.server.RunEventLoop(ls.newSession)
	return nil
}

func (ls *TcpListenerService) newSession(session getty.Session) (err error) {
	// session.SetCompressType(getty.CompressZip)

	tcpConn, ok := session.Conn().(*net.TCPConn)
	if !ok {
		panic(fmt.Sprintf("newSession: %s, session.conn{%#v} is not tcp connection", session.Stat(), session.Conn()))
	}

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

	session.SetName("hello")
	session.SetMaxMsgLen(128 * 1024) // max message package length is 128k
	session.SetReadTimeout(time.Second)
	session.SetWriteTimeout(5 * time.Second)
	session.SetCronPeriod(int(hello.CronPeriod / 1e6))
	session.SetWaitTime(time.Second)

	session.SetPkgHandler(NewPackageHandler(ls))
	session.SetEventListener(NewServerPackageHandler(ls))
	return nil
}
