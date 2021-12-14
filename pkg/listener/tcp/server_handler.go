package tcp

import getty "github.com/apache/dubbo-getty"

// inspired by dubbogo/remoting/getty
type ServerHandler struct {
	ls *TcpListenerService
}

func NewServerPackageHandler(ls *TcpListenerService) *ServerHandler {
	return &ServerHandler{ls}
}

func (h *ServerHandler) OnOpen(session getty.Session) error {

}

func (h *ServerHandler) OnError(session getty.Session, err error) {

}

func (h *ServerHandler) OnClose(session getty.Session) {

}

func (h *ServerHandler) OnMessage(session getty.Session, pkg interface{}) {

}

func (h *ServerHandler) OnCron(session getty.Session) {

}
