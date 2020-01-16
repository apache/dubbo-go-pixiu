package service

import (
	"context"
	"github.com/pantianying/dubbo-go-proxy/dubbo"
)

type ProxyContext interface {
	//base
	NextFilter() Filter

	Match() (*MetadataIdentifier, *MetadataInfo, int)

	InvokeData() *dubbo.InvokeData

	//SetInvokeResp()
	//GetInvokeResp()
}
type BaseContext struct {
	filter []Filter
	ctx    context.Context
}

func NewBaseContext(filter []Filter) *BaseContext {
	return &BaseContext{
		filter: filter,
		ctx:    context.Background(),
	}
}
func (h *BaseContext) NextFilter() Filter {
	if len(h.filter) > 0 {
		f := h.filter[0]
		h.filter = h.filter[1:]
		return f
	}
	return nil
}
