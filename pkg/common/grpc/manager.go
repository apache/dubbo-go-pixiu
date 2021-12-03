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
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"net/http"
)

// HttpConnectionManager network filter for http
type GrpcConnectionManager struct {
	config            *model.GRPCConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
	filterManager     *filter.FilterManager

	serverClientConn map[string][ggrpc.ClientConn]
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

func (gcm *GrpcConnectionManager) handler(srv interface{}, ss ggrpc.ServerStream) error {

	fullMethodName, ok := ggrpc.MethodFromServerStream(ss)
	if !ok {
		return ggrpc.Errorf(codes.Internal, "handler not found method name")
	}

	err := gcm.findRoute()
	if err != nil {
		return err
	}

	ctx := ss.Context()
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		logger.Infof("GrpcConnectionManager handler grpc FromIncomingContext not ok")
	}

	outCtx, _ := context.WithCancel(ctx)
	outCtx = metadata.NewOutgoingContext(outCtx, md.Copy())



	return nil
}

func (gcm *GrpcConnectionManager) findRoute(hc *pch.HttpContext) error {
	ra, err := gcm.routerCoordinator.Route(hc)
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
