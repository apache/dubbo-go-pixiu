package http2

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/gorilla/mux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net"
	"net/http"
	"strconv"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeHTTP2, newHttp2ListenerService)
}

type (
	// ListenerService the facade of a listener
	Http2ListenerService struct {
		Config  *model.Listener
		listener         net.Listener
		server  *http.Server
		fc *filterchain.FilterChain
	}
)

type handleWrapper struct {
	router *mux.Router
}

type h2cWrapper struct {
	w *handleWrapper
	h http.Handler
}

func (h *h2cWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.h.ServeHTTP(w, r)
}
// 转发 grpc 的请求
func (h *handleWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)

}

func newHttp2ListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	fc := filterchain.CreateFilterChain(lc.FilterChain, bs)
	return &Http2ListenerService {
		Config:  lc,
		fc: fc,
	}, nil
}

func (g Http2ListenerService) GetNetworkFilter() filter.NetworkFilter {
	panic("implement me")
}

func (ls Http2ListenerService) Start() error {

	sa := ls.Config.Address.SocketAddress
	addr := resolveAddress(sa.Address + ":" + strconv.Itoa(sa.Port))

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	ls.listener = l

	var h http.Handler
	router := mux.NewRouter()
	h = &handleWrapper{router}
	h2s := &http2.Server{}
	h = &h2cWrapper{
		w: h.(*handleWrapper),
		h: h2c.NewHandler(h, h2s),
	}

	ls.server = &http.Server{
		Addr:         addr,
		Handler:      h,
	}

	go ls.server.Serve(ls.listener)
	return nil
}

func resolveAddress(addr string) string {
	if addr == "" {
		logger.Debug("Addr is undefined. Using port :8080 by default")
		return ":8080"
	}

	return addr
}