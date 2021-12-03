package grpc

import (
	"context"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	pch "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/pkg/errors"
	"net/http"
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
func (hcm *GrpcConnectionManager) OnData(hc *pch.HttpContext) error {
	hc.Ctx = context.Background()
	err := hcm.findRoute(hc)
	if err != nil {
		return err
	}

	return nil
}

func (hcm *GrpcConnectionManager) findRoute(hc *pch.HttpContext) error {
	ra, err := hcm.routerCoordinator.Route(hc)
	if err != nil {
		if _, err := hc.WriteWithStatus(http.StatusNotFound, constant.Default404Body); err != nil {
			logger.Warnf("WriteWithStatus error %s", err)
		}
		hc.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		e := errors.Errorf("Requested URL %s not found", hc.GetUrl())
		logger.Debug(e.Error())
		return e
		// return 404
	}
	hc.RouteEntry(ra)
	return nil
}
