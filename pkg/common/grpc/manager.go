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

package grpc

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"io/ioutil"
	"net"
	stdHttp "net/http"
)

import (
	"github.com/pkg/errors"
	"golang.org/x/net/http2"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

// HttpConnectionManager network filter for http
type GrpcConnectionManager struct {
	config            *model.GRPCConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
}

// CreateHttpConnectionManager create http connection manager
func CreateGrpcConnectionManager(hcmc *model.GRPCConnectionManagerConfig, bs *model.Bootstrap) *GrpcConnectionManager {
	hcm := &GrpcConnectionManager{config: hcmc}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(&hcmc.RouteConfig)
	return hcm
}

func (gcm *GrpcConnectionManager) OnData(data []byte) (interface{}, int, error) {
	return nil, 0, nil
}

func (gcm *GrpcConnectionManager) ServeHTTP(w stdHttp.ResponseWriter, r *stdHttp.Request) {

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
		bt, _ := json.Marshal(http.ErrResponse{Message: "pixiu forward error"})
		w.WriteHeader(stdHttp.StatusServiceUnavailable)
		w.Write(bt)
		return
	}

	if err := gcm.response(w, res); err != nil {
		logger.Info("GrpcConnectionManager response  error %v", err)
	}
}

func (gcm *GrpcConnectionManager) response(w stdHttp.ResponseWriter, res *stdHttp.Response) error {
	defer res.Body.Close()

	copyHeader(w.Header(), res.Header)
	w.WriteHeader(res.StatusCode)

	byts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	w.Write(byts)

	for k, vv := range res.Trailer {
		k = stdHttp.TrailerPrefix + k
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}
	return nil
}

func (gcm *GrpcConnectionManager) newHttpForwarder() *HttpForwarder {
	transport := &http2.Transport{
		DialTLS: func(network, addr string, _ *tls.Config) (net.Conn, error) {
			return net.Dial(network, addr)
		},
		AllowHTTP: true,
	}
	return &HttpForwarder{transport: transport}
}

func copyHeader(dst, src stdHttp.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
