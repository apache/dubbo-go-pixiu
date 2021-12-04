package http

import (
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	h "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/pkg/errors"
	"golang.org/x/crypto/acme/autocert"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeHTTP, newHttpListenerService)
}

type (
	// ListenerService the facade of a listener
	HttpListenerService struct {
		listener.BaseListenerService
		srv *http.Server
	}

	// DefaultHttpListener
	DefaultHttpWorker struct {
		pool sync.Pool
		ls   *HttpListenerService
	}
)

func newHttpListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	hcm := createHttpManager(lc, bs)

	return &HttpListenerService{
		listener.BaseListenerService{Config: lc, NetworkFilter: *hcm},
		nil,
	}, nil
}

func createHttpManager(lc *model.Listener, bs *model.Bootstrap) *filter.NetworkFilter {
	p, err := filter.GetNetworkFilterPlugin(constant.HTTPConnectManagerFilter)
	if err != nil {
		panic(err)
	}

	hcmc := findHttpManager(lc)
	hcm, err := p.CreateFilter(hcmc, bs)
	if err != nil {
		panic(err)
	}
	return &hcm
}

func findHttpManager(l *model.Listener) *model.HttpConnectionManagerConfig {
	for _, f := range l.FilterChain.Filters {
		if f.Name == constant.HTTPConnectManagerFilter {
			hcmc := &model.HttpConnectionManagerConfig{}
			if err := yaml.ParseConfig(hcmc, f.Config); err != nil {
				return nil
			}

			return hcmc
		}
	}

	panic("http listener filter chain don't contain http connection manager")
}

func (ls *HttpListenerService) GetNetworkFilter() filter.NetworkFilter {
	return ls.NetworkFilter
}

// Start start the listener
func (ls *HttpListenerService) Start() error {
	switch ls.Config.Protocol {
	case model.ProtocolTypeHTTP:
		ls.httpListener()
	case model.ProtocolTypeHTTPS:
		ls.httpsListener()
	default:
		return errors.New(fmt.Sprintf("unsupported protocol start: %d", ls.Config.Protocol))
	}
	return nil
}

func (ls *HttpListenerService) httpsListener() {
	hl := createDefaultHttpWorker(ls)
	hl.pool.New = func() interface{} {
		return ls.allocateContext()
	}
	// user customize http config
	var hc *model.HttpConfig
	hc = model.MapInStruct(ls.Config)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hl.ServeHTTP)
	m := &autocert.Manager{
		Cache:      autocert.DirCache(ls.Config.Address.SocketAddress.CertsDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(ls.Config.Address.SocketAddress.Domains...),
	}
	ls.srv = &http.Server{
		Addr:           ":https",
		Handler:        mux,
		ReadTimeout:    resolveStr2Time(hc.ReadTimeoutStr, 20*time.Second),
		WriteTimeout:   resolveStr2Time(hc.WriteTimeoutStr, 20*time.Second),
		IdleTimeout:    resolveStr2Time(hc.IdleTimeoutStr, 20*time.Second),
		MaxHeaderBytes: resolveInt2IntProp(hc.MaxHeaderBytes, 1<<20),
		TLSConfig:      m.TLSConfig(),
	}
	autoLs := autocert.NewListener(ls.Config.Address.SocketAddress.Domains...)
	logger.Infof("[dubbo-go-server] httpsListener start at : %s", ls.srv.Addr)
	err := ls.srv.Serve(autoLs)
	logger.Info("[dubbo-go-server] httpsListener result:", err)
}

func (ls *HttpListenerService) httpListener() {
	hl := createDefaultHttpWorker(ls)
	hl.pool.New = func() interface{} {
		return ls.allocateContext()
	}

	// user customize http config
	var hc *model.HttpConfig
	hc = model.MapInStruct(ls.Config)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hl.ServeHTTP)

	sa := ls.Config.Address.SocketAddress
	ls.srv = &http.Server{
		Addr:           resolveAddress(sa.Address + ":" + strconv.Itoa(sa.Port)),
		Handler:        mux,
		ReadTimeout:    resolveStr2Time(hc.ReadTimeoutStr, 20*time.Second),
		WriteTimeout:   resolveStr2Time(hc.WriteTimeoutStr, 20*time.Second),
		IdleTimeout:    resolveStr2Time(hc.IdleTimeoutStr, 20*time.Second),
		MaxHeaderBytes: resolveInt2IntProp(hc.MaxHeaderBytes, 1<<20),
	}

	logger.Infof("[dubbo-go-server] httpListener start at : %s", ls.srv.Addr)

	log.Println(ls.srv.ListenAndServe())
}

func (ls *HttpListenerService) allocateContext() *h.HttpContext {
	return &h.HttpContext{
		Listener: ls.Config,
		Params:   make(map[string]interface{}),
	}
}

// createDefaultHttpWorker create http listener
func createDefaultHttpWorker(ls *HttpListenerService) *DefaultHttpWorker {
	return &DefaultHttpWorker{
		pool: sync.Pool{},
		ls:   ls,
	}
}

// ServeHTTP http request entrance.
func (s *DefaultHttpWorker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hc := s.pool.Get().(*h.HttpContext)
	defer s.pool.Put(hc)

	hc.Request = r
	hc.ResetWritermen(w)
	hc.Reset()

	// now only one filter http_connection_manager, so just get it and call
	err := s.ls.NetworkFilter.OnData(hc)

	if err != nil {
		s.pool.Put(hc)
		logger.Errorf("ServeHTTP %v", err)
	}
}

func resolveInt2IntProp(currentV, defaultV int) int {
	if currentV == 0 {
		return defaultV
	}

	return currentV
}

func resolveStr2Time(currentV string, defaultV time.Duration) time.Duration {
	if currentV == "" {
		return defaultV
	} else {
		if duration, err := time.ParseDuration(currentV); err != nil {
			return 20 * time.Second
		} else {
			return duration
		}
	}
}

func resolveAddress(addr string) string {
	if addr == "" {
		logger.Debug("Addr is undefined. Using port :8080 by default")
		return ":8080"
	}

	return addr
}
