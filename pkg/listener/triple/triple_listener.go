package triple

import (
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	tripleCommon "github.com/dubbogo/triple/pkg/common"
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
	triConfig "github.com/dubbogo/triple/pkg/config"
	"github.com/dubbogo/triple/pkg/triple"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeTriple, newTripleListenerService)
}

type (
	// ListenerService the facade of a listener
	TripleListenerService struct {
		listener.BaseListenerService
		server triple.TripleServer
	}
)

func newTripleListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {

	opts := []triConfig.OptionFunction{
		triConfig.WithCodecType(tripleConstant.HessianCodecName),
		triConfig.WithLocation(lc.Address.SocketAddress.GetAddress()),
		triConfig.WithLogger(logger.GetLogger()),
	}

	triOption := triConfig.NewTripleOption(opts...)
	tripleService := &dubbo3.UnaryService{proxyImpl: invoker}

	fc := filterchain.CreateFilterChain(lc.FilterChain, bs)
	return &TripleListenerService{
		BaseListenerService: listener.BaseListenerService{
			Config:      lc,
			FilterChain: fc,
		},
	}, nil
}

func (ls *TripleListenerService) Start() error {

}
