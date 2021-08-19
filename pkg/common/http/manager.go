package http

import (
	"context"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"net/http"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	pch "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

import (
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
	"github.com/pkg/errors"
)

type HttpConnectionManager struct {
	config            *model.HttpConnectionManager
	routerCoordinator *router2.RouterCoordinator
	clusterManager    *server.ClusterManager
	filters           []extension.HttpFilter
}

func CreateHttpConnectionManager(hcmc *model.HttpConnectionManager, bs *model.Bootstrap) *HttpConnectionManager {
	hcm := &HttpConnectionManager{config: hcmc}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(hcmc)
	hcm.filters = initFilterIfNeed(hcm, hcmc, bs)
	return hcm
}

func (hcm *HttpConnectionManager) OnData(hc *pch.HttpContext) error {
	hc.Ctx = context.Background()
	hcm.findRoute(hc)
	hcm.addFilter(hc)
	hcm.handleHTTPRequest(hc)
	return nil
}

func (hcm *HttpConnectionManager) findRoute(hc *pch.HttpContext) {
	ra, err := hcm.routerCoordinator.Route(hc)
	if err != nil {
		// return 404
	}
	hc.SetRouteEntry(ra)
}

func (hcm *HttpConnectionManager) handleHTTPRequest(c *pch.HttpContext) {
	if len(c.BaseContext.Filters) > 0 {
		c.Next()
		c.WriteHeaderNow()
		return
	}

	// TODO redirect
}

func (hcm *HttpConnectionManager) addFilter(ctx *pch.HttpContext) {
	for _, f := range hcm.filters {
		f.PrepareFilterChain(ctx)
	}
}

func (hcm *HttpConnectionManager) routeRequest(ctx *pch.HttpContext, req *http.Request) (router.API, error) {
	apiDiscSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	api, err := apiDiscSrv.GetAPI(req.URL.Path, fc.HTTPVerb(req.Method))
	if err != nil {
		ctx.WriteWithStatus(http.StatusNotFound, constant.Default404Body)
		ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		e := errors.Errorf("Requested URL %s not found", req.URL.Path)
		logger.Debug(e.Error())
		return router.API{}, e
	}
	if !api.Method.OnAir {
		ctx.WriteWithStatus(http.StatusNotAcceptable, constant.Default406Body)
		ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		e := errors.Errorf("Requested API %s %s does not online", req.Method, req.URL.Path)
		logger.Debug(e.Error())
		return router.API{}, e
	}
	ctx.API(api)
	return api, nil
}

func initFilterIfNeed(hcm *HttpConnectionManager, hcmc *model.HttpConnectionManager, bs *model.Bootstrap) []extension.HttpFilter {
	var hfs []extension.HttpFilter

	for _, f := range hcmc.HTTPFilters {
		hp, err := extension.GetHttpFilterPlugin(f.Name)
		if err != nil {
			logger.Error("initFilterIfNeed get plugin error %s", err)
		}

		hf, err := hp.CreateFilter(hcm, f.Config, bs)
		if err != nil {
			logger.Error("initFilterIfNeed create filter error %s", err)
		}
		hfs = append(hfs, hf)
	}
	return hfs
}
