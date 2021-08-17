package routerfilter

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	http2 "github.com/apache/dubbo-go-pixiu/pkg/common/http"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"net"
	http3 "net/http"
	"time"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPRouterFilter
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

func (rf *RouterFilter) Handle(c *http.HttpContext) {
	rEntry := c.GetRouteEntry()
	if rEntry == nil {
		panic("no route entry")
	}

	clusterName := rEntry.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)

	r, w := c.Request, c.SourceResp

	var errPrefix string
	defer func() {
		if err := recover(); err != nil {

			w.SetStatusCode(http3.StatusServiceUnavailable)
			// NOTE: We don't use stringtool.Cat because err needs
			// the internal reflection of fmt.Sprintf.
			ctx.AddTag(fmt.Sprintf("remoteFilterErr: %s: %v", errPrefix, err))
			result = resultFailed
		}
	}()

	var (
		req *http.Request
		err error
	)

	if rf.spec.timeout > 0 {
		timeoutCtx, cancelFunc := stdcontext.WithTimeout(stdcontext.Background(), rf.spec.timeout)
		defer cancelFunc()
		req, err = http.NewRequestWithContext(timeoutCtx, http.MethodPost, rf.spec.URL, bytes.NewReader(ctxBuff))
	} else {
		req, err = http.NewRequest(http.MethodPost, rf.spec.URL, bytes.NewReader(ctxBuff))
	}

	if err != nil {
		logger.Errorf("BUG: new request failed: %v", err)
		w.SetStatusCode(http.StatusInternalServerError)
		ctx.AddTag(stringtool.Cat("remoteFilterBug: ", err.Error()))
		return resultFailed
	}

	errPrefix = "do request"
	resp, err := globalClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		panic(fmt.Errorf("not 2xx status code: %d", resp.StatusCode))
	}

	errPrefix = "read remote body"
	ctxBuff = rf.limitRead(resp.Body, maxContextBytes)

	errPrefix = "unmarshal context"
	rf.unmarshalHTTPContext(ctxBuff, ctx)

	if resp.StatusCode == 205 {
		return resultResponseAlready
	}

	return ""

}
