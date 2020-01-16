package filter

import (
	"github.com/dubbogo/dubbo-go-proxy/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/service"
)

func init() {
	service.SetFilter(constant.MatchFilterName, NewMatchFilter)
}

type MatchFilter struct{}

func NewMatchFilter() service.Filter {
	return &MatchFilter{}
}
func (f *MatchFilter) OnRequest(ctx service.ProxyContext) (ret int) {
	_, _, ret = ctx.Match()
	return
}
