package grpc

import (
	"context"
	"crypto/tls"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	pch "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
	"io/ioutil"
	stdHttp "net/http"
	"strings"
	"net"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"encoding/json"
	"time"
)

// HttpConnectionManager network filter for http
type GrpcConnectionManager struct {
	config            *model.GRPCConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
	filterManager     *filter.FilterManager
}

// CreateHttpConnectionManager create http connection manager
func CreateGrpcConnectionManager(hcmc *model.GRPCConnectionManagerConfig, bs *model.Bootstrap) *GrpcConnectionManager {
	hcm := &GrpcConnectionManager{config: hcmc}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(&hcmc.RouteConfig)
	return hcm
}

// OnData receive data from listener
func (gcm *GrpcConnectionManager) OnData(hc *pch.HttpContext) error {
	panic("grpc connection manager OnData function shouldn't be called")
}

func (gcm *GrpcConnectionManager) ServeHTTP(w stdHttp.ResponseWriter, r *stdHttp.Request) {
	logger.Info("GrpcConnectionManager ServeHTTP")

	ra, err := gcm.routerCoordinator.RouteByPathAndName(r.RequestURI, r.Method)
	if err != nil {
		w.WriteHeader(stdHttp.StatusNotFound)
		if _, err := w.Write(constant.Default404Body); err != nil {
			logger.Warnf("WriteWithStatus error %v", err)
		}
		w.Header().Add(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		e := errors.Errorf("Requested URL %s not found", r.RequestURI)
		logger.Debug(e.Error())
		// return 404
	}

	logger.Debugf("[dubbo-go-pixiu] client choose endpoint from cluster :%v", ra.Cluster)

	clusterName := ra.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		bt, _ := json.Marshal(http.ErrResponse{Message: "cluster not found endpoint"})
		w.WriteHeader(stdHttp.StatusServiceUnavailable)
		w.Write(bt)
		return
	}

	outreq := new(stdHttp.Request)
	*outreq = *r
	outreq.URL.Scheme = "http"
	outreq.URL.Host = endpoint.Address.GetAddress()
	r2 := outreq.WithContext(context.Background())
	*outreq = *r2

	// todo: need cache?
	forwarder := gcm.newHttpForwarder()
	res, err := forwarder.Forward(outreq)

	if err != nil {
		logger.Info("GrpcConnectionManager forward request error %v", err)
	}

	inres := new(stdHttp.Response)
	*inres = *res
	inres.StatusCode = res.StatusCode
	inres.ContentLength = res.ContentLength

	if err := gcm.response(w, inres); err != nil {

	}
}

func (gmc *GrpcConnectionManager) response(w stdHttp.ResponseWriter, res *stdHttp.Response) error{
	copyHeader(w.Header(), res.Header)

	announcedTrailers := len(res.Trailer)
	if announcedTrailers > 0 {
		trailerKeys := make([]string, 0, len(res.Trailer))
		for k := range res.Trailer {
			trailerKeys = append(trailerKeys, k)
		}
		w.Header().Add("Trailer", strings.Join(trailerKeys, ", "))
	}

	w.WriteHeader(res.StatusCode)
	byts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	w.Write(byts)
	return nil
}

func (gcm *GrpcConnectionManager) newHttpForwarder() * HttpForwarder{
	transport := &http2.Transport{
		DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
		AllowHTTP: true,
	}
	return &HttpForwarder{ transport: transport}
}

func copyHeader(dst, src stdHttp.Header) {

	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}