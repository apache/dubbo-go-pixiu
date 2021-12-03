package grpc

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	pixiu_grpc "github.com/apache/dubbo-go-pixiu/pkg/common/grpc"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"google.golang.org/grpc"
	"net"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeHTTP, newGrpcListenerService)
}

type (
	// ListenerService the facade of a listener
	GrpcListenerService struct {
		Config  *model.Listener
		manager *pixiu_grpc.GrpcConnectionManager
		server  *grpc.Server
	}
)

func newGrpcListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	hcm := createGrpcManager(lc, bs)

	return &GrpcListenerService{
		Config:  lc,
		manager: hcm,
	}, nil
}

func createGrpcManager(lc *model.Listener, bs *model.Bootstrap) *pixiu_grpc.GrpcConnectionManager {
	p, err := filter.GetNetworkFilterPlugin(constant.GRPCConnectManagerFilter)
	if err != nil {
		panic(err)
	}

	hcmc := findGrpcManagerConfig(lc)
	hcm, err := p.CreateFilter(hcmc, bs)
	if err != nil {
		panic(err)
	}
	return hcm.(pixiu_grpc.GrpcConnectionManager)
}

func findGrpcManagerConfig(l *model.Listener) *model.GRPCConnectionManagerConfig {
	for _, fc := range l.FilterChains {
		for _, f := range fc.Filters {
			if f.Name == constant.GRPCConnectManagerFilter {
				hcmc := &model.GRPCConnectionManagerConfig{}
				if err := yaml.ParseConfig(hcmc, f.Config); err != nil {
					return nil
				}

				return hcmc
			}
		}
	}

	panic("http listener filter chain don't contain http connection manager")
}

func (g GrpcListenerService) GetNetworkFilter() filter.NetworkFilter {
	panic("implement me")
}

func (g GrpcListenerService) Start() error {
	ln, err := net.Listen("tcp", g.Config.Address.SocketAddress.Address)
	if err != nil {
		logger.Errorf("GrpcListenerService Start %s", err)
	}
	logger.Infof("start grpc server")
	p := NewProxy()
	g.server = NewGrpcServerWithProxy(p)
	RegisterProxy(g.server, p, g)
	return g.server.Serve(ln)
}

func (g *GrpcListenerService) handler(srv interface{}, ss grpc.ServerStream) error {
	g.manager
}

func NewGrpcServerWithProxy(g GrpcListenerService) *grpc.Server {
	return grpc.NewServer(
		grpc.UnknownServiceHandler(g.handler),
	)
}

func RegisterProxy(server *grpc.Server, g GrpcListenerService, listener GrpcListenerService) {

	serviceName := "dddd"

	proxyDesc := &grpc.ServiceDesc{
		ServiceName: serviceName,
		HandlerType: (*interface{})(nil),
	}

	methodNames := []string{"sssss"}

	for _, m := range methodNames {
		sd := grpc.StreamDesc{
			StreamName:    m,
			Handler:       g.handler,
			ServerStreams: true,
			ClientStreams: true,
		}
		proxyDesc.Streams = append(proxyDesc.Streams, sd)
	}
	server.RegisterService(proxyDesc, g)
}
