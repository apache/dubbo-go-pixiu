package httpproxy

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	http3 "net/http"
	"net/url"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	http2 "github.com/apache/dubbo-go-pixiu/pkg/common/http"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPPROXYFilter
)

// All RemoteFilter instances use one globalClient in order to reuse
// some resounces such as keepalive connections.
var globalClient = &http3.Client{
	// NOTE: Timeout could be no limit, real client or server could cancel it.
	Timeout: 0,
	Transport: &http3.Transport{
		Proxy: http3.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 60 * time.Second,
			DualStack: true,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			// NOTE: Could make it an paramenter,
			// when the requests need cross WAN.
			InsecureSkipVerify: true,
		},
		DisableCompression: false,
		// NOTE: The large number of Idle Connections can
		// reduce overhead of building connections.
		MaxIdleConns:          10240,
		MaxIdleConnsPerHost:   512,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	},
}

func init() {
	extension.RegisterHttpFilter(&RouterPlugin{})
}

type (
	// RouterPlugin is http filter plugin.
	RouterPlugin struct {
	}
	// RouterFilter is http filter instance
	RouterFilter struct {
		hcm *http2.HttpConnectionManager
	}
	// Config describe the config of RouterFilter
	Config struct{}
)

func (rp *RouterPlugin) Kind() string {
	return Kind
}

func (rp *RouterPlugin) CreateFilter(hcm *http2.HttpConnectionManager, config interface{}, bs *model.Bootstrap) (extension.HttpFilter, error) {
	return &RouterFilter{hcm: hcm}, nil
}

func (rf *RouterFilter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(rf.Handle)
	return nil
}

func (rf *RouterFilter) Handle(hc *http.HttpContext) {
	rEntry := hc.GetRouteEntry()
	if rEntry == nil {
		panic("no route entry")
	}

	clusterName := rEntry.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)

	r := hc.Request

	var errPrefix string
	defer func() {
		if err := recover(); err != nil {
			bt, _ := json.Marshal(extension.ErrResponse{Message: fmt.Sprintf("remoteFilterErr: %s: %v", errPrefix, err)})
			hc.SourceResp = bt
			hc.TargetResp = &client.Response{Data: bt}
			hc.WriteJSONWithStatus(http3.StatusServiceUnavailable, bt)
			hc.Abort()
			return
		}
	}()

	var (
		req *http3.Request
		err error
	)

	parsedURL := url.URL{
		Host:     endpoint.Address.GetHost(),
		Scheme:   r.URL.Scheme,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}

	req, err = http3.NewRequest(r.Method, parsedURL.String(), r.Body)
	req.Header = r.Header

	if err != nil {
		bt, _ := json.Marshal(extension.ErrResponse{Message: fmt.Sprintf("BUG: new request failed: %v", err)})
		hc.SourceResp = bt
		hc.TargetResp = &client.Response{Data: bt}
		hc.WriteJSONWithStatus(http3.StatusInternalServerError, bt)
		hc.Abort()
		return
	}

	errPrefix = "do request"
	resp, err := globalClient.Do(req)
	if err != nil {
		panic(err)
	}

	hc.SourceResp = resp
	// response write in response filter.
	hc.Next()
}
